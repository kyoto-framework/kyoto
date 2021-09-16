package main

type ComponentContent struct {
	Title       string
	Subtitle    string
	Description string
	Image       string
	Links       []ComponentContentLink
	CenterLinks bool
	SideSwitch  bool
}

type ComponentContentLink struct {
	Href  string
	Title string
}
