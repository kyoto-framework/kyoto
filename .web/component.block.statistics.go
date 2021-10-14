package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yuriizinets/kyoto"
)

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
	resp, err := http.Get("https://api.github.com/repos/" + c.Repo)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("%v", resp.StatusCode)
	}
	repo := map[string]interface{}{}
	err = json.NewDecoder(resp.Body).Decode(&repo)
	if err != nil {
		return err
	}
	c.Stars = int(repo["stargazers_count"].(float64))
	c.Forks = int(repo["forks_count"].(float64))
	c.Contributors = 3
	c.Sponsors = 1
	return nil
}
