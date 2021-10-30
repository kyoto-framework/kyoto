package kyoto

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"reflect"
	"strings"
	"sync"
)

// SSA Component Store is a storage for component types.
// When SSA is called, page's general lifecycle components store is not available (we have dummy page instead).
var ssacstore = map[string]reflect.Type{}
var ssacstorerw = &sync.RWMutex{}

// SSAParameters represents parameters, needed for handling SSA request
type SSAParameters struct {
	Component string
	Action    string
	State     string // In encoded state
	Args      []interface{}
}

// RenderSSA is a low-level component rendering function for SSA. Responsible for rendering and components SSA lifecycle
func RenderSSA(w io.Writer, dp *DummyPage, p SSAParameters) {
	// Async specific state
	var wg sync.WaitGroup
	var err = make(chan error, 1000)
	// Find component type in store
	ssacstorerw.RLock()
	ctype, found := ssacstore[p.Component]
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
	state, _ := base64.StdEncoding.DecodeString(p.State)
	if err := json.Unmarshal(state, &component); err != nil {
		panic(err)
	}
	// Call action
	if _component, ok := component.(ImplementsActions); ok {
		_component.Actions(dp)[p.Action](p.Args...)
	} else if _component, ok := component.(ImplementsActionsWithoutPage); ok {
		_component.Actions()[p.Action](p.Args...)
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
		t = template.Must(t.Parse(fmt.Sprintf(`{{ template "%s" . }}`, p.Component)))
		// Render component
		terr := t.Execute(w, component)
		if terr != nil {
			panic(terr)
		}
	}
}

// SSAHandler is an opinionated SSA net/http handler.
// Context:
// - internal:rw - http.ResponseWriter
// - internal:r - *http.Request
func SSAHandler(tb TemplateBuilder) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Init dummy page
		dp := &DummyPage{
			TemplateBuilder: tb,
		}
		// Extract SSA parameters
		params := SSAParameters{}
		tokens := strings.Split(r.URL.Path, "/")
		var args []interface{}
		json.Unmarshal([]byte(r.PostFormValue("Args")), &args)
		params.Component = tokens[2]
		params.Action = tokens[3]
		params.State = r.PostFormValue("State")
		params.Args = args
		// Set context
		SetContext(dp, "internal:rw", rw)
		SetContext(dp, "internal:r", r)
		// Render SSA
		RenderSSA(rw, dp, params)
		// Del context
		DelContext(dp, "")
	}
}

// Deprecated: define own handler using low-level RenderSSA, or use SSAHandler instead
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
	dp := &DummyPage{
		TemplateBuilder: tb,
	}
	// Set context
	for k, v := range context {
		SetContext(dp, k, v)
	}
	// Return handler
	return func(rw http.ResponseWriter, r *http.Request) {
		// Extract SSA parameters
		params := SSAParameters{}
		tokens := strings.Split(r.URL.Path, "/")
		var args []interface{}
		json.Unmarshal([]byte(r.PostFormValue("Args")), &args)
		params.Component = tokens[2]
		params.Action = tokens[3]
		params.State = r.PostFormValue("State")
		params.Args = args
		// Set context
		SetContext(dp, "internal:rw", rw)
		SetContext(dp, "internal:r", r)
		// Render SSA
		RenderSSA(rw, dp, params)
		// Del context
		DelContext(dp, "")
	}
}
