package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Profile contains academic/personal info for template injection
type Profile struct {
	Name         string `json:"name"`
	Title        string `json:"title"`       // e.g., "PhD Candidate", "Professor"
	Institution  string `json:"institution"` // e.g., "Stanford University"
	Department   string `json:"department"`  // e.g., "Department of Psychology"
	Email        string `json:"email"`
	Instructions string `json:"instructions"` // Common instructions for AI
}

type Config struct {
	DefaultDirectory string  `json:"default_directory"`
	Profile          Profile `json:"profile"`
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

func GetProfile() Profile {
	cfg, err := Load()
	if err != nil {
		return Profile{}
	}
	return cfg.Profile
}

func SetProfile(profile Profile) error {
	cfg, err := Load()
	if err != nil {
		cfg = &Config{}
	}
	cfg.Profile = profile
	return cfg.Save()
}

func ClearProfile() error {
	return SetProfile(Profile{})
}

func ClearDefaultDirectory() error {
	return SetDefaultDirectory("")
}

// HasProfile returns true if any profile fields are set
func HasProfile() bool {
	p := GetProfile()
	return p.Name != "" || p.Institution != "" || p.Title != ""
}
