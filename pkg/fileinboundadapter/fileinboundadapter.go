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
func PollFolder(inDir string, outDir string, failedDir string, interval int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second) // Ensures precise polling
	defer ticker.Stop()

	for range ticker.C {
		startTime := time.Now()
		fmt.Println("\n--- New Polling Event ---")

		// Process failed_files.txt if available
		processFailedFiles(inDir, failedDir)

		// Get the list of files at the start of this polling event
		files, err := scanDirectory(inDir)
		if err != nil {
			log.Printf("Error scanning directory %s: %v", inDir, err)
			continue
		}

		// Process each file
		for _, file := range files {
			//have to test is it safe to make go routines for each file arbitrarily
			// A solution may be put a threshold (eg:- 100 files) and then make go routines for each file.If the number of files is greater than the threshold make only upper limit (threshold) of go routines
			go ProcessFile(file)
		}

		// Ensure accurate polling interval
		elapsed := time.Since(startTime)
		if elapsed < time.Duration(interval)*time.Second {
			time.Sleep(time.Duration(interval)*time.Second - elapsed)
		}
	}
}

func ProcessFile(file string) {
	metadata, err := ReadFile(file)
	if err != nil {
		log.Printf("Error reading file %s: %v", file, err)
		return
	}
	fmt.Printf("Processed: %s\nMetadata: %+v\nContent:\n%s\n", file, metadata, metadata.Context)

}


// ReadFile reads a file and returns metadata & content
func ReadFile(filePath string) (*models.ExtractedFileData, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file %s: %v", filePath, err)
		return nil, err
	}
	defer file.Close()

	// Get file metadata
	info, err := file.Stat()
	if err != nil {
		log.Printf("Error getting file info for %s: %v", filePath, err)
		return nil, err
	}

	metadata := &models.FileMetadata{
		Name:     info.Name(),
		Size:     info.Size(),
		FileType: filepath.Ext(filePath), // Get file extension
	}

	// Read file content
	content, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Error reading file %s: %v", filePath, err)
		//return metadata, "", err
		return &models.ExtractedFileData{FileMetadata: *metadata, Context: ""}, err
	}

	return &models.ExtractedFileData{FileMetadata: *metadata, Context: string(content)}, nil
}

// Moves failed files from `test/in/` to `test/failed/`. test/failed/failed_files.txt contains the list of failed files.
func processFailedFiles(inDir, failedDir string) {
	// Path to failed_files.txt
	failedFilePath := filepath.Join(failedDir, "failed_files.txt")

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
		if !entry.IsDir() { // Ignore directories
			files = append(files, filepath.Join(folderPath, entry.Name()))
		}
	}

	return files, nil
}

// Scan a directory and return the list of files matching the given pattern
func scanDirectoryWithPattern(folderPath, pattern string) ([]string, error) {
	var files []string

	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			matched, err := filepath.Match(pattern, entry.Name()) // Match file pattern
			if err != nil {
				log.Printf("Error matching pattern %s: %v", pattern, err)
				continue
			}
			if matched {
				files = append(files, filepath.Join(folderPath, entry.Name()))
			}
		}
			files = append(files, filepath.Join(folderPath, entry.Name()))
	}
	return files, nil
}

type CoreInterface interface {
	//ReceiveRequests
	ReceiveRequests()
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
	PollFolder(f.FileURI, f.MoveAfterProcess, f.MoveAfterFailure, f.Interval)

}

func (f  *FileInboundAdapter) ReceiveResults() {
	//reveive results
}

func (f  *FileInboundAdapter) Start() {
	//start polling
	go f.StartPolling() //used go routine since there may be another functionalities in fileinbound in furture improvements

}

func (f  *FileInboundAdapter) Stop() {	
	//stop process
}