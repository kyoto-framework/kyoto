package lifecycle

import (
	"github.com/kyoto-framework/kyoto"
	"github.com/kyoto-framework/scheduler"
)

// Init is a function to define an initialization job.
func Init(b *kyoto.Core, init func()) {
	b.Scheduler.Dispatch(&scheduler.Job{
		Group:  "init",
		Before: []string{"render"},
		Func: func() error {
			init()
			return nil
		},
	})
}

// Async is a function to define an asynchronous job.
// Will be executed after initialization step.
func Async(b *kyoto.Core, async func() error) {
	b.Scheduler.Dispatch(&scheduler.Job{
		Group:  "async",
		Before: []string{"render"},
		After:  []string{"init"},
		Func:   async,
	})
}

// AfterAsync is a function to define a job,
// that will be executed after asynchronous step.
func AfterAsync(b *kyoto.Core, afterasync func() error) {
	b.Scheduler.Dispatch(&scheduler.Job{
		Group:  "afterasync",
		Before: []string{"render"},
		After:  []string{"init", "async"},
		Func:   afterasync,
	})
}
