// Package optix contains the CLI commands for the Optix file processor.
// This file implements the 'search' command for pattern matching in files.
package process

import (
	"fmt"
	"path/filepath"

	"github.com/kcansari/optix/cmd"
	"github.com/kcansari/optix/internal/processor"
	"github.com/kcansari/optix/internal/reader"
	"github.com/kcansari/optix/internal/validator"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command.
// This command searches for patterns in files with regex support.
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for patterns in files",
	Long: `Search for text patterns in files with advanced regex support.

The search command supports:
  - Regular expressions and literal text matching
  - Case-sensitive and case-insensitive searches
  - Whole word matching
  - Context lines around matches
  - Multiple file processing with glob patterns

Examples:
  optix search --pattern "error" --files "*.log"
  optix search --pattern "user\d+" --regex --files "data.txt"
  optix search --pattern "TODO" --context 2 --files "*.go"
  optix search --pattern "config" --whole-word --files "*.json"`,

	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flag values
		pattern, _ := cmd.Flags().GetString("pattern")
		files, _ := cmd.Flags().GetString("files")
		regexMode, _ := cmd.Flags().GetBool("regex")
		caseSensitive, _ := cmd.Flags().GetBool("case-sensitive")
		wholeWord, _ := cmd.Flags().GetBool("whole-word")
		contextLines, _ := cmd.Flags().GetInt("context")

		// Validate required flags
		if pattern == "" {
			return fmt.Errorf("pattern is required (use --pattern flag)")
		}
		if files == "" {
			return fmt.Errorf("files pattern is required (use --files flag)")
		}

		// Find matching files
		matchingFiles, err := filepath.Glob(files)
		if err != nil {
			return fmt.Errorf("invalid file pattern '%s': %w", files, err)
		}

		if len(matchingFiles) == 0 {
			return fmt.Errorf("no files found matching pattern '%s'", files)
		}

		// Create processor strategy
		processorStrategy := processor.NewTextProcessorStrategy()
		readerStrategy := reader.NewFileReaderStrategy()
		validatorStrategy := validator.NewValidatorStrategy(validator.NewBasicFileValidator())

		totalMatches := 0
		totalFiles := 0

		fmt.Printf("🔍 Searching for pattern: %s\n", pattern)
		fmt.Printf("📁 Files: %s\n", files)
		if regexMode {
			fmt.Printf("🔧 Mode: Regular Expression\n")
		} else {
			fmt.Printf("🔧 Mode: Literal Text\n")
		}
		fmt.Printf("📊 Case Sensitive: %t\n", caseSensitive)
		if wholeWord {
			fmt.Printf("🔤 Whole Word: %t\n", wholeWord)
		}
		if contextLines > 0 {
			fmt.Printf("📄 Context Lines: %d\n", contextLines)
		}
		fmt.Println("─────────────────────────────────────────────────────")

		// Process each file
		for _, fileName := range matchingFiles {
			// Validate file
			if err := validatorStrategy.ValidateFile(fileName); err != nil {
				fmt.Printf("❌ Skipping '%s': %v\n", fileName, err)
				continue
			}

			// Read file content
			content, err := readerStrategy.ReadFile(fileName)
			if err != nil {
				fmt.Printf("❌ Failed to read '%s': %v\n", fileName, err)
				continue
			}

			// Prepare processing options
			options := processor.ProcessOptions{
				Pattern:       pattern,
				RegexMode:     regexMode,
				CaseSensitive: caseSensitive,
				WholeWord:     wholeWord,
				ContextLines:  contextLines,
				FileName:      fileName,
			}

			// Process the file
			result, err := processorStrategy.ProcessText("search", content, options)
			if err != nil {
				fmt.Printf("❌ Search failed for '%s': %v\n", fileName, err)
				continue
			}

			// Display results
			if result.MatchesFound > 0 {
				fmt.Printf("\n📄 %s (%d matches)\n", fileName, result.MatchesFound)
				totalMatches += result.MatchesFound
				totalFiles++

				// For now, we'll display a summary. In a full implementation,
				// we'd want to return the actual SearchResult objects and display them
				fmt.Printf("   ✅ Found %d matches in %d lines\n", result.MatchesFound, result.LinesProcessed)
			}
		}

		// Display summary
		fmt.Println("\n─────────────────────────────────────────────────────")
		fmt.Printf("📊 Search Summary:\n")
		fmt.Printf("   🎯 Total matches: %d\n", totalMatches)
		fmt.Printf("   📁 Files with matches: %d\n", totalFiles)
		fmt.Printf("   📝 Files processed: %d\n", len(matchingFiles))

		if totalMatches == 0 {
			fmt.Printf("   ℹ️  No matches found for pattern '%s'\n", pattern)
		}

		return nil
	},
}

// init function registers the search command and its flags.
func init() {
	cmd.RootCmd.AddCommand(searchCmd)

	// Add flags for search options
	searchCmd.Flags().StringP("pattern", "p", "", "Search pattern (required)")
	searchCmd.Flags().StringP("files", "f", "", "File pattern to search (required, supports glob)")
	searchCmd.Flags().BoolP("regex", "r", false, "Use regular expression mode")
	searchCmd.Flags().BoolP("case-sensitive", "c", false, "Case sensitive search")
	searchCmd.Flags().BoolP("whole-word", "w", false, "Match whole words only")
	searchCmd.Flags().IntP("context", "C", 0, "Number of context lines to show around matches")

	// Mark required flags
	searchCmd.MarkFlagRequired("pattern")
	searchCmd.MarkFlagRequired("files")
}
