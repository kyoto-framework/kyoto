package main

import "github.com/yuriizinets/go-ssc"

type ComponentSampleParent struct {
	Value                string
	ComponentSampleChild ssc.Component
}

func (c *ComponentSampleParent) Init(p ssc.Page) {
	c.Value = "None"
	c.ComponentSampleChild = ssc.RegC(p, &ComponentSampleChild{})
}

func (*ComponentSampleParent) Async() error { return nil }

func (*ComponentSampleParent) AfterAsync() {}

func (c *ComponentSampleParent) Actions() ssc.ActionsMap {
	return ssc.ActionsMap{
		"SetValue": func(args map[string]interface{}) {
			c.Value = args["Value"].(string)
		},
	}
}
