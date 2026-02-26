package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"unicode"

	"github.com/jackmorganxyz/projectsCLI/internal/config"
	"github.com/jackmorganxyz/projectsCLI/internal/project"
	"golang.org/x/text/unicode/norm"
)

// openFile opens a file or directory with the OS default application.
func openFile(path string) error {
	var cmd string
	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
	case "windows":
		cmd = "start"
	default: // linux, freebsd, etc.
		cmd = "xdg-open"
	}
	return exec.Command(cmd, path).Start()
}

var slugRegexp = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

// ValidateSlug checks that a project slug is valid (lowercase, hyphens, no spaces).
func ValidateSlug(slug string) error {
	if slug == "" {
		return fmt.Errorf("slug cannot be empty")
	}
	if len(slug) > 64 {
		return fmt.Errorf("slug too long (max 64 characters)")
	}
	if !slugRegexp.MatchString(slug) {
		return fmt.Errorf("invalid slug %q: must be lowercase alphanumeric with hyphens (e.g. my-project)", slug)
	}
	return nil
}

// Slugify converts a human-readable title into a valid slug.
// e.g. "My Cool Project!" → "my-cool-project"
func Slugify(title string) string {
	// Normalize unicode to decomposed form, then drop non-ASCII.
	t := norm.NFD.String(title)
	var sb strings.Builder
	for _, r := range t {
		if r <= unicode.MaxASCII && (unicode.IsLetter(r) || unicode.IsDigit(r)) {
			sb.WriteRune(unicode.ToLower(r))
		} else if r == ' ' || r == '-' || r == '_' {
			sb.WriteRune('-')
		}
	}
	// Collapse consecutive hyphens and trim leading/trailing hyphens.
	slug := regexp.MustCompile(`-{2,}`).ReplaceAllString(sb.String(), "-")
	slug = strings.Trim(slug, "-")
	if len(slug) > 64 {
		slug = slug[:64]
		slug = strings.TrimRight(slug, "-")
	}
	return slug
}

// writeJSON encodes v as indented JSON to w.
func writeJSON(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

// extractField extracts a field value from a struct or map based on dot-notation path.
// For example: "dir", "meta.title", "meta.slug"
func extractField(v any, path string) (any, error) {
	// Convert struct/slice to map via JSON
	var m map[string]any
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("marshal to JSON: %w", err)
	}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("unmarshal from JSON: %w", err)
	}

	parts := strings.Split(path, ".")
	current := any(m)

	for _, part := range parts {
		switch c := current.(type) {
		case map[string]any:
			val, ok := c[part]
			if !ok {
				return nil, fmt.Errorf("field %q not found", path)
			}
			current = val
		default:
			return nil, fmt.Errorf("cannot access field %q on type %T", part, current)
		}
	}
	return current, nil
}

// findProject locates a project by slug, searching the top-level projects directory
// and all configured folders. If folderHint is non-empty, only that folder is searched.
func findProject(cfg config.Config, slug string, folderHint string) (*project.Project, error) {
	if folderHint != "" {
		// Search only in the specified folder.
		folderDir := filepath.Join(cfg.ProjectsDir, folderHint)
		proj, err := project.FindProject(folderDir, slug)
		if err != nil {
			return nil, fmt.Errorf("project %q not found in folder %q", slug, folderHint)
		}
		proj.Folder = folderHint
		return proj, nil
	}

	// Search top-level first.
	proj, err := project.FindProject(cfg.ProjectsDir, slug)
	if err == nil {
		return proj, nil
	}

	// Search configured folders.
	for _, f := range cfg.Folders {
		folderDir := filepath.Join(cfg.ProjectsDir, f.Name)
		proj, err := project.FindProject(folderDir, slug)
		if err == nil {
			proj.Folder = f.Name
			return proj, nil
		}
	}

	return nil, fmt.Errorf("project %q not found", slug)
}

// listAllProjects lists projects from the top-level and all configured folders.
// If folderHint is non-empty, only that folder is listed.
func listAllProjects(cfg config.Config, folderHint string) ([]*project.Project, error) {
	if folderHint != "" {
		folderDir := filepath.Join(cfg.ProjectsDir, folderHint)
		projects, err := project.ListProjects(folderDir)
		if err != nil {
			return nil, err
		}
		for _, p := range projects {
			p.Folder = folderHint
		}
		return projects, nil
	}

	// Collect top-level projects.
	all, err := project.ListProjects(cfg.ProjectsDir)
	if err != nil {
		return nil, err
	}

	// Collect projects from each configured folder.
	for _, f := range cfg.Folders {
		folderDir := filepath.Join(cfg.ProjectsDir, f.Name)
		projects, err := project.ListProjects(folderDir)
		if err != nil {
			continue
		}
		for _, p := range projects {
			p.Folder = f.Name
			all = append(all, p)
		}
	}

	sort.Slice(all, func(i, j int) bool {
		return all[i].Meta.Slug < all[j].Meta.Slug
	})

	return all, nil
}

// folderForProject returns the folder config for a project, or nil if it's top-level.
func folderForProject(cfg config.Config, proj *project.Project) *config.Folder {
	if proj.Folder == "" {
		return nil
	}
	return cfg.FolderByName(proj.Folder)
}
