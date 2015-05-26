package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	"github.com/andrew-d/wolf"
)

// Test that the request ID is added and non-empty.
func TestRequestID(t *testing.T) {
	a := wolf.New()
	a.Use(RequestID)

	var run bool
	a.Get("/", func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		run = true
		assert.True(t, len(GetReqID(ctx)) > 0)
	})

	var w http.ResponseWriter = httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)
	a.ServeHTTP(w, r)

	assert.True(t, run)
}

// Test that retrieving a non-existant request ID doesn't fail.
func TestNonExistantRequestID(t *testing.T) {
	a := wolf.New()

	var run bool
	a.Get("/", func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		run = true
		assert.Equal(t, "", GetReqID(ctx))
	})

	var w http.ResponseWriter = httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)
	a.ServeHTTP(w, r)

	assert.True(t, run)
}
