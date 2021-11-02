package main

import (
	"time"

	"github.com/yuriizinets/kyoto"
)

type ComponentDemoLoading struct{}

func (c *ComponentDemoLoading) Actions() kyoto.ActionMap {
	return kyoto.ActionMap{
		"Start": func(args ...interface{}) {
			time.Sleep(3 * time.Second)
		},
	}
}
