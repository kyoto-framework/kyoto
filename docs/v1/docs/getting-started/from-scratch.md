# From Scratch

Everything begins somewhere.
And we are starting from scratch, with project setup.
This short "quick start" guide will include a basic project setup with just a simple index page, components and few examples.

## Installation

Everything is starting from package installation.

```bash
go get github.com/kyoto-framework/kyoto@master
```

## Entry Point

Let's continue with serving foundations, basis for every web application.

=== "main.go"

	```go
	package main

	import (
	    "net/http"
	    "log"
	    "os"
	)

	func main() {
	    // Init serve mux
	    mux := http.NewServeMux()

	    // Routes
	    // ...

	    // Run
	    if os.Getenv("PORT") == "" {
	        log.Println("Listening on localhost:25025")
	        http.ListenAndServe("localhost:25025", mux)
	    } else {
	        log.Println("Listening on 0.0.0.0:" + os.Getenv("PORT"))
	        http.ListenAndServe(":"+os.Getenv("PORT"), mux)
	    }
	}
	```

## Page

Now, we can define our page.

For detailed pages explanation, check [Basics • Pages](/basics/pages) documentation category.

=== "page.index.go"

	```go
	package main

	import (
	    "html/template"

	    "github.com/kyoto-framework/kyoto"
	    "github.com/kyoto-framework/kyoto/render"
	)

	func PageIndex(core *kyoto.Core) {
		render.Template(core, func() *template.Template {
			return template.Must(template.New("page.index.html").Funcs(render.FuncMap(core)).ParseGlob("*.html"))
		})
	}
	```

=== "page.index.html"

	```html
	<html>
		<head> <title>Home page</title> </head>
		<body>
			Example page
		</body>
	</html>
	```

!!! note ""
    You can define bootstrap functions for easier template definitions. For example:
    ```go
    func newtemplate(core *kyoto.Core, page string) *template.Template {
        return template.Must(template.New(page).Funcs(render.FuncMap(core)).ParseGlob("*.html"))
    }
    ```

## Page routing

For attaching your page, you can simply use the built-in page handler (`render.PageHandler`), right below the Routes comment in your main function.

For detailed routing explanation, check [Basics • Routing](/basics/routing) documentation category.

=== "main.go"

	```go
	...
	mux.HandleFunc("/", render.PageHandler(PageIndex))
	...
	```
	

## Running

Your can run your app with the usual:

```bash
go run .
```

For setting ports or exposing on a local network, you can run with the following command:

```bash
PORT=25025 go run .
```

## Adding a component

Let's define a simple component, which fetches a UUID from httpbin API.
This example is good demonstration of asynchronous library functionality.

For detailed components explanation, check [Basics • Components](/basics/components) documentation category.  
For detailed lifecycle explanation, check [Basics • Lifecycle](/basics/lifecycle) documentation category.

=== "component.uuid.go"

	```go
	package main

	import (
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

	```html
	{{ define "ComponentUUID" }}
	<div>
		httpbin.org uuid: {{ .UUID }}
	</div>
	{{ end }}
	```

After component definition, let's use it multiple times in our index page.

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
			core.Component("UUID1", ComponentUUID)
			core.Component("UUID2", ComponentUUID)
		})
		render.Template(core, func() *template.Template {
			return template.Must(template.New("page.index.html").Funcs(render.FuncMap(core)).ParseGlob("*.html"))
		})
	}
	```

=== "page.index.html"

	```html
	<html>
	    <head> <title>Home page</title> </head>
	    <body>
	        {{ template "ComponentUUID" .UUID1 }}
	        {{ template "ComponentUUID" .UUID2 }}
	    </body>
	</html>
	```

Thanks to asynchronous lifecycle, data fetching is concurrent without any goroutines hassle and page rendering happens much sooner.

!!! note
	For components rendering you can use `render` template function instead of built-it `template` function.
	In this way, you will also have an option to define own rendering logic for component with `render.Writer`.

	```html
	...
	<body>
		{{ render .UUID1 }}
		{{ render .UUID2 }}
	</body>
	```

	For more details, check [Basics • Components](/basics/components) documentation category.
