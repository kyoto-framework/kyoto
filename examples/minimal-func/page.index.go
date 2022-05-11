package main

import (
	"html/template"

	"github.com/kyoto-framework/kyoto"
	"github.com/kyoto-framework/kyoto/lifecycle"
	"github.com/kyoto-framework/kyoto/render"
	"github.com/kyoto-framework/kyoto/state"
)

func PageIndex(c *kyoto.Core) {
	// Define state
	state.New(c, "Title", "Kyoto in a functional way")

	// Define lifecycle
	lifecycle.Init(c, func() {
		c.State.Set("Title", "Kyoto in a functional way")
		c.Component("UUID1", ComponentUUID("First UUID"))
		c.Component("UUID2", ComponentUUID("Second UUID"))
		c.Component("UserAgent", ComponentUserAgent)
	})

	// Define rendering
	render.Template(c, func() *template.Template {
		return template.Must(template.New("page.index.html").Funcs(render.FuncMap(c)).ParseGlob("*.html"))
	})
}
