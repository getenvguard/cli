package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Environment struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type Project struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Slug         string        `json:"slug"`
	Environments []Environment `json:"environments"`
}

type ProjectsResponse struct {
	Projects []Project `json:"projects"`
}

type Client struct {
	APIHost    string
	Token      string
	HTTPClient *http.Client
}

func NewClient(apiHost, token string) *Client {
	if apiHost == "" {
		apiHost = "https://getenvguard.com"
	}
	apiHost = strings.TrimSuffix(apiHost, "/")

	return &Client{
		APIHost: apiHost,
		Token:   token,
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *Client) FetchProjects() ([]Project, error) {
	url := fmt.Sprintf("%s/api/projects", c.APIHost)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
		req.Header.Set("X-Api-Key", c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("unauthorized (401). Please run 'envg login' or pass --token")
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned HTTP %d: %s", resp.StatusCode, string(body))
	}

	var res ProjectsResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("failed to decode projects JSON: %w", err)
	}

	return res.Projects, nil
}

func (c *Client) ExportSecretsRaw(projectID, envID, format string) (string, error) {
	url := fmt.Sprintf("%s/api/projects/%s/environments/%s/export?format=%s", c.APIHost, projectID, envID, format)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
		req.Header.Set("X-Api-Key", c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("export API returned HTTP %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (c *Client) ExportSecretsJSON(projectID, envID string) (map[string]string, error) {
	raw, err := c.ExportSecretsRaw(projectID, envID, "json")
	if err != nil {
		return nil, err
	}

	var res struct {
		Secrets map[string]string `json:"secrets"`
	}
	if err := json.Unmarshal([]byte(raw), &res); err != nil {
		return nil, fmt.Errorf("failed to parse JSON secrets: %w", err)
	}

	return res.Secrets, nil
}
