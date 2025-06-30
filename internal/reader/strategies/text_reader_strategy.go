package strategies

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/kcansari/optix/internal/types"
)

type TextFileReader struct{}

func (r *TextFileReader) Read(filename string) (*types.FileContent, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open text file '%s': %w", filename, err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info for '%s': %w", filename, err)
	}

	scanner := bufio.NewScanner(file)

	var lines []string
	var contentBuilder strings.Builder
	var wordCount int

	for scanner.Scan() {
		line := scanner.Text()

		lines = append(lines, line)

		contentBuilder.WriteString(line)
		contentBuilder.WriteString("\n")

		wordCount += len(strings.Fields(line))
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading text file '%s': %w", filename, err)
	}

	content := contentBuilder.String()

	return &types.FileContent{
		Content:   content,
		Lines:     lines,
		FileType:  "txt",
		Size:      fileInfo.Size(),
		LineCount: len(lines),
		WordCount: wordCount,
	}, nil
}

func (r *TextFileReader) SupportsFileType(extension string) bool {
	for _, ext := range r.SupportedExtensions() {
		if strings.ToLower(extension) == ext {
			return true
		}
	}
	return false
}

func (r *TextFileReader) SupportedExtensions() []string {
	return []string{".txt", ".text", ".log"}
}
