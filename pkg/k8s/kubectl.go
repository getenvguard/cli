package k8s

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func ApplyYAML(yamlContent, namespace string) (string, error) {
	if namespace != "" {
		yamlContent = strings.ReplaceAll(yamlContent, "namespace: default", "namespace: "+namespace)
	}

	cmd := exec.Command("kubectl", "apply", "-f", "-")
	cmd.Stdin = strings.NewReader(yamlContent)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("kubectl apply failed: %w\n%s", err, stderr.String())
	}

	return stdout.String(), nil
}
