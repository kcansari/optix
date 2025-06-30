package strategies

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kcansari/optix/internal/types"
)

type CSVFileReader struct{}

func (r *CSVFileReader) Read(filename string) (*types.FileContent, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file '%s': %w", filename, err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info for '%s': %w", filename, err)
	}

	bufferedReader := bufio.NewReader(file)
	csvReader := csv.NewReader(bufferedReader)

	var contentBuilder strings.Builder
	var lines []string
	var wordCount int
	var recordCount int

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV record in file '%s': %w", filename, err)
		}

		line := strings.Join(record, ",")
		lines = append(lines, line)
		contentBuilder.WriteString(line)
		contentBuilder.WriteString("\n")

		// Count words in this record
		for _, field := range record {
			wordCount += len(strings.Fields(field))
		}
		recordCount++
	}

	content := contentBuilder.String()

	return &types.FileContent{
		Content:   content,
		Lines:     lines,
		FileType:  "csv",
		Size:      fileInfo.Size(),
		LineCount: recordCount,
		WordCount: wordCount,
	}, nil
}

func (r *CSVFileReader) SupportsFileType(extension string) bool {
	for _, ext := range r.SupportedExtensions() {
		if strings.ToLower(extension) == ext {
			return true
		}
	}
	return false
}

func (r *CSVFileReader) SupportedExtensions() []string {
	return []string{".csv", ".tsv"}
}
