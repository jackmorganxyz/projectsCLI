package config

import (
	"os"
	"path/filepath"
)

// AppDir returns the root config directory (~/.projects).
func AppDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".projects"), nil
}

// ProjectsDir returns the default projects directory (~/.projects/projects).
func ProjectsDir() (string, error) {
	root, err := AppDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "projects"), nil
}

// OpenClawDir returns the openclaw directory (~/.openclaw) if it exists.
// Returns the path and true if found, empty string and false otherwise.
func OpenClawDir() (string, bool) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", false
	}
	dir := filepath.Join(home, ".openclaw")
	if info, err := os.Stat(dir); err == nil && info.IsDir() {
		return dir, true
	}
	return "", false
}

// OpenClawProjectsDir returns the openclaw projects directory path.
func OpenClawProjectsDir() string {
	if dir, ok := OpenClawDir(); ok {
		return filepath.Join(dir, "projects")
	}
	return ""
}

// ConfigPath returns the path to config.toml (~/.projects/config.toml).
func ConfigPath() (string, error) {
	root, err := AppDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "config.toml"), nil
}

// EnsureDirs creates all required directories if they don't exist.
func EnsureDirs() error {
	for _, fn := range []func() (string, error){AppDir, ProjectsDir} {
		dir, err := fn()
		if err != nil {
			return err
		}
		if err := os.MkdirAll(dir, 0700); err != nil {
			return err
		}
	}
	return nil
}
