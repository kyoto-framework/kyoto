package helpers

import (
	"testing"
)

func testFunction() {}

type testComponent struct{}

func TestComponentID(t *testing.T) {
	// Create test objects
	obj1 := map[int]int{}
	obj2 := map[int]int{}
	// Check ids are not empty
	if ComponentID(obj1) == "" || ComponentID(obj2) == "" {
		t.Error("ComponentID returned empty string")
	}
	// Compare object ids
	if ComponentID(obj1) == ComponentID(obj2) {
		t.Error("ComponentID is not producing unique ids")
	}
}

func TestComponentName(t *testing.T) {
	// Create test objects
	obj1 := map[string]interface{}{"internal:name": "testcomponent"}
	// Check functions
	if ComponentName(testFunction) != "testFunction" {
		t.Errorf("ComponentName is not working correctly for functions. %s != %s", ComponentName(testFunction), "testFunction")
	}
	// Check structs
	if ComponentName(&testComponent{}) != "testComponent" {
		t.Error("ComponentName is not working correctly for structs")
	}
	// Check maps
	if ComponentName(obj1) != "testcomponent" {
		t.Error("ComponentName is not working correctly for maps")
	}
}

func TestComponentNamePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	// Trigger panic with wrong type
	ComponentName(testComponent{})
}

func TestComponentSerialize(t *testing.T) {
	// Create test objects
	obj1 := map[string]interface{}{
		"internal:name": "testcomponent",
	}
	obj2 := map[string]interface{}{
		"Foo": "Bar",
	}
	// Check serialization
	if ComponentSerialize(obj1) == "eyJpbnRlcm5hbDpuYW1lIjoidGVzdGNvbXBvbmVudCJ9" {
		t.Error("ComponentSerialize is not clearing internal variables")
	}
	if ComponentSerialize(obj2) != "eyJGb28iOiJCYXIifQ==" {
		t.Error("ComponentSerialize is not working correctly")
	}
}

func TestComponentSerializePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	// Create test objects
	obj1 := map[string]interface{}{
		"Foo": func() {},
	}
	// Trigger panic
	ComponentSerialize(obj1)
}
