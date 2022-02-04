package render

import (
	"errors"
	"html/template"
	"io"
	"net/http"

	"github.com/kyoto-framework/kyoto"
	"github.com/kyoto-framework/scheduler"
)

// PageHandler is a function, responsible for page rendering.
// Returns http.HandlerFunc for registering in your router.
func PageHandler(page func(*kyoto.Core)) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Initialize the core
		core := kyoto.NewCore()
		// Set core context
		core.Context.Set("internal:rw", rw)
		core.Context.Set("internal:r", r)
		// Call core receiver
		page(core)
		// Check if page implements render
		if core.Context.Get("internal:render:tb") == nil && core.Context.Get("internal:render:cm") == nil {
			panic("Rendering is not specified for page")
		}
		// Schedule a render job
		core.Scheduler.Add(&scheduler.Job{
			Group: "render",
			Func: func() error {
				if redirect := core.Context.Get("internal:render:redirect"); redirect != nil { // Check redirect
					code := 302
					if redirectCode := core.Context.Get("internal:render:redirectCode"); redirectCode != nil {
						code = redirectCode.(int)
					}
					http.Redirect(rw, r, redirect.(string), code)
					return nil
				}
				if renderer := core.Context.Get("internal:render:cm"); renderer != nil { // Check renderer
					return renderer.(func(io.Writer) error)(rw) // Call renderer
				} else if tbuilder := core.Context.Get("internal:render:tb"); tbuilder != nil { // Check template builder
					return tbuilder.(func() *template.Template)().Execute(rw, core.State.Export()) // Execute template
				}
				return errors.New("no renderer or template builder specified") // Error if no renderer or template builder
			},
		})
		// Execute scheduler
		core.Execute()
	}
}
