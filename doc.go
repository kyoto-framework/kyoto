/*
	Extensible Go library for creating fast, SSR-first frontend avoiding vanilla templating downsides.

	Motivation

	Creating asynchronous and dynamic layout parts is a complex problem for larger projects using `html/template`.
	Library tries to simplify this process.

	Components

	Kyoto provides a way to define components.
	It's a very common approach for modern libraries to manage frontend parts.
	In kyoto each component is a context receiver, which returns it's state.
	Each component becomes a part of the page or top-level component,
	which executes component asynchronously and gets a state future object.
	In that way your components are executing in a non-blocking way.

		func CUUID(ctx *kyoto.Context) (state CUUIDState) {
			// Fetch uuid data
			resp, _ := http.Get("http://httpbin.org/uuid")
			data := map[string]string{}
			json.NewDecoder(resp.Body).Decode(&data)
			// Set state
			state.UUID = data["uuid"]
		}

	Context

	Kyoto provides a context,
	which holds common objects like http.ResponseWriter, *http.Request, etc.

		func Component(ctx *kyoto.Context) (state ComponentState) {
			log.Println(ctx.Request.UserAgent())
			...
		}

	Template

	Kyoto provides a set of parameters and functions
	to provide a comfortable template building process.

		func Page(ctx *kyoto.Context) (state PageState) {
			// By default it will:
			// - use kyoto.FuncMap as a FuncMap
			// - parse everything in the current directory with a .ParseGlob("*.html")
			// - render a template with a given name
			kyoto.Template(ctx, "page.index.html")
		}

	HTTP

	Kyoto provides a simple net/http handlers and function wrappers
	to handle pages rendering and serving.

		func main() {
			kyoto.HandlePage("/foo", PageFoo)
			kyoto.HandlePage("/bar", PageBar)

			kyoto.Serve(":8000")
		}

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
