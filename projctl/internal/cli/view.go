package cli

import (
	"fmt"
	"strings"

	"github.com/jackpmorgan/projctl/internal/project"
	"github.com/jackpmorgan/projctl/internal/tui"
	"github.com/spf13/cobra"
)

// NewViewCmd displays project details.
func NewViewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view <slug>",
		Short: "View project details",
		Long:  "Display project metadata and content. Launches scrollable TUI in interactive mode.",
		Args:  cobra.ExactArgs(1),
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

			if tui.IsJSON() {
				return writeJSON(cmd.OutOrStdout(), proj)
			}

			// If interactive, launch the detail TUI.
			if tui.IsInteractive() {
				m := tui.NewDetailModel(proj)
				_, err := tui.RunProgram(m)
				return err
			}

			// Plain text view.
			w := cmd.OutOrStdout()
			fmt.Fprintln(w, tui.Header(proj.Meta.Title))
			fmt.Fprintln(w)
			fmt.Fprintln(w, tui.FormatField("Slug", proj.Meta.Slug))
			fmt.Fprintln(w, tui.FormatField("Status", tui.StatusColor(proj.Meta.Status)))
			fmt.Fprintln(w, tui.FormatField("Created", proj.Meta.CreatedAt))
			fmt.Fprintln(w, tui.FormatField("Updated", proj.Meta.UpdatedAt))
			if proj.Meta.Description != "" {
				fmt.Fprintln(w, tui.FormatField("Description", proj.Meta.Description))
			}
			if len(proj.Meta.Tags) > 0 {
				fmt.Fprintln(w, tui.FormatField("Tags", strings.Join(proj.Meta.Tags, ", ")))
			}
			if proj.Meta.GitRemote != "" {
				fmt.Fprintln(w, tui.FormatField("Remote", proj.Meta.GitRemote))
			}
			fmt.Fprintln(w, tui.FormatField("Directory", proj.Dir))

			if proj.Body != "" {
				fmt.Fprintln(w)
				fmt.Fprintln(w, proj.Body)
			}
			return nil
		},
	}

	return cmd
}
