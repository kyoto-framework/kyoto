package main

import (
	"strconv"

	ssc "github.com/yuriizinets/ssceng"
)

type ComponentDemoCalc struct {
	FirstValue  string
	SecondValue string
	Result      string
}

func (c *ComponentDemoCalc) Init(p ssc.Page) {
	if c.FirstValue == "" {
		c.FirstValue = "5"
	}
	if c.SecondValue == "" {
		c.SecondValue = "5"
	}
	if c.Result == "" {
		c.Result = "10"
	}
}

func (c *ComponentDemoCalc) Actions() ssc.ActionMap {
	return ssc.ActionMap{
		"Calculate": func(args ...interface{}) {
			fv, err := strconv.Atoi(c.FirstValue)
			if err != nil {
				c.Result = "Can't calculate"
			}
			sv, err := strconv.Atoi(c.SecondValue)
			if err != nil {
				c.Result = "Can't calculate"
			}
			c.Result = strconv.Itoa(fv + sv)
		},
	}
}
