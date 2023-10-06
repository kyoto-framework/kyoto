package component

import (
	"fmt"
	"html/template"
)

// FuncMap holds a library predefined template functions.
// You have to include it in your template building to use kyoto properly.
var FuncMap = template.FuncMap{
	// Marshal allows to marshal state to string.
	"marshal": func(state State) string {
		return state.Marshal()
	},
	// Component composes HTML attributes for component state and identification.
	// Allows to avoid explicit component declaration syntax.
	"component": func(state State) template.HTMLAttr {
		return template.HTMLAttr(
			fmt.Sprintf(`component="%s" state="%s"`,
				state.GetName(),
				state.Marshal()),
		)
	},
}
