package main

import ssc "github.com/yuriizinets/ssceng"

type ComponentDemoNestingFirst struct {
	Nested ssc.Component
}

func (c *ComponentDemoNestingFirst) Init(p ssc.Page) {
	c.Nested = ssc.RegC(p, &ComponentDemoNestingSecond{})
}
