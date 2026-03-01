package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/jackmorganxyz/projectsCLI/internal/config"
	"github.com/jackmorganxyz/projectsCLI/internal/git"
	"github.com/jackmorganxyz/projectsCLI/internal/tui"
	"github.com/spf13/cobra"
)

// NewFolderCmd creates the folder command group for managing account-based folders.
func NewFolderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "folder",
		Short: "Manage project folders for different GitHub accounts",
		Long: `Manage named folders that group projects by GitHub account.

Each folder maps to a subdirectory under your projects directory and is
associated with a specific GitHub account. When you push a project that
lives in a folder, the CLI automatically switches gh auth to the right account.`,
	}

	cmd.AddCommand(
		newFolderAddCmd(),
		newFolderListCmd(),
		newFolderRemoveCmd(),
	)

	return cmd
}

func newFolderAddCmd() *cobra.Command {
	var account string

	cmd := &cobra.Command{
		Use:   "add <name>",
		Short: "Add a new folder with a GitHub account",
		Long: `Add a named folder associated with a GitHub account.

The folder name is used as a subdirectory under your projects directory.
Projects created with --folder <name> will live in this directory and
push using the associated GitHub account.

If --account is omitted and gh is authenticated, you'll be prompted to
pick from your logged-in accounts.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			runtime, ok := RuntimeFromContext(cmd.Context())
			if !ok {
				return fmt.Errorf("missing runtime context")
			}

			name := args[0]

			// Validate the folder name using slug rules.
			if err := ValidateSlug(name); err != nil {
				return fmt.Errorf("invalid folder name: %w", err)
			}

			// Check for duplicate folder name.
			if runtime.Config.FolderByName(name) != nil {
				return fmt.Errorf("folder %q already exists", name)
			}

			// Resolve the account: flag > interactive picker > error.
			if account == "" {
				picked, err := pickGHAccount(cmd)
				if err != nil {
					return err
				}
				account = picked
			}

			// Warn (don't block) if gh isn't set up — the folder is still useful
			// as config, and auth can be sorted out before the first push.
			if !git.HasGHCLI() {
				fmt.Fprintln(cmd.ErrOrStderr(), tui.WarningMessage("gh CLI not found — install it and run 'gh auth login' before pushing"))
			} else if accounts := git.ListAuthAccounts(); len(accounts) > 0 && !git.IsAuthAccount(account) {
				fmt.Fprintln(cmd.ErrOrStderr(), tui.WarningMessage(
					fmt.Sprintf("account %q not found in gh auth (have: %s) — run 'gh auth login' to add it",
						account, strings.Join(accounts, ", "))))
			}

			// Create the folder directory.
			folderDir := filepath.Join(runtime.Config.ProjectsDir, name)
			if err := os.MkdirAll(folderDir, 0755); err != nil {
				return fmt.Errorf("create folder directory: %w", err)
			}

			// Add to config and save.
			runtime.Config.Folders = append(runtime.Config.Folders, config.Folder{
				Name:          name,
				GitHubAccount: account,
			})

			if err := config.SaveToPath(runtime.Config, runtime.ConfigPath); err != nil {
				return fmt.Errorf("save config: %w", err)
			}

			if tui.IsJSON() {
				return writeJSON(cmd.OutOrStdout(), map[string]string{
					"status":         "created",
					"folder":         name,
					"github_account": account,
					"path":           folderDir,
				})
			}

			w := cmd.OutOrStdout()
			fmt.Fprintln(w, tui.SuccessMessage(fmt.Sprintf("Folder %s created — %s", tui.Slug(name), tui.RandomFolderCheer())))
			fmt.Fprintln(w, tui.FormatField("GitHub account", tui.Slug(account)))
			fmt.Fprintln(w, tui.FormatField("Path", tui.Path(folderDir)))
			fmt.Fprintln(w)
			fmt.Fprintln(w, tui.InfoMessage(fmt.Sprintf("Create projects here with: projects create <slug> --folder %s", name)))
			return nil
		},
	}

	cmd.Flags().StringVar(&account, "account", "", "GitHub username/account for this folder")

	return cmd
}

// pickGHAccount tries to interactively pick a gh account. Falls back to
// a clear error message if non-interactive or gh isn't available.
func pickGHAccount(cmd *cobra.Command) (string, error) {
	if !git.HasGHCLI() {
		return "", fmt.Errorf("--account is required (gh CLI not available for interactive selection)")
	}

	accounts := git.ListAuthAccounts()
	if len(accounts) == 0 {
		return "", fmt.Errorf("--account is required (no accounts found in gh auth)\n\nRun 'gh auth login' first, or pass --account <username>")
	}

	if !tui.IsInteractive() {
		return "", fmt.Errorf("--account is required in non-interactive mode (available: %s)", strings.Join(accounts, ", "))
	}

	// Single account — just use it.
	if len(accounts) == 1 {
		fmt.Fprintln(cmd.ErrOrStderr(), tui.Muted(fmt.Sprintf("Using GitHub account %q", accounts[0])))
		return accounts[0], nil
	}

	// Multiple accounts — let the user pick.
	var selected string
	options := make([]huh.Option[string], len(accounts))
	for i, a := range accounts {
		options[i] = huh.NewOption(a, a)
	}

	err := huh.NewSelect[string]().
		Title("Which GitHub account for this folder?").
		Options(options...).
		Value(&selected).
		Run()

	if err != nil {
		return "", fmt.Errorf("selection cancelled")
	}
	if selected == "" {
		return "", fmt.Errorf("no account selected")
	}

	return selected, nil
}

func newFolderListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all configured folders",
		RunE: func(cmd *cobra.Command, _ []string) error {
			runtime, ok := RuntimeFromContext(cmd.Context())
			if !ok {
				return fmt.Errorf("missing runtime context")
			}

			folders := runtime.Config.Folders

			if tui.IsJSON() {
				return writeJSON(cmd.OutOrStdout(), folders)
			}

			if len(folders) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), tui.Muted("No folders configured. Use 'projects folder add <name> --account <gh-user>' to get started."))
				return nil
			}

			w := cmd.OutOrStdout()
			fmt.Fprintln(w, tui.Header("📂 Your Folders"))
			fmt.Fprintln(w)

			headers := []string{"Name", "GitHub Account", "Path"}
			var rows [][]string
			for _, f := range folders {
				path := filepath.Join(runtime.Config.ProjectsDir, f.Name)
				rows = append(rows, []string{f.Name, f.GitHubAccount, path})
			}

			fmt.Fprintln(w, tui.Table(headers, rows))
			return nil
		},
	}

	return cmd
}

func newFolderRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove <name>",
		Aliases: []string{"rm"},
		Short:   "Remove a folder configuration",
		Long:    "Remove a folder from the config. Does not delete the directory or its projects.",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			runtime, ok := RuntimeFromContext(cmd.Context())
			if !ok {
				return fmt.Errorf("missing runtime context")
			}

			name := args[0]

			// Find and remove the folder from config.
			found := false
			var remaining []config.Folder
			for _, f := range runtime.Config.Folders {
				if f.Name == name {
					found = true
					continue
				}
				remaining = append(remaining, f)
			}

			if !found {
				return fmt.Errorf("folder %q not found", name)
			}

			runtime.Config.Folders = remaining
			if err := config.SaveToPath(runtime.Config, runtime.ConfigPath); err != nil {
				return fmt.Errorf("save config: %w", err)
			}

			if tui.IsJSON() {
				return writeJSON(cmd.OutOrStdout(), map[string]string{
					"status": "removed",
					"folder": name,
				})
			}

			fmt.Fprintln(cmd.OutOrStdout(), tui.SuccessMessage(fmt.Sprintf("Folder %s removed from config", tui.Slug(name))))
			fmt.Fprintln(cmd.OutOrStdout(), tui.Muted("  "+tui.RandomFolderRemoveQuip()))
			fmt.Fprintln(cmd.OutOrStdout(), tui.InfoMessage("Directory and projects were not deleted."))
			return nil
		},
	}

	return cmd
}
