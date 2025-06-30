package strategies

import (
	"github.com/kcansari/optix/internal/processor"
)

func NewDefaultTextProcessorStrategy() *processor.TextProcessorStrategy {
	strategy := processor.NewTextProcessorStrategy()

	// Register all available text processors
	strategy.AddProcessor(&FilterProcessorStrategy{})
	strategy.AddProcessor(&ReplaceProcessorStrategy{})
	strategy.AddProcessor(&FilterProcessorStrategy{})
	strategy.AddProcessor(&TransformProcessorStrategy{})

	return strategy
}
