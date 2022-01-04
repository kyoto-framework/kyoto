package smode

import "html/template"

type Action func(args ...interface{})
type ActionMap map[string]Action

type Page interface{}
type Component interface{}

type ImplementsTemplate interface {
	Template() *template.Template
}

type ImplementsTemplateWithPage interface {
	Template(p Page) *template.Template
}

type ImplementsInit interface {
	Init()
}

type ImplementsInitWithPage interface {
	Init(p Page)
}

type ImplementsAsync interface {
	Async() error
}

type ImplementsAsyncWithPage interface {
	Async(p Page) error
}

type ImplementsAfterAsync interface {
	AfterAsync()
}

type ImplementsAfterAsyncWithPage interface {
	AfterAsync(p Page)
}
