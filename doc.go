/*
	Extensible Go library for creating fast, SSR-first frontend avoiding vanilla templating downsides.

	Motivation

	Creating asynchronous and dynamic layout parts is a complex problem for larger projects using `html/template`.
	Library tries to simplify this process.

	Quick start

	Let's go straight into a simple example.
	Then, we will dig into details, step by step, how it works.

		package main

		import (
			"html/template"
			"encoding/json"

			"github.com/kyoto-framework/kyoto"
		)

		// This example demonstrates main advantage of kyoto library - asynchronous lifecycle.
		// Multiple UUIDs will be fetched from httpbin in asynchronous way, without explicitly touching goroutines
		// and synchronization tools like sync.WaitGroup.

		type CUUIDState struct {
			UUID string
		}

		// Let's assume markup of this component is stored in 'component.uuid.html'
		//
		// {{ define "CUUID" }}
		//  <div>UUID: {{ state.UUID }}</div>
		// {{ end }}
		func CUUID(ctx *kyoto.Context) (state CUUIDState) {
			// Fetch uuid data
			resp, _ := http.Get("http://httpbin.org/uuid")
			data := map[string]string{}
			json.NewDecoder(resp.Body).Decode(&data)
			// Set state
			state.UUID = data["uuid"]
		}

		type PIndexState struct {
			UUID1 kyoto.Component[CUUIDState]
			UUID1 kyoto.Component[CUUIDState]
		}

		// Let's assume markup of this page is stored in 'page.index.html'
		//
		// <!DOCTYPE html>
		// <html lang="en">
		// <head>
		// 	<meta charset="UTF-8">
		// 	<meta http-equiv="X-UA-Compatible" content="IE=edge">
		// 	<meta name="viewport" content="width=device-width, initial-scale=1.0">
		// 	<title>Example</title>
		// </head>
		// <body>
		// 	{{ template "CUUID" .UUID1 }}
		// 	{{ template "CUUID" .UUID2 }}
		// </body>
		// </html>
		func PIndex(ctx *kyoto.Context) (state PIndexState) {
			// Define rendering
			render.Template(ctx, "page.index.html")
			// Attach components
			state.UUID1 = kyoto.Use(ctx, CUUID)
			state.UUID2 = kyoto.Use(ctx, CUUID)
		}

		func main() {
			// Register page
			kyoto.HandlePage("/", PIndex)
			// Serve
			kyoto.Serve(":8080")
		}

	Handling requests

	Kyoto provides a simple net/http handlers and function wrappers
	to handle pages rendering and serving.

	See functions inside of nethttp.go file for details and advanced usage.

	Example:

		func main() {
			kyoto.HandlePage("/foo", PageFoo)
			kyoto.HandlePage("/bar", PageBar)

			kyoto.Serve(":8000")
		}

	Components

	Kyoto provides a way to define components.
	It's a very common approach for modern libraries to manage frontend parts.
	In kyoto each component is a context receiver, which returns it's state.
	Each component becomes a part of the page or top-level component,
	which executes component asynchronously and gets a state future object.
	In that way your components are executing in a non-blocking way.

	Pages are just top-level components, where you can configure rendering and page related stuff.

	Example:

		// Component is a context receiver, that returns it's state.
		// State can be whatever you want (simple type, struct, slice, map, etc).
		func CUUID(ctx *kyoto.Context) (state CUUIDState) {
			// Fetch uuid data
			resp, _ := http.Get("http://httpbin.org/uuid")
			data := map[string]string{}
			json.NewDecoder(resp.Body).Decode(&data)
			// Set state
			state.UUID = data["uuid"]
		}

		// Page is just a top-level component, which attaches components and defines rendering
		func PExample(ctx *kyoto.Context) (state PExampleState) {
			// Define rendering
			kyoto.Template(ctx, "page.example.html")
			// Attach components
			state.UUID1 = kyoto.Use(ctx, CUUID)
			state.UUID2 = kyoto.Use(ctx, CUUID)
		}

	Context

	Kyoto provides a context,
	which holds common objects like http.ResponseWriter, *http.Request, etc.

	See kyoto.Context for details.

	Example:

		func Component(ctx *kyoto.Context) (state ComponentState) {
			log.Println(ctx.Request.UserAgent())
			...
		}

	Template

	Kyoto provides a set of parameters and functions
	to provide a comfortable template building process.
	You can configure template building parameters with
	kyoto.TemplateConf configuration.

	See template.go for available functions
	and kyoto.TemplateConfiguration for configuration details.

	Example:

		func Page(ctx *kyoto.Context) (state PageState) {
			// By default it will:
			// - use kyoto.FuncMap as a FuncMap
			// - parse everything in the current directory with a .ParseGlob("*.html")
			// - render a template with a given name
			kyoto.Template(ctx, "page.index.html")
			...
		}

	Actions

	Kyoto provides a way to simplify building dynamic UIs.
	For this purpose it has a feature named actions.
	Logic is pretty simple.
	Client calls an action (sends a request to the server).
	Action is executing on server side and
	server is sending updated component markup to the client
	which will be morphed into DOM.
	That's it.

	To use actions, you need to go through a few steps.
	You'll need to include a client into page (JS functions for communication)
	and register an actions handler for a needed component.

	Let's start from including a client.

		<html>
			<head>
				...
			</head>
			<body>
				...
				{{ client }}
			</body>
		</html>

	Then, let's register an actions handler for a needed component.

		func main() {
			...
			kyoto.HandleAction(Component)
			...
		}

	That's all!
	Now we ready to use actions to provide a dynamic UI.

	Example:

		...

		type CUUIDState struct {
			UUID string
		}

		// Let's assume markup of this component is stored in 'component.uuid.html'
		//
		//	{{ define "CUUID" }}
		//	<div {{ state . }}>
		//		<div>UUID: {{ state.UUID }}</div>
		//		<button onclick="Action(this, 'Reload')">Reload</button>
		//	</div>
		//	{{ end }}
		func CUUID(ctx *kyoto.Component) (state CUUIDState) {
			// Define uuid loader
			uuid := func() string {
				resp, _ := http.Get("http://httpbin.org/uuid")
				data := map[string]string{}
				json.NewDecoder(resp.Body).Decode(&data)
				return data["uuid"]
			}
			// Handle action
			handled := kyoto.Action(ctx, "Reload", func(args ...any) {
				// We will just set a new uuid and will print a log
				// It's not makes a lot of sense now, but it's just a demonstration example
				state.UUID = uuid()
				log.Println("New uuid was issued:", state.UUID)
			})
			// Prevent further execution if action handled
			if handled {
				return
			}
			// Default loading behavior
			state.UUID = uuid()
		}

		type PIndexState struct {
			UUID1 kyoto.Component[CUUIDState]
			UUID2 kyoto.Component[CUUIDState]
		}

		// Let's assume markup of this page is stored in 'page.index.html'
		//
		//	<!DOCTYPE html>
		//	<html lang="en">
		//	<head>
		//		<meta charset="UTF-8">
		//		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		//		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		//		<title>Example</title>
		//	</head>
		//	<body>
		//		{{ template "CUUID" .UUID1 }}
		//		{{ template "CUUID" .UUID2 }}
		//		{{ client }}
		//	</body>
		//	</html>
		func PIndex(ctx *kyoto.Context) (state PIndexState) {
			// Define rendering
			kyoto.Template(ctx, "page.index.html")
			// Attach components
			state.UUID1 = kyoto.Use(ctx, CUUID)
			state.UUID2 = kyoto.Use(ctx, CUUID)
		}

		func main() {
			kyoto.HandlePage("/", PIndex)
			kyoto.HandleAction(CUUID)
			kyoto.Serve(":8000")
		}

	In this example you can see provided modifications to the quick start example.

	First, we've added a state into our components' markup.
	In this way we are saving our components' state between actions and find a component root.

	Second, we've added a reload button with onclick function call.
	We're using a function Action provided by a client.
	Action triggering will be described in details later.

	Third, we've added an action handler inside of our component.
	This handler will be executed when a client calls an action with a corresponding name.

	Action triggering

	Documentation is not ready yet.

	Action flow control

	Documentation is not ready yet.

	Action rendering options

	Documentation is not ready yet.

*/
package kyoto
