package cli

import (
	"fmt"

	"github.com/jackmorganxyz/projectsCLI/internal/project"
	"github.com/jackmorganxyz/projectsCLI/internal/tui"
	"github.com/spf13/cobra"
)

// NewListCmd lists all projects.
func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all projects",
		RunE: func(cmd *cobra.Command, _ []string) error {
			runtime, ok := RuntimeFromContext(cmd.Context())
			if !ok {
				return fmt.Errorf("missing runtime context")
			}

			projects, err := project.ListProjects(runtime.Config.ProjectsDir)
			if err != nil {
				return err
			}

			if tui.IsJSON() {
				return writeJSON(cmd.OutOrStdout(), projects)
			}

			if len(projects) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), tui.Muted(tui.RandomEmptyState()))
				return nil
			}

			// If interactive, launch the dashboard TUI.
			if tui.IsInteractive() {
				return runDashboard(projects)
			}

			// Plain text table.
			headers := []string{"Slug", "Title", "Status", "Created"}
			var rows [][]string
			for _, p := range projects {
				created := p.Meta.CreatedAt
				if len(created) > 10 {
					created = created[:10]
				}
				rows = append(rows, []string{
					p.Meta.Slug,
					p.Meta.Title,
					p.Meta.Status,
					created,
				})
			}

			fmt.Fprintln(cmd.OutOrStdout(), tui.Table(headers, rows))
			return nil
		},
	}

	return cmd
}

// runDashboard launches the interactive dashboard TUI.
func runDashboard(projects []*project.Project) error {
	m := tui.NewDashboardModel(projects)
	_, err := tui.RunProgram(m)
	return err
}
