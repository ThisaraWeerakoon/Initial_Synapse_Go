package main

import (
        "fmt"
        "net/http"
)

func main() {
        mux := http.NewServeMux()
        
        mux.HandleFunc("/product/{id}/orders/{orderId}", func(w http.ResponseWriter, r *http.Request) {

				// Accessing the path parameter (e.g., "id")
                productId := r.PathValue("id")
				orderId := r.PathValue("orderId")
                queryParams := r.URL.Query()

                // Accessing a specific query parameter (e.g., "color")
                color := queryParams.Get("color")

                fmt.Fprintf(w, "displaying properties for product %s order %s\n", productId, orderId)
                if color != "" {
                        fmt.Fprintf(w, "Color: %s\n", color)
                }

                // Accessing all query parameters
                for key, values := range queryParams {
                        fmt.Printf("Query Parameter: %s = %v\n", key, values)
                }
        })
        fmt.Println("Server is listening on port 8000")
        http.ListenAndServe(":8000", mux)
}