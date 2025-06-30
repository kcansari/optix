package processor

import (
	"fmt"
	"strings"

	"github.com/kcansari/optix/internal/reader"
	"github.com/kcansari/optix/internal/types"
)

type TextProcessor = types.TextProcessor
type ProcessOptions = types.ProcessOptions
type ProcessingResult = types.ProcessingResult

type TextProcessorStrategy struct {
	processors map[string]TextProcessor
}

func NewTextProcessorStrategy() *TextProcessorStrategy {
	return &TextProcessorStrategy{
		processors: make(map[string]TextProcessor),
	}
}

func (tps *TextProcessorStrategy) AddProcessor(processor TextProcessor) {
	tps.processors[processor.GetOperationType()] = processor
}

func (tps *TextProcessorStrategy) ProcessText(operationType string, content *reader.FileContent, options ProcessOptions) (*ProcessingResult, error) {
	processor, exists := tps.processors[operationType]
	if !exists {
		return nil, fmt.Errorf("unsupported operation type '%s'. Available types: %s",
			operationType, strings.Join(tps.GetSupportedOperations(), ", "))
	}

	return processor.Process(content, options)
}

func (tps *TextProcessorStrategy) GetSupportedOperations() []string {
	var operations []string
	for op := range tps.processors {
		operations = append(operations, op)
	}
	return operations
}

func (tps *TextProcessorStrategy) GetProcessor(operationType string) TextProcessor {
	return tps.processors[operationType]
}
