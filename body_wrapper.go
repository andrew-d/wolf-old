package wolf

import (
	"io"

	"golang.org/x/net/context"
)

// Wrapper struct that lets us store our context within an incoming
// http.Request's Body field.
type bodyWrapper struct {
	ctx        context.Context
	underlying io.ReadCloser
}

func (w *bodyWrapper) Read(buf []byte) (int, error) {
	return w.underlying.Read(buf)
}

func (w *bodyWrapper) Close() error {
	return w.underlying.Close()
}

// Static type-checking
var _ io.ReadCloser = &bodyWrapper{}
