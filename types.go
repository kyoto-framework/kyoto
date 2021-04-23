package ssc

import (
	"html/template"
)

// Aliases
type Action func(args ...interface{})
type ActionsMap map[string]Action

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

// Page initialization extension
type PageInit interface {
	Init()
}

// Page meta extension
type PageMeta interface {
	Meta() Meta
}

// Basic component, only must-have methods
type Component interface{}

// Component initialization extension
type ComponentInit interface {
	Init(Page)
}

// Component lifecycle extension
type ComponentAsync interface {
	Async() error
}

// Component lifecycle extension
type ComponentAfterAsync interface {
	AfterAsync()
}

// Component SSA extension
type ComponentActions interface {
	Actions() ActionsMap
}
