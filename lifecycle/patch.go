package lifecycle

import (
	"github.com/kyoto-framework/kyoto"
	"github.com/kyoto-framework/scheduler"
)

// Patch is a function to inject lifecycle functionality into kyoto.Core.
// Usually you don't need to call this function directly. Patch is applied automatically
// when you call one of lifecycle adapters.
func Patch(b *kyoto.Core) {
	// Check if patching is needed
	if b.Context.Get("internal:lifecycle") != nil {
		return
	}
	// Set patch flag
	b.Context.Set("internal:lifecycle", true)
	// Add empty jobs to schedule (for groups definitions)
	// Without empty jobs, "Depends" clause will not work as expected
	b.Scheduler.Dispatch(&scheduler.Job{
		Group: "init",
		Name:  "placeholder",
		Func:  func() error { return nil },
	})
	b.Scheduler.Dispatch(&scheduler.Job{
		Group: "async",
		Name:  "placeholder",
		After: []string{"init"},
		Func:  func() error { return nil },
	})
	b.Scheduler.Dispatch(&scheduler.Job{
		Group:  "afterasync",
		Name:   "placeholder",
		After:  []string{"async"},
		Before: []string{"render"},
		Func:   func() error { return nil },
	})
}
