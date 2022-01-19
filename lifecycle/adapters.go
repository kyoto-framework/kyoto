package lifecycle

import (
	"github.com/kyoto-framework/kyoto"
	"github.com/kyoto-framework/scheduler"
)

// Init is a function to define an initialization job.
func Init(b *kyoto.Core, init func()) {
	Patch(b)
	b.Scheduler.Add(&scheduler.Job{
		Group: "init",
		Func: func() error {
			init()
			return nil
		},
	})
}

// Async is a function to define an asynchronous job.
// Will be executed after initialization step.
func Async(b *kyoto.Core, async func() error) {
	Patch(b)
	b.Scheduler.Add(&scheduler.Job{
		Group:   "async",
		Depends: []string{"init"},
		Func:    async,
	})
}

// AfterAsync is a function to define a job,
// that will be executed after asynchronous step.
func AfterAsync(b *kyoto.Core, afterasync func() error) {
	Patch(b)
	b.Scheduler.Add(&scheduler.Job{
		Group:   "afterasync",
		Depends: []string{"async"},
		Func:    afterasync,
	})
}
