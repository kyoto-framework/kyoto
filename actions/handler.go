package actions

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/kyoto-framework/kyoto"
)

// Handler is a generic actions handler, responsible for component rendering
// on action call. Please note, you also need to register your dynamic components with Register method.
func Handler(tb func(c *kyoto.Core) *template.Template) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Set headers
		rw.Header().Set("Content-Type", "plain/html")
		rw.Header().Set("Transfer-Encoding", "chunked")
		rw.Header().Set("Cache-Control", "no-store")
		// Extract SSA parameters
		parameters, err := ParseParameters(r)
		if err != nil {
			panic(err)
		}
		// Prepare core
		core := kyoto.NewCore()
		// Inject component name into state
		core.State.Set("internal:name", parameters.Component)
		// Set context
		core.Context.Set("internal:rw", rw)
		core.Context.Set("internal:r", r)
		core.Context.Set("internal:render:p", &parameters)
		// Find component and apply to core
		registryrw.RLock()
		if component, found := registry[parameters.Component]; found {
			component(core)
		} else {
			panic("Component not found in registry")
		}
		registryrw.RUnlock()
		// If no custom render, set template builder
		if core.State.Get("internal:render:writer") == nil {
			tbc := func() *template.Template {
				return template.Must(tb(core).Parse(fmt.Sprintf(`{{ template "%s" . }}`, parameters.Component)))
			}
			core.Context.Set("internal:render:tbuilder", tbc)
			core.Context.Set("internal:render:template", tbc())
		}
		// Patch scheduler with state population, action and flush jobs
		Patch(core, parameters)
		// Execute
		core.Execute()
	}
}
