package formurlencodedbuilder

import (
	"log"
	"net/http"
	"net/url"
	"strings"
)

// FormDataBuilder defines methods for parsing and building form-urlencoded data.
type FormDataBuilder interface {
	// Parse extracts form values from the given HTTP request.
	Parse(r *http.Request) (url.Values, error)
	// Build creates an application/x-www-form-urlencoded payload from provided data.
	Build(data url.Values) ([]byte, error)
}

// DefaultFormDataBuilder is a simple implementation that uses Go's standard library.
type DefaultFormDataBuilder struct{}

// Parse calls r.ParseForm() and returns the parsed form values.
func (b *DefaultFormDataBuilder) Parse(r *http.Request) (url.Values, error) {
	// Ensure the form data is parsed. This parses both URL query parameters and POST form data.
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	return r.Form, nil
}

// Build encodes the provided url.Values into a form-urlencoded byte slice.
func (b *DefaultFormDataBuilder) Build(data url.Values) ([]byte, error) {
	encoded := data.Encode()
	return []byte(encoded), nil
}

// simulateIncomingFormRequest creates a dummy HTTP POST request with an application/x-www-form-urlencoded body.
func simulateIncomingFormRequest() *http.Request {
	// Create sample form data.
	form := url.Values{}
	form.Add("username", "john")
	form.Add("email", "john@example.com")
	form.Add("roles", "admin")
	form.Add("roles", "user")
	encodedForm := form.Encode()

	// Create an HTTP request with the form payload.
	req, err := http.NewRequest("POST", "http://example.com/form", strings.NewReader(encodedForm))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

