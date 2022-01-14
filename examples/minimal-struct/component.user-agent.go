package main

import (
	"net/http"

	"github.com/kyoto-framework/kyoto/smode"
)

type ComponentUserAgent struct {
	UserAgent string
}

func (c *ComponentUserAgent) Init(p smode.Page) {
	request := smode.GetContext(p, "internal:r").(*http.Request)
	c.UserAgent = request.UserAgent()
}
