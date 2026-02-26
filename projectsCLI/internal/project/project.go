package project

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// ProjectMeta is the YAML frontmatter stored in PROJECT.md.
type ProjectMeta struct {
	Title       string   `yaml:"title" json:"title"`
	Slug        string   `yaml:"slug" json:"slug"`
	Status      string   `yaml:"status" json:"status"`
	Tags        []string `yaml:"tags,omitempty" json:"tags,omitempty"`
	Description string   `yaml:"description,omitempty" json:"description,omitempty"`
	CreatedAt   string   `yaml:"created_at" json:"created_at"`
	UpdatedAt   string   `yaml:"updated_at" json:"updated_at"`
	GitRemote   string   `yaml:"git_remote,omitempty" json:"git_remote,omitempty"`
}

// Project is a fully loaded project with its metadata, body, and filesystem path.
type Project struct {
	Meta   ProjectMeta `json:"meta"`
	Body   string      `json:"body,omitempty"`
	Dir    string      `json:"dir"`
	Folder string      `json:"folder,omitempty"`
}

// ProjectFilePath returns the PROJECT.md path for a project directory.
func ProjectFilePath(dir string) string {
	return filepath.Join(dir, "PROJECT.md")
}

// ParseProjectFile reads and parses a PROJECT.md file.
func ParseProjectFile(path string) (*Project, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	meta, body, err := parseFrontmatter(string(data))
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}

	return &Project{
		Meta: *meta,
		Body: body,
		Dir:  filepath.Dir(path),
	}, nil
}

// LoadProject loads a project from its directory.
func LoadProject(dir string) (*Project, error) {
	return ParseProjectFile(ProjectFilePath(dir))
}

// WriteProjectFile writes a PROJECT.md with YAML frontmatter and markdown body.
func WriteProjectFile(dir string, meta ProjectMeta, body string) error {
	yamlBytes, err := yaml.Marshal(meta)
	if err != nil {
		return fmt.Errorf("marshal metadata: %w", err)
	}

	var sb strings.Builder
	sb.WriteString("---\n")
	sb.Write(yamlBytes)
	sb.WriteString("---\n")
	if body != "" {
		sb.WriteString("\n")
		sb.WriteString(body)
		if !strings.HasSuffix(body, "\n") {
			sb.WriteString("\n")
		}
	}

	return os.WriteFile(ProjectFilePath(dir), []byte(sb.String()), 0644)
}

// NewMeta creates a ProjectMeta with defaults filled in.
func NewMeta(slug, title string) ProjectMeta {
	now := time.Now().UTC().Format(time.RFC3339)
	if title == "" {
		title = slug
	}
	return ProjectMeta{
		Title:     title,
		Slug:      slug,
		Status:    "active",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// parseFrontmatter splits a document into YAML frontmatter and markdown body.
func parseFrontmatter(content string) (*ProjectMeta, string, error) {
	scanner := bufio.NewScanner(strings.NewReader(content))

	// First line must be "---"
	if !scanner.Scan() || strings.TrimSpace(scanner.Text()) != "---" {
		return nil, "", fmt.Errorf("missing opening frontmatter delimiter")
	}

	var yamlLines []string
	foundEnd := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "---" {
			foundEnd = true
			break
		}
		yamlLines = append(yamlLines, line)
	}

	if !foundEnd {
		return nil, "", fmt.Errorf("missing closing frontmatter delimiter")
	}

	var meta ProjectMeta
	if err := yaml.Unmarshal([]byte(strings.Join(yamlLines, "\n")), &meta); err != nil {
		return nil, "", fmt.Errorf("unmarshal frontmatter: %w", err)
	}

	// Remaining content is the body.
	var bodyLines []string
	for scanner.Scan() {
		bodyLines = append(bodyLines, scanner.Text())
	}
	body := strings.Join(bodyLines, "\n")
	body = strings.TrimPrefix(body, "\n")

	return &meta, body, nil
}
