// Package optix contains the CLI commands for the Optix file processor.
// This file implements the 'filter' command for extracting matching lines.
package process

import (
	"fmt"

	"github.com/kcansari/optix/cmd"
	"github.com/kcansari/optix/internal/processor"
	"github.com/kcansari/optix/internal/reader"
	"github.com/kcansari/optix/internal/validator"
	"github.com/spf13/cobra"
)

// filterCmd represents the filter command.
// This command extracts lines matching specific criteria.
var filterCmd = &cobra.Command{
	Use:   "filter",
	Short: "Filter lines from files based on patterns",
	Long: `Filter and extract lines from files based on text patterns.

The filter command supports:
  - Regular expressions and literal text matching
  - Inverted matching (lines that don't match)
  - Extract only matching parts or entire lines
  - Case-sensitive and case-insensitive filtering
  - Output to file or console

Examples:
  optix filter --contains "WARNING" --input app.log --output warnings.log
  optix filter --pattern "error\d+" --regex --input system.log
  optix filter --contains "TODO" --invert --input code.go
  optix filter --pattern "user" --only-matching --input data.txt`,

	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flag values
		pattern, _ := cmd.Flags().GetString("pattern")
		contains, _ := cmd.Flags().GetString("contains")
		inputFile, _ := cmd.Flags().GetString("input")
		outputFile, _ := cmd.Flags().GetString("output")
		regexMode, _ := cmd.Flags().GetBool("regex")
		caseSensitive, _ := cmd.Flags().GetBool("case-sensitive")
		invertMatch, _ := cmd.Flags().GetBool("invert")
		onlyMatching, _ := cmd.Flags().GetBool("only-matching")

		// Determine the search pattern
		searchPattern := pattern
		if contains != "" {
			if pattern != "" {
				return fmt.Errorf("cannot use both --pattern and --contains flags")
			}
			searchPattern = contains
		}

		// Validate required flags
		if searchPattern == "" {
			return fmt.Errorf("search criteria is required (use --pattern or --contains flag)")
		}
		if inputFile == "" {
			return fmt.Errorf("input file is required (use --input flag)")
		}

		// Create processor strategy
		processorStrategy := processor.NewTextProcessorStrategy()
		readerStrategy := reader.NewFileReaderStrategy()
		validatorStrategy := validator.NewValidatorStrategy(validator.NewBasicFileValidator())

		// Validate file
		if err := validatorStrategy.ValidateFile(inputFile); err != nil {
			return fmt.Errorf("file validation failed: %v", err)
		}

		// Read file content
		content, err := readerStrategy.ReadFile(inputFile)
		if err != nil {
			return fmt.Errorf("failed to read file: %v", err)
		}

		// Prepare processing options
		options := processor.ProcessOptions{
			Pattern:       searchPattern,
			RegexMode:     regexMode || (pattern != ""), // Use regex mode if --pattern flag was used
			CaseSensitive: caseSensitive,
			InvertMatch:   invertMatch,
			OnlyMatching:  onlyMatching,
			FileName:      inputFile,
			OutputFile:    outputFile,
		}

		// Display operation info
		fmt.Printf("📋 Filter Operation\n")
		fmt.Printf("📄 Input: %s\n", inputFile)
		fmt.Printf("🔍 Pattern: %s\n", searchPattern)
		if regexMode || (pattern != "") {
			fmt.Printf("🔧 Mode: Regular Expression\n")
		} else {
			fmt.Printf("🔧 Mode: Literal Text (contains)\n")
		}
		fmt.Printf("📊 Case Sensitive: %t\n", caseSensitive)
		if invertMatch {
			fmt.Printf("🔄 Invert Match: %t (lines that DON'T match)\n", invertMatch)
		}
		if onlyMatching {
			fmt.Printf("✂️  Only Matching: %t (extract matching parts only)\n", onlyMatching)
		}
		if outputFile != "" {
			fmt.Printf("📤 Output: %s\n", outputFile)
		} else {
			fmt.Printf("📤 Output: Console\n")
		}
		fmt.Println("─────────────────────────────────────────────────────")

		// Process the file
		result, err := processorStrategy.ProcessText("filter", content, options)
		if err != nil {
			return fmt.Errorf("filter operation failed: %v", err)
		}

		// Display filtered content if no output file specified
		if outputFile == "" && result.ModifiedContent != "" {
			fmt.Printf("📋 Filtered Content:\n")
			fmt.Println("─────────────────────────────────────────────────────")
			fmt.Print(result.ModifiedContent)
			fmt.Println("─────────────────────────────────────────────────────")
		}

		// Display results summary
		fmt.Printf("✅ Filter operation completed successfully\n")
		fmt.Printf("📊 Results:\n")
		fmt.Printf("   🎯 Matching lines: %d\n", result.MatchesFound)
		fmt.Printf("   📝 Total lines processed: %d\n", result.LinesProcessed)
		fmt.Printf("   ⏱️  Execution time: %v\n", result.ExecutionTime)

		if outputFile != "" {
			fmt.Printf("   📄 Output written to: %s\n", outputFile)
		}

		if result.MatchesFound == 0 {
			if invertMatch {
				fmt.Printf("   ℹ️  All lines matched the pattern '%s'\n", searchPattern)
			} else {
				fmt.Printf("   ℹ️  No lines matched the pattern '%s'\n", searchPattern)
			}
		}

		return nil
	},
}

// init function registers the filter command and its flags.
func init() {
	cmd.RootCmd.AddCommand(filterCmd)

	// Add flags for filter options
	filterCmd.Flags().StringP("pattern", "p", "", "Regular expression pattern to match")
	filterCmd.Flags().String("contains", "", "Literal text that lines must contain")
	filterCmd.Flags().StringP("input", "i", "", "Input file to filter (required)")
	filterCmd.Flags().StringP("output", "o", "", "Output file for filtered results (optional)")
	filterCmd.Flags().BoolP("regex", "r", false, "Use regular expression mode (auto-enabled with --pattern)")
	filterCmd.Flags().BoolP("case-sensitive", "c", false, "Case sensitive filtering")
	filterCmd.Flags().BoolP("invert", "v", false, "Invert match (select lines that DON'T match)")
	filterCmd.Flags().Bool("only-matching", false, "Output only the matching parts of lines")

	// Mark required flags
	filterCmd.MarkFlagRequired("input")
}
