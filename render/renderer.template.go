package render

import (
	"embed"
	"html/template"
	"io"

	"go.kyoto.codes/v3/component"
)

// Template is a html/template renderer.
// Use Raw to provide handmade template,
// or provide template building parameters (Name, Glob, etc.).
type Template struct {
	Raw  *template.Template `json:"-"`
	Name string             // Resolved from component name by default

	Glob    string           `json:"-"` // *.html by default
	EmbedFS *embed.FS        `json:"-"` // nil by default
	FuncMap template.FuncMap `json:"-"` // render.FuncMap by default
}

func (t *Template) Render(state component.State, w io.Writer) error {
	// Defaults
	if t.Name == "" {
		t.Name = state.GetName()
	}
	if t.Glob == "" {
		t.Glob = "*.html"
	}
	if t.FuncMap == nil {
		t.FuncMap = FuncMapAll
	}
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
		} else {
			// Parse from disk
			tmpl = template.Must(tmpl.ParseGlob(t.Glob))
		}
	}
	// Render
	return tmpl.Execute(w, state)
}
