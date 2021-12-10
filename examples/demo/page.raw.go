package main

import (
	"io/ioutil"
	"strings"

	"github.com/kyoto-framework/kyoto"
)

type PageRaw struct {
	DemoCounter kyoto.Component
}

func (p *PageRaw) Init() {
	p.DemoCounter = kyoto.RegC(p, &ComponentDemoCounterRaw{})
}

func (p *PageRaw) Render() string {
	// Read page markup
	cbts, err := ioutil.ReadFile("page.raw.html")
	if err != nil {
		panic(err)
	}
	// Convert to string
	dom := string(cbts)
	// Replace dynamics
	dom = strings.ReplaceAll(dom, "{component-raw}", kyoto.TRender(p.DemoCounter))
	dom = strings.ReplaceAll(dom, "{dynamics}", string(kyoto.TDynamics("/demo/ssa")))
	// Return
	return dom
}
