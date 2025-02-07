package multipartbuilder

import (
	"fmt"
	"log"
	"os"
)

func MultipartBuilderRunner() {
	// Instantiate our builder.
	var builder MultipartBuilder = &DefaultMultipartBuilder{}

	// ----- Incoming Scenario -----
	fmt.Println("=== Incoming Scenario ===")
	req := simulateIncomingMultipartRequest()
	parsedData, err := builder.Parse(req)
	if err != nil {
		log.Fatalf("Error parsing multipart: %v", err)
	}
	fmt.Println("Parsed Fields:")
	for key, vals := range parsedData.Fields {
		fmt.Printf("  %s: %v\n", key, vals)
	}
	fmt.Println("Parsed Files:")
	for key, parts := range parsedData.Files {
		for _, part := range parts {
			fmt.Printf("  %s: Filename=%s, Size=%d bytes\n", key, part.Filename, len(part.Content))
			fmt.Printf("     Content: %s\n", string(part.Content))
		}
	}

	// ----- Outgoing Scenario -----
	fmt.Println("\n=== Outgoing Scenario ===")
	// Prepare outgoing data.
	outData := &MultipartData{
		Fields: map[string][]string{
			"action": {"create"},
			"token":  {"abc123"},
		},
		Files: make(map[string][]FilePart),
	}
	// Create a temporary file to simulate a file to upload.
	tempFile, err := os.CreateTemp("", "sample")
	if err != nil {
		log.Fatalf("Error creating temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	sampleContent := "Sample file content for outgoing multipart."
	tempFile.WriteString(sampleContent)
	tempFile.Close()

	// Read the temporary file content.
	content, err := os.ReadFile(tempFile.Name())
	if err != nil {
		log.Fatalf("Error reading temp file: %v", err)
	}
	// Add the file part to the outgoing data.
	outData.Files["upload"] = append(outData.Files["upload"], FilePart{
		Filename: "sample.txt",
		Content:  content,
	})

	// Build the outgoing multipart/form-data payload.
	ct, body, err := builder.Build(outData)
	if err != nil {
		log.Fatalf("Error building multipart: %v", err)
	}
	fmt.Println("Outgoing multipart/form-data Content-Type:")
	fmt.Println(ct)
	fmt.Println("Outgoing multipart/form-data Body:")
	fmt.Println(string(body))
}
