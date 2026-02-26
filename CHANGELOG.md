# Changelog

All notable changes to projects will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/), and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Changed
- **Binary renamed from `projectsCLI` to `projects`** — all commands are now `projects <command>`
- **`create` slug is now optional** — provide `--title` and the slug is auto-generated (e.g. `--title "My Cool Project"` → `my-cool-project`)
- **`edit` now has an interactive file browser + editor picker** — browse any file in the project, choose from installed editors (Cursor, VS Code, Vim, etc.), choice is saved to config. Use `--editor` flag to override. Falls back to PROJECT.md + saved editor in non-interactive mode.
- **`github_org` config field renamed to `github_username`** in `~/.projects/config.toml`
- First-run setup now interactively prompts for **GitHub username** and **auto-git-init** preference

### Added
- **Interactive editor detection** — `edit` auto-detects installed GUI editors (macOS via Spotlight, Linux/Windows via PATH) and terminal editors (nvim, vim, nano, emacs, micro, hx)
- **`--editor` flag** on `edit` command — bypass the editor picker for a single invocation
- **`editor` package** (`internal/editor/`) — reusable cross-platform editor detection and launch logic
- **Multi-account folders** — organize projects by GitHub account with `folder add`, `folder list`, `folder remove`
- **`move <slug>` command** — move projects between folders or back to top level
- **`--folder` global flag** — scope any command to a specific folder
- **Automatic `gh auth switch`** on `push` — projects in folders automatically switch to the folder's GitHub account
- **Interactive GitHub account picker** — `folder add` without `--account` shows a picker from `gh auth` accounts
- **Folder column** in `list` and `status` output when folders are configured
- **`folder` field** in JSON output for projects, status health, create, and move commands
- **`update <slug>` command** — update project metadata (title, description, status, tags) without opening an editor
- **`--field` flag** on `list`, `view`, and `status` — extract specific fields with dot-notation (e.g. `--field meta.title`) without needing `jq`
- `open <slug>` command — opens the project folder in Finder (macOS), Explorer (Windows), or default file manager (Linux)
- `Slugify()` helper — converts titles to valid slugs with unicode normalization
- Tests for `Slugify` and `ValidateSlug`
- README.md, README_4_HUMANS.md, README_4_AGENTS.md documentation
- LICENSE (MIT)
- CONTRIBUTING.md with development setup and code conventions
- CODE_OF_CONDUCT.md
- SECURITY.md
- CHANGELOG.md
- Makefile with build, test, lint, and release targets
- GitHub issue templates (bug report, feature request)
- GitHub pull request template

## [0.1.x] — 2025

### Added
- `create` command — scaffold projects with opinionated directory structure
- `list` / `ls` command — interactive TUI dashboard and JSON output
- `view` command — scrollable project detail view
- `edit` command — open PROJECT.md in configured editor
- `load` command — export project data as JSON, shell exports, or bash variables
- `delete` / `rm` command — delete with confirmation prompts
- `status` command — health check across all projects
- `push` command — full git workflow with GitHub repo creation
- Interactive TUI with Charmbracelet (Bubble Tea, Lip Gloss)
- Auto JSON output mode (detected via TTY)
- TOML configuration (`~/.projects/config.toml`)
- Project scaffolding with docs, memory, context, tasks directories
- Auto-generated PROJECTS.md registry
- Random quips, celebrations, tips, and personality
- Homebrew tap (`jackmorganxyz/tap/projects`)
- Install script (`install.sh`)
- Automated release pipeline with GoReleaser and GitHub Actions
