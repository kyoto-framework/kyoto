package ssc

import (
	"html/template"
)

// SSA aliases

type Action func(args ...interface{})
type ActionMap map[string]Action
type TemplateBuilder func() *template.Template

// SSA page placeholder

type dummypage struct {
	TemplateBuilder func() *template.Template
}

func (p *dummypage) Template() *template.Template {
	return p.TemplateBuilder()
}

// Meta

type Hreflang struct {
	Lang string
	Href string
}

type Meta struct {
	Title       string
	Description string
	Canonical   string
	Hreflangs   []Hreflang
	Additional  []map[string]string
}

// Basic page, only must-have methods
type Page interface {
	// Template builder
	Template() *template.Template
}

// Basic component, only must-have methods
type Component interface{}

// Extensions

type ImplementsInit interface {
	Init()
}

type ImplementsNestedInit interface {
	Init(Page)
}

type ImplementsAsync interface {
	Async() error
}

type ImplementsAfterAsync interface {
	AfterAsync()
}

type ImplementsActions interface {
	Actions() ActionMap
}

type ImplementsMeta interface {
	Meta() Meta
}
