package kyoto

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"reflect"
	"strings"
	"time"
)

// TFuncMap is responsible for integration of library functionality into template rendering
// You need to use this function while template building (or mix with your own)
func TFuncMap() template.FuncMap {
	return template.FuncMap{
		"meta":           TMeta,
		"dynamics":       TDynamics,
		"json":           TJSON,
		"componentattrs": TComponentAttrs,
		"action":         TAction,
		"bind":           TBind,
		"formsubmit":     TFormSubmit,
	}
}

func TRender(c Component) string {
	// Get component insights
	insights := GetInsights(c)
	st := time.Now()
	// Render
	dom := ""
	if _c, ok := c.(ImplementsRender); ok {
		dom = _c.Render()
	} else {
		panic("Trying to render a component which is not implemented Render")
	}
	// Save render insights
	if insights != nil {
		insights.Update(InsightsTiming{
			Render: time.Since(st),
		})

	}
	// Return
	return dom
}

func TMeta(p Page) template.HTML {
	builder := strings.Builder{}

	meta := Meta{}
	if p, ok := p.(ImplementsMeta); ok {
		meta = p.Meta()
	}
	if meta.Title != "" {
		builder.WriteString("<title>" + meta.Title + "</title>\n")
	}
	if meta.Canonical != "" {
		builder.WriteString(`<link rel="canonical" href="` + meta.Canonical + `">` + "\n")
	}
	if len(meta.Hreflangs) != 0 {
		for _, hreflang := range meta.Hreflangs {
			builder.WriteString(`<link rel="alternate" hreflang="` + hreflang.Lang + `" href="` + hreflang.Href + `">` + "\n")
		}
	}
	return template.HTML(builder.String())
}

func TDynamics(path ...string) template.HTML {
	if len(path) == 0 {
		path = append(path, "/SSA")
	}
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("<script>const ssapath = \"%s\"</script>", path[0]))
	builder.WriteString(ssaclient)
	return template.HTML(builder.String())
}

func TJSON(data interface{}) string {
	d, _ := json.Marshal(data)
	return string(d)
}

func TComponentAttrs(c Component) template.HTMLAttr {
	// Extract component data
	name := reflect.ValueOf(c).Elem().Type().Name()
	statebytes, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	state := strings.ReplaceAll(base64.StdEncoding.EncodeToString(statebytes), "/", "-")
	// Build attributes
	builder := fmt.Sprintf(`name='%s' state='%s'`, name, state)
	return template.HTMLAttr(builder)
}

func TAction(action string, args ...interface{}) template.JS {
	var formattedargs []string
	for _, arg := range args {
		b, _ := json.Marshal(arg)
		formattedargs = append(formattedargs, string(b))
	}

	return template.JS(fmt.Sprintf("Action(this, '%s', %s)", action, strings.Join(formattedargs, ", ")))
}

func TBind(field string) template.JS {
	return template.JS(fmt.Sprintf("Bind(this, '%s')", field))
}

func TFormSubmit() template.JS {
	return "FormSubmit(this, event)"
}
