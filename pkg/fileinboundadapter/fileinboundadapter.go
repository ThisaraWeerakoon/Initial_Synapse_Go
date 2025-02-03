package fileinboundadapter

import (
	"time"
	"os"
	"fmt"
	"log"
	"path/filepath"
	"io"
	"bufio"
	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/pkg/models"
)

// MoveFile moves a file from source to destination.
func MoveFile(source, destination string) error {
	return os.Rename(source, destination)
}

// PollFolder continuously reads files at a given interval
func PollFolder(folderPath string, interval time.Duration) {
	ticker := time.NewTicker(interval) // Ensures precise polling
	defer ticker.Stop()

	for range ticker.C {
		startTime := time.Now()
		fmt.Println("\n--- New Polling Event ---")

		// Process failed_files.txt if available
		processFailedFiles(folderPath)

		// Get the list of files at the start of this polling event
		files, err := scanDirectory(folderPath)
		if err != nil {
			log.Printf("Error scanning directory %s: %v", folderPath, err)
			continue
		}

		// Process each file
		for _, file := range files {
			metadata, content, err := ReadFile(file)
			if err != nil {
				log.Printf("Error reading file %s: %v", file, err)
				continue
			}
			fmt.Printf("Processed: %s\nMetadata: %+v\nContent:\n%s\n", file, metadata, content)
		}

		// Ensure accurate polling interval
		elapsed := time.Since(startTime)
		if elapsed < interval {
			time.Sleep(interval - elapsed)
		}
	}
}


// FileMetadata stores metadata of the file
type FileMetadata struct {
	Name     string // File name
	Size     int64  // File size in bytes
	FileType string // File extension
}

// ReadFile reads a file and returns metadata & content
func ReadFile(filePath string) (*FileMetadata, string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file %s: %v", filePath, err)
		return nil, "", err
	}
	defer file.Close()

	// Get file metadata
	info, err := file.Stat()
	if err != nil {
		log.Printf("Error getting file info for %s: %v", filePath, err)
		return nil, "", err
	}

	metadata := &FileMetadata{
		Name:     info.Name(),
		Size:     info.Size(),
		FileType: filepath.Ext(filePath), // Get file extension
	}

	// Read file content
	content, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Error reading file %s: %v", filePath, err)
		return metadata, "", err
	}

	return metadata, string(content), nil
}

// Moves failed files from `test/in/` to `test/failed/`
func processFailedFiles(failedFilePath string) {
	baseDir := "test" // Root directory for test folders
	inDir := filepath.Join(baseDir, "in")
	failedDir := filepath.Join(baseDir, "failed")

	// Open failed_files.txt if it exists
	file, err := os.Open(failedFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No failed_files.txt found, skipping.")
			return
		}
		log.Printf("Error opening %s: %v", failedFilePath, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		failedFile := scanner.Text()
		sourcePath := filepath.Join(inDir, failedFile)
		destPath := filepath.Join(failedDir, failedFile)

		// Move the failed file to the failed directory
		err := os.Rename(sourcePath, destPath)
		if err != nil {
			log.Printf("Error moving %s to failed folder: %v", sourcePath, err)
		} else {
			fmt.Printf("Moved %s to failed folder\n", failedFile)
		}
	}

	// Handle scanning errors
	if err := scanner.Err(); err != nil {
		log.Printf("Error reading failed_files.txt: %v", err)
	}

	// Remove failed_files.txt after processing
	err = os.Remove(failedFilePath)
	if err != nil {
		log.Printf("Error deleting failed_files.txt: %v", err)
	}
}

// Scan a directory and return the list of files present at the start of polling
func scanDirectory(folderPath string) ([]string, error) {
	var files []string

	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && entry.Name() != "failed_files.txt" { // Ignore directories and failed_files.txt
			files = append(files, filepath.Join(folderPath, entry.Name()))
		}
	}

	return files, nil
}


type CoreInterface interface {
	//ReceiveRequests
}

type FileInboundAdapter struct{
	models.Configurations

}

func NewFileInboundAdapter(config models.Configurations) *FileInboundAdapter {
	return &FileInboundAdapter{
		Configurations: config,
	}
}

func (f  *FileInboundAdapter) StartPolling() {
	//start polling
	PollFolder(f.FileURI, time.Duration(f.Interval)*time.Second)

}

func (f  *FileInboundAdapter) ReceiveResults() {
	//reveive results
}

func (f  *FileInboundAdapter) Start() {
	//start process
	PollFolder(f.FileURI, time.Duration(f.Interval)*time.Second)

}

func (f  *FileInboundAdapter) Stop() {	
	//stop process
}