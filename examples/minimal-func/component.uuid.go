package main

import (
	"encoding/json"
	"net/http"

	"github.com/kyoto-framework/kyoto"
	"github.com/kyoto-framework/kyoto/actions"
	"github.com/kyoto-framework/kyoto/lifecycle"
)

func ComponentUUID(title string) func(*kyoto.Core) {
	return func(c *kyoto.Core) {
		// Define UUID loader
		var loader = func() error {
			// Execute request
			resp, err := http.Get("http://httpbin.org/uuid")
			if err != nil {
				c.State.Set("UUID", "Error while retrieving UUID")
				return nil
			}
			// Defer closing of response body
			defer resp.Body.Close()
			// Decode response
			data := map[string]string{}
			json.NewDecoder(resp.Body).Decode(&data)
			// Set state
			c.State.Set("UUID", data["uuid"])
			// Return
			return nil
		}

		// Initialize empty state
		lifecycle.Init(c, func() {
			c.State.Set("Title", title)
			c.State.Set("UUID", "")
		})
		// Load UUID after initialization
		lifecycle.Async(c, loader)
		// Define reload action
		actions.Define(c, "Reload", func(args ...interface{}) {
			loader()
		})
	}
}
