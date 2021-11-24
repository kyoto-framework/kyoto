
# From scratch

This guide will only include base project setup with just index page, without details, components and examples. For advanced setup with example, check [example with guide](example-with-guide.md).

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
    You can define bootstrap function for easier template definition. For example:
    ```go
    func newtemplate(page string) *template.Template {
        return template.Must(template.New(page).Funcs(kyoto.Funcs()).ParseGlob("*.html"))
    }
    ```

## Page routing

For attaching your page, now you can simply use built-in page handler (`kyoto.PageHandler`), right bellow Routes comment in your main function.  

```go
...
mux.HandleFunc("/", kyoto.PageHandler(&PageIndex{}))
...
```

## Running

Your can run your app with usual:

```bash
go run .
```

For setting custom port, or exposing on local network, you can run in that way:

```bash
PORT=25025 go run .
```
