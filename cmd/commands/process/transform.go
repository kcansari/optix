// Package optix contains the CLI commands for the Optix file processor.
// This file implements the 'transform' command for text transformations.
package process

import (
	"fmt"
	"strings"

	"github.com/kcansari/optix/cmd"
	"github.com/kcansari/optix/internal/processor"
	"github.com/kcansari/optix/internal/reader"
	"github.com/kcansari/optix/internal/validator"
	"github.com/spf13/cobra"
)

// transformCmd represents the transform command.
// This command performs text transformations like case conversion and whitespace cleanup.
var transformCmd = &cobra.Command{
	Use:   "transform",
	Short: "Transform text content (case conversion, whitespace cleanup)",
	Long: `Transform text content with various operations.

The transform command supports:
  - Case conversion (upper, lower, title)
  - Whitespace cleanup (trim)
  - Output to file or overwrite original
  - Dry run mode to preview changes

Available transformations:
  - upper: Convert all text to uppercase
  - lower: Convert all text to lowercase
  - title: Convert text to title case
  - trim:  Remove leading/trailing whitespace from each line

Examples:
  optix transform --type upper --file document.txt
  optix transform --type lower --file README.md --output readme.md
  optix transform --type trim --file data.csv --dry-run
  optix transform --type title --file notes.txt`,

	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flag values
		transformType, _ := cmd.Flags().GetString("type")
		fileName, _ := cmd.Flags().GetString("file")
		outputFile, _ := cmd.Flags().GetString("output")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		// Validate required flags
		if transformType == "" {
			return fmt.Errorf("transformation type is required (use --type flag)")
		}
		if fileName == "" {
			return fmt.Errorf("file is required (use --file flag)")
		}

		// Validate transformation type
		validTypes := []string{"upper", "lower", "title", "trim"}
		isValid := false
		for _, validType := range validTypes {
			if strings.ToLower(transformType) == validType {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("invalid transformation type '%s'. Valid types: %s",
				transformType, strings.Join(validTypes, ", "))
		}

		// Create processor strategy
		processorStrategy := processor.NewTextProcessorStrategy()
		readerStrategy := reader.NewFileReaderStrategy()
		validatorStrategy := validator.NewValidatorStrategy(validator.NewBasicFileValidator())

		// Validate file
		if err := validatorStrategy.ValidateFile(fileName); err != nil {
			return fmt.Errorf("file validation failed: %v", err)
		}

		// Read file content
		content, err := readerStrategy.ReadFile(fileName)
		if err != nil {
			return fmt.Errorf("failed to read file: %v", err)
		}

		// Prepare processing options
		options := processor.ProcessOptions{
			TransformType: strings.ToLower(transformType),
			FileName:      fileName,
			OutputFile:    outputFile,
			DryRun:        dryRun,
		}

		// Display operation info
		fmt.Printf("🔄 Transform Operation\n")
		fmt.Printf("📄 File: %s\n", fileName)
		fmt.Printf("🔧 Transform Type: %s\n", transformType)
		if dryRun {
			fmt.Printf("🧪 Dry Run: Enabled (no changes will be made)\n")
		}
		if outputFile != "" {
			fmt.Printf("📤 Output File: %s\n", outputFile)
		} else {
			fmt.Printf("📤 Output: Overwrite original file\n")
		}
		fmt.Println("─────────────────────────────────────────────────────")

		// Process the file
		result, err := processorStrategy.ProcessText("transform", content, options)
		if err != nil {
			return fmt.Errorf("transform operation failed: %v", err)
		}

		// Display preview for dry run
		if dryRun {
			fmt.Printf("🧪 Dry Run Preview:\n")
			fmt.Println("─────────────────────────────────────────────────────")

			// Show first few lines of transformed content
			lines := strings.Split(result.ModifiedContent, "\n")
			previewLines := 10
			if len(lines) < previewLines {
				previewLines = len(lines)
			}

			for i := 0; i < previewLines; i++ {
				if lines[i] != "" {
					fmt.Printf("%3d: %s\n", i+1, lines[i])
				}
			}

			if len(lines) > previewLines {
				fmt.Printf("... and %d more lines\n", len(lines)-previewLines)
			}

			fmt.Println("─────────────────────────────────────────────────────")
		}

		// Display results
		fmt.Printf("✅ Transform operation completed successfully\n")
		fmt.Printf("📊 Results:\n")
		fmt.Printf("   📝 Lines processed: %d\n", result.LinesProcessed)
		fmt.Printf("   ⏱️  Execution time: %v\n", result.ExecutionTime)

		if dryRun {
			fmt.Printf("   🧪 Dry run completed - no changes were made\n")
			fmt.Printf("   ℹ️  Run without --dry-run to apply transformation\n")
		} else {
			outputTarget := fileName
			if outputFile != "" {
				outputTarget = outputFile
			}
			fmt.Printf("   📄 Transformed file: %s\n", outputTarget)

			// Show transformation summary
			switch strings.ToLower(transformType) {
			case "upper":
				fmt.Printf("   🔤 All text converted to UPPERCASE\n")
			case "lower":
				fmt.Printf("   🔤 All text converted to lowercase\n")
			case "title":
				fmt.Printf("   🔤 All text converted to Title Case\n")
			case "trim":
				fmt.Printf("   ✂️  Whitespace trimmed from all lines\n")
			}
		}

		return nil
	},
}

// init function registers the transform command and its flags.
func init() {
	cmd.RootCmd.AddCommand(transformCmd)

	// Add flags for transform options
	transformCmd.Flags().StringP("type", "t", "", "Transformation type: upper, lower, title, trim (required)")
	transformCmd.Flags().String("file", "", "File to transform (required)")
	transformCmd.Flags().StringP("output", "o", "", "Output file (default: overwrite input file)")
	transformCmd.Flags().Bool("dry-run", false, "Preview transformation without modifying files")

	// Mark required flags
	transformCmd.MarkFlagRequired("type")
	transformCmd.MarkFlagRequired("file")
}
