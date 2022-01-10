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
	// Return new core
	return &Core{
		State:     _state,
		Context:   _context,
		Scheduler: _scheduler,
	}
}

func (core *Core) Component(alias string, component func(*Core)) {
	// Create custom core for component to scope state
	_core := NewCore()
	_core.Scheduler = core.Scheduler
	_core.Context = core.Context
	// Inject component name into state
	_core.State.Set("internal:name", helpers.ComponentName(component))
	// Execute core receiver
	component(_core)
	// Schedule a job for state gathering
	core.Scheduler.Add(&scheduler.Job{
		Group:   "state",
		Depends: []string{"afterasync"},
		Func: func() error {
			// Gather state
			core.State.Set(alias, _core.State.Export())
			return nil
		},
	})
}

// Execute is a method to trigger scheduler execution
func (core *Core) Execute() {
	// Execute scheduler
	core.Scheduler.Execute()
	// Analyze errors
	for job, res := range core.Scheduler.Results {
		if res != nil {
			panic("Error in job " + job + ": " + res.Error())
		}
	}
}
