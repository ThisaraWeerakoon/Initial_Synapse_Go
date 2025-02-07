package octetstreambuilder

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func OctetStreamBuilderRunner() {
	builder := &DefaultOctetStreamBuilder{}

	// Set up an HTTP handler that uses the builder.
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		// Parse the incoming octet-stream.
		data, err := builder.Parse(r)
		if err != nil {
			http.Error(w, "Error parsing octet-stream", http.StatusBadRequest)
			return
		}
		// For demonstration, simply echo back the received data.
		responseData, err := builder.Build(data)
		if err != nil {
			http.Error(w, "Error building response", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(responseData)
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
	outgoingPayload := "This is the outgoing raw binary payload."
	resp, err := http.Post("http://localhost:8080/upload", "application/octet-stream", strings.NewReader(outgoingPayload))
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

	// Block forever (or until interrupted) to keep the server running.
	// select {}
}
