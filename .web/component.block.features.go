package main

import "html/template"

type ComponentBlockFeatures struct {
	Entries []ComponentBlockFeaturesEntry
}

type ComponentBlockFeaturesEntry struct {
	Image       template.HTML
	Title       string
	Description string
}
