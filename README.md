<h1 align="center">projects</h1>

<p align="center">
  <em>Less chaos, more shipping. Built for humans and agents. ✨</em>
</p>

<p align="center">
  <a href="https://github.com/jackmorganxyz/projectsCLI/releases"><img src="https://img.shields.io/github/v/release/jackmorganxyz/projectsCLI" alt="Release"></a>
  <a href="./LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue" alt="License: MIT"></a>
  <a href="https://go.dev/"><img src="https://img.shields.io/github/go-mod/go-version/jackmorganxyz/projectsCLI?filename=projectsCLI/go.mod" alt="Go Version"></a>
</p>

---

A terminal-native project manager with a gorgeous TUI, JSON output for automation, and just enough personality to make you smile. Scaffold projects, track metadata, push to GitHub — all from one tool.

## 🚀 Quick Start

```sh
brew install jackmorganxyz/tap/projects
projects create --title "My Project"
projects ls
```

Three commands. You now have a scaffolded project with docs, memory files, task tracking, and a git repo.

## 📖 Documentation

This project has two READMEs because it serves two audiences:

| Audience | Guide |
|----------|-------|
| 🧑‍💻 **Humans** | [README_4_HUMANS.md](./README_4_HUMANS.md) — Tutorial walkthrough, command reference, tips, personality |
| 🤖 **Agents** | [README_4_AGENTS.md](./README_4_AGENTS.md) — Schemas, flag tables, JSON output specs, integration patterns |

## ✨ Features

- **Project scaffolding** — Opinionated directory structure with metadata, docs, tasks, and memory files
- **Interactive TUI dashboard** — Navigate your projects with a beautiful terminal UI ([Charmbracelet](https://charm.sh) stack)
- **Auto JSON mode** — Pipe any command and output switches from TUI to clean JSON automatically. Machines have feelings too
- **One-command GitHub push** — `push` handles git init, commit, repo creation, and push in a single step. Yes, really
- **Multi-account folders** — Organize projects by GitHub account. Push from work or personal — the CLI switches `gh auth` automatically
- **Portfolio health checks** — `status` shows git state, remotes, and uncommitted changes across all projects
- **Smart slug generation** — Just provide a title and the slug is auto-generated for you
- **Personality included** — Random quips, celebrations, and tips because dev tools should spark joy, not existential dread

## 🎯 Commands at a Glance

| Command | What it does |
|---------|-------------|
| `create [slug]` | Scaffold a new project — slug auto-generated from `--title` if omitted |
| `list` / `ls` | Dashboard of all projects (gorgeous TUI or clean JSON) |
| `view <slug>` | Project details in a scrollable, styled view |
| `edit <slug>` | Browse project files and open in your preferred editor |
| `open <slug>` | Open the project folder in Finder / Explorer / file manager |
| `load <slug>` | Export project data for scripts (`--json`, `--export`, `--bash`) |
| `delete <slug>` / `rm` | Delete a project (with appropriately dramatic confirmation prompts) |
| `status` | Health check across all projects — your morning standup, minus the standing |
| `update <slug>` | Update project metadata (title, description, status, tags) |
| `push <slug>` | Full git workflow: init → commit → create repo → push. One command to rule them all |
| `folder add/list/remove` | Manage folders for multi-account GitHub setups |
| `move <slug>` | Move a project between folders |

## 📦 Install

**Homebrew** (recommended):
```sh
brew install jackmorganxyz/tap/projects
```

**Quick install**:
```sh
curl -sSL https://raw.githubusercontent.com/jackmorganxyz/projectsCLI/main/install.sh | sh
```

**From source**:
```sh
git clone https://github.com/jackmorganxyz/projectsCLI.git
cd projectsCLI && make build && make install
```

## ⚙️ Configuration

```toml
# ~/.projects/config.toml
projects_dir = "~/.projects/projects"
editor = "cursor"                 # saved automatically on first `edit`
github_username = "my-username"
auto_git_init = true

# Multi-account folders (optional)
[[folders]]
name = "work"
github_account = "work-username"

[[folders]]
name = "personal"
github_account = "personal-username"
```

All fields are optional. Sensible defaults are built in — we're not here to make you configure things. `github_username` and `auto_git_init` are prompted during first-run setup. Folders are added via `projects folder add`.

## 🤖 Agent Skill

projects ships with an [Agent Skill](https://agentskills.io) — a portable instruction set that teaches AI agents how to use the CLI. Compatible with Claude Code, and any agent that supports the [Agent Skills format](https://github.com/agentskills/agentskills).

**Install the skill (Claude Code):**

```sh
claude install-skill https://github.com/jackmorganxyz/projectsCLI/tree/main/skill/projects-cli
```

**Or copy it manually** into your agent's skills directory:

```sh
# Clone just the skill
git clone --depth 1 --filter=blob:none --sparse https://github.com/jackmorganxyz/projectsCLI.git
cd projectsCLI && git sparse-checkout set skill/projects-cli
cp -r skill/projects-cli ~/.claude/skills/projects-cli
```

The skill lives in [`skill/projects-cli/`](./skill/projects-cli/) and follows the [Agent Skills specification](https://agentskills.io/specification).

## 🤝 Contributing

Contributions welcome! See [CONTRIBUTING.md](./CONTRIBUTING.md) for development setup and guidelines.

## 📄 License

[MIT](./LICENSE) — go nuts.
