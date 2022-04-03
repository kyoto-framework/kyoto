package actions

import (
	"bytes"
	"html/template"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kyoto-framework/kyoto"
	"github.com/kyoto-framework/kyoto/lifecycle"
)

func testHandlerComponent(c *kyoto.Core) {
	lifecycle.Init(c, func() {
		c.State.Set("Foo", "Bar")
	})
	Define(c, "Baz", func(args ...interface{}) {
		c.State.Set("Foo", "Baz")
	})
}

func TestHandler(t *testing.T) {
	// Register component
	RegisterWithName("testHandlerComponent", testHandlerComponent)
	// Define handler
	handler := Handler(func() *template.Template {
		return template.Must(template.New("").Parse(
			`{{ define "testHandlerComponent" }}` +
				`Content is {{ .Foo }}` +
				`{{ end }}`,
		))
	})
	// Build a test request
	reqb := bytes.Buffer{}
	reqw := multipart.NewWriter(&reqb)
	reqw.WriteField("State", `{"Foo":"Bar"}`)
	reqw.WriteField("Args", "[]")
	reqw.Close()
	request, _ := http.NewRequest("POST", "/internal/actions/testHandlerComponent/Baz", &reqb)
	request.Header.Set("Content-Type", reqw.FormDataContentType())
	// Make a request to the handler
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	// Check response results
	if recorder.Code != 200 {
		t.Error("Expected 200, got", recorder.Code)
	}
	if recorder.Body.String() != "Content is Baz" {
		t.Error("Expected Baz, got", recorder.Body.String())
	}
}
