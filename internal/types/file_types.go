// Package types contains shared types used across the optix application.
// This package helps avoid circular dependencies by providing common types.
package types

// FileContent represents the content and metadata of a file.
// This struct holds all the information we extract from a file.
type FileContent struct {
	// Content holds the raw content of the file as a string
	Content string

	// Lines contains each line of the file as separate string elements
	// This is useful for line-by-line processing
	Lines []string

	// FileType indicates what type of file this is (txt, csv, json)
	FileType string

	// Size is the file size in bytes
	Size int64

	// LineCount is the total number of lines in the file
	LineCount int

	// WordCount is the total number of words in the file
	WordCount int
}

// FileReader defines the interface that all file readers must implement.
// This is the Strategy interface in the Strategy Pattern.
type FileReader interface {
	// Read reads a file and returns its content and metadata.
	// It takes a filename and returns a FileContent struct with all the extracted information.
	Read(filename string) (*FileContent, error)

	// SupportsFileType checks if this reader can handle the given file extension.
	// For example, a TXT reader would return true for ".txt" files.
	SupportsFileType(extension string) bool

	// SupportedExtensions returns a slice of file extensions this reader supports.
	// This removes hardcoding and allows dynamic discovery of supported types.
	SupportedExtensions() []string
}
