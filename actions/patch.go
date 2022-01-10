package actions

import (
	"github.com/kyoto-framework/kyoto"
	"github.com/kyoto-framework/scheduler"
)

func Patch(core *kyoto.Core, params Parameters) {
	// Check template builder
	if core.Context.Get("internal:render:tb") == nil {
		panic("No template specified for page")
	}
	// Patch and cleanup existing jobs
	jobs := []*scheduler.Job{}
	for _, job := range core.Scheduler.Jobs {
		if job.Name == "empty" {
			jobs = append(jobs, job)
			continue
		}
		if job.Group == "init" {
			jobs = append(jobs, job)
			continue
		}
		if job.Group == "state" {
			job.Depends = []string{"action"}
			jobs = append(jobs, job)
			continue
		}
	}
	core.Scheduler.Jobs = jobs
	// Add state population job
	core.Scheduler.Add(&scheduler.Job{
		Group:   "populate",
		Depends: []string{"init"},
		Func: func() error {
			for k, v := range params.State {
				core.State.Set(k, v)
			}
			return nil
		},
	})
	// Add action job
	core.Scheduler.Add(&scheduler.Job{
		Group:   "action",
		Depends: []string{"populate"},
		Func: func() error {
			GetActions(core)[params.Action](params.Args...)
			return nil
		},
	})
	// Add final flush job
	core.Scheduler.Add(&scheduler.Job{
		Group:   "flush",
		Depends: []string{"action"},
		Func: func() error {
			Flush(core)
			return nil
		},
	})
}
