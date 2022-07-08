package kyoto

import (
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.sr.ht/~kyoto-framework/zen"
)

// Common test items

type testNethttpComponentState struct {
	Value string
}

func testNethttpComponent(ctx *Context) (state testNethttpComponentState) {
	state.Value = "This is a value"
	return
}

type testNethttpPageState struct {
	Component *ComponentF[testNethttpComponentState]
}

func testNethttpPage(ctx *Context) (state testNethttpPageState) {
	// Define rendering
	TemplateRaw(ctx, template.Must(template.New("nethttp.page.html").Funcs(FuncMap).Parse(
		`{{ define "nethttpComponent" }}{{ .Value }}{{ end }}`+
			`<html>{{ template "nethttpComponent" await .Component }}</html>`,
	)))
	// Attach component
	state.Component = Use(ctx, testNethttpComponent)
	// Return
	return
}

func testNethttpPageErr(ctx *Context) (state testNethttpPageState) {
	// Define rendering with an error (nethttpComponent is not defined)
	TemplateRaw(ctx, template.Must(template.New("nethttp.page.html").Funcs(FuncMap).Parse(
		`<html>{{ template "nethttpComponent" await .Component }}</html>`,
	)))
	// Attach component
	state.Component = Use(ctx, testNethttpComponent)
	// Return
	return
}

func TestPageHandler(t *testing.T) {
	// Start a local HTTP server with a testNethttpPage handler
	server := httptest.NewServer(HandlerPage(testNethttpPage))
	// Close the server when test finishes
	defer server.Close()

	// Make a request
	res, err := http.Get(fmt.Sprintf("%s/", server.URL))
	if err != nil {
		t.Fatal(err)
	}

	// Check status code
	if res.StatusCode != http.StatusOK {
		t.Errorf("status not OK")
	}

	// Check response body
	defer res.Body.Close()
	if body := zen.Response(res).Text(); body != "<html>This is a value</html>" {
		t.Error("Expected <html>This is a value</html>, got", body)
	}
}

func TestHandlePage(t *testing.T) {
	// Start a local HTTP server with a testNethttpPage automatic handler
	HandlePage("/", testNethttpPage)
	// Start serve
	go func() {
		Serve(":8080")
	}()

	// Make a request
	res, err := http.Get(fmt.Sprintf("%s/", "http://localhost:8080"))
	if err != nil {
		t.Fatal(err)
	}

	// Check status code
	if res.StatusCode != http.StatusOK {
		t.Errorf("status not OK")
	}

	// Check response body
	defer res.Body.Close()
	if body := zen.Response(res).Text(); body != "<html>This is a value</html>" {
		t.Error("Expected <html>This is a value</html>, got", body)
	}
}

func TestPageHandlerError(t *testing.T) {
	// Define recovery
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Expected recovery")
		}
	}()

	// Start a local HTTP server with a testNethttpPageErr handler
	server := httptest.NewServer(HandlerPage(testNethttpPageErr))
	// Close the server when test finishes
	defer server.Close()

	// Get request to page with wrong template
	_, err := http.Get(fmt.Sprintf("%s/", server.URL))
	if err == nil {
		t.Fatal("Expected error, got OK")
	}
}

func TestServeError(t *testing.T) {
	// Define recovery
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Expected recovery")
		}
	}()
	// Serve with invalid port
	Serve("/")

	// Error if not panicked
	t.Error("Expected error not occurred")
}
