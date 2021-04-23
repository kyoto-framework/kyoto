package ssc

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"net/url"
	"reflect"
	"sync"
)

// Components store
var cstore = map[Page][]Component{}

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
	// Save to store
	if _, ok := cstore[p]; !ok {
		cstore[p] = []Component{}
	}
	cstore[p] = append(cstore[p], c)
	// Trigger component init
	if c, ok := c.(ComponentInit); ok {
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
	if p, ok := p.(PageInit); ok {
		p.Init()
	}
	// Trigger async in goroutines
	for _, component := range cstore[p] {
		if component, ok := component.(ComponentAsync); ok {
			wg.Add(1)
			go func(wg *sync.WaitGroup, err chan error, c ComponentAsync) {
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
		if component, ok := component.(ComponentAfterAsync); ok {
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

// RenderComponent is a minor entrypoint of rendering.
func RenderComponent(w io.Writer, c Component, t *template.Template, d string) {
	// Init dummy page for rendering
	t, _ = t.Parse(`{{ template "` + d + `" . }}`)
	page := &DummyPage{
		DTemplate:  t,
		DComponent: c,
	}
	// Dender dummy page with component
	RenderPage(w, page)
}

func RenderComponentString(c Component, t *template.Template, d string) string {
	var b bytes.Buffer
	RenderComponent(&b, c, t, d)
	return b.String()
}

func HandleSSA(w io.Writer, t *template.Template, componentname string, state string, action string, argsstr string, clist []Component) {
	// Find component
	var found bool = false
	var component Component
	for _, c := range clist {
		if reflect.ValueOf(c).Elem().Type().Name() == componentname {
			component = c
			found = true
		}
	}
	if !found {
		panic("Can't find component. Perhaps, you forgot to register it while calling HandleSSA")
	}
	// Init
	if component, ok := component.(ComponentInit); ok {
		component.Init(&DummyPage{})
	}
	// Init component with state
	state, _ = url.QueryUnescape(state)
	if err := json.Unmarshal([]byte(state), &component); err != nil {
		panic(err)
	}
	// Extract arguments
	var args []interface{}
	json.Unmarshal([]byte(argsstr), &args)
	// Call action
	if component, ok := component.(ComponentActions); ok {
		component.Actions()[action](args...)
	}
	// Render component
	err := t.Execute(w, component)
	if err != nil {
		panic(err)
	}
}
