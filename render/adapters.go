package render

import (
	"html/template"
	"io"

	"github.com/kyoto-framework/kyoto"
)

func Template(b *kyoto.Core, builder func() *template.Template) {
	b.Context.Set("internal:render:tb", builder)
}

func Custom(b *kyoto.Core, renderer func(io.Writer) error) {
	b.Context.Set("internal:render:cm", renderer)
}

func Redirect(b *kyoto.Core, target string, code int) {
	b.Context.Set("internal:render:redirect", target)
	b.Context.Set("internal:render:redirectCode", code)
}
