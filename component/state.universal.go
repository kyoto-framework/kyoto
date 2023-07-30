package component

import (
	"encoding/base64"
	"encoding/json"
	"net/url"

	"go.kyoto.codes/zen/v3/errorsx"
	"go.kyoto.codes/zen/v3/jsonx"
)

// Universal is a universal component state implementation.
// It uses combination of JSON, base64 and URI encoding
// to marshal and unmarshal the state.
type Universal struct {
	Name string
}

func (c *Universal) SetName(name string) {
	c.Name = name
}

func (c *Universal) GetName() string {
	return c.Name
}

func (*Universal) Marshal(state any) string {
	// Marshal into json
	stateJson := jsonx.String(state)
	// Encode to URI representation to avoid html breaking
	stateJsonUri := url.PathEscape(stateJson)
	// Encode to base64
	stateJsonUriBase64 := base64.StdEncoding.EncodeToString([]byte(stateJsonUri))
	// Return
	return stateJsonUriBase64
}

func (*Universal) Unmarshal(str string, state any) error {
	// Decode from base64
	stateJsonUri := errorsx.Must(base64.StdEncoding.DecodeString(str))
	// Decode from URI representation
	stateJson := errorsx.Must(url.PathUnescape(string(stateJsonUri)))
	// Unmarshal from json
	return json.Unmarshal([]byte(stateJson), state)
}
