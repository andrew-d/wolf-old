package main

import (
	"fmt"
	"net/http"

	"github.com/andrew-d/wolf"
	"github.com/zenazn/goji/graceful"
	"golang.org/x/net/context"
)

func main() {
	m := wolf.New()

	m.Use(func(h http.Handler) http.Handler {
		fmt.Println("Basic middleware constructed")
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Basic middleware run")
			h.ServeHTTP(w, r)
		})
	})

	m.Use(func(ctx *context.Context, h http.Handler) http.Handler {
		fmt.Println("Full middleware constructed")
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Full middleware run")
			h.ServeHTTP(w, r)
		})
	})

	m.Get("/", func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		fmt.Println("Index page handler run")
		fmt.Fprintln(w, "Hello world")
	})

	m.Get("/:param", func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		s, _ := wolf.ParamFrom(ctx, "param")
		fmt.Printf("Param handler run with param = %s\n", s)
		fmt.Fprintf(w, "You gave me: %s\n", s)
	})

	graceful.ListenAndServe(":3001", m)
}
