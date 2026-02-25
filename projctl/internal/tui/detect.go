package tui

import (
	"os"
	"sync/atomic"

	"github.com/mattn/go-isatty"
)

var jsonOverride atomic.Bool

// SetJSON explicitly toggles JSON output mode.
func SetJSON(enabled bool) {
	jsonOverride.Store(enabled)
}

// IsInteractive reports whether stdin/stdout are attached to a terminal.
func IsInteractive() bool {
	stdinTTY := isatty.IsTerminal(os.Stdin.Fd()) || isatty.IsCygwinTerminal(os.Stdin.Fd())
	stdoutTTY := isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
	return stdinTTY && stdoutTTY
}

// IsJSON reports whether command output should be JSON.
//
// JSON mode is enabled when:
//  1. --json is set, or
//  2. output is non-interactive (stdout is not a TTY).
func IsJSON() bool {
	if jsonOverride.Load() {
		return true
	}

	for _, arg := range os.Args[1:] {
		switch {
		case arg == "--json", arg == "--json=true":
			return true
		case arg == "--json=false":
			return false
		}
	}

	stdoutTTY := isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
	return !stdoutTTY
}
