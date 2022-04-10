
# Alternative rendering

Kyoto ships with an ability to define own rendering, avoiding built-in `html/template`.
To use this option, you will need to use `render.Writer` function, 
provide own rendering to `io.Writer` and use `render` template function to render a component.

This example shows how to define custom rendering for a single component
and use this component in pair with classic templates.

=== "component.uuid.go"

    ```go
    package main

	import (
	    "net/http"
	    "encoding/json"
        "io"

	    "github.com/kyoto-framework/kyoto"
	    "github.com/kyoto-framework/kyoto/render"
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
        render.Writer(core, func(w io.Writer) error {
            fmt.Sprintf(w, `
                <div>
                    Component UUID: %s
                </div>
            `, core.State.Get("UUID"))
        })
	}
    ```

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
	<!DOCTYPE html>
	<html>
	<head>
	    <title>Page Index</title>
	</head>
	<body>
	    {{ render .UUID1 }}
	    {{ render .UUID2 }}
	</body>
	</html>
	```

You can use `render` function also with a classic templates too.  
First, this function tries to find a custom rendering function. 
If not found, it dynamicaly resolves a component name and doing inline dynamic rendering with a classic template.
In this way you can simplify your code, ommit writing a template definition name and your rendering becomes dynamic.

!!! note ""
    With `render` function you can avoid an error while calling other template with a dynamic name.  
    Wondering what I'm talking about?  
    Check this: [https://stackoverflow.com/questions/20716726/call-other-templates-with-dynamic-name](https://stackoverflow.com/questions/20716726/call-other-templates-with-dynamic-name)