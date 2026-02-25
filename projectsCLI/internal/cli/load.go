package cli

import (
	"fmt"
	"strings"

	"github.com/jackpmorgan/projectsCLI/internal/project"
	"github.com/jackpmorgan/projectsCLI/internal/tui"
	"github.com/spf13/cobra"
)

// NewLoadCmd outputs project data for agent consumption.
func NewLoadCmd() *cobra.Command {
	var (
		export bool
		bash   bool
	)

	cmd := &cobra.Command{
		Use:   "load <slug>",
		Short: "Load project data for agents",
		Long:  "Output project metadata for agent consumption.\nUse --export for shell variable exports, --bash for eval-able script, or --json for structured data.",
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

			if tui.IsJSON() && !export && !bash {
				return writeJSON(cmd.OutOrStdout(), proj)
			}

			if bash {
				return writeBashVars(cmd, proj)
			}

			if export {
				return writeExports(cmd, proj)
			}

			// Default: JSON.
			return writeJSON(cmd.OutOrStdout(), proj)
		},
	}

	cmd.Flags().BoolVar(&export, "export", false, "output as shell export statements")
	cmd.Flags().BoolVar(&bash, "bash", false, "output as eval-able bash variables")

	return cmd
}

func writeExports(cmd *cobra.Command, proj *project.Project) error {
	w := cmd.OutOrStdout()
	fmt.Fprintf(w, "export PROJECT_SLUG=%q\n", proj.Meta.Slug)
	fmt.Fprintf(w, "export PROJECT_TITLE=%q\n", proj.Meta.Title)
	fmt.Fprintf(w, "export PROJECT_STATUS=%q\n", proj.Meta.Status)
	fmt.Fprintf(w, "export PROJECT_DIR=%q\n", proj.Dir)
	fmt.Fprintf(w, "export PROJECT_DESCRIPTION=%q\n", proj.Meta.Description)
	if len(proj.Meta.Tags) > 0 {
		fmt.Fprintf(w, "export PROJECT_TAGS=%q\n", strings.Join(proj.Meta.Tags, ","))
	}
	if proj.Meta.GitRemote != "" {
		fmt.Fprintf(w, "export PROJECT_GIT_REMOTE=%q\n", proj.Meta.GitRemote)
	}
	return nil
}

func writeBashVars(cmd *cobra.Command, proj *project.Project) error {
	w := cmd.OutOrStdout()
	fmt.Fprintf(w, "PROJECT_SLUG=%q\n", proj.Meta.Slug)
	fmt.Fprintf(w, "PROJECT_TITLE=%q\n", proj.Meta.Title)
	fmt.Fprintf(w, "PROJECT_STATUS=%q\n", proj.Meta.Status)
	fmt.Fprintf(w, "PROJECT_DIR=%q\n", proj.Dir)
	fmt.Fprintf(w, "PROJECT_DESCRIPTION=%q\n", proj.Meta.Description)
	if len(proj.Meta.Tags) > 0 {
		fmt.Fprintf(w, "PROJECT_TAGS=%q\n", strings.Join(proj.Meta.Tags, ","))
	}
	if proj.Meta.GitRemote != "" {
		fmt.Fprintf(w, "PROJECT_GIT_REMOTE=%q\n", proj.Meta.GitRemote)
	}
	return nil
}
