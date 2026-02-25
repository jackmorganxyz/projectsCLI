package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

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
