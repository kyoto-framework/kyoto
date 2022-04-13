package smode

import "testing"

type testStructmap struct {
	Foo string
	Bar string `json:"-"`
}

func TestStructmap(t *testing.T) {
	// Initialize struct
	obj := testStructmap{
		Foo: "Bar",
		Bar: "Baz",
	}
	// Create map
	objmap := structmap(obj)
	// Check Foo is in map
	if objmap["Foo"] != "Bar" {
		t.Error("Map value is missing")
	}
	// Check Bar is not in map
	if objmap["Bar"] != nil {
		t.Error("structmap is not filtering accoring to json tag")
	}
}
