package main

import "github.com/kyoto-framework/kyoto"

type ComponentDemoCounter struct {
	Count int
}

func (c *ComponentDemoCounter) Actions() kyoto.ActionMap {
	return kyoto.ActionMap{
		"Increment": func(args ...interface{}) {
			c.Count++
		},
	}
}
