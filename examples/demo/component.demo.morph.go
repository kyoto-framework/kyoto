package main

import (
	"strconv"

	"github.com/kyoto-framework/kyoto"
)

type ComponentDemoMorph struct {
	Count    int
	ColorNum string
}

func (c *ComponentDemoMorph) Actions() kyoto.ActionMap {
	return kyoto.ActionMap{
		"Increment": func(args ...interface{}) {
			c.Count++
			c.ColorNum = strconv.Itoa(c.Count + 1)
			c.ColorNum = c.ColorNum[len(c.ColorNum)-1:]
		},
	}
}
