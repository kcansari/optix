// Package optix contains the CLI commands for the Optix file processor.
// This file implements the 'replace' command for text replacement operations.
package process

import (
	"fmt"

	"github.com/kcansari/optix/cmd"
	"github.com/kcansari/optix/internal/processor"
	"github.com/kcansari/optix/internal/reader"
	"github.com/kcansari/optix/internal/validator"
	"github.com/spf13/cobra"
)

// replaceCmd represents the replace command.
// This command performs search and replace operations with backup support.
var replaceCmd = &cobra.Command{
	Use:   "replace",
	Short: "Search and replace text in files",
	Long: `Search and replace text in files with backup support.

The replace command supports:
  - Regular expressions and literal text replacement
  - Automatic backup creation before modification
  - Dry run mode to preview changes
  - Case-sensitive and case-insensitive replacement
  - Whole word matching

Examples:
  optix replace --find "old_url" --replace "new_url" --file config.txt
  optix replace --find "user\d+" --replace "customer$0" --regex --file data.txt
  optix replace --find "TODO" --replace "DONE" --file notes.txt --backup
  optix replace --find "debug" --replace "info" --file app.log --dry-run`,

	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flag values
		findPattern, _ := cmd.Flags().GetString("find")
		replaceWith, _ := cmd.Flags().GetString("replace")
		fileName, _ := cmd.Flags().GetString("file")
		regexMode, _ := cmd.Flags().GetBool("regex")
		caseSensitive, _ := cmd.Flags().GetBool("case-sensitive")
		wholeWord, _ := cmd.Flags().GetBool("whole-word")
		createBackup, _ := cmd.Flags().GetBool("backup")
		backupDir, _ := cmd.Flags().GetString("backup-dir")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		outputFile, _ := cmd.Flags().GetString("output")

		// Validate required flags
		if findPattern == "" {
			return fmt.Errorf("find pattern is required (use --find flag)")
		}
		if replaceWith == "" {
			return fmt.Errorf("replacement text is required (use --replace flag)")
		}
		if fileName == "" {
			return fmt.Errorf("file is required (use --file flag)")
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
			Pattern:       findPattern,
			ReplaceWith:   replaceWith,
			RegexMode:     regexMode,
			CaseSensitive: caseSensitive,
			WholeWord:     wholeWord,
			CreateBackup:  createBackup,
			BackupDir:     backupDir,
			DryRun:        dryRun,
			FileName:      fileName,
			OutputFile:    outputFile,
		}

		// Display operation info
		fmt.Printf("🔄 Replace Operation\n")
		fmt.Printf("📄 File: %s\n", fileName)
		fmt.Printf("🔍 Find: %s\n", findPattern)
		fmt.Printf("🔄 Replace: %s\n", replaceWith)
		if regexMode {
			fmt.Printf("🔧 Mode: Regular Expression\n")
		} else {
			fmt.Printf("🔧 Mode: Literal Text\n")
		}
		fmt.Printf("📊 Case Sensitive: %t\n", caseSensitive)
		if wholeWord {
			fmt.Printf("🔤 Whole Word: %t\n", wholeWord)
		}
		if createBackup {
			fmt.Printf("💾 Backup: Enabled\n")
			if backupDir != "" {
				fmt.Printf("📁 Backup Directory: %s\n", backupDir)
			}
		}
		if dryRun {
			fmt.Printf("🧪 Dry Run: Enabled (no changes will be made)\n")
		}
		if outputFile != "" {
			fmt.Printf("📤 Output File: %s\n", outputFile)
		}
		fmt.Println("─────────────────────────────────────────────────────")

		// Process the file
		result, err := processorStrategy.ProcessText("replace", content, options)
		if err != nil {
			return fmt.Errorf("replace operation failed: %v", err)
		}

		// Display results
		fmt.Printf("✅ Replace operation completed successfully\n")
		fmt.Printf("📊 Results:\n")
		fmt.Printf("   🎯 Matches found: %d\n", result.MatchesFound)
		fmt.Printf("   📝 Lines processed: %d\n", result.LinesProcessed)
		fmt.Printf("   ⏱️  Execution time: %v\n", result.ExecutionTime)

		if result.BackupPath != "" {
			fmt.Printf("   💾 Backup created: %s\n", result.BackupPath)
		}

		if dryRun {
			fmt.Printf("   🧪 Dry run completed - no changes were made\n")
			if result.MatchesFound > 0 {
				fmt.Printf("   ℹ️  Run without --dry-run to apply changes\n")
			}
		} else {
			outputTarget := fileName
			if outputFile != "" {
				outputTarget = outputFile
			}
			fmt.Printf("   📄 Modified file: %s\n", outputTarget)
		}

		if result.MatchesFound == 0 {
			fmt.Printf("   ℹ️  No matches found for pattern '%s'\n", findPattern)
		}

		return nil
	},
}

// init function registers the replace command and its flags.
func init() {
	cmd.RootCmd.AddCommand(replaceCmd)

	// Add flags for replace options
	replaceCmd.Flags().StringP("find", "f", "", "Text pattern to find (required)")
	replaceCmd.Flags().StringP("replace", "r", "", "Replacement text (required)")
	replaceCmd.Flags().String("file", "", "File to process (required)")
	replaceCmd.Flags().Bool("regex", false, "Use regular expression mode")
	replaceCmd.Flags().BoolP("case-sensitive", "c", false, "Case sensitive replacement")
	replaceCmd.Flags().BoolP("whole-word", "w", false, "Match whole words only")
	replaceCmd.Flags().BoolP("backup", "b", false, "Create backup before modification")
	replaceCmd.Flags().String("backup-dir", "", "Directory for backup files (default: same as original)")
	replaceCmd.Flags().Bool("dry-run", false, "Preview changes without modifying files")
	replaceCmd.Flags().StringP("output", "o", "", "Output file (default: overwrite input file)")

	// Mark required flags
	replaceCmd.MarkFlagRequired("find")
	replaceCmd.MarkFlagRequired("replace")
	replaceCmd.MarkFlagRequired("file")
}
