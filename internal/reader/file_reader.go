// Package reader provides file reading functionality using the Strategy Pattern.
// This package allows different file types to be read using different strategies,
// making it easy to add new file types without modifying existing code.
package reader

import (
	"bufio"         // Package for buffered I/O operations
	"encoding/csv"  // Package for reading CSV files
	"encoding/json" // Package for JSON encoding/decoding
	"fmt"           // Package for formatted I/O operations
	"io"            // Package for I/O primitives
	"os"            // Package for operating system interface
	"path/filepath" // Package for file path manipulation
	"strings"       // Package for string operations
)

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

// FileReader is the Strategy interface that defines how different file types should be read.
// Any type that implements this interface can be used as a file reading strategy.
// This follows the Strategy Pattern from our CLAUDE.md instructions.
type FileReader interface {
	// Read takes a filename and returns FileContent and error.
	// The error will be non-nil if something goes wrong during reading.
	Read(filename string) (*FileContent, error)

	// SupportsFileType checks if this reader can handle the given file extension.
	// For example, a TXT reader would return true for ".txt" files.
	SupportsFileType(extension string) bool

	// SupportedExtensions returns a slice of file extensions this reader supports.
	// This removes hardcoding and allows dynamic discovery of supported types.
	SupportedExtensions() []string
}

// TextFileReader is a concrete implementation of FileReader for text files.
// This struct implements the FileReader interface for handling .txt files.
type TextFileReader struct{}

// Read implements the FileReader interface for text files.
// Uses buffered reading to handle large files efficiently without loading everything into memory.
func (r *TextFileReader) Read(filename string) (*FileContent, error) {
	// Open the file for reading
	// os.Open returns a file handle and an error
	file, err := os.Open(filename)
	if err != nil {
		// Use error wrapping with %w for better error handling (Go 1.13+)
		return nil, fmt.Errorf("failed to open text file '%s': %w", filename, err)
	}
	// defer ensures the file is closed when this function returns
	// This is Go's way of ensuring cleanup happens automatically
	defer file.Close()

	// Get file information (size, modification time, etc.)
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info for '%s': %w", filename, err)
	}

	// Create a buffered scanner to read the file line by line
	// bufio.Scanner provides a convenient way to read input line by line
	scanner := bufio.NewScanner(file)

	// Initialize slices to store lines and build content
	// In Go, slices are dynamic arrays that can grow as needed
	var lines []string
	var contentBuilder strings.Builder // Efficient way to build strings
	var wordCount int                  // Count words as we process each line for efficiency

	// Read file line by line using buffered approach
	for scanner.Scan() {
		// scanner.Text() returns the current line without the newline character
		line := scanner.Text()

		// Add line to our lines slice
		// append() adds elements to a slice and returns the new slice
		lines = append(lines, line)

		// Build the complete content string
		// WriteString is more efficient than string concatenation
		contentBuilder.WriteString(line)
		contentBuilder.WriteString("\n") // Add newline back

		// Count words in this line immediately to avoid processing entire content later
		wordCount += len(strings.Fields(line))
	}

	// Check if scanner encountered any errors while reading
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading text file '%s': %w", filename, err)
	}

	// Get the complete content as a string
	content := contentBuilder.String()

	// Create and return the FileContent struct
	// The &FileContent{} syntax creates a pointer to a new FileContent struct
	return &FileContent{
		Content:   content,
		Lines:     lines,
		FileType:  "txt",
		Size:      fileInfo.Size(),
		LineCount: len(lines),
		WordCount: wordCount,
	}, nil
}

// SupportsFileType checks if this reader can handle the given file extension.
func (r *TextFileReader) SupportsFileType(extension string) bool {
	// Check if the extension is in our supported list
	for _, ext := range r.SupportedExtensions() {
		if strings.ToLower(extension) == ext {
			return true
		}
	}
	return false
}

// SupportedExtensions returns the file extensions supported by TextFileReader.
func (r *TextFileReader) SupportedExtensions() []string {
	return []string{".txt", ".text", ".log"} // More comprehensive support
}

// CSVFileReader is a concrete implementation of FileReader for CSV files.
type CSVFileReader struct{}

// Read implements the FileReader interface for CSV files.
// Uses streaming record-by-record reading for better memory efficiency with large CSV files.
func (r *CSVFileReader) Read(filename string) (*FileContent, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file '%s': %w", filename, err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info for '%s': %w", filename, err)
	}

	// Create a buffered CSV reader for efficient processing
	bufferedReader := bufio.NewReader(file)
	csvReader := csv.NewReader(bufferedReader)

	// Build content string from CSV records using streaming approach
	var contentBuilder strings.Builder
	var lines []string
	var wordCount int
	var recordCount int

	// Stream records one by one to avoid loading entire CSV into memory
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break // End of file reached
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV record in file '%s': %w", filename, err)
		}

		// Process record immediately
		line := strings.Join(record, ",")
		lines = append(lines, line)
		contentBuilder.WriteString(line)
		contentBuilder.WriteString("\n")

		// Count words in this record
		for _, field := range record {
			wordCount += len(strings.Fields(field))
		}
		recordCount++
	}

	content := contentBuilder.String()

	return &FileContent{
		Content:   content,
		Lines:     lines,
		FileType:  "csv",
		Size:      fileInfo.Size(),
		LineCount: recordCount,
		WordCount: wordCount,
	}, nil
}

// SupportsFileType checks if this reader can handle CSV files.
func (r *CSVFileReader) SupportsFileType(extension string) bool {
	for _, ext := range r.SupportedExtensions() {
		if strings.ToLower(extension) == ext {
			return true
		}
	}
	return false
}

// SupportedExtensions returns the file extensions supported by CSVFileReader.
func (r *CSVFileReader) SupportedExtensions() []string {
	return []string{".csv", ".tsv"} // Support both CSV and TSV
}

// JSONFileReader is a concrete implementation of FileReader for JSON files.
type JSONFileReader struct{}

// Read implements the FileReader interface for JSON files.
// Uses buffered reading and streaming validation for efficient processing of large JSON files.
func (r *JSONFileReader) Read(filename string) (*FileContent, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open JSON file '%s': %w", filename, err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info for '%s': %w", filename, err)
	}

	// Create a buffered reader for efficient processing
	bufferedReader := bufio.NewReader(file)

	// Use buffered approach similar to TextFileReader
	// Read the file line by line while building content
	var lines []string
	var contentBuilder strings.Builder
	var wordCount int

	// Create a scanner to read line by line
	scanner := bufio.NewScanner(bufferedReader)

	// Read file line by line for memory efficiency
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
		contentBuilder.WriteString(line)
		contentBuilder.WriteString("\n")

		// Count words in this line immediately
		wordCount += len(strings.Fields(line))
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading JSON file '%s': %w", filename, err)
	}

	// Get the complete content as string
	contentStr := contentBuilder.String()

	// Validate JSON using streaming decoder from the content we just read
	// This is more efficient than reading the file again
	decoder := json.NewDecoder(strings.NewReader(contentStr))

	var jsonData interface{}
	if err := decoder.Decode(&jsonData); err != nil {
		return nil, fmt.Errorf("file '%s' contains invalid JSON: %w", filename, err)
	}

	return &FileContent{
		Content:   contentStr,
		Lines:     lines,
		FileType:  "json",
		Size:      fileInfo.Size(),
		LineCount: len(lines),
		WordCount: wordCount,
	}, nil
}

// SupportsFileType checks if this reader can handle JSON files.
func (r *JSONFileReader) SupportsFileType(extension string) bool {
	for _, ext := range r.SupportedExtensions() {
		if strings.ToLower(extension) == ext {
			return true
		}
	}
	return false
}

// SupportedExtensions returns the file extensions supported by JSONFileReader.
func (r *JSONFileReader) SupportedExtensions() []string {
	return []string{".json", ".jsonl", ".ndjson"} // Support multiple JSON formats
}

// FileReaderStrategy is the Context class in the Strategy Pattern.
// It holds a reference to a FileReader strategy and delegates the reading operation to it.
// Renamed from ReaderStrategy for better clarity and specificity.
type FileReaderStrategy struct {
	// readers is a slice of available file readers
	// We can have multiple readers and choose the appropriate one based on file type
	readers []FileReader
}

// NewFileReaderStrategy creates a new FileReaderStrategy with default readers.
// This is a constructor function - Go doesn't have constructors like other languages,
// but by convention we create New* functions that return initialized structs.
func NewFileReaderStrategy() *FileReaderStrategy {
	return &FileReaderStrategy{
		readers: []FileReader{
			&TextFileReader{}, // Create instances of each reader type
			&CSVFileReader{},
			&JSONFileReader{},
		},
	}
}

// NewReaderStrategy creates a new FileReaderStrategy with default readers.
// This function is kept for backward compatibility.
// DEPRECATED: Use NewFileReaderStrategy instead.
func NewReaderStrategy() *FileReaderStrategy {
	return NewFileReaderStrategy()
}

// AddReader allows adding a new file reader strategy at runtime.
// This makes our system extensible - we can add new file types without changing existing code.
func (frs *FileReaderStrategy) AddReader(reader FileReader) {
	frs.readers = append(frs.readers, reader)
}

// ReadFile reads a file using the appropriate strategy based on file extension.
// This is the main method that clients will use - it automatically selects the right strategy.
func (frs *FileReaderStrategy) ReadFile(filename string) (*FileContent, error) {
	// Extract the file extension using filepath.Ext
	// filepath.Ext returns the file name extension (including the dot)
	extension := filepath.Ext(filename)

	// Find the appropriate reader for this file type
	for _, reader := range frs.readers {
		if reader.SupportsFileType(extension) {
			// Found a reader that can handle this file type
			return reader.Read(filename)
		}
	}

	// Provide helpful error message with supported types
	supportedTypes := frs.GetSupportedTypes()
	return nil, fmt.Errorf("unsupported file type '%s' for file '%s'. Supported types: %s",
		extension, filename, strings.Join(supportedTypes, ", "))
}

// GetSupportedTypes returns a list of supported file extensions.
// This method now dynamically discovers supported types from readers,
// removing hardcoded extension lists for better extensibility.
func (frs *FileReaderStrategy) GetSupportedTypes() []string {
	var types []string
	extensionSet := make(map[string]bool) // Use map to avoid duplicates

	// Dynamically gather supported extensions from all readers
	for _, reader := range frs.readers {
		for _, ext := range reader.SupportedExtensions() {
			if !extensionSet[ext] {
				extensionSet[ext] = true
				types = append(types, ext)
			}
		}
	}

	return types
}

// GetReaderForExtension returns the reader that supports the given extension.
// This is useful for getting reader-specific information or capabilities.
func (frs *FileReaderStrategy) GetReaderForExtension(extension string) FileReader {
	for _, reader := range frs.readers {
		if reader.SupportsFileType(extension) {
			return reader
		}
	}
	return nil
}

// GetReaderCount returns the number of registered readers.
func (frs *FileReaderStrategy) GetReaderCount() int {
	return len(frs.readers)
}
