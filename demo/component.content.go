package main

type ComponentContent struct {
	Title       string
	Subtitle    string
	Description string
	Image       string
	Links       []ComponentContentLink
	LinksCenter bool
	SideSwitch  bool
}

type ComponentContentLink struct {
	Href  string
	Title string
}
