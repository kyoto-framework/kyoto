package kyoto

import (
	"fmt"

	"github.com/kyoto-framework/kyoto/helpers"
	"github.com/kyoto-framework/scheduler"
)

// Core is a start point of kyoto.
// It cointains only basics, while external packages
// are using core for providing all the kyoto functionality.
//
// Core consists of state, context and scheduler.
// Each component is a receiver of core.
// Use "adapter" functions to inject functionality into kyoto core:
// (f.e. lifecycle.Init(core, ...) ).
type Core struct {
	State     State
	Context   Context
	Scheduler *scheduler.Scheduler
}

// NewCore is a constructor for Core.
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

// Component is a method to inject nested component.
// Under the hood it composes custom Core for state separation and nesting.
func (core *Core) Component(alias string, component func(*Core)) {
	// Create custom core for component to scope state
	_core := NewCore()
	_core.Scheduler = core.Scheduler
	_core.Context = core.Context
	// Inject component name into state
	_core.State.Set("internal:name", helpers.ComponentName(component))
	// Execute core receiver
	component(_core)
	// Bind state to the root state
	// Map acts like a reference, so it's safe to continue modify state
	core.State.Set(alias, _core.State.Export())
}

// Execute is a method to trigger scheduler execution.
func (core *Core) Execute() {
	// Execute scheduler
	core.Scheduler.Execute()
	// Analyze errors
	for job, res := range core.Scheduler.Results {
		if res != nil {
			fmt.Printf("%s failed -> %s", job, res.Error())
		}
	}
}
