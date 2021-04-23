package main

import "github.com/yuriizinets/go-ssc"

type ComponentCounter struct {
	Count int
}

func (*ComponentCounter) Init(p ssc.Page) {

}

func (*ComponentCounter) Async() error {
	return nil
}

func (*ComponentCounter) AfterAsync() {

}

func (c *ComponentCounter) Actions() ssc.ActionsMap {
	return ssc.ActionsMap{
		"Increment": func(args ...interface{}) {
			c.Count++
		},
	}
}
