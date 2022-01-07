
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
	Extendable Go library for creating fast, SSR-first frontend avoiding plain <code>html/template</code> approach downsides.
</p>


> **Disclaimer**  
> This project is in early development, don't use in production! In case of any issues/proposals, feel free to open an [issue](https://github.com/kyoto-framework/kyoto/issues/new)

## Why I need this?

- Get rid of spaghetti inside of handlers
- Organize code into configurable components structure
- Simple and straightforward page rendering lifecycle
- Asynchronous DTO without goroutine mess
- Built-in dynamics like Hotwire or Laravel Livewire
- Everything on top of well-known `html/template`
- Get control over project setup: 0 external dependencies, just `kyoto` itself

## Why not?

- In active development (not production ready)
- Not situable for pretty dynamic frontends
- You want to develop SPA/PWA
- You're just feeling OK with JS frameworks

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

<a target="_blank" href="https://www.buymeacoffee.com/yuriizinets"><img alt="Buy me a Coffee" src="https://github.com/egonelbre/gophers/blob/master/.thumb/animation/buy-morning-coffee-3x.gif?raw=true"></a>


Or directly with Bitcoin: `bc1qgxe4u799f8pdyzk65sqpq28xj0yc6g05ckhvkk`
