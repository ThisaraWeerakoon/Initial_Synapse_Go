package main

import (
	"log"
	"net/http"
	"time"
)

// loggingMiddleware is our custom middleware.
// It takes an http.Handler (the next handler in the chain)
// and returns a new http.Handler.
func loggingMiddleware(next http.Handler) http.Handler {
	// The returned http.Handler is an http.HandlerFunc,
	// which is a type that allows us to use an ordinary function
	// as an HTTP handler.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Our middleware logic starts here
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		// Call the next handler in the chain. This could be another
		// middleware or our actual application handler.
		next.ServeHTTP(w, r)

		// Our middleware logic continues after the next handler has finished
		log.Printf("Completed %s in %v", r.URL.Path, time.Since(start))
	})
}

// helloHandler is our actual application handler.
func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Hello, World!"))
}

// aboutHandler is another application handler.
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is the about page."))
}

func main() {
	// Our main application handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/about", aboutHandler)

	// Wrap our main router (mux) with the loggingMiddleware.
	// All requests going to mux will now pass through loggingMiddleware first.
	loggedMux := loggingMiddleware(mux)

	log.Println("Server starting on :8080...")
	// Start the server with our middleware-wrapped handler
	if err := http.ListenAndServe(":8080", loggedMux); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}