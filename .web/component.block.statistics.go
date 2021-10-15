package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yuriizinets/kyoto"
)

// Statistics component cache
var statisticsCache = ComponentBlockStatistics{}
var statisticsTimestamp = time.Time{}

// Repository statistics
type ComponentBlockStatistics struct {
	Title string
	Repo  string

	// Internal
	Stars        int
	Forks        int
	Contributors int
	Sponsors     int
}

func (c *ComponentBlockStatistics) Init(p kyoto.Page) {
	if c.Repo == "" {
		panic("ComponentBlockStatistics: Repo is required")
	}
}

func (c *ComponentBlockStatistics) Async() error {
	// Update cache (once in 3 hours)
	if time.Since(statisticsTimestamp) > (3 * time.Hour) {
		resp, err := http.Get("https://api.github.com/repos/" + c.Repo)
		if err != nil {
			return err
		}
		if resp.StatusCode != 200 {
			return fmt.Errorf("%v", resp.StatusCode)
		}
		data := map[string]interface{}{}
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return err
		}
		statisticsCache.Stars = int(data["stargazers_count"].(float64))
		statisticsCache.Forks = int(data["forks_count"].(float64))
		statisticsCache.Contributors = 3
		statisticsCache.Sponsors = 1
		statisticsTimestamp = time.Now()
	}
	// Set fields from cache
	c.Stars = statisticsCache.Stars
	c.Forks = statisticsCache.Forks
	c.Contributors = statisticsCache.Contributors
	c.Sponsors = statisticsCache.Sponsors
	// Return
	return nil
}
