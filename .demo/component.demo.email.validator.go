package main

import (
	"regexp"

	ssc "github.com/yuriizinets/ssceng"
)

var emailregex = regexp.MustCompile(`.+\@.+\..+`)

type ComponentDemoEmailValidator struct {
	Email   string
	Message string
	Color   string
}

func (c *ComponentDemoEmailValidator) Actions() ssc.ActionMap {
	return ssc.ActionMap{
		"Submit": func(args ...interface{}) {
			if emailregex.MatchString(c.Email) {
				c.Message = "Provided email is valid"
				c.Color = "green"
			} else {
				c.Message = "Provided email is not valid"
				c.Color = "red"
			}
		},
	}
}
