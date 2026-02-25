package cli

import (
	"fmt"
	"path/filepath"

	"github.com/jackpmorgan/projctl/internal/git"
	"github.com/jackpmorgan/projctl/internal/project"
	"github.com/jackpmorgan/projctl/internal/tui"
	"github.com/spf13/cobra"
)

// projectHealth holds health check results for a project.
type projectHealth struct {
	Slug         string `json:"slug"`
	Title        string `json:"title"`
	Status       string `json:"status"`
	HasGit       bool   `json:"has_git"`
	HasRemote    bool   `json:"has_remote"`
	Uncommitted  bool   `json:"uncommitted"`
	HasProjectMD bool   `json:"has_project_md"`
}

// NewStatusCmd shows health check across all projects.
func NewStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show project health check",
		Long:  "Display a health summary of all projects including git status and file integrity.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			runtime, ok := RuntimeFromContext(cmd.Context())
			if !ok {
				return fmt.Errorf("missing runtime context")
			}

			projects, err := project.ListProjects(runtime.Config.ProjectsDir)
			if err != nil {
				return err
			}

			if len(projects) == 0 {
				if tui.IsJSON() {
					return writeJSON(cmd.OutOrStdout(), []projectHealth{})
				}
				fmt.Fprintln(cmd.OutOrStdout(), tui.Muted(tui.RandomEmptyState()))
				return nil
			}

			var health []projectHealth
			for _, p := range projects {
				dir := filepath.Join(runtime.Config.ProjectsDir, p.Meta.Slug)
				h := projectHealth{
					Slug:         p.Meta.Slug,
					Title:        p.Meta.Title,
					Status:       p.Meta.Status,
					HasGit:       git.IsRepo(dir),
					HasProjectMD: true,
				}

				if h.HasGit {
					h.HasRemote = git.HasRemote(dir)
					uncommitted, _ := git.HasUncommitted(dir)
					h.Uncommitted = uncommitted
				}

				health = append(health, h)
			}

			if tui.IsJSON() {
				return writeJSON(cmd.OutOrStdout(), health)
			}

			headers := []string{"Slug", "Status", "Git", "Remote", "Clean"}
			var rows [][]string
			for _, h := range health {
				gitStatus := tui.Muted("no")
				if h.HasGit {
					gitStatus = tui.SuccessMessage("yes")
				}

				remoteStatus := tui.Muted("-")
				if h.HasGit && h.HasRemote {
					remoteStatus = tui.SuccessMessage("yes")
				} else if h.HasGit {
					remoteStatus = tui.WarningMessage("no")
				}

				cleanStatus := tui.Muted("-")
				if h.HasGit && !h.Uncommitted {
					cleanStatus = tui.SuccessMessage("clean")
				} else if h.HasGit && h.Uncommitted {
					cleanStatus = tui.WarningMessage("dirty")
				}

				rows = append(rows, []string{
					h.Slug,
					tui.StatusColor(h.Status),
					gitStatus,
					remoteStatus,
					cleanStatus,
				})
			}

			fmt.Fprintln(cmd.OutOrStdout(), tui.Header("🩺 Project Health"))
			fmt.Fprintln(cmd.OutOrStdout())
			fmt.Fprintln(cmd.OutOrStdout(), tui.Table(headers, rows))
			return nil
		},
	}

	return cmd
}
