package htmx

import (
	"fmt"
	"go.kyoto.codes/v3/component"
	"html/template"
)

// FuncMap holds a library predefined template functions.
// You have to include it in your template building to use kyoto properly.
var FuncMap = template.FuncMap{
	// hxstate returns a hidden input with the state marshaled as a value.
	"hxstate": func(state any) template.HTML {
		_state := state.(component.State)
		return template.HTML(fmt.Sprintf(
			`<input type="hidden" name="hx-state" value="%s">`,
			_state.Marshal()))
	},
}
