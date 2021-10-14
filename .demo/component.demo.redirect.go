package main

import "github.com/yuriizinets/kyoto"

type ComponentDemoRedirect struct {
	Page kyoto.Page `json:"-"`
}

func (c *ComponentDemoRedirect) Init(p kyoto.Page) {
	c.Page = p
}

func (c *ComponentDemoRedirect) Actions() kyoto.ActionMap {
	return kyoto.ActionMap{
		"Redirect": func(args ...interface{}) {
			target := args[0].(string)
			kyoto.Redirect(&kyoto.RedirectParameters{
				Page:              c.Page,
				ResponseWriterKey: "internal:rw",
				RequestKey:        "internal:r",
				Target:            target,
			})
		},
	}
}
