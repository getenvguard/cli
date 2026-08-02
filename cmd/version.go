package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version = "v3.0.0"
	Commit  = "none"
	Date    = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print EnvGuard CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("🔐 EnvGuard CLI %s (Go Native, commit: %s, built: %s)\n", Version, Commit, Date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
