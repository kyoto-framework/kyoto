package templates

import "html/template"

// FuncMap is responsible for integration of library functionality into template rendering
// You need to use this function while template building (or mix with your own)
func FuncMap() template.FuncMap {
	return template.FuncMap{
		"meta":           Meta,
		"dynamics":       Dynamics,
		"json":           JSON,
		"componentattrs": ComponentAttrs,
		"action":         Action,
		"bind":           Bind,
		"formsubmit":     FormSubmit,
	}
}
