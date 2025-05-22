package main

import (
	"net/http"

	"github.com/rs/cors"
)

func main() {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://foo.com","http://example.com"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"X-Requested-With", "Content-Type"},
		ExposedHeaders: []string{"X-My-Custom-Header"},
		AllowCredentials: true,
		MaxAge: 3600,
		Debug: true,
	})
	mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte("{\"hello\": \"world\"}"))
    })

	handler := c.Handler(mux)
    http.ListenAndServe(":8080", handler)

	// handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.Write([]byte("{\"hello\": \"world\"}"))
	// })

	// http.ListenAndServe(":8080", c.Handler(handler))
}