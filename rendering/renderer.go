package rendering

import (
	"io"

	"go.kyoto.codes/v3/component"
)

// Renderer defines requirements for rendering implementations.
type Renderer interface {
	// Render component into io.Writer.
	Render(state component.State, out io.Writer) error
}
