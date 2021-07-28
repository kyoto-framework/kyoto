package main

import (
	"encoding/json"
	"net/http"

	"github.com/yuriizinets/go-ssc"
)

type ComponentHttpbinUUID struct {
	UUID string
}

func (c *ComponentHttpbinUUID) Async() error {
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

func (c *ComponentHttpbinUUID) Actions() ssc.ActionMap {
	return ssc.ActionMap{
		"Reload": func(args ...interface{}) {
			c.Async()
		},
	}
}
