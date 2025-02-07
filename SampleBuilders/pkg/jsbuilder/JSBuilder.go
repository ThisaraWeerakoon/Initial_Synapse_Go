package jsbuilder

import (
	"io"
	"log"
	"net/http"
	"strings"
)

// JSBuilder defines methods for handling text/javascript content.
type JSBuilder interface {
	// Parse reads the HTTP request body and returns the JavaScript code as a string.
	Parse(r *http.Request) (string, error)
	// Build creates a payload from given JavaScript code.
	Build(js string) ([]byte, error)
}

// DefaultJSBuilder is a basic implementation using standard library functions.
type DefaultJSBuilder struct{}

// Parse reads the entire request body as a string.
func (b *DefaultJSBuilder) Parse(r *http.Request) (string, error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(bodyBytes), nil
}

// Build converts the given JavaScript code into a byte slice.
func (b *DefaultJSBuilder) Build(js string) ([]byte, error) {
	return []byte(js), nil
}

// simulateIncomingJSRequest constructs a dummy HTTP request with text/javascript content.
func simulateIncomingJSRequest() *http.Request {
	jsCode := `function hello() { console.log("Hello, world!"); }`
	req, err := http.NewRequest("POST", "http://example.com/js", strings.NewReader(jsCode))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "text/javascript")
	return req
}


