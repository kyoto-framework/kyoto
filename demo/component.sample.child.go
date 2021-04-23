package main

import "github.com/yuriizinets/go-ssc"

type ComponentSampleChild struct {
	Value string
}

func (c *ComponentSampleChild) Init(p ssc.Page) {
	c.Value = "Child's component value"
}
