package render

import (
	"io"

	"github.com/kyoto-framework/kyoto/v3/component"
)

// Renderer defines requirements for rendering implementations.
type Renderer interface {
	// Render component into io.Writer.
	Render(state component.State, out io.Writer) error
}
