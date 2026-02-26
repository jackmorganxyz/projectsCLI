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

// ListAuthAccounts returns the GitHub usernames authenticated via gh auth.
// Returns nil if gh is not installed or no accounts are logged in.
func ListAuthAccounts() []string {
	out, err := output(".", "gh", "auth", "status", "--format", "{{range .}}{{.account}}\n{{end}}")
	if err != nil {
		return nil
	}
	var accounts []string
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			accounts = append(accounts, line)
		}
	}
	return accounts
}

// IsAuthAccount checks if the given account is authenticated via gh auth.
func IsAuthAccount(account string) bool {
	for _, a := range ListAuthAccounts() {
		if strings.EqualFold(a, account) {
			return true
		}
	}
	return false
}
