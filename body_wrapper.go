package wolf

import (
	"io"

	"golang.org/x/net/context"
)

// Wrapper struct that lets us store our context within an incoming
// http.Request's Body field
type bodyWrapper struct {
	ctx        context.Context
	underlying io.ReadCloser
}

func (w *bodyWrapper) Read([]byte) (int, error) {
	panic("should not be called")
}

func (w *bodyWrapper) Close() error {
	panic("should not be called")
}

// Static type-checking
var _ io.ReadCloser = &bodyWrapper{}
