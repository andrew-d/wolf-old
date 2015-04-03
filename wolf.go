package wolf

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
)

// Internal private type for context key
type private struct{}

// App is the base type for wolf.  It allows defining routes and adding
// middleware, and implements the http.Handler interface.
type App struct {
	router *httprouter.Router
	stack  middlewareStack

	// RootContext is the root context for this App.  Middleware functions'
	// context pointer defaults to pointing to this.
	RootContext context.Context
}

// New creates a new App with a background context.
func New() *App {
	ret := &App{
		router: httprouter.New(),
		stack: middlewareStack{
			funcs: make([]resolvedMiddlewareType, 0),
		},
		RootContext: context.Background(),
	}
	ret.stack.app = ret
	ret.stack.resetPool()
	return ret
}

// Use appends a middleware function to the set of middleware on this App.
func (a *App) Use(m interface{}) {
	a.stack.Push(m)
}

// Handle registers a new request handler with the given path and method.
//
// The app also provides shortcut methods for common HTTP methods (e.g. GET,
// POST, DELETE, etc.)
func (a *App) Handle(method, path string, handler HandlerType) {
	a.router.Handle(method, path, a.wrapHandler(handler))
}

// Delete is a shortcut for app.Handle("DELETE", path, handler)
func (a *App) Delete(path string, handler HandlerType) {
	a.Handle("DELETE", path, handler)
}

// Get is a shortcut for app.Handle("GET", path, handler)
func (a *App) Get(path string, handler HandlerType) {
	a.Handle("GET", path, handler)
}

// Head is a shortcut for app.Handle("HEAD", path, handler)
func (a *App) Head(path string, handler HandlerType) {
	a.Handle("HEAD", path, handler)
}

// Options is a shortcut for app.Handle("OPTIONS", path, handler)
func (a *App) Options(path string, handler HandlerType) {
	a.Handle("OPTIONS", path, handler)
}

// Patch is a shortcut for app.Handle("PATCH", path, handler)
func (a *App) Patch(path string, handler HandlerType) {
	a.Handle("PATCH", path, handler)
}

// Post is a shortcut for app.Handle("POST", path, handler)
func (a *App) Post(path string, handler HandlerType) {
	a.Handle("POST", path, handler)
}

// Put is a shortcut for app.Handle("PUT", path, handler)
func (a *App) Put(path string, handler HandlerType) {
	a.Handle("PUT", path, handler)
}

// ServeHTTP makes this App implement the http.Handler interface.
func (a *App) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	m := a.stack.get()
	m.ServeHTTP(w, req)
	a.stack.release(m)
}
