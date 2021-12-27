package main

import (
	"encoding/json"
	"net/http"

	"github.com/kyoto-framework/kyoto"
)

func ComponentUUID(title string) func(*kyoto.Builder) {
	return func(b *kyoto.Builder) {
		// Define UUID loader
		var loader = func() error {
			// Execute request
			resp, err := http.Get("http://httpbin.org/uuid")
			if err != nil {
				return err
			}
			// Defer closing of response body
			defer resp.Body.Close()
			// Decode response
			data := map[string]string{}
			json.NewDecoder(resp.Body).Decode(&data)
			// Set state
			b.State.Set("UUID", data["uuid"])
			// Return
			return nil
		}

		// Initialize empty state
		b.Init(func() {
			b.State.Set("Title", title)
			b.State.Set("UUID", "")
		})
		// Load UUID after initialization
		b.Async(loader)
		// Define reload action
		b.Action("Reload", func(args ...interface{}) {
			loader()
		})
	}
}
