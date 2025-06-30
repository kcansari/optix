package strategies

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/kcansari/optix/internal/reader"
	"github.com/kcansari/optix/internal/types"
)

type ReplaceProcessorStrategy struct{}

func (rp *ReplaceProcessorStrategy) Process(content *reader.FileContent, options types.ProcessOptions) (*types.ProcessingResult, error) {
	startTime := time.Now()

	if err := rp.ValidateOptions(options); err != nil {
		return nil, fmt.Errorf("invalid replace options: %w", err)
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
		if options.WholeWord {
			escapedPattern = `\b` + escapedPattern + `\b`
		}
		flags := ""
		if !options.CaseSensitive {
			flags = "(?i)"
		}
		pattern, err = regexp.Compile(flags + escapedPattern)
		if err != nil {
			return nil, fmt.Errorf("failed to compile replace pattern: %w", err)
		}
	}

	var backupPath string
	if options.CreateBackup && !options.DryRun {
		backupPath, err = rp.createBackup(options.FileName, options.BackupDir)
		if err != nil {
			return nil, fmt.Errorf("failed to create backup: %w", err)
		}
	}

	originalContent := content.Content
	modifiedContent := pattern.ReplaceAllString(originalContent, options.ReplaceWith)

	matches := pattern.FindAllString(originalContent, -1)
	matchCount := len(matches)

	result := &types.ProcessingResult{
		FileName:        options.FileName,
		Operation:       "replace",
		MatchesFound:    matchCount,
		LinesProcessed:  content.LineCount,
		Success:         true,
		BackupPath:      backupPath,
		ExecutionTime:   time.Since(startTime),
		ModifiedContent: modifiedContent,
	}

	if !options.DryRun {
		outputFile := options.OutputFile
		if outputFile == "" {
			outputFile = options.FileName
		}

		err = os.WriteFile(outputFile, []byte(modifiedContent), 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to write modified content: %w", err)
		}
	}

	return result, nil
}

func (rp *ReplaceProcessorStrategy) createBackup(fileName, backupDir string) (string, error) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to read original file: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	baseName := filepath.Base(fileName)
	backupName := fmt.Sprintf("%s.backup_%s", baseName, timestamp)

	var backupPath string
	if backupDir != "" {
		err = os.MkdirAll(backupDir, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to create backup directory: %w", err)
		}
		backupPath = filepath.Join(backupDir, backupName)
	} else {
		backupPath = fileName + ".backup_" + timestamp
	}

	err = os.WriteFile(backupPath, content, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write backup file: %w", err)
	}

	return backupPath, nil
}

func (rp *ReplaceProcessorStrategy) GetOperationType() string {
	return "replace"
}

func (rp *ReplaceProcessorStrategy) ValidateOptions(options types.ProcessOptions) error {
	if options.Pattern == "" {
		return fmt.Errorf("search pattern cannot be empty")
	}
	if options.ReplaceWith == "" {
		return fmt.Errorf("replacement text cannot be empty")
	}
	return nil
}
