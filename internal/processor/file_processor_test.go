package processor_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/kcansari/optix/internal/processor/strategies"
	"github.com/kcansari/optix/internal/reader"
	"github.com/kcansari/optix/internal/types"
)

// createTestFileContent creates a test FileContent for testing purposes.
func createTestFileContent(content string) *reader.FileContent {
	lines := strings.Split(content, "\n")
	// Remove last empty line if present
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	wordCount := 0
	for _, line := range lines {
		wordCount += len(strings.Fields(line))
	}

	return &reader.FileContent{
		Content:   content,
		Lines:     lines,
		FileType:  "txt",
		Size:      int64(len(content)),
		LineCount: len(lines),
		WordCount: wordCount,
	}
}

func TestSearchProcessor(t *testing.T) {
	processor := &strategies.SearchProcessorStrategy{}

	// Test data
	testContent := `This is line 1 with error message
This is line 2 with normal text
This is line 3 with ERROR in caps
This is line 4 with debug info`

	content := createTestFileContent(testContent)

	tests := []struct {
		name            string
		options         types.ProcessOptions
		expectedMatches int
		expectError     bool
	}{
		{
			name: "Case sensitive search",
			options: types.ProcessOptions{
				Pattern:       "error",
				CaseSensitive: true,
				FileName:      "test.txt",
			},
			expectedMatches: 1,
			expectError:     false,
		},
		{
			name: "Case insensitive search",
			options: types.ProcessOptions{
				Pattern:       "error",
				CaseSensitive: false,
				FileName:      "test.txt",
			},
			expectedMatches: 2,
			expectError:     false,
		},
		{
			name: "Regex search",
			options: types.ProcessOptions{
				Pattern:   "line \\d+",
				RegexMode: true,
				FileName:  "test.txt",
			},
			expectedMatches: 4,
			expectError:     false,
		},
		{
			name: "Empty pattern",
			options: types.ProcessOptions{
				Pattern:  "",
				FileName: "test.txt",
			},
			expectedMatches: 0,
			expectError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := processor.Process(content, tt.options)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result.MatchesFound != tt.expectedMatches {
				t.Errorf("Expected %d matches, got %d", tt.expectedMatches, result.MatchesFound)
			}

			if result.Operation != "search" {
				t.Errorf("Expected operation 'search', got '%s'", result.Operation)
			}
		})
	}
}

func TestReplaceProcessor(t *testing.T) {
	processor := &strategies.ReplaceProcessorStrategy{}

	// Test data
	testContent := `This is old text
Replace old with new
Keep old values here`

	content := createTestFileContent(testContent)

	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "test_replace_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Write test content to temp file
	if _, err := tmpFile.WriteString(testContent); err != nil {
		t.Fatalf("Failed to write test content: %v", err)
	}
	tmpFile.Close()

	tests := []struct {
		name            string
		options         types.ProcessOptions
		expectedMatches int
		expectedContent string
		expectError     bool
	}{
		{
			name: "Simple replace",
			options: types.ProcessOptions{
				Pattern:     "old",
				ReplaceWith: "new",
				FileName:    tmpFile.Name(),
				DryRun:      true,
			},
			expectedMatches: 3,
			expectedContent: `This is new text
Replace new with new
Keep new values here`,
			expectError: false,
		},
		{
			name: "Regex replace",
			options: types.ProcessOptions{
				Pattern:     "\\bold\\b",
				ReplaceWith: "new",
				RegexMode:   true,
				FileName:    tmpFile.Name(),
				DryRun:      true,
			},
			expectedMatches: 3,
			expectedContent: `This is new text
Replace new with new
Keep new values here`,
			expectError: false,
		},
		{
			name: "Empty pattern",
			options: types.ProcessOptions{
				Pattern:     "",
				ReplaceWith: "new",
				FileName:    tmpFile.Name(),
				DryRun:      true,
			},
			expectedMatches: 0,
			expectError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := processor.Process(content, tt.options)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result.MatchesFound != tt.expectedMatches {
				t.Errorf("Expected %d matches, got %d", tt.expectedMatches, result.MatchesFound)
			}

			if result.Operation != "replace" {
				t.Errorf("Expected operation 'replace', got '%s'", result.Operation)
			}

			if result.ModifiedContent != tt.expectedContent {
				t.Errorf("Expected content:\n%s\nGot:\n%s", tt.expectedContent, result.ModifiedContent)
			}
		})
	}
}

func TestFilterProcessor(t *testing.T) {
	processor := &strategies.FilterProcessorStrategy{}

	// Test data
	testContent := `INFO: Application started
ERROR: Database connection failed
DEBUG: Processing user data
ERROR: Invalid input format
INFO: Application stopped`

	content := createTestFileContent(testContent)

	tests := []struct {
		name            string
		options         types.ProcessOptions
		expectedMatches int
		expectedLines   int
		expectError     bool
	}{
		{
			name: "Filter ERROR lines",
			options: types.ProcessOptions{
				Pattern:  "ERROR",
				FileName: "test.txt",
			},
			expectedMatches: 2,
			expectedLines:   2,
			expectError:     false,
		},
		{
			name: "Invert filter (non-ERROR lines)",
			options: types.ProcessOptions{
				Pattern:     "ERROR",
				InvertMatch: true,
				FileName:    "test.txt",
			},
			expectedMatches: 3,
			expectedLines:   3,
			expectError:     false,
		},
		{
			name: "Case insensitive filter",
			options: types.ProcessOptions{
				Pattern:       "info",
				CaseSensitive: false,
				FileName:      "test.txt",
			},
			expectedMatches: 2,
			expectedLines:   2,
			expectError:     false,
		},
		{
			name: "Empty pattern",
			options: types.ProcessOptions{
				Pattern:  "",
				FileName: "test.txt",
			},
			expectedMatches: 0,
			expectError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := processor.Process(content, tt.options)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result.MatchesFound != tt.expectedMatches {
				t.Errorf("Expected %d matches, got %d", tt.expectedMatches, result.MatchesFound)
			}

			if result.Operation != "filter" {
				t.Errorf("Expected operation 'filter', got '%s'", result.Operation)
			}

			// Count lines in modified content
			lines := strings.Split(strings.TrimSpace(result.ModifiedContent), "\n")
			if result.ModifiedContent == "" {
				lines = []string{}
			}
			actualLines := len(lines)
			if result.ModifiedContent == "" {
				actualLines = 0
			}

			if actualLines != tt.expectedLines {
				t.Errorf("Expected %d lines in output, got %d", tt.expectedLines, actualLines)
			}
		})
	}
}

func TestTransformProcessor(t *testing.T) {
	processor := &strategies.TransformProcessorStrategy{}

	// Test data
	testContent := `  Hello World
  UPPER CASE TEXT
  lower case text  `

	content := createTestFileContent(testContent)

	tests := []struct {
		name            string
		options         types.ProcessOptions
		expectedContent string
		expectError     bool
	}{
		{
			name: "Transform to uppercase",
			options: types.ProcessOptions{
				TransformType: "upper",
				FileName:      "test.txt",
				DryRun:        true,
			},
			expectedContent: `  HELLO WORLD
  UPPER CASE TEXT
  LOWER CASE TEXT  `,
			expectError: false,
		},
		{
			name: "Transform to lowercase",
			options: types.ProcessOptions{
				TransformType: "lower",
				FileName:      "test.txt",
				DryRun:        true,
			},
			expectedContent: `  hello world
  upper case text
  lower case text  `,
			expectError: false,
		},
		{
			name: "Trim whitespace",
			options: types.ProcessOptions{
				TransformType: "trim",
				FileName:      "test.txt",
				DryRun:        true,
			},
			expectedContent: `Hello World
UPPER CASE TEXT
lower case text
`,
			expectError: false,
		},
		{
			name: "Invalid transform type",
			options: types.ProcessOptions{
				TransformType: "invalid",
				FileName:      "test.txt",
				DryRun:        true,
			},
			expectedContent: "",
			expectError:     true,
		},
		{
			name: "Empty transform type",
			options: types.ProcessOptions{
				TransformType: "",
				FileName:      "test.txt",
				DryRun:        true,
			},
			expectedContent: "",
			expectError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := processor.Process(content, tt.options)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result.Operation != "transform" {
				t.Errorf("Expected operation 'transform', got '%s'", result.Operation)
			}

			if result.ModifiedContent != tt.expectedContent {
				t.Errorf("Expected content:\n%q\nGot:\n%q", tt.expectedContent, result.ModifiedContent)
			}
		})
	}
}

func TestTextProcessorStrategy(t *testing.T) {
	strategy := strategies.NewDefaultTextProcessorStrategy()

	// Test data
	testContent := createTestFileContent("Hello World\nTest Line")

	tests := []struct {
		name          string
		operationType string
		expectError   bool
	}{
		{
			name:          "Valid search operation",
			operationType: "search",
			expectError:   false,
		},
		{
			name:          "Valid replace operation",
			operationType: "replace",
			expectError:   false,
		},
		{
			name:          "Valid filter operation",
			operationType: "filter",
			expectError:   false,
		},
		{
			name:          "Valid transform operation",
			operationType: "transform",
			expectError:   false,
		},
		{
			name:          "Invalid operation",
			operationType: "invalid",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := types.ProcessOptions{
				Pattern:       "Hello",
				ReplaceWith:   "Hi",
				TransformType: "upper",
				FileName:      "test.txt",
				DryRun:        true,
			}

			_, err := strategy.ProcessText(tt.operationType, testContent, options)

			if tt.expectError && err == nil {
				t.Errorf("Expected error for operation '%s' but got none", tt.operationType)
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for operation '%s': %v", tt.operationType, err)
			}
		})
	}

	// Test supported operations
	supportedOps := strategy.GetSupportedOperations()
	expectedOps := []string{"search", "replace", "filter", "transform"}

	if len(supportedOps) != len(expectedOps) {
		t.Errorf("Expected %d supported operations, got %d", len(expectedOps), len(supportedOps))
	}

	// Test adding custom processor
	customProcessor := &strategies.SearchProcessorStrategy{} // Using SearchProcessor as example
	strategy.AddProcessor(customProcessor)

	processor := strategy.GetProcessor("search")
	if processor == nil {
		t.Errorf("Failed to retrieve search processor")
	}
}

func TestProcessingResultTiming(t *testing.T) {
	processor := &strategies.SearchProcessorStrategy{}
	content := createTestFileContent("test content")

	options := types.ProcessOptions{
		Pattern:  "test",
		FileName: "test.txt",
	}

	start := time.Now()
	result, err := processor.Process(content, options)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.ExecutionTime <= 0 {
		t.Errorf("Expected positive execution time, got %v", result.ExecutionTime)
	}

	if result.ExecutionTime > duration {
		t.Errorf("Execution time %v is greater than actual duration %v", result.ExecutionTime, duration)
	}
}
