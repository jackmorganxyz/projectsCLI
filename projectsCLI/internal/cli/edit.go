package cli

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/jackmorganxyz/projectsCLI/internal/project"
	"github.com/spf13/cobra"
)

// openFile opens a file with the OS default application.
func openFile(path string) error {
	var cmd string
	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
	case "windows":
		cmd = "start"
	default: // linux, freebsd, etc.
		cmd = "xdg-open"
	}
	return exec.Command(cmd, path).Start()
}

// NewEditCmd opens a project's PROJECT.md in the OS default application.
func NewEditCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit <slug>",
		Short: "Edit a project's PROJECT.md",
		Long:  "Open a project's PROJECT.md in the default application for your OS (e.g. TextEdit on macOS).",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rt, ok := RuntimeFromContext(cmd.Context())
			if !ok {
				return fmt.Errorf("missing runtime context")
			}

			slug := args[0]
			proj, err := project.FindProject(rt.Config.ProjectsDir, slug)
			if err != nil {
				return err
			}

			projectFile := filepath.Join(proj.Dir, "PROJECT.md")
			return openFile(projectFile)
		},
	}

	return cmd
}
