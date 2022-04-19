package lifecycle

import (
	"testing"

	"github.com/kyoto-framework/kyoto"
)

// TestAdapters will ensure that:
// - Jobs are provisioning as expected
// - Adapters are calling core patch
func TestAdapters(t *testing.T) {
	// Initialize core
	core := kyoto.NewCore()
	// Initialize lifecycle handlers
	Init(core, func() {})
	Async(core, func() error { return nil })
	AfterAsync(core, func() error { return nil })
	// Check jobs are actually added
	if len(core.Scheduler.Jobs) != 3 {
		t.Error("Something went wrong and jobs count is not correct")
	}
	// Check jobs groups
	init := false
	async := false
	afterasync := false
	for _, job := range core.Scheduler.Jobs {
		if job.Group == "init" {
			init = true
		} else if job.Group == "async" {
			async = true
		} else if job.Group == "afterasync" {
			afterasync = true
		}
	}
	if !init {
		t.Error("Lifecycle adapter is not provisioning an init job")
	}
	if !async {
		t.Error("Lifecycle adapter is not provisioning an async job")
	}
	if !afterasync {
		t.Error("Lifecycle adapter is not provisioning an afterasync job")
	}
}

// TestAdaptersTiming will ensure that job groups are executing in the correct order
func TestAdaptersTiming(t *testing.T) {
	// Initialize core
	core := kyoto.NewCore()
	// Initialize lifecycle handlers with tests
	counter := 0
	Init(core, func() {
		if counter != 0 {
			t.Error("Init job is not executing first")
		}
		counter++
	})
	Async(core, func() error {
		if counter != 1 {
			t.Error("Async job is not executing second")
		}
		counter++
		return nil
	})
	AfterAsync(core, func() error {
		if counter != 2 {
			t.Error("AfterAsync job is not executing third")
		}
		counter++
		return nil
	})
	// Execute
	core.Execute()
}
