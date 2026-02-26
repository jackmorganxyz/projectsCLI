package config

import (
	"os"
	"path/filepath"
	"strings"

	toml "github.com/pelletier/go-toml/v2"
)

// Folder represents a named project folder associated with a GitHub account.
type Folder struct {
	Name          string `toml:"name"`
	GitHubAccount string `toml:"github_account"`
}

// Config holds all projectsCLI configuration fields.
type Config struct {
	ProjectsDir    string   `toml:"projects_dir,omitempty"`
	Editor         string   `toml:"editor,omitempty"`
	GitHubUsername string   `toml:"github_username,omitempty"`
	AutoGitInit    bool     `toml:"auto_git_init"`
	Folders        []Folder `toml:"folders,omitempty"`
}

// FolderByName returns the folder with the given name, or nil if not found.
func (c Config) FolderByName(name string) *Folder {
	for i := range c.Folders {
		if c.Folders[i].Name == name {
			return &c.Folders[i]
		}
	}
	return nil
}

// FolderNames returns a slice of all configured folder names.
func (c Config) FolderNames() []string {
	names := make([]string, len(c.Folders))
	for i, f := range c.Folders {
		names[i] = f.Name
	}
	return names
}

// Defaults returns a Config with sensible default values.
func Defaults() Config {
	projDir, _ := ProjectsDir()
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	return Config{
		ProjectsDir: projDir,
		Editor:      editor,
		AutoGitInit: true,
	}
}

// Load reads config from the default path, returning defaults for missing files.
func Load() (Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return Defaults(), err
	}
	return LoadFromPath(path)
}

// LoadFromPath reads config from a specific path.
func LoadFromPath(path string) (Config, error) {
	cfg := Defaults()

	expanded, err := expandPath(path)
	if err != nil {
		return cfg, err
	}

	data, err := os.ReadFile(expanded)
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return cfg, err
	}

	if err := toml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}

	// Expand ~ in projects_dir if present.
	if cfg.ProjectsDir != "" {
		cfg.ProjectsDir, _ = expandPath(cfg.ProjectsDir)
	}

	return cfg, nil
}

// Save writes config to the default path.
func Save(cfg Config) error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}
	return SaveToPath(cfg, path)
}

// SaveToPath writes config to a specific path.
func SaveToPath(cfg Config, path string) error {
	expanded, err := expandPath(path)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(expanded), 0700); err != nil {
		return err
	}

	data, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(expanded, data, 0600)
}

func expandPath(path string) (string, error) {
	if path == "~" || strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		if path == "~" {
			return home, nil
		}
		return filepath.Join(home, strings.TrimPrefix(path, "~/")), nil
	}
	return path, nil
}
