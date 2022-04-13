package actions

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"testing"
)

func buildActionRequest(state string, args string, link string) *http.Request {
	mpb := bytes.Buffer{}
	mpw := multipart.NewWriter(&mpb)
	mpw.WriteField("State", state)
	mpw.WriteField("Args", args)
	mpw.Close()
	req, err := http.NewRequest("POST", link, &mpb)
	req.Header.Set("Content-Type", mpw.FormDataContentType())
	if err != nil {
		panic(err)
	}
	return req
}

// TestParseParameters ensures that the ParseParameters function works as expected
func TestParseParameters(t *testing.T) {

	// Common request
	req1 := buildActionRequest(`{"Foo":"Bar"}`, `[]`, "http://localhost:8080/internal/actions/ComponentFoo/ActionName")
	// Custom route request
	req2 := buildActionRequest(`{"Foo":"Bar"}`, `[]`, "http://localhost:8080/custom-route/ComponentFoo/ActionName")

	// Validate common request parameters
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

	// Validate custom route request parameters
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

func TestParseParametersErrors(t *testing.T) {
	// Request with empty state
	req1 := buildActionRequest(``, `[]`, "http://localhost:8080/internal/actions/ComponentFoo/ActionName")
	// Request with empty args
	req2 := buildActionRequest(`{"Foo":"Bar"}`, ``, "http://localhost:8080/internal/actions/ComponentFoo/ActionName")
	// Request with invalid state
	req3 := buildActionRequest(`{"Foo":"Bar"`, `[]`, "http://localhost:8080/internal/actions/ComponentFoo/ActionName")
	// Request with invalid args
	req4 := buildActionRequest(`{"Foo":"Bar"}`, `[]]`, "http://localhost:8080/internal/actions/ComponentFoo/ActionName")

	for i, req := range []*http.Request{req1, req2, req3, req4} {
		_, err := ParseParameters(req)
		if err == nil {
			t.Errorf("ParseParameters(%s%d) should have returned error", "req", i)
		}
	}
}
