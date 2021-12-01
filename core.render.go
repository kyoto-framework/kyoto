package kyoto

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"reflect"
	"sync"
	"time"
)

// Component Store Lifecycle is a temporary storage for components processing
// Will be cleared in the end of lifecycle
var csl = map[Page][]Component{}
var cslrw = &sync.RWMutex{}

// RegisterComponent is used while defining components in the Init() section
func RegisterComponent(p Page, c Component) Component {
	// Extract insights
	insights := GetInsights(p)
	// Init csl store
	cslrw.Lock()
	if _, ok := csl[p]; !ok {
		csl[p] = []Component{}
	}
	cslrw.Unlock()
	// Save type to SSA store
	ssacstorerw.Lock()
	if _, ok := ssacstore[reflect.ValueOf(c).Elem().Type().Name()]; !ok {
		// Extract component type
		ctype := reflect.ValueOf(c).Elem().Type()
		// Save to store
		ssacstore[reflect.ValueOf(c).Elem().Type().Name()] = ctype
	}
	ssacstorerw.Unlock()
	// Save component to lifecycle store
	cslrw.Lock()
	csl[p] = append(csl[p], c)
	cslrw.Unlock()
	// Trigger component init
	if _c, ok := c.(ImplementsInit); ok {
		st := time.Now()
		_c.Init(p)
		if insights != nil {
			GetInsights(p).GetOrCreateNested(_c).Update(InsightsTiming{
				Init: time.Since(st),
			})
		}
	} else if _c, ok := c.(ImplementsInitWithoutPage); ok {
		st := time.Now()
		_c.Init()
		if insights != nil {
			GetInsights(p).GetOrCreateNested(_c).Update(InsightsTiming{
				Init: time.Since(st),
			})
		}
	}
	// Return component for external assignment
	return c
}

// RegC is an alias for RegisterComponent
var RegC = RegisterComponent

// RenderPage is a main entrypoint of rendering. Responsible for rendering and components lifecycle
func RenderPage(w io.Writer, p Page) {
	// Init insights (if enabled)
	var insights *Insights
	if INSIGHTS {
		insights = NewInsights(p)
	}
	// Async specific state
	var wg sync.WaitGroup
	var err = make(chan error, 1000)
	// Trigger init
	if p, ok := p.(ImplementsInitWithoutPage); ok {
		st := time.Now()
		p.Init()
		if insights != nil {
			insights.Update(InsightsTiming{
				Init: time.Since(st),
			})
		}
	}
	// Check redirect
	if redirected := GetContext(p, "internal:redirect"); redirected != nil {
		return
	}
	// Trigger async in goroutines
	st := time.Now()
	subset := 0
	for {
		cslrw.RLock()
		regc := csl[p][subset:]
		cslrw.RUnlock()
		subset += len(regc)
		if len(regc) == 0 {
			break
		}
		for _, component := range regc {
			var cinsights *Insights
			if insights != nil {
				cinsights = insights.GetOrCreateNested(component)
			}
			if _component, ok := component.(ImplementsAsync); ok {
				wg.Add(1)
				go func(wg *sync.WaitGroup, err chan error, c ImplementsAsync, p Page) {
					defer wg.Done()
					st := time.Now()
					_err := c.Async(p)
					if cinsights != nil {
						cinsights.Update(InsightsTiming{
							Async: time.Since(st),
						})
					}
					if _err != nil {
						err <- _err
					}
				}(&wg, err, _component, p)
			} else if _component, ok := component.(ImplementsAsyncWithoutPage); ok {
				wg.Add(1)
				go func(wg *sync.WaitGroup, err chan error, c ImplementsAsyncWithoutPage) {
					defer wg.Done()
					st := time.Now()
					_err := c.Async()
					if cinsights != nil {
						cinsights.Update(InsightsTiming{
							Async: time.Since(st),
						})
					}
					if _err != nil {
						err <- _err
					}
				}(&wg, err, _component)
			}
		}
		wg.Wait()
	}
	if insights != nil {
		insights.Update(InsightsTiming{
			Async: time.Since(st),
		})
	}
	// Trigger aftersync
	st = time.Now()
	cslrw.RLock()
	for _, component := range csl[p] {
		var cinsights *Insights
		if insights != nil {
			cinsights = insights.GetOrCreateNested(component)
		}
		if _component, ok := component.(ImplementsAfterAsync); ok {
			st := time.Now()
			_component.AfterAsync(p)
			if cinsights != nil {
				cinsights.Update(InsightsTiming{
					AfterAsync: time.Since(st),
				})
			}
		} else if _component, ok := component.(ImplementsAfterAsyncWithoutPage); ok {
			st := time.Now()
			_component.AfterAsync()
			if cinsights != nil {
				cinsights.Update(InsightsTiming{
					AfterAsync: time.Since(st),
				})
			}
		}
	}
	cslrw.RUnlock()
	if insights != nil {
		insights.Update(InsightsTiming{
			AfterAsync: time.Since(st),
		})
	}
	// Clear components store (not needed more)
	cslrw.Lock()
	delete(csl, p)
	cslrw.Unlock()
	// Extract flags
	redirected := GetContext(p, "internal:redirect")
	// Execute template
	if redirected == nil && len(err) == 0 {
		st := time.Now()
		if _p, ok := p.(ImplementsTemplateWithoutPage); ok {
			err := _p.Template().Execute(w, _p)
			if err != nil {
				panic(err)
			}
		} else if _p, ok := p.(ImplementsRender); ok {
			w.Write([]byte(_p.Render()))
		}
		if insights != nil {
			insights.Update(InsightsTiming{
				Render: time.Since(st),
			})
		}
	}
	// Print insights
	if INSIGHTS && INSIGHTS_CLI {
		if INSIGHTS_CLI_JSON {
			jsonInsights, err := json.Marshal(insights)

			if err != nil {
				panic(err)
			}

			log.Printf(" ---------------- insights %s %s", insights.ID, insights.Name)
			log.Printf(string(jsonInsights))
		} else {
			log.Printf(" ---------------- insights %s %s", insights.ID, insights.Name)
			log.Printf("i:%s a:%s aa:%s r:%s", insights.Init, insights.Async, insights.AfterAsync, insights.Render)
			for _, ci := range insights.Nested {
				log.Printf("--- id:%s n:%s i:%s a:%s aa:%s r:%s", ci.ID, ci.Name, ci.Init, ci.Async, ci.AfterAsync, ci.Render)
			}
		}
	}
}

// Redirect is a wrapper around http.Redirect for correct work inside of SSC
func Redirect(rp *RedirectParameters) {
	// Write redirect value
	SetContext(rp.Page, "internal:redirect", rp.Target)
	// Extract r/rw
	rw := rp.ResponseWriter
	r := rp.Request
	if rp.ResponseWriterKey != "" {
		rw = GetContext(rp.Page, rp.ResponseWriterKey).(http.ResponseWriter)
	}
	if rp.RequestKey != "" {
		r = GetContext(rp.Page, rp.RequestKey).(*http.Request)
	}
	// Do actual redirect in case of usual response
	if _, ssa := rp.Page.(*DummyPage); !ssa {
		http.Redirect(rw, r, rp.Target, rp.StatusCode)
	}
}

// PageHandler is an opinionated page net/http handler.
// Context:
// - internal:rw - http.ResponseWriter
// - internal:r - *http.Request
func PageHandler(p Page) http.HandlerFunc {
	// Return handler
	return func(rw http.ResponseWriter, r *http.Request) {
		// Init new page instance
		_p := reflect.New(reflect.TypeOf(p).Elem()).Interface().(Page)
		// Set context
		SetContext(_p, "internal:rw", rw)
		SetContext(_p, "internal:r", r)
		// Render page
		RenderPage(rw, _p)
		// Clear context
		DelContext(_p, "")
	}
}
