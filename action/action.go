package action

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"go.kyoto.codes/zen/v3/errorsx"
)

// Action parameters and state,
// passed through context.
type Action struct {
	Component string
	Action    string
	State     string // marshalled component state
	Args      []any

	// handled flag
	// allows to avoid recursive action call
	// while attaching recursive component.
	handled bool
	// rendered flag
	// allows to avoid double-rendering
	// or headers double-writing.
	// Please note, redirection is also considered "rendered".
	rendered bool
}

// UnmarshalHttpRequest allows to load action parameters
// from http.Request.
func (a *Action) UnmarshalHttpRequest(r *http.Request) error {
	// Validate request format
	if r.FormValue("State") == "" {
		return errors.New("state is empty")
	}
	if r.FormValue("Args") == "" {
		return errors.New("args is empty")
	}
	// Split path into tokens
	tokens := strings.Split(r.URL.Path, "/")
	// Extract component state
	a.State = r.FormValue("State")
	// Extract component arguments
	errorsx.Must(0, json.Unmarshal([]byte(r.FormValue("Args")), &a.Args))
	// Extract component & action names
	a.Component = tokens[len(tokens)-2]
	a.Action = tokens[len(tokens)-1]
	// Return
	return nil
}
