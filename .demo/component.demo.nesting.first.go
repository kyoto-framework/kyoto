package main

import "github.com/yuriizinets/kyoto"

type ComponentDemoNestingFirst struct {
	Nested kyoto.Component
}

func (c *ComponentDemoNestingFirst) Init(p kyoto.Page) {
	c.Nested = kyoto.RegC(p, &ComponentDemoNestingSecond{})
}
