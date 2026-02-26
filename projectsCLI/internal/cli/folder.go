package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jackmorganxyz/projectsCLI/internal/config"
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
push using the associated GitHub account.`,
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

			if account == "" {
				return fmt.Errorf("--account is required: specify the GitHub username for this folder")
			}

			// Check for duplicate folder name.
			if runtime.Config.FolderByName(name) != nil {
				return fmt.Errorf("folder %q already exists", name)
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
			fmt.Fprintln(w, tui.SuccessMessage(fmt.Sprintf("Folder %q created", name)))
			fmt.Fprintln(w, tui.FormatField("GitHub account", account))
			fmt.Fprintln(w, tui.FormatField("Path", folderDir))
			fmt.Fprintln(w)
			fmt.Fprintln(w, tui.Muted(fmt.Sprintf("Create projects here with: projects create <slug> --folder %s", name)))
			return nil
		},
	}

	cmd.Flags().StringVar(&account, "account", "", "GitHub username/account for this folder (required)")
	_ = cmd.MarkFlagRequired("account")

	return cmd
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

			headers := []string{"Name", "GitHub Account", "Path"}
			var rows [][]string
			for _, f := range folders {
				path := filepath.Join(runtime.Config.ProjectsDir, f.Name)
				rows = append(rows, []string{f.Name, f.GitHubAccount, path})
			}

			fmt.Fprintln(cmd.OutOrStdout(), tui.Table(headers, rows))
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

			fmt.Fprintln(cmd.OutOrStdout(), tui.SuccessMessage(fmt.Sprintf("Folder %q removed from config", name)))
			fmt.Fprintln(cmd.OutOrStdout(), tui.Muted("Directory and projects were not deleted."))
			return nil
		},
	}

	return cmd
}
