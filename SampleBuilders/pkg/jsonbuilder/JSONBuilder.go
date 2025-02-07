package jsonbuilder

import (
	"encoding/json"
)

// JSONBuilder defines the interface for custom JSON builders.
type JSONBuilder interface {
	// Parse the incoming JSON payload and return an internal representation.
	Parse(input []byte) (interface{}, error)
	// Build converts the internal representation to an outgoing JSON payload.
	Build(data interface{}) ([]byte, error)
}

// DefaultJSONBuilder is a simple implementation that uses encoding/json.
type DefaultJSONBuilder struct{}

// Parse implements JSONBuilder.Parse by unmarshaling the JSON.
func (b *DefaultJSONBuilder) Parse(input []byte) (interface{}, error) {
	var result interface{}
	err := json.Unmarshal(input, &result)
	return result, err
}

// Build implements JSONBuilder.Build by marshaling the object to JSON.
func (b *DefaultJSONBuilder) Build(data interface{}) ([]byte, error) {
	// For pretty output, we can use MarshalIndent.
	return json.MarshalIndent(data, "", "  ")
}

