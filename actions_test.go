package kyoto

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
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
