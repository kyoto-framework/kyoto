package main

import "github.com/yuriizinets/kyoto"

type ComponentDemoIntersect struct {
	Count int
}

func (c *ComponentDemoIntersect) Actions() kyoto.ActionMap {
	return kyoto.ActionMap{
		"Increment": func(args ...interface{}) {
			c.Count++
		},
	}
}
