// Package optix contains the CLI commands for the Optix file processor.
// This file implements the 'show' command that displays file contents.
package optix

import (
	"fmt" // Package for formatted I/O operations

	"github.com/kcansari/optix/internal/reader"    // Our file reader package
	"github.com/kcansari/optix/internal/validator" // Our file validator package
	"github.com/spf13/cobra"                       // CLI framework
)

// showCmd represents the show command.
// This command displays the contents of a file using the appropriate reader strategy.
// The command follows the pattern: optix show <filename>
var showCmd = &cobra.Command{
	Use:   "show [filename]",                // Use describes the command syntax
	Short: "Display the contents of a file", // Short description for help
	Long: `Display the contents of a file with support for multiple formats.

Supported file types:
  - .txt  (Text files)
  - .csv  (Comma-separated values)
  - .json (JSON files)

Examples:
  optix show myfile.txt     # Display a text file
  optix show data.csv       # Display a CSV file  
  optix show config.json    # Display a JSON file`,

	// Args validates the number of command line arguments
	// cobra.ExactArgs(1) means this command requires exactly 1 argument
	Args: cobra.ExactArgs(1),

	// RunE is the function that executes when the command is called
	// The 'E' suffix means it can return an error
	RunE: func(cmd *cobra.Command, args []string) error {
		// args[0] contains the filename passed to the command
		filename := args[0]

		// Step 1: Validate the file exists and is readable
		// Create a file validator using our strategy pattern
		fileValidator := validator.NewBasicFileValidator()
		validatorStrategy := validator.NewValidatorStrategy(fileValidator)

		// Validate the file before trying to read it
		if err := validatorStrategy.ValidateFile(filename); err != nil {
			// If validation fails, return a user-friendly error message
			return fmt.Errorf("file validation failed: %v", err)
		}

		// Step 2: Read the file using our improved reader strategy
		// Create a reader strategy that can handle multiple file types
		readerStrategy := reader.NewFileReaderStrategy()

		// Read the file - the strategy will automatically choose the right reader
		content, err := readerStrategy.ReadFile(filename)
		if err != nil {
			// If reading fails, return an error with context
			return fmt.Errorf("failed to read file: %v", err)
		}

		// Step 3: Display the file information and contents
		// Print a header with file information
		fmt.Printf("📄 File: %s\n", filename)
		fmt.Printf("📊 Type: %s\n", content.FileType)
		fmt.Printf("📏 Size: %d bytes\n", content.Size)
		fmt.Printf("📝 Lines: %d\n", content.LineCount)
		fmt.Printf("🔤 Words: %d\n", content.WordCount)
		fmt.Println("📖 Content:")
		fmt.Println("─────────────────────────────────────────────────────")

		// Print the actual file content
		fmt.Print(content.Content)

		// Add a separator line at the end for better readability
		fmt.Println("─────────────────────────────────────────────────────")
		fmt.Printf("✅ Successfully displayed %s (%s file)\n", filename, content.FileType)

		// Return nil to indicate success
		return nil
	},
}

// init function is called automatically when the package is imported.
// We use it to register the show command with the root command.
func init() {
	// Add the show command to the root command
	// This makes it available as 'optix show'
	rootCmd.AddCommand(showCmd)
}
