package main

import (
	"io/ioutil"
	"net/http"

	"github.com/yuriizinets/go-ssc"
)

type ComponentHttpbinUUID struct {
	UUID string
}

func (*ComponentHttpbinUUID) Init() {

}

func (c *ComponentHttpbinUUID) Async() error {
	resp, err := http.Get("http://httpbin.org/uuid")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	c.UUID = string(data)
	return nil
}

func (*ComponentHttpbinUUID) AfterAsync() {

}

func (c *ComponentHttpbinUUID) Actions() ssc.ActionsMap {
	return ssc.ActionsMap{
		"Reload": func(args map[string]interface{}) {
			c.Async()
		},
	}
}
