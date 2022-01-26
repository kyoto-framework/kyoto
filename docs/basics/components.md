
# Components

Components approach is very common for modern web frameworks.
It's not very popular in template engines.
Even more, in most cases you can't define own functionality for each component.
Kyoto library tries to combine this approach with template engines.

To define a component, we are using core receiver, just like we did with pages.

=== "component.uuid.go"

	```
	package main

	import (
	    "net/http"
	    "encoding/json"

	    "github.com/kyoto-framework/kyoto"
	    "github.com/kyoto-framework/kyoto/lifecycle"
	)

	func ComponentUUID(core *kyoto.Core) {
	    lifecycle.Init(core, func() {
	        core.State.Set("UUID", "")
	    })
	    lifecycle.Async(core, func() error {
	        resp, _ := http.Get("http://httpbin.org/uuid")
	        data := map[string]string{}
	        json.NewDecoder(resp.Body).Decode(&data)
	        c.State.Set("UUID", data["uuid"])
	        return nil
	    })
	}
	```

=== "component.uuid.html"

	```html title="component.uuid.html"
	{{ define "ComponentUUID" }}
	<div>
	    UUID: {{ .UUID }}
	</div>
	{{ end }}
	```

And now let's attach this component to the page multiple times.

=== "page.index.go"

	```go
	...

	func PageIndex(core *kyoto.Core) {
	    lifecycle.Init(core, func() {
	        core.Component("UUID1", ComponentUUID)
	        core.Component("UUID2", ComponentUUID)
	    })
	    ...
	}
	```

Here you can see multiple things:

- Component definition
- Component template
- Component attaching to page

Our component just defines 2 lifecycle functions:
init with setting empty state and async with getting UUID from httpbin.org.
Lifecycle functionality will be described in a "Lifecycle" documentation category.

To define a component template, we are using "define" template functionality.
Please note, definition name must to be the same as a component name.

Last thing in this list is creating a component instance and attaching to a page.
To do this, we are using `core.Component` function.

To create a parameterized component, you can use a wrapped core receiver.

```go
func ComponentFoo(param1, param2 string) func(*kyoto.Core) {
    return func(core *kyoto.Core) {
        ...
    }
}
```

Attach it to a page in the same way, but with a call.

```go
core.Component("Foo", ComponentFoo("param1", "param2"))
```
