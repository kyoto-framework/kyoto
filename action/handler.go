package action

import (
	"net/http"

	"github.com/kyoto-framework/kyoto/v3/component"
)

// Handler is a http.HandlerFunc that have to be registered in your project
// for dynamic components rendering during action.
func Handler(w http.ResponseWriter, r *http.Request) {
	// Set headers
	w.Header().Set("Content-Type", "text/html")    // Sending component HTML
	w.Header().Set("Transfer-Encoding", "chunked") // Sending chunked data
	w.Header().Set("Cache-Control", "no-store")    // We don't need to cache this
	// Create context
	ctx := component.NewContext(w, r)
	// Prepare context for action
	context(ctx)
	// Get action
	action := ctx.Store.Get("Action").(Action)
	// Find component
	c, ok := components[action.Component]
	if !ok {
		panic("component not registered, probably you need to use action.Register")
	}
	// Build component state tree
	state := c(ctx)
	// If state doesn't have name, try to resolve it
	if state.GetName() == "" {
		state.SetName(c.GetName())
	}
	// Update action info
	action = ctx.Store.Get("Action").(Action)
	// Trigger flush, if not rendered
	if !action.rendered {
		Flush(ctx, state)
	}
}
