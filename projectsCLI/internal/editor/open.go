package editor

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// Open launches the given editor to open filePath.
// Terminal editors run in the foreground (blocking); GUI editors start in the background.
func Open(ed Editor, filePath string) error {
	switch ed.Type {
	case Terminal:
		cmd := exec.Command(ed.Command, filePath)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()

	case GUI:
		if runtime.GOOS == "darwin" {
			return exec.Command("open", "-a", ed.Name, filePath).Start()
		}
		return exec.Command(ed.Command, filePath).Start()

	default:
		return fmt.Errorf("unknown editor type %d", ed.Type)
	}
}

// OpenByCommand finds the editor matching command among detected editors
// and opens filePath with it. If no match is found, it falls back to a
// terminal-style launch (blocking with stdin/stdout attached).
func OpenByCommand(command, filePath string) error {
	editors := Detect()
	if ed := FindByCommand(editors, command); ed != nil {
		return Open(*ed, filePath)
	}

	// Fallback: treat unknown command as a terminal editor.
	cmd := exec.Command(command, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
