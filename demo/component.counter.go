package main

import "github.com/yuriizinets/go-ssc"

type ComponentCounter struct {
	Count int
}

func (*ComponentCounter) Init() {

}

func (*ComponentCounter) Async() error {
	return nil
}

func (*ComponentCounter) AfterAsync() {

}

func (c *ComponentCounter) Actions() ssc.ActionsMap {
	return ssc.ActionsMap{
		"Increment": func(args map[string]interface{}) {
			c.Count++
		},
	}
}
