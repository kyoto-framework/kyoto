package kyoto

import (
	"fmt"
	"reflect"
	"sync"
	"time"
)

// Insights storage
var insights = []*Insights{}
var insightsrw = sync.RWMutex{}

type Insights struct {
	InsightsTiming

	ID     string
	Name   string
	Nested []*Insights
}

type InsightsTiming struct {
	Init       time.Duration
	Async      time.Duration
	AfterAsync time.Duration
	Render     time.Duration
}

// NewInsights creates new insights instance for provided object pointer,
//  saves insights pointer to local store and returns it.
// Oldest insights are cutted in case of store overflow (INSIGHTS_LIMIT config)
func NewInsights(p interface{}) *Insights {
	// Init new insights
	i := &Insights{
		ID:   InsightsID(p),
		Name: InsightsName(p),
	}
	// Add new insights to storage
	insightsrw.Lock()
	insights = append(insights, i)
	insightsrw.Unlock()
	// Cut insights storage in case of overflow
	if len(insights) > INSIGHTS_LIMIT {
		insights = insights[1:]
	}
	return i
}

// GetInsights returns insights pointer by given object pointer
func GetInsights(p interface{}) *Insights {
	for _, i := range insights {
		if i.ID == InsightsID(p) {
			return i
		}
	}
	return nil
}

// GetInsightsByID returns insights pointer by insights id
func GetInsightsByID(id string) *Insights {
	for _, i := range insights {
		if i.ID == id {
			return i
		}
	}
	return nil
}

func (i *Insights) Update(t InsightsTiming) {
	if t.Init != 0 {
		i.InsightsTiming.Init = t.Init
	}
	if t.Async != 0 {
		i.InsightsTiming.Async = t.Async
	}
	if t.AfterAsync != 0 {
		i.InsightsTiming.AfterAsync = t.AfterAsync
	}
	if t.Render != 0 {
		i.InsightsTiming.Render = t.Render
	}
}

func (i *Insights) Component(c Component) *Insights {
	// Try to return existing component insights
	for _, ci := range i.Nested {
		if ci.ID == InsightsID(c) {
			return ci
		}
	}
	// Init new component insights
	ci := &Insights{
		ID:   InsightsID(c),
		Name: InsightsName(c),
	}
	i.Nested = append(i.Nested, ci)
	// Return new component insights
	return ci
}

// InsightsID is a function to generate ID from pointer
func InsightsID(p interface{}) string {
	return fmt.Sprintf("%p", p)
}

// InsightsName is a function that extracts type name from pointer
func InsightsName(p interface{}) string {
	return reflect.ValueOf(p).Elem().Type().Name()
}
