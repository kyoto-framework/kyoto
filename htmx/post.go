package htmx

import "go.kyoto.codes/v3/component"

// Post is a helper function that simplifies the handling of stateful htmx POST requests.
func Post(ctx *component.Context, state component.State, handler func()) {
	// We are only interested in POST requests here
	if ctx.Request.Method == "POST" {
		// Parse the form to get the state
		ctx.Request.ParseForm()
		// If no state is present in the form, we ignore.
		// Porbably this is a regular POST request, not related to htmx.
		if ctx.Request.FormValue("hx-state") == "" {
			return
		}
		// If the state is disposable, we panic.
		// This is a safety measure to prevent misuse of disposable components.
		if ctx.Request.FormValue("hx-state") == "disposable" {
			panic("incorrect use of disposable component")
		}
		// Unmarshal the state from the form
		state.Unmarshal(state, ctx.Request.FormValue("hx-state"))
		// Call the handler
		handler()
	}
}
