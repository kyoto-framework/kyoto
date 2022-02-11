package render

import "testing"

// TestFuncMap ensures that funcmap is not returning empty map
func TestFuncMap(t *testing.T) {
	if len(FuncMap()) == 0 {
		t.Error("FuncMap is empty")
	}
}
