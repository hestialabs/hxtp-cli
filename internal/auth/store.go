package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zalando/go-keyring"
)

const (
	Service      = "io.hxtp.cli"
	DefaultToken = "main-token"
	ConfigDir    = ".hxtp"
	ConfigFile   = "config.json"
)

type Config struct {
	ApiUrl     string `json:"api_url"`
	TenantId   string `json:"tenant_id"`
	ClientId   string `json:"client_id"`
	DeviceId   string `json:"device_id"`
	Secret     string `json:"secret"`
	LastLogin  string `json:"last_login"`
}

// GetConfigPath returns the cross-platform path to the config JSON.
func GetConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ConfigDir, ConfigFile)
}

// SaveToken secures the API token in the system keychain.
func SaveToken(token string) error {
	return keyring.Set(Service, DefaultToken, token)
}

// GetToken retrieves the token from the system keychain.
func GetToken() (string, error) {
	// Check env first for CI/CD support
	if envToken := os.Getenv("HXTP_TOKEN"); envToken != "" {
		return envToken, nil
	}
	return keyring.Get(Service, DefaultToken)
}

// SaveConfig saves non-sensitive data to the config file.
func SaveConfig(cfg *Config) error {
	path := GetConfigPath()
	_ = os.MkdirAll(filepath.Dir(path), 0700)

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// LoadConfig reads metadata from the config file.
func LoadConfig() (*Config, error) {
	path := GetConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return &Config{ApiUrl: "https://api.hestialabs.in/api/v1"}, nil
	}
	var cfg Config
	err = json.Unmarshal(data, &cfg)
	return &cfg, err
}

// RequireAuth ensures a token exists in the keychain.
func RequireAuth() (string, *Config, error) {
	token, err := GetToken()
	if err != nil {
		return "", nil, fmt.Errorf("Not authenticated. Run 'hxtpctl login'")
	}
	cfg, _ := LoadConfig()
	return token, cfg, nil
}
