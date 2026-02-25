package config

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Theme colors (mirrored from tui to avoid import cycle).
const (
	colorPrimary = "#8B5CF6"
	colorSuccess = "#34D399"
	colorMuted   = "#9CA3AF"
	colorInfo    = "#60A5FA"
)

var (
	primaryStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorPrimary)).Bold(true)
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorSuccess)).Bold(true)
	mutedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color(colorMuted))
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color(colorInfo)).Bold(true)
	boxStyle     = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(colorPrimary)).
			Padding(0, 1)
)

// NeedsSetup returns true if this is the first run (no config file exists).
func NeedsSetup() bool {
	path, err := ConfigPath()
	if err != nil {
		return true
	}
	expanded, err := expandPath(path)
	if err != nil {
		return true
	}
	_, err = os.Stat(expanded)
	return os.IsNotExist(err)
}

// RunSetup handles the first-run experience. Detects ~/.openclaw and lets the
// user choose where to store projects. Returns the config that was created.
func RunSetup() (Config, error) {
	cfg := Defaults()

	// Always print the welcome banner.
	fmt.Println()
	fmt.Println(boxStyle.Render(
		primaryStyle.Render("projectsCLI") + mutedStyle.Render(" — first time? nice."),
	))
	fmt.Println()

	openclawDir, hasOpenclaw := OpenClawDir()

	if hasOpenclaw {
		// They're an openclaw user — offer the choice.
		fmt.Println(infoStyle.Render("  👀 Well, well, well..."))
		fmt.Println(mutedStyle.Render("  Looks like you've got an openclaw setup at " + openclawDir))
		fmt.Println()
		fmt.Println(mutedStyle.Render("  projectsCLI can store projects inside your openclaw folder"))
		fmt.Println(mutedStyle.Render("  (cozy roommates) or keep its own space (independent vibes)."))
		fmt.Println()

		choice, err := runLocationPicker(openclawDir)
		if err != nil {
			return cfg, fmt.Errorf("setup cancelled: %w", err)
		}

		cfg.ProjectsDir = choice
	} else {
		// No openclaw — straightforward setup.
		fmt.Println(mutedStyle.Render("  Setting up your projects home at ~/.projects/"))
		fmt.Println(mutedStyle.Render("  Config, projects, everything in one tidy spot."))
	}

	// Ensure the directories exist.
	if err := EnsureDirs(); err != nil {
		return cfg, err
	}

	// If they picked the openclaw path, ensure that dir exists too.
	if cfg.ProjectsDir != "" {
		expanded, err := expandPath(cfg.ProjectsDir)
		if err == nil {
			_ = os.MkdirAll(expanded, 0700)
		}
	}

	// Save the config.
	if err := Save(cfg); err != nil {
		return cfg, fmt.Errorf("save config: %w", err)
	}

	fmt.Println()
	fmt.Println(successStyle.Render("  ✨ All set!"))
	projDisplay := cfg.ProjectsDir
	if projDisplay == "" {
		projDisplay, _ = ProjectsDir()
	}
	fmt.Println(mutedStyle.Render("  Projects will live at: ") + primaryStyle.Render(shortenHome(projDisplay)))
	fmt.Println(mutedStyle.Render("  Config saved to:       ") + primaryStyle.Render("~/.projects/config.toml"))
	fmt.Println()

	return cfg, nil
}

// shortenHome replaces the home directory prefix with ~.
func shortenHome(path string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	if len(path) > len(home) && path[:len(home)] == home {
		return "~" + path[len(home):]
	}
	return path
}

// --- Location picker TUI ---

type locationPickerModel struct {
	choices  []string
	labels   []string
	cursor   int
	selected string
	done     bool
}

func newLocationPicker(openclawDir string) locationPickerModel {
	defaultDir, _ := ProjectsDir()
	openclawProjects := filepath.Join(openclawDir, "projects")
	return locationPickerModel{
		choices: []string{defaultDir, openclawProjects},
		labels: []string{
			"~/.projects/projects/  (independent — my own space)",
			shortenHome(openclawProjects) + "  (nested inside openclaw — cozy roommates)",
		},
		cursor: 0,
	}
}

func (m locationPickerModel) Init() tea.Cmd { return nil }

func (m locationPickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.done {
		return m, nil
	}

	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch keyMsg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.choices)-1 {
			m.cursor++
		}
	case "enter":
		m.selected = m.choices[m.cursor]
		m.done = true
		return m, tea.Quit
	case "esc", "ctrl+c":
		m.done = true
		return m, tea.Quit
	}

	return m, nil
}

func (m locationPickerModel) View() string {
	if m.done {
		return ""
	}

	header := primaryStyle.Render("  Where should projects live?") + "\n\n"

	var lines string
	for i, label := range m.labels {
		cursor := "  "
		style := mutedStyle
		if i == m.cursor {
			cursor = primaryStyle.Render("▸ ")
			style = lipgloss.NewStyle().Bold(true)
		}
		lines += "  " + cursor + style.Render(label) + "\n"
	}

	footer := "\n" + mutedStyle.Render("  ↑/↓ to move, enter to select")

	return header + lines + footer
}

func runLocationPicker(openclawDir string) (string, error) {
	m := newLocationPicker(openclawDir)
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	result := finalModel.(locationPickerModel)
	if result.selected == "" {
		return "", fmt.Errorf("no selection made")
	}

	return result.selected, nil
}
