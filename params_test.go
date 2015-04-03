package wolf

import (
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestParams(t *testing.T) {
	ctx := context.Background()
	p := httprouter.Params{
		{Key: "foo", Value: "bar"},
		{Key: "foo", Value: "other"},
		{Key: "asdf", Value: "1234"},
	}

	newCtx := setParamsInContext(ctx, p)

	var (
		val  string
		vals []string
		ok   bool
	)
	val, ok = ParamFrom(newCtx, "foo")
	assert.True(t, ok)
	assert.Equal(t, "bar", val)

	val, ok = ParamFrom(newCtx, "asdf")
	assert.True(t, ok)
	assert.Equal(t, "1234", val)

	vals, ok = AllParamsFrom(newCtx, "foo")
	assert.True(t, ok)
	assert.Equal(t, []string{"bar", "other"}, vals)

	_, ok = ParamFrom(newCtx, "notfound")
	assert.False(t, ok)

	_, ok = AllParamsFrom(newCtx, "notfound")
	assert.False(t, ok)
}
