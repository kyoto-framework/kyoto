package main

import (
	"encoding/json"
	"net/http"

	"github.com/kyoto-framework/kyoto/smode"
)

type ComponentUUIDStruct struct {
	Title string
	UUID  string
}

func (c *ComponentUUIDStruct) Load() error {
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

func (c *ComponentUUIDStruct) Async() error {
	return c.Load()
}

func (c *ComponentUUIDStruct) Actions() smode.ActionMap {
	return smode.ActionMap{
		"Reload": func(args ...interface{}) {
			c.Load()
		},
	}
}
