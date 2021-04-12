package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var cpgotemplate = `package main

import (
	"html/template"
	"github.com/yuriizinets/go-ssc"
)

type {name} struct{}

func (*{name}) Template() *template.Template {
	return template.Must(template.New("{htmlfile}").Funcs(funcmap()).ParseGlob("*.html"))
}

func (p *{name}) Init() {}

func (*{name}) Meta() ssc.Meta {
	return ssc.Meta{}
}
`

var cphtmltemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    {{ meta . }}
    {{ dynamics }}
</head>
<body>
    
</body>
</html>

`

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
	fsnewpage := flag.NewFlagSet("new:page", flag.ExitOnError)
	fsnpname := fsnewpage.String("name", "", "Page name (required)")
	fsnpgofile := fsnewpage.String("gofile", "", "Go file path (required)")
	fsnphtmlfile := fsnewpage.String("htmlfile", "", "HTML file path (required)")

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
	case "new:page":
		fsnewpage.Parse(os.Args[2:])
	default:
		fallback()
	}

	if fsnewpage.Parsed() {
		if *fsnpname == "" || *fsnpgofile == "" || *fsnphtmlfile == "" {
			fmt.Println("")
			fsnewpage.PrintDefaults()
			fmt.Println("")
			os.Exit(0)
		}
		newpage(
			*fsnpname,
			*fsnpgofile,
			*fsnphtmlfile,
		)
	} else if fsnewcomponent.Parsed() {
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
		"- new:page",
		"- new:component",
	)
	fmt.Println("")
	os.Exit(0)
}

func newpage(name, gofile, htmlfile string) {
	// Compile and save Go template
	gotemplate := strings.ReplaceAll(cpgotemplate, "{name}", name)
	ioutil.WriteFile(gofile, []byte(gotemplate), 0644)
	// Compile and save HTML tempalte
	htmltemplate := strings.ReplaceAll(cphtmltemplate, "{name}", name)
	htmltemplate = strings.ReplaceAll(htmltemplate, "{htmlfile}", htmlfile)
	ioutil.WriteFile(htmlfile, []byte(htmltemplate), 0644)
}

func newcomponent(name, gofile, htmlfile string) {
	// Compile and save Go template
	gotemplate := strings.ReplaceAll(cmgotemplate, "{name}", name)
	ioutil.WriteFile(gofile, []byte(gotemplate), 0644)
	// Compile and save HTML tempalte
	htmltemplate := strings.ReplaceAll(cmhtmltemplate, "{name}", name)
	ioutil.WriteFile(htmlfile, []byte(htmltemplate), 0644)
}
