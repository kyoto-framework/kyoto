<p align="center">
    <img width="200" src="https://raw.githubusercontent.com/kyoto-framework/kyoto/master/docs/v1/docs/assets/kyoto.svg" />
</p>

<h1 align="center">kyoto</h1>

<p align="center">
    Go library for creating fast, SSR-first frontend avoiding vanilla templating downsides.
</p>

<p align="center">
    <a href="https://goreportcard.com/report/github.com/kyoto-framework/kyoto">
        <img src="https://goreportcard.com/badge/github.com/kyoto-framework/kyoto">
    </a>
    <a href="https://codecov.io/gh/kyoto-framework/kyoto">
        <img src="https://codecov.io/gh/kyoto-framework/kyoto/branch/master/graph/badge.svg?token=XVLKT20DP8">
    </a href="https://pkg.go.dev/github.com/kyoto-framework/kyoto">
        <img src="https://pkg.go.dev/badge/github.com/kyoto-framework/kyoto.svg">
    </a>
    <a href="https://opencollective.com/kyoto-framework">
        <img src="https://img.shields.io/opencollective/all/kyoto-framework?label=backers%20%26%20sponsors">
    </a>
    <img src="https://img.shields.io/github/license/kyoto-framework/kyoto">
</p>

<p align="center">
    <a href="https://kyoto.codes/getting-started/">Getting started</a>&nbsp;&bull; <a href="https://kyoto.codes/basics/">Basics</a>&nbsp;&bull; <a href="https://kyoto.codes/features/">Features</a>&nbsp;&bull; <a href="https://v0.kyoto.codes/">v0 Documentation</a>&nbsp;&bull; <a href="https://github.com/kyoto-framework/kyoto#donate">Donate</a>
</p>

## Motivation

Creating asynchronous and dynamic layout parts is a complex problem for larger projects using `html/template`.
Library tries to simplify this process.

## What kyoto proposes?

- Organize code into configurable and standalone components structure
- Get rid of spaghetti inside of handlers
- Simple asynchronous rendering lifecycle
- Built-in dynamics like Hotwire or Laravel Livewire
- Built-in rendering based on `html/template`
- Full control over project setup (minimal dependencies)
- 0kb JS payload without actions client (~8kb when using actions)

## Reasons to opt out

- API may change drastically between major versions
- You want to develop SPA/PWA
- You're just feeling OK with JS frameworks
- Not situable for a frontend with a lot of client-side logic

## Installation

As simple as `go get github.com/kyoto-framework/kyoto@master`  
Check latest version of documentation page for quick start: [https://kyoto.codes/getting-started/](https://kyoto.codes/getting-started/)

## Usage

Kyoto project setup may seem complicated and unusual at first sight.  
It's highly recommended to follow documentation while using library: [https://kyoto.codes/getting-started/](https://kyoto.codes/getting-started/)  

```go
package main

import (
    "html/template"
    "encoding/json"
    "net/http"

    "github.com/kyoto-framework/kyoto"
)

// This example demonstrates main advantage of kyoto library - asynchronous lifecycle.
// Multiple UUIDs will be fetched in asynchronous way, without explicitly touching goroutines 
// and synchronization tools like sync.WaitGroup.

type CUUIDState struct {
    UUID string
}

// Let's assume markup of this component is stored in 'component.uuid.html'
func CUUID(ctx *kyoto.Context) (state CUUIDState) {
    // Fetch uuid data
    resp, _ := http.Get("http://httpbin.org/uuid")
    data := map[string]string{}
    json.NewDecoder(resp.Body).Decode(&data)
    // Set state
    state.UUID = data["uuid"]
}

type PIndexState struct {
    UUID1 kyoto.Component[CUUIDState]
    UUID1 kyoto.Component[CUUIDState]
}

// Let's assume markup of this page is stored in 'page.index.html'
func PIndex(ctx *kyoto.Context) (state PIndexState) {
    // Define rendering
    render.Template(ctx, "page.index.html")
    // Attach components
    state.UUID1 = kyoto.Use(ctx, CUUID)
    state.UUID2 = kyoto.Use(ctx, CUUID)
}

func main() {
    // Register page
    kyoto.HandlePage("/", PIndex)
    // Serve
    kyoto.Serve(":8080")
}

```

## Donate

Any project support is appreciated! Donations will help us to keep high updates frequency. If you would like to avoid using listed methods, contact us directly with [info@kyoto.codes](mailto:info@kyoto.codes)  

Bitcoin: `bc1qgxe4u799f8pdyzk65sqpq28xj0yc6g05ckhvkk`  
Ethereum: `0xEB2f24e830223bE081264e0c81fb5FD4DDD2B7B0`

Open Collective: [https://opencollective.com/kyoto-framework](https://opencollective.com/kyoto-framework)