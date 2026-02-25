package config

import (
	"os"
	"path/filepath"
)

// OpenClawDir returns the root directory (~/.openclaw).
func OpenClawDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".openclaw"), nil
}

// ProjectsDir returns the projects directory (~/.openclaw/projects).
func ProjectsDir() (string, error) {
	root, err := OpenClawDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "projects"), nil
}

// ConfigPath returns the path to config.toml (~/.openclaw/config.toml).
func ConfigPath() (string, error) {
	root, err := OpenClawDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "config.toml"), nil
}

// EnsureDirs creates all required directories if they don't exist.
func EnsureDirs() error {
	for _, fn := range []func() (string, error){OpenClawDir, ProjectsDir} {
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
