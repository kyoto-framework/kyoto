package smode

import (
	"github.com/kyoto-framework/kyoto"
	"github.com/kyoto-framework/kyoto/helpers"
	"github.com/kyoto-framework/kyoto/lifecycle"
	"github.com/kyoto-framework/kyoto/render"
	"github.com/kyoto-framework/scheduler"
)

var (
	cmap = map[interface{}]*kyoto.Core{}
	pmap = map[Component]Page{}
)

func Adapt(item interface{}) func(*kyoto.Core) {
	return func(core *kyoto.Core) {
		// Save core mapping
		cmap[item] = core
		// Extract page
		page := pmap[item]
		// Adapt template
		if _item, ok := item.(ImplementsTemplate); ok {
			render.Template(core, _item.Template)
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
		// Schedule state export
		core.Scheduler.Add(scheduler.Job{
			Group:   "state",
			Depends: []string{"afterasync"},
			Func: func() error {
				for k, v := range structmap(item) {
					core.State.Set(k, v)
				}
				return nil
			},
		})
		// Schedule global cmap cleanup
		core.Scheduler.Add(scheduler.Job{
			Group:   "cleanup",
			Depends: []string{"render"},
			Func: func() error {
				delete(cmap, item)
				return nil
			},
		})
	}
}

func RegC(page Page, component Component) Component {
	// Save component/page mapping to temporary store
	pmap[component] = page
	// Create custom core for component to scope state
	_core := kyoto.NewCore()
	_core.Scheduler = cmap[page].Scheduler
	_core.Context = cmap[page].Context
	// Inject component name into state
	_core.State.Set("internal:name", helpers.ComponentName(component))
	// Execute builder receiver
	switch c := component.(type) {
	case func(*kyoto.Core): // In case of builder receiver, just call it
		c(_core)
	case interface{}: // In case of struct component, adapt and call it
		Adapt(c)(_core)
	}
	// Remove component/page mapping from temporary store
	delete(pmap, component)
	// Return a state of the component
	return _core.State.Export()
}
