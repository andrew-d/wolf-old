package wolf

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
)

// HandlerType is an alias for interface{}, but is documented here for clarity.
// wolf will accept a handler of one of the following types, and will convert
// it to the Handler interface that is used internally.
//
//	- types that implement http.Handler
//	- types that implement Handler
//	- func(http.ResponseWriter, *http.Request)
//	- func(context.Context, http.ResponseWriter, *http.Request)
type HandlerType interface{}

// Handler is similar to net/http's http.Handler, but accepts a Context from
// x/net/context as the first parameter.
type Handler interface {
	ServeHTTPCtx(context.Context, http.ResponseWriter, *http.Request)
}

// HandlerFunc is similar to net/http's http.HandlerFunc, but accepts a Context
// object.  It implements both the Handler interface and http.Handler.
type HandlerFunc func(context.Context, http.ResponseWriter, *http.Request)

// ServeHTTP implements http.Handler, allowing HandlerFuncs to be used with
// net/http and other routers.  When used this way, the underlying function
// will be passed a Background context.
func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(context.Background(), w, r)
}

// ServeHTTPCtx implements Handler.
func (f HandlerFunc) ServeHTTPCtx(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	f(ctx, w, r)
}

// netHTTPWrap is a helper to turn a http.Handler into our Handler
type netHTTPWrap struct {
	http.Handler
}

func (h netHTTPWrap) ServeHTTPCtx(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	h.ServeHTTP(w, r)
}

// MakeHandler turns a HandlerType into something that implements our Handler
// interface.  It will panic if the input is not a valid HandlerType.
func MakeHandler(h HandlerType) Handler {
	// Convert the handler type
	switch f := h.(type) {
	case Handler:
		return f
	case http.Handler:
		return netHTTPWrap{f}
	case func(context.Context, http.ResponseWriter, *http.Request):
		return HandlerFunc(f)
	case func(http.ResponseWriter, *http.Request):
		return netHTTPWrap{http.HandlerFunc(f)}
	default:
		panic("") // TODO
	}
}

// wrapHandler turns something that implements our Handler interface into a
// function that implements httprouter's interface.
func (a *App) wrapHandler(v HandlerType) httprouter.Handle {
	h := MakeHandler(v)
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var ctx context.Context

		// Get context that was modified by the middleware
		wrapper := r.Body.(*bodyWrapper)
		ctx = wrapper.ctx
		r.Body = wrapper.underlying

		// Unpack the request params
		ctx = setParamsInContext(ctx, p)

		// TODO: do we want to save w&r in the context?

		// Call the underlying handler
		h.ServeHTTPCtx(ctx, w, r)
	}
}
