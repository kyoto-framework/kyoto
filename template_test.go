package kyoto

import (
	"html/template"
	"testing"
)

// TestComposeFuncMap ensures that funcmaps are composed correctly
func TestComposeFuncMap(t *testing.T) {
	// Define test funcmaps
	fm1 := template.FuncMap{
		"Foo": func() int { return 1 },
	}
	fm2 := template.FuncMap{
		"Bar": func() int { return 2 },
	}
	// Compose
	fm3 := ComposeFuncMap(fm1, fm2)
	// Check len
	if len(fm3) != 2 {
		t.Error("FuncMap is not composed correctly")
	}
	// Check functions are correct
	if fm3["Foo"].(func() int)() != 1 {
		t.Error("FuncMap is not composed correctly")
	}
	if fm3["Bar"].(func() int)() != 2 {
		t.Error("FuncMap is not composed correctly")
	}
}

func TestTemplate(t *testing.T) {
	// Create context
	c := Context{}
	// Create template
	Template(&c, "test")
	// Check template is set
	if c.Template == nil {
		t.Error("Template is not set")
	}
}

func TestTemplateInline(t *testing.T) {
	// Create context
	c := Context{}
	// Create template
	TemplateInline(&c, "<html><head><title>Test</title></head><body><h1>Test</h1><p>This is a test.</p></body></html>")
	// Check template is set
	if c.Template == nil {
		t.Error("Template is not set")
	}
}

func TestTemplateRaw(t *testing.T) {
	// Create context
	c := Context{}
	// Create template
	tmpl, _ := template.New("test").Parse("<html><head><title>Test</title></head><body><h1>Test</h1><p>This is a test.</p></body></html>")
	TemplateRaw(&c, tmpl)
	// Check template is set
	if c.Template == nil {
		t.Error("Template is not set")
	}
}
