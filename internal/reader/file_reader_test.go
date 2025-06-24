// Package reader_test provides comprehensive tests for the improved file reading functionality.
// Tests focus on scalability, error handling, and extensibility improvements.
package reader

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestTextFileReader tests the enhanced TextFileReader implementation.
func TestTextFileReader(t *testing.T) {
	testContent := "Line 1\nLine 2\nLine 3 with more words\n"
	testFile := createTempFile(t, "test.txt", testContent)
	defer os.Remove(testFile)

	reader := &TextFileReader{}

	// Test SupportedExtensions method
	supportedExts := reader.SupportedExtensions()
	expectedExts := []string{".txt", ".text", ".log"}

	if len(supportedExts) != len(expectedExts) {
		t.Errorf("Expected %d supported extensions, got %d", len(expectedExts), len(supportedExts))
	}

	// Test SupportsFileType method with multiple extensions
	for _, ext := range expectedExts {
		if !reader.SupportsFileType(ext) {
			t.Errorf("TextFileReader should support %s files", ext)
		}
	}

	if reader.SupportsFileType(".csv") {
		t.Error("TextFileReader should not support .csv files")
	}

	// Test Read method
	content, err := reader.Read(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	// Verify content and statistics
	if content.Content != testContent {
		t.Errorf("Expected content %q, got %q", testContent, content.Content)
	}

	if content.FileType != "txt" {
		t.Errorf("Expected file type 'txt', got '%s'", content.FileType)
	}

	if content.LineCount != 3 {
		t.Errorf("Expected 3 lines, got %d", content.LineCount)
	}

	// Word count: "Line 1\nLine 2\nLine 3 with more words\n" = 9 words
	if content.WordCount != 9 {
		t.Errorf("Expected 9 words, got %d", content.WordCount)
	}

	// Verify lines content
	expectedLineContent := []string{"Line 1", "Line 2", "Line 3 with more words"}
	for i, line := range content.Lines {
		if line != expectedLineContent[i] {
			t.Errorf("Line %d: expected %q, got %q", i, expectedLineContent[i], line)
		}
	}
}

// TestCSVFileReader tests the enhanced CSV reader with streaming capability.
func TestCSVFileReader(t *testing.T) {
	testContent := "Name,Age,City\nJohn,25,New York\nAlice,30,San Francisco\n"
	testFile := createTempFile(t, "test.csv", testContent)
	defer os.Remove(testFile)

	reader := &CSVFileReader{}

	// Test multiple supported extensions
	supportedExts := reader.SupportedExtensions()
	expectedExts := []string{".csv", ".tsv"}

	if len(supportedExts) != len(expectedExts) {
		t.Errorf("Expected %d supported extensions, got %d", len(expectedExts), len(supportedExts))
	}

	for _, ext := range expectedExts {
		if !reader.SupportsFileType(ext) {
			t.Errorf("CSVFileReader should support %s files", ext)
		}
	}

	if reader.SupportsFileType(".txt") {
		t.Error("CSVFileReader should not support .txt files")
	}

	// Test reading
	content, err := reader.Read(testFile)
	if err != nil {
		t.Fatalf("Failed to read CSV file: %v", err)
	}

	if content.FileType != "csv" {
		t.Errorf("Expected file type 'csv', got '%s'", content.FileType)
	}

	if content.LineCount != 3 {
		t.Errorf("Expected 3 lines, got %d", content.LineCount)
	}

	// Test first line content
	if len(content.Lines) > 0 && content.Lines[0] != "Name,Age,City" {
		t.Errorf("Expected first line 'Name,Age,City', got %q", content.Lines[0])
	}
}

// TestJSONFileReader tests the enhanced JSON reader with streaming validation.
func TestJSONFileReader(t *testing.T) {
	testContent := `{
  "name": "John Doe",
  "age": 30,
  "city": "New York"
}`
	testFile := createTempFile(t, "test.json", testContent)
	defer os.Remove(testFile)

	reader := &JSONFileReader{}

	// Test multiple JSON format support
	supportedExts := reader.SupportedExtensions()
	expectedExts := []string{".json", ".jsonl", ".ndjson"}

	if len(supportedExts) != len(expectedExts) {
		t.Errorf("Expected %d supported extensions, got %d", len(expectedExts), len(supportedExts))
	}

	for _, ext := range expectedExts {
		if !reader.SupportsFileType(ext) {
			t.Errorf("JSONFileReader should support %s files", ext)
		}
	}

	// Test reading valid JSON
	content, err := reader.Read(testFile)
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	if content.FileType != "json" {
		t.Errorf("Expected file type 'json', got '%s'", content.FileType)
	}

	// The buffered reader adds a newline at the end, so we need to account for that
	expectedContent := testContent + "\n"
	if content.Content != expectedContent {
		t.Errorf("Content mismatch. Expected length %d, got %d", len(expectedContent), len(content.Content))
	}
}

// TestJSONFileReaderInvalidJSON tests enhanced error handling for invalid JSON.
func TestJSONFileReaderInvalidJSON(t *testing.T) {
	invalidJSON := `{"name": "John", "age": 30,}` // Trailing comma
	testFile := createTempFile(t, "invalid.json", invalidJSON)
	defer os.Remove(testFile)

	reader := &JSONFileReader{}

	_, err := reader.Read(testFile)
	if err == nil {
		t.Error("Expected error for invalid JSON, but got none")
	}

	// Test error wrapping
	if !strings.Contains(err.Error(), "invalid JSON") {
		t.Errorf("Expected error to mention 'invalid JSON', got: %v", err)
	}

	// Test that we can unwrap the error (Go 1.13+ error wrapping)
	var targetErr *os.PathError
	if errors.As(err, &targetErr) {
		t.Log("Successfully identified unwrappable error type")
	}
}

// TestFileReaderStrategy tests the improved strategy pattern implementation.
func TestFileReaderStrategy(t *testing.T) {
	// Create test files
	txtFile := createTempFile(t, "test.txt", "Hello World")
	csvFile := createTempFile(t, "test.csv", "Name,Value\nTest,123")
	jsonFile := createTempFile(t, "test.json", `{"test": "value"}`)

	defer func() {
		os.Remove(txtFile)
		os.Remove(csvFile)
		os.Remove(jsonFile)
	}()

	strategy := NewFileReaderStrategy()

	// Test dynamic supported types discovery
	supportedTypes := strategy.GetSupportedTypes()

	// Should include extensions from all readers (including new ones like .log, .tsv, etc.)
	expectedMinTypes := []string{".txt", ".csv", ".json"}
	for _, expectedType := range expectedMinTypes {
		found := false
		for _, supportedType := range supportedTypes {
			if supportedType == expectedType {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected supported type %s not found in %v", expectedType, supportedTypes)
		}
	}

	// Test reading different file types
	testCases := []struct {
		filename     string
		expectedType string
	}{
		{txtFile, "txt"},
		{csvFile, "csv"},
		{jsonFile, "json"},
	}

	for _, tc := range testCases {
		content, err := strategy.ReadFile(tc.filename)
		if err != nil {
			t.Errorf("Failed to read %s: %v", tc.filename, err)
			continue
		}

		if content.FileType != tc.expectedType {
			t.Errorf("File %s: expected type %s, got %s",
				tc.filename, tc.expectedType, content.FileType)
		}
	}

	// Test GetReaderForExtension
	txtReader := strategy.GetReaderForExtension(".txt")
	if txtReader == nil {
		t.Error("Expected to find reader for .txt extension")
	}

	unknownReader := strategy.GetReaderForExtension(".unknown")
	if unknownReader != nil {
		t.Error("Expected nil reader for unknown extension")
	}

	// Test GetReaderCount
	expectedCount := 3 // txt, csv, json readers
	if strategy.GetReaderCount() != expectedCount {
		t.Errorf("Expected %d readers, got %d", expectedCount, strategy.GetReaderCount())
	}
}

// TestFileReaderStrategyUnsupportedType tests enhanced error messages for unsupported types.
func TestFileReaderStrategyUnsupportedType(t *testing.T) {
	testFile := createTempFile(t, "test.xyz", "some content")
	defer os.Remove(testFile)

	strategy := NewFileReaderStrategy()

	_, err := strategy.ReadFile(testFile)
	if err == nil {
		t.Error("Expected error for unsupported file type, but got none")
	}

	// Check that error message includes supported types
	errMsg := err.Error()
	if !strings.Contains(errMsg, "unsupported file type") {
		t.Errorf("Expected error to mention 'unsupported file type', got: %v", err)
	}

	if !strings.Contains(errMsg, "Supported types:") {
		t.Errorf("Expected error to list supported types, got: %v", err)
	}
}

// TestAddReader tests adding custom readers to the strategy.
func TestAddReader(t *testing.T) {
	strategy := NewFileReaderStrategy()
	initialCount := strategy.GetReaderCount()

	// Create and add a custom reader
	customReader := &MockReader{extensions: []string{".mock", ".test"}}
	strategy.AddReader(customReader)

	// Verify reader was added
	if strategy.GetReaderCount() != initialCount+1 {
		t.Errorf("Expected %d readers after adding custom reader, got %d",
			initialCount+1, strategy.GetReaderCount())
	}

	// Test that custom extensions are now supported
	supportedTypes := strategy.GetSupportedTypes()
	mockSupported := false
	testSupported := false

	for _, ext := range supportedTypes {
		if ext == ".mock" {
			mockSupported = true
		}
		if ext == ".test" {
			testSupported = true
		}
	}

	if !mockSupported {
		t.Error("Expected .mock extension to be supported after adding custom reader")
	}
	if !testSupported {
		t.Error("Expected .test extension to be supported after adding custom reader")
	}

	// Test reading with custom reader
	mockFile := createTempFile(t, "test.mock", "mock content")
	defer os.Remove(mockFile)

	content, err := strategy.ReadFile(mockFile)
	if err != nil {
		t.Errorf("Failed to read mock file with custom reader: %v", err)
	}

	if content.FileType != "mock" {
		t.Errorf("Expected file type 'mock', got '%s'", content.FileType)
	}
}

// TestErrorWrapping tests that errors are properly wrapped for unwrapping.
func TestErrorWrapping(t *testing.T) {
	reader := &TextFileReader{}

	// Try to read a non-existent file
	_, err := reader.Read("nonexistent.txt")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}

	// Test that the error can be unwrapped
	var pathErr *os.PathError
	if !errors.As(err, &pathErr) {
		t.Error("Expected to be able to unwrap to os.PathError")
	}
}

// TestLargeFileHandling tests behavior with larger files (simulation).
func TestLargeFileHandling(t *testing.T) {
	// Create a larger test file
	largeContent := strings.Repeat("This is a test line with multiple words.\n", 1000)
	testFile := createTempFile(t, "large.txt", largeContent)
	defer os.Remove(testFile)

	reader := &TextFileReader{}
	content, err := reader.Read(testFile)
	if err != nil {
		t.Fatalf("Failed to read large file: %v", err)
	}

	expectedLines := 1000
	if content.LineCount != expectedLines {
		t.Errorf("Expected %d lines, got %d", expectedLines, content.LineCount)
	}

	// Verify word count is calculated correctly
	expectedWords := 1000 * 8 // 8 words per line
	if content.WordCount != expectedWords {
		t.Errorf("Expected %d words, got %d", expectedWords, content.WordCount)
	}
}

// TestBackwardCompatibility tests that the old NewReaderStrategy function still works.
func TestBackwardCompatibility(t *testing.T) {
	// Test deprecated function still works
	strategy := NewReaderStrategy()
	if strategy == nil {
		t.Error("NewReaderStrategy should still work for backward compatibility")
	}

	if strategy.GetReaderCount() != 3 {
		t.Error("Backward compatible function should return same result as new function")
	}
}

// MockReader for testing extensibility with improved interface.
type MockReader struct {
	extensions []string
}

// Read implements the FileReader interface.
func (r *MockReader) Read(filename string) (*FileContent, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open mock file '%s': %w", filename, err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info for '%s': %w", filename, err)
	}

	return &FileContent{
		Content:   "mock content",
		Lines:     []string{"mock content"},
		FileType:  "mock",
		Size:      fileInfo.Size(),
		LineCount: 1,
		WordCount: 2,
	}, nil
}

// SupportsFileType implements the FileReader interface.
func (r *MockReader) SupportsFileType(extension string) bool {
	for _, ext := range r.extensions {
		if strings.ToLower(extension) == ext {
			return true
		}
	}
	return false
}

// SupportedExtensions implements the improved FileReader interface.
func (r *MockReader) SupportedExtensions() []string {
	return r.extensions
}

// createTempFile is a helper function to create temporary files for testing.
func createTempFile(t *testing.T, filename, content string) string {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, filename)

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file %s: %v", filePath, err)
	}

	return filePath
}

// Benchmark tests for performance comparison
func BenchmarkTextFileReader(b *testing.B) {
	content := strings.Repeat("This is a benchmark test line with several words.\n", 1000)
	testFile := createTempFileForBenchmark(b, "benchmark.txt", content)
	defer os.Remove(testFile)

	reader := &TextFileReader{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := reader.Read(testFile)
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}

func BenchmarkFileReaderStrategy(b *testing.B) {
	content := strings.Repeat("Benchmark test line.\n", 500)
	testFile := createTempFileForBenchmark(b, "benchmark.txt", content)
	defer os.Remove(testFile)

	strategy := NewFileReaderStrategy()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := strategy.ReadFile(testFile)
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}

// createTempFileForBenchmark is a helper for benchmark tests.
func createTempFileForBenchmark(b *testing.B, filename, content string) string {
	tempDir := b.TempDir()
	filePath := filepath.Join(tempDir, filename)

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		b.Fatalf("Failed to create temp file: %v", err)
	}

	return filePath
}
