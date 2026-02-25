package cli

import "testing"

func TestSlugify(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"My Cool Project", "my-cool-project"},
		{"Hello World!", "hello-world"},
		{"  spaces  everywhere  ", "spaces-everywhere"},
		{"UPPER CASE", "upper-case"},
		{"already-a-slug", "already-a-slug"},
		{"café résumé", "cafe-resume"},
		{"lots---of---dashes", "lots-of-dashes"},
		{"under_score", "under-score"},
		{"123 Numbers First", "123-numbers-first"},
		{"", ""},
		{"!!!", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := Slugify(tt.input)
			if got != tt.want {
				t.Errorf("Slugify(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidateSlug(t *testing.T) {
	// Valid slugs.
	for _, s := range []string{"my-project", "a", "hello-world-123"} {
		if err := ValidateSlug(s); err != nil {
			t.Errorf("ValidateSlug(%q) unexpected error: %v", s, err)
		}
	}

	// Invalid slugs.
	for _, s := range []string{"", "Has Spaces", "UPPER", "trailing-"} {
		if err := ValidateSlug(s); err == nil {
			t.Errorf("ValidateSlug(%q) expected error, got nil", s)
		}
	}
}
