package render

import (
	"html/template"
	"testing"

	"github.com/kyoto-framework/kyoto"
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

// TestFuncMap ensures that funcmap is not returning empty map
func TestFuncMap(t *testing.T) {
	if len(FuncMap(kyoto.NewCore())) == 0 {
		t.Error("FuncMap is empty")
	}
}
