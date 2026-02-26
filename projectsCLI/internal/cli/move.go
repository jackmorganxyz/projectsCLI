package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jackmorganxyz/projectsCLI/internal/project"
	"github.com/jackmorganxyz/projectsCLI/internal/tui"
	"github.com/spf13/cobra"
)

// NewMoveCmd moves a project into or out of a folder.
func NewMoveCmd() *cobra.Command {
	var folder string

	cmd := &cobra.Command{
		Use:   "move <slug>",
		Short: "Move a project to a different folder",
		Long: `Move an existing project into a folder, out of a folder, or between folders.

Use --folder <name> to move into a folder. Use --folder "" to move to the top level.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			runtime, ok := RuntimeFromContext(cmd.Context())
			if !ok {
				return fmt.Errorf("missing runtime context")
			}

			if !cmd.Flags().Changed("folder") {
				return fmt.Errorf("--folder is required: specify the target folder (or --folder \"\" for top level)")
			}

			slug := args[0]

			// Find the project wherever it currently lives.
			proj, err := findProject(runtime.Config, slug, "")
			if err != nil {
				return err
			}

			// Determine the destination.
			var destDir string
			if folder == "" {
				// Moving to top level.
				destDir = filepath.Join(runtime.Config.ProjectsDir, slug)
			} else {
				if runtime.Config.FolderByName(folder) == nil {
					return fmt.Errorf("folder %q not configured; run 'projects folder add %s --account <gh-user>' first", folder, folder)
				}
				destDir = filepath.Join(runtime.Config.ProjectsDir, folder, slug)
			}

			// Check we're not moving to the same place.
			if proj.Dir == destDir {
				if folder == "" {
					return fmt.Errorf("project %q is already at the top level", slug)
				}
				return fmt.Errorf("project %q is already in folder %q", slug, folder)
			}

			// Check destination doesn't already exist.
			if _, err := os.Stat(destDir); err == nil {
				return fmt.Errorf("destination already exists: %s", destDir)
			}

			// Ensure parent directory exists.
			if err := os.MkdirAll(filepath.Dir(destDir), 0755); err != nil {
				return fmt.Errorf("create destination directory: %w", err)
			}

			// Move the project directory.
			if err := os.Rename(proj.Dir, destDir); err != nil {
				return fmt.Errorf("move project: %w", err)
			}

			// Regenerate registry.
			_ = project.WriteRegistry(runtime.Config.ProjectsDir)

			if tui.IsJSON() {
				result := map[string]any{
					"status": "moved",
					"slug":   slug,
					"from":   proj.Dir,
					"to":     destDir,
				}
				if proj.Folder != "" {
					result["from_folder"] = proj.Folder
				}
				if folder != "" {
					result["to_folder"] = folder
				}
				return writeJSON(cmd.OutOrStdout(), result)
			}

			w := cmd.OutOrStdout()
			fromLabel := "top level"
			if proj.Folder != "" {
				fromLabel = fmt.Sprintf("folder %q", proj.Folder)
			}
			toLabel := "top level"
			if folder != "" {
				toLabel = fmt.Sprintf("folder %q", folder)
			}
			fmt.Fprintln(w, tui.SuccessMessage(fmt.Sprintf("Moved %q from %s to %s", slug, fromLabel, toLabel)))
			fmt.Fprintln(w, tui.FormatField("New path", destDir))
			return nil
		},
	}

	cmd.Flags().StringVar(&folder, "folder", "", "target folder (empty string for top level)")

	return cmd
}
