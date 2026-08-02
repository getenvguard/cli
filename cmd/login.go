package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/getenvguard/cli/pkg/config"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to EnvGuard via Personal Access Token (PAT)",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🔐 EnvGuard CLI Authentication")
		
		token := flagToken
		if token == "" {
			fmt.Print("Enter your Personal Access Token (PAT) [eg_pat_...]: ")
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			token = strings.TrimSpace(input)
		}

		if token == "" {
			return fmt.Errorf("token cannot be empty")
		}

		if err := config.SaveCredentials(token, flagAPIHost); err != nil {
			return fmt.Errorf("failed to save credentials: %w", err)
		}

		path, _ := config.GetCredentialsFilePath()
		fmt.Printf("✔ Successfully authenticated & credentials saved to %s\n", path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
