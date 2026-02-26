package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/jackmorganxyz/projectsCLI/internal/project"
	"github.com/jackmorganxyz/projectsCLI/internal/tui"
	"github.com/spf13/cobra"
)

// NewUpdateCmd updates project metadata.
func NewUpdateCmd() *cobra.Command {
	var (
		title       string
		description string
		status      string
		tags        string
	)

	cmd := &cobra.Command{
		Use:   "update <slug>",
		Short: "Update project metadata",
		Long: `Update project metadata including title, description, status, and tags.

Use flags to update specific fields. The updated_at timestamp is automatically set.`,
		Args: cobra.ExactArgs(1),
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

			// Update fields if provided
			updated := false
			if title != "" {
				proj.Meta.Title = title
				updated = true
			}
			if description != "" {
				proj.Meta.Description = description
				updated = true
			}
			if status != "" {
				// Validate status
				if status != "active" && status != "paused" && status != "archived" {
					return fmt.Errorf("invalid status %q: must be active, paused, or archived", status)
				}
				proj.Meta.Status = status
				updated = true
			}
			if tags != "" {
				proj.Meta.Tags = strings.Split(tags, ",")
				for i := range proj.Meta.Tags {
					proj.Meta.Tags[i] = strings.TrimSpace(proj.Meta.Tags[i])
				}
				updated = true
			}

			if !updated {
				return fmt.Errorf("no fields to update. Use --title, --description, --status, or --tags")
			}

			// Update timestamp
			proj.Meta.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

			// Write updated project file
			if err := project.WriteProjectFile(proj.Dir, proj.Meta, proj.Body); err != nil {
				return fmt.Errorf("write project file: %w", err)
			}

			// Regenerate registry
			_ = project.WriteRegistry(runtime.Config.ProjectsDir)

			if tui.IsJSON() {
				return writeJSON(cmd.OutOrStdout(), map[string]any{
					"status":     "updated",
					"slug":       proj.Meta.Slug,
					"updated_at": proj.Meta.UpdatedAt,
				})
			}

			w := cmd.OutOrStdout()
			fmt.Fprintln(w, tui.SuccessMessage(fmt.Sprintf("Updated project %q", slug)))
			if title != "" {
				fmt.Fprintln(w, tui.FormatField("Title", proj.Meta.Title))
			}
			if description != "" {
				fmt.Fprintln(w, tui.FormatField("Description", proj.Meta.Description))
			}
			if status != "" {
				fmt.Fprintln(w, tui.FormatField("Status", proj.Meta.Status))
			}
			if tags != "" {
				fmt.Fprintln(w, tui.FormatField("Tags", strings.Join(proj.Meta.Tags, ", ")))
			}
			fmt.Fprintln(w, tui.FormatField("Updated", proj.Meta.UpdatedAt))
			return nil
		},
	}

	cmd.Flags().StringVar(&title, "title", "", "update title")
	cmd.Flags().StringVar(&description, "description", "", "update description")
	cmd.Flags().StringVar(&status, "status", "", "update status (active/paused/archived)")
	cmd.Flags().StringVar(&tags, "tags", "", "update tags (comma-separated)")

	return cmd
}
