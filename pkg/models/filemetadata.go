package models

// FileMetadata stores metadata of the file
type FileMetadata struct {
	Name     string // File name
	Size     int64  // File size in bytes
	FileType string // File extension
	FilePath string // File path
}