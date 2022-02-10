package lifecycle

import (
	"testing"

	"github.com/kyoto-framework/kyoto"
)

func TestPatch(t *testing.T) {
	// Initialize core
	core := kyoto.NewCore()

	// Patch core
	Patch(core)

	// Validate context
	if core.Context.Get("internal:lifecycle") == nil {
		t.Error("Lifecycle patch is not setting context flag")
	}

	// Validate placeholder jobs count
	if len(core.Scheduler.Jobs) != 3 {
		t.Error("Lifecycle patch is not providing placeholder jobs")
	}
	// Validate placeholder jobs groups
	init := false
	async := false
	afterasync := false
	for _, job := range core.Scheduler.Jobs {
		if job.Group == "init" {
			init = true
		} else if job.Group == "async" {
			async = true
		} else if job.Group == "afterasync" {
			afterasync = true
		}
	}
	if !init {
		t.Error("Lifecycle patch is not providing placeholder job for init group")
	}
	if !async {
		t.Error("Lifecycle patch is not providing placeholder job for async group")
	}
	if !afterasync {
		t.Error("Lifecycle patch is not providing placeholder job for afterasync group")
	}
}
