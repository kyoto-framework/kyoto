package kyoto

import (
	"html/template"
	"net/http"
)

// SSA aliases

// Action is the
type Action func(args ...interface{})

// ActionMap is a map of Actions stored by string
type ActionMap map[string]Action

// TemplateBuilder is used to build templates.
type TemplateBuilder func(p Page) *template.Template

// SSA page placeholder

type dummypage struct {
	TemplateBuilder TemplateBuilder
}

// Template returns a template
func (p *dummypage) Template() *template.Template {
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
type Page interface {
	// Template builder
	Template() *template.Template
}

// Component containts only must-have methods
type Component interface{}

// Extensions

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
