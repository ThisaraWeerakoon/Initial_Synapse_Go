package octetstreambuilder

import (
	"io"
	"log"
	"net/http"
	"strings"
)

// OctetStreamBuilder defines methods to handle application/octet-stream payloads.
type OctetStreamBuilder interface {
	// Parse reads the raw bytes from the HTTP request body.
	Parse(r *http.Request) ([]byte, error)
	// Build prepares the outgoing payload from the provided raw bytes.
	Build(data []byte) ([]byte, error)
}

// DefaultOctetStreamBuilder implements the OctetStreamBuilder interface.
type DefaultOctetStreamBuilder struct{}

// Parse uses io.ReadAll to read the entire request body.
func (b *DefaultOctetStreamBuilder) Parse(r *http.Request) ([]byte, error) {
	return io.ReadAll(r.Body)
}

// Build simply returns the provided data.
func (b *DefaultOctetStreamBuilder) Build(data []byte) ([]byte, error) {
	return data, nil
}

// simulateIncomingRequest constructs a dummy HTTP POST request with an octet-stream body.
func simulateIncomingRequest() *http.Request {
	payload := "This is raw binary data for incoming request."
	req, err := http.NewRequest("POST", "http://localhost:8080/upload", strings.NewReader(payload))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	return req
}
