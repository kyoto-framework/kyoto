package ssc

import (
	"html/template"
)

// Aliases
type Action func(args map[string]interface{})
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

type Page interface {
	// Template builder
	Template() *template.Template
	// Entrypoint for registering components
	Init()
	// Meta info
	Meta() Meta
}

type Component interface {
	// Entrypoint for registering nested components
	Init()
	// Parts of component lifecycle
	Async() error
	AfterAsync()
	// Dynamic actions
	Actions() ActionsMap
}
