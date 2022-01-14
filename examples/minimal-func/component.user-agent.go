package main

import (
	"github.com/kyoto-framework/kyoto"
	"github.com/kyoto-framework/kyoto/lifecycle"
)

func ComponentUserAgent(c *kyoto.Core) {
	lifecycle.Async(c, func() error {
		request := c.Context.GetRequest()
		c.State.Set("UserAgent", request.UserAgent())
		return nil
	})
}
