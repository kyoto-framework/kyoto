package lifecycle

import (
	"github.com/kyoto-framework/kyoto"
	"github.com/kyoto-framework/scheduler"
)

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

func Async(b *kyoto.Core, async func() error) {
	Patch(b)
	b.Scheduler.Add(&scheduler.Job{
		Group:   "async",
		Depends: []string{"init"},
		Func:    async,
	})
}

func AfterAsync(b *kyoto.Core, afterasync func() error) {
	Patch(b)
	b.Scheduler.Add(&scheduler.Job{
		Group:   "afterasync",
		Depends: []string{"async"},
		Func:    afterasync,
	})
}
