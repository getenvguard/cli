package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/getenvguard/cli/pkg/api"
	"github.com/getenvguard/cli/pkg/config"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run -- <command> [args...]",
	Short: "Inject secrets in-memory and execute process (Zero disk footprint!)",
	Args:  cobra.MinimumNArgs(1),
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

		fmt.Printf("🔐 Fetching in-memory secrets from EnvGuard (%s)...\n", client.APIHost)
		secrets, err := client.ExportSecretsJSON(targetProject.ID, targetEnv.ID)
		if err != nil {
			return err
		}

		cmdToRun := args[0]
		cmdArgs := args[1:]

		fmt.Printf("🚀 Injecting %d secrets in-memory & spawning: %s %s\n\n", len(secrets), cmdToRun, strings.Join(cmdArgs, " "))

		child := exec.Command(cmdToRun, cmdArgs...)
		child.Stdout = os.Stdout
		child.Stderr = os.Stderr
		child.Stdin = os.Stdin

		// Build injected environment
		envMap := make(map[string]string)
		for _, envStr := range os.Environ() {
			parts := strings.SplitN(envStr, "=", 2)
			if len(parts) == 2 {
				envMap[parts[0]] = parts[1]
			}
		}
		for k, v := range secrets {
			envMap[k] = v
		}

		var newEnv []string
		for k, v := range envMap {
			newEnv = append(newEnv, fmt.Sprintf("%s=%s", k, v))
		}
		child.Env = newEnv

		if err := child.Run(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				os.Exit(exitErr.ExitCode())
			}
			return err
		}

		return nil
	},
}

func init() {
	runCmd.Flags().StringVarP(&flagProject, "project", "p", "", "Target project name, slug, or ID")
	runCmd.Flags().StringVarP(&flagEnv, "env", "e", "", "Target environment (development, staging, production)")
	rootCmd.AddCommand(runCmd)
}
