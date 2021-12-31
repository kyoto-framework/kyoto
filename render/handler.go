package render

import (
	"html/template"
	"net/http"

	"github.com/kyoto-framework/kyoto"
	"github.com/kyoto-framework/scheduler"
)

func PageHandler(page func(*kyoto.Core)) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		core := kyoto.NewCore()
		core.Context.Set("internal:rw", rw)
		core.Context.Set("internal:r", r)
		page(core)
		if core.Context.Get("internal:render:tb") == nil {
			panic("No template specified for page")
		}
		groups := []string{}
		for _, job := range core.Scheduler.Jobs {
			groups = append(groups, job.Group)
		}
		core.Scheduler.Add(scheduler.Job{
			Group:   "render",
			Depends: groups,
			Func: func() error {
				tb := core.Context.Get("internal:render:tb").(func() *template.Template)
				return tb().Execute(rw, core.State.Export())
			},
		})
		core.Execute()
	}
}
