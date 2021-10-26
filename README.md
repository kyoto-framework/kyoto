
<p align="center">
    <img width="400" src="https://raw.githubusercontent.com/yuriizinets/kyoto/master/.web/docs/.vuepress/public/kyoto.svg" />
</p>

<h1 align="center">kyoto</h1>
<p align="center">
    An HTML render engine concept that brings frontend-like components experience to the server side with native <code>html/template</code> on steroids. Supports any serving basis (<code>net/http</code>/<code>gin</code>/etc), that provides <code>io.Writer</code> in response.
</p>

![License](https://img.shields.io/github/license/yuriizinets/kyoto)
[![Go Report Card](https://goreportcard.com/badge/github.com/yuriizinets/kyoto)](https://goreportcard.com/report/github.com/yuriizinets/kyoto)
[![Go Reference](https://pkg.go.dev/badge/github.com/yuriizinets/kyoto.svg)](https://pkg.go.dev/github.com/yuriizinets/kyoto)



> **Disclaimer**  
> This project in early development, don't use in production! In case of any issues/proposals, feel free to open an [issue](https://github.com/yuriizinets/kyoto/issues/new)

## Installation

As simple as `go get github.com/yuriizinets/kyoto`  
Check documentation page for quick start: [https://kyoto.codes/docs/get-started.html](https://kyoto.codes/docs/get-started.html)

## Usage

Kyoto project setup may seem complicated and unusual at first sight.  
It's highly recommended to follow documentation while using library: [https://kyoto.codes/docs/get-started.html](https://kyoto.codes/docs/get-started.html)  

This example is not completely independent and just shows what the code looks like when using kyoto:

```go
package main

import (
	"html/template"

	"github.com/yuriizinets/kyoto"
	"github.com/yuriizinets/kyoto-uikit/twui"
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
			{Text: "Kyoto", Href: "https://github.com/yuriizinets/kyoto"},
			{Text: "UIKit", Href: "https://github.com/yuriizinets/kyoto-uikit"},
			{Text: "Charts", Href: "https://github.com/yuriizinets/kyoto-charts"},
			{Text: "Starter", Href: "https://github.com/yuriizinets/kyoto-starter"},
		},
		Profile: twui.AppUINavNavbarProfile{
			Enabled: true,
			Avatar: `
					<svg class="w-6 h-6 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"></path></svg>
				`,
			Links: []twui.AppUINavNavbarLink{
				{Text: "GitHub", Href: "https://github.com/yuriizinets/kyoto/discussions/40"},
				{Text: "Telegram", Href: "https://t.me/yuriizinets"},
				{Text: "Email", Href: "mailto:yurii.zinets@icloud.com"},
			},
		},
    })
}

```


## References

Project main page: https://kyoto.codes/  
Documentation: https://kyoto.codes/docs/  
UIKit: https://github.com/yuriizinets/kyoto-uikit  
Demo project, Hacker News client made with kyoto: https://hn.kyoto.codes/  
Demo project, features overview: https://github.com/yuriizinets/kyoto/tree/master/.demo  


## Support

<a target="_blank" href="https://www.buymeacoffee.com/yuriizinets"><img alt="Buy me a Coffee" src="https://github.com/egonelbre/gophers/blob/master/.thumb/animation/buy-morning-coffee-3x.gif?raw=true"></a>


Or directly with Bitcoin: `18WY5KRWVKVxjWKFzJLqqALsRrsDh4snqg`
