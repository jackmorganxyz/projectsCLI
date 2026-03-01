package cli

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jackmorganxyz/projectsCLI/internal/tui"
	"github.com/spf13/cobra"
)

// githubRelease represents the subset of the GitHub Releases API response we need.
type githubRelease struct {
	TagName string `json:"tag_name"`
}

// upgradeResult holds the JSON output for the upgrade command.
type upgradeResult struct {
	Status         string `json:"status"`
	CurrentVersion string `json:"current_version"`
	LatestVersion  string `json:"latest_version"`
	Method         string `json:"method,omitempty"`
	Message        string `json:"message,omitempty"`
}

// NewUpgradeCmd upgrades the CLI binary to the latest release.
func NewUpgradeCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade to the latest version",
		Long:  "Check for a newer release on GitHub and upgrade the projects binary in place.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runUpgrade(cmd, version)
		},
	}
	return cmd
}

func runUpgrade(cmd *cobra.Command, currentVersion string) error {
	w := cmd.OutOrStdout()

	// Fetch latest release from GitHub.
	if !tui.IsJSON() {
		fmt.Fprintln(w, tui.InfoMessage("Checking for updates... 🔍"))
	}

	latest, err := fetchLatestRelease()
	if err != nil {
		return fmt.Errorf("check for updates: %w", err)
	}
	latestVersion := strings.TrimPrefix(latest.TagName, "v")

	// Dev builds can't be upgraded.
	if currentVersion == "dev" {
		if tui.IsJSON() {
			return writeJSON(w, upgradeResult{
				Status:         "dev_build",
				CurrentVersion: currentVersion,
				LatestVersion:  latestVersion,
				Message:        "Running a dev build, skipping upgrade",
			})
		}
		fmt.Fprintln(w, tui.WarningMessage("Running a dev build — upgrade skipped."))
		fmt.Fprintln(w, tui.Muted(fmt.Sprintf("  Latest release: v%s", latestVersion)))
		return nil
	}

	// Already up to date.
	if currentVersion == latestVersion {
		if tui.IsJSON() {
			return writeJSON(w, upgradeResult{
				Status:         "up_to_date",
				CurrentVersion: currentVersion,
				LatestVersion:  latestVersion,
			})
		}
		fmt.Fprintln(w, tui.SuccessMessage(fmt.Sprintf("v%s is the latest — %s", currentVersion, tui.RandomAlreadyLatest())))
		return nil
	}

	// Resolve current binary path.
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("find executable: %w", err)
	}
	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return fmt.Errorf("resolve executable path: %w", err)
	}

	if isHomebrew(execPath) {
		return upgradeViaHomebrew(cmd, currentVersion, latestVersion)
	}
	return upgradeViaBinary(cmd, execPath, currentVersion, latestVersion, latest.TagName)
}

// fetchLatestRelease calls the GitHub Releases API.
func fetchLatestRelease() (*githubRelease, error) {
	resp, err := http.Get("https://api.github.com/repos/jackmorganxyz/projectsCLI/releases/latest")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("parse release: %w", err)
	}
	return &release, nil
}

// isHomebrew returns true if the binary path looks like a Homebrew installation.
func isHomebrew(execPath string) bool {
	lower := strings.ToLower(execPath)
	return strings.Contains(lower, "homebrew") || strings.Contains(lower, "cellar") || strings.Contains(lower, "linuxbrew")
}

// upgradeViaHomebrew delegates to brew upgrade.
func upgradeViaHomebrew(cmd *cobra.Command, currentVersion, latestVersion string) error {
	w := cmd.OutOrStdout()

	if tui.IsJSON() {
		return writeJSON(w, upgradeResult{
			Status:         "upgraded",
			CurrentVersion: currentVersion,
			LatestVersion:  latestVersion,
			Method:         "homebrew",
			Message:        "Upgraded via Homebrew",
		})
	}

	fmt.Fprintln(w, tui.InfoMessage(fmt.Sprintf("Homebrew detected — upgrading v%s → v%s 🍺", currentVersion, latestVersion)))
	fmt.Fprintln(w)

	brewCmd := exec.Command("brew", "upgrade", "jackmorganxyz/tap/projects")
	brewCmd.Stdout = cmd.OutOrStdout()
	brewCmd.Stderr = cmd.ErrOrStderr()
	if err := brewCmd.Run(); err != nil {
		return fmt.Errorf("brew upgrade failed: %w", err)
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, tui.SuccessMessage(tui.RandomUpgradeCheer()))
	return nil
}

// upgradeViaBinary downloads and replaces the binary directly.
func upgradeViaBinary(cmd *cobra.Command, execPath, currentVersion, latestVersion, tag string) error {
	w := cmd.OutOrStdout()

	// Check if the binary's directory is writable by creating a temp file.
	dir := filepath.Dir(execPath)
	tmpFile, err := os.CreateTemp(dir, ".projects-upgrade-*")
	if err != nil {
		return fmt.Errorf("cannot write to %s — try running with sudo: %w", dir, err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	if !tui.IsJSON() {
		fmt.Fprintln(w, tui.InfoMessage(fmt.Sprintf("Downloading v%s... 📦", latestVersion)))
	}

	// Build the archive URL.
	archiveName := fmt.Sprintf("projects_%s_%s_%s.tar.gz", latestVersion, runtime.GOOS, runtime.GOARCH)
	archiveURL := fmt.Sprintf("https://github.com/jackmorganxyz/projectsCLI/releases/download/%s/%s", tag, archiveName)

	// Download the archive.
	resp, err := http.Get(archiveURL)
	if err != nil {
		return fmt.Errorf("download release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: HTTP %d (is %s available for %s/%s?)", resp.StatusCode, tag, runtime.GOOS, runtime.GOARCH)
	}

	// Extract the "projects" binary from the tar.gz.
	binaryData, err := extractBinaryFromTarGz(resp.Body, "projects")
	if err != nil {
		return err
	}

	// Write to temp file, then atomic rename.
	if err := os.WriteFile(tmpPath, binaryData, 0755); err != nil {
		return fmt.Errorf("write binary: %w", err)
	}
	if err := os.Rename(tmpPath, execPath); err != nil {
		return fmt.Errorf("replace binary: %w", err)
	}

	if tui.IsJSON() {
		return writeJSON(w, upgradeResult{
			Status:         "upgraded",
			CurrentVersion: currentVersion,
			LatestVersion:  latestVersion,
			Method:         "binary",
			Message:        fmt.Sprintf("Upgraded from v%s to v%s", currentVersion, latestVersion),
		})
	}

	fmt.Fprintln(w, tui.SuccessMessage(fmt.Sprintf("Upgraded to v%s — %s", latestVersion, tui.RandomUpgradeCheer())))
	if tip := tui.MaybeTip(); tip != "" {
		fmt.Fprintln(w, tip)
	}
	return nil
}

// extractBinaryFromTarGz reads a tar.gz stream and returns the contents of the named file.
func extractBinaryFromTarGz(r io.Reader, name string) ([]byte, error) {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("decompress: %w", err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil, fmt.Errorf("binary %q not found in archive", name)
		}
		if err != nil {
			return nil, fmt.Errorf("read archive: %w", err)
		}
		if filepath.Base(hdr.Name) == name && hdr.Typeflag == tar.TypeReg {
			data, err := io.ReadAll(tr)
			if err != nil {
				return nil, fmt.Errorf("read binary from archive: %w", err)
			}
			return data, nil
		}
	}
}
