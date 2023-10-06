package action

import (
	"fmt"
	"html/template"
	"strings"
)

var FuncMap = template.FuncMap{
	"client": func(path string) template.HTML {
		builder := strings.Builder{}
		builder.WriteString(fmt.Sprintf("<script>const actionpath = \"%s\"; const actionterminator = \"%s\"</script>", path, terminator))
		builder.WriteString(fmt.Sprintf("<script>%s</script>", client))
		return template.HTML(builder.String())
	},
}
