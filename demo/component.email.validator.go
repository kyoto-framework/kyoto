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

func (c *ComponentEmailValidator) Actions() ssc.ActionMap {
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
