# Example with Guide

This guide is an extended version of the "[From Scratch](from-scratch.md)" page and will show a minimal setup with: a page, multiple component instances, lifecycle integration and `net/http` setup. This guide will rely on the demo project setup found [here](https://github.com/kyoto-framework/kyoto/tree/master/examples/demo).

## Entry Point

Firstly, we need to setup the serving foundations.

```go title="main.go"
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
A page is represented by a struct which implements `Page` interface.
The Page's required method is for returning a ready-to-use template. In this example, we will store our mark-up in `page.index.html`.
`kyoto.Funcs` is a function that returns FuncMap. This funcmap is required for the correct working of some `kyoto` features.

```go title="page.index.go"
package main

import (
    "html/template"
    "github.com/kyoto-framework/kyoto"
)

type PageIndex struct {}

func (p *PageIndex) Template() *template.Template {
    return template.Must(template.New("page.index.html").Funcs(kyoto.Funcs()).ParseGlob("*.html"))
}
```

!!! note
You can define bootstrap functions for easier template definitions. For example:
`go func newtemplate(page string) *template.Template { return template.Must(template.New(page).Funcs(kyoto.Funcs()).ParseGlob("*.html")) } `

```html title="page.index.html"
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Kyoto Quick Start</title>
  </head>

  <body>
    ...
  </body>
</html>
```

## Component

Let's define a sample component, which fetches a UUID from httpbin page.  
The component is represented by a struct which implements `Component` interface.
By default, the Component interface doesn't have any required methods. Instead of having all-in-one, we have multiple interfaces with separate functionality.
This approach also applies to pages. In this example, we will implement `ImplementsAsync` interface.
This method will be called as a goroutine in the page rendering lifecycle.
In that way, all needed async data will be fetched concurrently. In this example, component's mark-up will be stored in `component.uuid.html`

```go title="component.uuid.go"
package main

import (
    "net/http"
    "encoding/json"
)

type ComponentUUID struct {
    UUID string
}

func (c *ComponentUUID) Async() error {
    resp, err := http.Get("http://httpbin.org/uuid")
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    data := map[string]string{}
    json.NewDecoder(resp.Body).Decode(&data)
    c.UUID = data["uuid"]
    return nil
}
```

```html title="component.uuid.html"
{{ define "ComponentUUID" }}
<div>httpbin.org uuid: {{ .UUID }}</div>
{{ end }}
```

## Attaching Component

For using the component, you need to define some page fields for storing component objects and an `Init` method for initialization and registration of components.
Inside of init, use `kyoto.RegC` to register your components. This will include the component in the page render lifecycle.
After that you need to pass the component to a template in your page mark-up.

```go title="page.index.html"
...
type PageIndex struct {
    DemoUUID1 kyoto.Component
    DemoUUID2 kyoto.Component
    DemoUUID3 kyoto.Component
    DemoUUID4 kyoto.Component
}

...

func (p *PageIndex) Init() {
    p.DemoUUID1 = kyoto.RegC(p, &ComponentUUID{})
    p.DemoUUID2 = kyoto.RegC(p, &ComponentUUID{})
    p.DemoUUID3 = kyoto.RegC(p, &ComponentUUID{})
    p.DemoUUID4 = kyoto.RegC(p, &ComponentUUID{})
}
```

```html title="page.index.html"
...
<body>
  {{ template "ComponentUUID" .DemoUUID1 }} {{ template "ComponentUUID"
  .DemoUUID2 }} {{ template "ComponentUUID" .DemoUUID3 }} {{ template
  "ComponentUUID" .DemoUUID4 }}
</body>
...
```

## Page Routing

For attaching your page you can simply use the built-in page handler (`kyoto.PageHandler`) found below the `Routes` comment in your main function.

```go
...
mux.HandleFunc("/", kyoto.PageHandler(&PageIndex{}))
...
```

## Running

Ready! Great! Now your can run your app with the usual:

```bash
go run .
```

For setting custom ports or exposing on a local network, you can run the following:

```bash
PORT=25025 go run .
```
