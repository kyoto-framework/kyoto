package main

type ComponentBlockStatistics struct {
	Entries []ComponentBlockStatisticsEntry
}

type ComponentBlockStatisticsEntry struct {
	Title string
	Count int
}
