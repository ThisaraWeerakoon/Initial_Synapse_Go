package jsbuilder

import (
	"fmt"
	"log"
)

func JSONBuilderRunner() {
	builder := &DefaultJSBuilder{}

	// ----- Incoming Scenario -----
	fmt.Println("=== Incoming Scenario ===")
	req := simulateIncomingJSRequest()
	parsedJS, err := builder.Parse(req)
	if err != nil {
		log.Fatalf("Error parsing text/javascript: %v", err)
	}
	fmt.Println("Parsed JavaScript code:")
	fmt.Println(parsedJS)

	// ----- Outgoing Scenario -----
	fmt.Println("\n=== Outgoing Scenario ===")
	// Prepare a JavaScript payload to send.
	jsToSend := `var greeting = "Hello, world!"; alert(greeting);`
	payload, err := builder.Build(jsToSend)
	if err != nil {
		log.Fatalf("Error building text/javascript payload: %v", err)
	}
	// In a real HTTP response you would set the Content-Type header to "text/javascript".
	fmt.Println("Outgoing text/javascript payload:")
	fmt.Println(string(payload))
}