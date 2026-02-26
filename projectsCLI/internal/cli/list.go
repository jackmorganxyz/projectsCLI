package cli

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/jackmorganxyz/projectsCLI/internal/project"
	"github.com/jackmorganxyz/projectsCLI/internal/tui"
	"github.com/spf13/cobra"
)

// NewListCmd lists all projects.
func NewListCmd() *cobra.Command {
	var field string

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all projects",
		RunE: func(cmd *cobra.Command, _ []string) error {
			runtime, ok := RuntimeFromContext(cmd.Context())
			if !ok {
				return fmt.Errorf("missing runtime context")
			}

			projects, err := listAllProjects(runtime.Config, runtime.Folder)
			if err != nil {
				return err
			}

			// Handle --field flag for field extraction
			if field != "" {
				if len(projects) == 0 {
					return nil
				}
				for _, p := range projects {
					val, err := extractField(p, field)
					if err != nil {
						return err
					}
					fmt.Fprintln(cmd.OutOrStdout(), val)
				}
				return nil
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
				return runDashboard(cmd, projects)
			}

			// Plain text table — include Folder column if folders are configured.
			hasFolders := len(runtime.Config.Folders) > 0
			if hasFolders {
				headers := []string{"Slug", "Folder", "Title", "Status", "Created"}
				var rows [][]string
				for _, p := range projects {
					created := p.Meta.CreatedAt
					if len(created) > 10 {
						created = created[:10]
					}
					folderDisplay := p.Folder
					if folderDisplay == "" {
						folderDisplay = "-"
					}
					rows = append(rows, []string{
						p.Meta.Slug,
						folderDisplay,
						p.Meta.Title,
						p.Meta.Status,
						created,
					})
				}
				fmt.Fprintln(cmd.OutOrStdout(), tui.Table(headers, rows))
			} else {
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
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&field, "field", "", "extract specific field from JSON output (e.g. --field dir, --field meta.title)")

	return cmd
}

// runDashboard launches the interactive dashboard TUI and, if a project is
// selected, shows a command picker and executes the chosen command.
func runDashboard(cmd *cobra.Command, projects []*project.Project) error {
	m := tui.NewDashboardModel(projects)
	finalModel, err := tui.RunProgram(m)
	if err != nil {
		return err
	}

	dm := finalModel.(tui.DashboardModel)
	if !dm.WasSelected() {
		return nil
	}

	slug := dm.SelectedSlug()
	if slug == "" {
		return nil
	}

	command, err := pickProjectCommand(slug)
	if err != nil || command == "" {
		return err
	}

	root := cmd.Root()
	root.SetArgs([]string{command, slug})
	return root.Execute()
}

// pickProjectCommand shows a command dropdown for the selected project.
func pickProjectCommand(slug string) (string, error) {
	type commandOption struct {
		label   string
		command string
	}

	commands := []commandOption{
		{"View details", "view"},
		{"Edit a file", "edit"},
		{"Open in file manager", "open"},
		{"Check status", "status"},
		{"Push to remote", "push"},
		{"Update metadata", "update"},
		{"Move to folder", "move"},
		{"Delete project", "delete"},
	}

	options := make([]huh.Option[string], len(commands))
	for i, c := range commands {
		options[i] = huh.NewOption(c.label, c.command)
	}

	var selected string
	theme := huh.ThemeBase()
	theme.Focused.Title = theme.Focused.Title.Foreground(lipgloss.Color(tui.ColorPrimary))
	theme.Focused.SelectSelector = theme.Focused.SelectSelector.Foreground(lipgloss.Color(tui.ColorPrimary))

	err := huh.NewSelect[string]().
		Title(fmt.Sprintf("What would you like to do with %q?", slug)).
		Options(options...).
		Value(&selected).
		WithTheme(theme).
		Run()

	if err != nil {
		return "", nil // Ctrl+C / Esc
	}
	return selected, nil
}
