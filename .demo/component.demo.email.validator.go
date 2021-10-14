package main

import (
	"regexp"

	"github.com/yuriizinets/kyoto"
)

var emailregex = regexp.MustCompile(`.+\@.+\..+`)

type ComponentDemoEmailValidator struct {
	Email   string
	Message string
	Color   string
}

func (c *ComponentDemoEmailValidator) Actions() kyoto.ActionMap {
	return kyoto.ActionMap{
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
