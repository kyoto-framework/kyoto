package ssc

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

func PageHandler(p Page) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		SetContext(p, "internal:request", r)
		SetContext(p, "internal:rwriter", rw)
		RenderPage(rw, p)
		DelContext(p, "")
	}
}

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
		component.Init(dummypage)
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
	if component, ok := component.(ImplementsActions); ok {
		component.Actions()[aname](args...)
	} else {
		panic("Component not implements Actions, unexpected behavior")
	}
	// Prepare template
	t := ssacs.TemplateBuilder()
	t = template.Must(t.Parse(fmt.Sprintf(`{{ template "%s" . }}`, cname)))
	// Render component
	terr := t.Execute(rw, component)
	if terr != nil {
		panic(terr)
	}
	// Clear context
	DelContext(dummypage, "")
}
