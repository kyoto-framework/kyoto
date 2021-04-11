package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var cmgotemplate = `package main

import "github.com/yuriizinets/go-ssc"

type {name} struct{}

func (*{name}) Init(p ssc.Page) {}

func (*{name}) Async() error { return nil }

func (*{name}) AfterAsync() {}

func (c *{name}) Actions() ssc.ActionsMap {
	return ssc.ActionsMap{}
}
`

var cmhtmltemplate = `{{ define "{name}" }}
<div {{ componentattrs . }}>
</div>
{{ end }}
`

func main() {
	fsnewcomponent := flag.NewFlagSet("new:component", flag.ExitOnError)
	fsncname := fsnewcomponent.String("name", "", "Component name (required)")
	fsncgofile := fsnewcomponent.String("gofile", "", "Go file path (required)")
	fsnchtmlfile := fsnewcomponent.String("htmlfile", "", "HTML file path (required)")

	if len(os.Args) < 2 {
		fallback()
	}

	switch os.Args[1] {
	case "new:component":
		fsnewcomponent.Parse(os.Args[2:])
	default:
		fallback()
	}

	if fsnewcomponent.Parsed() {
		if *fsncname == "" || *fsncgofile == "" || *fsnchtmlfile == "" {
			fmt.Println("")
			fsnewcomponent.PrintDefaults()
			fmt.Println("")
			os.Exit(0)
		}
		newcomponent(
			*fsncname,
			*fsncgofile,
			*fsnchtmlfile,
		)
	}
}

func fallback() {
	fmt.Println("")
	fmt.Println(
		"Command is required. Supported commands:\n",
		"- new:component",
	)
	fmt.Println("")
	os.Exit(0)
}

func newcomponent(name, gofile, htmlfile string) {
	// Compile and save Go template
	gotemplate := strings.ReplaceAll(cmgotemplate, "{name}", name)
	ioutil.WriteFile(gofile, []byte(gotemplate), 0644)
	// Compile and save HTML tempalte
	htmltemplate := strings.ReplaceAll(cmhtmltemplate, "{name}", name)
	ioutil.WriteFile(htmlfile, []byte(htmltemplate), 0644)
}
