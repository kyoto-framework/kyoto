package ssc

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/url"
	"reflect"
	"strings"
)

// Code, responsible for dynamics, like Server Side Actions, bindings, etc.
var dynamics = `
<script>

function Action(self, action, ...args) {
	// Determine depth
	let depth = (action.split('').filter(x => x === '$') || []).length
	action = action.replaceAll('$', '')
	// Find component root
	let root = self
	let dcount = 0
	while (true) {
		if (!root.getAttribute('state')) {
			root = root.parentElement
		} else {
			if (dcount != depth) {
				root = root.parentElement
				dcount++
			} else {
				break
			}
		}
	}
	// Prepare form data
	let formdata = new FormData()
	formdata.set('State', root.getAttribute('state'))
	formdata.set('Args', JSON.stringify(args))
	// Make request
	fetch("/SSA/"+root.getAttribute('name')+"/"+action, {
		method: 'POST',
		body: formdata
	}).then(resp => {
		return resp.text()
	}).then(data => {
		root.outerHTML = data
	}).catch(err => {
		console.log(err)
	})
}

function Bind(self, field) {
	// Find component root
	let root = self
	while (true) {
		if (!root.getAttribute('state')) {
			root = root.parentElement
		} else {
			break
		}
	}
	// Load state
	let state = JSON.parse(decodeURIComponent(root.getAttribute('state')))
	console.log(state)
	// Set value
	state[field] = self.value
	// Set state
	root.setAttribute('state', JSON.stringify(state))
}
</script>
`

func Funcs() template.FuncMap {
	return template.FuncMap{
		"meta": func(p Page) template.HTML {
			builder := ""
			meta := Meta{}
			if p, ok := p.(ImplementsMeta); ok {
				meta = p.Meta()
			}
			if meta.Title != "" {
				builder += "<title>" + meta.Title + "</title>\n"
			}
			if meta.Canonical != "" {
				builder += `<link rel="canonical" href="` + meta.Canonical + `">` + "\n"
			}
			if len(meta.Hreflangs) != 0 {
				for _, hreflang := range meta.Hreflangs {
					builder += `<link rel="alternate" hreflang="` + hreflang.Lang + `" href="` + hreflang.Href + `">` + "\n"
				}
			}
			return template.HTML(builder)
		},
		"dynamics": func() template.HTML {
			return template.HTML(dynamics)
		},
		"json": func(data interface{}) string {
			d, _ := json.Marshal(data)
			return string(d)
		},
		"componentattrs": func(c Component) template.HTMLAttr {
			// Extract component data
			name := reflect.ValueOf(c).Elem().Type().Name()
			statebytes, err := json.Marshal(c)
			if err != nil {
				panic(err)
			}
			state := string(statebytes)
			state = url.QueryEscape(state)
			// Build attributes
			builder := fmt.Sprintf(`name='%s' state='%s'`, name, state)
			return template.HTMLAttr(builder)
		},
		"action": func(action string, args ...interface{}) template.JS {
			formattedargs := []string{}
			for _, arg := range args {
				b, _ := json.Marshal(arg)
				formattedargs = append(formattedargs, string(b))
			}

			return template.JS(fmt.Sprintf("Action(this, '%s', %s)", action, strings.Join(formattedargs, ", ")))
		},
		"bind": func(field string) template.JS {
			return template.JS(fmt.Sprintf("Bind(this, '%s')", field))
		},
	}
}
