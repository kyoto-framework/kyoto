package rendering

import (
	"net/http"

	"go.kyoto.codes/v3/component"
)

// Handler builds a http.HandlerFunc that renders provided component.
func Handler(c component.Component) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create context
		ctx := component.NewContext(w, r)
		// Build page state tree
		state := c(ctx)
		// Inject component name, unless it's already set
		if state.GetName() == "" {
			state.SetName(c.GetName())
		}
		// Ensure state implements render
		if _, ok := state.(Renderer); !ok {
			panic("The component does not implement rendering")
		}
		// Check if we need to skip rendering
		if state.(Renderer).RenderSkip() {
			return
		}
		// Render
		if err := state.(Renderer).Render(state, ctx.ResponseWriter); err != nil {
			panic(err)
		}
	}
}
