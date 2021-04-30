package main

import "github.com/yuriizinets/go-ssc"

type ComponentPСCChild struct {
	Value string
}

func (c *ComponentPСCChild) Init(p ssc.Page) {
	c.Value = "Child's component value"
}
