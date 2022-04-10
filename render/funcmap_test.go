package render

import (
	"testing"

	"github.com/kyoto-framework/kyoto"
)

// TestFuncMap ensures that funcmap is not returning empty map
func TestFuncMap(t *testing.T) {
	if len(FuncMap(kyoto.NewCore())) == 0 {
		t.Error("FuncMap is empty")
	}
}
