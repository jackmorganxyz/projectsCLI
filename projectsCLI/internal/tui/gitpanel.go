package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// GitStatusMsg carries git status output to the TUI.
type GitStatusMsg struct {
	Status string
	Branch string
	Remote string
	Err    error
}

// GitPanelModel displays git status in a scrollable panel.
type GitPanelModel struct {
	viewport viewport.Model
	status   string
	branch   string
	remote   string
	err      error
	ready    bool
	quitting bool
}

// NewGitPanelModel creates a git status display panel.
func NewGitPanelModel() GitPanelModel {
	return GitPanelModel{}
}

func (m GitPanelModel) Init() tea.Cmd { return nil }

func (m GitPanelModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		headerHeight := 4
		footerHeight := 2
		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-headerHeight-footerHeight)
			m.viewport.SetContent(m.renderContent())
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - headerHeight - footerHeight
			m.viewport.SetContent(m.renderContent())
		}

	case GitStatusMsg:
		m.status = msg.Status
		m.branch = msg.Branch
		m.remote = msg.Remote
		m.err = msg.Err
		if m.ready {
			m.viewport.SetContent(m.renderContent())
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m GitPanelModel) View() string {
	if m.quitting {
		return ""
	}
	if !m.ready {
		return "Loading..."
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorPrimary)).
		Bold(true)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorMuted))

	header := titleStyle.Render("Git Status")
	if m.branch != "" {
		header += " " + lipgloss.NewStyle().Foreground(lipgloss.Color(ColorInfo)).Render("("+m.branch+")")
	}

	footer := helpStyle.Render("  q quit")

	return header + "\n\n" + m.viewport.View() + "\n" + footer
}

func (m GitPanelModel) renderContent() string {
	if m.err != nil {
		return ErrorMessage(m.err.Error())
	}

	var sb strings.Builder

	if m.branch != "" {
		sb.WriteString(FormatField("Branch", m.branch) + "\n")
	}
	if m.remote != "" {
		sb.WriteString(FormatField("Remote", m.remote) + "\n")
	}

	sb.WriteString("\n")

	if strings.TrimSpace(m.status) == "" {
		sb.WriteString(SuccessMessage("Working tree clean") + "\n")
	} else {
		sb.WriteString(lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorWarning)).
			Bold(true).
			Render("Changes:") + "\n\n")

		for _, line := range strings.Split(m.status, "\n") {
			if strings.TrimSpace(line) == "" {
				continue
			}
			styled := styleGitLine(line)
			sb.WriteString(fmt.Sprintf("  %s\n", styled))
		}
	}

	return sb.String()
}

// styleGitLine applies color to a git status line based on its prefix.
func styleGitLine(line string) string {
	if len(line) < 2 {
		return line
	}

	prefix := line[:2]
	rest := line[2:]

	switch {
	case strings.Contains(prefix, "M"):
		return lipgloss.NewStyle().Foreground(lipgloss.Color(ColorWarning)).Render(prefix) + rest
	case strings.Contains(prefix, "A"):
		return lipgloss.NewStyle().Foreground(lipgloss.Color(ColorSuccess)).Render(prefix) + rest
	case strings.Contains(prefix, "D"):
		return lipgloss.NewStyle().Foreground(lipgloss.Color(ColorError)).Render(prefix) + rest
	case strings.Contains(prefix, "?"):
		return lipgloss.NewStyle().Foreground(lipgloss.Color(ColorMuted)).Render(prefix) + rest
	default:
		return line
	}
}
