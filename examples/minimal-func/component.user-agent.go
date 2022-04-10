package main

import (
	"fmt"
	"io"

	"github.com/kyoto-framework/kyoto"
	"github.com/kyoto-framework/kyoto/lifecycle"
	"github.com/kyoto-framework/kyoto/render"
)

func ComponentUserAgent(c *kyoto.Core) {
	lifecycle.Async(c, func() error {
		request := c.Context.GetRequest()
		c.State.Set("UserAgent", request.UserAgent())
		return nil
	})
	render.Writer(c, func(w io.Writer) error {
		fmt.Fprintf(w, `<div>User-Agent: %s</div>`, c.State.Get("UserAgent"))
		return nil
	})
}
