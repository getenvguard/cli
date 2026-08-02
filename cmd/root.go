package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	flagAPIHost string
	flagToken   string

	rootCmd = &cobra.Command{
		Use:   "envg",
		Short: "🔐 EnvGuard CLI — Declarative Secret Management for Modern Teams",
		Long: `EnvGuard CLI (envg) is a native, zero-dependency command line tool for 
EnvGuard (https://getenvguard.com).

Seamlessly pull secrets, sync directly to Kubernetes clusters, or inject 
encrypted secrets in-memory into process runtimes with zero disk footprint.`,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagAPIHost, "api-host", "https://getenvguard.com", "EnvGuard API host URL")
	rootCmd.PersistentFlags().StringVar(&flagToken, "token", "", "Personal Access Token (PAT) for authentication")
}
