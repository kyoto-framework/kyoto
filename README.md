
<p align="center">
    <img src="https://raw.githubusercontent.com/yuriizinets/go-ssc/master/demo/static/ssclogo.png" width="250">
</p>
<h1 align="center">Go SSC (Go Server Side Components)</h1>

HTML render engine concept, that brings frontend-like components experience to the server side with native `html/template` on steroids. Supports any serving basis (nethttp/Gin/etc), that provides `io.Writer` for response.  

**Disclaimer 1**  
> Under heavy development, not stable **(!!!)**

**Disclaimer 2**  
> **I'm not Golang "rockstar"**, I'm just a regular developer. If you see any problems in the project - feel free to open new Issue.

## TOC

- [TOC](#toc)
- [Why?](#why)
- [What problems it solves? Why not plain GoKit?](#what-problems-it-solves-why-not-plain-gokit)
- [Zen](#zen)
- [Features](#features)
- [Quick start (simple page)](#quick-start-simple-page)
- [Basic concepts](#basic-concepts)
- [Pages](#pages)
- [Components](#components)
- [Server Side Actions](#server-side-actions)
- [Lifecycle](#lifecycle)

## Why?

I'm trying to minimize usage of popular SPA/PWA frameworks where it's not needed at all because it adds a lot of complexity and overhead. I don't want to bring large runtime, VirtualDOM and webpack into small landing project with minimal dynamic behavior.  
This project proves posibility to keep most of the logic on the server side.

## What problems it solves? Why not plain GoKit?

While developing website's frontend with Go I realised some downsides of such approach:  

- With plain `html/template` your're starting to repeat yourself. It's harder to define reusable parts
- You must to repeat DTO calls for each page, where you're using reusable parts
- With Go's routines approach it's hard to make async-like DTO calls in the handlers
- For dynamic things, you still need to use JS and client-side DOM modification

Complexity is much higher when all of them combined.

This engine tries to bring components and async experience to the traditional server side rendering.

## Zen

- Don't replace Golang's features, that already exist
- Don't do work that's already done
- Don't bind developer with specific solutions (Gin/Chi/GORM/sqlx/etc), let developer choose
- Use server for rendering, no JS specifics or client-side only behavior

## Features

- Component approach in mix with `html/template`
- Asynchronous operations
- Component methods, that can be called from client side (Server Side Actions, SSA)
- Different types of components communication (parent, cross)

## Quick start (simple page)

Basic page (on Gin basis)  
  
```go
package main

import(
    "html/template"

    "github.com/gin-gonic/gin"
    "github.com/yuriizinets/go-ssc"
)

// PageIndex is an implementation of ssc.Page interface
type PageIndex struct{}

// Template is required method. It tells about template configuration
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

Documentation not ready yet. Try to explore [demo](https://github.com/yuriizinets/go-ssc/tree/master/demo) project for features.

## Pages

Documentation not ready yet. Try to explore [demo](https://github.com/yuriizinets/go-ssc/tree/master/demo) project for features.

## Components

*Reference component is [here](https://github.com/yuriizinets/go-ssc/blob/master/demo/component.httpbin.uuid.go). Check [demo](https://github.com/yuriizinets/go-ssc/tree/master/demo) for full review.*  

To create new component, you need to do next steps:

- Create new struct with decided component name (f.e. ComponentCounter)
- Create new template with `{{ define "{name}" }} ... {{ end }}` inside. You need to use same name as struct name

To attach created component to the page, follow this steps:

- Create component field in the page struct with `ssc.Component` type
- Register component in the `Init()` method and assign it to the page struct field
- Use it in the template with passing Go component as template parameter  

Example of component, that fetches and displays UUID response from httpbin.org  

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

```html
{{ define "ComponentHttpbinUUID" }}
<div>
    <div>{{ .UUID }}</div>
</div>
{{ end }}
```

And that's how it can be attached to the index page  

```go
package main

import (
    "html/template"

    "github.com/yuriizinets/go-ssc"
)

type PageIndex struct {
    ComponentHttpbinUUID   ssc.Component
}

func (*PageIndex) Template() *template.Template {
    tmpl, _ := template.New("index.html").Funcs(ssc.Funcs()).ParseGlob("*.html")
    return tmpl
}

func (p *PageIndex) Init() {
    p.ComponentHttpbinUUID = ssc.RegC(p, &ComponentHttpbinUUID{})
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
</head>
<body>
    <h1>UUID Example</h1>
    {{ template "ComponentHttpbinUUID" .ComponentHttpbinUUID }}
</body>
</html>

```

## Server Side Actions

Documentation not ready yet. Try to explore [demo](https://github.com/yuriizinets/go-ssc/tree/master/demo) project for features.

## Lifecycle

Page's lifecycle is hidden under render functions and looks like this:

- Defining shared variables (waitgroup, errors channel)
- Triggering page's `Init()` to initialize and register components
- Running all component's `Async()` functions in separate goroutines
- Wait till all asynchronous operations will be completed
- Call `AfterAsync()` for each component
- Cleaning up registered components (not needed more for internal usage)
- Getting page's template and render
