package action

import (
	"fmt"
	"net/http"

	"go.kyoto.codes/v3/component"
)

// Redirect is a function to trigger redirect during action handling.
func Redirect(ctx *component.Context, href string, status int) {
	// Initialize flusher
	flusher := ctx.ResponseWriter.(http.Flusher)
	// Compose command
	cmd := fmt.Sprintf("action:redirect=%s", href)
	// Append terminator sequence and write to stream
	if _, err := fmt.Fprint(ctx.ResponseWriter, cmd+string(terminator)); err != nil {
		panic(err)
	}
	// Set redirected flag
	action := ctx.Get("Action").(Action)
	action.rendered = true
	ctx.Set("Action", action)
	// Flush
	flusher.Flush()
}
