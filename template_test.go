package kyoto

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"os"
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

func TestTemplateDefault(t *testing.T) {
	// Create context
	c := &Context{}
	// Create a template file for test
	ioutil.WriteFile("template_test.html", []byte("Placeholder for tests"), 0644)
	// Create template
	Template(c, "template_test.html")
	// Check template is set
	if c.Template == nil {
		t.Error("Template is not set")
		os.Remove("template_test.html")
		return
	}
	// Check template is rendering
	buf := &bytes.Buffer{}
	if err := c.Template.Execute(buf, nil); err != nil || buf.String() != "Placeholder for tests" {
		t.Error("Something wrong with template rendering")
		os.Remove("template_test.html")
		return
	}
	// Cleanup
	os.Remove("template_test.html")
}

func TestTemplateGlob(t *testing.T) {
	// Create context
	c := &Context{}
	// Create a template file for test
	os.Mkdir("tmp", 0700)
	err := ioutil.WriteFile("tmp/template_test.html", []byte("Placeholder for tests"), 0644)
	if err != nil {
		panic(err)
	}
	// Set config
	TemplateConf.ParseGlob = "tmp/*.html"
	// Create template
	Template(c, "template_test.html")
	// Check template is set
	if c.Template == nil {
		t.Error("Template is not set")
		os.RemoveAll("tmp")
		return
	}
	// Check template is rendering
	buf := &bytes.Buffer{}
	if err := c.Template.Execute(buf, nil); err != nil || buf.String() != "Placeholder for tests" {
		t.Error("Something wrong with template rendering")
		os.RemoveAll("tmp")
		return
	}
	// Cleanup
	os.RemoveAll("tmp")
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
