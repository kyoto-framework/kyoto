package main

import (
	"html/template"

	"github.com/yuriizinets/go-ssc"
)

type PageIndex struct {
	ComponentHttpbinUUID   ssc.Component
	ComponentCounter       ssc.Component
	ComponentSampleBinding ssc.Component
}

func (*PageIndex) Template() *template.Template {
	fmap := funcmap()
	tmpl, _ := template.New("page.index.html").Funcs(fmap).ParseGlob("*.html")
	return tmpl
}

func (p *PageIndex) Init() {
	p.ComponentHttpbinUUID = ssc.RegC(p, &ComponentHttpbinUUID{})
	p.ComponentCounter = ssc.RegC(p, &ComponentCounter{})
	p.ComponentSampleBinding = ssc.RegC(p, &ComponentSampleBinding{})
}

func (*PageIndex) Meta() ssc.Meta {
	return ssc.Meta{
		Title: "SSC Example",
	}
}
