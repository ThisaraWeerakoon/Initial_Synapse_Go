package formurlencodedbuilder

import (
	"fmt"
	"log"
	"net/url"
)

func FormUrlEncodedBuilderRunner() {
	builder := &DefaultFormDataBuilder{}

	// ----- Incoming Scenario -----
	fmt.Println("=== Incoming Scenario ===")
	req := simulateIncomingFormRequest()
	parsedValues, err := builder.Parse(req)
	if err != nil {
		log.Fatalf("Error parsing form data: %v", err)
	}
	fmt.Println("Parsed form values from request:")
	for key, vals := range parsedValues {
		fmt.Printf("  %s: %v\n", key, vals)
	}

	// ----- Outgoing Scenario -----
	fmt.Println("\n=== Outgoing Scenario ===")
	// Prepare data to be sent as application/x-www-form-urlencoded.
	outData := url.Values{}
	outData.Add("status", "success")
	outData.Add("message", "Operation completed successfully")
	outData.Add("id", "12345")
	encodedPayload, err := builder.Build(outData)
	if err != nil {
		log.Fatalf("Error building form data: %v", err)
	}
	// In a real HTTP response, you would set the Content-Type header to "application/x-www-form-urlencoded" and write encodedPayload.
	fmt.Println("Outgoing application/x-www-form-urlencoded payload:")
	fmt.Println(string(encodedPayload))
}
