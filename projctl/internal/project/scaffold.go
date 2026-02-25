package project

import (
	"fmt"
	"os"
	"path/filepath"
)

// Scaffold creates the full directory tree and template files for a new project.
func Scaffold(projectsDir string, meta ProjectMeta) (string, error) {
	dir := filepath.Join(projectsDir, meta.Slug)

	if _, err := os.Stat(dir); err == nil {
		return "", fmt.Errorf("project directory already exists: %s", dir)
	}

	// Create directory tree.
	dirs := []string{
		dir,
		filepath.Join(dir, "docs"),
		filepath.Join(dir, "memory"),
		filepath.Join(dir, "context"),
		filepath.Join(dir, "tasks"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return "", fmt.Errorf("create directory %s: %w", d, err)
		}
	}

	// Write PROJECT.md with frontmatter.
	body := fmt.Sprintf("# %s\n\n%s\n", meta.Title, meta.Description)
	if err := WriteProjectFile(dir, meta, body); err != nil {
		return "", fmt.Errorf("write PROJECT.md: %w", err)
	}

	// Write template files.
	templates := map[string]string{
		filepath.Join(dir, "memory", "MEMORY.md"):   fmt.Sprintf("# %s - Memory\n\nPersistent notes and context for this project.\n", meta.Title),
		filepath.Join(dir, "context", "CONTEXT.md"): fmt.Sprintf("# %s - Context\n\nProject context, decisions, and architecture notes.\n", meta.Title),
		filepath.Join(dir, "tasks", "TODO.md"):      fmt.Sprintf("# %s - Tasks\n\n- [ ] Initial setup\n", meta.Title),
		filepath.Join(dir, "docs", "README.md"):     fmt.Sprintf("# %s\n\n%s\n", meta.Title, meta.Description),
	}

	for path, content := range templates {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return "", fmt.Errorf("write %s: %w", filepath.Base(path), err)
		}
	}

	return dir, nil
}
