package strategies

import (
	"github.com/kcansari/optix/internal/reader"
)

func NewDefaultFileReaderStrategy() *reader.FileReaderStrategy {
	strategy := reader.NewFileReaderStrategy()

	strategy.AddReader(&TextFileReader{})
	strategy.AddReader(&CSVFileReader{})
	strategy.AddReader(&JSONFileReader{})

	return strategy
}
