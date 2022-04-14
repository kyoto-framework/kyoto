package smode

import (
	"encoding/json"
	"html/template"
	"log"
	"reflect"
	"sync"

	"github.com/kyoto-framework/kyoto"
	"github.com/kyoto-framework/kyoto/actions"
	"github.com/kyoto-framework/kyoto/helpers"
	"github.com/kyoto-framework/kyoto/lifecycle"
	"github.com/kyoto-framework/kyoto/render"
	"github.com/kyoto-framework/scheduler"
)

var (
	cmap  = map[interface{}]*kyoto.Core{}
	pmap  = map[Component]Page{}
	cmapm = sync.Mutex{}
	pmapm = sync.Mutex{}
)

// Adapt transforms passed struct component to a classic functional component.
//
// Usage:
//  core.Component("Component", smode.Adapt(&MyComponent{}))
func Adapt(item interface{}) func(*kyoto.Core) {
	return func(core *kyoto.Core) {
		// Aquire locks
		cmapm.Lock()
		pmapm.Lock()
		defer cmapm.Unlock()
		defer pmapm.Unlock()
		// If no page, need to create a reference
		if cmap[item] == nil {
			cmap[item] = core
		}
		if pmap[item] == nil {
			pmap[item] = item
		}
		// In case of page, need to create a new instance
		if _, ispage := item.(ImplementsTemplate); ispage {
			item = reflect.New(reflect.TypeOf(item).Elem()).Interface().(Page)
		} else if _, ispage := item.(ImplementsTemplateWithPage); ispage {
			item = reflect.New(reflect.TypeOf(item).Elem()).Interface().(Page)
		}
		// Inject component name
		core.State.Set("internal:name", helpers.ComponentName(item))
		// Save core mapping
		cmap[item] = core
		// Extract page
		page := pmap[item]
		// Adapt rendering
		if _item, ok := item.(ImplementsTemplate); ok {
			render.Template(core, _item.Template)
		}
		if _item, ok := item.(ImplementsWriter); ok {
			render.Writer(core, _item.Writer)
		}
		// Adapt lifecycle
		if _item, ok := item.(ImplementsInit); ok {
			lifecycle.Init(core, _item.Init)
		}
		if _item, ok := item.(ImplementsInitWithPage); ok {
			lifecycle.Init(core, func() {
				_item.Init(page)
			})
		}
		if _item, ok := item.(ImplementsAsync); ok {
			lifecycle.Async(core, _item.Async)
		}
		if _item, ok := item.(ImplementsAsyncWithPage); ok {
			lifecycle.Async(core, func() error {
				return _item.Async(page)
			})
		}
		if _item, ok := item.(ImplementsAfterAsync); ok {
			lifecycle.AfterAsync(core, func() error {
				_item.AfterAsync()
				return nil
			})
		}
		if _item, ok := item.(ImplementsAfterAsyncWithPage); ok {
			lifecycle.AfterAsync(core, func() error {
				_item.AfterAsync(page)
				return nil
			})
		}
		// Adapt actions
		adaptactions := func(amap ActionMap) {
			wrapaction := func(action Action) func(...interface{}) {
				return func(args ...interface{}) {
					statebts, _ := json.Marshal(core.State.Export())
					json.Unmarshal(statebts, item)
					action(args...)
				}
			}
			for name, action := range amap {
				// Wrap and register an action.
				// This "hack" is required because Core receiver is called before action patch,
				// so we can't override state population
				actions.Define(core, name, wrapaction(action))
			}
		}
		if _item, ok := item.(ImplementsActions); ok {
			adaptactions(_item.Actions())
		}
		if _item, ok := item.(ImplementsActionsWithPage); ok {
			adaptactions(_item.Actions(page))
		}

		// Schedule state export
		core.Scheduler.Add(&scheduler.Job{
			Group:  "state",
			After:  []string{"afterasync", "action"}, // Export state only after "afterasync" or "action", otherwise it will be executed immediately
			Before: []string{"render"},
			Func: func() error {
				for k, v := range structmap(item) {
					core.State.Set(k, v)
				}
				return nil
			},
		})
		// Schedule global cmap cleanup
		core.Scheduler.Add(&scheduler.Job{
			Group: "cleanup",
			After: []string{"render"},
			Func: func() error {
				cmapm.Lock()
				log.Printf("Calling cleanup for %s %p", helpers.ComponentName(item), item)
				delete(cmap, item)
				delete(pmap, item)
				cmapm.Unlock()
				return nil
			},
		})
	}
}

// UseC registers passed component in lifecycle and returns it.
//
// Usage:
//  page.Component = smode.UseC(page, &MyComponent{})
//
// TODO: Not working for functional components (hash of unhashable type func(*kyoto.Core))
func UseC(page Page, component Component) Component {
	// Save component/page mapping to temporary store
	pmap[component] = page
	// Create custom core for component to scope state
	_core := kyoto.NewCore()
	_core.Scheduler = cmap[page].Scheduler
	_core.Context = cmap[page].Context
	// Inject component name into state
	_core.State.Set("internal:name", helpers.ComponentName(component))
	// Execute builder receiver
	Adapt(component)(_core)
	// Remove component/page mapping from temporary store
	pmapm.Lock()
	delete(pmap, component)
	pmapm.Unlock()
	// Return a state of the component
	return _core.State.Export()
}

// Register registers dynamic components in the actions system.
// Adapted to be compatible with struct components.
func Register(components ...interface{}) {
	for _, component := range components {
		switch c := component.(type) {
		case func(*kyoto.Core):
			// Register builder
			actions.RegisterWithName(helpers.ComponentName(component), c)
		case Component:
			// Register struct with adapt
			actions.RegisterWithName(helpers.ComponentName(component), Adapt(c))
		}
	}
}

// Redirect is a wrapper around `render.Redirect`.
// See original function for details.
func Redirect(page Page, target string, code int) {
	// Extract core
	core := cmap[page]
	// Redirect
	render.Redirect(core, target, code)
}

// FuncMap is a wrapper around `render.FuncMap`.
// See original function for details.
func FuncMap(p Page) template.FuncMap {
	// Extract kyoto.Core by page
	core := cmap[p]
	if core == nil {
		panic("Can't find core for a provided page")
	}
	// Return funcmap with core
	return render.FuncMap(core)
}
