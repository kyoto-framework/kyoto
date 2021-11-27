package kyoto

import (
	"html/template"
	"net/http"
)

// SSA aliases

// Action is used to pass arguments onto a component to be used
type Action func(args ...interface{})

// ActionMap is a map of Actions stored by string
type ActionMap map[string]Action

// TemplateBuilder is used to build templates.
type TemplateBuilder func(p Page) *template.Template

// DummyPage is an SSA page placeholder
type DummyPage struct {
	TemplateBuilder TemplateBuilder
}

// Template returns a template
func (p *DummyPage) Template() *template.Template {
	return p.TemplateBuilder(p)
}

// Meta

// Hreflang stores the lang and href
type Hreflang struct {
	Lang string
	Href string
}

// Meta data
type Meta struct {
	Title       string
	Description string
	Canonical   string
	Hreflangs   []Hreflang
	Additional  []map[string]string
}

// RedirectParameters is a helper for the Redirect function
type RedirectParameters struct {
	Page              Page
	ResponseWriter    http.ResponseWriter
	Request           *http.Request
	ResponseWriterKey string
	RequestKey        string
	Target            string
	StatusCode        int
}

// Page contains only must-have methods
type Page interface{}

// Component containts only must-have methods
type Component interface{}

// Extensions

type ImplementsTemplate interface {
	Template(Page) *template.Template
}

type ImplementsTemplateWithoutPage interface {
	Template() *template.Template
}

type ImplementsRender interface {
	Render() string
}

// ImplementsInit interface for implementing init
type ImplementsInit interface {
	Init(Page)
}

// ImplementsInitWithoutPage interface type
type ImplementsInitWithoutPage interface {
	Init()
}

// ImplementsAsync interface for asynchronous functions
type ImplementsAsync interface {
	Async(Page) error
}

// ImplementsAsyncWithoutPage interface for Asynchronous functions without a page
type ImplementsAsyncWithoutPage interface {
	Async() error
}

// ImplementsAfterAsync interface for after asynchronous functions
type ImplementsAfterAsync interface {
	AfterAsync(Page)
}

// ImplementsAfterAsyncWithoutPage interface type
type ImplementsAfterAsyncWithoutPage interface {
	AfterAsync()
}

// ImplementsActions interface for implementing actions
type ImplementsActions interface {
	Actions(Page) ActionMap
}

// ImplementsActionsWithoutPage interface type
type ImplementsActionsWithoutPage interface {
	Actions() ActionMap
}

// ImplementsMeta interface type
type ImplementsMeta interface {
	Meta() Meta
}
