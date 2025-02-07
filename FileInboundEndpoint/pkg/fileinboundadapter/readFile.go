package fileinboundadapter

import(
	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/FileInboundEndpoint/models"
	"io"
	"log"
	"os"
	"syscall"
	"net/url"
	"path/filepath"
	"strings"
	"fmt"
)

// ConvertFileURIToPath converts a file:// URI to an absolute file path.
func ConvertFileURIToPath(fileURI string) (string, error) {
	// Parse the file URI
	parsedURI, err := url.Parse(fileURI)
	if err != nil {
		return "", fmt.Errorf("invalid file URI: %v", err)
	}

	// Ensure scheme is "file"
	if parsedURI.Scheme != "file" {
		return "", fmt.Errorf("unsupported URI scheme: %s", parsedURI.Scheme)
	}

	// Get the file path and decode any URL encoding (e.g., spaces as `%20`)
	filePath := parsedURI.Path
	filePath = filepath.Clean(filePath)
	filePath = strings.ReplaceAll(filePath, "%20", " ") // Handle spaces

	return filePath, nil
}

// ReadFile reads a file, extracts metadata, locks it, and returns the extracted data.
func ReadFile(fileURI string) (*models.ExtractedFileDataFromFileAdapter, error) {
	// Convert file URI to absolute file path
	filePath, err := ConvertFileURIToPath(fileURI)
	if err != nil {
		log.Printf("Error converting file URI to path: %v", err)
		return nil, err
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file %s: %v", filePath, err)
		return nil, err
	}
	defer file.Close()

	// Lock the file to prevent modifications
	err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX) // Exclusive lock
	if err != nil {
		log.Printf("Error locking file %s: %v", filePath, err)
		return nil, err
	}

	// Get file metadata
	info, err := file.Stat()
	if err != nil {
		log.Printf("Error getting file info for %s: %v", filePath, err)
		return nil, err
	}

	// Extract metadata into ContextHeader
	header := models.ContextHeader{
		FILE_LENGTH:   float64(info.Size()),
		LAST_MODIFIED: float64(info.ModTime().Unix()), // Convert to Unix timestamp
		FILE_URI:      fileURI,                        // Keep original FILE_URI
		FILE_PATH:     filePath,                       // Derived file path
		FILE_NAME:     info.Name(),
	}

	// Read file content
	content, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Error reading file %s: %v", filePath, err)
		return &models.ExtractedFileDataFromFileAdapter{ContextHeader: header, Context: ""}, err
	}

	return &models.ExtractedFileDataFromFileAdapter{ContextHeader: header, Context: string(content)}, nil
}

