package kyoto

import (
	"html/template"
	"net/http"
)

// SSA aliases

// Deprecated: Action is used to pass arguments onto a component to be used
type Action func(args ...interface{})

// Deprecated: ActionMap is a map of Actions stored by string
type ActionMap map[string]Action

// Deprecated: TemplateBuilder is used to build templates.
type TemplateBuilder func(p Page) *template.Template

// Deprecated: DummyPage is an SSA page placeholder
type DummyPage struct {
	TemplateBuilder TemplateBuilder
}

// Deprecated: Template returns a template
func (p *DummyPage) Template() *template.Template {
	return p.TemplateBuilder(p)
}

// Meta

// Deprecated: Hreflang stores the lang and href
type Hreflang struct {
	Lang string
	Href string
}

// Deprecated: Meta data
type Meta struct {
	Title       string
	Description string
	Canonical   string
	Hreflangs   []Hreflang
	Additional  []map[string]string
}

// Deprecated RedirectParameters is a helper for the Redirect function
type RedirectParameters struct {
	Page              Page
	ResponseWriter    http.ResponseWriter
	Request           *http.Request
	ResponseWriterKey string
	RequestKey        string
	Target            string
	StatusCode        int
}

// Deprecated: Page contains only must-have methods
type Page interface{}

// Deprecated: Component contains only must-have methods
type Component interface{}

// Extensions

// Deprecated
type ImplementsTemplate interface {
	Template(Page) *template.Template
}

// Deprecated
type ImplementsTemplateWithoutPage interface {
	Template() *template.Template
}

// Deprecated
type ImplementsRender interface {
	Render() string
}

// Deprecated: ImplementsInit interface for implementing init
type ImplementsInit interface {
	Init(Page)
}

// Deprecated: ImplementsInitWithoutPage interface type
type ImplementsInitWithoutPage interface {
	Init()
}

// Deprecated: ImplementsAsync interface for asynchronous functions
type ImplementsAsync interface {
	Async(Page) error
}

// Deprecated: ImplementsAsyncWithoutPage interface for Asynchronous functions without a page
type ImplementsAsyncWithoutPage interface {
	Async() error
}

// Deprecated: ImplementsAfterAsync interface for after asynchronous functions
type ImplementsAfterAsync interface {
	AfterAsync(Page)
}

// Deprecated: ImplementsAfterAsyncWithoutPage interface type
type ImplementsAfterAsyncWithoutPage interface {
	AfterAsync()
}

// Deprecated: ImplementsActions interface for implementing actions
type ImplementsActions interface {
	Actions(Page) ActionMap
}

// Deprecated: ImplementsActionsWithoutPage interface type
type ImplementsActionsWithoutPage interface {
	Actions() ActionMap
}

// Deprecated: ImplementsMeta interface type
type ImplementsMeta interface {
	Meta() Meta
}
