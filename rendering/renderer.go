package rendering

import (
	"io"

	"go.kyoto.codes/v3/component"
)

// Renderer defines requirements for rendering implementations.
type Renderer interface {
	// Define if rendering must to be skipped.
	// Needed for cases like redirects.
	RenderSkip() bool
	// Render component into io.Writer.
	Render(state component.State, out io.Writer) error
}
