package binarybuilder

import (
	"io"
	"log"
	"net/http"
	"strings"
)

// BinaryBuilder defines methods for handling application/binary payloads.
type BinaryBuilder interface {
	// Parse reads the raw bytes from the HTTP request body.
	Parse(r *http.Request) ([]byte, error)
	// Build returns the raw byte payload for an outgoing message.
	Build(data []byte) ([]byte, error)
}

// DefaultBinaryBuilder is a simple implementation that uses the standard library.
type DefaultBinaryBuilder struct{}

// Parse uses io.ReadAll to read the entire request body.
func (b *DefaultBinaryBuilder) Parse(r *http.Request) ([]byte, error) {
	return io.ReadAll(r.Body)
}

// Build simply returns the provided data without modification.
func (b *DefaultBinaryBuilder) Build(data []byte) ([]byte, error) {
	return data, nil
}

// simulateIncomingBinaryRequest creates a dummy HTTP POST request with an application/binary body.
func simulateIncomingBinaryRequest() *http.Request {
	payload := "This is raw binary data for application/binary."
	req, err := http.NewRequest("POST", "http://localhost:8080/binary", strings.NewReader(payload))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/binary")
	return req
}


