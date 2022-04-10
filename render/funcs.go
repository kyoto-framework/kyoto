package render

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"strings"

	"github.com/kyoto-framework/kyoto"
	"github.com/kyoto-framework/kyoto/actions"
	"github.com/kyoto-framework/kyoto/helpers"
)

// Render is a function to render a component.
// TODO: Not implemented.
func Render(c *kyoto.Core, state map[string]interface{}) template.HTML {
	// Define buffer
	buf := &bytes.Buffer{}
	// Check if state have a renderer
	if renderer, ok := state["internal:render:wr"]; ok {
		err := renderer.(func(io.Writer) error)(buf)
		if err != nil {
			panic(err)
		}
	} else { // Render with template
		tbuilder := c.Context.Get("internal:render:tb").(func() *template.Template)
		tmpl := template.Must(tbuilder().Parse(fmt.Sprintf(`{{ template "%s" . }}`, helpers.ComponentName(state))))
		err := tmpl.Execute(buf, state)
		if err != nil {
			panic(err)
		}
	}
	// Return wrapped buffer
	return template.HTML(buf.String())
}

// Dynamics is a function to integrate dynamic kyoto functionality (actions).
func Dynamics(path ...string) template.HTML {
	if len(path) == 0 {
		path = append(path, "/internal/actions/")
	}
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("<script>const ssapath = \"%s\"</script>", path[0]))
	builder.WriteString(actions.Client)
	return template.HTML(builder.String())
}

// ComponentAttrs is a function to serialize and inject component data into page.
func ComponentAttrs(component interface{}) template.HTMLAttr {
	return template.HTMLAttr(fmt.Sprintf(
		`cid='%s' name='%s' state='%s'`,
		helpers.ComponentID(component),
		helpers.ComponentName(component),
		helpers.ComponentSerialize(component),
	))
}

// Action is a wrapper around JS function for calling server side actions.
func Action(action string, args ...interface{}) template.JS {
	var formattedargs []string
	for _, arg := range args {
		b, _ := json.Marshal(arg)
		formattedargs = append(formattedargs, string(b))
	}

	return template.JS(fmt.Sprintf("Action(this, '%s', %s)", action, strings.Join(formattedargs, ", ")))
}

// Bind is a wrapper around JS function to bind input to the component state.
func Bind(field string) template.JS {
	return template.JS(fmt.Sprintf("Bind(this, '%s')", field))
}

// FormSubmit is a wrapper around JS function to submit a form as an action.
func FormSubmit() template.JS {
	return "FormSubmit(this, event)"
}
