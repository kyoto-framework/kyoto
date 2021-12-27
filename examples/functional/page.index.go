package main

import (
	"html/template"

	"github.com/kyoto-framework/kyoto"
	"github.com/kyoto-framework/kyoto/templates"
)

func PageIndex(b *kyoto.Builder) {
	b.Template(func() *template.Template {
		return template.Must(template.New("page.index.html").Funcs(templates.FuncMap()).ParseGlob("*.html"))
	})
	b.Init(func() {
		b.State.Set("Title", "Kyoto in a functional way")
		b.Component("UUID1", ComponentUUID("First UUID"))
		b.Component("UUID2", ComponentUUID("Second UUID"))
	})
}
