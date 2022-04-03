package actions

import (
	"testing"

	"github.com/kyoto-framework/kyoto"
	"github.com/kyoto-framework/kyoto/helpers"
)

func testRegisterComponent(c *kyoto.Core) {}

func TestRegister(t *testing.T) {
	// Register component
	Register(testRegisterComponent)
	// Check component is registered
	if _, ok := registry[helpers.ComponentName(testRegisterComponent)]; !ok {
		t.Error("Component is not registered")
	}
}

func TestRegisterWithName(t *testing.T) {
	// Register component
	RegisterWithName("fooBar", testRegisterComponent)
	// Check component is registered
	if _, ok := registry["fooBar"]; !ok {
		t.Error("Component is not registered")
	}
}
