package render

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kyoto-framework/kyoto"
	"github.com/kyoto-framework/kyoto/lifecycle"
)

// TestPageHandler ensures that the page handler actually executes core and renders the page
func TestPageHandler(t *testing.T) {
	// Initialize test page
	page := func(c *kyoto.Core) {
		Template(c, func() *template.Template {
			return template.Must(template.New("test").Parse("{{ .Foo }}"))
		})
		lifecycle.Init(c, func() {
			c.State.Set("Foo", "Bar")
		})
	}
	// Make a handler
	handler := PageHandler(page)
	// Make a test request
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/", nil)
	handler.ServeHTTP(recorder, request)
	// Check response results
	if recorder.Code != 200 {
		t.Error("Expected 200, got", recorder.Code)
	}
	if recorder.Body.String() != "Bar" {
		t.Error("Expected Bar, got", recorder.Body.String())
	}
}

// TestPageHandlerCustomRender ensures that page custom render inside of page handler actually works
func TestPageHandlerWriterRender(t *testing.T) {
	// Initialize test page
	page := func(c *kyoto.Core) {
		Writer(c, func(w io.Writer) error {
			w.Write([]byte(fmt.Sprintf(
				"%s", c.State.Get("Foo"),
			)))
			return nil
		})
		Template(c, func() *template.Template {
			return template.Must(template.New("test").Parse("{{ .Foo }}"))
		})
		lifecycle.Init(c, func() {
			c.State.Set("Foo", "Bar")
		})
	}
	// Make a handler
	handler := PageHandler(page)
	// Make a test request
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/", nil)
	handler.ServeHTTP(recorder, request)
	// Check response results
	if recorder.Code != 200 {
		t.Error("Expected 200, got", recorder.Code)
	}
	if recorder.Body.String() != "Bar" {
		t.Error("Expected Bar, got", recorder.Body.String())
	}
}

// TestPageHandlerRedirect ensures that redirect inside of page handler actually works
func TestPageHandlerRedirect(t *testing.T) {
	// Initialize test page
	page := func(c *kyoto.Core) {
		Template(c, func() *template.Template {
			return template.Must(template.New("test").Parse("{{ .Foo }}"))
		})
		lifecycle.Init(c, func() {
			Redirect(c, "/foo", 302)
		})
	}
	// Make a handler
	handler := PageHandler(page)
	// Make a test request
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/", nil)
	handler.ServeHTTP(recorder, request)
	// Check response results
	if recorder.Code != 302 {
		t.Error("Expected 302, got", recorder.Code)
	}
	if recorder.Header().Get("Location") != "/foo" {
		t.Error("Expected /foo, got", recorder.Header().Get("Location"))
	}
}
