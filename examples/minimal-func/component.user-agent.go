package main

import (
	"fmt"
	"io"

	"github.com/kyoto-framework/kyoto"
	"github.com/kyoto-framework/kyoto/lifecycle"
	"github.com/kyoto-framework/kyoto/render"
	"github.com/kyoto-framework/kyoto/state"
)

func ComponentUserAgent(c *kyoto.Core) {
	// Define state
	useragent := state.New(c, "UserAgent", "")

	// Define lifecycle
	lifecycle.Async(c, func() error {
		request := c.Context.GetRequest()
		useragent.Set(request.UserAgent())
		return nil
	})

	// Define rendering
	render.Writer(c, func(w io.Writer) error {
		fmt.Fprintf(w, `<div>User-Agent: %s</div>`, c.State.Get("UserAgent"))
		return nil
	})
}
