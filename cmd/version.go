package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version information (will be set during build)
var (
	Version   = "0.1.0"
	GitCommit = "development"
	BuildDate = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information of go-env-cli",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("go-env-cli %s\n", Version)
		fmt.Printf("Commit: %s\n", GitCommit)
		fmt.Printf("Built: %s\n", BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
