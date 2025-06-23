package validator

import (
	"os"
	"testing"
)

func TestFileValidation(t *testing.T) {
	validator := NewBasicFileValidator()
	strategy := NewValidatorStrategy(validator)

	// Create a test file
	testFile := "test_file.txt"
	file, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.WriteString("test content")
	file.Close()

	// Test valid file
	err = strategy.ValidateFile(testFile)
	if err != nil {
		t.Errorf("Expected no error for valid file, got: %v", err)
	}

	// Test non-existent file
	err = strategy.ValidateFile("non_existent_file.txt")
	if err == nil {
		t.Errorf("Expected error for non-existent file")
	}

	// Test empty filename
	err = strategy.ValidateFile("")
	if err == nil {
		t.Errorf("Expected error for empty filename")
	}

	// Cleanup
	os.Remove(testFile)
}

func TestBasicFileValidator(t *testing.T) {
	validator := NewBasicFileValidator()

	// Test empty filename
	err := validator.Validate("")
	if err == nil {
		t.Error("Expected error for empty filename")
	}

	// Test non-existent file
	err = validator.Validate("non_existent_file.txt")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}
