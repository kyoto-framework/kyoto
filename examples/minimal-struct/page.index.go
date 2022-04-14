package main

import (
	"html/template"

	"github.com/kyoto-framework/kyoto/smode"
)

type PageIndex struct {
	Title     string
	UUID1     smode.Component
	UUID2     smode.Component
	UserAgent smode.Component
}

func (p *PageIndex) Template() *template.Template {
	return template.Must(template.New("page.index.html").Funcs(smode.FuncMap(p)).ParseGlob("*.html"))
}

func (p *PageIndex) Init() {
	p.Title = "Kyoto in a struct way"
	p.UUID1 = smode.UseC(p, &ComponentUUID{
		Title: "First UUID",
	})
	p.UUID2 = smode.UseC(p, &ComponentUUID{
		Title: "Second UUID",
	})
	p.UserAgent = smode.UseC(p, &ComponentUserAgent{})
}
