package render

import (
	"html/template"

	"github.com/kyoto-framework/kyoto"
)

func Template(b *kyoto.Core, builder func() *template.Template) {
	b.Context.Set("internal:render:tb", builder)
}
