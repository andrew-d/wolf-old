package wolf

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	//"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

type dummyHandler struct{}

func (d dummyHandler) ServeHTTPCtx(ctx context.Context, w http.ResponseWriter, r *http.Request) {}

type dummyStdHandler struct{}

func (d dummyStdHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

func TestMakeHandler(t *testing.T) {
	// Our handler type
	assert.NotNil(t, MakeHandler(dummyHandler{}))

	// http.Handler
	var stdHandler http.Handler = dummyStdHandler{}
	assert.NotNil(t, MakeHandler(stdHandler))

	// net/http HandlerFunc style
	stdFn := func(w http.ResponseWriter, r *http.Request) {}
	assert.NotNil(t, MakeHandler(stdFn))

	// Our HandlerFunc type
	fn := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {}
	assert.NotNil(t, MakeHandler(fn))

	// Another, incompatible type
	assert.Panics(t, func() {
		MakeHandler(func(i int) int {
			return i + 1
		})
	})
}

// Test that wrapping a handler actually calls through to the underlying
// Handler, that it properly decodes parameters, and that the original body is
// readable and contains the same content
func TestWrapHandler(t *testing.T) {
	var w http.ResponseWriter = httptest.NewRecorder()

	const bodyString = `foo bar`
	b := bytes.NewBufferString(bodyString)
	r, err := http.NewRequest("GET", "/foo", b)
	assert.NoError(t, err)

	var run bool
	fn := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		run = true

		// Test body
		b, err := ioutil.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.Equal(t, bodyString, string(b))

		// Test parameters
		val, ok := ParamFrom(ctx, "param")
		assert.True(t, ok)
		assert.Equal(t, "foo", val)
	}

	a := New()
	a.Get("/:param", fn)
	a.ServeHTTP(w, r)
	assert.True(t, run)
}

// Test that we can add a http.Handler and it'll be run
func TestStdHandler(t *testing.T) {
	var w http.ResponseWriter = httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/foo", nil)
	assert.NoError(t, err)

	var run bool
	var h http.Handler
	h = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		run = true
	})

	a := New()
	a.Get("/:param", h)
	a.ServeHTTP(w, r)
	assert.True(t, run)
}

// Test that we can use our Handler type in the standard http.Handler location
func TestOurHandlerStd(t *testing.T) {
	var w http.ResponseWriter = httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/foo", nil)
	assert.NoError(t, err)

	var run bool
	h := HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		run = true
		assert.NotNil(t, ctx)
	})

	h.ServeHTTP(w, r)
	assert.True(t, run)
}
