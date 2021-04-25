
<p align="center">
    <img src="https://raw.githubusercontent.com/yuriizinets/go-ssc/master/demo/static/ssclogo.png" width="250">
</p>
<h1 align="center">Go SSC (Go Server Side Components)</h1>

An HTML render engine concept that brings frontend-like components experience to the server side with native `html/template` on steroids. Supports any serving basis (nethttp/Gin/etc), that provides `io.Writer` in response.  

**Disclaimer 1**  
> Under heavy development, not stable **(!!!)**

**Disclaimer 2**  
> **I'm not Golang "rockstar"**, and code may be not so good quality as you may expect. If you see any problems in the project - feel free to open new Issue.

## TOC

- [TOC](#toc)
- [Why?](#why)
- [What problems it solves? Why not plain GoKit?](#what-problems-it-solves-why-not-plain-gokit)
- [Zen](#zen)
- [Features](#features)
- [Quick start (simple page)](#quick-start-simple-page)
- [Basic concepts](#basic-concepts)
  - [Lifecycle](#lifecycle)
- [Pages](#pages)
  - [Example of page](#example-of-page)
- [Components](#components)
  - [Component example](#component-example)
- [Server Side Actions](#server-side-actions)

## Why?

I am trying to minimize the usage of popular SPA/PWA frameworks where it's not needed because it adds a lot of complexity and overhead. I don't want to bring significant runtime, VirtualDOM, and Webpack into the project with minimal dynamic frontend behavior. 

This project proves the possibility of keeping most of the logic on the server's side.

## What problems does it solve? Why not using plain GoKit?

While developing the website's frontend with Go, I discovered some of the downsides of this approach:

- With plain html/template you're starting to repeat yourself. It's harder to define reusable parts.
- You must repeat DTO calls for each page, where you're using reusable parts.
- With Go's routines approach it's hard to make async-like DTO calls in the handlers.
- For dynamic things, you still need to use JS and client-side DOM modification.

Complexity is much higher when all of them get combined.

This engine tries to bring components and async experience to the traditional server-side rendering.

## Zen

- Don't replace Go features that exist already.
- Don't do work that's already done
- Don't force developers to use a specific solution (Gin/Chi/GORM/sqlx/etc). Let them choose
- Rely on the server to do the rendering, no JS specifics or client-side only behavior

## Features

- Component approach in mix with `html/template`
- Asynchronous operations
- Component methods that can be called from the client side (Server Side Actions, SSA)
- Different types of component communication (parent, cross)

## Quick start (simple page)

Basic page (based on Gin)  
  
```go
package main

import(
    "html/template"

    "github.com/gin-gonic/gin"
    "github.com/yuriizinets/go-ssc"
)

// PageIndex is an implementation of ssc.Page interface
type PageIndex struct{}

// Template is a required page method. It tells about template configuration
func (*PageIndex) Template() *template.Template {
    // Template body is located in index.html
    // <html>
    //   <body>The most basic example</body>
    // </html>
    tmpl, _ := template.New("index.html").ParseGlob("*.html")
    return tmpl
}

func main() {
    g := gin.Default()

    g.GET("/", func(c *gin.Context) {
        ssc.RenderPage(c.Writer, &PageIndex{})
    })

    g.Run("localhost:25025")
}
```

## Basic concepts

Each page or component is represented by its own structure. For implementing specific functionality, you can use structure's methods with a predefined declaration (f.e. `Init(p ssc.Page)`). You need to follow declaration rules to match the interfaces required (you can find all interfaces in `types.go`).  
Before implementing any method, you need to understand the rendering lifecycle.

### Lifecycle

Each page's lifecycle is hidden under the render function and follows this steps:

- Defining shared variables (waitgroup, errors channel)
- Triggering the page's `Init()` to initialize and register components
- Running all component's `Async()` functions in separate goroutines
- Waiting untill all asynchronous operations are completed
- Calling `AfterAsync()` for each component
- Cleaning up registered components (not needed more for internal usage)
- Getting page's template and render

> Even though methods like `Init()` or `Async()` can handle your business logic like forms processing, please, try to avoid that. Keep your app's business logic inside tje handlers, and use this library only for page rendering.

## Pages

To implement a page, you need to declare its structure with `Template() *template.Template` method. This is the only requirements. Also, each page has these optional methods:

- `Init()` - used to initialize page, f.e. components registering or providing default values
- `Meta() ssc.Meta` - used to provide advanced page meta, like title, description, hreflangs, etc.

### Example of page

*Reference page is [here](https://github.com/yuriizinets/go-ssc/blob/master/demo/page.index.go). Check [demo](https://github.com/yuriizinets/go-ssc/tree/master/demo) for full review.*  

```go
package main

import (
    "html/template"

    "github.com/yuriizinets/go-ssc"
)

type PageIndex struct {
    ComponentHttpbinUUID   ssc.Component
    ComponentCounter       ssc.Component
    ComponentSampleBinding ssc.Component
    ComponentSampleParent  ssc.Component
}

func (*PageIndex) Template() *template.Template {
    return template.Must(template.New("page.index.html").Funcs(funcmap()).ParseGlob("*.html"))
}

func (p *PageIndex) Init() {
    p.ComponentHttpbinUUID = ssc.RegC(p, &ComponentHttpbinUUID{})
    p.ComponentCounter = ssc.RegC(p, &ComponentCounter{})
    p.ComponentSampleBinding = ssc.RegC(p, &ComponentSampleBinding{})
    p.ComponentSampleParent = ssc.RegC(p, &ComponentSampleParent{})
}

func (*PageIndex) Meta() ssc.Meta {
    return ssc.Meta{
        Title: "SSC Example",
    }
}

```

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    {{ meta . }}
    {{ dynamics }}
    <link href="https://unpkg.com/tailwindcss@^2/dist/tailwind.min.css" rel="stylesheet">
</head>
<body>
    <div class="pt-24"></div>
    <div>
        <img src="/static/ssclogo.png" alt="Logo" width="400" class="m-auto">
    </div>
    <h1 class="mt-4 text-5xl text-center">Go SSC Demo Page</h1>
    <div class="pt-16"></div>
    <h2 class="text-3xl text-center">Httpbin UUID</h2>
    <p class="text-center">UUID, fetched on the server side, asynchronously, from httpbin.org</p>
    <div class="mt-2 text-center">
        {{ template "ComponentHttpbinUUID" .ComponentHttpbinUUID }}
    </div>
    <div class="pt-16"></div>
    <h2 class="text-3xl text-center">Counter</h2>
    <p class="text-center">Counter, fully implemented on server side (Server Side Actions, SSA)</p>
    <div class="mt-2 text-center">
        {{ template "ComponentCounter" .ComponentCounter }}
    </div>
    <div class="pt-16"></div>
    <h2 class="text-3xl text-center">Binding</h2>
    <p class="text-center">Demo of Client Side state binding with Server Side calculation</p>
    <div class="mt-2 text-center">
        {{ template "ComponentSampleBinding" .ComponentSampleBinding }}
    </div>
    <div class="pt-16"></div>
    <h2 class="text-3xl text-center">Parent Component Communication</h2>
    <p class="text-center">SSA, triggered from child component</p>
    <div class="mt-2 text-center">
        {{ template "ComponentSampleParent" .ComponentSampleParent }}
    </div>
    <div class="pt-16"></div>
</body>
</html>

```

## Components

To implement a component, you just need to declare its structure. There are no requirements for declaring a component. Also, each component has these optional methods:

- `Init(p ssc.Page)` - used to initialize component, f.e. nested components registering or providing default values
- `Async() error` - method is called asynchronously with goroutines and processed concurrently during lifecycle. You can use it for fetching information from DB or API
- `AfterAsync()` - method is called after all finishing all async operations
- `Actions() ActionsMap` - used for providing SSA. Check [Server Side Actions](#server-side-actions) for details

### Component example

*Reference component is [here](https://github.com/yuriizinets/go-ssc/blob/master/demo/component.httpbin.uuid.go). Check [demo](https://github.com/yuriizinets/go-ssc/tree/master/demo) for full review.*  

Example of a component that fetches and displays UUID response from httpbin.org  

```go
package main

import (
    "io/ioutil"
    "net/http"

    "github.com/yuriizinets/go-ssc"
)

type ComponentHttpbinUUID struct {
    UUID string
}

// Async method is handled by library under the hood
// Each async method is called asynchronously with goroutines and processed concurrently
func (c *ComponentHttpbinUUID) Async() error {
    resp, err := http.Get("http://httpbin.org/uuid")
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    data, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    c.UUID = string(data)
    return nil
}
```

For component usage you can check [example of page](#example-of-page).

## Server Side Actions

The documentation is not ready yet. Try to explore [the demo](https://github.com/yuriizinets/go-ssc/tree/master/demo) project for gettting familiar with all the features.
