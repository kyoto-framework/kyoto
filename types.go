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
	Actions() ActionsMap
}

type ImplementsMeta interface {
	Meta() Meta
}
