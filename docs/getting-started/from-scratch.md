# From Scratch

This guide will include a basic project setup with just a basic index page, components and examples. For advanced setup with examples, check [Example with Guide](example-with-guide.md).

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
    ```go
    func newtemplate(page string) *template.Template {
        return template.Must(template.New(page).Funcs(kyoto.Funcs()).ParseGlob("*.html"))
    }
    ```

## Page routing

For attaching your page, you can simply use the built-in page handler (`kyoto.PageHandler`), right below the Routes comment in your main function.

```go
...
mux.HandleFunc("/", kyoto.PageHandler(&PageIndex{}))
...
```

## Running

Your can run your app with the usual:

```bash
go run .
```

For setting custom ports or exposing on a local network, you can run with the following:

```bash
PORT=25025 go run .
```
