
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
