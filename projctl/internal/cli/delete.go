package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jackpmorgan/projctl/internal/project"
	"github.com/jackpmorgan/projctl/internal/tui"
	"github.com/spf13/cobra"
)

// NewDeleteCmd deletes a project.
func NewDeleteCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:     "delete <slug>",
		Aliases: []string{"rm"},
		Short:   "Delete a project",
		Long:    "Delete a project and its directory. Use --force to skip confirmation.",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			runtime, ok := RuntimeFromContext(cmd.Context())
			if !ok {
				return fmt.Errorf("missing runtime context")
			}

			slug := args[0]
			proj, err := project.FindProject(runtime.Config.ProjectsDir, slug)
			if err != nil {
				return err
			}

			if !force && tui.IsInteractive() {
				confirmed, err := tui.RunConfirm(fmt.Sprintf("Delete project %q? This cannot be undone.", slug))
				if err != nil {
					return err
				}
				if !confirmed {
					fmt.Fprintln(cmd.OutOrStdout(), tui.Muted("Cancelled."))
					return nil
				}
			} else if !force {
				return fmt.Errorf("use --force to delete without confirmation in non-interactive mode")
			}

			dir := filepath.Join(runtime.Config.ProjectsDir, proj.Meta.Slug)
			if err := os.RemoveAll(dir); err != nil {
				return fmt.Errorf("remove project directory: %w", err)
			}

			// Regenerate registry.
			_ = project.WriteRegistry(runtime.Config.ProjectsDir)

			if tui.IsJSON() {
				return writeJSON(cmd.OutOrStdout(), map[string]string{
					"status": "deleted",
					"slug":   slug,
				})
			}

			fmt.Fprintln(cmd.OutOrStdout(), tui.SuccessMessage(fmt.Sprintf("Deleted project %q", slug)))
			return nil
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "skip confirmation prompt")

	return cmd
}
