package textplainbuilder

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func TextPlainBuilderRunner() {
	builder := &DefaultTextPlainBuilder{}

	// Set up an HTTP handler that uses the builder.
	http.HandleFunc("/text", func(w http.ResponseWriter, r *http.Request) {
		// Parse the incoming text/plain payload.
		text, err := builder.Parse(r)
		if err != nil {
			http.Error(w, "Error reading text payload", http.StatusBadRequest)
			return
		}
		fmt.Printf("Server parsed text: %s\n", text)

		// For demonstration, echo back the text with a prefix.
		responseText := "Server received: " + text
		out, err := builder.Build(responseText)
		if err != nil {
			http.Error(w, "Error building text response", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write(out)
	})

	// Start the HTTP server in a goroutine.
	go func() {
		fmt.Println("Starting HTTP server on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait a bit for the server to start.
	time.Sleep(2 * time.Second)

	// ----- Outgoing Client Scenario -----
	fmt.Println("=== Outgoing Client Scenario ===")
	// Simulate a client POST request with text/plain payload.
	// req := simulateIncomingTextRequest()
	resp, err := http.Post("http://localhost:8080/text", "text/plain", strings.NewReader("Hello from client!"))
	if err != nil {
		log.Fatalf("Error making client POST: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading client response: %v", err)
	}
	fmt.Println("Client received response:")
	fmt.Println(string(body))

	// Block indefinitely so the server keeps running.
	select {}
}
