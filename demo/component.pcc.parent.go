package main

import "github.com/yuriizinets/go-ssc"

type ComponentPCCParent struct {
	Value string
	Child ssc.Component
}

func (c *ComponentPCCParent) Init(p ssc.Page) {
	c.Value = "None"
	c.Child = ssc.RegC(p, &ComponentPСCChild{})
}

func (c *ComponentPCCParent) Actions() ssc.ActionMap {
	return ssc.ActionMap{
		"SetValue": func(args ...interface{}) {
			c.Value = args[0].(string)
		},
	}
}
