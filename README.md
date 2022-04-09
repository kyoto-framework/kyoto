
<p align="center">
    <img width="200" src="https://raw.githubusercontent.com/kyoto-framework/kyoto/master/docs/v1/docs/assets/kyoto.svg" />
</p>

<h1 align="center">kyoto</h1>

<p align="center">
    Extensible Go library for creating fast, SSR-first frontend avoiding vanilla templating downsides.
</p>

<p align="center">
    <img src="https://img.shields.io/github/license/kyoto-framework/kyoto">
    <img src="https://goreportcard.com/badge/github.com/kyoto-framework/kyoto">
    <img src="https://pkg.go.dev/badge/github.com/kyoto-framework/kyoto.svg">
</p>

<p align="center">
    <a href="https://v1.kyoto.codes/getting-started/">Getting started</a>&nbsp;&bull; <a href="https://v1.kyoto.codes/basics/">Basics</a>&nbsp;&bull; <a href="https://v1.kyoto.codes/features/">Features</a>&nbsp;&bull; <a href="https://v0.kyoto.codes/">v0 Documentation</a>&nbsp;&bull; <a href="https://github.com/kyoto-framework/kyoto#donate">Donate</a>
</p>

> **Disclaimers**  
> - High entry threshold. You must understand what problems are you trying to solve before using.
> - Due to current situation in Ukraine development has slowed down a bit, but it won't stop. Слава Україні 🇺🇦
> - This project is in early development, don't use in production! In case of any issues/proposals, feel free to open an [issue](https://github.com/kyoto-framework/kyoto/issues/new)

## Motivation

When our team migrated from JS framework to vanilla Go templating (we had a lot of reasons), we faced a set of tipical inconveniences during development. That's why we decided to bring the best from two worlds and create this small library.

## What kyoto proposes?

- Organize code into configurable and standalone components structure
- Get rid of spaghetti inside of handlers
- Simple asynchronous page rendering lifecycle
- Built-in dynamics like Hotwire or Laravel Livewire
- Built-in rendering based on `html/template`
- Full control over project setup (minimal dependencies)
- 0kb JS payload without dynamics (~8kb when using dynamics)
- Extensible architecture (everyone can create own extensions for library)

## Reasons to opt out

- In active development (not production ready)
- You want to develop SPA/PWA
- You're just feeling OK with JS frameworks
- Not situable for a frontend with a lot of client-side logic

## Installation

As simple as `go get github.com/kyoto-framework/kyoto@master`  
Check latest version of documentation page for quick start: [https://v1.kyoto.codes/getting-started/](https://v1.kyoto.codes/getting-started/)

## Usage

Kyoto project setup may seem complicated and unusual at first sight.  
It's highly recommended to follow documentation while using library: [https://v1.kyoto.codes/getting-started/](https://v1.kyoto.codes/getting-started/)  

```go
package main

import (
    "html/template"
    "encoding/json"
    "net/http"

    "github.com/kyoto-framework/kyoto"
    "github.com/kyoto-framework/kyoto/render"
    "github.com/kyoto-framework/kyoto/lifecycle"
)

// This example demonstrates main advantage of kyoto library - asynchronous lifecycle.
// Multiple UUIDs will be fetched in asynchronous way, without even touching goroutines and synchronization tools like sync.WaitGroup.

// Let's assume markup of this component is stored in 'component.uuid.html'
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

// Let's assume markup of this page is stored in 'page.index.html'
func PageIndex(core *kyoto.Core) {
    lifecycle.Init(core, func() {
        core.Component("UUID1", ComponentUUID)
        core.Component("UUID2", ComponentUUID)
    })
    render.Template(c, func() *template.Template {
        return template.Must(template.New("page.index.html").Funcs(render.FuncMap(core)).ParseGlob("*.html"))
    })
}

func main() {
    http.HandleFunc("/", render.PageHandler(PageIndex))
    http.ListenAndServe(":8080", nil)
}

```

## References

Documentation: [https://v1.kyoto.codes/](https://v1.kyoto.codes/)  
UIKit: [https://github.com/kyoto-framework/uikit](https://github.com/kyoto-framework/uikit)  
Demo project, Hacker News client made with kyoto: [https://hn.kyoto.codes/](https://hn.kyoto.codes/)  
Demo project, features overview: [https://github.com/kyoto-framework/kyoto/tree/master/examples/demo](https://github.com/kyoto-framework/kyoto/tree/master/examples/demo)  

## Donate

Any project support is appreciated! Donations will help us to keep high updates frequency. If you would like to avoid using listed methods, contact us directly with [info@kyoto.codes](mailto:info@kyoto.codes)  

Bitcoin: `bc1qgxe4u799f8pdyzk65sqpq28xj0yc6g05ckhvkk`  
Ethereum: `0xEB2f24e830223bE081264e0c81fb5FD4DDD2B7B0`

Open Collective: [https://opencollective.com/kyoto-framework](https://opencollective.com/kyoto-framework)
