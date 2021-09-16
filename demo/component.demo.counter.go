package main

import ssc "github.com/yuriizinets/ssceng"

type ComponentDemoCounter struct {
	Count int
}

func (c *ComponentDemoCounter) Actions() ssc.ActionMap {
	return ssc.ActionMap{
		"Increment": func(args ...interface{}) {
			c.Count++
		},
	}
}
