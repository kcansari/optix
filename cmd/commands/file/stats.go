// Package optix contains the CLI commands for the Optix file processor.
// This file implements the 'stats' command that displays detailed file statistics.
package file

import (
	"fmt"     // Package for formatted I/O operations
	"strings" // Package for string operations

	"github.com/kcansari/optix/cmd"
	"github.com/kcansari/optix/internal/reader"            // Our file reader package
	"github.com/kcansari/optix/internal/reader/strategies" // Reader strategies
	"github.com/kcansari/optix/internal/validator"         // Our file validator package
	"github.com/spf13/cobra"                               // CLI framework
)

// statsCmd represents the stats command.
// This command displays detailed statistics about a file including:
// - File size, line count, word count
// - Character count (with and without spaces)
// - Average words per line
// - File type specific statistics
var statsCmd = &cobra.Command{
	Use:   "stats [filename]",                         // Command syntax
	Short: "Display detailed statistics about a file", // Brief description
	Long: `Display comprehensive statistics about a file including:

General Statistics:
  - File size in bytes
  - Number of lines
  - Number of words
  - Number of characters (with and without spaces)
  - Average words per line

File Type Specific:
  - CSV: Number of records and fields
  - JSON: Validation status and structure info
  - TXT: Line length analysis

Supported file types: .txt, .csv, .json

Examples:
  optix stats document.txt   # Show statistics for a text file
  optix stats data.csv       # Show statistics for a CSV file
  optix stats config.json    # Show statistics for a JSON file`,

	// Require exactly one argument (the filename)
	Args: cobra.ExactArgs(1),

	// RunE executes the command and can return an error
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := args[0]

		// Step 1: Validate the file
		fileValidator := validator.NewBasicFileValidator()
		validatorStrategy := validator.NewValidatorStrategy(fileValidator)

		if err := validatorStrategy.ValidateFile(filename); err != nil {
			return fmt.Errorf("file validation failed: %v", err)
		}

		// Step 2: Read the file to get content for analysis
		readerStrategy := strategies.NewDefaultFileReaderStrategy()
		content, err := readerStrategy.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("failed to read file for statistics: %v", err)
		}

		// Step 3: Calculate additional statistics
		stats := calculateDetailedStats(content)

		// Step 4: Display comprehensive statistics
		displayStats(filename, content, stats)

		return nil
	},
}

// DetailedStats holds additional calculated statistics.
// This struct extends the basic FileContent with more detailed analysis.
type DetailedStats struct {
	// CharCount is the total number of characters including spaces
	CharCount int

	// CharCountNoSpaces is the total number of characters excluding whitespace
	CharCountNoSpaces int

	// AvgWordsPerLine is the average number of words per line
	AvgWordsPerLine float64

	// LongestLine contains the length of the longest line
	LongestLine int

	// ShortestLine contains the length of the shortest line (excluding empty lines)
	ShortestLine int

	// EmptyLines is the count of completely empty lines
	EmptyLines int
}

// calculateDetailedStats performs additional statistical analysis on file content.
// This function demonstrates Go's approach to data processing and analysis.
func calculateDetailedStats(content *reader.FileContent) *DetailedStats {
	// Initialize our statistics struct
	stats := &DetailedStats{
		CharCount:         len(content.Content),
		CharCountNoSpaces: len(strings.ReplaceAll(content.Content, " ", "")),
		ShortestLine:      -1, // We'll update this with the first non-empty line
	}

	// Calculate average words per line
	// We need to handle the case where there are no lines to avoid division by zero
	if content.LineCount > 0 {
		// float64() converts integers to floating point for division
		stats.AvgWordsPerLine = float64(content.WordCount) / float64(content.LineCount)
	}

	// Analyze each line for length statistics
	for _, line := range content.Lines {
		lineLength := len(line)

		// Check for empty lines
		// strings.TrimSpace removes leading and trailing whitespace
		if strings.TrimSpace(line) == "" {
			stats.EmptyLines++
			continue // Skip empty lines for min/max length calculation
		}

		// Update longest line
		if lineLength > stats.LongestLine {
			stats.LongestLine = lineLength
		}

		// Update shortest line (excluding empty lines)
		if stats.ShortestLine == -1 || lineLength < stats.ShortestLine {
			stats.ShortestLine = lineLength
		}
	}

	// Handle case where all lines are empty
	if stats.ShortestLine == -1 {
		stats.ShortestLine = 0
	}

	return stats
}

// displayStats presents the statistics in a well-formatted, user-friendly way.
// This function demonstrates Go's fmt package capabilities for formatted output.
func displayStats(filename string, content *reader.FileContent, stats *DetailedStats) {
	// Print header with file information
	fmt.Printf("📊 File Statistics for: %s\n", filename)
	fmt.Println("═════════════════════════════════════════════════════")

	// Basic file information
	fmt.Printf("📄 File Type:           %s\n", strings.ToUpper(content.FileType))
	fmt.Printf("📏 File Size:           %d bytes\n", content.Size)

	// Line statistics
	fmt.Println("\n📝 Line Statistics:")
	fmt.Printf("   Total Lines:         %d\n", content.LineCount)
	fmt.Printf("   Empty Lines:         %d\n", stats.EmptyLines)
	fmt.Printf("   Non-empty Lines:     %d\n", content.LineCount-stats.EmptyLines)
	fmt.Printf("   Longest Line:        %d characters\n", stats.LongestLine)
	fmt.Printf("   Shortest Line:       %d characters\n", stats.ShortestLine)

	// Word and character statistics
	fmt.Println("\n🔤 Word & Character Statistics:")
	fmt.Printf("   Total Words:         %d\n", content.WordCount)
	fmt.Printf("   Total Characters:    %d\n", stats.CharCount)
	fmt.Printf("   Chars (no spaces):   %d\n", stats.CharCountNoSpaces)

	// Average calculations with formatting
	// %.2f formats a float to 2 decimal places
	fmt.Printf("   Avg Words/Line:      %.2f\n", stats.AvgWordsPerLine)

	// Calculate and display additional averages
	if content.LineCount > 0 {
		avgCharsPerLine := float64(stats.CharCount) / float64(content.LineCount)
		fmt.Printf("   Avg Chars/Line:      %.2f\n", avgCharsPerLine)
	}

	if content.WordCount > 0 {
		avgCharsPerWord := float64(stats.CharCountNoSpaces) / float64(content.WordCount)
		fmt.Printf("   Avg Chars/Word:      %.2f\n", avgCharsPerWord)
	}

	// File type specific statistics
	displayFileTypeSpecificStats(content)

	// Summary
	fmt.Println("\n✅ Statistics Summary:")
	fmt.Printf("   📊 %d lines, %d words, %d characters in %s file\n",
		content.LineCount, content.WordCount, stats.CharCount, content.FileType)
}

// displayFileTypeSpecificStats shows statistics specific to each file type.
// This demonstrates Go's switch statement and type-specific processing.
func displayFileTypeSpecificStats(content *reader.FileContent) {
	fmt.Printf("\n📋 %s Specific Statistics:\n", strings.ToUpper(content.FileType))

	// Use switch statement to handle different file types
	// Go's switch statements don't fall through by default (unlike C/Java)
	switch content.FileType {
	case "csv":
		displayCSVStats(content)
	case "json":
		displayJSONStats(content)
	case "txt":
		displayTextStats(content)
	default:
		fmt.Printf("   No specific statistics available for %s files\n", content.FileType)
	}
}

// displayCSVStats shows CSV-specific statistics.
func displayCSVStats(content *reader.FileContent) {
	if len(content.Lines) == 0 {
		fmt.Println("   Empty CSV file")
		return
	}

	// Estimate number of fields by looking at the first line
	// In a real application, you might want to parse the CSV more thoroughly
	firstLine := content.Lines[0]
	estimatedFields := len(strings.Split(firstLine, ","))

	fmt.Printf("   Records (rows):      %d\n", content.LineCount)
	fmt.Printf("   Estimated Fields:    %d (based on first row)\n", estimatedFields)
	fmt.Printf("   Estimated Cells:     %d\n", content.LineCount*estimatedFields)
}

// displayJSONStats shows JSON-specific statistics.
func displayJSONStats(content *reader.FileContent) {
	// Count braces and brackets for structure analysis
	openBraces := strings.Count(content.Content, "{")
	closeBraces := strings.Count(content.Content, "}")
	openBrackets := strings.Count(content.Content, "[")
	closeBrackets := strings.Count(content.Content, "]")

	fmt.Printf("   Objects ({}):        %d pairs\n", openBraces)
	fmt.Printf("   Arrays ([]):         %d pairs\n", openBrackets)
	fmt.Printf("   Bracket Balance:     %s\n", getBracketBalanceStatus(openBraces, closeBraces, openBrackets, closeBrackets))

	// Count commas as a rough estimate of JSON elements
	commas := strings.Count(content.Content, ",")
	fmt.Printf("   Estimated Elements:  %d (based on commas)\n", commas+1)
}

// displayTextStats shows text-specific statistics.
func displayTextStats(content *reader.FileContent) {
	// Count sentences (rough estimate based on sentence-ending punctuation)
	sentences := strings.Count(content.Content, ".") +
		strings.Count(content.Content, "!") +
		strings.Count(content.Content, "?")

	// Count paragraphs (double newlines)
	paragraphs := strings.Count(content.Content, "\n\n") + 1

	fmt.Printf("   Estimated Sentences: %d\n", sentences)
	fmt.Printf("   Estimated Paragraphs: %d\n", paragraphs)

	if sentences > 0 {
		avgWordsPerSentence := float64(content.WordCount) / float64(sentences)
		fmt.Printf("   Avg Words/Sentence:  %.2f\n", avgWordsPerSentence)
	}
}

// getBracketBalanceStatus checks if JSON brackets are properly balanced.
// This is a helper function that demonstrates Go's approach to utility functions.
func getBracketBalanceStatus(openB, closeB, openBr, closeBr int) string {
	if openB == closeB && openBr == closeBr {
		return "✅ Balanced"
	}
	return "❌ Unbalanced"
}

// init registers the stats command with the root command.
func init() {
	cmd.RootCmd.AddCommand(statsCmd)
}
