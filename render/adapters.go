package render

import (
	"html/template"
	"net/http"

	"github.com/kyoto-framework/kyoto"
)

func Template(b *kyoto.Core, builder func() *template.Template) {
	b.Context.Set("internal:render:tb", builder)
}

func Renderer(b *kyoto.Core, renderer func(rw http.ResponseWriter) error) {
	b.Context.Set("internal:render:rnd", renderer)
}

func Redirect(b *kyoto.Core, url string, code int) {
	b.Context.Set("internal:render:redirect", url)
	b.Context.Set("internal:render:redirectCode", code)
}
