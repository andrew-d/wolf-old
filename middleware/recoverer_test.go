package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	"github.com/andrew-d/wolf"
)

// Test that the Recoverer will recover from panics
func TestRecoverer(t *testing.T) {
	a := wolf.New()
	a.Use(RequestID)
	a.Use(Recoverer)

	var run bool
	a.Get("/", func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		run = true
		panic("foo bar")
	})

	recorder := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)
	a.ServeHTTP(recorder, r)

	// The recoverer should have caught the panic (i.e. this code will run), and
	// return a 500 error.
	assert.True(t, run)
	assert.Equal(t, 500, recorder.Code)
}

// Test that the CustomRecoverer will pass appropriate information to the callback.
func TestCustomRecoverer(t *testing.T) {
	var (
		info  RecoverInformation
		cbRun bool
	)
	cb := func(w http.ResponseWriter, r *http.Request, i RecoverInformation) {
		cbRun = true
		info = i
	}

	a := wolf.New()
	a.Use(RequestID)
	a.Use(CustomRecoverer(cb))

	var run bool
	a.Get("/", func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		run = true
		panic("foo bar")
	})

	recorder := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)
	a.ServeHTTP(recorder, r)

	// Both the handler and the callback should have run
	assert.True(t, cbRun)
	assert.True(t, run)

	// The information should be valid - proper error, request ID, etc.
	assert.NotEmpty(t, info.RequestID)
	assert.NotEmpty(t, info.Stack)

	assert.Equal(t, "foo bar", info.Error)
}
