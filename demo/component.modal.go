package main

import "github.com/yuriizinets/go-ssc"

type ComponentModal struct {
	Show bool
}

func (*ComponentModal) Init(p ssc.Page) {}

func (c *ComponentModal) Async() error { return nil }

func (*ComponentModal) AfterAsync() {}

func (c *ComponentModal) Actions() ssc.ActionsMap {
	return ssc.ActionsMap{
		"Open": func(args map[string]interface{}) {
			c.Show = true
		},
		"Close": func(args map[string]interface{}) {
			c.Show = false
		},
	}
}
