package main

import "github.com/yuriizinets/go-ssc"

type ComponentCounter struct {
	Count int
}

func (c *ComponentCounter) Actions() ssc.ActionsMap {
	return ssc.ActionsMap{
		"Increment": func(args ...interface{}) {
			c.Count++
		},
	}
}
