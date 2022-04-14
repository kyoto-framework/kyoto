package smode

import (
	"html/template"
	"io"
	"testing"

	"github.com/kyoto-framework/kyoto"
)

type testPage struct{}

func (c *testPage) Template() *template.Template {
	return template.New("test")
}
func (c *testPage) Init()             {}
func (c *testPage) Async() error      { return nil }
func (c *testPage) AfterAsync() error { return nil }
func (c *testPage) Actions() ActionMap {
	return ActionMap{
		"Foo": func(args ...interface{}) {},
	}
}

type testComponent struct {
	Foo string
}

func (c *testComponent) Writer(w io.Writer) error { return nil }
func (c *testComponent) Init(p Page)              {}
func (c *testComponent) Async(p Page) error       { return nil }
func (c *testComponent) AfterAsync(p Page) error  { return nil }
func (c *testComponent) Actions(p Page) ActionMap {
	return ActionMap{
		"Bar": func(args ...interface{}) {},
	}
}

func TestAdaptPage(t *testing.T) {
	// Initialize core and apply page
	core := kyoto.NewCore()
	Adapt(&testPage{})(core)
	// Test name
	if name := core.State.Get("internal:name"); name == nil || name != "testPage" {
		t.Errorf("Expected name to be 'testPage', got '%s'", name)
	}
	// Test template
	if tbuilder := core.Context.Get("internal:render:tbuilder"); tbuilder == nil {
		t.Errorf("Expected template builder to be set")
	}
	// Test lifecycle
	flaginit := false
	flagasync := false
	flagafter := false
	for _, job := range core.Scheduler.Jobs {
		switch job.Group {
		case "init":
			flaginit = true
		case "async":
			flagasync = true
		case "afterasync":
			flagafter = true
		}
	}
	if !flaginit {
		t.Errorf("Expected init job to be set")
	}
	if !flagasync {
		t.Errorf("Expected async job to be set")
	}
	if !flagafter {
		t.Errorf("Expected afterasync job to be set")
	}
	// Cleanup
	core.Execute()
}

func TestAdaptComponent(t *testing.T) {
	// Initialize core and apply component
	core := kyoto.NewCore()
	Adapt(&testComponent{})(core)
	// Test name
	if name := core.State.Get("internal:name"); name == nil || name != "testComponent" {
		t.Errorf("Expected name to be 'testComponent', got '%s'", name)
	}
	// Test lifecycle
	flaginit := false
	flagasync := false
	flagafter := false
	for _, job := range core.Scheduler.Jobs {
		switch job.Group {
		case "init":
			flaginit = true
		case "async":
			flagasync = true
		case "afterasync":
			flagafter = true
		}
	}
	if !flaginit {
		t.Errorf("Expected init job to be set")
	}
	if !flagasync {
		t.Errorf("Expected async job to be set")
	}
	if !flagafter {
		t.Errorf("Expected afterasync job to be set")
	}
	// Cleanup
	core.Execute()
}

func TestRegC(t *testing.T) {
	// Initialize page and component
	page := &testPage{}
	component := &testComponent{}
	// Initialize root core
	rootcore := kyoto.NewCore()
	// Register page
	Adapt(page)(rootcore)
	// Register component
	UseC(page, component)
	// Get component core
	core, exists := cmap[component]
	if !exists {
		t.Error("Expected component core to be registered")
		return
	}
	// Check component is adapted
	if name := core.State.Get("internal:name"); name == nil || name != "testComponent" {
		t.Error("Empty core. Seems like component is not adapted")
	}
	// Check globals cleaned up
	core.Execute()
	if len(cmap) != 0 {
		t.Errorf("Expected component map to be empty, got %d", len(cmap))
	}
	if len(pmap) != 0 {
		t.Errorf("Expected page map to be empty, got %v", len(pmap))
	}
}

func TestRedirect(t *testing.T) {
	// Initialize core, page and apply it to core
	core := kyoto.NewCore()
	page := &testPage{}
	Adapt(page)(core)
	// Apply redirect
	Redirect(page, "/foo", 307)
	// Test redirect applied
	if core.Context.Get("internal:render:redirect") == nil {
		t.Error("Redirect is not set")
	}
	if core.Context.Get("internal:render:redirectCode") == nil {
		t.Error("Redirect code is not set")
	}
}

func TestFuncMap(t *testing.T) {
	// Initialize core, page and apply it to core
	core := kyoto.NewCore()
	page := &testPage{}
	Adapt(page)(core)
	// Generate funcmap with page
	funcmap := FuncMap(page)
	// Test funcmap
	if len(funcmap) == 0 {
		t.Error("Empty funcmap")
	}
	// Cleanup
	core.Execute()
}
