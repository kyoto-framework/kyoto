package main

import (
	"html/template"

	"github.com/kyoto-framework/kyoto/render"
	"github.com/kyoto-framework/kyoto/smode"
)

type PageSMode struct {
	Title string
	UUID1 smode.Component
	UUID2 smode.Component
}

func (p *PageSMode) Template() *template.Template {
	return template.Must(template.New("page.smode.html").Funcs(render.FuncMap()).ParseGlob("*.html"))
}

func (p *PageSMode) Init() {
	p.Title = "Kyoto in a structure way"
	p.UUID1 = smode.RegC(p, ComponentUUID("First UUID"))
	p.UUID2 = smode.RegC(p, &ComponentUUIDStruct{Title: "Second UUID"})
}
