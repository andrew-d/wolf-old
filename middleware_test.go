package wolf

import (
	"net/http"
	"net/http/httptest"
	"testing"

	//"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestMiddlewareTypes(t *testing.T) {
	a := New()

	a.Use(func(ctx *context.Context, h http.Handler) http.Handler { return nil })
	a.Use(func(h http.Handler) http.Handler { return nil })

	assert.Panics(t, func() {
		a.Use(func(i int) int { return i + 1 })
	})
}

func TestMiddlewareOrder(t *testing.T) {
	a := New()

	var calls []string

	// Verify that our middleware type works
	a.Use(func(ctx *context.Context, h http.Handler) http.Handler {
		wrap := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls = append(calls, "one")
			h.ServeHTTP(w, r)
		})
		return wrap
	})

	// Standard library-ish middleware type should work too
	a.Use(func(h http.Handler) http.Handler {
		wrap := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls = append(calls, "two")
			h.ServeHTTP(w, r)
		})
		return wrap
	})

	var run bool
	a.Get("/:param", func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		run = true
		assert.NotNil(t, ctx)
	})

	var w http.ResponseWriter = httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/foo", nil)
	assert.NoError(t, err)
	a.ServeHTTP(w, r)

	assert.True(t, run)
	assert.Equal(t, []string{"one", "two"}, calls)
}
