package actions

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kyoto-framework/kyoto"
	"github.com/kyoto-framework/kyoto/lifecycle"
)

func testHandlerComponentClassic(c *kyoto.Core) {
	lifecycle.Init(c, func() {
		c.State.Set("Foo", "Bar")
	})
	Define(c, "Baz", func(args ...interface{}) {
		c.State.Set("Foo", "Baz")
	})
}

func testHandlerComponentWriter(c *kyoto.Core) {
	lifecycle.Init(c, func() {
		c.State.Set("Foo", "Bar")
	})
	Define(c, "Baz", func(args ...interface{}) {
		c.State.Set("Foo", "Baz")
	})
	// Injected directly to avoid circular dependency
	c.State.Set("internal:render:writer", func(w io.Writer) error {
		w.Write([]byte("Content is "))
		w.Write([]byte(c.State.Get("Foo").(string)))
		return nil
	})
}

func testHandlerComponentWriterErr(c *kyoto.Core) {
	lifecycle.Init(c, func() {
		c.State.Set("Foo", "Bar")
	})
	Define(c, "Baz", func(args ...interface{}) {
		c.State.Set("Foo", "Baz")
	})
	// Injected directly to avoid circular dependency
	c.State.Set("internal:render:writer", func(w io.Writer) error {
		return errors.New("test error")
	})
}

func TestHandlerClassic(t *testing.T) {
	// Register component
	RegisterWithName("testHandlerComponentClassic", testHandlerComponentClassic)
	// Define handler
	handler := Handler(func(c *kyoto.Core) *template.Template {
		return template.Must(template.New("").Parse(
			`{{ define "testHandlerComponentClassic" }}` +
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
	request, _ := http.NewRequest("POST", "/internal/actions/testHandlerComponentClassic/Baz", &reqb)
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

func TestHandlerWriter(t *testing.T) {
	// Register component
	RegisterWithName("testHandlerComponentWriter", testHandlerComponentWriter)
	// Define handler
	handler := Handler(func(c *kyoto.Core) *template.Template {
		return template.Must(template.New("").Parse(
			`{{ define "testHandlerComponentWriter" }}` +
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
	request, _ := http.NewRequest("POST", "/internal/actions/testHandlerComponentWriter/Baz", &reqb)
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

func TestHandlerWriterErr(t *testing.T) {
	// Register component
	RegisterWithName("testHandlerComponentWriterErr", testHandlerComponentWriterErr)
	// Define handler
	handler := Handler(func(c *kyoto.Core) *template.Template {
		return template.Must(template.New("").Parse(
			`{{ define "testHandlerComponentWriter" }}` +
				`Content is {{ .Foo }}` +
				`{{ end }}`,
		))
	})
	// Build a test request
	request := buildActionRequest(`{"Foo":"Bar"}`, `[]`, "http://localhost:8080/internal/actions/testHandlerComponentWriterErr/Baz")
	// Make a request to the handler
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	// Check response status code
	if recorder.Code != 500 {
		t.Error("Expected 500, got", recorder.Code)
	}
}
