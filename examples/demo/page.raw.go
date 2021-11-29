package main

import "io/ioutil"

type PageRaw struct{}

func (*PageRaw) Render() string {
	p, err := ioutil.ReadFile("page.raw.html")
	if err != nil {
		panic(err)
	}
	return string(p)
}
