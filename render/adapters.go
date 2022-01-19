package render

import (
	"html/template"
	"io"

	"github.com/kyoto-framework/kyoto"
)

// Template is a function to define a template for page.
func Template(b *kyoto.Core, builder func() *template.Template) {
	b.Context.Set("internal:render:tb", builder)
}

// Custom is a function to define a custom renderer for page.
func Custom(b *kyoto.Core, renderer func(io.Writer) error) {
	b.Context.Set("internal:render:cm", renderer)
}

// Redirect is a function to redirect.
func Redirect(b *kyoto.Core, target string, code int) {
	b.Context.Set("internal:render:redirect", target)
	b.Context.Set("internal:render:redirectCode", code)
}
