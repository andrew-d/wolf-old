package wolf

import (
	"fmt"
	"net/http"
	"sync"

	"golang.org/x/net/context"
)

// MiddlewareType is an alias for interface{}, but is documented here for
// clarity.  wolf will accept middleware of one of the following types, and
// will convert it to the internal middleware type.
//
//	- func(*context.Context, http.Handler) http.Handler
//	- func(http.Handler) http.Handler
type MiddlewareType interface{}

type canonicalMiddleware func(ctx *context.Context, h http.Handler) http.Handler

// middlewareStack is an entire middleware stack.  It contains an array of
// middleware functions (outermost first) protected by a mutex, and a cache of
// pre-built stack instances.
type middlewareStack struct {
	funcs []canonicalMiddleware
	mu    sync.Mutex
	cache *sync.Pool // cache of pre-built middleware functions
	app   *App       // the app that this stack belongs to
}

func (m *middlewareStack) Push(fn MiddlewareType) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Typecheck and append this function
	var resolvedFn canonicalMiddleware
	switch f := fn.(type) {
	case func(http.Handler) http.Handler:
		resolvedFn = func(ctx *context.Context, h http.Handler) http.Handler {
			return f(h)
		}
	case func(*context.Context, http.Handler) http.Handler:
		resolvedFn = f
	default:
		msg := fmt.Sprintf(`Invalid middleware type '%T'.  See `+
			`https://godoc.org/github.com/andrew-d/wolf#MiddlewareType for a `+
			`list of valid middleware types`, fn)
		panic(msg)
	}

	m.funcs = append(m.funcs, resolvedFn)

	// Invalidate the existing cache
	m.resetPool()
}

func (m *middlewareStack) resetPool() {
	m.cache = &sync.Pool{
		New: m.newResolved,
	}
}

func (m *middlewareStack) get() http.Handler {
	return m.cache.Get().(http.Handler)
}

func (m *middlewareStack) release(h http.Handler) {
	m.cache.Put(h)
}

// Apply all middleware funcs to our final routing function
func (m *middlewareStack) newResolved() interface{} {
	ctx := m.app.RootContext
	wrapper := bodyWrapper{}

	// This is the final routing function - it just dispatches to our router
	var finalFunc http.Handler
	finalFunc = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Save the context
		wrapper.ctx = ctx
		wrapper.underlying = r.Body
		r.Body = &wrapper

		// Dispatch to router
		m.app.router.ServeHTTP(w, r)
	})

	// Apply middleware
	for i := len(m.funcs) - 1; i >= 0; i-- {
		finalFunc = m.funcs[i](&ctx, finalFunc)
	}

	return finalFunc
}
