package kyoto

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
)

// SSA Component Store is a storage for component types.
// When SSA is called, page's general lifecycle components store is not available (we have dummy page instead).
var ssacstore = map[string]reflect.Type{}
var ssacstorerw = &sync.RWMutex{}

// SSAHandlerFactory is a factory for building Server Side Action handler.
// Check documentation for lifecycle details (different comparing to page's).
// Example of usage:
// func ssatemplate(p ssc.Page) *template.Template {
// 	return template.Must(template.New("SSA").Funcs(ssc.Funcs()).ParseGlob("*.html"))
// }
// func ssahandler() http.HandlerFunc {
//     return func(rw http.ResponseWriter, r *http.Request) {
// 	       ssc.SSAHandlerFactory(ssatemplate, map[string]interface{}{
//	           "internal:rw": rw,
//             "internal:r": r,
//         })(rw, r)
//     }
// }
func SSAHandlerFactory(tb TemplateBuilder, context map[string]interface{}) http.HandlerFunc {
	// Init dummy page
	dp := &dummypage{
		TemplateBuilder: tb,
	}
	// Set context
	for k, v := range context {
		SetContext(dp, k, v)
	}
	// Return handler
	return func(rw http.ResponseWriter, r *http.Request) {
		// Async specific state
		var wg sync.WaitGroup
		var err = make(chan error, 1000)
		// Extract component action and name from route
		tokens := strings.Split(r.URL.Path, "/")
		cname := tokens[2]
		aname := tokens[3]
		// Find component type in store
		ssacstorerw.RLock()
		ctype, found := ssacstore[cname]
		ssacstorerw.RUnlock()
		// Panic, if not found
		if !found {
			panic("Can't find component. Seems like it's not registered")
		}
		// Create component
		component := reflect.New(ctype).Interface().(Component)
		// Init
		if _component, ok := component.(ImplementsInit); ok {
			_component.Init(dp)
		} else if _component, ok := component.(ImplementsInitWithoutPage); ok {
			_component.Init()
		}
		// Populate component state
		state, _ := url.QueryUnescape(r.PostFormValue("State"))
		if err := json.Unmarshal([]byte(state), &component); err != nil {
			panic(err)
		}
		// Extract arguments
		var args []interface{}
		json.Unmarshal([]byte(r.PostFormValue("Args")), &args)
		// Call action
		if _component, ok := component.(ImplementsActions); ok {
			_component.Actions(dp)[aname](args...)
		} else if _component, ok := component.(ImplementsActionsWithoutPage); ok {
			_component.Actions()[aname](args...)
		} else {
			panic("Component not implements Actions, unexpected behavior")
		}
		// If new components registered, trigger async
		subset := 0
		for {
			cslrw.RLock()
			regc := csl[dp][subset:]
			cslrw.RUnlock()
			subset += len(regc)
			if len(regc) == 0 {
				break
			}
			for _, component := range regc {
				if _component, ok := component.(ImplementsAsync); ok {
					wg.Add(1)
					go func(wg *sync.WaitGroup, err chan error, c ImplementsAsync, dp Page) {
						defer wg.Done()
						_err := c.Async(dp)
						if _err != nil {
							err <- _err
						}
					}(&wg, err, _component, dp)
				} else if _component, ok := component.(ImplementsAsyncWithoutPage); ok {
					wg.Add(1)
					go func(wg *sync.WaitGroup, err chan error, c ImplementsAsyncWithoutPage) {
						defer wg.Done()
						_err := c.Async()
						if _err != nil {
							err <- _err
						}
					}(&wg, err, _component)
				}
			}
			wg.Wait()
		}
		// Extact flags
		redirected := GetContext(dp, "internal:redirected")
		// Render page
		if redirected == nil {
			// Prepare template
			t := dp.Template()
			t = template.Must(t.Parse(fmt.Sprintf(`{{ template "%s" . }}`, cname)))
			// Render component
			terr := t.Execute(rw, component)
			if terr != nil {
				panic(terr)
			}
		}
		// Clear context
		DelContext(dp, "")
	}
}
