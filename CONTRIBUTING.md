# Contributing to projectsCLI

Thanks for your interest in contributing! Whether it's a bug fix, new feature, or documentation improvement, we appreciate the help.

## Development Setup

```bash
# Clone the repo
git clone https://github.com/jackmorganxyz/projectsCLI.git
cd projectsCLI

# Build the binary
make build

# Run it
make run ARGS="--help"

# Or run with specific commands
make run ARGS="ls"
make run ARGS="create test-project --title 'Test Project'"
```

**Requirements:**
- Go 1.23+
- `make`
- `gh` CLI (optional, for GitHub integration testing)

## Project Structure

```
projectsCLI/
├── cmd/projectsCLI/
│   └── main.go              # Entry point, Cobra root command
├── internal/
│   ├── cli/                  # Command implementations
│   │   ├── context.go        # Runtime context (config passing)
│   │   ├── create.go         # create command
│   │   ├── delete.go         # delete/rm command
│   │   ├── edit.go           # edit command
│   │   ├── helpers.go        # JSON output, slug validation
│   │   ├── list.go           # list/ls command
│   │   ├── load.go           # load command (agent data export)
│   │   ├── push.go           # push command (git workflow)
│   │   ├── status.go         # status command (health check)
│   │   └── view.go           # view command
│   ├── config/               # TOML configuration
│   ├── git/                  # Git and GitHub CLI wrappers
│   ├── project/              # Project CRUD, scaffolding, registry
│   └── tui/                  # Charmbracelet TUI components
│       ├── components.go     # Spinner, confirmation prompt
│       ├── dashboard.go      # Project list table
│       ├── detail.go         # Scrollable project view
│       ├── detect.go         # TTY/JSON mode detection
│       ├── flair.go          # Random messages, quips, emoji
│       ├── gitpanel.go       # Git status visualization
│       ├── selector.go       # Project selector
│       ├── theme.go          # Color palette and styles
│       └── wizard.go         # Interactive setup
```

## Making Changes

1. **Fork** the repo and create a feature branch
2. Make your changes
3. Run checks: `make check` (formats code, runs vet and tests)
4. Submit a PR

## Code Conventions

- Follow standard Go conventions (`gofmt`, `go vet`)
- Every command must handle **both** TUI and JSON output modes
- User-facing messages go through `tui/` helpers (for emoji/color support)
- Fun strings (quips, cheers, tips) live in `tui/flair.go`

## Adding a New Command

1. Create `internal/cli/<command>.go`
2. Follow the pattern in existing commands (e.g., `create.go`):
   - `NewXxxCmd()` returns a `*cobra.Command`
   - Handle JSON mode via `tui.IsJSON()`
   - Handle interactive mode via `tui.IsInteractive()`
3. Register it in `cmd/projectsCLI/main.go` via `rootCmd.AddCommand()`
4. Add flair strings to `tui/flair.go` for any celebratory messages

## Reporting Bugs

Use the [bug report template](https://github.com/jackmorganxyz/projectsCLI/issues/new?template=bug_report.yml) on GitHub.

## Feature Requests

Use the [feature request template](https://github.com/jackmorganxyz/projectsCLI/issues/new?template=feature_request.yml) on GitHub.
