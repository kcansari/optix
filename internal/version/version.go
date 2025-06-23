package version

import "fmt"

var (
	Version   = "1.0.0"
	BuildDate = "2024-06-23"
	Commit    = "dev"
)

func PrintVersion() {
	fmt.Printf("Optix File Processor\n")
	fmt.Printf("Version: %s\n", Version)
	fmt.Printf("Build Date: %s\n", BuildDate)
	fmt.Printf("Commit: %s\n", Commit)
}
