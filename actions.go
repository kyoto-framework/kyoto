package kyoto

import (
	"bytes"
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

// ActionConfiguration holds a global actions configuration.
type ActionConfiguration struct {
	Path string // Configure a path prefix for action calls
}

// ActionConf is a global configuration that will be used during actions handling.
// See ActionConfiguration for more details.
var ActionConf = ActionConfiguration{
	Path: "/internal/actions/",
}

// ****************
// Action parameters representation
// ****************

// ActionParameters is a Go representation of an action request.
type ActionParameters struct {
	Component string
	Action    string
	State     string
	Args      []any

	processed bool
}

// Parse extracts action data from a provided request.
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

// Action is a function that handles an action request.
// Returns an execution flag (true if action was executed).
// You can use a flag to prevent farther execution of a component.
//
// Example:
//
//	func Foo(ctx *kyoto.Context) (state FooState) {
//		// Handle action
//		bar := kyoto.Action(ctx, "Bar", func(args ...any) {
//			// Do something
//		})
//		// Prevent further execution
//		if bar {
//			return
//		}
//		// Default non-action behavior
//		// ...
//	}
//
func Action(c *Context, name string, action func(args ...any)) bool {
	// This will allow to avoid recursive action call
	// while attaching recursive component
	if c.Action.processed {
		return false
	}
	if c.Action.Action == name {
		action(c.Action.Args...)
		c.Action.processed = true
		return true
	}
	return false
}

// ActionPreload unpacks a component state from an action request.
// Executing only in case of an action request.
//
// Example:
//
//	func CompFoo(ctx *kyoto.Context) (state CompFooState) {
//		// Preload state
//		kyoto.ActionPreload(ctx, &state)
//		// Handle actions
//		...
//	}
func ActionPreload[T any](c *Context, state T) {
	// Pass if not an action
	if c.Action.Component == "" {
		return
	}
	// Unmarshal state
	UnmarshalState(c.Action.State, state)
}

// ActionFlush allows to push multiple component UI updates during single action call.
// Call it when you need to push an updated component markup to the client.
//
// Example:
//
//	func CompFoo(ctx *kyoto.Context) (state CompFooState) {
//		...
//		// Handle example action
//		kyoto.Action(ctx, "Bar", func(args ...any) {
//			// Do something with a state
//			state.Content = "Bar"
//			// Push updated UI to the client
//			kyoto.ActionFlush(ctx, state)
//			// Do something else with a state
//			state.Content = "Baz"
//			// Push updated UI to the client
//			kyoto.ActionFlush(ctx, state)
//		})
//		...
// }
func ActionFlush(c *Context, state any) {
	// Initialize flusher
	flusher := c.ResponseWriter.(http.Flusher)
	// Clone prepared template
	tmpl, _ := c.Template.Clone()
	// Render template into buffer
	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, state); err != nil {
		panic(err)
	}
	// Write to stream
	if _, err := fmt.Fprint(c.ResponseWriter, buf.String()); err != nil {
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

// HandleAction registers a component action handler with a predefined pattern in the DefaultServeMux.
// It's a wrapper around http.HandlePage, but accepts a component instead of usual http.HandlerFunc.
//
// Example:
//
//	kyoto.HandleAction(CompFoo) // Register a usual component
//	kyoto.handleAction(CompBar("")) // Register a component which accepts arguments and returns wrapped function
func HandleAction[T any](component Component[T]) {
	pattern := ActionConf.Path + ComponentName(component) + "/"
	log.Printf("Registering '%s' component action handler under '%s'", ComponentName(component), pattern)
	http.HandleFunc(pattern, HandlerAction(component))
}

// HandlerAction returns a http.HandlerFunc that handles an action request for a specified component.
// Pattern still must to correspond to the provided component.
// It's recommended to use HandleAction instead.
//
// Example:
//
//	http.HandleFunc("/internal/actions/Foo/", kyoto.HandlerAction(Foo))
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
