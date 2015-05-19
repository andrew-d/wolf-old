package wolf

import (
	"net/http"
	"testing"

	"golang.org/x/net/context"
)

// Verifies that we can add handlers with our various helper functions
func TestWolfHelpers(t *testing.T) {
	a := New()

	fn := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {}

	a.Delete("/", fn)
	a.Get("/", fn)
	a.Head("/", fn)
	a.Options("/", fn)
	a.Patch("/", fn)
	a.Post("/", fn)
	a.Put("/", fn)

	a.Compile()
}
