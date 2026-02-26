package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jackmorganxyz/projectsCLI/internal/project"
)

// DashboardModel displays the project list with colored status.
type DashboardModel struct {
	table    table.Model
	projects []*project.Project
	quitting bool
	selected bool
}

// NewDashboardModel creates a dashboard from a list of projects.
func NewDashboardModel(projects []*project.Project) DashboardModel {
	columns := []table.Column{
		{Title: "Slug", Width: 24},
		{Title: "Title", Width: 30},
		{Title: "Status", Width: 12},
		{Title: "Created", Width: 12},
		{Title: "Tags", Width: 20},
	}

	var rows []table.Row
	for _, p := range projects {
		created := p.Meta.CreatedAt
		if len(created) > 10 {
			created = created[:10]
		}
		tags := ""
		if len(p.Meta.Tags) > 0 {
			for i, t := range p.Meta.Tags {
				if i > 0 {
					tags += ", "
				}
				tags += t
			}
		}
		rows = append(rows, table.Row{
			p.Meta.Slug,
			p.Meta.Title,
			p.Meta.Status,
			created,
			tags,
		})
	}

	tbl := NewStyledTable(columns, rows)

	return DashboardModel{
		table:    tbl,
		projects: projects,
	}
}

func (m DashboardModel) Init() tea.Cmd { return nil }

func (m DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			m.quitting = true
			m.selected = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m DashboardModel) View() string {
	if m.quitting {
		return ""
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorPrimary)).
		Bold(true).
		MarginBottom(1)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorMuted)).
		MarginTop(1)

	view := titleStyle.Render("📋 Projects") + "\n\n"
	view += m.table.View() + "\n"
	view += helpStyle.Render(fmt.Sprintf("  %d projects | j/k navigate | enter select | q quit", len(m.projects)))

	return view
}

// WasSelected reports whether the user pressed Enter (vs q/esc to quit).
func (m DashboardModel) WasSelected() bool {
	return m.selected
}

// SelectedSlug returns the slug of the currently selected project.
func (m DashboardModel) SelectedSlug() string {
	row := m.table.SelectedRow()
	if len(row) > 0 {
		return row[0]
	}
	return ""
}
