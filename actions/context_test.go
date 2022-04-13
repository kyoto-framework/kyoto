package actions

import (
	"fmt"
	"testing"

	"github.com/kyoto-framework/kyoto"
	"github.com/kyoto-framework/kyoto/helpers"
)

// TestActionsGetSet ensures that the actions get/set for context works as expected
func TestActionsGetSet(t *testing.T) {
	// Innitialize core
	core := kyoto.NewCore()
	// Determine actions address
	contextaddr := fmt.Sprintf("internal:actions:%s", helpers.ComponentID(core))
	// Set actions
	SetActions(core, ActionMap{
		`test`: func(args ...interface{}) {},
	})
	// Test if actions are set
	if core.Context.Get(contextaddr).(ActionMap) == nil {
		t.Error("Actions setting is not working")
	}
	// Get actions
	if len(GetActions(core)) == 0 {
		t.Error("Actions getting is not working")
	}
	// Ensure default is working
	if GetActions(kyoto.NewCore()) == nil {
		t.Error("Actions getting with default is not working")
	}
}
