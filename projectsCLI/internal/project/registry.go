package project

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ListProjects scans the projects directory and returns all valid projects.
func ListProjects(projectsDir string) ([]*Project, error) {
	entries, err := os.ReadDir(projectsDir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read projects dir: %w", err)
	}

	var projects []*Project
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		projDir := filepath.Join(projectsDir, entry.Name())
		projFile := ProjectFilePath(projDir)

		if _, err := os.Stat(projFile); err != nil {
			continue // skip directories without PROJECT.md
		}

		proj, err := LoadProject(projDir)
		if err != nil {
			continue // skip unparseable projects
		}

		projects = append(projects, proj)
	}

	sort.Slice(projects, func(i, j int) bool {
		return projects[i].Meta.Slug < projects[j].Meta.Slug
	})

	return projects, nil
}

// FindProject finds a project by slug in the projects directory.
func FindProject(projectsDir, slug string) (*Project, error) {
	dir := filepath.Join(projectsDir, slug)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, fmt.Errorf("project %q not found", slug)
	}
	return LoadProject(dir)
}

// WriteRegistry regenerates PROJECTS.md from the filesystem.
func WriteRegistry(projectsDir string) error {
	projects, err := ListProjects(projectsDir)
	if err != nil {
		return err
	}

	var sb strings.Builder
	sb.WriteString("# Projects\n\n")
	sb.WriteString("Auto-generated registry of all projects.\n\n")

	if len(projects) == 0 {
		sb.WriteString("No projects yet. Run `projctl create <slug>` to get started.\n")
	} else {
		sb.WriteString("| Slug | Title | Status | Created |\n")
		sb.WriteString("|------|-------|--------|---------|\n")
		for _, p := range projects {
			created := p.Meta.CreatedAt
			if len(created) > 10 {
				created = created[:10]
			}
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
				p.Meta.Slug, p.Meta.Title, p.Meta.Status, created))
		}
	}

	sb.WriteString("\n")

	registryPath := filepath.Join(projectsDir, "PROJECTS.md")
	return os.WriteFile(registryPath, []byte(sb.String()), 0644)
}
