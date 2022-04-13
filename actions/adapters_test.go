package actions

import (
	"testing"

	"github.com/kyoto-framework/kyoto"
)

// TestDefine ensures that action definitions are working as expected
func TestDefine(t *testing.T) {
	// Initialize core
	core := kyoto.NewCore()
	// Define test action
	Define(core, "test", func(args ...interface{}) {})
	// Test if action is set
	if GetActions(core)["test"] == nil {
		t.Error("Actions setting is not working")
	}
}
