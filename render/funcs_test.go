package render

import (
	"fmt"
	"html/template"
	"io"
	"strings"
	"testing"

	"github.com/kyoto-framework/kyoto"
	"github.com/kyoto-framework/kyoto/helpers"
)

// TestRender ensures that render function correctly triggers a Writer render
func TestRenderWriter(t *testing.T) {
	// Initialize core
	core := kyoto.NewCore()
	// Apply Writer
	Writer(core, func(w io.Writer) error {
		_, err := w.Write([]byte("Foo"))
		return err
	})
	// Generate and check render output
	output := Render(core, core.State.Export())
	if output != "Foo" {
		t.Error("Render output is incorrect")
	}
}

// TestRenderFallback ensures that render function correctly
// fallbacks to template rendering in case of no Writer
func TestRenderFallback(t *testing.T) {
	// Initialize core
	core := kyoto.NewCore()
	// Inject component name (usually done by core.Component method)
	core.State.Set("internal:name", "testRenderBar")
	// Inject template and template builder
	Template(core, func() *template.Template {
		return template.Must(template.New("test").Parse(`
			{{ define "testRenderBar" }}Bar{{ end }}
		`))
	})
	// Generate and check render output
	output := Render(core, core.State.Export())
	if output != "Bar" {
		t.Error("Render output is incorrect")
	}
}

// TestDynamics ensures that dynamic integration is working as expected
func TestDynamics(t *testing.T) {
	// Generate dynamics output
	output := Dynamics("/custom/route/")
	// Check script tag
	if !strings.Contains(string(output), "<script>") {
		t.Error("Script tag is not set")
	}
	// Check ssapath
	if !strings.Contains(string(output), "ssapath") {
		t.Error("SSAPath is not set")
	}
	// Check default route
	if !strings.Contains(string(Dynamics()), "/internal/actions") {
		t.Error("Default route for dynamics is not set")
	}
}

// TestComponentAttrs ensures that component attributes are generating in the correct way
func TestComponentAttrs(t *testing.T) {
	// Initialize component state
	state := map[string]interface{}{
		"internal:name": "TestComponent",
		"Foo":           "Bar",
	}
	// Generate component attrs output
	output := ComponentAttrs(state)
	// Check component id
	if !strings.Contains(string(output), fmt.Sprintf(`cid='%s'`, helpers.ComponentID(state))) {
		t.Error("Component id is not set")
	}
	// Check component name
	if !strings.Contains(string(output), fmt.Sprintf(`name='%s'`, helpers.ComponentName(state))) {
		t.Error("Component name is not set")
	}
	// Check component state
	if !strings.Contains(string(output), fmt.Sprintf(`state='%s'`, helpers.ComponentSerialize(state))) {
		t.Error("Component state is not set")
	}
}

// TestActions ensures that action call is generating in the correct way
func TestAction(t *testing.T) {
	// Generate action output
	output := Action("Foo", "Bar")
	// Check action output
	if !strings.Contains(string(output), fmt.Sprintf(`Action(this, '%s', "%s")`, "Foo", "Bar")) {
		t.Error("Action output is incorrect")
	}
}

// TestBind ensures that binding call is generating in the correct way
func TestBind(t *testing.T) {
	// Generate bind output
	output := Bind("Foo")
	// Check bind output
	if !strings.Contains(string(output), fmt.Sprintf("Bind(this, '%s')", "Foo")) {
		t.Error("Bind output is incorrect")
	}
}

// TestFormSubmit ensures that form submit call is generating in the correct way
func TestFormSubmit(t *testing.T) {
	if FormSubmit() != "FormSubmit(this, event)" {
		t.Error("FormSubmit output is incorrect")
	}
}
