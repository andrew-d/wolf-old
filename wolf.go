package wolf

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
)

// Internal private type for context key
type private struct{}

type Handler func(context.Context, http.ResponseWriter, *http.Request)

type App struct {
	router      *httprouter.Router
	stack       middlewareStack
	RootContext context.Context
}

func New() *App {
	ret := &App{
		router: httprouter.New(),
		stack: middlewareStack{
			funcs: make([]resolvedMiddlewareType, 0),
		},
	}
	ret.stack.app = ret
	ret.stack.resetPool()
	return ret
}

func (a *App) wrapHandler(h Handler) httprouter.Handle {
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
		h(ctx, w, r)
	}
}

func (a *App) Use(m interface{}) {
	a.stack.Push(m)
}

func (a *App) Handle(method, path string, handler Handler) {
	a.router.Handle(method, path, a.wrapHandler(handler))
}

func (a *App) Delete(path string, handler Handler) {
	a.router.DELETE(path, a.wrapHandler(handler))
}

func (a *App) Get(path string, handler Handler) {
	a.router.GET(path, a.wrapHandler(handler))
}

func (a *App) Head(path string, handler Handler) {
	a.router.HEAD(path, a.wrapHandler(handler))
}

func (a *App) Options(path string, handler Handler) {
	a.router.OPTIONS(path, a.wrapHandler(handler))
}

func (a *App) Patch(path string, handler Handler) {
	a.router.PATCH(path, a.wrapHandler(handler))
}

func (a *App) Post(path string, handler Handler) {
	a.router.POST(path, a.wrapHandler(handler))
}

func (a *App) Put(path string, handler Handler) {
	a.router.PUT(path, a.wrapHandler(handler))
}

func (a *App) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	m := a.stack.get()
	m.ServeHTTP(w, req)
	a.stack.release(m)
}
