package render

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/kyoto-framework/kyoto"
	"github.com/kyoto-framework/scheduler"
)

func PageHandler(page func(*kyoto.Core)) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Initialize the core
		core := kyoto.NewCore()
		// Set core context
		core.Context.Set("internal:rw", rw)
		core.Context.Set("internal:r", r)
		// Call core receiver
		page(core)
		// Check if page implements template
		if core.Context.Get("internal:render:tb") == nil {
			panic("No template specified for page")
		}
		// Collect all job groups to use them as dependencies
		groups := []string{}
		for _, job := range core.Scheduler.Jobs {
			// Avoid cycle dependencies
			if !strings.Contains(strings.Join(job.Depends, ","), "render") {
				groups = append(groups, job.Group)
			}
		}
		// Schedule a render job
		core.Scheduler.Add(&scheduler.Job{
			Group:   "render",
			Depends: groups,
			Func: func() error {
				tb := core.Context.Get("internal:render:tb").(func() *template.Template)
				return tb().Execute(rw, core.State.Export())
			},
		})
		// Execute scheduler
		core.Execute()
	}
}
