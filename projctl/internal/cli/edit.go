package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jackpmorgan/projctl/internal/project"
	"github.com/spf13/cobra"
)

// NewEditCmd opens a project's PROJECT.md in the configured editor.
func NewEditCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit <slug>",
		Short: "Edit a project's PROJECT.md",
		Long:  "Open a project's PROJECT.md in your configured editor ($EDITOR or config).",
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

			editor := runtime.Config.Editor
			projectFile := filepath.Join(proj.Dir, "PROJECT.md")

			editorCmd := exec.Command(editor, projectFile)
			editorCmd.Stdin = os.Stdin
			editorCmd.Stdout = os.Stdout
			editorCmd.Stderr = os.Stderr

			return editorCmd.Run()
		},
	}

	return cmd
}
