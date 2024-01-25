package htmx

import "go.kyoto.codes/v3/component"

func StateOf[T component.State](ctx *component.Context) T {
	// Extract encoded state
	encoded := ctx.Request.FormValue("hx-state")
	if encoded == "" {
		encoded = ctx.Request.URL.Query().Get("hx-state")
	}
	// Decode state
	var state T
	state.Unmarshal(encoded)
	// Return
	return state
}
