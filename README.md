
***I'm not Golang "rockstar"**, I'm just a regular developer. If you see any problems in the project - feel free to open new Issue. And if you would like to see features you want or just support our bright feature, please, support me with [Patreon](https://www.patreon.com/yuriizinets) or Bitcoin (39extGa1qoaGx2kpZ6RHPVkenNP9RGzVEB)*

# GoFR (Go Frontend, or Go Framework)

Micro framework, that brings frontend-like components experience to the server side with native `html/template`. Supports any serving basis (nethttp/Gin/etc), that provides `io.Writer` for response.  

**Under heavy development, not stable (!!!)**

## TOC

- [GoFR (Go Frontend, or Go Framework)](#gofr-go-frontend-or-go-framework)
  - [TOC](#toc)
  - [Why?](#why)
  - [What problems it solves? Why not plain GoKit?](#what-problems-it-solves-why-not-plain-gokit)
  - [Basic concepts](#basic-concepts)
  - [Quick start](#quick-start)
  - [Lifecycle](#lifecycle)
  - [Roadmap](#roadmap)

## Why?

Because "website" is not the same as "web application". But nowadays trends are saying otherwise. I'm trying to minimize usage of popular SPA/PWA frameworks where it's not needed at all because it adds a lot of complexity and overhead. I don't want to bring large runtime, VirtualDOM and webpack into small landing project.  

## What problems it solves? Why not plain GoKit?

While developing website's frontend with Go I realised some downsides of such approach:  

- With plain `html/template` your're starting to repeat yourself. It's harder to define reusable parts
- You must to repeat DTO calls for each page, where you're using reusable parts
- With Go's routines approach it's hard to make async-like DTO calls in the handlers

Complexity is much higher when all of them combined.

This micro framework tries to bring components and async experience to the traditional server side rendering.

## Basic concepts

- Don't replace Golang's features, that already exist
- Don't do work that's already done
- Don't bind developer with specific solutions (Gin/Chi/GORM/sqlx/etc), let developer choose

## Quick start

Basic page (on Gin basis)  
  
```go
package main

import(
    "html/template"

    "github.com/gin-gonic/gin"
    "github.com/yuriizinets/gofr"
)

// PageIndex is an implementation of gofr.Page interface
// and must to implement all required methods (even if not needed)
type PageIndex struct{}

func (*PageIndex) Template() *template.Template {
    // Template body is located in index.html
    // <html>
    //   <body>The most basic example</body>
    // </html>
    tmpl, _ := template.New("index.html").ParseGlob("*.html")
    return tmpl
}

func (*PageIndex) Meta() gofr.Meta {
    return gofr.Meta{}
}

func (p *PageIndex) Init() {}

func main() {
    g := gin.Default()

    g.GET("/", func(c *gin.Context) {
        gofr.RenderPage(c.Writer, &PageIndex{})
    })

    g.Run("localhost:25025")
}
```

Here is more complex example, that demonstrates meta usage and component definition with asynchronous httpbin data fetching. In the real-world application it's important to have separate files for pages/components/etc.

```go
package main

import (
    "encoding/json"
    "html/template"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/yuriizinets/gofr"
)

type ComponentHttpbinUUID struct {
    UUID string
}

// Component's definition, same as in the template
func (*ComponentHttpbinUUID) Definition() string {
    // Component definition is located in httpbin.uuid.html
    // and parsed in page's Template() call
    // {{ define "httpbin.uuid" }}
    //   {{ .UUID }}
    // {{ end }}
    return "httpbin.uuid"
}

// This function will be called asynchronously
// Errors processing is missed to minimize example's code
func (c *ComponentHttpbinUUID) Async() error {
    resp, _ := http.Get("http://httpbin.org/uuid")
    defer resp.Body.Close()
    var data map[string]string
    json.NewDecoder(resp.Body).Decode(&data)
    c.UUID = data["uuid"]
    return nil
}

// You can use AfterAsync to finalize data operations or whatever you want
func (c *ComponentHttpbinUUID) AfterAsync() {}

// Use gofr.Component for components
type PageIndex struct {
    UUID gofr.Component
}

func (*PageIndex) Template() *template.Template {
    // Template body is located in index.html
    // gofr.Funcs() - returns additional framework functions (we need "meta" for this case)
    // <html>
    //   <head> {{ meta . }} </head>
    //   <body> {{ template "httbin.uuid" .UUID }} </body>
    // </html>
    tmpl, _ := template.New("index.html").Funcs(gofr.Funcs()).ParseGlob("*.html")
    return tmpl
}

// Define page meta tags
func (*PageIndex) Meta() gofr.Meta {
    return gofr.Meta{
        Title: "Complex Example - GoFR",
    }
}

// Initialize page data and nested components
// Use gofr.RegisterComponent (or just gofr.RC) as a wrapper around components initialization to include this component into lifecycle
func (p *PageIndex) Init() {
    p.UUID = gofr.RegC(p, &ComponentHttpbinUUID{})
}

func main() {
    g := gin.Default()

    g.GET("/", func(c *gin.Context) {
        gofr.RenderPage(c.Writer, &PageIndex{})
    })

    g.Run("localhost:25025")
}

```

## Lifecycle

Page's lifecycle is hidden under render functions and looks like this:

- Defining shared variables (waitgroup, errors channel)
- Triggering page's `Init()` to initialize and register components
- Running all component's `Async()` functions in separate goroutines
- Wait till all asynchronous operations will be completed
- Call `AfterAsync()` for each component
- Cleaning up registered components (not needed more for internal usage)
- Getting page's template and render

## Roadmap

- [x] Basic pages interface
- [x] Basic components interface
- [x] Render pages
- [x] Components lifecycle
- [x] Basic Meta processing
- [ ] Separate RenderComponent
- [ ] Separate RenderComponentString
- [ ] Advanced Meta processing
- [ ] Separate examples projects
- [ ] Better documentation
- [ ] Dynamic Components
