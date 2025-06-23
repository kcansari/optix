package optix

import (
	"github.com/kcansari/optix/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information",
	Long:  "Display version information including build date and commit hash",
	Run: func(cmd *cobra.Command, args []string) {
		version.PrintVersion()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
