package smode

import (
	"testing"

	"github.com/kyoto-framework/kyoto"
)

func TestContext(t *testing.T) {
	// Initialize core, page and apply it to core
	core := kyoto.NewCore()
	page := &testPage{}
	Adapt(page)(core)

	// Validate setting works as expected
	SetContext(page, "test", "value")
	if core.Context.Get("test") == nil {
		t.Error("SetContext not working")
	}

	// Validate getting works as expected
	if GetContext(page, "test") != "value" {
		t.Error("GetContext not working")
	}

	// Validate delete works as expected
	DelContext(page, "test")
	if core.Context.Get("test") != nil {
		t.Error("DelContext not working")
	}
}
