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
	"time"
)

// Debug flag
var BENCH_HANDLERS = false

// Page handler helper
// Simple wrapper around RenderPage with context setting
func PageHandler(p Page) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		SetContext(p, "internal:request", r)
		SetContext(p, "internal:rwriter", rw)
		RenderPage(rw, p)
		DelContext(p, "")
	}
}

// SSA handler
func SSAHandler(rw http.ResponseWriter, r *http.Request) {
	// Init dummy page and set context
	dummypage := &DummyPage{}
	SetContext(dummypage, "internal:request", r)
	SetContext(dummypage, "internal:rwriter", rw)
	// Extract component action and name from route
	tokens := strings.Split(r.URL.Path, "/")
	cname := tokens[2]
	aname := tokens[3]
	// Find ssacomponentstore
	ssacs, found := ssastore[cname]
	// Panic, if not found
	if !found {
		panic("Can't find component. Seems like it's not registered")
	}
	// Create component
	component := reflect.New(ssacs.ComponentType).Interface().(Component)
	// Init
	if component, ok := component.(ImplementsNestedInit); ok {
		st := time.Now()
		component.Init(dummypage)
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
	// Prepare template
	st = time.Now()
	t := ssacs.TemplateBuilder()
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
		log.Println("Execute time", reflect.TypeOf(component), et)
	}
	// Clear context
	DelContext(dummypage, "")
}

// SSA handler with passed page, instead of dummy
func SSAHandlerWithPage(page Page) func(rw http.ResponseWriter, r *http.Request) {
	p := page
	return func(rw http.ResponseWriter, r *http.Request) {
		// Init dummy page and set context
		SetContext(p, "internal:request", r)
		SetContext(p, "internal:rwriter", rw)
		// Extract component action and name from route
		tokens := strings.Split(r.URL.Path, "/")
		cname := tokens[2]
		aname := tokens[3]
		// Find ssacomponentstore
		ssacs, found := ssastore[cname]
		// Panic, if not found
		if !found {
			panic("Can't find component. Seems like it's not registered")
		}
		// Create component
		component := reflect.New(ssacs.ComponentType).Interface().(Component)
		// Init
		if component, ok := component.(ImplementsNestedInit); ok {
			st := time.Now()
			component.Init(p)
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
		// Prepare template
		st = time.Now()
		t := ssacs.TemplateBuilder()
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
			log.Println("Execute time", reflect.TypeOf(component), et)
		}
		// Clear context
		DelContext(p, "")
	}
}
