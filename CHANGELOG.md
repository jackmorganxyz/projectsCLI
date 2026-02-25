# Changelog

All notable changes to projectsCLI will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/), and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Added
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
- Homebrew tap (`jackpmorgan/tap/projectsCLI`)
- Install script (`install.sh`)
- Automated release pipeline with GoReleaser and GitHub Actions
