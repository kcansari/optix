package reader

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/kcansari/optix/internal/types"
)

type FileContent = types.FileContent

type FileReader = types.FileReader

type FileReaderStrategy struct {
	readers []FileReader
}

func NewFileReaderStrategy() *FileReaderStrategy {
	return &FileReaderStrategy{
		readers: []FileReader{},
	}
}

func (frs *FileReaderStrategy) AddReader(reader FileReader) {
	frs.readers = append(frs.readers, reader)
}

func (frs *FileReaderStrategy) ReadFile(filename string) (*FileContent, error) {
	extension := filepath.Ext(filename)

	for _, reader := range frs.readers {
		if reader.SupportsFileType(extension) {

			return reader.Read(filename)
		}
	}

	supportedTypes := frs.GetSupportedTypes()
	return nil, fmt.Errorf("unsupported file type '%s' for file '%s'. Supported types: %s",
		extension, filename, strings.Join(supportedTypes, ", "))
}

func (frs *FileReaderStrategy) GetSupportedTypes() []string {
	var types []string
	extensionSet := make(map[string]bool)

	for _, reader := range frs.readers {
		for _, ext := range reader.SupportedExtensions() {
			if !extensionSet[ext] {
				extensionSet[ext] = true
				types = append(types, ext)
			}
		}
	}

	return types
}

func (frs *FileReaderStrategy) GetReaderForExtension(extension string) FileReader {
	for _, reader := range frs.readers {
		if reader.SupportsFileType(extension) {
			return reader
		}
	}
	return nil
}

func (frs *FileReaderStrategy) GetReaderCount() int {
	return len(frs.readers)
}
