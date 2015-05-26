package middleware

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"golang.org/x/net/context"
)

// RecoverInformation contains information about a recovered panic.
type RecoverInformation struct {
	// Error is the value that was passed to panic()
	Error interface{}

	// Stack is the call stack as of the call to recover().
	Stack []byte

	// RequestID is the ID of this request.  It may be empty if the corresponding
	// middleware was not included.
	RequestID string
}

// RecoverFunc is the function type for a callback that can handle recovers.
type RecoverFunc func(w http.ResponseWriter, r *http.Request, info RecoverInformation)

// Recoverer is a middleware that recovers from panics, prints the panic (and
// a backtrace), and then returns a HTTP 500 (Internal Server Error) status to
// the client, if possible.
//
// Recoverer will also include the request ID if one is provided.
func Recoverer(ctx *context.Context, h http.Handler) http.Handler {
	return CustomRecoverer(defaultRecoverFunc)(ctx, h)
}

// CustomRecoverer creates a middleware that recovers from panics, as with
// Recoverer, but passes information about the panic to a user-defined function
// that can take whatever action is necessary.
func CustomRecoverer(cb RecoverFunc) func(*context.Context, http.Handler) http.Handler {
	middlewareFunc := func(ctx *context.Context, h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			requestId := GetReqID(*ctx)

			defer func() {
				if err := recover(); err != nil {
					info := RecoverInformation{
						Error:     err,
						Stack:     debug.Stack(),
						RequestID: requestId,
					}
					cb(w, r, info)
				}
			}()

			h.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}

	return middlewareFunc
}

// The default RecoverFunc just prints to the screen
func defaultRecoverFunc(w http.ResponseWriter, r *http.Request, info RecoverInformation) {
	var buf bytes.Buffer

	if info.RequestID != "" {
		fmt.Fprintf(&buf, "[%s] ", info.RequestID)
	}
	fmt.Fprintf(&buf, "panic: %+v", info.Error)

	// Print the error to the screen
	log.Printf(buf.String())

	// Print the stack.
	os.Stderr.Write(info.Stack)

	// Return an error
	http.Error(w, http.StatusText(500), 500)
}
