package ssc

import (
	"html/template"
	"io"
	"log"
	"reflect"
	"sync"
	"time"
)

// Debug flag
var BENCH_LOWLEVEL = false

// Temporary components store. Will be cleared in the end of lifecycle
var cstore = map[Page][]Component{}

// Persistent components store, used by SSA. Components are stored in pair with page template
var ssastore = map[string]ssacomponentstore{}

// Part of persistent components store
type ssacomponentstore struct {
	TemplateBuilder func() *template.Template
	ComponentType   reflect.Type
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
		st := time.Now()
		p.Init()
		et := time.Since(st)
		if BENCH_LOWLEVEL {
			log.Println("Init time", reflect.TypeOf(p), et)
		}
	}
	// Trigger async in goroutines
	for _, component := range cstore[p] {
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
	// Wait for async completion
	wg.Wait()
	// Trigger aftersync
	for _, component := range cstore[p] {
		if component, ok := component.(ImplementsAfterAsync); ok {
			st := time.Now()
			component.AfterAsync()
			et := time.Since(st)
			if BENCH_LOWLEVEL {
				log.Println("AfterAsync time", reflect.TypeOf(component), et)
			}
		}
	}
	// Clear components store (not needed more)
	delete(cstore, p)
	// Execute template
	st := time.Now()
	terr := p.Template().Execute(w, reflect.ValueOf(p).Elem())
	et := time.Since(st)
	if BENCH_LOWLEVEL {
		log.Println("Execute time", et)
	}
	if terr != nil {
		panic(terr)
	}
}
