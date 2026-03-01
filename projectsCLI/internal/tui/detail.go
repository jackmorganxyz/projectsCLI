package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jackmorganxyz/projectsCLI/internal/project"
)

// DetailModel displays a scrollable project view.
type DetailModel struct {
	viewport viewport.Model
	project  *project.Project
	ready    bool
	quitting bool
}

// NewDetailModel creates a detail view for a project.
func NewDetailModel(proj *project.Project) DetailModel {
	return DetailModel{project: proj}
}

func (m DetailModel) Init() tea.Cmd { return nil }

func (m DetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		headerHeight := 3
		footerHeight := 2
		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-headerHeight-footerHeight)
			m.viewport.SetContent(m.renderContent())
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - headerHeight - footerHeight
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

func (m DetailModel) View() string {
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

	header := titleStyle.Render("📂 " + m.project.Meta.Title)
	footer := helpStyle.Render(fmt.Sprintf("  scroll %d%% | q quit", int(m.viewport.ScrollPercent()*100)))

	return header + "\n\n" + m.viewport.View() + "\n" + footer
}

func (m DetailModel) renderContent() string {
	var sb strings.Builder
	p := m.project

	sb.WriteString(FormatField("Slug", Slug(p.Meta.Slug)) + "\n")
	sb.WriteString(FormatField("Status", StatusEmoji(p.Meta.Status)+StatusColor(p.Meta.Status)) + "\n")
	sb.WriteString(FormatField("Created", p.Meta.CreatedAt) + "\n")
	sb.WriteString(FormatField("Updated", p.Meta.UpdatedAt) + "\n")

	if p.Meta.Description != "" {
		sb.WriteString(FormatField("Description", p.Meta.Description) + "\n")
	}
	if len(p.Meta.Tags) > 0 {
		sb.WriteString(FormatField("Tags", TagList(p.Meta.Tags)) + "\n")
	}
	if p.Meta.GitRemote != "" {
		sb.WriteString(FormatField("Remote", Path(p.Meta.GitRemote)) + "\n")
	}
	sb.WriteString(FormatField("Directory", Path(p.Dir)) + "\n")

	if p.Body != "" {
		sb.WriteString("\n")
		sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(ColorPrimary)).Bold(true).Render("📝 Content") + "\n")
		sb.WriteString(Divider(40) + "\n")
		sb.WriteString(p.Body)
	}

	return sb.String()
}
