package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Init initializes a git repository in the given directory.
func Init(dir string) error {
	return run(dir, "git", "init")
}

// AddAll stages all changes.
func AddAll(dir string) error {
	return run(dir, "git", "add", "-A")
}

// Commit creates a commit with the given message.
func Commit(dir string, message string) error {
	return run(dir, "git", "commit", "-m", message)
}

// Push pushes to the remote.
func Push(dir string) error {
	return run(dir, "git", "push")
}

// PushSetUpstream pushes and sets the upstream branch.
func PushSetUpstream(dir string, remote, branch string) error {
	return run(dir, "git", "push", "-u", remote, branch)
}

// Status returns the git status output for a directory.
func Status(dir string) (string, error) {
	return output(dir, "git", "status", "--short")
}

// IsRepo checks if a directory is a git repository.
func IsRepo(dir string) bool {
	err := run(dir, "git", "rev-parse", "--is-inside-work-tree")
	return err == nil
}

// HasRemote checks if the repository has a remote configured.
func HasRemote(dir string) bool {
	out, err := output(dir, "git", "remote")
	return err == nil && strings.TrimSpace(out) != ""
}

// RemoteURL returns the URL of the first remote.
func RemoteURL(dir string) (string, error) {
	return output(dir, "git", "remote", "get-url", "origin")
}

// CurrentBranch returns the current branch name.
func CurrentBranch(dir string) (string, error) {
	return output(dir, "git", "branch", "--show-current")
}

// HasUncommitted checks if there are uncommitted changes.
func HasUncommitted(dir string) (bool, error) {
	out, err := Status(dir)
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(out) != "", nil
}

// run executes a command in the given directory.
func run(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg != "" {
			return fmt.Errorf("%s: %s", err, msg)
		}
		return err
	}
	return nil
}

// output executes a command and returns its stdout.
func output(dir string, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg != "" {
			return "", fmt.Errorf("%s: %s", err, msg)
		}
		return "", err
	}
	return strings.TrimSpace(stdout.String()), nil
}
