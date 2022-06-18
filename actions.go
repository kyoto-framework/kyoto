/*
	-

	Actions

	Kyoto provides a way to simplify building dynamic UIs.
	For this purpose it has a feature named actions.
	Logic is pretty simple.
	Action is executing on server side,
	server is sending updated component markup to the client
	which will be morphed into DOM.
	That's it.
*/
package kyoto

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

// ****************
// Action configuration
// ****************

type ActionConfiguration struct {
	Path string
}

var ActionConf = ActionConfiguration{
	Path: "/internal/actions/",
}

// ****************
// Action parameters representation
// ****************

type ActionParameters struct {
	Component string
	Action    string
	State     string
	Args      []any
}

func (p *ActionParameters) Parse(r *http.Request) error {
	// Validate request format
	if r.FormValue("State") == "" {
		return errors.New("state is empty")
	}
	if r.FormValue("Args") == "" {
		return errors.New("args is empty")
	}
	// Split path into tokens
	tokens := strings.Split(r.URL.Path, "/")
	// Extract component state as a raw json string (will be decoded with ActionPreload)
	p.State = r.FormValue("State")
	// Extract component arguments
	err := json.Unmarshal([]byte(r.FormValue("Args")), &p.Args)
	if err != nil {
		return errors.New("something wrong with arguments")
	}
	// Extract component & action names
	p.Component = tokens[len(tokens)-2]
	p.Action = tokens[len(tokens)-1]
	// Return
	return nil
}

// ****************
// Action methods
// ****************

func Action(c *Context, name string, action func(args ...any)) bool {
	if c.Action.Action == name {
		action(c.Action.Args...)
		return true
	}
	return false
}

func ActionPreload[T any](c *Context, state T) {
	// Pass if not an action
	if c.Action.Component == "" {
		return
	}
	// Unmarshal state
	UnmarshalState(c.Action.State, state)
}

func ActionFlush(c *Context, state any) {
	// Initialize flusher
	flusher := c.ResponseWriter.(http.Flusher)
	// Clone prepared template
	tmpl, _ := c.Template.Clone()
	// Render template into flusher
	if err := tmpl.Execute(c.ResponseWriter, state); err != nil {
		panic(err)
	}
	// Flush
	flusher.Flush()
}

// ****************
// Action template functions
// ****************

func actionFuncState(state any) template.HTMLAttr {
	return template.HTMLAttr(fmt.Sprintf(`state="%s"`, MarshalState(state)))
}

func actionFuncClient() template.HTML {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("<script>const ssapath = \"%s\"</script>", ActionConf.Path))
	builder.WriteString(ActionClient)
	return template.HTML(builder.String())
}

// ****************
// Action handling
// ****************

func HandleAction[T any](component Component[T]) {
	pattern := ActionConf.Path + ComponentName(component) + "/"
	log.Printf("Registering '%s' component action handler under '%s'", ComponentName(component), pattern)
	http.HandleFunc(pattern, HandlerAction(component))
}

func HandlerAction[T any](component Component[T]) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set headers
		w.Header().Set("Content-Type", "plain/html")
		w.Header().Set("Transfer-Encoding", "chunked")
		w.Header().Set("Cache-Control", "no-store")
		// Extract action parameters
		a := ActionParameters{}
		if err := a.Parse(r); err != nil {
			panic(err)
		}
		// Prepare context
		ctx := &Context{
			ResponseWriter: w,
			Request:        r,
			Action:         a,
		}
		// Prepare template
		TemplateInline(ctx, fmt.Sprintf(`{{ template "%s" . }}`, a.Component))
		// Trigger building
		state := component(ctx)
		// Trigger flush
		ActionFlush(ctx, state)
	}
}
