package tui

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// emojiEnabled reports whether the terminal likely supports emoji.
// Respects NO_EMOJI=1 and TERM=dumb as opt-outs.
func emojiEnabled() bool {
	if os.Getenv("NO_EMOJI") != "" {
		return false
	}
	if os.Getenv("TERM") == "dumb" {
		return false
	}
	return true
}

// pick returns a random element from the slice.
func pick(items []string) string {
	return items[rand.IntN(len(items))]
}

// --- Emoji prefixes ---

// SuccessEmoji returns the emoji prefix for success messages.
func SuccessEmoji() string {
	if emojiEnabled() {
		return "✨ "
	}
	return "ok "
}

// ErrorEmoji returns the emoji prefix for error messages.
func ErrorEmoji() string {
	if emojiEnabled() {
		return "💥 "
	}
	return "err "
}

// WarningEmoji returns the emoji prefix for warning messages.
func WarningEmoji() string {
	if emojiEnabled() {
		return "⚠️  "
	}
	return "warn "
}

// InfoEmoji returns the emoji prefix for info messages.
func InfoEmoji() string {
	if emojiEnabled() {
		return "💡 "
	}
	return "-- "
}

// --- Message pools ---

var createCheers = []string{
	"Fresh project, who dis?",
	"Another masterpiece begins.",
	"Off to a great start!",
	"The world needs this.",
	"Built different.",
	"Let's gooo.",
}

// RandomCreateCheer returns a random creation celebration.
func RandomCreateCheer() string { return pick(createCheers) }

var deleteFarewells = []string{
	"Gone but not forgotten.",
	"And just like that... poof.",
	"Freed up some disk karma.",
	"Making room for the next big thing.",
}

// RandomDeleteFarewell returns a random deletion farewell.
func RandomDeleteFarewell() string { return pick(deleteFarewells) }

var deleteCancelled = []string{
	"Crisis averted.",
	"Living to see another day.",
	"Good call, that one's a keeper.",
	"Phew, close one.",
}

// RandomDeleteCancelled returns a random cancellation quip.
func RandomDeleteCancelled() string { return pick(deleteCancelled) }

var deleteConfirms = []string{
	"Really delete %q? No backsies.",
	"Nuke %q from orbit? This cannot be undone.",
	"Say goodbye to %q? This cannot be undone.",
	"Delete %q? It's had a good run. This cannot be undone.",
}

// RandomDeleteConfirm returns a random delete confirmation prompt.
func RandomDeleteConfirm(slug string) string {
	return fmt.Sprintf(pick(deleteConfirms), slug)
}

var pushCheers = []string{
	"Shipped it!",
	"To the cloud and beyond!",
	"Another one for the history books.",
	"Your code is free now.",
	"Chef's kiss.",
}

// RandomPushCheer returns a random push/commit celebration.
func RandomPushCheer() string { return pick(pushCheers) }

var celebrations = []string{
	"Nice work!",
	"Smooth operator.",
	"Like a pro.",
	"Clean.",
	"Nailed it.",
}

// RandomCelebration returns a random short celebration phrase.
func RandomCelebration() string { return pick(celebrations) }

var noChangesQuips = []string{
	"Already looking good — nothing to commit.",
	"Squeaky clean.",
	"All caught up. Maybe go touch grass?",
}

// RandomNoChanges returns a random no-changes quip.
func RandomNoChanges() string { return pick(noChangesQuips) }

var emptyStateMessages = []string{
	"It's quiet in here... too quiet. Run 'projects create <slug>' to fix that.",
	"No projects yet — the world is your oyster. Try 'projects create <slug>'.",
	"A blank canvas awaits. Zero projects, infinite potential. Run 'projects create <slug>'.",
}

// RandomEmptyState returns a random empty-state message.
func RandomEmptyState() string { return pick(emptyStateMessages) }

var spinnerMessages = []string{
	"Doing the thing...",
	"Crunching bytes...",
	"Almost there, probably...",
	"Summoning the code elves...",
	"Reticulating splines...",
	"Polishing pixels...",
	"Warming up the flux capacitor...",
}

// RandomSpinnerMessage returns a fun waiting message.
func RandomSpinnerMessage() string { return pick(spinnerMessages) }

// SpinnerMessages returns the full spinner message pool for rotation.
func SpinnerMessages() []string { return spinnerMessages }

var tips = []string{
	"Tip: Use 'projects ls' for a quick overview of all your projects.",
	"Tip: 'projects push <slug>' handles git init, commit, and GitHub in one step.",
	"Tip: Add tags with --tags to stay organized.",
	"Tip: 'projects status' gives you a health check across all projects.",
	"Tip: Set $EDITOR to customize which editor 'projects edit' opens.",
}

// MaybeTip returns a styled tip ~30% of the time, or empty string.
func MaybeTip() string {
	if rand.IntN(10) < 3 {
		return DefaultTheme.mutedStyle.Render(pick(tips))
	}
	return ""
}

// --- Banner ---

var taglines = []string{
	"Your projects, organized.",
	"Build cool stuff.",
	"Less chaos, more shipping.",
	"Projects under control.",
	"Tidy projects, tidy mind.",
}

// Banner returns a compact styled banner for the root command.
func Banner() string {
	tagline := pick(taglines)

	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorPrimary)).
		Bold(true).
		Render("projects")

	tag := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorMuted)).
		Italic(true).
		Render(tagline)

	content := fmt.Sprintf("  %s — %s", title, tag)

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(ColorPrimary)).
		Padding(0, 1)

	emoji := ""
	if emojiEnabled() {
		emoji = " 🚀"
	}

	return box.Render(content + emoji) + "\n"
}

// --- Status emoji ---

// StatusEmoji returns an emoji for a project status string.
func StatusEmoji(status string) string {
	if !emojiEnabled() {
		return ""
	}
	switch strings.ToLower(status) {
	case "active":
		return "🟢 "
	case "paused":
		return "🟡 "
	case "archived":
		return "📦 "
	case "error", "broken":
		return "🔴 "
	default:
		return ""
	}
}
