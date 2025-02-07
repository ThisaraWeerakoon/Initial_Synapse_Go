package textplainbuilder

import (
	"io"
	"log"
	"net/http"
	"strings"
)

// TextPlainBuilder defines methods for handling text/plain content.
type TextPlainBuilder interface {
	// Parse reads the request body and returns its content as a string.
	Parse(r *http.Request) (string, error)
	// Build converts a given string into a byte slice.
	Build(text string) ([]byte, error)
}

// DefaultTextPlainBuilder implements TextPlainBuilder using standard library functions.
type DefaultTextPlainBuilder struct{}

// Parse uses io.ReadAll to read the entire request body.
func (b *DefaultTextPlainBuilder) Parse(r *http.Request) (string, error) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Build converts the provided string into a byte slice.
func (b *DefaultTextPlainBuilder) Build(text string) ([]byte, error) {
	return []byte(text), nil
}

// simulateIncomingTextRequest constructs a dummy HTTP POST request with text/plain content.
func simulateIncomingTextRequest() *http.Request {
	payload := "This is a sample text payload from the client."
	req, err := http.NewRequest("POST", "http://localhost:8080/text", strings.NewReader(payload))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "text/plain")
	return req
}
