package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Credentials struct {
	Token  string `json:"token"`
	Name   string `json:"name,omitempty"`
	APIUrl string `json:"apiUrl,omitempty"`
}

func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(homeDir, ".envguard")
	return configDir, nil
}

func GetCredentialsFilePath() (string, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "credentials.json"), nil
}

func LoadCredentials() (*Credentials, error) {
	path, err := GetCredentialsFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, err
	}

	return &creds, nil
}

func SaveCredentials(token, apiUrl string) error {
	dir, err := GetConfigDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	path, err := GetCredentialsFilePath()
	if err != nil {
		return err
	}

	creds := Credentials{
		Token:  token,
		Name:   "Developer",
		APIUrl: apiUrl,
	}

	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

func GetEffectiveToken(flagToken string) string {
	if flagToken != "" {
		return flagToken
	}
	if envToken := os.Getenv("ENVGUARD_TOKEN"); envToken != "" {
		return envToken
	}
	if envKey := os.Getenv("ENVGUARD_API_KEY"); envKey != "" {
		return envKey
	}
	creds, err := LoadCredentials()
	if err == nil && creds != nil && creds.Token != "" {
		return creds.Token
	}
	return ""
}
