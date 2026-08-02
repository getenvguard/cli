package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/getenvguard/cli/pkg/api"
	"github.com/getenvguard/cli/pkg/config"
	"github.com/spf13/cobra"
)

var (
	flagProject string
	flagEnv     string
	flagFormat  string
	flagOut     string
)

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull environment variables to local workspace (.env, k8s, json)",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.GetEffectiveToken(flagToken)
		if token == "" {
			return fmt.Errorf("unauthorized. Please run 'envg login' or pass --token")
		}

		client := api.NewClient(flagAPIHost, token)
		projects, err := client.FetchProjects()
		if err != nil {
			return err
		}

		if len(projects) == 0 {
			return fmt.Errorf("no projects found in your workspace")
		}

		targetProject := projects[0]
		if flagProject != "" {
			for _, p := range projects {
				if strings.EqualFold(p.Name, flagProject) || strings.EqualFold(p.Slug, flagProject) || p.ID == flagProject {
					targetProject = p
					break
				}
			}
		}

		if len(targetProject.Environments) == 0 {
			return fmt.Errorf("no environments found for project '%s'", targetProject.Name)
		}

		targetEnv := targetProject.Environments[0]
		for _, e := range targetProject.Environments {
			if strings.EqualFold(e.Name, "development") || strings.EqualFold(e.Slug, "development") {
				targetEnv = e
				break
			}
		}
		if flagEnv != "" {
			for _, e := range targetProject.Environments {
				if strings.EqualFold(e.Name, flagEnv) || strings.EqualFold(e.Slug, flagEnv) || e.ID == flagEnv {
					targetEnv = e
					break
				}
			}
		}

		fmt.Printf("⬇️  Fetching secrets for project '%s' [%s]...\n", targetProject.Name, targetEnv.Name)

		rawContent, err := client.ExportSecretsRaw(targetProject.ID, targetEnv.ID, flagFormat)
		if err != nil {
			return err
		}

		if flagOut != "" {
			if err := os.WriteFile(flagOut, []byte(rawContent), 0600); err != nil {
				return fmt.Errorf("failed to write output file: %w", err)
			}
			fmt.Printf("✔ Successfully exported secrets to %s\n", flagOut)
		} else {
			fmt.Println(rawContent)
		}

		return nil
	},
}

func init() {
	pullCmd.Flags().StringVarP(&flagProject, "project", "p", "", "Target project name, slug, or ID")
	pullCmd.Flags().StringVarP(&flagEnv, "env", "e", "", "Target environment (development, staging, production)")
	pullCmd.Flags().StringVarP(&flagFormat, "format", "f", "env", "Output format: env, k8s, json")
	pullCmd.Flags().StringVarP(&flagOut, "out", "o", "", "Output filename (default: stdout)")
	rootCmd.AddCommand(pullCmd)
}
