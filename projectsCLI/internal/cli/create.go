package cli

import (
	"fmt"
	"time"

	"github.com/jackmorganxyz/projectsCLI/internal/git"
	"github.com/jackmorganxyz/projectsCLI/internal/project"
	"github.com/jackmorganxyz/projectsCLI/internal/tui"
	"github.com/spf13/cobra"
)

// NewCreateCmd creates a new project scaffold.
func NewCreateCmd() *cobra.Command {
	var (
		title       string
		description string
		tags        []string
		status      string
	)

	cmd := &cobra.Command{
		Use:   "create [slug]",
		Short: "Create a new project",
		Long: `Create a new project scaffold with directory structure and template files.

If slug is omitted, it is generated from --title.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			runtime, ok := RuntimeFromContext(cmd.Context())
			if !ok {
				return fmt.Errorf("missing runtime context")
			}

			var slug string
			if len(args) == 1 {
				slug = args[0]
			} else if title != "" {
				slug = Slugify(title)
				if slug == "" {
					return fmt.Errorf("could not generate a valid slug from title %q", title)
				}
			} else {
				return fmt.Errorf("provide a slug argument or --title to auto-generate one")
			}

			if err := ValidateSlug(slug); err != nil {
				return err
			}

			meta := project.NewMeta(slug, title)
			if description != "" {
				meta.Description = description
			}
			if len(tags) > 0 {
				meta.Tags = tags
			}
			if status != "" {
				meta.Status = status
			}

			dir, err := project.Scaffold(runtime.Config.ProjectsDir, meta)
			if err != nil {
				return err
			}

			// Auto-init git if configured.
			if runtime.Config.AutoGitInit {
				if err := git.Init(dir); err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "warning: git init failed: %v\n", err)
				} else {
					_ = git.AddAll(dir)
					_ = git.Commit(dir, "Initial project scaffold")
				}
			}

			// Regenerate registry.
			_ = project.WriteRegistry(runtime.Config.ProjectsDir)

			if tui.IsJSON() {
				return writeJSON(cmd.OutOrStdout(), map[string]any{
					"status":     "created",
					"slug":       slug,
					"dir":        dir,
					"created_at": meta.CreatedAt,
				})
			}

			w := cmd.OutOrStdout()
			fmt.Fprintln(w, tui.SuccessMessage(fmt.Sprintf("Created project %q — %s", slug, tui.RandomCreateCheer())))
			fmt.Fprintln(w, tui.FormatField("Directory", dir))
			fmt.Fprintln(w, tui.FormatField("Created", time.Now().Format("2006-01-02")))
			if tip := tui.MaybeTip(); tip != "" {
				fmt.Fprintln(w)
				fmt.Fprintln(w, tip)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&title, "title", "", "project title (defaults to slug)")
	cmd.Flags().StringVar(&description, "description", "", "project description")
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "project tags (comma-separated)")
	cmd.Flags().StringVar(&status, "status", "active", "project status")

	return cmd
}
