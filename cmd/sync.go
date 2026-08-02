package cmd

import (
	"fmt"
	"strings"

	"github.com/getenvguard/cli/pkg/api"
	"github.com/getenvguard/cli/pkg/config"
	"github.com/getenvguard/cli/pkg/k8s"
	"github.com/spf13/cobra"
)

var (
	flagNamespace string
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Directly sync secrets from EnvGuard to Kubernetes Cluster via kubectl",
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

		k8sYAML, err := client.ExportSecretsRaw(targetProject.ID, targetEnv.ID, "k8s")
		if err != nil {
			return err
		}

		out, err := k8s.ApplyYAML(k8sYAML, flagNamespace)
		if err != nil {
			return err
		}

		fmt.Println("✔ Successfully synced secrets to Kubernetes Cluster!")
		if strings.TrimSpace(out) != "" {
			fmt.Printf("  %s\n", strings.TrimSpace(out))
		}

		return nil
	},
}

func init() {
	syncCmd.Flags().StringVarP(&flagProject, "project", "p", "", "Target project name, slug, or ID")
	syncCmd.Flags().StringVarP(&flagEnv, "env", "e", "", "Target environment (development, staging, production)")
	syncCmd.Flags().StringVarP(&flagNamespace, "namespace", "n", "", "Target Kubernetes namespace (default: default)")
	rootCmd.AddCommand(syncCmd)
}
