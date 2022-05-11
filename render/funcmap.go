package render

import (
	"html/template"

	"github.com/kyoto-framework/kyoto"
	"github.com/kyoto-framework/kyoto/helpers"
)

// ComposeFuncMap is a function for composing multiple FuncMap instances into one
func ComposeFuncMap(fmaps ...template.FuncMap) template.FuncMap {
	var result = template.FuncMap{}
	for _, fmap := range fmaps {
		for k, v := range fmap {
			result[k] = v
		}
	}
	return result
}

// FuncMap is responsible for integration of library functionality into template rendering
// You need to use this function while template building (or mix with your own)
func FuncMap(c *kyoto.Core) template.FuncMap {
	return template.FuncMap{
		// Component operations
		"componentid":    helpers.ComponentID,
		"componentname":  helpers.ComponentName,
		"componentattrs": ComponentAttrs,
		// Code injection
		"dynamics": Dynamics,
		// Rendering operations
		"render": func(state map[string]interface{}) template.HTML { return Render(c, state) },
		// Actions
		"action":     Action,
		"bind":       Bind,
		"formsubmit": FormSubmit,
	}
}
