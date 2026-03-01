package cli

import (
	"fmt"

	"github.com/jackmorganxyz/projectsCLI/internal/git"
	"github.com/jackmorganxyz/projectsCLI/internal/tui"
	"github.com/spf13/cobra"
)

// projectHealth holds health check results for a project.
type projectHealth struct {
	Slug         string `json:"slug"`
	Folder       string `json:"folder,omitempty"`
	Title        string `json:"title"`
	Status       string `json:"status"`
	HasGit       bool   `json:"has_git"`
	HasRemote    bool   `json:"has_remote"`
	Uncommitted  bool   `json:"uncommitted"`
	HasProjectMD bool   `json:"has_project_md"`
}

// NewStatusCmd shows health check across all projects.
func NewStatusCmd() *cobra.Command {
	var field string

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show project health check",
		Long:  "Display a health summary of all projects including git status and file integrity.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			runtime, ok := RuntimeFromContext(cmd.Context())
			if !ok {
				return fmt.Errorf("missing runtime context")
			}

			projects, err := listAllProjects(runtime.Config, runtime.Folder)
			if err != nil {
				return err
			}

			if len(projects) == 0 {
				if tui.IsJSON() || field != "" {
					return writeJSON(cmd.OutOrStdout(), []projectHealth{})
				}
				fmt.Fprintln(cmd.OutOrStdout(), tui.Muted(tui.RandomEmptyState()))
				return nil
			}

			var health []projectHealth
			for _, p := range projects {
				h := projectHealth{
					Slug:         p.Meta.Slug,
					Folder:       p.Folder,
					Title:        p.Meta.Title,
					Status:       p.Meta.Status,
					HasGit:       git.IsRepo(p.Dir),
					HasProjectMD: true,
				}

				if h.HasGit {
					h.HasRemote = git.HasRemote(p.Dir)
					uncommitted, _ := git.HasUncommitted(p.Dir)
					h.Uncommitted = uncommitted
				}

				health = append(health, h)
			}

			// Handle --field flag for field extraction
			if field != "" {
				for _, h := range health {
					val, err := extractField(h, field)
					if err != nil {
						return err
					}
					fmt.Fprintln(cmd.OutOrStdout(), val)
				}
				return nil
			}

			if tui.IsJSON() {
				return writeJSON(cmd.OutOrStdout(), health)
			}

			hasFolders := len(runtime.Config.Folders) > 0
			if hasFolders {
				headers := []string{"Slug", "Folder", "Status", "Git", "Remote", "Clean"}
				var rows [][]string
				for _, h := range health {
					folderDisplay := h.Folder
					if folderDisplay == "" {
						folderDisplay = "-"
					}
					rows = append(rows, []string{
						h.Slug,
						folderDisplay,
						tui.StatusColor(h.Status),
						gitIcon(h.HasGit),
						remoteIcon(h.HasGit, h.HasRemote),
						cleanIcon(h.HasGit, h.Uncommitted),
					})
				}
				fmt.Fprintln(cmd.OutOrStdout(), tui.Header(tui.RandomStatusHeader()))
				fmt.Fprintln(cmd.OutOrStdout())
				fmt.Fprintln(cmd.OutOrStdout(), tui.Table(headers, rows))
			} else {
				headers := []string{"Slug", "Status", "Git", "Remote", "Clean"}
				var rows [][]string
				for _, h := range health {
					rows = append(rows, []string{
						h.Slug,
						tui.StatusColor(h.Status),
						gitIcon(h.HasGit),
						remoteIcon(h.HasGit, h.HasRemote),
						cleanIcon(h.HasGit, h.Uncommitted),
					})
				}
				fmt.Fprintln(cmd.OutOrStdout(), tui.Header(tui.RandomStatusHeader()))
				fmt.Fprintln(cmd.OutOrStdout())
				fmt.Fprintln(cmd.OutOrStdout(), tui.Table(headers, rows))
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&field, "field", "", "extract specific field from JSON output (e.g. --field slug, --field status)")

	return cmd
}

func gitIcon(hasGit bool) string {
	if hasGit {
		return tui.SuccessMessage("yes")
	}
	return tui.Muted("no")
}

func remoteIcon(hasGit, hasRemote bool) string {
	if hasGit && hasRemote {
		return tui.SuccessMessage("yes")
	}
	if hasGit {
		return tui.WarningMessage("no")
	}
	return tui.Muted("-")
}

func cleanIcon(hasGit, uncommitted bool) string {
	if hasGit && !uncommitted {
		return tui.SuccessMessage("clean")
	}
	if hasGit && uncommitted {
		return tui.WarningMessage("dirty")
	}
	return tui.Muted("-")
}
