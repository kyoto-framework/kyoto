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
)

func Adapt(item interface{}) func(*kyoto.Core) {
	return func(core *kyoto.Core) {
		// Save core mapping
		cmap[item] = core
		// Adapt template
		if _page, ok := item.(ImplementsTemplate); ok {
			render.Template(core, _page.Template)
		}
		// Adapt lifecycle
		if _page, ok := item.(ImplementsInit); ok {
			lifecycle.Init(core, _page.Init)
		}
		if _page, ok := item.(ImplementsAsync); ok {
			lifecycle.Async(core, _page.Async)
		}
		if _page, ok := item.(ImplementsAfterAsync); ok {
			lifecycle.AfterAsync(core, func() error {
				_page.AfterAsync()
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
	// Return a state of the component
	return _core.State.Export()
}
