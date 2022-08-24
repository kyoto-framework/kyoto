package kyoto

import (
	"html/template"
	"net/http"
)

// Context is the context of the current request.
// It is passed to the pages and components.
type Context struct {
	// Handler
	ResponseWriter http.ResponseWriter
	Request        *http.Request

	// Rendering
	Template     *template.Template
	TemplateConf *TemplateConfiguration

	// Action
	Action ActionParameters
}
