package kyoto

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

// ****************
// Action configuration
// ****************

// ActionConfiguration holds a global actions configuration.
type ActionConfiguration struct {
	Path       string // Configure a path prefix for action calls
	Terminator string // Configure a terminator sequence which responsible for chunk separation
}

// ActionConf is a global configuration that will be used during actions handling.
// See ActionConfiguration for more details.
var ActionConf = ActionConfiguration{
	Path:       "/internal/actions/",
	Terminator: "=!EOC!=", // EOC (End Of Component)
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
	redirected bool
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
func ActionFlush(ctx *Context, state any) {
	// Exit if redirected.
	// ActionFlush must not to be triggered after redirection,
	// but still adding this for safety purposes.
	// (f.e. manual ActionFlush call after redirect without return)
	if ctx.Action.redirected {
		return
	}
	// Initialize flusher
	flusher := ctx.ResponseWriter.(http.Flusher)
	// Clone prepared template
	tmpl, _ := ctx.Template.Clone()
	// Render template into buffer
	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, state); err != nil {
		panic(err)
	}
	// Append terminator sequence
	// Details: https://todo.sr.ht/~kyoto-framework/kyoto-framework/10
	buf.Write([]byte(ActionConf.Terminator))
	// Write to stream
	if _, err := fmt.Fprint(ctx.ResponseWriter, buf.String()); err != nil {
		panic(err)
	}
	// Flush
	flusher.Flush()
}

// ActionRedirect is a function to trigger redirect during action handling.
//
// Example:
//
//	func CompFoo(ctx *kyoto.Context) (state CompFooState) {
//		...
//		// Handle example action
//		kyoto.Action(ctx, "Bar", func(args ...any) {
//			// Redirect to the home page
//			kyoto.ActionRedirect(ctx, "/")
//		})
//		...
//	}
func ActionRedirect(ctx *Context, location string) {
	// Initialize flusher
	flusher := ctx.ResponseWriter.(http.Flusher)
	// Create command
	cmd := fmt.Sprintf("ssa:redirect=%s", location)
	// Append terminator sequence and write to stream
	// Details: https://todo.sr.ht/~kyoto-framework/kyoto-framework/10
	if _, err := fmt.Fprint(ctx.ResponseWriter, cmd + ActionConf.Terminator); err != nil {
		panic(err)
	}
	// Set redirected flag
	ctx.Action.redirected = true
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
	builder.WriteString(fmt.Sprintf("<script>const actionpath = \"%s\"; const actionterminator = \"%s\"</script>", ActionConf.Path, ActionConf.Terminator))
	builder.WriteString(ActionClient)
	return template.HTML(builder.String())
}

// ****************
// Action handling
// ****************

// HandleAction registers a component action handler with a predefined pattern in the DefaultServeMux.
// It's a wrapper around http.HandleFunc, but accepts a component instead of usual http.HandlerFunc.
//
// Example:
//
//	kyoto.HandleAction(CompFoo) // Register a usual component
//	kyoto.HandleAction(CompBar("")) // Register a component which accepts arguments and returns wrapped function
func HandleAction[T any](component Component[T], ctx ...*Context) {
	pattern := ActionConf.Path + ComponentName(component) + "/"
	logf("Registering component action handler '%s':\t'%s'", ComponentName(component), pattern)
	http.HandleFunc(pattern, HandlerAction(component, ctx...))
}

// HandlerAction returns a http.HandlerFunc that handles an action request for a specified component.
// Pattern still must to correspond to the provided component.
// It's recommended to use HandleAction instead.
//
// Example:
//
//	http.HandleFunc("/internal/actions/Foo/", kyoto.HandlerAction(Foo))
func HandlerAction[T any](component Component[T], _ctx ...*Context) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set headers
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Transfer-Encoding", "chunked")
		w.Header().Set("Cache-Control", "no-store")
		// Extract action parameters
		action := ActionParameters{}
		if err := action.Parse(r); err != nil {
			panic(err)
		}
		// Initialize context
		ctx := &Context{}
		if len(_ctx) > 0 {
			ctx = _ctx[0]
		}
		// Set context parameters
		ctx.Request = r
		ctx.ResponseWriter = w
		ctx.Action = action
		// Prepare template
		TemplateInline(ctx, fmt.Sprintf(`{{ template "%s" . }}`, action.Component))
		// Trigger building
		state := component(ctx)
		// Trigger flush (if not redirected)
		if !ctx.Action.redirected {
			ActionFlush(ctx, state)
		}
	}
}
