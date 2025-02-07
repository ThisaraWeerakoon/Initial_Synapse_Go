package xmlbuilder

import (
	"encoding/xml"
)

// XMLBuilder defines the interface for custom XML builders.
type XMLBuilder interface {
	// Parse the incoming XML payload into an internal representation.
	Parse(input []byte) (interface{}, error)
	// Build converts an internal representation to an XML payload.
	Build(data interface{}) ([]byte, error)
}

// For demonstration, define a simple Order struct.
type Order struct {
	XMLName  xml.Name `xml:"order"`
	ID       string   `xml:"id"`
	Customer string   `xml:"customer"`
	Total    float64  `xml:"total"`
}

// DefaultXMLBuilder is a basic implementation using encoding/xml.
type DefaultXMLBuilder struct{}

// Parse unmarshals the XML into an Order.
func (b *DefaultXMLBuilder) Parse(input []byte) (interface{}, error) {
	var order Order
	err := xml.Unmarshal(input, &order)
	return order, err
}

// Build marshals the Order struct into an XML payload.
func (b *DefaultXMLBuilder) Build(data interface{}) ([]byte, error) {
	// Here we pretty-print the XML output.
	return xml.MarshalIndent(data, "", "  ")
}

