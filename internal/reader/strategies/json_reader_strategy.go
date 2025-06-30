package strategies

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/kcansari/optix/internal/types"
)

type JSONFileReader struct{}

func (r *JSONFileReader) Read(filename string) (*types.FileContent, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open JSON file '%s': %w", filename, err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info for '%s': %w", filename, err)
	}

	bufferedReader := bufio.NewReader(file)

	var lines []string
	var contentBuilder strings.Builder
	var wordCount int

	scanner := bufio.NewScanner(bufferedReader)

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
		contentBuilder.WriteString(line)
		contentBuilder.WriteString("\n")

		wordCount += len(strings.Fields(line))
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading JSON file '%s': %w", filename, err)
	}

	contentStr := contentBuilder.String()

	decoder := json.NewDecoder(strings.NewReader(contentStr))

	var jsonData interface{}
	if err := decoder.Decode(&jsonData); err != nil {
		return nil, fmt.Errorf("file '%s' contains invalid JSON: %w", filename, err)
	}

	return &types.FileContent{
		Content:   contentStr,
		Lines:     lines,
		FileType:  "json",
		Size:      fileInfo.Size(),
		LineCount: len(lines),
		WordCount: wordCount,
	}, nil
}

func (r *JSONFileReader) SupportsFileType(extension string) bool {
	for _, ext := range r.SupportedExtensions() {
		if strings.ToLower(extension) == ext {
			return true
		}
	}
	return false
}

func (r *JSONFileReader) SupportedExtensions() []string {
	return []string{".json", ".jsonl", ".ndjson"}
}
