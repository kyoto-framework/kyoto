package kyoto

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestActionComponent(t *testing.T) {
	// Define component state
	type testActionComponentState struct {
		Foo string
	}

	// Create component
	component := func(ctx *Context) (state testActionComponentState) {
		// Preload state
		ActionPreload(ctx, &state)

		// Execute action
		executed := Action(ctx, "Baz", func(args ...any) {
			state.Foo = "Baz"
		})

		// Interrupt if executed
		if executed {
			return state
		}

		// Default behavior
		return state
	}

	// Handle action
	HandleAction(component)
	// Handler action
	HandlerAction(component)
}

func TestActionNotExecuted(t *testing.T) {
	// Define component state
	type testActionComponentState struct {
		Foo string
	}

	// Create component
	component := func(ctx *Context) (state testActionComponentState) {
		// Preload state
		ActionPreload(ctx, &state)
		// Execute action
		executed := Action(ctx, "Baz", func(args ...any) {
			state.Foo = "Baz"
		})

		// Interrupt if executed
		if executed {
			return state
		}

		// Default behavior
		return state
	}

	// Define Page state
	type testPageState struct {
		Foo *ComponentF[testActionComponentState]
	}

	// Create Page
	page := func(ctx *Context) (state testPageState) {
		state.Foo = Use(ctx, component)
		return state
	}

	page(&Context{
		Request:        &http.Request{},
		ResponseWriter: httptest.NewRecorder(),
		Action: ActionParameters{
			Component: "",
			Action:    "",
			Args:      []any{},
		},
	})
}

func TestAction(t *testing.T) {
	// Define component state
	type testActionComponentState struct {
		Bar string
	}

	// Create component
	component := func(ctx *Context) (state testActionComponentState) {
		// Preload action
		ActionPreload(ctx, &state)
		// Execute action
		executed := Action(ctx, "Baz", func(args ...any) {
			state.Bar = "Baz"
		})
		// Interrupt if executed
		if executed {
			return state
		}
		// Default behavior
		state.Bar = "Bar"
		return state
	}
	// Prepare request
	wstate := testActionComponentState{Bar: "Bar"}
	reqb := bytes.Buffer{}
	reqw := multipart.NewWriter(&reqb)
	reqw.WriteField("State", MarshalState(wstate))
	reqw.WriteField("Args", "[]")
	reqw.WriteField("Component", "component")
	reqw.WriteField("Action", "Baz")
	reqw.Close()
	request, _ := http.NewRequest("POST", "/internal/actions/component/Baz", &reqb)
	request.Header.Set("Content-Type", reqw.FormDataContentType())
	// Make a request to the handler
	recorder := httptest.NewRecorder()
	handler := HandlerAction(component)
	handler(recorder, request)
	// Check response results
	if recorder.Code != 200 {
		t.Error("Expected 200, got", recorder.Code)
	}
	if recorder.Body.String() != `Baz` {
		t.Error("Expected `Baz`, got", recorder.Body.String())
	}
}

func TestActionEmptyStateError(t *testing.T) {
	// Define recovery
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in TestActionEmptyStateError")
		}
	}()

	// Define component state
	type testActionComponentState struct {
		Bar string
	}

	// Create component
	component := func(ctx *Context) (state testActionComponentState) {
		// Preload action
		ActionPreload(ctx, &state)
		// Execute action
		executed := Action(ctx, "Baz", func(args ...any) {
			state.Bar = "Baz"
		})
		// Interrupt if executed
		if executed {
			return state
		}
		// Default behavior
		state.Bar = "Bar"
		return state
	}
	// Prepare request
	reqb := bytes.Buffer{}
	reqw := multipart.NewWriter(&reqb)
	reqw.WriteField("Args", "[]")
	reqw.WriteField("Component", "component")
	reqw.WriteField("Action", "Baz")
	reqw.Close()
	request, _ := http.NewRequest("POST", "/internal/actions/component/Baz", &reqb)
	request.Header.Set("Content-Type", reqw.FormDataContentType())
	// Make a request to the handler
	recorder := httptest.NewRecorder()
	handler := HandlerAction(component)
	handler(recorder, request)

	t.Error("Expected panic, got", recorder.Code)
}

func TestActionEmptyArgsError(t *testing.T) {
	// Define recovery
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in TestActionEmptyArgsError")
		}
	}()

	// Define component state
	type testActionComponentState struct {
		Bar string
	}

	// Create component
	component := func(ctx *Context) (state testActionComponentState) {
		// Preload action
		ActionPreload(ctx, &state)
		// Execute action
		executed := Action(ctx, "Baz", func(args ...any) {
			state.Bar = "Baz"
		})
		// Interrupt if executed
		if executed {
			return state
		}
		// Default behavior
		state.Bar = "Bar"
		return state
	}
	// Prepare request
	wstate := testActionComponentState{Bar: "Bar"}
	reqb := bytes.Buffer{}
	reqw := multipart.NewWriter(&reqb)
	reqw.WriteField("State", MarshalState(wstate))
	reqw.WriteField("Component", "component")
	reqw.WriteField("Action", "Baz")
	reqw.Close()
	request, _ := http.NewRequest("POST", "/internal/actions/component/Baz", &reqb)
	request.Header.Set("Content-Type", reqw.FormDataContentType())
	// Make a request to the handler
	recorder := httptest.NewRecorder()
	handler := HandlerAction(component)
	handler(recorder, request)

	t.Error("Expected panic, got", recorder.Code)
}

func TestActionWrongArgsError(t *testing.T) {
	// Define recovery
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in TestActionWrongArgsError")
		}
	}()

	// Define component state
	type testActionComponentState struct {
		Bar string
	}

	// Create component
	component := func(ctx *Context) (state testActionComponentState) {
		// Preload action
		ActionPreload(ctx, &state)
		// Execute action
		executed := Action(ctx, "Baz", func(args ...any) {
			state.Bar = "Baz"
		})
		// Interrupt if executed
		if executed {
			return state
		}
		// Default behavior
		state.Bar = "Bar"
		return state
	}
	// Prepare request
	wstate := testActionComponentState{Bar: "Bar"}
	reqb := bytes.Buffer{}
	reqw := multipart.NewWriter(&reqb)
	reqw.WriteField("State", MarshalState(wstate))
	reqw.WriteField("Args", "Baz")
	reqw.WriteField("Component", "component")
	reqw.WriteField("Action", "Baz")
	reqw.Close()
	request, _ := http.NewRequest("POST", "/internal/actions/component/Baz", &reqb)
	request.Header.Set("Content-Type", reqw.FormDataContentType())
	// Make a request to the handler
	recorder := httptest.NewRecorder()
	handler := HandlerAction(component)
	handler(recorder, request)

	t.Error("Expected panic, got", recorder.Code)
}

func TestActionFuncState(t *testing.T) {
	// Define component state
	type componentState struct {
		Bar string
	}

	state := componentState{Bar: "Bar"}

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

func TestActionFuncClient(t *testing.T) {
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
