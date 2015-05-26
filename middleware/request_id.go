package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	"golang.org/x/net/context"
)

var (
	requestIdKey private
	prefix       string
	reqid        uint64
)

func init() {
	hostname, err := os.Hostname()
	if hostname == "" || err != nil {
		hostname = "localhost"
	}

	var buf [12]byte
	var b64 string
	for len(b64) < 10 {
		rand.Read(buf[:])
		b64 = base64.StdEncoding.EncodeToString(buf[:])
		b64 = strings.NewReplacer("+", "", "/", "").Replace(b64)
	}

	prefix = fmt.Sprintf("%s/%s", hostname, b64[0:10])
}

// RequestID is a middleware that injects a request ID into the context of each
// request. A request ID is a string of the form "host.example.com/random-0001",
// where "random" is a base62 random string that uniquely identifies this go
// process, and where the last number is an atomically incremented request
// counter.
//
// Note: this middleware is borrowed from goji:
//	https://github.com/zenazn/goji/blob/master/web/middleware/request_id.go
func RequestID(ctx *context.Context, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctr := atomic.AddUint64(&reqid, 1)
		id := fmt.Sprintf("%s-%06d", prefix, ctr)

		*ctx = context.WithValue(*ctx, requestIdKey, id)
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// GetReqID returns a request ID from the given context if one is present.
// Returns the empty string if a request ID cannot be found.
func GetReqID(ctx context.Context) string {
	val := ctx.Value(requestIdKey)
	if val == nil {
		return ""
	}

	// This should never fail, since nobody else can use our key.
	return val.(string)
}
