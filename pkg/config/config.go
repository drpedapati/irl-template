package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DefaultDirectory string `json:"default_directory"`
}

var configPath string

func init() {
	home, _ := os.UserHomeDir()
	configPath = filepath.Join(home, ".irl", "config.json")
}

func Load() (*Config, error) {
	cfg := &Config{}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil // Return empty config if file doesn't exist
		}
		return nil, err
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Save() error {
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

func GetDefaultDirectory() string {
	cfg, err := Load()
	if err != nil {
		return ""
	}
	return cfg.DefaultDirectory
}

func SetDefaultDirectory(dir string) error {
	cfg, err := Load()
	if err != nil {
		cfg = &Config{}
	}
	cfg.DefaultDirectory = dir
	return cfg.Save()
}
