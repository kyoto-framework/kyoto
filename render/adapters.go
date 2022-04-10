package render

import (
	"html/template"
	"io"

	"github.com/kyoto-framework/kyoto"
)

// Template is a function to define a template builder for a page.
func Template(c *kyoto.Core, builder func() *template.Template) {
	// Save a template builder, just in case
	c.Context.Set("internal:render:tbuilder", builder)
	// Save prepared template for future optimizations
	c.Context.Set("internal:render:template", builder())
}

// Writer is a function to define a custom renderer for page.
func Writer(c *kyoto.Core, renderer func(io.Writer) error) {
	c.State.Set("internal:render:writer", renderer)
}

// Redirect is a function to redirect.
func Redirect(c *kyoto.Core, target string, code int) {
	c.Context.Set("internal:render:redirect", target)
	c.Context.Set("internal:render:redirectCode", code)
}
