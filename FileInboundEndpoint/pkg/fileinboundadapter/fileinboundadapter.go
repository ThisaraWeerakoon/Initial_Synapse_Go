package fileinboundadapter

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/ThisaraWeerakoon/Initial_Synapse_Go/FileInboundEndpoint/models"
)

// PollFolder continuously reads files at a given interval
func (f  *FileInboundAdapter) PollFolder(ctx context.Context,inDir string, outDir string, failedDir string, interval int, pattern string) {
	newCtx, cancel := context.WithCancel(ctx)
	var wg sync.WaitGroup
	defer cancel()
	ticker := time.NewTicker(time.Duration(interval) * time.Second) // Ensures precise polling
	defer ticker.Stop()

	for range ticker.C {
		startTime := time.Now()
		fmt.Println("\n--- New Polling Event ---")

		// Process failed_files.txt if available
		processFailedFiles(inDir, failedDir)

		// Get the list of files at the start of this polling event
		files, err := scanDirectoryWithPattern(inDir, pattern)
		if err != nil {
			log.Printf("Error scanning directory %s: %v", inDir, err)
			continue
		}

		// Process each file
		for _, file := range files {
			wg.Add(1)
			//have to test is it safe to make go routines for each file arbitrarily
			// A solution may be put a threshold (eg:- 100 files) and then make go routines for each file.If the number of files is greater than the threshold make only upper limit (threshold) of go routines
			go f.ProcessFile(newCtx,&wg,file)
		}

		// Ensure accurate polling interval
		elapsed := time.Since(startTime)
		if elapsed < time.Duration(interval)*time.Second {
			time.Sleep(time.Duration(interval)*time.Second - elapsed)
		}
	}
	wg.Wait()
}

func (f *FileInboundAdapter) ProcessFile(ctx context.Context,parentWg *sync.WaitGroup,file string) {
	defer parentWg.Done()
	extractedFileData, err := ReadFile(file)
	if err != nil {
		log.Printf("Error reading file %s: %v", file, err)
		return
	}

	//Attention : Here I implemented considering same reading go routine taking care the receiving results and it is needed to reconsider the design. Here the design is simpple but think the situation where some of the processed results of the previous iteration coming and there might be not enough threads.
	f.CallCore(*extractedFileData) //Finally f.ReceiveRequests() is called


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
	ReceiveRequests(*models.ExtractedFileDataFromFileAdapter)
}

type FileInboundAdapter struct{
	models.Configurations
	core CoreInterface

}

func NewFileInboundAdapter(config models.Configurations, core CoreInterface) *FileInboundAdapter {
	return &FileInboundAdapter{
		Configurations: config,
		core: core,
	}
}

func (f  *FileInboundAdapter) StartPolling(ctx context.Context, parentWg *sync.WaitGroup) {
	defer parentWg.Done()

	//start polling
	f.PollFolder(ctx,f.FileURI, f.MoveAfterProcess, f.MoveAfterFailure, f.Interval, f.FileNamePattern)

}

//After receiving the results from core this function will move the file if it's success or write to failed_files.txt if it's failure
func (f  *FileInboundAdapter) ReceiveResults(processedMessageFromCore models.ProcessedMessageFromCore){
	//reveive results
	if processedMessageFromCore.IsSuccess {
		// Move file to success directory
		err := MoveFile(processedMessageFromCore.FilePath, f.MoveAfterProcess)
		if err != nil {
			log.Printf("Error moving file %s to %s: %v", processedMessageFromCore.FilePath, f.MoveAfterProcess, err)
		} else {
			fmt.Printf("Moved %s to %s\n", processedMessageFromCore.FilePath, f.MoveAfterProcess)
		}
	} else {
		// Write to failed_files.txt
		failedFilePath := filepath.Join(f.MoveAfterFailure, "failed_files.txt")
		file, err := os.OpenFile(failedFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Error opening failed_files.txt: %v", err)
			return
		}
		defer file.Close()

		if _, err := file.WriteString(processedMessageFromCore.FilePath + "\n"); err != nil {
			log.Printf("Error writing to failed_files.txt: %v", err)
		} else {
			fmt.Printf("Wrote %s to failed_files.txt\n", processedMessageFromCore.FilePath)
		}
	}

}

func (f  *FileInboundAdapter) Start(ctx context.Context, parentWg *sync.WaitGroup) {
	defer parentWg.Done()
	var wg sync.WaitGroup
	newCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	wg.Add(1)
	//start polling
	go f.StartPolling(newCtx,&wg) //used go routine since there may be another functionalities in fileinbound in furture improvements
	// go f.ReceiveResults() //used go routine since there may be another functionalities in fileinbound in furture improvements
	wg.Wait()

}

func (f  *FileInboundAdapter) Stop() {	
	//stop process
}