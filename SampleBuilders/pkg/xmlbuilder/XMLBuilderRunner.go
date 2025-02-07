package xmlbuilder

import (
	"fmt"
	"log"
)

func XMLBuilderRunner() {
	// Instantiate our builder.
	var builder XMLBuilder = &DefaultXMLBuilder{}

	// ----- Incoming scenario -----
	// Imagine this XML payload is received by your MI.
	incomingXML := `<order>
  <id>9876</id>
  <customer>John Doe</customer>
  <total>123.45</total>
</order>`
	fmt.Println("Incoming XML payload:")
	fmt.Println(incomingXML)

	// Parse the incoming XML.
	parsedData, err := builder.Parse([]byte(incomingXML))
	if err != nil {
		log.Fatalf("Error parsing XML: %v", err)
	}
	fmt.Println("\nParsed internal representation:")
	fmt.Printf("%+v\n\n", parsedData)

	// ----- Outgoing scenario -----
	// Create an internal Order object to send out.
	outgoingOrder := Order{
		ID:       "54321",
		Customer: "Alice Smith",
		Total:    987.65,
	}
	xmlPayload, err := builder.Build(outgoingOrder)
	if err != nil {
		log.Fatalf("Error building XML: %v", err)
	}
	fmt.Println("Outgoing XML payload:")
	fmt.Println(string(xmlPayload))
}