package component

import (
	"html/template"
)

// FuncMap holds a library predefined template functions.
// You have to include it in your template building to use kyoto properly.
var FuncMap = template.FuncMap{
	// Marshal allows to marshal state to string.
	"marshal": func(state State) string {
		return state.Marshal(state)
	},
}
