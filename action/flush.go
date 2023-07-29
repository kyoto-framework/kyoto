package action

import (
	"bytes"
	"io"
	"net/http"

	"go.kyoto.codes/v3/component"
)

// We're duplicating this from render package to avoid cycle import.
// There's probably a better solution to this.
type renderer interface {
	Render(state component.State, out io.Writer) error
}

// Flush renders a component UI and returns it to the client side.
// Call it when you need to push an updated component markup to the client.
// Might be called multiple times during single action call.
func Flush(ctx *component.Context, state component.State) {
	// Get action
	action := ctx.Store.Get("Action").(Action)
	// Exit if rendering is already done
	if action.rendered {
		return
	}
	// Ensure state implements render
	if _, ok := state.(renderer); !ok {
		panic("component does not implement rendering")
	}
	// Initialize flusher
	flusher := ctx.ResponseWriter.(http.Flusher)
	// Render into buffer
	buf := &bytes.Buffer{}
	if err := state.(renderer).Render(state, buf); err != nil {
		panic(err)
	}
	// Append terminator sequence.
	// We are using buffer and terminator sequence to ensure integrity.
	// Sometimes, chunk becomes split which leads to broken render.
	// This is a workaround solution, actual reason of such behavior wasn't found.
	buf.Write(terminator)
	// Write buffer to response
	if _, err := buf.WriteTo(ctx.ResponseWriter); err != nil {
		panic(err)
	}
	// Flush
	flusher.Flush()
}
