package main

import (
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/kyoto-framework/kyoto"
)

type ComponentDemoCounterRaw struct {
	Count int
}

func (c *ComponentDemoCounterRaw) Actions() kyoto.ActionMap {
	return kyoto.ActionMap{
		"Increment": func(args ...interface{}) {
			c.Count++
		},
	}
}

func (c *ComponentDemoCounterRaw) Render() string {
	// Read component markup
	cbts, err := ioutil.ReadFile("component.demo.counter.raw.html")
	if err != nil {
		panic(err)
	}
	// Convert to string
	dom := string(cbts)
	// Replace dynamics
	dom = strings.ReplaceAll(dom, "{attrs}", string(kyoto.TComponentAttrs(c)))
	dom = strings.ReplaceAll(dom, "{count}", strconv.Itoa(c.Count))
	dom = strings.ReplaceAll(dom, "{increment}", string(kyoto.TAction("Increment")))
	// Return
	return string(dom)
}
