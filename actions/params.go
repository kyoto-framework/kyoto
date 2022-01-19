package actions

import (
	"encoding/base64"
	"encoding/json"
	"strings"
)

// Parameters represents parameters, needed for handling Actions request.
type Parameters struct {
	Component string
	Action    string
	State     map[string]interface{}
	Args      []interface{}
}

// ParseParameters is a function that parses parameters from request path.
func ParseParameters(path string) (Parameters, error) {
	params := Parameters{}
	tokens := strings.Split(path, "/")
	var _state, _args []byte
	_state, _ = base64.StdEncoding.DecodeString(strings.ReplaceAll(tokens[len(tokens)-3], "-", "/"))
	_args, _ = base64.StdEncoding.DecodeString(strings.ReplaceAll(tokens[len(tokens)-1], "-", "/"))
	state := map[string]interface{}{}
	if err := json.Unmarshal(_state, &state); err != nil {
		return params, err
	}
	args := []interface{}{}
	if err := json.Unmarshal(_args, &args); err != nil {
		return params, err
	}
	params.Component = tokens[len(tokens)-4]
	params.State = state
	params.Action = tokens[len(tokens)-2]
	params.Args = args
	return params, nil
}
