package main

import (
	"html/template"

	"github.com/kyoto-framework/kyoto"
	"github.com/kyoto-framework/kyoto/lifecycle"
	"github.com/kyoto-framework/kyoto/render"
)

func PageIndex(c *kyoto.Core) {
	render.Template(c, func() *template.Template {
		return template.Must(template.New("page.index.html").Funcs(render.FuncMap()).ParseGlob("*.html"))
	})
	lifecycle.Init(c, func() {
		c.State.Set("Title", "Kyoto in a functional way")
		c.Component("UUID1", ComponentUUID("First UUID"))
		c.Component("UUID2", ComponentUUID("Second UUID"))
	})
}
