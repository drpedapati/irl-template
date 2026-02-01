package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Profile contains academic/personal info for template injection
type Profile struct {
	Name         string `json:"name"`
	Title        string `json:"title"`        // e.g., "PhD Candidate", "Professor"
	Institution  string `json:"institution"`  // e.g., "Stanford University"
	Department   string `json:"department"`   // e.g., "Department of Psychology"
	Email        string `json:"email"`
	Instructions string `json:"instructions"` // Common instructions for AI
}

type Config struct {
	DefaultDirectory string  `json:"default_directory"`
	Profile          Profile `json:"profile"`
	FavoriteEditors  []string `json:"favorite_editors,omitempty"` // Editor cmd names (e.g., "cursor", "code")
	PlanEditor       string   `json:"plan_editor,omitempty"`      // Plan editor: "nano", "vim", "code", "cursor", "auto"
	PlanEditorType   string   `json:"plan_editor_type,omitempty"` // "terminal" or "gui"
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

// ClearProfile removes all saved profile data
func ClearProfile() error {
	return SetProfile(Profile{})
}

// ClearDefaultDirectory removes the saved default directory
func ClearDefaultDirectory() error {
	return SetDefaultDirectory("")
}

// HasProfile returns true if any profile fields are set
func HasProfile() bool {
	p := GetProfile()
	return p.Name != "" || p.Institution != "" || p.Title != ""
}

// GetFavoriteEditors returns the list of favorite editor command names
func GetFavoriteEditors() []string {
	cfg, err := Load()
	if err != nil {
		return nil
	}
	return cfg.FavoriteEditors
}

// SetFavoriteEditors saves the list of favorite editor command names
func SetFavoriteEditors(editors []string) error {
	cfg, err := Load()
	if err != nil {
		cfg = &Config{}
	}
	cfg.FavoriteEditors = editors
	return cfg.Save()
}

// ToggleFavoriteEditor adds or removes an editor from favorites
func ToggleFavoriteEditor(cmd string) error {
	favorites := GetFavoriteEditors()

	// Check if already a favorite
	for i, f := range favorites {
		if f == cmd {
			// Remove it
			favorites = append(favorites[:i], favorites[i+1:]...)
			return SetFavoriteEditors(favorites)
		}
	}

	// Add it
	favorites = append(favorites, cmd)
	return SetFavoriteEditors(favorites)
}

// IsFavoriteEditor returns true if the editor is a favorite
func IsFavoriteEditor(cmd string) bool {
	favorites := GetFavoriteEditors()
	for _, f := range favorites {
		if f == cmd {
			return true
		}
	}
	return false
}

// ClearFavoriteEditors removes all favorite editors
func ClearFavoriteEditors() error {
	return SetFavoriteEditors(nil)
}

// GetPlanEditor returns the configured plan editor command
func GetPlanEditor() string {
	cfg, err := Load()
	if err != nil {
		return ""
	}
	return cfg.PlanEditor
}

// GetPlanEditorType returns the configured plan editor type ("terminal" or "gui")
func GetPlanEditorType() string {
	cfg, err := Load()
	if err != nil {
		return ""
	}
	return cfg.PlanEditorType
}

// SetPlanEditor saves the plan editor preference
func SetPlanEditor(editor, editorType string) error {
	cfg, err := Load()
	if err != nil {
		cfg = &Config{}
	}
	cfg.PlanEditor = editor
	cfg.PlanEditorType = editorType
	return cfg.Save()
}

// ClearPlanEditor removes the plan editor preference
func ClearPlanEditor() error {
	return SetPlanEditor("", "")
}
