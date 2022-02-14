package actions

import "testing"

// TestParseParameters ensures that the ParseParameters function works as expected
func TestParseParameters(t *testing.T) {
	path1 := "/internal/actions/ComponentFoo/eyJGb28iOiJCYXIifQ==/ActionName/W10="
	path2 := "/custom-route/ComponentFoo/eyJGb28iOiJCYXIifQ==/ActionName/W10="

	params1, err := ParseParameters(path1)
	if err != nil {
		t.Errorf("ParseParameters(%s) returned error: %v", path1, err)
	}
	if params1.Component != "ComponentFoo" {
		t.Errorf("ParseParameters(%s) returned wrong component: %s", path1, params1.Component)
	}
	if params1.Action != "ActionName" {
		t.Errorf("ParseParameters(%s) returned wrong action: %s", path1, params1.Action)
	}
	if params1.State["Foo"] != "Bar" {
		t.Errorf("ParseParameters(%s) returned wrong state: %v", path1, params1.State)
	}

	params2, err := ParseParameters(path2)
	if err != nil {
		t.Errorf("ParseParameters(%s) returned error: %v", path2, err)
	}
	if params2.Component != "ComponentFoo" {
		t.Errorf("ParseParameters(%s) returned wrong component: %s", path2, params2.Component)
	}
	if params2.Action != "ActionName" {
		t.Errorf("ParseParameters(%s) returned wrong action: %s", path2, params2.Action)
	}
	if params2.State["Foo"] != "Bar" {
		t.Errorf("ParseParameters(%s) returned wrong state: %v", path2, params2.State)
	}
}
