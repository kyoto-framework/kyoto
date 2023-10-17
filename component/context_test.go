package component_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.kyoto.codes/v3/component"
)

func TestContext(t *testing.T) {
	t.Parallel()
	// Create new context
	ctx := component.NewContext(nil, nil)
	// Assert field values
	assert.Nil(t, ctx.Request)
	assert.Nil(t, ctx.ResponseWriter)
	// Assert new store value
	ctx.Set("key", "value")
	assert.Equal(t, "value", ctx.Get("key"))
}
