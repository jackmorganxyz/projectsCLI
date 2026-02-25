package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SelectorItem represents a selectable project in the picker.
type SelectorItem struct {
	slug        string
	title       string
	description string
}

func (i SelectorItem) Title() string       { return i.slug }
func (i SelectorItem) Description() string { return i.title + " - " + i.description }
func (i SelectorItem) FilterValue() string { return i.slug + " " + i.title }

// SelectorModel is an fzf-style project picker.
type SelectorModel struct {
	list     list.Model
	selected string
	quitting bool
}

// NewSelectorModel creates a project selector from slug/title pairs.
func NewSelectorModel(items []SelectorItem) SelectorModel {
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = item
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color(ColorPrimary)).
		BorderForeground(lipgloss.Color(ColorPrimary))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(lipgloss.Color(ColorMuted)).
		BorderForeground(lipgloss.Color(ColorPrimary))

	l := list.New(listItems, delegate, 60, 20)
	l.Title = "Select Project"
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorPrimary)).
		Bold(true)
	l.SetFilteringEnabled(true)

	return SelectorModel{list: l}
}

func (m SelectorModel) Init() tea.Cmd { return nil }

func (m SelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			if item, ok := m.list.SelectedItem().(SelectorItem); ok {
				m.selected = item.slug
			}
			m.quitting = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m SelectorModel) View() string {
	if m.quitting {
		return ""
	}
	return m.list.View()
}

// Selected returns the slug of the selected project (empty if cancelled).
func (m SelectorModel) Selected() string {
	return m.selected
}

// RunSelector runs the project selector and returns the selected slug.
func RunSelector(slugs []string, titles []string, descriptions []string) (string, error) {
	items := make([]SelectorItem, len(slugs))
	for i, slug := range slugs {
		title := slug
		if i < len(titles) {
			title = titles[i]
		}
		desc := ""
		if i < len(descriptions) {
			desc = descriptions[i]
		}
		items[i] = SelectorItem{slug: slug, title: title, description: desc}
	}

	// Strip empty descriptions.
	for i := range items {
		items[i].description = strings.TrimSpace(items[i].description)
	}

	m := NewSelectorModel(items)
	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}
	return finalModel.(SelectorModel).Selected(), nil
}
