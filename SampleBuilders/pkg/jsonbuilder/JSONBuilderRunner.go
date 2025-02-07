package jsonbuilder

import (
	"fmt"
	"log"
)

func JSONBuilderRunner() {
	// Instantiate our builder.
	var builder JSONBuilder = &DefaultJSONBuilder{}

	// ----- Incoming scenario -----
	// Imagine this JSON payload was received on a message.
	incomingPayload := `{
  "user": {
    "name": "Alice",
    "email": "alice@example.com"
  },
  "active": true,
  "roles": ["admin", "user"]
}`
	fmt.Println("Incoming JSON payload:")
	fmt.Println(incomingPayload)

	// Parse the incoming JSON.
	parsedData, err := builder.Parse([]byte(incomingPayload))
	if err != nil {
		log.Fatalf("Error parsing incoming JSON: %v", err)
	}
	fmt.Println("\nParsed internal representation:")
	fmt.Printf("%+v\n\n", parsedData)

	// ----- Outgoing scenario -----
	// Prepare an internal representation that we want to send out.
	outgoingData := map[string]interface{}{
		"status":  "success",
		"code":    200,
		"message": "Operation completed successfully",
		"data": map[string]interface{}{
			"id":    1001,
			"items": []string{"item1", "item2", "item3"},
		},
	}

	// Build (serialize) the outgoing JSON payload.
	outgoingPayload, err := builder.Build(outgoingData)
	if err != nil {
		log.Fatalf("Error building outgoing JSON: %v", err)
	}
	fmt.Println("Outgoing JSON payload:")
	fmt.Println(string(outgoingPayload))
}
