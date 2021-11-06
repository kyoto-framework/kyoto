package main

import (
	"time"

	"github.com/yuriizinets/kyoto"
)

type ComponentDemoFlush struct {
	Status string
}

func (c *ComponentDemoFlush) Actions(p kyoto.Page) kyoto.ActionMap {
	return kyoto.ActionMap{
		"Trigger": func(args ...interface{}) {
			c.Status = "Preparing ..."
			kyoto.SSAFlush(p, c)
			time.Sleep(2 * time.Second)
			c.Status = "Loading ..."
			kyoto.SSAFlush(p, c)
			time.Sleep(3 * time.Second)
			c.Status = "Finished!"
		},
	}
}
