package main

import (
	"log"
	"net/http"
	"time"
)

// loggingMiddleware logs method, path, and duration
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// anotherMiddleware
func otherMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("other middleware activated")
		next.ServeHTTP(w, r)
	})
}

func main() {
	// create a mux
	mux := http.NewServeMux()

	// a regular path
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// the other path, our middlewares expect an http.handler interface, so we wrap our unnamed function with http.HandlerFunc
	mux.Handle("GET /other", otherMiddleware(http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello Other World!"))
	})))

	// attach the loggingMiddleware to every request by wrapping our router (mux) inside of it
	log.Println("Launching server on port 3000")
	log.Fatal(http.ListenAndServe(":3000", loggingMiddleware(mux)))
}
