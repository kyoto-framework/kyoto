package actions

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/kyoto-framework/kyoto"
)

func Handler(tb func() *template.Template) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Set server-sent events headers
		rw.Header().Set("Content-Type", "text/event-stream")
		rw.Header().Set("Cache-Control", "no-cache")
		rw.Header().Set("Connection", "keep-alive")
		// Extract SSA parameters
		params, err := ParseParameters(r.URL.Path)
		if err != nil {
			panic(err)
		}
		// Prepare core
		core := kyoto.NewCore()
		// Inject component name into state
		core.State.Set("internal:name", params.Component)
		// Set context
		core.Context.Set("internal:rw", rw)
		core.Context.Set("internal:r", r)
		core.Context.Set("internal:render:p", &params)
		// Set template builder
		core.Context.Set("internal:render:tb", func() *template.Template {
			return template.Must(tb().Parse(fmt.Sprintf(`{{ template "%s" . }}`, params.Component)))
		})
		// Find component and apply to core
		registryrw.RLock()
		if component, found := registry[params.Component]; found {
			// Switch behavior based on type of the component
			switch component := component.(type) {
			case func(*kyoto.Core): // Function component
				component(core)
			default: // Not supported
				panic("Component type is not supported")
			}
		} else {
			panic("Component not found in registry")
		}
		registryrw.RUnlock()
		// Patch scheduler with state population, action and flush jobs
		Patch(core, params)
		// Execute
		core.Execute()
	}
}
