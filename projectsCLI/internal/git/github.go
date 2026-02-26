package git

import (
	"fmt"
	"strings"
)

// CreateRepo creates a GitHub repository using the gh CLI.
func CreateRepo(dir, name, org string, private bool) (string, error) {
	args := []string{"repo", "create"}

	repoName := name
	if org != "" {
		repoName = org + "/" + name
	}
	args = append(args, repoName)

	if private {
		args = append(args, "--private")
	} else {
		args = append(args, "--public")
	}

	args = append(args, "--source", dir, "--push")

	out, err := output(dir, "gh", args...)
	if err != nil {
		return "", fmt.Errorf("gh repo create: %w", err)
	}
	return strings.TrimSpace(out), nil
}

// HasGHCLI checks if the gh CLI is available.
func HasGHCLI() bool {
	_, err := output(".", "gh", "version")
	return err == nil
}

// SwitchAuth switches the active GitHub CLI account to the given user.
func SwitchAuth(account string) error {
	return run(".", "gh", "auth", "switch", "--user", account)
}
