package rendering

import (
	"embed"
	"html/template"
	"io"

	"go.kyoto.codes/v3/component"
)

// Global template configuration defaults.
// We're providing them to make it easier to configure rendering defaults across the project.
var (
	TEMPLATE_GLOB              = "*.html"
	TEMPLATE_FUNCMAP           = FuncMapAll
	TEMPLATE_EMBEDFS *embed.FS = nil
)

// Template is a html/template renderer.
// Use Raw to provide handmade template,
// or provide template building parameters (Name, Glob, etc.).
type Template struct {
	Raw  *template.Template `json:"-"` // Raw template will be used instead if provided
	Name string             // Resolved from component name by default
	Skip bool               `json:"-"` // false by default

	Glob    string           `json:"-"` // *.html by default
	EmbedFS *embed.FS        `json:"-"` // nil by default
	FuncMap template.FuncMap `json:"-"` // render.FuncMap by default
}

func (t *Template) RenderSkip() bool {
	return t.Skip
}

func (t *Template) Render(state component.State, w io.Writer) error {
	// Defaults
	if t.Name == "" {
		t.Name = state.GetName()
	}
	if t.Glob == "" {
		t.Glob = TEMPLATE_GLOB
	}
	if t.FuncMap == nil {
		t.FuncMap = TEMPLATE_FUNCMAP
	}
	// Define template
	tmpl := t.Raw
	if tmpl == nil {
		// Base
		tmpl = template.New(t.Name)
		// Functions
		tmpl = tmpl.Funcs(t.FuncMap)
		// Parse
		if t.EmbedFS != nil {
			// Parse embedded
			tmpl = template.Must(tmpl.ParseFS(t.EmbedFS, t.Glob))
		} else if TEMPLATE_EMBEDFS != nil {
			// Parse embedded
			tmpl = template.Must(tmpl.ParseFS(TEMPLATE_EMBEDFS, t.Glob))
		} else {
			// Parse from disk
			tmpl = template.Must(tmpl.ParseGlob(t.Glob))
		}
	}
	// Render
	return tmpl.Execute(w, state)
}
