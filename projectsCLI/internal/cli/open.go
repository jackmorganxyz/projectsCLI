package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewOpenCmd opens a project's directory in the OS file manager.
func NewOpenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open <slug>",
		Short: "Open a project folder in the file manager",
		Long:  "Open a project's directory in Finder (macOS), Explorer (Windows), or the default file manager (Linux).",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rt, ok := RuntimeFromContext(cmd.Context())
			if !ok {
				return fmt.Errorf("missing runtime context")
			}

			slug := args[0]
			proj, err := findProject(rt.Config, slug, rt.Folder)
			if err != nil {
				return err
			}

			return openFile(proj.Dir)
		},
	}

	return cmd
}
