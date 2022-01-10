
<p align="center">
    <img width="400" src="https://raw.githubusercontent.com/kyoto-framework/kyoto/master/docs/assets/kyoto.svg" />
</p>

<h1 align="center">kyoto</h1>

<p align="center">
	<img src="https://img.shields.io/github/license/kyoto-framework/kyoto">
	<img src="https://goreportcard.com/badge/github.com/kyoto-framework/kyoto">
	<img src="https://pkg.go.dev/badge/github.com/kyoto-framework/kyoto.svg">
</p>

<p align="center">
	Extensible Go library for creating fast, SSR-first frontend avoiding vanilla templating downsides.
</p>

> **Disclaimer №1**  
> High entry threshold. You must understand what problems are you trying to solve before using.

> **Disclaimer №2**  
> This project is in early development, don't use in production! In case of any issues/proposals, feel free to open an [issue](https://github.com/kyoto-framework/kyoto/issues/new)

## Motivation

When our team migrated from JS framework to vanilla Go templating (we had a lot of reasons), we faced a set of tipical inconveniences during development. Copy-paste of DTO calls, markups, spaghetti inside of handlers, etc. Also, the goroutine system might be quite verbose in some places. That's why we decided to bring the best from two worlds and create a small library.

## What kyoto proposes?

- Organize code into configurable and standalone components structure
- Get rid of spaghetti inside of handlers
- Simple asynchronous page rendering lifecycle
- Built-in dynamics like Hotwire or Laravel Livewire
- Built-in rendering based on `html/template`
- Full control over project setup (minimal dependencies)
- Extensible architecture (everyone can create own extensions for library)

## Reasons to opt out

- In active development (not production ready)
- You want to develop SPA/PWA
- You're just feeling OK with JS frameworks
- Not situable for a frontend with a lot of client-side logic

## Installation

As simple as `go get github.com/kyoto-framework/kyoto`  
Check documentation page for quick start: [https://kyoto.codes/getting-started/](https://kyoto.codes/getting-started/)

## Usage

Kyoto project setup may seem complicated and unusual at first sight.  
It's highly recommended to follow documentation while using library: [https://kyoto.codes/getting-started/](https://kyoto.codes/getting-started/)  

This example is not completely independent and just shows what the code looks like when using kyoto:

```go
package main

import (
	"html/template"

	"github.com/kyoto-framework/kyoto"
	"github.com/kyoto-framework/uikit/twui"
)

type PageIndex struct {
	Navbar kyoto.Component
}

func (p *PageIndex) Template() *template.Template {
	return mktemplate("page.index.html")
}

func (p *PageIndex) Init() {
	p.Navbar = kyoto.RegC(p, &twui.AppUINavNavbar{
		Logo: `<img src="/static/img/kyoto.svg" class="h-8 w-8 scale-150" />`,
		Links: []twui.AppUINavNavbarLink{
			{Text: "Kyoto", Href: "https://github.com/kyoto-framework/kyoto"},
			{Text: "UIKit", Href: "https://github.com/kyoto-framework/uikit"},
			{Text: "Charts", Href: "https://github.com/kyoto-framework/kyoto-charts"},
			{Text: "Starter", Href: "https://github.com/kyoto-framework/starter"},
		},
		Profile: twui.AppUINavNavbarProfile{
			Enabled: true,
			Avatar: `
					<svg class="w-6 h-6 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"></path></svg>
				`,
			Links: []twui.AppUINavNavbarLink{
				{Text: "GitHub", Href: "https://github.com/kyoto-framework/kyoto/discussions/40"},
				{Text: "Telegram", Href: "https://t.me/yuriizinets"},
				{Text: "Email", Href: "mailto:yurii.zinets@icloud.com"},
			},
		},
    })
}

```

## References

Documentation: [https://kyoto.codes/](https://kyoto.codes/)  
UIKit: [https://github.com/kyoto-framework/uikit](https://github.com/kyoto-framework/uikit)  
Demo project, Hacker News client made with kyoto: [https://hn.kyoto.codes/](https://hn.kyoto.codes/)  
Demo project, features overview: [https://github.com/kyoto-framework/kyoto/tree/master/examples/demo](https://github.com/kyoto-framework/kyoto/tree/master/examples/demo)  

## Support

Any project support is appreciated! Donations will help us to keep high updates frequency. If you would like to avoid using listed methods, contact us directly with [info@kyoto.codes](mailto:info@kyoto.codes)  

Bitcoin: `bc1qgxe4u799f8pdyzk65sqpq28xj0yc6g05ckhvkk`  
Ethereum: `0xEB2f24e830223bE081264e0c81fb5FD4DDD2B7B0`

Open Collective: [https://opencollective.com/kyoto-framework](https://opencollective.com/kyoto-framework)
