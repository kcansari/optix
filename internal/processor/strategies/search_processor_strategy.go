// Package processor provides text processing functionality using the Strategy Pattern.
// This file implements the SearchProcessor for pattern matching operations.
package strategies

import (
	"fmt"
	"regexp"
	"time"

	"github.com/kcansari/optix/internal/reader"
	"github.com/kcansari/optix/internal/types"
)

type SearchProcessorStrategy struct{}

func (sp *SearchProcessorStrategy) Process(content *reader.FileContent, options types.ProcessOptions) (*types.ProcessingResult, error) {
	startTime := time.Now()

	if err := sp.ValidateOptions(options); err != nil {
		return nil, fmt.Errorf("invalid search options: %w", err)
	}

	var pattern *regexp.Regexp
	var err error

	if options.RegexMode {
		// Use regex pattern directly
		flags := ""
		if !options.CaseSensitive {
			flags = "(?i)"
		}
		pattern, err = regexp.Compile(flags + options.Pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid regex pattern '%s': %w", options.Pattern, err)
		}
	} else {
		// Escape special regex characters for literal search
		escapedPattern := regexp.QuoteMeta(options.Pattern)
		if options.WholeWord {
			escapedPattern = `\b` + escapedPattern + `\b`
		}
		flags := ""
		if !options.CaseSensitive {
			flags = "(?i)"
		}
		pattern, err = regexp.Compile(flags + escapedPattern)
		if err != nil {
			return nil, fmt.Errorf("failed to compile search pattern: %w", err)
		}
	}

	var results []types.SearchResult
	lines := content.Lines

	for i, line := range lines {
		if pattern.MatchString(line) {
			match := pattern.FindString(line)

			// Get context lines if requested
			var context []string
			if options.ContextLines > 0 {
				start := max(0, i-options.ContextLines)
				end := min(len(lines), i+options.ContextLines+1)
				context = lines[start:end]
			}

			results = append(results, types.SearchResult{
				FileName:   options.FileName,
				LineNumber: i + 1,
				Line:       line,
				Match:      match,
				Context:    context,
			})
		}
	}

	result := &types.ProcessingResult{
		FileName:       options.FileName,
		Operation:      "search",
		MatchesFound:   len(results),
		LinesProcessed: len(lines),
		Success:        true,
		ExecutionTime:  time.Since(startTime),
	}

	return result, nil
}

func (sp *SearchProcessorStrategy) GetOperationType() string {
	return "search"
}

func (sp *SearchProcessorStrategy) ValidateOptions(options types.ProcessOptions) error {
	if options.Pattern == "" {
		return fmt.Errorf("search pattern cannot be empty")
	}
	if options.ContextLines < 0 {
		return fmt.Errorf("context lines cannot be negative")
	}
	return nil
}
