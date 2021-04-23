package main

import "github.com/yuriizinets/go-ssc"

type ComponentModal struct {
	Show bool
}

func (c *ComponentModal) Actions() ssc.ActionsMap {
	return ssc.ActionsMap{
		"Open": func(args ...interface{}) {
			c.Show = true
		},
		"Close": func(args ...interface{}) {
			c.Show = false
		},
	}
}
