package actions

import (
	"log"

	"github.com/kyoto-framework/kyoto"
	"github.com/kyoto-framework/scheduler"
)

func Patch(b *kyoto.Builder, params Parameters) {
	// Check template builder
	if b.TB == nil {
		panic("Template builder is not set")
	}
	// Patch and cleanup existing jobs
	jobs := []scheduler.Job{}
	for _, job := range b.Scheduler.Jobs {
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
	b.Scheduler.Jobs = jobs
	// Add state population job
	b.Scheduler.Add(scheduler.Job{
		Group:   "populate",
		Depends: []string{"init"},
		Func: func() error {
			log.Println("populating state")
			for k, v := range params.State {
				b.State.Set(k, v)
			}
			return nil
		},
	})
	// Add action job
	b.Scheduler.Add(scheduler.Job{
		Group:   "action",
		Depends: []string{"populate"},
		Func: func() error {
			log.Println("calling an action")
			b.ActionMap[params.Action](params.Args...)
			return nil
		},
	})
	// Add final flush job
	b.Scheduler.Add(scheduler.Job{
		Group:   "flush",
		Depends: []string{"action"},
		Func: func() error {
			log.Println("final flush")
			Flush(b)
			return nil
		},
	})
}
