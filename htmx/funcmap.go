package htmx

import (
	"fmt"
	"go.kyoto.codes/v3/component"
	"html/template"
)

var FuncMap = template.FuncMap{
	"hxstate": func(state any) template.HTML {
		_state := state.(component.State)
		return template.HTML(fmt.Sprintf(
			`<input type="hidden" name="hx-state" value="%s">`,
			_state.Marshal()))
	},
}
