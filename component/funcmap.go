package component

import (
	"fmt"
	"html/template"
	"strings"
)

// FuncMap holds a library predefined template functions.
// You have to include it in your template building to use kyoto properly.
var FuncMap = template.FuncMap{
	// Marshal allows to marshal state to string.
	"marshal": func(state State) string {
		return state.Marshal(state)
	},
	// Component composes HTML attributes for component state and identification.
	// Allows to avoid explicit component declaration syntax.
	"component": func(state State) template.HTMLAttr {
		builder := strings.Builder{}
		// Add state
		builder.WriteString(fmt.Sprintf(`state="%s"`, state.Marshal(state)))
		// Add name
		if state, ok := state.(interface{ GetName() string }); ok {
			builder.WriteString(fmt.Sprintf(` name="%s"`, state.GetName()))
		}
		// Return
		return template.HTMLAttr(builder.String())
	},
}
