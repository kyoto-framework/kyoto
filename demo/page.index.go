package main

import (
	"html/template"

	"github.com/yuriizinets/go-ssc"
)

type PageIndex struct {
	ComponentHttpbinUUID   ssc.Component
	ComponentCounter       ssc.Component
	ComponentSampleBinding ssc.Component
	ComponentSampleParent  ssc.Component
}

func (*PageIndex) Template() *template.Template {
	return template.Must(template.New("page.index.html").Funcs(funcmap()).ParseGlob("*.html"))
}

func (p *PageIndex) Init() {
	p.ComponentHttpbinUUID = ssc.RegC(p, &ComponentHttpbinUUID{})
	p.ComponentCounter = ssc.RegC(p, &ComponentCounter{})
	p.ComponentSampleBinding = ssc.RegC(p, &ComponentSampleBinding{})
	p.ComponentSampleParent = ssc.RegC(p, &ComponentSampleParent{})
}

func (*PageIndex) Meta() ssc.Meta {
	return ssc.Meta{
		Title: "SSC Example",
	}
}
