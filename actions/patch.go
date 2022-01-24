package actions

import (
	"errors"

	"github.com/kyoto-framework/kyoto"
	"github.com/kyoto-framework/scheduler"
)

// Patch is a function that patches scheduler with state population, action and flush jobs.
// Also it removes all jobs, not needed for action call (like "async" and "afterasync").
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
			actions := GetActions(core)
			if action, found := actions[params.Action]; found {
				action(params.Args...)
			} else {
				return errors.New("action " + params.Action + " not found for component " + params.Component)
			}
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
