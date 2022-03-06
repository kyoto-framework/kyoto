package kyoto

import (
	"fmt"
	"reflect"
	"sync"
	"time"
)

// Deprecated: Insights storage
var (
	insights   = []*Insights{}
	insightsrw = sync.RWMutex{}
)

// Deprecated: Insights data type
type Insights struct {
	InsightsTiming

	ID     string      `json:"id"`
	Name   string      `json:"name"`
	Nested []*Insights `json:"nested"`
}

// Deprecated: InsightsTiming data type
type InsightsTiming struct {
	Init       time.Duration `json:"i"`
	Async      time.Duration `json:"a"`
	AfterAsync time.Duration `json:"aa"`
	Render     time.Duration `json:"r"`
}

// Deprecated: NewInsights creates new insights instance for provided object pointer,
// saves insights pointer to local store and returns it
// The oldest insights are cut in case of store overflow (INSIGHTS_LIMIT config)
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

// Deprecated: GetInsights returns insights pointer by given object pointer
func GetInsights(p interface{}) *Insights {
	for _, i := range insights {
		if i.ID == InsightsID(p) {
			return i
		}
		for _, ci := range i.Nested {
			if ci.ID == InsightsID(p) {
				return ci
			}
		}
	}
	return nil
}

// Deprecated: GetInsightsByID returns insights pointer by insights id
func GetInsightsByID(id string) *Insights {
	for _, i := range insights {
		if i.ID == id {
			return i
		}
	}
	return nil
}

// Deprecated: Update the Insights value
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

// Deprecated: GetOrCreateNested attempts to return existing nested insights, or returns new ones
func (i *Insights) GetOrCreateNested(p interface{}) *Insights {
	// Try to return existing nested insights
	for _, ci := range i.Nested {
		if ci.ID == InsightsID(p) {
			return ci
		}
	}
	// Init new nested insights
	ci := &Insights{
		ID:   InsightsID(p),
		Name: InsightsName(p),
	}
	i.Nested = append(i.Nested, ci)
	// Return new nested insights
	return ci
}

// Deprecated: InsightsID is a function to generate ID from pointer
func InsightsID(p interface{}) string {
	return fmt.Sprintf("%p", p)
}

// Deprecated: InsightsName is a function that extracts type name from pointer
func InsightsName(p interface{}) string {
	return reflect.ValueOf(p).Elem().Type().Name()
}
