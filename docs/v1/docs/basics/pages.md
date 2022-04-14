
# Pages

Page definition is an entry point of out page rendering.
The most basic way to define a page is to use template builder.
Let's see how it works:

=== "page.index.go"

	```go
	package main

	import (
	    "html/template"

	    "github.com/kyoto-framework/kyoto"
	    "github.com/kyoto-framework/kyoto/render"
	    "github.com/kyoto-framework/kyoto/lifecycle"
	)

	func PageIndex(core *kyoto.Core) {
	    lifecycle.Init(core, func() {
	        core.State.Set("Content", "Hello, Kyoto!")
	    })
	    render.Template(core, func() *template.Template {
	        return template.Must(template.New("page.index.html").Funcs(render.FuncMap(core)).ParseGlob("*.html"))
	    })
	}
	```

=== "page.index.html"

	```html
	<!DOCTYPE html>
	<html>
	<head>
	    <title>Page Index</title>
	</head>
	<body>
	    <h1>{{ .Content }}</h1>
	</body>
	</html>
	```

First, we will get the terms right:

- Core receiver - function, receiving a `kyoto.Core` as an argument.
- Adapter - function, that interacts with a core.

We need to understand is that everything in this library built around `kyoto.Core`.
Every page, every component, every module interacts with a core.
Kyoto is not a single package.
It's a set of packages that are built around a core.

Let's return to the example.
You can see a core receiver here.
Every page or component is a core receiver.
To add page functionality to this core receiver, we are using `render` module.
With a help of `render.Template` we are specifying a template builder.

Also, as you can see, there is a `lifecycle.Init` adapter.
This adapter is responsible for defining of initializing function.
Lifecycle functionality will be described in a [Basics • Lifecycle](/basics/lifecycle) documentation category.

Another thing you may notice in this example is setting of a state.
To interact with a state, we are using a `kyoto.State` instance inside of a core.
There are 3 methods for state modification: `Set`, `Get`, `Del`.
For more information, check [Basics • State](/basics/state) documentation category.
