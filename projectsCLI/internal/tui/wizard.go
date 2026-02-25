package tui

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// WizardResult holds the collected values from the create wizard.
type WizardResult struct {
	Slug        string
	Title       string
	Description string
	Tags        string
	Status      string
}

// RunCreateWizard runs the interactive create form and returns the result.
func RunCreateWizard() (*WizardResult, error) {
	result := &WizardResult{
		Status: "active",
	}

	theme := huh.ThemeBase()
	theme.Focused.Title = theme.Focused.Title.Foreground(lipgloss.Color(ColorPrimary))
	theme.Focused.SelectedOption = theme.Focused.SelectedOption.Foreground(lipgloss.Color(ColorPrimary))

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Project Slug").
				Description("Lowercase identifier (e.g. my-project)").
				Placeholder("my-project").
				Value(&result.Slug),

			huh.NewInput().
				Title("Title").
				Description("Human-readable project name").
				Placeholder("My Project").
				Value(&result.Title),

			huh.NewText().
				Title("Description").
				Description("Brief project description").
				Value(&result.Description),
		),

		huh.NewGroup(
			huh.NewInput().
				Title("Tags").
				Description("Comma-separated tags (e.g. go,cli,tools)").
				Placeholder("go, cli").
				Value(&result.Tags),

			huh.NewSelect[string]().
				Title("Status").
				Options(
					huh.NewOption("Active", "active"),
					huh.NewOption("Paused", "paused"),
					huh.NewOption("Archived", "archived"),
				).
				Value(&result.Status),
		),
	).WithTheme(theme)

	if err := form.Run(); err != nil {
		return nil, err
	}

	return result, nil
}
