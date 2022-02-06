
# Context

During the development process, we may need to access global things or internals like request or response objects.
In the same time, kyoto modules sometimes need to communicate with each other or temporary save some shared stuff.
That's why context was needed to solve that common problems.

Let's check this example.

=== "main.go"

	```go

	func ComponentTest(core *kyoto.Core) {
		lifecycle.Init(core, func() {
			req := core.Context.GetRequest()
			println("User-Agent:", req.UserAgent())
			println("Context Example:", core.Context.Get("foo").(string))
		})
	}

	func PageIndex(core *kyoto.Core) {
		lifecycle.Init(core, func() {
			core.Context.Set("foo", "bar")
			core.Component("c1", ComponentTest)
		})
	}
	```

Context is: 

- Based on `kyoto.Store`, a simple atomic wrapper around a map.
- Context value is accessible globally after setting.
- Lives during page processing and dies on processing completion.
  It's not stored anywhere, unlike state in dynamic components.

For generic data, like in the state, you have 3 methods: `Get`, `Set`, `Del`.
You'll need to use type casting because of the generic nature of the Store.
Also, there are 2 additional getters for common objects: `GetRequest`, `GetResponseWriter`.
