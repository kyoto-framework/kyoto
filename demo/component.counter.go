package main

import "github.com/yuriizinets/go-ssc"

type ComponentCounter struct {
	Count int
}

func (c *ComponentCounter) Actions() ssc.ActionMap {
	return ssc.ActionMap{
		"Increment": func(args ...interface{}) {
			c.Count++
		},
	}
}
