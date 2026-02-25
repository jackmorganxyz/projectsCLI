package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// RunProgram runs a bubbletea program and returns the final model.
func RunProgram(m tea.Model) (tea.Model, error) {
	p := tea.NewProgram(m, tea.WithAltScreen())
	return p.Run()
}

// RunConfirm runs a confirmation prompt and returns the result.
func RunConfirm(question string) (bool, error) {
	m := NewConfirmationPrompt(question)
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return false, err
	}
	return finalModel.(ConfirmationPromptModel).Confirmed(), nil
}

// SpinnerModel is a shared Bubble Tea model for async operations.
type SpinnerModel struct {
	spinner spinner.Model
	message string
	done    bool
	err     error
}

// DoneMsg signals that an async operation completed.
type DoneMsg struct {
	Err error
}

// NewSpinnerModel creates a themed spinner for long-running work.
func NewSpinnerModel(message string) SpinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorPrimary))

	if strings.TrimSpace(message) == "" {
		message = "Working..."
	}

	return SpinnerModel{spinner: s, message: message}
}

func (m SpinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m SpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch typed := msg.(type) {
	case DoneMsg:
		m.done = true
		m.err = typed.Err
		return m, tea.Quit
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(typed)
		return m, cmd
	case tea.KeyMsg:
		if typed.String() == "ctrl+c" {
			m.done = true
			m.err = fmt.Errorf("cancelled")
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m SpinnerModel) View() string {
	if m.done {
		if m.err != nil {
			return ErrorMessage(m.err.Error())
		}
		return SuccessMessage("Done")
	}
	return fmt.Sprintf("%s %s", m.spinner.View(), m.message)
}

// ConfirmationPromptModel is a yes/no Bubble Tea prompt.
type ConfirmationPromptModel struct {
	question    string
	yesSelected bool
	done        bool
}

// NewConfirmationPrompt returns a prompt model with default answer = yes.
func NewConfirmationPrompt(question string) ConfirmationPromptModel {
	if strings.TrimSpace(question) == "" {
		question = "Continue?"
	}
	return ConfirmationPromptModel{question: question, yesSelected: true}
}

func (m ConfirmationPromptModel) Init() tea.Cmd { return nil }

func (m ConfirmationPromptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.done {
		return m, nil
	}

	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch keyMsg.String() {
	case "left", "h", "y", "Y":
		m.yesSelected = true
	case "right", "l", "n", "N":
		m.yesSelected = false
	case "enter":
		m.done = true
		return m, tea.Quit
	case "esc", "ctrl+c":
		m.done = true
		m.yesSelected = false
		return m, tea.Quit
	}

	return m, nil
}

func (m ConfirmationPromptModel) View() string {
	yes := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorMuted)).Render("Yes")
	no := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorMuted)).Render("No")

	if m.yesSelected {
		yes = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorSuccess)).Bold(true).Render("Yes")
	} else {
		no = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorError)).Bold(true).Render("No")
	}

	return Header(m.question) + "\n" + fmt.Sprintf("[%s] / [%s]", yes, no)
}

// Confirmed reports the selected answer once the prompt exits.
func (m ConfirmationPromptModel) Confirmed() bool {
	return m.done && m.yesSelected
}
