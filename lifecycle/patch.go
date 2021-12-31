package lifecycle

import (
	"github.com/kyoto-framework/kyoto"
	"github.com/kyoto-framework/scheduler"
)

func Patch(b *kyoto.Core) {
	// Check if patching is needed
	if b.Context.Get("internal:lifecycle") != nil {
		return
	}
	// Set patch flag
	b.Context.Set("internal:lifecycle", true)
	// Add empty jobs to schedule (for groups definitions)
	// Without empty jobs, "Depends" clause will not work as expected
	b.Scheduler.Add(scheduler.Job{
		Group: "init",
		Name:  "empty",
		Func:  func() error { return nil },
	})
	b.Scheduler.Add(scheduler.Job{
		Group:   "async",
		Name:    "empty",
		Depends: []string{"init"},
		Func:    func() error { return nil },
	})
	b.Scheduler.Add(scheduler.Job{
		Group:   "afterasync",
		Name:    "empty",
		Depends: []string{"async"},
		Func:    func() error { return nil },
	})
	b.Scheduler.Add(scheduler.Job{
		Group:   "state",
		Name:    "empty",
		Depends: []string{"afterasync"},
		Func:    func() error { return nil },
	})
}
