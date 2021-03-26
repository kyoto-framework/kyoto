package gofr

import (
	"io"
	"reflect"
	"sync"
)

var cstore = map[Page][]Component{}

// RegisterComponent is used while defining components in the Init() section of Page
func RegisterComponent(p Page, c Component) Component {
	if _, ok := cstore[p]; !ok {
		cstore[p] = []Component{}
	}
	cstore[p] = append(cstore[p], c)
	return c
}

// RegC - Shortcut for RegisterComponent
func RegC(p Page, c Component) Component {
	return RegisterComponent(p, c)
}

// RenderPage is a main entrypoint of rendering. Responsible for rendering and components lifecycle
func RenderPage(w io.Writer, p Page) {
	var wg sync.WaitGroup
	var err = make(chan error, 1000)
	p.Init()
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
	wg.Wait()
	for _, component := range cstore[p] {
		component.AfterAsync()
	}
	delete(cstore, p)
	p.Template().Execute(w, reflect.ValueOf(p).Elem())
}

func RenderComponent(w io.Writer, c Component) {

}

func RenderComponentString(c Component) string {
	return ""
}
