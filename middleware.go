package wolf

import (
	"net/http"
	"sync"

	"golang.org/x/net/context"
)

/*
General plans:
	- Middleware stack can be 'Push'-ed to
	- Pushing to a stack appends a new middleware func and invalidates any
		existing cache
	- We use a sync.Pool to cache the "applied" middleware stack (name?)
	- An "applied" middleware stack consists of an underlying handler function,
		with all middleware functions called on it
*/

type resolvedMiddlewareType func(ctx *context.Context, h http.Handler) http.Handler

// middlewareStack is an entire middleware stack.  It contains an array of
// middleware functions (outermost first) protected by a mutex, and a cache of
// pre-built stack instances.
type middlewareStack struct {
	funcs []resolvedMiddlewareType
	mu    sync.Mutex
	cache *sync.Pool // cache of pre-built middleware functions
	app   *App       // the app that this stack belongs to
}

func (m *middlewareStack) Push(fn interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Typecheck and append this function
	var resolvedFn resolvedMiddlewareType
	switch f := fn.(type) {
	case func(http.Handler) http.Handler:
		resolvedFn = func(ctx *context.Context, h http.Handler) http.Handler {
			return f(h)
		}
	case func(*context.Context, http.Handler) http.Handler:
		resolvedFn = f
	default:
		panic("TODO: real error message")
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
