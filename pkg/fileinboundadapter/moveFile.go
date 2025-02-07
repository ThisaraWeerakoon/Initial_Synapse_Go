package fileinboundadapter

import (
	"os"
	"path/filepath"
	"fmt"
	"io"
)

// MoveFile moves a file from source to destination, handling cross-filesystem moves.
func MoveFile(source, destination string) error {
	// Check if the destination is a directory
	info, err := os.Stat(destination)
	if err == nil && info.IsDir() {
		// Extract the filename from the source path
		filename := filepath.Base(source)
		// Append the filename to the destination directory
		destination = filepath.Join(destination, filename)
	}

	// Try to rename first (fast path)
	err = os.Rename(source, destination)
	if err == nil {
		return nil // Successfully renamed (same file system)
	}

	// If rename fails, assume it's due to a cross-filesystem issue and proceed with copy & delete
	err = copyFile(source, destination)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	// Remove the original file after successful copy
	err = os.Remove(source)
	if err != nil {
		return fmt.Errorf("failed to remove source file: %w", err)
	}

	return nil
}

// copyFile copies the contents of the source file to the destination file.
func copyFile(source, destination string) error {
	srcFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	// Ensure the new file has the same permissions as the source file
	srcInfo, err := os.Stat(source)
	if err != nil {
		return err
	}
	return os.Chmod(destination, srcInfo.Mode())
}
