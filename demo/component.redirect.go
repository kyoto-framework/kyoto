package main

import "github.com/yuriizinets/go-ssc"

type ComponentRedirect struct {
	Page ssc.Page `json:"-"`
}

func (c *ComponentRedirect) Init(p ssc.Page) {
	c.Page = p
}

func (c *ComponentRedirect) Actions() ssc.ActionMap {
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
