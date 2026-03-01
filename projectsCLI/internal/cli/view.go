package cli

import (
	"fmt"

	"github.com/jackmorganxyz/projectsCLI/internal/tui"
	"github.com/spf13/cobra"
)

// NewViewCmd displays project details.
func NewViewCmd() *cobra.Command {
	var field string

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
			proj, err := findProject(runtime.Config, slug, runtime.Folder)
			if err != nil {
				return err
			}

			// Handle --field flag for field extraction
			if field != "" {
				val, err := extractField(proj, field)
				if err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), val)
				return nil
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
			fmt.Fprintln(w, tui.Header("📂 "+proj.Meta.Title))
			fmt.Fprintln(w, tui.Muted("  "+tui.RandomViewGreeting()+" "+tui.Slug(proj.Meta.Slug)))
			fmt.Fprintln(w)
			fmt.Fprintln(w, tui.FormatField("Slug", tui.Slug(proj.Meta.Slug)))
			fmt.Fprintln(w, tui.FormatField("Status", tui.StatusEmoji(proj.Meta.Status)+tui.StatusColor(proj.Meta.Status)))
			fmt.Fprintln(w, tui.FormatField("Created", proj.Meta.CreatedAt))
			fmt.Fprintln(w, tui.FormatField("Updated", proj.Meta.UpdatedAt))
			if proj.Meta.Description != "" {
				fmt.Fprintln(w, tui.FormatField("Description", proj.Meta.Description))
			}
			if len(proj.Meta.Tags) > 0 {
				fmt.Fprintln(w, tui.FormatField("Tags", tui.TagList(proj.Meta.Tags)))
			}
			if proj.Meta.GitRemote != "" {
				fmt.Fprintln(w, tui.FormatField("Remote", tui.Path(proj.Meta.GitRemote)))
			}
			fmt.Fprintln(w, tui.FormatField("Directory", tui.Path(proj.Dir)))

			if proj.Body != "" {
				fmt.Fprintln(w)
				fmt.Fprintln(w, tui.Divider(40))
				fmt.Fprintln(w, proj.Body)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&field, "field", "", "extract specific field from JSON output (e.g. --field dir, --field meta.title)")

	return cmd
}
