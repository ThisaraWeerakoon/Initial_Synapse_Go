package binarybuilder

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func BinaryBuilderRunner() {
	builder := &DefaultBinaryBuilder{}

	// HTTP handler that uses the builder to parse and echo binary data.
	http.HandleFunc("/binary", func(w http.ResponseWriter, r *http.Request) {
		// Parse the incoming binary payload.
		data, err := builder.Parse(r)
		if err != nil {
			http.Error(w, "Error parsing binary data", http.StatusBadRequest)
			return
		}
		// For demonstration, echo back the same data.
		responseData, err := builder.Build(data)
		if err != nil {
			http.Error(w, "Error building binary response", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/binary")
		w.Write(responseData)
	})

	// Start the HTTP server in a goroutine.
	go func() {
		fmt.Println("Starting HTTP server on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait a couple of seconds to ensure the server is running.
	time.Sleep(2 * time.Second)

	// ----- Outgoing Client Scenario -----
	fmt.Println("=== Outgoing Client Scenario ===")
	// Prepare an outgoing binary payload.
	outgoingPayload := "This is the outgoing application/binary payload."
	resp, err := http.Post("http://localhost:8080/binary", "application/binary", strings.NewReader(outgoingPayload))
	if err != nil {
		log.Fatalf("Error making client POST: %v", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading client response: %v", err)
	}
	resp.Body.Close()
	fmt.Println("Client received response:")
	fmt.Println(string(body))

	// Block forever to keep the server running.
	select {}
}