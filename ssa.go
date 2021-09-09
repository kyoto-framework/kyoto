package ssc

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"time"
)

// Component Store SSA is a storage for component types.
// When SSA is called, page's general lifecycle components store is not available (we have dummy page instead).
var csSSA = map[string]reflect.Type{}
var csSSALock = &sync.Mutex{}

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
		// Extract component action and name from route
		tokens := strings.Split(r.URL.Path, "/")
		cname := tokens[2]
		aname := tokens[3]
		// Find component type in store
		ctype, found := csSSA[cname]
		// Panic, if not found
		if !found {
			panic("Can't find component. Seems like it's not registered")
		}
		// Create component
		component := reflect.New(ctype).Interface().(Component)
		// Init
		if component, ok := component.(ImplementsNestedInit); ok {
			st := time.Now()
			component.Init(dp)
			et := time.Since(st)
			if BENCH_HANDLERS {
				log.Println("Init time", reflect.TypeOf(component), et)
			}
		}
		// Populate component state
		st := time.Now()
		state, _ := url.QueryUnescape(r.PostFormValue("State"))
		if err := json.Unmarshal([]byte(state), &component); err != nil {
			panic(err)
		}
		et := time.Since(st)
		if BENCH_HANDLERS {
			log.Println("Populate time", reflect.TypeOf(component), et)
		}
		// Extract arguments
		st = time.Now()
		var args []interface{}
		json.Unmarshal([]byte(r.PostFormValue("Args")), &args)
		et = time.Since(st)
		if BENCH_HANDLERS {
			log.Println("Extract args time", reflect.TypeOf(component), et)
		}
		// Call action
		st = time.Now()
		if component, ok := component.(ImplementsActions); ok {
			component.Actions()[aname](args...)
		} else {
			panic("Component not implements Actions, unexpected behavior")
		}
		et = time.Since(st)
		if BENCH_HANDLERS {
			log.Println("Action time", reflect.TypeOf(component), et)
		}
		// Extact flags
		redirected := GetContext(dp, "internal:redirected")
		// Render page
		if redirected == nil {
			// Prepare template
			st = time.Now()
			t := dp.Template()
			t = template.Must(t.Parse(fmt.Sprintf(`{{ template "%s" . }}`, cname)))
			et = time.Since(st)
			if BENCH_HANDLERS {
				log.Println("Template prepare time", reflect.TypeOf(component), et)
			}
			// Render component
			st = time.Now()
			terr := t.Execute(rw, component)
			if terr != nil {
				panic(terr)
			}
			et = time.Since(st)
			if BENCH_HANDLERS {
				log.Println("Executiton time", reflect.TypeOf(component), et)
			}
		}
		// Clear context
		DelContext(dp, "")
	}
}

// SSA client side code
var ssaclient = `
<script>

function Action(self, action, ...args) {
    // Determine target component, if provided
    let root
    if (action.includes(':')) {
        let rootid = action.split(':')[0]
        action = action.split(':')[1]
        root = document.getElementById(rootid)
    } else {
        let depth = (action.split('').filter(x => x === '$') || []).length
        action = action.replaceAll('$', '')
        root = self
        let dcount = 0
        while (true) {
            if (!root.getAttribute('state')) {
                root = root.parentElement
            } else {
                if (dcount != depth) {
                    root = root.parentElement
                    dcount++
                } else {
                    break
                }
            }
        }
    }
	// Prepare form data
	let formdata = new FormData()
	formdata.set('State', root.getAttribute('state'))
	formdata.set('Args', JSON.stringify(args))
	// Make request
	fetch("/SSA/"+root.getAttribute('name')+"/"+action, {
		method: 'POST',
		body: formdata
	}).then(resp => {
        if (resp.headers.get('X-Redirect')) {
            window.location.href = resp.headers.get('X-Redirect')
			return ''
        }
		return resp.text()
	}).then(data => {
		if (data) {
			root.outerHTML = data
		}
	}).catch(err => {
		console.log(err)
	})
}

function Bind(self, field) {
	// Find component root
	let root = self
	while (true) {
		if (!root.getAttribute('state')) {
			root = root.parentElement
		} else {
			break
		}
	}
	// Load state
	let state = JSON.parse(decodeURIComponent(root.getAttribute('state')))
	// Set value
	state[field] = self.value
	// Set state
	root.setAttribute('state', JSON.stringify(state))
}

function FormSubmit(self, e) {
	// Prevent default submit
	e.preventDefault()
	// Override state with form data
	let state = JSON.parse(decodeURIComponent(self.getAttribute('state')))
	let form = new FormData(e.target)
	let formdata = Object.fromEntries(form.entries())
	Object.entries(formdata).forEach(pair => {
		state[pair[0]] = pair[1]
	})
	self.setAttribute('state', JSON.stringify(state))
	// Trigger "Submit" action
	Action(self, 'Submit')
	return false
}
</script>
`
