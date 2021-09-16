package main

import "html/template"

type ComponentFooter struct {
	Left      []ComponentFooterHref
	Center    []ComponentFooterHref
	Copyright string
}

type ComponentFooterHref struct {
	Title template.HTML
	Image template.HTML
	Href  string
}
