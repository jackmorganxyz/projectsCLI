package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
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

// writeJSON encodes v as indented JSON to w.
func writeJSON(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
