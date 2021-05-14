package ssc

import "html/template"

// Dummy page for component rendering
type DummyPage struct {
	DTemplate  *template.Template
	DComponent Component
}

func (p *DummyPage) Template() *template.Template {
	return p.DTemplate
}

func (p *DummyPage) Init() {
	p.DComponent = RegC(p, p.DComponent)
}

func (*DummyPage) Meta() Meta {
	return Meta{}
}
