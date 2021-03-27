package gofr

import (
	"html/template"
)

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
	// Meta info
	Meta() Meta
	// Entrypoint for registering components
	Init()
}

type Component interface {
	// Parts of component lifecycle
	Async() error
	AfterAsync()
	// Entrypoint for registering nested components
	Init()
}
