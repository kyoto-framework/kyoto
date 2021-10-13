package main

import (
	"html/template"

	ssc "github.com/yuriizinets/ssceng"
)

type PageIndex struct {
	Hero                          ssc.Component
	DescriptionDemoUUID           ssc.Component
	DemoUUID                      ssc.Component
	DescriptionDemoCounter        ssc.Component
	DemoCounter                   ssc.Component
	DescriptionDemoCalc           ssc.Component
	DemoCalc                      ssc.Component
	DescriptionDemoAutocomplete   ssc.Component
	DemoAutocomplete              ssc.Component
	DescriptionDemoEmailValidator ssc.Component
	DemoEmailValidator            ssc.Component
	DescriptionDemoRedirect       ssc.Component
	DemoRedirect                  ssc.Component
	DescriptionDemoNesting        ssc.Component
	DemoNesting                   ssc.Component
}

func (p *PageIndex) Template() *template.Template {
	return template.Must(template.New("page.index.html").Funcs(ssc.Funcs()).ParseGlob("*.html"))
}

func (p *PageIndex) Init() {
	p.Hero = ssc.RegC(p, &ComponentHero{
		Title:    "SSC Engine Demo",
		Subtitle: "Demo project with demonstration of SSC features",
	})
	p.DescriptionDemoUUID = ssc.RegC(p, &ComponentContent{
		Title:       "Async method",
		Description: "With async method you're able to fetch all needed data concurrently without worring about goroutines. All needed async methods are triggered on page render as separate goroutines.",
		Links: []ComponentContentLink{
			{
				Title: "component.demo.uuid.go",
				Href:  "https://github.com/yuriizinets/ssceng/blob/master/.demo/component.demo.uuid.go",
			},
			{
				Title: "component.demo.uuid.html",
				Href:  "https://github.com/yuriizinets/ssceng/blob/master/.demo/component.demo.uuid.html",
			},
		},
	})
	p.DemoUUID = ssc.RegC(p, &ComponentDemoUUID{})
	p.DescriptionDemoCounter = ssc.RegC(p, &ComponentContent{
		Title:       "Server Side Actions (SSA)",
		Description: "Component methods, executed and rendered entirely on server side. Frontend only gets ready for use HTML response.",
		Links: []ComponentContentLink{
			{
				Title: "component.demo.counter.go",
				Href:  "https://github.com/yuriizinets/ssceng/blob/master/.demo/component.demo.counter.go",
			},
			{
				Title: "component.demo.counter.html",
				Href:  "https://github.com/yuriizinets/ssceng/blob/master/.demo/component.demo.counter.html",
			},
		},
	})
	p.DemoCounter = ssc.RegC(p, &ComponentDemoCounter{})
	p.DescriptionDemoCalc = ssc.RegC(p, &ComponentContent{
		Title:       "State binding",
		Description: "Not all actions can be done on server side. Some things needs to be done on client side, like state binding. SSC library provides some primitives to make life easier.",
		Links: []ComponentContentLink{
			{
				Title: "component.demo.calc.go",
				Href:  "https://github.com/yuriizinets/ssceng/blob/master/.demo/component.demo.calc.go",
			},
			{
				Title: "component.demo.calc.html",
				Href:  "https://github.com/yuriizinets/ssceng/blob/master/.demo/component.demo.calc.html",
			},
		},
	})
	p.DemoCalc = ssc.RegC(p, &ComponentDemoCalc{})
	p.DescriptionDemoAutocomplete = ssc.RegC(p, &ComponentContent{
		Title:       "Combining events",
		Description: "Example, that combines state binding and server action.",
		Links: []ComponentContentLink{
			{
				Title: "component.demo.autocomplete.go",
				Href:  "https://github.com/yuriizinets/ssceng/blob/master/.demo/component.demo.autocomplete.go",
			},
			{
				Title: "component.demo.autocomplete.html",
				Href:  "https://github.com/yuriizinets/ssceng/blob/master/.demo/component.demo.autocomplete.html",
			},
		},
	})
	p.DemoAutocomplete = ssc.RegC(p, &ComponentDemoAutocomplete{
		Placeholder: "Select browser ...",
		Items: []string{
			"Edge",
			"Firefox",
			"Chrome",
			"Safari",
			"Opera",
		},
	})
	p.DescriptionDemoEmailValidator = ssc.RegC(p, &ComponentContent{
		Title:       "SSA form rocessing",
		Description: "Example of using formsubmit shortcut to simplify server-side form processing.",
		Links: []ComponentContentLink{
			{
				Title: "component.demo.email.validator.go",
				Href:  "https://github.com/yuriizinets/ssceng/blob/master/.demo/component.demo.email.validator.go",
			},
			{
				Title: "component.demo.email.validator.html",
				Href:  "https://github.com/yuriizinets/ssceng/blob/master/.demo/component.demo.email.validator.html",
			},
		},
	})
	p.DemoEmailValidator = ssc.RegC(p, &ComponentDemoEmailValidator{})
	p.DescriptionDemoRedirect = ssc.RegC(p, &ComponentContent{
		Title:       "SSA redirect",
		Description: "SSA has own redirect method, because usually SSA doesn't have any impact on frontend behavior rather than HTML structure.",
		Links: []ComponentContentLink{
			{
				Title: "component.demo.redirect.go",
				Href:  "https://github.com/yuriizinets/ssceng/blob/master/.demo/component.demo.redirect.go",
			},
			{
				Title: "component.demo.redirect.html",
				Href:  "https://github.com/yuriizinets/ssceng/blob/master/.demo/component.demo.redirect.html",
			},
		},
	})
	p.DemoRedirect = ssc.RegC(p, &ComponentDemoRedirect{})
	p.DescriptionDemoNesting = ssc.RegC(p, &ComponentContent{
		Title:       "Component nesting",
		Description: "Small example of component nesting. Registering component inside of another component.",
		Links: []ComponentContentLink{
			{
				Title: "component.demo.nesting.first.go",
				Href:  "https://github.com/yuriizinets/ssceng/blob/master/.demo/component.demo.nesting.first.go",
			},
			{
				Title: "component.demo.nesting.first.html",
				Href:  "https://github.com/yuriizinets/ssceng/blob/master/.demo/component.demo.nesting.first.html",
			},
		},
	})
	p.DemoNesting = ssc.RegC(p, &ComponentDemoNestingFirst{})
}
