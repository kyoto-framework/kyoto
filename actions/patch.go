package actions

import (
	"errors"

	"github.com/kyoto-framework/kyoto"
	"github.com/kyoto-framework/scheduler"
)

// Patch is a function that patches scheduler with state population, action and flush jobs.
// Also it removes all jobs, not needed for action call (like "async" and "afterasync").
func Patch(core *kyoto.Core, params Parameters) {
	// Patch and cleanup existing jobs
	jobs := []*scheduler.Job{}
	for _, job := range core.Scheduler.Jobs {
		if job.Name == "placeholder" { // Need lifecycle placeholders
			jobs = append(jobs, job)
			continue
		}
		if job.Group == "init" { // Need initialization
			jobs = append(jobs, job)
			continue
		}
		if job.Group == "state" { // Need state export (smode)
			jobs = append(jobs, job)
			continue
		}
	}
	core.Scheduler.Jobs = jobs
	// Add state population job
	core.Scheduler.Dispatch(&scheduler.Job{
		Group:  "populate",
		After:  []string{"init"},
		Before: []string{"action"},
		Func: func() error {
			// Iterate over values
			for k, v := range params.State {
				if v, ok := v.(map[string]any); ok && v["State"] != nil && v["State"].(bool) == true {
					core.State.Get(k).(interface {
						SetAny(any)
					}).SetAny(v["Value"])
				} else {
					core.State.Set(k, v)
				}
			}
			return nil
		},
	})
	// Add action job
	core.Scheduler.Dispatch(&scheduler.Job{
		Group:  "action",
		Before: []string{"render"},
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
	core.Scheduler.Dispatch(&scheduler.Job{
		Group: "render",
		Func: func() error {
			Flush(core)
			return nil
		},
	})
}
