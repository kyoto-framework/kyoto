/*
Package kyoto was made for creating fast, server side frontend avoiding vanilla templating downsides.

It tries to address complexities in frontend domain like
responsibility separation, components structure, asynchronous load
and hassle-free dynamic layout updates.
These issues are common for frontends written with Go.

The library provides you with primitives for pages and components creation,
state and rendering management, dynamic layout updates (with external packages integration),
utility functions and asynchronous components out of the box.
Still, it bundles with minimal dependencies
and tries to utilize built-ins as much as possible.

You would probably want to opt out from this library in few cases, like,
if you're not ready for drastic API changes between major version,
you want to develop SPA/PWA and/or complex client-side logic,
or you're just feeling OK with your current setup.
Please, don't compare kyoto with a popular JS libraries like React, Vue or Svelte.
I know you will have such a desire, but most likely you will be wrong.
Use cases and underlying principles are just too different.

If you want to get an idea of what a typical static component would look like, here's some sample code.
It's very ascetic and simplistic, as we don't want to overload you with implementation details.
Markup is also not included here (it's just a well-known `html/template`).

	// State is declared separately from component itself
	type ComponentState struct {
		// We're providing component with some abilities here
		component.Universal // This component uses universal state, that can be (un)marshalled with both server and client

		// Actual component state is just struct fields
		Foo string
		Bar string
	}

	// Component is a function, that returns configured state.
	// To be able to provide additional arguments to the component on initialization,
	// you have to wrap component with additional function that will handle args and return actual component.
	// Until then, you may keep component declaration as-is.
	func Component(ctx *component.Context) component.State {
		// Initialize state here.
		// As far as component.Universal provided in declaration,
		// we're implementing component.State interface.
		state := &ComponentState{}
		// Do whatever you want with a state
		state.Foo = "foo"
		state.Bar = "bar"
		// Done
		return state
	}

For details, please check project's website on https://kyoto.codes.
Also, you may check the library index to explore available sub-packages
and https://pkg.go.dev for Go'ish documentation style.

# Quick start

We don't want you to deal with boilerplate code on your own,
so you can proceed with our simple starter project.

	git clone https://github.com/kyoto-framework/new <your-new-project>
	rm -r <your-new-project>/.git

Feel free to use it as an example for your own setup.

# Components

Components is a common approach for modern libraries to manage frontend parts.
Kyoto's components are trying to be mostly independent (but configurable) part of the project.

To create component, it would be enough to implement component.Component.
It's a function, a context receiver which returns a component state.
State is an implementation of component.State,
which is easy to implement with nesting one of the state implementations (options will be described later).

	package main

	type ComponentState struct {
		component.Disposable // You're providing component with some abilities here
	}

	func Component(ctx *component.Context) component.State {
		state := &ComponentState{}
		return state
	}

	...

Each component becomes a part of the page or top-level component,
which executes component function asynchronously and gets a state future object.
In that way your components are executing in a non-blocking way.

Pages are just top-level components, where you can configure rendering and page related stuff.

# Components with state

Stateful components are pretty similar to stateless ones,
but they are actually implementing marshal/unmarshal interface instead of mocking it.

You have multiple state options to choose from: universal or server.

Universal state is a state, that can be marshalled and unmarshalled both on server and client.
It's a common state option without functionality limitations.
On the other hand, the whole state must be sent and received,
which applies some limitations on the state size.

	package main

	type ComponentState struct {
		component.Universal // This state allows you to operate with data on both server and client
	}

	func Component(ctx *component.Context) component.State {
		state := &ComponentState{}
		return state
	}

Server state can be marshalled and unmarshalled only on server.
It's a good option for components, that are not supposed to be updated on client side (f.e. no inputs).
Also, it's a good option for components with lots of state data.

	package main

	type ComponentState struct {
		component.Server // This state allows you to operate with data on server only
	}

	func Component(ctx *component.Context) component.State {
		state := &ComponentState{}
		return state
	}

# Components with arguments

Sometimes you may want to pass some arguments to the component.
It's easy to do with wrapping component with additional function.

	package main

	type ComponentState struct {
		component.Universal

		Data string
	}

	func Component(data string) component.Component {
		return func(ctx *component.Context) component.State {
			state := &ComponentState{}
			state.Data = data // We are passing arg to the component state, but it's not a requirement.
			return state
		}
	}

# Context

You have an access to the context inside the component.
It includes request and response objects, as well as some other useful stuff like store.

	package main

	...

	func Component(ctx *component.Context) component.State {
		...
		ctx.Request // http.Request
		ctx.Response // http.ResponseWriter
		ctx.Set("k", "v") // Store arbitrary data in the context
		v := ctx.Get("k").(string) // Get arbitrary data from the context
		...
	}

# Routing

This library doesn't provide you with routing out of the box.
You can use any router you want, built-in one is not a bad option for basic needs.

	package main

	...

	func main() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", rendering.Handler(Page))
		http.ListenAndServe(":8080", mux)
	}

# Rendering

Rendering might be tricky, but we're trying to make it as simple as possible.
By default, we're using `html/template` as a rendering engine.
It's a well-known built-in package, so you don't have to learn anything new.

Out of the box we're parsing all templates in root directory with `*.html` glob.
You can change this behavior with `TEMPLATE_GLOB` global variable.
Don't rely on file names while working with template names,
use `define` entry for each your component.

To provide your components with ability to be rendered, you have to do some basic steps.
First, you have to nest one of the rendering implementations into your component state (f.e. `rendering.Template`).

	package main

	type ComponentState struct {
		component.Disposable
		rendering.Template // This line allows you to render your component with html/template
	}

	...

You can customize rendering with providing values to the rendering implementation.
If you need to modify these values for the entire project,
we recommend looking at the global settings or creating a builder function for rendering object.

	package main

	...

	func Component(ctx *component.Context) component.State {
		state := &ComponentState{}
		state.Template.Name = "CustomName" // Set custom template name
		...
	}

By default, render handler will use a component name as a template name.
So, you have to define a template with the same name as your component
(not the filename, but "define" entry).

	{{ define "Component" }}
		...
	{{ end }}

That's enough to be rendered by `rendering.Handler`.

For rendering a nested component, use built-in `template` function.
Provide a resolved future object as a template argument in this way.
Nested components are not obligated to have rendering implementation if you're using them in this way.

	<div>{{ template "component" call .Component }}</div>

As an alternative, you can nest rendering implementation (e.g. `rendering.Template`) into your nested component.
In this way you can use `render` function to simplify your code.
Please, don't use this approach heavily now, as it affects rendering performance.

	<div>{{ render .Component }}</div>

# HTMX

HTMX is a frontend library, that allows you to update your page layout dynamically.
It perfectly fits into kyoto, which focuses on components and server side rendering.
Thanks to the component structure, there is no need to define separate rendering logic specially for HTMX.

# HTMX Setup

Please, check https://htmx.org/docs/#installing for installation instructions.
In addition to this, you must register HTMX handlers for your dynamic components.

	package main

	...

	func main() {
		// Initialize mux
		mux := http.NewServeMux()
		// Register pages
		mux.HandleFunc("/", rendering.Handler(Page))
		// Register components
		mux.HandleFunc("/htmx/component", rendering.Handler(Component))
		// Serve
		http.ListenAndServe(":8080", mux)
	}

# HTMX Usage

This is a basic example of HTMX usage.
Please, check https://htmx.org/docs/ for more details.

In this example we're defining a form component, that is updating itself on submit.

	{{ define "Component" }}
	<form hx-post="/htmx/component" hx-target="this" hx-swap="outerHTML">
		<input type="text" name="foo" value="{{ .Foo }}">
		<input type="text" name="bar" value="{{ .Bar }}">
		<button type="submit">Submit</button>
	</form>
	{{ end }}

And this is how you can define a component, that will handle this request.

	package main

	type ComponentState struct {
		component.Disposable // We're not using any stored state here, so we're using disposable
		rendering.Template   // We're using template rendering for this component, just like in pages

		Foo string
		Bar string
	}

	func Component(ctx *component.Context) component.State {
		// Initialize state
		state := &ComponentState{}
		// We're getting request data from context and passing it to the state
		if ctx.Request.Method == http.MethodPost {
			ctx.Request.ParseForm()
			state.Foo = ctx.Request.FormValue("foo")
			state.Bar = ctx.Request.FormValue("bar")
		}
		// Done
		return state
	}

# HTMX State

Sometimes it might be useful to have a component state,
which will persist between requests and will be stored without any actual usage in the client side presentation.

	<form hx-post="/htmx/component" hx-target="this" hx-swap="outerHTML">
		{{ hxstate . }}
		<div>Cursor: {{ .Cursor }}</div>
		<button type="submit">Submit</button>
	</form>

This function injects a hidden input field with a serialized state.
Let's check how it works on the server side.

	package main

	type ComponentState struct {
		component.Universal // We're using server state here
		rendering.Template  // We're using template rendering for this component, just like in pages

		Cursor string
	}

	func Component(ctx *component.Context) component.State {
		// Initialize state
		state := &ComponentState{}
		// Unmarshal state on post request
		if ctx.Request.Method == http.MethodPost {
			ctx.Request.ParseForm()
			state.Unmarshal(ctx.Request.FormValue("hx-state"))
		}
		// Initialize cursor if it's empty
		if state.Cursor == "" {
			state.Cursor = "..."
		}
		// Done
		return state
	}

As a result, we have a component with a persistent state between requests.
*/
package kyoto
