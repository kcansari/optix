package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "optix",
	Short: "A powerful file processing CLI tool",
	Long: `Optix is a Go-based file processing CLI tool designed to handle text, CSV, and JSON file operations
with advanced features like batch processing, concurrency, and data transformation.`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
