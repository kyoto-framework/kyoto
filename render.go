package gofr

import (
	"bytes"
	"html/template"
	"io"
	"reflect"
	"sync"
)

// Components store
var cstore = map[interface{}][]Component{}

// RegisterComponent is used while defining components in the Init() section
func RegisterComponent(p interface{}, c Component) Component {
	// Save to store
	if _, ok := cstore[p]; !ok {
		cstore[p] = []Component{}
	}
	cstore[p] = append(cstore[p], c)
	// Trigger component init
	c.Init()
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
	p.Init()
	// Trigger async in goroutines
	for _, component := range cstore[p] {
		wg.Add(1)
		go func(wg *sync.WaitGroup, err chan error, c Component) {
			defer wg.Done()
			_err := c.Async()
			if _err != nil {
				err <- _err
			}
		}(&wg, err, component)
	}
	// Wait for async completion
	wg.Wait()
	// Trigger aftersync
	for _, component := range cstore[p] {
		component.AfterAsync()
	}
	// Clear components store (not needed more)
	delete(cstore, p)
	// Execute template
	p.Template().Execute(w, reflect.ValueOf(p).Elem())
}

// RenderComponent is a minor entrypoint of rendering.
func RenderComponent(w io.Writer, c Component, t *template.Template, d string) {
	// Async specific state
	var wg sync.WaitGroup
	var err = make(chan error, 1000)
	// Trigger init
	c.Init()
	// Trigger async in goroutines
	for _, component := range cstore[c] {
		wg.Add(1)
		go func(wg *sync.WaitGroup, err chan error, c Component) {
			defer wg.Done()
			_err := c.Async()
			if _err != nil {
				err <- _err
			}
		}(&wg, err, component)
	}
	// Wait for async completion
	wg.Wait()
	// Trigger aftersync
	for _, component := range cstore[c] {
		component.AfterAsync()
	}
	// Clear components store (not needed more)
	delete(cstore, c)
	// Render component
	t, _ = t.Parse(`{{ template "` + d + `" . }}`)
	t.Execute(w, reflect.ValueOf(c).Elem())
}

func RenderComponentString(c Component, t *template.Template, d string) string {
	var b bytes.Buffer
	RenderComponent(&b, c, t, d)
	return b.String()
}
