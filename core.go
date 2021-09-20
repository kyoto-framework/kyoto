package ssc

import (
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
var cslLock = &sync.RWMutex{}

// RegisterComponent is used while defining components in the Init() section
func RegisterComponent(p Page, c Component) Component {
	// Init csl store
	cslLock.Lock()
	if _, ok := csl[p]; !ok {
		csl[p] = []Component{}
	}
	cslLock.Unlock()
	// Save type to SSA store
	csSSALock.Lock()
	if _, ok := csSSA[reflect.ValueOf(c).Elem().Type().Name()]; !ok {
		// Extract component type
		ctype := reflect.ValueOf(c).Elem().Type()
		// Save to store
		csSSA[reflect.ValueOf(c).Elem().Type().Name()] = ctype
	}
	csSSALock.Unlock()
	// Save component to lifecycle store
	cslLock.Lock()
	csl[p] = append(csl[p], c)
	cslLock.Unlock()
	// Trigger component init
	if c, ok := c.(ImplementsNestedInit); ok {
		c.Init(p)
	}
	// Return component for external assignment
	return c
}

// Alias for RegisterComponent
var RegC = RegisterComponent

// RenderPage is a main entrypoint of rendering. Responsible for rendering and components lifecycle
func RenderPage(w io.Writer, p Page) {
	// Async specific state
	var wg sync.WaitGroup
	var err = make(chan error, 1000)
	// Trigger init
	if p, ok := p.(ImplementsInit); ok {
		st := time.Now()
		p.Init()
		et := time.Since(st)
		if BENCH_LOWLEVEL {
			log.Println("Init time", reflect.TypeOf(p), et)
		}
	}
	// Trigger async in goroutines
	subset := 0
	for {
		cslLock.RLock()
		regc := csl[p][subset:]
		cslLock.RUnlock()
		subset += len(regc)
		if len(regc) == 0 {
			break
		}
		for _, component := range regc {
			if component, ok := component.(ImplementsAsync); ok {
				wg.Add(1)
				go func(wg *sync.WaitGroup, err chan error, c ImplementsAsync) {
					defer wg.Done()
					st := time.Now()
					_err := c.Async()
					et := time.Since(st)
					if BENCH_LOWLEVEL {
						log.Println("Async time", reflect.TypeOf(component), et)
					}
					if _err != nil {
						err <- _err
					}
				}(&wg, err, component)
			}
		}
		wg.Wait()
	}
	// Trigger aftersync
	cslLock.RLock()
	for _, component := range csl[p] {
		if component, ok := component.(ImplementsAfterAsync); ok {
			st := time.Now()
			component.AfterAsync()
			et := time.Since(st)
			if BENCH_LOWLEVEL {
				log.Println("AfterAsync time", reflect.TypeOf(component), et)
			}
		}
	}
	cslLock.RUnlock()
	// Clear components store (not needed more)
	cslLock.Lock()
	delete(csl, p)
	cslLock.Unlock()
	// Extract flags
	redirected := GetContext(p, "internal:redirected")
	// Execute template
	if redirected == nil {
		st := time.Now()
		terr := p.Template().Execute(w, reflect.ValueOf(p).Elem())
		et := time.Since(st)
		if BENCH_LOWLEVEL {
			log.Println("Execution time", et)
		}
		if terr != nil {
			panic(terr)
		}
	}
}

// Redirect is a wrapper around http.Redirect for correct work inside of SSC
func Redirect(rp *RedirectParameters) {
	// Write redirected flag
	SetContext(rp.Page, "internal:redirected", true)
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
	if _, ssa := rp.Page.(*dummypage); !ssa {
		http.Redirect(rw, r, rp.Target, rp.StatusCode)
	} else { // Special header in case of SSA
		rw.Header().Add("X-Redirect", rp.Target)
	}
}

// PageHandlerFactory is a factory for building Page handler.
// Simple wrapper around RenderPage with context setting.
// Good for defining own project-level handler.
// Example of usage:
// func handle(p ssc.Page) http.HandlerFunc {
//     return func(rw http.ResponseWriter, r *http.Request) {
// 	       ssc.PageHandlerFactory(p, map[string]interface{}{
//	           "internal:rw": rw,
//             "internal:r": r,
//         })(rw, r)
//     }
// }
func PageHandlerFactory(p Page, context map[string]interface{}) http.HandlerFunc {
	// Make page instance
	var pi Page
	pptr := reflect.New(reflect.TypeOf(p).Elem())
	pi = pptr.Interface().(Page)
	// Set context
	for k, v := range context {
		SetContext(pi, k, v)
	}
	// Return handler
	return func(rw http.ResponseWriter, r *http.Request) {
		// Render page
		RenderPage(rw, pi)
		// Clear context
		DelContext(pi, "")
	}
}
