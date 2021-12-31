package kyoto

import (
	"github.com/kyoto-framework/kyoto/helpers"
	"github.com/kyoto-framework/scheduler"
)

type Core struct {
	State     State
	Context   Context
	Scheduler *scheduler.Scheduler
}

func NewCore() *Core {
	// Init state and context
	_state := NewStore()
	_context := NewStore()
	// Init schedule
	_scheduler := scheduler.New()
	_scheduler.Workers = 10
	// Return new builder
	return &Core{
		State:     _state,
		Context:   _context,
		Scheduler: _scheduler,
	}
}

// Workers is a method for setting scheduler workers count (10 by default)
func (b *Core) Workers(num int) {
	b.Scheduler.Workers = num
}

func (b *Core) Component(alias string, component interface{}) {
	// Create custom builder for component to scope state
	_b := NewCore()
	_b.Scheduler = b.Scheduler
	_b.Context = b.Context
	// Inject component name into state
	_b.State.Set("internal:name", helpers.ComponentName(component))
	// Switch behavior based on type of component
	switch component := component.(type) {
	case func(*Core): // Function component
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

// Execute is a method to trigger scheduler execution
func (b *Core) Execute() {
	// Execute scheduler
	b.Scheduler.Execute()
	// Analyze errors
	for job, res := range b.Scheduler.Results {
		if res != nil {
			panic("Error in job " + job + ": " + res.Error())
		}
	}
}
