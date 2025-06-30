// Package processor provides text processing functionality using the Strategy Pattern.
// This file contains shared types, interfaces, and data structures used across
// all text processing operations.
package types

import (
	"time"
)

// SearchResult represents a single search match with context information.
type SearchResult struct {
	FileName   string
	LineNumber int
	Line       string
	Match      string
	Context    []string
}

// ProcessingResult represents the outcome of a text processing operation.
type ProcessingResult struct {
	FileName        string
	Operation       string
	MatchesFound    int
	LinesProcessed  int
	Success         bool
	ErrorMessage    string
	BackupPath      string
	ExecutionTime   time.Duration
	ModifiedContent string
}

// TextProcessor defines the strategy interface for text processing operations.
// This interface allows different text processing strategies to be implemented
// and used interchangeably following the Strategy Pattern.
type TextProcessor interface {
	// Process performs the text processing operation on the given content
	Process(content *FileContent, options ProcessOptions) (*ProcessingResult, error)

	// GetOperationType returns the type of operation this processor performs
	GetOperationType() string

	// ValidateOptions validates the processing options for this specific processor
	ValidateOptions(options ProcessOptions) error
}

// ProcessOptions contains configuration for text processing operations.
type ProcessOptions struct {
	// Search options
	Pattern       string
	RegexMode     bool
	CaseSensitive bool
	WholeWord     bool
	ContextLines  int

	// Replace options
	ReplaceWith  string
	CreateBackup bool
	BackupDir    string

	// Filter options
	InvertMatch  bool
	OnlyMatching bool

	// Transform options
	TransformType string // "upper", "lower", "title", "trim"

	// General options
	FileName   string
	OutputFile string
	DryRun     bool
}
