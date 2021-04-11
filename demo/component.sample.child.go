package main

import "github.com/yuriizinets/go-ssc"

type ComponentSampleChild struct {
	Value string
}

func (*ComponentSampleChild) Init(p ssc.Page) {}

func (*ComponentSampleChild) Async() error { return nil }

func (*ComponentSampleChild) AfterAsync() {}

func (c *ComponentSampleChild) Actions() ssc.ActionsMap {
	return ssc.ActionsMap{}
}
