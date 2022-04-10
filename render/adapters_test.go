package render

import (
	"html/template"
	"io"
	"testing"

	"github.com/kyoto-framework/kyoto"
)

// TestRender ensures template integrating into core as expected
func TestTemplate(t *testing.T) {
	// Initialize core
	core := kyoto.NewCore()

	// Set template builder
	Template(core, func() *template.Template {
		return template.Must(template.New("").Parse("{{.}}"))
	})

	// Check template builder
	if core.Context.Get("internal:render:tb") == nil {
		t.Error("Template builder is not set")
	}
}

// TestWriter ensures writer renderer integrating into core as expected
func TestWriter(t *testing.T) {
	// Initialize core
	core := kyoto.NewCore()

	// Set writer renderer
	Writer(core, func(w io.Writer) error {
		w.Write([]byte("test"))
		return nil
	})

	// Check custom renderer
	if core.Context.Get("internal:render:cm") == nil {
		t.Error("Custom renderer is not set")
	}
}

// TestRedirect ensures redirect integrating into core as expected
func TestRedirect(t *testing.T) {
	// Initialize core
	core := kyoto.NewCore()

	// Set redirect
	Redirect(core, "/", 302)

	// Check redirect
	if core.Context.Get("internal:render:redirect") == nil {
		t.Error("Redirect is not set")
	}
	if core.Context.Get("internal:render:redirectCode") == nil {
		t.Error("Redirect code is not set")
	}
}
