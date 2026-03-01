package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

const (
	ColorPrimary = "#8B5CF6" // Violet (bright)
	ColorSuccess = "#34D399" // Emerald (bright)
	ColorError   = "#F87171" // Red (soft)
	ColorWarning = "#FBBF24" // Amber (bright)
	ColorMuted   = "#9CA3AF" // Gray (lighter)
	ColorInfo    = "#60A5FA" // Blue (bright)
)

// Theme contains shared styles for projectsCLI terminal rendering.
type Theme struct {
	headerStyle  lipgloss.Style
	successStyle lipgloss.Style
	errorStyle   lipgloss.Style
	warningStyle lipgloss.Style
	mutedStyle   lipgloss.Style
	cardStyle    lipgloss.Style
	cardTitle    lipgloss.Style
	keyLabel     lipgloss.Style
	keyValue     lipgloss.Style
	tableHeader  lipgloss.Style
	tableCell    lipgloss.Style
	tableBorder  lipgloss.Style
	bubbleTable  table.Styles
}

// DefaultTheme is the shared application style set.
var DefaultTheme = NewTheme()

// NewTheme builds the default projectsCLI theme.
func NewTheme() Theme {
	primary := lipgloss.Color(ColorPrimary)
	success := lipgloss.Color(ColorSuccess)
	errColor := lipgloss.Color(ColorError)
	warning := lipgloss.Color(ColorWarning)
	muted := lipgloss.Color(ColorMuted)

	styles := table.DefaultStyles()
	styles.Header = styles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(primary).
		Foreground(primary).
		Bold(true)
	styles.Cell = styles.Cell.Padding(0, 1)
	styles.Selected = styles.Selected.
		Foreground(lipgloss.Color("#ffffff")).
		Background(primary).
		Bold(true)

	return Theme{
		headerStyle:  lipgloss.NewStyle().Foreground(primary).Bold(true),
		successStyle: lipgloss.NewStyle().Foreground(success).Bold(true),
		errorStyle:   lipgloss.NewStyle().Foreground(errColor).Bold(true),
		warningStyle: lipgloss.NewStyle().Foreground(warning).Bold(true),
		mutedStyle:   lipgloss.NewStyle().Foreground(muted),
		cardStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primary).
			Padding(0, 1),
		cardTitle: lipgloss.NewStyle().Foreground(primary).Bold(true),
		keyLabel:  lipgloss.NewStyle().Foreground(muted),
		keyValue:  lipgloss.NewStyle().Bold(true),
		tableHeader: lipgloss.NewStyle().Foreground(primary).Bold(true),
		tableCell:   lipgloss.NewStyle(),
		tableBorder: lipgloss.NewStyle().Foreground(muted),
		bubbleTable: styles,
	}
}

// Header renders a section header.
func Header(text string) string { return DefaultTheme.headerStyle.Render(text) }

// SuccessMessage renders a success message with emoji prefix.
func SuccessMessage(text string) string {
	return DefaultTheme.successStyle.Render(SuccessEmoji() + text)
}

// ErrorMessage renders an error message with emoji prefix.
func ErrorMessage(text string) string {
	return DefaultTheme.errorStyle.Render(ErrorEmoji() + text)
}

// WarningMessage renders a warning message with emoji prefix.
func WarningMessage(text string) string {
	return DefaultTheme.warningStyle.Render(WarningEmoji() + text)
}

// InfoMessage renders an informational message with emoji prefix.
func InfoMessage(text string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(ColorInfo)).Bold(true).Render(InfoEmoji() + text)
}

// Muted renders muted/dimmed text.
func Muted(text string) string { return DefaultTheme.mutedStyle.Render(text) }

// InfoCard renders a titled informational card.
func InfoCard(title, body string) string {
	lines := []string{DefaultTheme.cardTitle.Render(title)}
	if strings.TrimSpace(body) != "" {
		lines = append(lines, body)
	}
	return DefaultTheme.cardStyle.Render(strings.Join(lines, "\n"))
}

// KeyValue renders a label: value pair.
func KeyValue(label, value string) string {
	return DefaultTheme.keyLabel.Render(label+":") + " " + DefaultTheme.keyValue.Render(value)
}

// BubbleTableStyles returns styles for bubbles/table model.
func BubbleTableStyles() table.Styles { return DefaultTheme.bubbleTable }

// Table renders a compact plain-text table.
func Table(headers []string, rows [][]string) string {
	if len(headers) == 0 {
		return ""
	}

	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i := range headers {
			if i < len(row) && len(row[i]) > widths[i] {
				widths[i] = len(row[i])
			}
		}
	}

	renderRow := func(cells []string, style lipgloss.Style) string {
		parts := make([]string, len(headers))
		for i := range headers {
			cell := ""
			if i < len(cells) {
				cell = cells[i]
			}
			parts[i] = style.Render(padRight(cell, widths[i]))
		}
		return strings.Join(parts, DefaultTheme.tableBorder.Render(" | "))
	}

	borderParts := make([]string, len(headers))
	for i := range headers {
		borderParts[i] = strings.Repeat("-", widths[i])
	}
	separator := DefaultTheme.tableBorder.Render(strings.Join(borderParts, "-+-"))

	lines := []string{
		renderRow(headers, DefaultTheme.tableHeader),
		separator,
	}
	for _, row := range rows {
		lines = append(lines, renderRow(row, DefaultTheme.tableCell))
	}

	return strings.Join(lines, "\n")
}

// StatusColor returns the styled status string.
func StatusColor(status string) string {
	switch strings.ToLower(status) {
	case "active":
		return DefaultTheme.successStyle.Render(status)
	case "archived":
		return DefaultTheme.mutedStyle.Render(status)
	case "paused":
		return DefaultTheme.warningStyle.Render(status)
	case "error", "broken":
		return DefaultTheme.errorStyle.Render(status)
	default:
		return status
	}
}

// NewStyledTable builds a themed bubbles/table component.
func NewStyledTable(columns []table.Column, rows []table.Row) table.Model {
	height := len(rows) + 2
	if height < 4 {
		height = 4
	}
	if height > 20 {
		height = 20
	}

	tbl := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithHeight(height),
		table.WithFocused(true),
	)
	tbl.SetStyles(BubbleTableStyles())
	return tbl
}

func padRight(value string, width int) string {
	if len(value) >= width {
		return value
	}
	return value + strings.Repeat(" ", width-len(value))
}

// FormatField renders a labeled field for detail views.
func FormatField(label, value string) string {
	return fmt.Sprintf("  %s %s", DefaultTheme.keyLabel.Render(label+":"), DefaultTheme.keyValue.Render(value))
}

// Slug renders a project slug in the primary accent color.
func Slug(text string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(ColorPrimary)).Bold(true).Render(text)
}

// Path renders a file/directory path in info blue.
func Path(text string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(ColorInfo)).Render(text)
}

// TagList renders tags as a colorful comma-separated string.
func TagList(tags []string) string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorSuccess))
	parts := make([]string, len(tags))
	for i, t := range tags {
		parts[i] = style.Render(t)
	}
	return strings.Join(parts, DefaultTheme.mutedStyle.Render(", "))
}

// Divider renders a styled horizontal divider.
func Divider(width int) string {
	if width <= 0 {
		width = 40
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color(ColorPrimary)).Render(strings.Repeat("─", width))
}
