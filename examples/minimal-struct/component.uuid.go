package main

import (
	"encoding/json"
	"net/http"

	"github.com/kyoto-framework/kyoto/smode"
)

type ComponentUUID struct {
	Title string
	UUID  string
}

func (c *ComponentUUID) Async() error {
	return c.load()
}

func (c *ComponentUUID) Actions() smode.ActionMap {
	return smode.ActionMap{
		"Reload": func(args ...interface{}) {
			c.load()
		},
	}
}

func (c *ComponentUUID) load() error {
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
	c.UUID = data["uuid"]
	// Return
	return nil
}
