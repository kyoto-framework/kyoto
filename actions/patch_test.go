package actions

import (
	"html/template"
	"testing"

	"github.com/kyoto-framework/kyoto"
	"github.com/kyoto-framework/kyoto/lifecycle"
)

func testPatchComponent(c *kyoto.Core) {
	lifecycle.Init(c, func() {
		c.State.Set("Foo", "Bar")
	})
	lifecycle.Async(c, func() error {
		c.State.Set("Foo", "Baz")
		return nil
	})
	Define(c, "Bar", func(args ...interface{}) {
		c.State.Set("Foo", "Bar")
	})
}

func TestPatch(t *testing.T) {
	// Initialize core
	core := kyoto.NewCore()
	// Apply component
	testPatchComponent(core)
	// Define render
	core.Context.Set("internal:render:tb", func() *template.Template {
		return template.Must(template.New("").Parse(`
			{{ define "testPatchComponent" }}
			...
			{{ end }}

			{{ template "testPatchComponent" . }}
		`))
	})
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
		if job.Name != "placeholder" &&
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
