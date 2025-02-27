package main 

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /product/{id}", func(w http.ResponseWriter, r *http.Request) {
    	productId := r.PathValue("id")
    	fmt.Fprintf(w, "displaying properties for product %s", productId)
	})

	mux.HandleFunc("POST /product/{id}", func(w http.ResponseWriter, r *http.Request) {
    	productId := r.PathValue("id")
    	fmt.Fprintf(w, "creating product with id %s", productId)
	})

	http.ListenAndServe(":8000", mux)
}
