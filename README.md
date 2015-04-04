wolf
====

[![GoDoc](https://godoc.org/github.com/andrew-d/wolf?status.svg)](https://godoc.org/github.com/andrew-d/wolf) [![Build Status](https://travis-ci.org/andrew-d/wolf.svg?branch=master)](https://travis-ci.org/andrew-d/wolf) [![Coverage Status](https://coveralls.io/repos/andrew-d/wolf/badge.svg?branch=master)](https://coveralls.io/r/andrew-d/wolf?branch=master)

A very simple (~200 LoC) web "framework" that simply integrates [httprouter][hr] and [x/net/context][ctx].

## Example

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/andrew-d/wolf"
	"golang.org/x/net/context"
)

func main() {
	m := wolf.New()

	m.Use(myMiddleware)

	m.Get("/", func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello world")
	})

	log.Println("Started")
	http.ListenAndServe(":3001", m)
}

func myMiddleware(ctx *context.Context, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got request to: %s", r.URL)
		h.ServeHTTP(w, r)
	})
}
```


## License

MIT


[hr]: https://github.com/julienschmidt/httprouter
[ctx]: https://godoc.org/golang.org/x/net/context
