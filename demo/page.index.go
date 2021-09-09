package main

import (
	"html/template"
	"net/http"

	"github.com/yuriizinets/go-ssc"
)

type PageIndex struct {
	Navbar ssc.Component
	Hero   ssc.Component
	// Overview
	WhatIsSSC   ssc.Component
	WhyNotJS    ssc.Component
	WhyNotGoKit ssc.Component
	// Demo
	ContentUUID             ssc.Component
	ComponentUUID           ssc.Component
	ContentCounter          ssc.Component
	ComponentCounter        ssc.Component
	ContentBinding          ssc.Component
	ComponentBinding        ssc.Component
	ContentPCC              ssc.Component
	ComponentPCC            ssc.Component
	ContentEmailValidator   ssc.Component
	ComponentEmailValidator ssc.Component
	ContentRedirect         ssc.Component
	ComponentRedirect       ssc.Component
	// Help prompt
	HelpPrompt ssc.Component
	// Sponsor prompt
	SponsorPrompt ssc.Component
}

func (p *PageIndex) Template() *template.Template {
	return template.Must(template.New("page.index.html").Funcs(ssc.Funcs()).ParseGlob("*.html"))
}

func (p *PageIndex) Init() {
	// General
	p.Navbar = ssc.RegC(p, &ComponentNavbar{})
	p.Hero = ssc.RegC(p, &ComponentHero{})
	// Overview
	p.WhatIsSSC = ssc.RegC(p, &ComponentContent{
		Title: "What is SSC?",
		Description: "" +
			"First of all, it's approach of responsibilities separation. " +
			"You're separating application's view logic on different components. " +
			"It's very similar to JS frameworks, like React/Vue/Angular, but has 1 big difference - all that view logic executes on the server side and not loaded by the client. " +
			"SSC gives you a convenient approach to maintain your pages and components, including asynchronous data fetching, components reloading, server side methods, and lots more.",
	})
	p.WhyNotJS = ssc.RegC(p, &ComponentContent{
		Title: "Why not JS framework?",
		Description: "" +
			"I am trying to minimize the usage of popular SPA/PWA frameworks where it's not needed because it adds a lot of complexity and overhead. " +
			"I don't want to bring significant runtime, VirtualDOM, Webpack, and view logic into the project with minimal dynamic frontend behavior. " +
			"People start to forget that average website is not application at all. This project proves the possibility of keeping most of the logic on the server's side. ",
	})
	p.WhyNotGoKit = ssc.RegC(p, &ComponentContent{
		Title: "Why not using plain GoKit?",
		Description: "" +
			"While developing the website's frontend with plain Go and built-in library, I discovered some of the downsides of this approach. " +
			"With plain html/template you're starting to repeat yourself. It's harder to define reusable parts. " +
			"You must repeat DTO calls for each page/handler, where you're using reusable parts. " +
			"With Go's routines approach it's hard to make async-like DTO calls in the handlers. " +
			"For dynamic things, you still need to use JS and client-side DOM modification without relying on server logic. ",
	})
	// Demo
	p.ContentUUID = ssc.RegC(p, &ComponentContent{
		Title: "httpbin.org UUID",
		Description: "" +
			"Component fetches UUID from httpbin.org. " +
			"All fetches and rendering occurs on server side. " +
			"Try to reload page, and check UUID change.",
		Links: []ComponentContentLink{
			{
				Href:  "https://github.com/yuriizinets/go-ssc/blob/master/demo/component.httpbin.uuid.go",
				Title: "component.httpbin.uuid.go",
			},
			{
				Href:  "https://github.com/yuriizinets/go-ssc/blob/master/demo/component.httpbin.uuid.html",
				Title: "component.httpbin.uuid.html",
			},
		},
	})
	p.ComponentUUID = ssc.RegC(p, &ComponentHttpbinUUID{})
	p.ContentCounter = ssc.RegC(p, &ComponentContent{
		Title: "Counter",
		Description: "" +
			"Simliest example of dynamic component behavior. " +
			"Logic is calculated on the server side, client just receives html result.",
		Links: []ComponentContentLink{
			{
				Href:  "https://github.com/yuriizinets/go-ssc/blob/master/demo/component.counter.go",
				Title: "component.counter.go",
			},
			{
				Href:  "https://github.com/yuriizinets/go-ssc/blob/master/demo/component.counter.html",
				Title: "component.counter.html",
			},
		},
	})
	p.ComponentCounter = ssc.RegC(p, &ComponentCounter{})
	p.ContentBinding = ssc.RegC(p, &ComponentContent{
		Title: "Data binding",
		Description: "" +
			"Example of input binding. " +
			"As in the counter component, logic is calculated on the server side, " +
			"but reactive state data binding works only on the the client side. " +
			"Keep in mind that calling Actions is quite expensive task due to communication with server. ",
		Links: []ComponentContentLink{
			{
				Href:  "https://github.com/yuriizinets/go-ssc/blob/master/demo/component.binding.go",
				Title: "component.binding.go",
			},
			{
				Href:  "https://github.com/yuriizinets/go-ssc/blob/master/demo/component.binding.html",
				Title: "component.binding.html",
			},
		},
	})
	p.ComponentBinding = ssc.RegC(p, &ComponentBinding{})
	p.ContentPCC = ssc.RegC(p, &ComponentContent{
		Title: "Parent Component Communication",
		Description: "" +
			"Example of communication between child and parent, " +
			"that shows how to trigger parent Action from child. ",
		Links: []ComponentContentLink{
			{
				Href:  "https://github.com/yuriizinets/go-ssc/blob/master/demo/component.pcc.parent.go",
				Title: "component.pcc.parent.go",
			},
			{
				Href:  "https://github.com/yuriizinets/go-ssc/blob/master/demo/component.pcc.parent.html",
				Title: "component.pcc.parent.html",
			},
			{
				Href:  "https://github.com/yuriizinets/go-ssc/blob/master/demo/component.pcc.child.go",
				Title: "component.pcc.child.go",
			},
			{
				Href:  "https://github.com/yuriizinets/go-ssc/blob/master/demo/component.pcc.child.html",
				Title: "component.pcc.child.html",
			},
		},
	})
	p.ComponentPCC = ssc.RegC(p, &ComponentPCCParent{})
	p.ContentEmailValidator = ssc.RegC(p, &ComponentContent{
		Title: "Form Submit Action",
		Description: "" +
			"Example of dynamic component behavior. " +
			"Form is handled and processed on the server side, without page reload",
		Links: []ComponentContentLink{
			{
				Href:  "https://github.com/yuriizinets/go-ssc/blob/master/demo/component.email.validator.go",
				Title: "component.email.validator.go",
			},
			{
				Href:  "https://github.com/yuriizinets/go-ssc/blob/master/demo/component.email.validator.html",
				Title: "component.email.validator.html",
			},
		},
	})
	p.ComponentEmailValidator = ssc.RegC(p, &ComponentEmailValidator{})
	p.ContentRedirect = ssc.RegC(p, &ComponentContent{
		Title:       "Redirect Function",
		Description: "You can initiate redirect from server side code, with own custom logic",
		Links: []ComponentContentLink{
			{
				Href:  "https://github.com/yuriizinets/go-ssc/blob/master/demo/component.redirect.go",
				Title: "component.redirect.go",
			},
			{
				Href:  "https://github.com/yuriizinets/go-ssc/blob/master/demo/component.redirect.html",
				Title: "component.redirect.html",
			},
		},
	})
	p.ComponentRedirect = ssc.RegC(p, &ComponentRedirect{})
	// Help prompt
	p.HelpPrompt = ssc.RegC(p, &ComponentContent{
		Title: "Need to communicate?",
		Description: "" +
			"Examples and documentation are not good enough? " +
			"Or you have awesome proposals? " +
			"I'm always open for communication.",
		Links: []ComponentContentLink{
			{
				Href:  "https://github.com/yuriizinets/go-ssc/issues/new",
				Title: "Open new Issue",
			},
			{
				Href:  "mailto:yurii.zinets@icloud.com",
				Title: "Email",
			},
			{
				Href:  "https://t.me/yuriizinets",
				Title: "Telegram",
			},
		},
		LinksCenter: true,
	})
	// Sposnor prompt
	p.SponsorPrompt = ssc.RegC(p, &ComponentContent{
		Title: "Sponsorship",
		Description: "" +
			"It wasn't so easy to develop this idea and concept. " +
			"If you really liked it, I would be wery glad to see you as sponsor of this project! " +
			"If you don't see any convenient option, just contact my directly.",
		Links: []ComponentContentLink{
			{
				Href:  "https://blockchain.com/btc/payment_request?address=1LtNRMiPFPwQsFrkQdZzPufDzw9MQCxD3Y&amount=0.00018383&message=For SSC Development",
				Title: "Bitcoin",
			},
			{
				Href:  "https://liberapay.com/yuriizinets",
				Title: "LibrePay",
			},
			{
				Href:  "https://paypal.me/yuriizinets",
				Title: "PayPal",
			},
		},
		LinksCenter: true,
	})
}

func PageIndexHandler(rw http.ResponseWriter, r *http.Request) {
	ssc.RenderPage(rw, &PageIndex{})
}
