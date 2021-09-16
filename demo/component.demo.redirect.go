package main

import ssc "github.com/yuriizinets/ssceng"

type ComponentDemoRedirect struct {
	Page ssc.Page `json:"-"`
}

func (c *ComponentDemoRedirect) Init(p ssc.Page) {
	c.Page = p
}

func (c *ComponentDemoRedirect) Actions() ssc.ActionMap {
	return ssc.ActionMap{
		"Redirect": func(args ...interface{}) {
			target := args[0].(string)
			ssc.Redirect(&ssc.RedirectParameters{
				Page:              c.Page,
				ResponseWriterKey: "internal:rw",
				RequestKey:        "internal:r",
				Target:            target,
			})
		},
	}
}
