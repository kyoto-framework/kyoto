package component_test

import (
	"fmt"
	"html/template"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.kyoto.codes/v3/component"
)

func TestFuncMap(t *testing.T) {
	t.Parallel()
	// New state
	state := component.Use(component.NewContext(nil, nil), testComponent)()
	// Extract and type assert function
	funcMapMarshal, success := component.FuncMap["marshal"].(func(component.State) string)
	// Ensure for successed type assertion
	assert.True(t, success)
	// Assert marshal result
	assert.Equal(t, state.Marshal(), funcMapMarshal(state))

	// Extract and type assert function
	componentFunc, success := component.FuncMap["component"].(func(component.State) template.HTMLAttr)
	// Ensure for successed type assertion
	assert.True(t, success)
	// Assert component result
	assert.Equal(t, template.HTMLAttr(fmt.Sprintf(`component="%s" state="%s"`, state.GetName(), state.Marshal())), componentFunc(state))
}
