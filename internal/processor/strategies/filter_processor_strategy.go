package strategies

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/kcansari/optix/internal/reader"
	"github.com/kcansari/optix/internal/types"
)

type FilterProcessorStrategy struct{}

// Process performs filtering operations on text content.
func (fp *FilterProcessorStrategy) Process(content *reader.FileContent, options types.ProcessOptions) (*types.ProcessingResult, error) {
	startTime := time.Now()

	if err := fp.ValidateOptions(options); err != nil {
		return nil, fmt.Errorf("invalid filter options: %w", err)
	}

	var pattern *regexp.Regexp
	var err error

	if options.RegexMode {
		flags := ""
		if !options.CaseSensitive {
			flags = "(?i)"
		}
		pattern, err = regexp.Compile(flags + options.Pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid regex pattern '%s': %w", options.Pattern, err)
		}
	} else {
		escapedPattern := regexp.QuoteMeta(options.Pattern)
		flags := ""
		if !options.CaseSensitive {
			flags = "(?i)"
		}
		pattern, err = regexp.Compile(flags + escapedPattern)
		if err != nil {
			return nil, fmt.Errorf("failed to compile filter pattern: %w", err)
		}
	}

	var filteredLines []string
	matchCount := 0

	for _, line := range content.Lines {
		matches := pattern.MatchString(line)

		// Apply invert match logic
		if options.InvertMatch {
			matches = !matches
		}

		if matches {
			if options.OnlyMatching {
				// Extract only the matching part
				match := pattern.FindString(line)
				if match != "" {
					filteredLines = append(filteredLines, match)
					matchCount++
				}
			} else {
				// Include the entire line
				filteredLines = append(filteredLines, line)
				matchCount++
			}
		}
	}

	filteredContent := strings.Join(filteredLines, "\n")
	if len(filteredLines) > 0 {
		filteredContent += "\n"
	}

	result := &types.ProcessingResult{
		FileName:        options.FileName,
		Operation:       "filter",
		MatchesFound:    matchCount,
		LinesProcessed:  content.LineCount,
		Success:         true,
		ExecutionTime:   time.Since(startTime),
		ModifiedContent: filteredContent,
	}

	// Write filtered content if output file is specified and not in dry run mode
	if options.OutputFile != "" && !options.DryRun {
		err = os.WriteFile(options.OutputFile, []byte(filteredContent), 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to write filtered content: %w", err)
		}
	}

	return result, nil
}

func (fp *FilterProcessorStrategy) GetOperationType() string {
	return "filter"
}

func (fp *FilterProcessorStrategy) ValidateOptions(options types.ProcessOptions) error {
	if options.Pattern == "" {
		return fmt.Errorf("filter pattern cannot be empty")
	}
	return nil
}
