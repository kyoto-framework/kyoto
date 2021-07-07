package ssc

import "html/template"

// Dummy page for component rendering
type DummyPage struct {
	TemplateBuilder func() *template.Template
	Component       Component
}

func (p *DummyPage) Template() *template.Template {
	return p.TemplateBuilder()
}

func (p *DummyPage) Init() {
	p.Component = RegC(p, p.Component)
}

func (*DummyPage) Meta() Meta {
	return Meta{}
}
