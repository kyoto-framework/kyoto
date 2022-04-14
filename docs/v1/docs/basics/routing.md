
# Routing

Our journey to explore library basics begins with routing.
This is also the easiest part to learn, as it is based on built-in router.
Let's see how it works:

=== "main.go"

	```go
	package main

	import (
	    "net/http"

	    "github.com/kyoto-framework/kyoto"
	    "github.com/kyoto-framework/kyoto/render"
	)

	func PageIndex(core *kyoto.Core) {
	    ...
	}

	func main() {
	    http.HandleFunc("/", render.PageHandler(PageIndex))
	    http.ListenAndServe(":8080", nil)
	}
	```

We will skip page functionality for now and focus on routing itself.
On this example, you can see how we are using `render.PageHandler` to wrap our page definition.
Under the hood, it makes a lot of things to render our page.
Also, you are free to use any framework you want, that supports `http.HandlerFunc` interface.

Without going off topic, let's also check how we can use redirects:

=== "main.go"

	```go
	...
	func PageIndex(core *kyoto.Core) {
		render.Redirect(core, "/", 307)
	}
	...
	```

This function call will inject a redirection to the core.
During rendering page handler will check if there is a redirection and will redirect with the given parameters.
It's important to use kyoto's redirect to avoid header rewriting error.
