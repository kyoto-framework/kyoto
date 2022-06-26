package kyoto

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"os"
	"testing"
)

// TestComposeFuncMap ensures that funcmaps are composed correctly.
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

// TestTemplateDefault tests Template function with a default TemplateConf.
func TestTemplateDefault(t *testing.T) {
	// Define a setup for a test
	setup := func() {
		// Create a template file for test
		ioutil.WriteFile("template_test.html", []byte("Placeholder for tests"), 0644)
	}
	// Define a cleanup for a test
	cleanup := func() {
		// Remove template file
		os.Remove("template_test.html")
	}
	// Setup
	setup()
	defer cleanup()
	// Create context
	c := &Context{}
	// Create template
	Template(c, "template_test.html")
	// Check template is set
	if c.Template == nil {
		t.Error("Template is not set")
		return
	}
	// Check template is rendering
	buf := &bytes.Buffer{}
	if err := c.Template.Execute(buf, nil); err != nil || buf.String() != "Placeholder for tests" {
		t.Error("Something wrong with template rendering")
		return
	}
}

// TestTemplateGlob tests Template function with a TemplateConf.ParseGlob setting.
func TestTemplateGlob(t *testing.T) {
	// Config store
	glob := ""
	// Define a setup for a test
	setup := func() {
		// Create a template file inside of temp directory for test
		os.Mkdir("tmp", 0700)
		err := ioutil.WriteFile("tmp/template_test.html", []byte("Placeholder for tests"), 0644)
		if err != nil {
			panic(err)
		}
		// Set configuration and save default
		glob = TemplateConf.ParseGlob
		TemplateConf.ParseGlob = "tmp/*.html"
	}
	// Define a cleanup for a test
	cleanup := func() {
		// Cleanup template
		os.RemoveAll("tmp")
		// Bring back configuration
		TemplateConf.ParseGlob = glob
	}
	// Setup
	setup()
	defer cleanup()
	// Create context
	c := &Context{}
	// Create template
	Template(c, "template_test.html")
	// Check template is set
	if c.Template == nil {
		t.Error("Template is not set")
		return
	}
	// Check template is rendering
	buf := &bytes.Buffer{}
	if err := c.Template.Execute(buf, nil); err != nil || buf.String() != "Placeholder for tests" {
		t.Error("Something wrong with template rendering")
		return
	}
}

// TestTemplateGlob tests TemplateInline function with a default TemplateConf.
func TestTemplateInline(t *testing.T) {
	setup := func() {
		// Create a template file for test
		ioutil.WriteFile("template_test.html", []byte("Placeholder for tests"), 0644)
	}
	cleanup := func() {
		// Remove template file
		os.Remove("template_test.html")
	}
	// Setup
	setup()
	defer cleanup()
	// Create context
	c := Context{}
	// Create template
	TemplateInline(&c, "<html>This is a test</html>")
	// Check template is set
	if c.Template == nil {
		t.Error("Template is not set")
	}
	// Check template is rendering
	buf := &bytes.Buffer{}
	if err := c.Template.Execute(buf, nil); err != nil || buf.String() != "<html>This is a test</html>" {
		t.Error("Something wrong with template rendering")
		return
	}
}

// TestTemplateGlob tests TemplateInline function with a default TemplateConf.
func TestTemplateRaw(t *testing.T) {
	// Create context
	c := Context{}
	// Create template
	tmpl, _ := template.New("test").Parse("<html>This is a test</html>")
	TemplateRaw(&c, tmpl)
	// Check template is set
	if c.Template == nil {
		t.Error("Template is not set")
		return
	}
	// Check template is rendering
	buf := &bytes.Buffer{}
	if err := c.Template.Execute(buf, nil); err != nil || buf.String() != "<html>This is a test</html>" {
		t.Error("Something wrong with template rendering")
		return
	}
}
