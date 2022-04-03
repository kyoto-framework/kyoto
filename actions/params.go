package actions

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

// Parameters represents parameters, needed for handling Actions request.
type Parameters struct {
	Component string
	Action    string
	State     map[string]interface{}
	Args      []interface{}
}

func ParseParameters(r *http.Request) (Parameters, error) {
	// Validate request format
	if r.FormValue("State") == "" {
		return Parameters{}, errors.New("State is empty")
	}
	if r.FormValue("Args") == "" {
		return Parameters{}, errors.New("Args is empty")
	}
	// Initialize parameters store
	parameters := Parameters{}
	// Split path into tokens
	tokens := strings.Split(r.URL.Path, "/")
	// Extract component state
	err := json.Unmarshal([]byte(r.FormValue("State")), &parameters.State)
	if err != nil {
		return parameters, errors.New("Something wrong with state")
	}
	// Extract component arguments
	err = json.Unmarshal([]byte(r.FormValue("Args")), &parameters.Args)
	if err != nil {
		return parameters, errors.New("Something wrong with arguments")
	}
	// Extract component & action names
	parameters.Component = tokens[len(tokens)-2]
	parameters.Action = tokens[len(tokens)-1]
	// Return
	return parameters, nil
}
