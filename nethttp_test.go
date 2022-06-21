package kyoto

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPageHandler(t *testing.T) {
	// Define component state
	type FooState struct {
		Bar string
	}

	// Define component
	Foo := func(ctx *Context) (state FooState) {
		state.Bar = "Baz"
		return state
	}

	// Define page state
	type PIndexState struct {
		Foo ComponentF[FooState]
	}

	// Define page
	PIndex := func(ctx *Context) (state PIndexState) {
		// Define rendering
		TemplateInline(ctx, `{{ template "Foo" await .Foo }}`)
		// Attach components
		state.Foo = Use(ctx, Foo)
		return state
	}
	// Start a local HTTP server
	server := httptest.NewServer(HandlerPage(PIndex))
	// Close the server when test finishes
	defer server.Close()

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

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	bodyString := string(bodyBytes)
	if bodyString != "Baz" {
		t.Errorf("Expected Baz, got %s", bodyString)
	}
}

func TestHandlePage(t *testing.T) {
	// Define component state
	type FooState struct {
		Bar string
	}

	// Define component
	Foo := func(ctx *Context) (state FooState) {
		state.Bar = "Baz"
		return state
	}

	// Define page state
	type PIndexState struct {
		Foo ComponentF[FooState]
	}

	// Define page
	PIndex := func(ctx *Context) (state PIndexState) {
		// Define rendering
		TemplateInline(ctx, `{{ template "Foo" await .Foo }}`)
		// Attach components
		state.Foo = Use(ctx, Foo)
		return state
	}

	// Start a local HTTP server
	HandlePage("/", PIndex)

	// Start serve
	go func() {
		Serve(":8080")
	}()

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

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	bodyString := string(bodyBytes)
	if bodyString != "Baz" {
		t.Errorf("Expected Baz, got %s", bodyString)
	}
}

func TestPageHandlerError(t *testing.T) {
	// Define recovery
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Expected recovery")
		}
	}()

	// Define component state
	type FooState struct {
		Bar string
	}

	// Define component
	Foo := func(ctx *Context) (state FooState) {
		state.Bar = "Baz"
		return state
	}

	// Define page state
	type PIndexState struct {
		Foo ComponentF[FooState]
	}

	// Define page
	PIndex := func(ctx *Context) (state PIndexState) {
		TemplateInline(ctx, `{{ template "Foo" await .Bar }}`)
		// Attach components
		state.Foo = Use(ctx, Foo)
		return state
	}

	// Start a local HTTP server
	server := httptest.NewServer(HandlerPage(PIndex))

	// Close the server when test finishes
	defer server.Close()

	// Get request to page with wrong template
	_, err := http.Get(fmt.Sprintf("%s/", server.URL))
	if err == nil {
		t.Fatal("Expected error")
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
	t.Error("Expected error not occured")
}
