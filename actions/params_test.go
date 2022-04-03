package actions

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"testing"
)

// TestParseParameters ensures that the ParseParameters function works as expected
func TestParseParameters(t *testing.T) {
	mpb1 := bytes.Buffer{}
	mpw1 := multipart.NewWriter(&mpb1)
	mpw1.WriteField("State", `{"Foo":"Bar"}`)
	mpw1.WriteField("Args", "[]")
	mpw1.Close()
	req1, err := http.NewRequest("POST", "http://localhost:8080/internal/actions/ComponentFoo/ActionName", &mpb1)
	req1.Header.Set("Content-Type", mpw1.FormDataContentType())
	if err != nil {
		panic(err)
	}

	mpb2 := bytes.Buffer{}
	mpw2 := multipart.NewWriter(&mpb2)
	mpw2.WriteField("State", `{"Foo":"Bar"}`)
	mpw2.WriteField("Args", "[]")
	mpw2.Close()
	req2, err := http.NewRequest("POST", "http://localhost:8080/custom-route/ComponentFoo/ActionName", &mpb2)
	req2.Header.Set("Content-Type", mpw2.FormDataContentType())
	if err != nil {
		panic(err)
	}

	params1, err := ParseParameters(req1)
	if err != nil {
		t.Errorf("ParseParameters(%s) returned error: %v", "req1", err)
	}
	if params1.Component != "ComponentFoo" {
		t.Errorf("ParseParameters(%s) returned wrong component: %s", "req1", params1.Component)
	}
	if params1.Action != "ActionName" {
		t.Errorf("ParseParameters(%s) returned wrong action: %s", "req1", params1.Action)
	}
	if params1.State["Foo"] != "Bar" {
		t.Errorf("ParseParameters(%s) returned wrong state: %v", "req1", params1.State)
	}

	params2, err := ParseParameters(req2)
	if err != nil {
		t.Errorf("ParseParameters(%s) returned error: %v", "req2", err)
	}
	if params2.Component != "ComponentFoo" {
		t.Errorf("ParseParameters(%s) returned wrong component: %s", "req2", params2.Component)
	}
	if params2.Action != "ActionName" {
		t.Errorf("ParseParameters(%s) returned wrong action: %s", "req2", params2.Action)
	}
	if params2.State["Foo"] != "Bar" {
		t.Errorf("ParseParameters(%s) returned wrong state: %v", "req2", params2.State)
	}
}
