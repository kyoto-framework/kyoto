package ssc

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
)

// Temporary components store. Will be cleared in the end of lifecycle
var cstore = map[Page][]Component{}

// Persistent components store, used by SSA. Components are stored in pair with page template
var ssastore = map[string]ssacomponentstore{}

// Part of persistent components store
type ssacomponentstore struct {
	TemplateBuilder func() *template.Template
	ComponentType   reflect.Type
}

// Dummy page for rendering component
type DummyPage struct {
	DTemplate  *template.Template
	DComponent Component
}

func (p *DummyPage) Template() *template.Template {
	return p.DTemplate
}

func (p *DummyPage) Init() {
	p.DComponent = RegC(p, p.DComponent)
}

func (*DummyPage) Meta() Meta {
	return Meta{}
}

// RegisterComponent is used while defining components in the Init() section
func RegisterComponent(p Page, c Component) Component {
	// Save to stores
	if _, ok := cstore[p]; !ok {
		cstore[p] = []Component{}
	}
	if _, ok := ssastore[reflect.ValueOf(c).Elem().Type().Name()]; !ok {
		// Extract component type
		ctype := reflect.ValueOf(c).Elem().Type()
		// Save to store
		ssastore[reflect.ValueOf(c).Elem().Type().Name()] = ssacomponentstore{
			TemplateBuilder: p.Template,
			ComponentType:   ctype,
		}
	}
	cstore[p] = append(cstore[p], c)
	// Trigger component init
	if c, ok := c.(ImplementsNestedInit); ok {
		c.Init(p)
	}
	// Return component for external assignment
	return c
}

// RegC - Shortcut for RegisterComponent
func RegC(p Page, c Component) Component {
	return RegisterComponent(p, c)
}

// RenderPage is a main entrypoint of rendering. Responsible for rendering and components lifecycle
func RenderPage(w io.Writer, p Page) {
	// Async specific state
	var wg sync.WaitGroup
	var err = make(chan error, 1000)
	// Trigger init
	if p, ok := p.(ImplementsInit); ok {
		p.Init()
	}
	// Trigger async in goroutines
	for _, component := range cstore[p] {
		if component, ok := component.(ImplementsAsync); ok {
			wg.Add(1)
			go func(wg *sync.WaitGroup, err chan error, c ImplementsAsync) {
				defer wg.Done()
				_err := c.Async()
				if _err != nil {
					err <- _err
				}
			}(&wg, err, component)
		}
	}
	// Wait for async completion
	wg.Wait()
	// Trigger aftersync
	for _, component := range cstore[p] {
		if component, ok := component.(ImplementsAfterAsync); ok {
			component.AfterAsync()
		}
	}
	// Clear components store (not needed more)
	delete(cstore, p)
	// Execute template
	terr := p.Template().Execute(w, reflect.ValueOf(p).Elem())
	if terr != nil {
		panic(terr)
	}
}

func SSAHandler(rw http.ResponseWriter, r *http.Request) {
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
		component.Init(&DummyPage{})
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
}
