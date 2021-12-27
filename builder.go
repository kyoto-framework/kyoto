package kyoto

import (
	"html/template"
	"io"

	"github.com/kyoto-framework/kyoto/helpers"
	"github.com/kyoto-framework/scheduler"
)

type Builder struct {
	TB        func() *template.Template
	State     State
	Context   Context
	Scheduler *scheduler.Scheduler

	ActionMap map[string]func(...interface{})
}

func NewBuilder() *Builder {
	// Init state and context
	_state := NewStore()
	_context := NewStore()
	// Init schedule
	_scheduler := scheduler.New()
	_scheduler.Workers = 10
	// Add empty jobs to schedule (for groups definitions)
	// Without empty jobs, "Depends" clause will not work as expected
	_scheduler.Add(scheduler.Job{
		Group: "init",
		Name:  "empty",
		Func:  func() error { return nil },
	})
	_scheduler.Add(scheduler.Job{
		Group:   "async",
		Name:    "empty",
		Depends: []string{"init"},
		Func:    func() error { return nil },
	})
	_scheduler.Add(scheduler.Job{
		Group:   "afterasync",
		Name:    "empty",
		Depends: []string{"async"},
		Func:    func() error { return nil },
	})
	_scheduler.Add(scheduler.Job{
		Group:   "state",
		Name:    "empty",
		Depends: []string{"afterasync"},
		Func:    func() error { return nil },
	})
	// Return new builder
	return &Builder{
		State:     _state,
		Context:   _context,
		Scheduler: _scheduler,
		ActionMap: map[string]func(...interface{}){},
	}
}

// Template is a method for setting page template builder
func (b *Builder) Template(tb func() *template.Template) {
	b.TB = tb
}

// Workers is a method for setting scheduler workers count (10 by default)
func (b *Builder) Workers(num int) {
	b.Scheduler.Workers = num
}

func (b *Builder) Action(alias string, action func(args ...interface{})) {
	b.ActionMap[alias] = action
}

func (b *Builder) Component(alias string, component interface{}) {
	// Create custom builder for component to scope state
	_b := NewBuilder()
	_b.Scheduler = b.Scheduler
	// Inject component name into state
	_b.State.Set("internal:name", helpers.ComponentName(component))
	// Switch behavior based on type of component
	switch component := component.(type) {
	case func(*Builder): // Function component
		// Execute builder receiver
		component(_b)
	default: // Not supported
		panic("Component type is not supported")
	}
	// Schedule a job for state gathering
	b.Scheduler.Add(scheduler.Job{
		Group:   "state",
		Depends: []string{"afterasync"},
		Func: func() error {
			// Gather state
			b.State.Set(alias, _b.State.Export())
			return nil
		},
	})
}

func (b *Builder) Init(init func()) {
	b.Scheduler.Add(scheduler.Job{
		Group: "init",
		Func: func() error {
			init()
			return nil
		},
	})
}

func (b *Builder) Async(async func() error) {
	b.Scheduler.Add(scheduler.Job{
		Group:   "async",
		Depends: []string{"init"},
		Func:    async,
	})
}

// Render is a method to add rendering job to scheduler
func (b *Builder) Render(w io.Writer) {
	// Check template builder
	if b.TB == nil {
		panic("Template builder is not set")
	}
	// Add rendering job
	b.Scheduler.Add(scheduler.Job{
		Group:   "render",
		Depends: []string{"init", "async", "afterasync", "state"},
		Func: func() error {
			return b.TB().Execute(w, b.State.Export())
		},
	})
}

// Execute is a method to trigger scheduler execution
func (b *Builder) Execute() {
	// Execute scheduler
	b.Scheduler.Execute()
	// Analyze errors
	for job, res := range b.Scheduler.Results {
		if res != nil {
			panic("Error in job " + job + ": " + res.Error())
		}
	}
}
