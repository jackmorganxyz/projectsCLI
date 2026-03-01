package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/jackmorganxyz/projectsCLI/internal/config"
	"github.com/jackmorganxyz/projectsCLI/internal/editor"
	"github.com/jackmorganxyz/projectsCLI/internal/tui"
	"github.com/spf13/cobra"
)

// textExtensions is the set of file extensions shown in the interactive browser.
var textExtensions = map[string]bool{
	".md": true, ".txt": true, ".go": true, ".js": true, ".ts": true,
	".jsx": true, ".tsx": true, ".py": true, ".rs": true, ".rb": true,
	".java": true, ".c": true, ".cpp": true, ".h": true, ".hpp": true,
	".css": true, ".html": true, ".xml": true, ".json": true, ".yaml": true,
	".yml": true, ".toml": true, ".sh": true, ".bash": true, ".zsh": true,
	".fish": true, ".env": true, ".gitignore": true, ".sql": true,
	".swift": true, ".kt": true, ".scala": true, ".lua": true, ".r": true,
	".csv": true, ".ini": true, ".cfg": true, ".conf": true,
	".dockerfile": true, ".makefile": true,
}

const backOption = "← Back"

// NewEditCmd opens a project file in a chosen editor.
func NewEditCmd() *cobra.Command {
	var editorFlag string
	var editorPicker bool

	cmd := &cobra.Command{
		Use:   "edit <slug>",
		Short: "Browse and edit a project file",
		Long: `Interactively browse files in a project directory and open the selected
file in your preferred editor.

On first run you'll be prompted to pick an editor from those installed on
your system. The choice is saved to config so subsequent runs open directly.

Use --editor to override the saved editor for a single invocation.
Use --editor-picker to re-show the editor selection prompt.
In non-interactive mode (piped stdin/stdout) the command defaults to
opening PROJECT.md with the saved editor.`,
		Args: cobra.ExactArgs(1),
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

			// Pick a file.
			var filePath string
			if tui.IsInteractive() {
				filePath, err = browseFiles(proj.Dir)
				if err != nil {
					return err
				}
				if filePath == "" {
					return nil // user cancelled
				}
			} else {
				filePath = filepath.Join(proj.Dir, "PROJECT.md")
			}

			// Resolve editor.
			editorCmd, err := resolveEditor(cmd, rt, editorFlag, editorPicker)
			if err != nil {
				return err
			}

			return editor.OpenByCommand(editorCmd, filePath)
		},
	}

	cmd.Flags().StringVar(&editorFlag, "editor", "", "editor command to use (bypasses picker)")
	cmd.Flags().BoolVar(&editorPicker, "editor-picker", false, "force the editor picker (ignore saved preference)")

	return cmd
}

// browseFiles presents an interactive file browser rooted at projectDir.
// Returns the selected file's absolute path, or "" if the user cancelled.
func browseFiles(projectDir string) (string, error) {
	currentDir := projectDir

	for {
		entries, err := os.ReadDir(currentDir)
		if err != nil {
			return "", fmt.Errorf("read directory: %w", err)
		}

		var dirs, files []os.DirEntry
		for _, e := range entries {
			name := e.Name()
			// Skip hidden files/dirs.
			if strings.HasPrefix(name, ".") {
				continue
			}
			if e.IsDir() {
				dirs = append(dirs, e)
			} else if isTextFile(e) {
				files = append(files, e)
			}
		}

		sort.Slice(dirs, func(i, j int) bool { return dirs[i].Name() < dirs[j].Name() })
		sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })

		var options []huh.Option[string]

		// Add back option when not at project root.
		if currentDir != projectDir {
			options = append(options, huh.NewOption(backOption, backOption))
		}

		for _, d := range dirs {
			label := "📁 " + d.Name() + "/"
			options = append(options, huh.NewOption(label, filepath.Join(currentDir, d.Name())))
		}
		for _, f := range files {
			label := "📄 " + f.Name()
			options = append(options, huh.NewOption(label, filepath.Join(currentDir, f.Name())))
		}

		if len(options) == 0 {
			return "", fmt.Errorf("no browsable files in %s", currentDir)
		}

		var selected string
		theme := huh.ThemeBase()
		theme.Focused.Title = theme.Focused.Title.Foreground(lipgloss.Color(tui.ColorPrimary))
		theme.Focused.SelectSelector = theme.Focused.SelectSelector.Foreground(lipgloss.Color(tui.ColorPrimary))

		err = huh.NewSelect[string]().
			Title("Choose a file").
			Options(options...).
			Value(&selected).
			WithTheme(theme).
			Run()

		if err != nil {
			// Ctrl+C / Esc
			return "", nil
		}

		if selected == backOption {
			currentDir = filepath.Dir(currentDir)
			continue
		}

		info, err := os.Stat(selected)
		if err != nil {
			return "", err
		}
		if info.IsDir() {
			currentDir = selected
			continue
		}

		return selected, nil
	}
}

// isTextFile reports whether a directory entry looks like a text file we should show.
func isTextFile(e os.DirEntry) bool {
	ext := strings.ToLower(filepath.Ext(e.Name()))
	if textExtensions[ext] {
		return true
	}
	// Files with no extension that are small (< 1MB) — likely Makefile, Dockerfile, etc.
	if ext == "" {
		info, err := e.Info()
		if err != nil {
			return false
		}
		return info.Size() < 1<<20
	}
	return false
}

// resolveEditor determines which editor command to use.
// Priority: --editor flag > saved config (if still installed) > interactive picker.
func resolveEditor(cmd *cobra.Command, rt RuntimeContext, flagValue string, forcePicker bool) (string, error) {
	if flagValue != "" {
		return flagValue, nil
	}

	if !tui.IsInteractive() {
		// Non-interactive: use config value (which defaults to $EDITOR or vim).
		if rt.Config.Editor != "" {
			return rt.Config.Editor, nil
		}
		return "vim", nil
	}

	// Interactive: use saved editor if it's still installed and picker isn't forced.
	if !forcePicker && rt.Config.Editor != "" && editor.IsInstalled(rt.Config.Editor) {
		return rt.Config.Editor, nil
	}

	// Detect installed editors and prompt.
	editors := editor.Detect()
	if len(editors) == 0 {
		return "", fmt.Errorf("no editors detected on this system")
	}

	options := make([]huh.Option[string], len(editors))
	for i, ed := range editors {
		label := ed.Name
		if ed.Type == editor.Terminal {
			label += "  (terminal)"
		} else {
			label += "  (GUI)"
		}
		if ed.Command == rt.Config.Editor {
			label += " *"
		}
		options[i] = huh.NewOption(label, ed.Command)
	}

	var selected string
	theme := huh.ThemeBase()
	theme.Focused.Title = theme.Focused.Title.Foreground(lipgloss.Color(tui.ColorPrimary))
	theme.Focused.SelectSelector = theme.Focused.SelectSelector.Foreground(lipgloss.Color(tui.ColorPrimary))

	err := huh.NewSelect[string]().
		Title("Choose an editor").
		Options(options...).
		Value(&selected).
		WithTheme(theme).
		Run()

	if err != nil {
		return "", fmt.Errorf("editor selection cancelled")
	}

	// Save choice to config.
	rt.Config.Editor = selected
	if saveErr := config.SaveToPath(rt.Config, rt.ConfigPath); saveErr != nil {
		fmt.Fprintln(cmd.ErrOrStderr(), tui.WarningMessage(fmt.Sprintf("could not save editor preference: %v", saveErr)))
	}

	return selected, nil
}
