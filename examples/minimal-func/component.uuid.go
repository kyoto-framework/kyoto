package main

import (
	"encoding/json"
	"net/http"

	"github.com/kyoto-framework/kyoto"
	"github.com/kyoto-framework/kyoto/actions"
	"github.com/kyoto-framework/kyoto/lifecycle"
	"github.com/kyoto-framework/kyoto/state"
)

// Arguments wrapper
func ComponentUUID(title string) func(*kyoto.Core) {

	// Core receiver
	return func(c *kyoto.Core) {
		// Define state
		state.New(c, "Title", title)
		uuid := state.New(c, "UUID", "")

		// Define UUID loader
		var loader = func() error {
			// Execute request
			resp, err := http.Get("http://httpbin.org/uuid")
			if err != nil {
				uuid.Set("Error while retrieving UUID")
				return nil
			}
			// Defer closing of response body
			defer resp.Body.Close()
			// Decode response
			data := map[string]string{}
			json.NewDecoder(resp.Body).Decode(&data)
			// Set state
			uuid.Set(data["uuid"])
			// Return
			return nil
		}

		// Load UUID after initialization
		lifecycle.Async(c, loader)

		// Define reload action
		actions.Define(c, "Reload", func(args ...interface{}) {
			loader()
		})
	}
}
