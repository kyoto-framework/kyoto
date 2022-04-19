package kyoto

import (
	"errors"
	"testing"

	"github.com/kyoto-framework/scheduler"
)

// TestCoreComponent ensures component functionality works as expected
func TestCoreComponent(t *testing.T) {
	// Initialize new core
	core := NewCore()

	// Component, that will test state separation
	core.Component("test1", func(_core *Core) {
		if core.State == _core.State {
			t.Error("State separation failed")
		}
	})

	// Component, that will ensure context is the same
	core.Component("test2", func(_core *Core) {
		if core.Context != _core.Context {
			t.Error("Context is not the same")
		}
	})

	// Component, that will ensure scheduler is the same
	core.Component("test3", func(_core *Core) {
		if core.Scheduler != _core.Scheduler {
			t.Error("Scheduler is not the same")
		}
	})

	// Validate components are attached to state
	if core.State.Get("test1") == nil {
		t.Error("Component 'test1' is not attached to state")
	}
}

// TestCoreScheduler ensures scheduler integrated as expected
func TestCoreScheduler(t *testing.T) {
	// Initialize new core
	core := NewCore()

	// Initialization flag
	flag := false

	// Add a job to scheduler
	core.Scheduler.Dispatch(&scheduler.Job{
		Group: "init",
		Func: func() error {
			flag = true
			return errors.New("ignore this\n")
		},
	})

	// Execute core
	core.Execute()

	// Validate flag
	if !flag {
		t.Error("Scheduler didn't executed")
	}
}
