package main

import (
	"encoding/json"
	"net/http"
)

type ComponentDemoUUID struct {
	UUID string
}

func (c *ComponentDemoUUID) Async() error {
	resp, err := http.Get("http://httpbin.org/uuid")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data := map[string]string{}
	json.NewDecoder(resp.Body).Decode(&data)
	c.UUID = data["uuid"]
	return nil
}
