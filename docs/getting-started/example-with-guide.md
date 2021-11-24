
# Example with guide

This guide is extended version of "[from scratch](from-scratch.md)" documentation and will show minimal setup with page, multiple component instances, lifecycle integration and `net/http` setup. This guide will rely on demo project setup, that can be found [here](https://github.com/yuriizinets/kyoto/tree/master/.demo).  

## Entry point

First, we need to setup serving basis.  

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
Page is represented by a struct which implements `Page` interface.
Page requires method, returning ready for use template. In this example, we will store our page markup in `page.index.html`.
`kyoto.Funcs` is a function, that returns FuncMap. This funcmap is required for correct work of some `kyoto` features.

```go title="page.index.go"
package main

import (
    "html/template"
    "github.com/yuriizinets/kyoto"
)

type PageIndex struct {}

func (p *PageIndex) Template() *template.Template {
    return template.Must(template.New("page.index.html").Funcs(kyoto.Funcs()).ParseGlob("*.html"))
}
```

!!! note
    You can define bootstrap function for easier template definition. For example:
    ```go
    func newtemplate(page string) *template.Template {
        return template.Must(template.New(page).Funcs(kyoto.Funcs()).ParseGlob("*.html"))
    }
    ```

```html title="page.index.html"
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Kyoto Quick Start</title>
</head>

<body>
    ...
</body>
</html>
```

## Component

Let's define sample component, which fetches UUID from httpbin page.  
Component is represented by a struct which implements `Component` interface.
By default, Component interface doesn't have any required methods. Instead of having all-in-one, we have multiple interfaces with functionality separation.
This approach also covers pages. In this example, we will implement `ImplementsAsync` interface.
This method will be called as goroutine in page rendering lifecycle.
In that way, all needed async data will be fetched concurrently. In this example, component's markup will be stored in `component.uuid.html`

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
<div>
    httpbin.org uuid: {{ .UUID }}
</div>
{{ end }}
```

## Attaching component

For using component, you need to define page fields for storing component objects and `Init` method for initialization and registration of components.
Inside of init, use `kyoto.RegC` for registering your components. In that way you're including component in page render lifecycle.
After that, you need to pass component object to template in your page markup.

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
    {{ template "ComponentUUID" .DemoUUID1 }}
    {{ template "ComponentUUID" .DemoUUID2 }}
    {{ template "ComponentUUID" .DemoUUID3 }}
    {{ template "ComponentUUID" .DemoUUID4 }}
</body>
...
```

## Page routing

For attaching your page, now you can simply use built-in page handler (`kyoto.PageHandler`), bellow `Routes` comment in your main function.

```go
...
mux.HandleFunc("/", kyoto.PageHandler(&PageIndex{}))
...
```

## Running

Ready! Now your can run your app with usual:

```bash
go run .
```

For setting custom port, or exposing on local network, you can run in that way:

```bash
PORT=25025 go run .
```
