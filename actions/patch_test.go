package actions

import (
	"html/template"
	"testing"

	"github.com/kyoto-framework/kyoto"
	"github.com/kyoto-framework/scheduler"
)

func TestPatch(t *testing.T) {
	// Initialize core
	core := kyoto.NewCore()
	// Define render
	core.Context.Set("internal:render:tb", func() *template.Template {
		return template.Must(template.New("").Parse(`TestPatch`))
	})
	// Add few jobs
	core.Scheduler.Add(&scheduler.Job{
		Group: "init",
		Func:  func() error { return nil },
	})
	core.Scheduler.Add(&scheduler.Job{
		Group: "async",
		Func:  func() error { return nil },
	})
	core.Scheduler.Add(&scheduler.Job{
		Group: "afterasync",
		Func:  func() error { return nil },
	})
	// Define an action
	Define(core, "Bar", func(args ...interface{}) {})
	// Patch
	Patch(core, Parameters{
		Component: "Foo",
		Action:    "Bar",
		State: map[string]interface{}{
			"Foo": "Bar",
		},
		Args: []interface{}{},
	})
	// Check scheduler is cleared
	for _, job := range core.Scheduler.Jobs {
		if job.Group != "placeholder" &&
			job.Group != "init" &&
			job.Group != "populate" &&
			job.Group != "state" &&
			job.Group != "action" &&
			job.Group != "render" {
			t.Error("Job group " + job.Group + " is not cleared")
		}
	}
	// Check scheduler has action and render job
	flagaction := false
	flagrender := false
	for _, job := range core.Scheduler.Jobs {
		if job.Group == "action" {
			flagaction = true
		}
		if job.Group == "render" {
			flagrender = true
		}
	}
	if !flagaction {
		t.Error("Action job is not added")
	}
	if !flagrender {
		t.Error("Render job is not added")
	}
}
