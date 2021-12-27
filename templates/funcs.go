package templates

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strings"

	"github.com/kyoto-framework/kyoto/actions"
	"github.com/kyoto-framework/kyoto/helpers"
)

func Render(component interface{}) string {
	return ""
}

func Meta(page interface{}) template.HTML {
	return ""
}

func Dynamics(path ...string) template.HTML {
	if len(path) == 0 {
		path = append(path, "/SSA")
	}
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("<script>const ssapath = \"%s\"</script>", path[0]))
	builder.WriteString(actions.Client)
	return template.HTML(builder.String())
}

func JSON(data interface{}) string {
	return ""
}

func ComponentAttrs(component interface{}) template.HTMLAttr {
	return template.HTMLAttr(fmt.Sprintf(
		`cid='%s' name='%s' state='%s'`,
		helpers.ComponentID(component),
		helpers.ComponentName(component),
		helpers.ComponentSerialize(component),
	))
}

func Action(action string, args ...interface{}) template.JS {
	var formattedargs []string
	for _, arg := range args {
		b, _ := json.Marshal(arg)
		formattedargs = append(formattedargs, string(b))
	}

	return template.JS(fmt.Sprintf("Action(this, '%s', %s)", action, strings.Join(formattedargs, ", ")))
}

func Bind(field string) template.JS {
	return template.JS(fmt.Sprintf("Bind(this, '%s')", field))
}

func FormSubmit() template.JS {
	return "FormSubmit(this, event)"
}
