package render

import "html/template"

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
func FuncMap() template.FuncMap {
	return template.FuncMap{
		"dynamics":       Dynamics,
		"componentattrs": ComponentAttrs,
		"action":         Action,
		"bind":           Bind,
		"formsubmit":     FormSubmit,
	}
}
