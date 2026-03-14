package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

const (
	EnvServerURL    = "ENTERPRISE_AI_PLATFORM_SERVER_URL"
	EnvAPIKey = "ENTERPRISE_AI_PLATFORM_API_KEY"
	EnvAgentID    = "ENTERPRISE_AI_PLATFORM_PERSONA_ID"
)

// EnterpriseAIPlatformCliConfig holds the CLI configuration.
type EnterpriseAIPlatformCliConfig struct {
	ServerURL        string `json:"server_url"`
	APIKey           string `json:"api_key"`
	DefaultAgentID int    `json:"default_persona_id"`
}

// DefaultConfig returns a config with default values.
func DefaultConfig() EnterpriseAIPlatformCliConfig {
	return EnterpriseAIPlatformCliConfig{
		ServerURL:        "https://cloud.your-platform.com",
		APIKey:           "",
		DefaultAgentID: 0,
	}
}

// IsConfigured returns true if the config has an API key.
func (c OnyxCliConfig) IsConfigured() bool {
	return c.APIKey != ""
}

// configDir returns ~/.config/enterprise-ai-platform-cli
func configDir() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "enterprise-ai-platform-cli")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", ".config", "enterprise-ai-platform-cli")
	}
	return filepath.Join(home, ".config", "enterprise-ai-platform-cli")
}

// ConfigFilePath returns the full path to the config file.
func ConfigFilePath() string {
	return filepath.Join(configDir(), "config.json")
}

// ConfigExists checks if the config file exists on disk.
func ConfigExists() bool {
	_, err := os.Stat(ConfigFilePath())
	return err == nil
}

// Load reads config from file and applies environment variable overrides.
func Load() OnyxCliConfig {
	cfg := DefaultConfig()

	data, err := os.ReadFile(ConfigFilePath())
	if err == nil {
		if jsonErr := json.Unmarshal(data, &cfg); jsonErr != nil {
			fmt.Fprintf(os.Stderr, "warning: config file %s is malformed: %v (using defaults)\n", ConfigFilePath(), jsonErr)
		}
	}

	// Environment overrides
	if v := os.Getenv(EnvServerURL); v != "" {
		cfg.ServerURL = v
	}
	if v := os.Getenv(EnvAPIKey); v != "" {
		cfg.APIKey = v
	}
	if v := os.Getenv(EnvAgentID); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			cfg.DefaultAgentID = id
		}
	}

	return cfg
}

// Save writes the config to disk, creating parent directories if needed.
func Save(cfg OnyxCliConfig) error {
	dir := configDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(ConfigFilePath(), data, 0o600)
}
