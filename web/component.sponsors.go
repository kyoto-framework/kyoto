package main

type ComponentSponsors struct {
	Title   string
	Entries []ComponentSponsorsEntry
}

type ComponentSponsorsEntry struct {
	Photo string
	Href  string
}
