# Changelog

All notable changes to projects will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/), and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Changed
- **Binary renamed from `projectsCLI` to `projects`** — all commands are now `projects <command>`
- **`create` slug is now optional** — provide `--title` and the slug is auto-generated (e.g. `--title "My Cool Project"` → `my-cool-project`)
- **`edit` now opens PROJECT.md in the OS default application** (TextEdit on macOS, Notepad on Windows, xdg-open on Linux) instead of `$EDITOR`/vim
- **`github_org` config field renamed to `github_username`** in `~/.projects/config.toml`
- First-run setup now interactively prompts for **GitHub username** and **auto-git-init** preference

### Added
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
