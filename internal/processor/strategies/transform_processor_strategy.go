package strategies

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/kcansari/optix/internal/reader"
	"github.com/kcansari/optix/internal/types"
)

type TransformProcessorStrategy struct{}

func (tp *TransformProcessorStrategy) Process(content *reader.FileContent, options types.ProcessOptions) (*types.ProcessingResult, error) {
	startTime := time.Now()

	if err := tp.ValidateOptions(options); err != nil {
		return nil, fmt.Errorf("invalid transform options: %w", err)
	}

	var transformedContent string

	switch strings.ToLower(options.TransformType) {
	case "upper":
		transformedContent = strings.ToUpper(content.Content)
	case "lower":
		transformedContent = strings.ToLower(content.Content)
	case "title":
		transformedContent = strings.Title(strings.ToLower(content.Content))
	case "trim":
		// Trim whitespace from each line
		var trimmedLines []string
		for _, line := range content.Lines {
			trimmedLines = append(trimmedLines, strings.TrimSpace(line))
		}
		transformedContent = strings.Join(trimmedLines, "\n")
		if len(trimmedLines) > 0 {
			transformedContent += "\n"
		}
	default:
		return nil, fmt.Errorf("unsupported transform type: %s", options.TransformType)
	}

	result := &types.ProcessingResult{
		FileName:        options.FileName,
		Operation:       "transform",
		MatchesFound:    1, // Transformation always affects the entire content
		LinesProcessed:  content.LineCount,
		Success:         true,
		ExecutionTime:   time.Since(startTime),
		ModifiedContent: transformedContent,
	}

	if !options.DryRun {
		outputFile := options.OutputFile
		if outputFile == "" {
			outputFile = options.FileName
		}

		err := os.WriteFile(outputFile, []byte(transformedContent), 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to write transformed content: %w", err)
		}
	}

	return result, nil
}

func (tp *TransformProcessorStrategy) GetOperationType() string {
	return "transform"
}

func (tp *TransformProcessorStrategy) ValidateOptions(options types.ProcessOptions) error {
	if options.TransformType == "" {
		return fmt.Errorf("transform type cannot be empty")
	}

	validTypes := []string{"upper", "lower", "title", "trim"}
	for _, validType := range validTypes {
		if strings.ToLower(options.TransformType) == validType {
			return nil
		}
	}

	return fmt.Errorf("invalid transform type '%s'. Valid types: %s",
		options.TransformType, strings.Join(validTypes, ", "))
}
