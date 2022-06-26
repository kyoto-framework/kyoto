package kyoto

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// Common test items

type testActionComponentState struct {
	Value string
}

func testActionSetup() {
	// Create a template file
	ioutil.WriteFile("action_test.html", []byte(`{{ define "testActionComponent" }}{{ .Value }}{{ end }}`), 0644)
}

func testActionCleanup() {
	// Remove template file
	os.Remove("action_test.html")
}

func testActionComponent(ctx *Context) (state testActionComponentState) {
	// Preload action state
	ActionPreload(ctx, &state)
	// Handle action
	handled := Action(ctx, "TestAction", func(args ...any) {
		state.Value = "Alternative value"
	})
	// Interrupt on handle
	if handled {
		return
	}
	// Default behavior
	state.Value = "Default value"
	// Return
	return
}

// TestAction ensures action call is working, not throws non-200 code and returns correct result.
func TestAction(t *testing.T) {
	// Setup
	testActionSetup()
	defer testActionCleanup()
	// Prepare request
	wstate := testActionComponentState{Value: "Default value"}
	reqb := bytes.Buffer{}
	reqw := multipart.NewWriter(&reqb)
	reqw.WriteField("State", MarshalState(wstate))
	reqw.WriteField("Args", "[]")
	reqw.Close()
	request, _ := http.NewRequest("POST", "/internal/actions/testActionComponent/TestAction", &reqb)
	request.Header.Set("Content-Type", reqw.FormDataContentType())
	// Make a request to the handler
	recorder := httptest.NewRecorder()
	handler := HandlerAction(testActionComponent)
	handler(recorder, request)
	// Check response results
	if recorder.Code != 200 {
		t.Error("Expected 200, got", recorder.Code)
	}
	if recorder.Body.String() != `Alternative value` {
		t.Errorf("Expected `Alternative value`, got `%s`", recorder.Body.String())
	}
}

// TestActionErrState ensures action will panic on empty state.
func TestActionErrState(t *testing.T) {
	// Define recovery
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in TestActionErrState")
		}
	}()
	// Setup
	testActionSetup()
	defer testActionCleanup()
	// Prepare request
	reqb := bytes.Buffer{}
	reqw := multipart.NewWriter(&reqb)
	reqw.WriteField("Args", "[]")
	reqw.Close()
	request, _ := http.NewRequest("POST", "/internal/actions/testActionComponent/TestAction", &reqb)
	request.Header.Set("Content-Type", reqw.FormDataContentType())
	// Make a request to the handler
	recorder := httptest.NewRecorder()
	handler := HandlerAction(testActionComponent)
	handler(recorder, request)
	// Panic expected
	t.Error("Expected panic, got", recorder.Code)
}

// TestActionErrArgs ensures action will panic on empty args ([] when no arguments).
func TestActionErrArgs(t *testing.T) {
	// Define recovery
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in TestActionErrArgs")
		}
	}()
	// Setup
	testActionSetup()
	defer testActionCleanup()
	// Prepare request
	wstate := testActionComponentState{Value: "Default value"}
	reqb := bytes.Buffer{}
	reqw := multipart.NewWriter(&reqb)
	reqw.WriteField("State", MarshalState(wstate))
	reqw.Close()
	request, _ := http.NewRequest("POST", "/internal/actions/testActionComponent/TestAction", &reqb)
	request.Header.Set("Content-Type", reqw.FormDataContentType())
	// Make a request to the handler
	recorder := httptest.NewRecorder()
	handler := HandlerAction(testActionComponent)
	handler(recorder, request)
	// Panic expected
	t.Error("Expected panic, got", recorder.Code)
}

// TestActionErrArgsCorrupted ensures action will panic on corrupted args.
func TestActionErrArgsCorrupted(t *testing.T) {
	// Define recovery
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in TestActionWrongArgsError")
		}
	}()
	// Setup
	testActionSetup()
	defer testActionCleanup()
	// Prepare request
	wstate := testActionComponentState{Value: "Default value"}
	reqb := bytes.Buffer{}
	reqw := multipart.NewWriter(&reqb)
	reqw.WriteField("State", MarshalState(wstate))
	reqw.WriteField("Args", "Baz")
	reqw.Close()
	request, _ := http.NewRequest("POST", "/internal/actions/testActionComponent/TestAction", &reqb)
	request.Header.Set("Content-Type", reqw.FormDataContentType())
	// Make a request to the handler
	recorder := httptest.NewRecorder()
	handler := HandlerAction(testActionComponent)
	handler(recorder, request)
	// Panic expected
	t.Error("Expected panic, got", recorder.Code)
}

// TestActionFuncState ensures state template function writes correct html attribute.
func TestActionFuncState(t *testing.T) {
	// Define a test state
	state := testActionComponentState{Value: "Default value"}
	// Compose into html attribute
	htmlattr := string(actionFuncState(state))
	// Check state key
	if !strings.Contains(htmlattr, "state") {
		t.Error("Expected state key, got", htmlattr)
	}
	// Check state value
	if !strings.Contains(htmlattr, MarshalState(state)) {
		t.Error("Expected state value, got", htmlattr)
	}
}

// TestActionFuncClient ensures client template function correctly writes JS client.
func TestActionFuncClient(t *testing.T) {
	// Get client
	client := actionFuncClient()
	// Check action client
	if !strings.Contains(string(client), ActionClient) {
		t.Error("Expected client, got", string(client))
	}
	// Check ssapath
	if !strings.Contains(string(client), ActionConf.Path) {
		t.Error("Expected path, got", string(client))
	}
}
