package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/andrew-d/wolf"
	"github.com/andrew-d/wolf/middleware"
	"github.com/zenazn/goji/graceful"
	"golang.org/x/net/context"
)

func main() {
	m := wolf.New()

	m.Use(middleware.RequestID)

	m.Use(func(h http.Handler) http.Handler {
		log.Println("Basic middleware constructed")
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("Basic middleware run")
			h.ServeHTTP(w, r)
		})
	})

	m.Use(func(ctx *context.Context, h http.Handler) http.Handler {
		log.Println("Full middleware constructed")
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("Full middleware run")
			h.ServeHTTP(w, r)
		})
	})

	m.Get("/", func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		log.Println("Index page handler run")
		fmt.Fprintf(w, "Hello world, our request ID is: %s\n", middleware.GetReqID(ctx))
	})

	m.Get("/:param", func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		s, _ := wolf.ParamFrom(ctx, "param")
		log.Printf("Param handler run with param = %s\n", s)
		fmt.Fprintf(w, "You gave me: %s\n", s)
	})

	log.Println("Started")
	graceful.ListenAndServe(":3001", m)
}
