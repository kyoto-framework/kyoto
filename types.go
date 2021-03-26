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
	Template() *template.Template
	Meta() Meta
	Init()
}

type Component interface {
	Definition() string
	Async() error
	AfterAsync()
}
