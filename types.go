package kyoto

import (
	"html/template"
	"net/http"
)

// SSA aliases

type Action func(args ...interface{})
type ActionMap map[string]Action
type TemplateBuilder func(p Page) *template.Template

// SSA page placeholder

type dummypage struct {
	TemplateBuilder TemplateBuilder
}

func (p *dummypage) Template() *template.Template {
	return p.TemplateBuilder(p)
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

// RedirectParameters is a helper for Redirect function
type RedirectParameters struct {
	Page              Page
	ResponseWriter    http.ResponseWriter
	Request           *http.Request
	ResponseWriterKey string
	RequestKey        string
	Target            string
	StatusCode        int
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
	Init(Page)
}

type ImplementsInitWithoutPage interface {
	Init()
}

type ImplementsAsync interface {
	Async(Page) error
}

type ImplementsAsyncWithoutPage interface {
	Async() error
}

type ImplementsAfterAsync interface {
	AfterAsync(Page)
}

type ImplementsAfterAsyncWithoutPage interface {
	AfterAsync()
}

type ImplementsActions interface {
	Actions(Page) ActionMap
}

type ImplementsActionsWithoutPage interface {
	Actions() ActionMap
}

type ImplementsMeta interface {
	Meta() Meta
}
