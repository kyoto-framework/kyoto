package main

import "github.com/kyoto-framework/kyoto"

type ComponentDemoPoll struct {
	Count int
}

func (c *ComponentDemoPoll) Actions() kyoto.ActionMap {
	return kyoto.ActionMap{
		"Increment": func(args ...interface{}) {
			c.Count++
		},
	}
}
