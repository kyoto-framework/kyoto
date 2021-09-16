package main

import "html/template"

type ComponentNavbar struct {
	Left  []ComponentNavbarHref
	Right []ComponentNavbarHref
}

type ComponentNavbarHref struct {
	Title template.HTML
	Image template.HTML
	Href  string
}
