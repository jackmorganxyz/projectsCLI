<h1 align="center">projectsCLI</h1>

<p align="center">
  <em>Your projects, organized. Built for humans and agents.</em>
</p>

<p align="center">
  <a href="https://github.com/jackpmorgan/projects-cli/releases"><img src="https://img.shields.io/github/v/release/jackpmorgan/projects-cli" alt="Release"></a>
  <a href="./LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue" alt="License: MIT"></a>
  <a href="https://go.dev/"><img src="https://img.shields.io/github/go-mod/go-version/jackpmorgan/projects-cli?filename=projectsCLI/go.mod" alt="Go Version"></a>
</p>

---

A terminal-native project manager with a gorgeous TUI, JSON output for automation, and just enough personality to make you smile. Scaffold projects, track metadata, push to GitHub — all from one tool.

## Quick Start

```sh
brew install jackpmorgan/tap/projectsCLI
projectsCLI create my-project --title "My Project"
projectsCLI ls
```

## Documentation

This project has two READMEs because it serves two audiences:

| Audience | Guide |
|----------|-------|
| **Developers & humans** | [README_4_HUMANS.md](./README_4_HUMANS.md) — Tutorial walkthrough, command reference, tips, personality |
| **AI agents & scripts** | [README_4_AGENTS.md](./README_4_AGENTS.md) — Schemas, flag tables, JSON output specs, integration patterns |

## Features

- **Project scaffolding** — Opinionated directory structure with metadata, docs, tasks, and memory files
- **Interactive TUI dashboard** — Navigate your projects with a beautiful terminal UI ([Charmbracelet](https://charm.sh) stack)
- **Auto JSON mode** — Pipe any command and output switches from TUI to clean JSON automatically
- **One-command GitHub push** — `push` handles git init, commit, repo creation, and push in a single step
- **Portfolio health checks** — `status` shows git state, remotes, and uncommitted changes across all projects
- **Personality included** — Random quips, celebrations, and tips because dev tools should spark joy

## Commands at a Glance

| Command | What it does |
|---------|-------------|
| `create <slug>` | Scaffold a new project with directory structure and metadata |
| `list` / `ls` | Dashboard of all projects (TUI or JSON) |
| `view <slug>` | Project details (scrollable TUI or JSON) |
| `edit <slug>` | Open PROJECT.md in your editor |
| `load <slug>` | Export project data for scripts (`--json`, `--export`, `--bash`) |
| `delete <slug>` / `rm` | Delete a project (with dramatic confirmation prompts) |
| `status` | Health check across all projects |
| `push <slug>` | Full git workflow: init, commit, create repo, push |

## Install

**Homebrew:**
```sh
brew install jackpmorgan/tap/projectsCLI
```

**Quick install:**
```sh
curl -sSL https://raw.githubusercontent.com/jackpmorgan/projects-cli/main/install.sh | sh
```

**From source:**
```sh
git clone https://github.com/jackpmorgan/projects-cli.git
cd projects-cli && make build && make install
```

## Configuration

```toml
# ~/.openclaw/config.toml
projects_dir = "~/.openclaw/projects"
editor = "nvim"
github_org = "my-org"
auto_git_init = true
```

All fields are optional. Sensible defaults are built in.

## Contributing

Contributions welcome! See [CONTRIBUTING.md](./CONTRIBUTING.md) for development setup and guidelines.

## License

[MIT](./LICENSE)
