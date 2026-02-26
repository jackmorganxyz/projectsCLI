package editor

import (
	"os/exec"
	"runtime"
	"strings"
)

// EditorType distinguishes terminal-based editors from GUI applications.
type EditorType int

const (
	Terminal EditorType = iota
	GUI
)

// Editor represents an installed text editor.
type Editor struct {
	Name    string     // Display name (e.g. "VS Code")
	Command string     // Launch key used for config/lookup (e.g. "code")
	Type    EditorType // Terminal or GUI
}

// candidate describes an editor to probe for during detection.
type candidate struct {
	Name    string
	Command string
	Type    EditorType
}

// Detect returns a list of editors actually installed on the current system.
func Detect() []Editor {
	switch runtime.GOOS {
	case "darwin":
		return detectDarwin()
	case "windows":
		return detectWindows()
	default: // linux, freebsd, etc.
		return detectLinux()
	}
}

// FindByCommand returns the first editor whose Command matches cmd, or nil.
func FindByCommand(editors []Editor, cmd string) *Editor {
	for i := range editors {
		if editors[i].Command == cmd {
			return &editors[i]
		}
	}
	return nil
}

// HasCLI reports whether a CLI program is available on PATH.
func HasCLI(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// IsInstalled reports whether an editor is available on the system.
// For CLI editors this checks PATH. For GUI editors on macOS it also
// checks via Spotlight (mdfind) since GUI apps aren't on PATH.
func IsInstalled(command string) bool {
	if HasCLI(command) {
		return true
	}
	if runtime.GOOS == "darwin" {
		for _, app := range guiAppsDarwin {
			if app.Command == command {
				out, err := exec.Command("mdfind", "kMDItemCFBundleIdentifier", "=", app.BundleID).Output()
				return err == nil && len(strings.TrimSpace(string(out))) > 0
			}
		}
	}
	return false
}

// --- macOS ---

// guiAppsDarwin maps macOS .app bundle identifiers to editor metadata.
var guiAppsDarwin = []struct {
	BundleID string
	Name     string
	Command  string
}{
	{"com.todesktop.230313mzl4w4u92", "Cursor", "cursor"},
	{"com.microsoft.VSCode", "VS Code", "code"},
	{"com.sublimetext.4", "Sublime Text", "subl"},
	{"com.sublimetext.3", "Sublime Text", "subl"},
	{"com.barebones.bbedit", "BBEdit", "bbedit"},
	{"dev.zed.Zed", "Zed", "zed"},
	{"com.apple.TextEdit", "TextEdit", "textedit"},
}

func detectDarwin() []Editor {
	var found []Editor

	// Probe GUI apps via mdfind (Spotlight).
	for _, app := range guiAppsDarwin {
		out, err := exec.Command("mdfind", "kMDItemCFBundleIdentifier", "=", app.BundleID).Output()
		if err == nil && len(strings.TrimSpace(string(out))) > 0 {
			// Skip duplicates (Sublime 3 vs 4 both map to "subl").
			if FindByCommand(found, app.Command) != nil {
				continue
			}
			found = append(found, Editor{Name: app.Name, Command: app.Command, Type: GUI})
		}
	}

	// Probe CLI editors.
	found = append(found, detectCLIEditors()...)
	return found
}

// --- Linux ---

var guiAppsLinux = []candidate{
	{"VS Code", "code", GUI},
	{"Cursor", "cursor", GUI},
	{"Sublime Text", "subl", GUI},
	{"Zed", "zed", GUI},
	{"Kate", "kate", GUI},
	{"gedit", "gedit", GUI},
	{"GNOME Text Editor", "gnome-text-editor", GUI},
}

func detectLinux() []Editor {
	var found []Editor
	for _, app := range guiAppsLinux {
		if HasCLI(app.Command) {
			found = append(found, Editor{Name: app.Name, Command: app.Command, Type: app.Type})
		}
	}
	found = append(found, detectCLIEditors()...)
	return found
}

// --- Windows ---

var guiAppsWindows = []candidate{
	{"VS Code", "code", GUI},
	{"Cursor", "cursor", GUI},
	{"Sublime Text", "subl", GUI},
	{"Notepad++", "notepad++", GUI},
	{"Notepad", "notepad", GUI},
}

func detectWindows() []Editor {
	var found []Editor
	for _, app := range guiAppsWindows {
		if HasCLI(app.Command) {
			found = append(found, Editor{Name: app.Name, Command: app.Command, Type: app.Type})
		}
	}
	found = append(found, detectCLIEditors()...)
	return found
}

// --- Shared CLI editors ---

var cliEditors = []candidate{
	{"Neovim", "nvim", Terminal},
	{"Vim", "vim", Terminal},
	{"Nano", "nano", Terminal},
	{"Emacs", "emacs", Terminal},
	{"Micro", "micro", Terminal},
	{"Helix", "hx", Terminal},
}

func detectCLIEditors() []Editor {
	var found []Editor
	for _, ed := range cliEditors {
		if HasCLI(ed.Command) {
			found = append(found, Editor{Name: ed.Name, Command: ed.Command, Type: ed.Type})
		}
	}
	return found
}
