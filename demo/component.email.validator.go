package main

import (
	"regexp"

	"github.com/yuriizinets/go-ssc"
)

var emailregex = regexp.MustCompile(`.+\@.+\..+`)

type ComponentEmailValidator struct {
	Email   string
	Message string
	Color   string
}

func (c *ComponentEmailValidator) Actions() ssc.ActionsMap {
	return ssc.ActionsMap{
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
