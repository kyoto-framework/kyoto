package main

import "html/template"

type ComponentBlockFAQ struct {
	Title   string
	Entries []ComponentBlockFAQEntry
}

type ComponentBlockFAQEntry struct {
	Question string
	Answer   template.HTML
}
