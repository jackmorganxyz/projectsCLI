---
name: projects-cli
description: Manage projects with the projects CLI — scaffold, list, view, edit, open, push, delete, and organize projects by GitHub account with folders. Supports AI agent integration (Claude Code, Codex CLI) for scaffolding and editing. Use this skill when the user wants to create a new project, organize existing projects, manage multi-account GitHub setups, check project health, push to GitHub, or work with projects commands and project metadata.
license: MIT
compatibility: Requires projects binary installed. Optional gh CLI for GitHub integration.
metadata:
  author: jackmorganxyz
  version: "1.1"
  repository: "https://github.com/jackmorganxyz/projectsCLI"
---

# projects Agent Skill

You are working with **projects**, a terminal-native project manager. It scaffolds projects with a consistent directory structure, tracks metadata in YAML frontmatter, provides a TUI dashboard for humans, and outputs clean JSON for agents and scripts.

## Key Concepts

- **Binary**: `projects`
- **Config**: `~/.projects/config.toml` (TOML)
- **Projects directory**: `~/.projects/projects/` (configurable)
- **Project source of truth**: `<project-dir>/PROJECT.md` (YAML frontmatter + Markdown body)
- **JSON mode**: Enabled automatically when stdout is piped, or explicitly with `--json`

## When to Use This Skill

Use projects when the user wants to:
- Create or scaffold a new project
- List, view, or search across their projects
- Check the health/status of projects (git state, remotes, uncommitted changes)
- Push a project to GitHub (init, commit, create repo, push — all in one command)
- Organize projects by GitHub account using folders
- Move projects between folders
- Load project metadata into shell variables for scripting
- Delete a project
- Update project metadata (title, status, tags, description) programmatically
- Edit project metadata or notes
- Open a project folder in the OS file manager

## Installation Check

Before using projects, verify it is installed:

```sh
which projects
```

If not installed, the user can install via:

```sh
# Homebrew (recommended)
brew install jackmorganxyz/tap/projects

# Shell script
curl -sSL https://raw.githubusercontent.com/jackmorganxyz/projectsCLI/main/install.sh | sh
```

## Commands Quick Reference

| Command | Alias | Purpose |
|---------|-------|---------|
| `create [slug]` | — | Scaffold a new project (slug auto-generated from `--title` if omitted). Optionally spawn an AI agent to fill out the scaffold. |
| `list` | `ls` | List all projects (TUI or JSON). In TUI mode, selecting a project shows a command picker. |
| `view <slug>` | — | View project details |
| `edit <slug>` | — | Browse project files, then choose manual edit (in an editor) or agent edit (AI-assisted via prompt). `--editor` to override, `--editor-picker` to re-pick editor. |
| `open <slug>` | — | Open project folder in OS file manager (Finder, Explorer, etc.) |
| `load <slug>` | — | Export project data (JSON, shell vars) |
| `delete <slug>` | `rm` | Delete a project |
| `status` | — | Health check across all projects |
| `update <slug>` | — | Update project metadata (title, description, status, tags) |
| `push <slug>` | — | Full git workflow: init, commit, GitHub repo, push |
| `folder add <name>` | — | Create a folder tied to a GitHub account |
| `folder list` | `folder ls` | List configured folders |
| `folder remove <name>` | `folder rm` | Remove a folder from config |
| `move <slug>` | — | Move a project to a different folder |

## Global Flags

All commands support:
- `--json` — Force JSON output (auto-enabled when piped)
- `--config <path>` — Override config file path
- `--folder <name>` — Target a specific folder (for multi-account setups)
- `--version` — Print version

## Agent Integration Guidelines

**Always use `--json` or pipe output** to get structured data instead of TUI output:

```sh
projects ls --json
projects view my-project --json
projects status --json
```

**For non-interactive deletion**, always pass `--force`:

```sh
projects delete my-project --force --json
```

**To load project data into the environment**:

```sh
eval $(projects load my-project --export)
# Now available: $PROJECT_SLUG, $PROJECT_TITLE, $PROJECT_STATUS, $PROJECT_DIR, etc.
```

## Creating a Project

```sh
projects create [slug] [flags] --json
```

**Slug rules**: lowercase alphanumeric + hyphens only, max 64 chars, regex `^[a-z0-9]+(?:-[a-z0-9]+)*$`

If slug is omitted, it is auto-generated from `--title` (e.g. `--title "My Cool Project"` produces slug `my-cool-project`).

**Flags**:
- `--title <string>` — Display name (defaults to slug if slug is provided)
- `--description <string>` — Short description
- `--tags <string>` — Comma-separated tags (e.g., `"go,api,backend"`)
- `--status <string>` — `active` (default), `paused`, or `archived`

**Examples**:
```sh
# Explicit slug
projects create my-api --title "My API" --tags "go,api" --description "REST API service" --json

# Auto-generate slug from title
projects create --title "My API" --tags "go,api" --json
```

**JSON response**:
```json
{"status": "created", "slug": "my-api", "dir": "/path/to/my-api", "created_at": "2025-02-25T00:00:00Z"}
```

**Side effects**: Creates directory tree (`docs/`, `memory/`, `context/`, `tasks/`, `code/`, `private/`), writes `PROJECT.md` and template files, optionally runs `git init`.

**AI agent integration**: In interactive mode, if Claude Code (`claude`) or Codex CLI (`codex`) is installed, the user is prompted to optionally spawn an AI agent to fill out the scaffolded files. The user provides a text prompt describing what the agent should do, and the agent runs in the project directory.

## Listing Projects

```sh
projects ls --json
```

**JSON response**: Array of project objects, each with `meta`, `body`, and `dir` fields. Note: `tags`, `description`, `git_remote`, and `body` are omitted from JSON when empty.

```json
[{"meta": {"title": "My API", "slug": "my-api", "status": "active", "tags": ["go"], "description": "...", "created_at": "...", "updated_at": "...", "git_remote": "..."}, "body": "# My API\n...", "dir": "/path/to/my-api"}]
```

**Interactive behavior**: In TUI mode, the dashboard displays a project table. Selecting a project (Enter) shows a command picker dropdown — view, edit, open, status, push, update, move, or delete — and auto-runs the chosen command on that project.

## Viewing a Project

```sh
projects view <slug> --json
```

Returns a single project object (same shape as one element from `list`).

## Checking Project Health

```sh
projects status --json
```

**JSON response**:
```json
[{"slug": "my-api", "title": "My API", "status": "active", "has_git": true, "has_remote": true, "uncommitted": false, "has_project_md": true}]
```

**Useful queries**:
```sh
# Projects with uncommitted changes
projects status --json | jq '.[] | select(.uncommitted == true) | .slug'

# Projects without a remote
projects status --json | jq '.[] | select(.has_remote == false and .has_git == true) | .slug'
```

## Pushing to GitHub

```sh
projects push <slug> -m "commit message" --json
```

**Flags**:
- `-m`, `--message <string>` — Commit message (default: `"Update project"`)
- `--private` — Create private GitHub repo (default: `true`)
- `--no-github` — Skip GitHub repo creation

**Workflow**: git init (if needed) -> git add -A -> git commit -> if project is in a folder: `gh auth switch` to the folder's account -> if no remote: `gh repo create --source --push` / if remote exists: `git push -u origin <branch>`

**Requires**: `gh` CLI installed and authenticated for GitHub repo creation. For folder-based projects, the associated account must be authenticated in `gh auth`.

**JSON response**:
```json
{"status": "pushed", "slug": "my-api", "remote": "https://github.com/user/my-api"}
```

## Updating Project Metadata

```sh
projects update <slug> [flags] --json
```

**Flags**:
- `--title <string>` — New title
- `--description <string>` — New description
- `--status <string>` — `active`, `paused`, or `archived`
- `--tags <string>` — Comma-separated tags (replaces existing)

At least one flag is required. The `updated_at` timestamp is set automatically.

**Examples**:
```sh
projects update my-api --status archived --json
projects update my-api --title "My API v2" --tags "go,api,v2" --json
```

**JSON response**:
```json
{"status": "updated", "slug": "my-api", "updated_at": "2025-02-25T00:00:00Z"}
```

**Side effects**: Updates `PROJECT.md` frontmatter and regenerates `PROJECTS.md` registry.

## Extracting Specific Fields

The `list`, `view`, and `status` commands support `--field` for extracting a specific field without needing `jq`:

```sh
# Get all project directories
projects ls --field dir

# Get a single project's title
projects view my-api --field meta.title

# Get all slugs from status
projects status --field slug
```

Supports dot-notation for nested fields (e.g. `meta.title`, `meta.slug`, `meta.status`).

## Loading Project Data

```sh
# JSON (default)
projects load <slug> --json

# Shell export statements
projects load <slug> --export

# Eval-able bash variables
projects load <slug> --bash
```

Exported variables: `PROJECT_SLUG`, `PROJECT_TITLE`, `PROJECT_STATUS`, `PROJECT_DIR`, `PROJECT_DESCRIPTION`. Additionally, `PROJECT_TAGS` and `PROJECT_GIT_REMOTE` are exported when non-empty.

## Deleting a Project

```sh
projects delete <slug> --force --json
```

**Important**: `--force` is required in non-interactive (piped/scripted) mode. Without it, the command errors.

**JSON response**:
```json
{"status": "deleted", "slug": "my-project"}
```

**This permanently removes the entire project directory.**

## Editing a Project File

```sh
projects edit <slug>                    # interactive file browser + edit mode picker
projects edit <slug> --editor vim       # skip edit mode and editor picker, use vim
projects edit <slug> --editor-picker    # force re-showing the editor picker
```

Interactively browse project files and choose an edit mode:

1. **Manual edit** — open the file in your preferred editor. On first run, auto-detects installed editors (Cursor, VS Code, TextEdit, Vim, Nano, etc.) and prompts the user to pick one. The choice is saved to `config.editor`. Editors are labeled as `(terminal)` or `(GUI)` in the picker.
2. **Agent edit** — spawn an AI agent (Claude Code or Codex CLI) to edit the file. The user provides a text prompt describing the changes, and the agent runs in the project directory. This option only appears when at least one AI agent is installed.

In non-interactive mode, defaults to opening `PROJECT.md` with the saved editor.

**Flags**:
- `--editor <command>` — override saved editor for a single invocation (skips edit mode picker)
- `--editor-picker` — force re-showing the editor selection even if a preference is saved

**Note:** This command is interactive — for agents, use direct file reads/writes on project files instead.

## Opening a Project Folder

```sh
projects open <slug>
```

Opens the project directory in the OS file manager (Finder on macOS, Explorer on Windows, xdg-open on Linux).

## Project Directory Structure

Every scaffolded project contains:

```
<slug>/
  PROJECT.md          # Source of truth (YAML frontmatter + markdown)
  USAGE.md            # Workspace guide for humans and agents
  docs/README.md      # Project documentation
  memory/MEMORY.md    # Persistent notes / AI agent context
  context/CONTEXT.md  # Architecture decisions
  tasks/TODO.md       # Task tracking
  code/               # Source code
  private/            # Gitignored (secrets, drafts)
  .gitignore
```

## PROJECT.md Schema

```yaml
---
title: "My Project"
slug: "my-project"
status: "active"          # active | paused | archived
tags: ["go", "cli"]
description: "Short description"
created_at: "2025-02-25T00:00:00Z"
updated_at: "2025-02-25T00:00:00Z"
git_remote: "https://github.com/user/my-project"
---

# My Project

Freeform markdown body.
```

## Multi-Account Folders

Folders let users organize projects by GitHub account. When pushing, the CLI automatically switches `gh auth` to the folder's account.

### Set up folders
```sh
# Interactive — picks from gh auth accounts
projects folder add work --json

# Explicit account
projects folder add personal --account my-gh-user --json
```

**JSON response**:
```json
{"status": "created", "folder": "work", "github_account": "work-org", "path": "/path/to/projects/work"}
```

### Create projects in a folder
```sh
projects create --title "Work API" --folder work --json
```

### Move existing projects
```sh
projects move my-project --folder work --json
```

**JSON response**:
```json
{"status": "moved", "slug": "my-project", "from": "/path/from", "to": "/path/to", "to_folder": "work"}
```

### List folders
```sh
projects folder list --json
```

### Remove a folder from config
```sh
projects folder remove work --json
```

Does not delete the directory or its projects.

## Configuration

Config at `~/.projects/config.toml`:

```toml
projects_dir = "~/.projects/projects"
editor = "cursor"    # auto-saved on first `edit` pick
github_username = ""
auto_git_init = true

# Multi-account folders (managed via `projects folder add`)
[[folders]]
name = "work"
github_account = "work-org"

[[folders]]
name = "personal"
github_account = "my-gh-user"
```

`github_username` and `auto_git_init` are prompted during first-run setup. Folders are managed via `projects folder add/remove`.

## Reading Project Files Directly

You can also read project files directly from the filesystem:

```sh
# Project metadata and notes
cat ~/.projects/projects/<slug>/PROJECT.md

# Persistent memory for agent context
cat ~/.projects/projects/<slug>/memory/MEMORY.md

# Architecture decisions
cat ~/.projects/projects/<slug>/context/CONTEXT.md

# Task list
cat ~/.projects/projects/<slug>/tasks/TODO.md
```

## Common Agent Workflows

### Create and push a new project
```sh
projects create --title "My Service" --tags "python,api" --json
projects push my-service -m "Initial scaffold" --json
```

### Audit all projects
```sh
projects status --json | jq '.'
```

### Find and fix dirty projects
```sh
for slug in $(projects status --json | jq -r '.[] | select(.uncommitted == true) | .slug'); do
  projects push "$slug" -m "Auto-commit changes" --json
done
```

### Get a project's directory path
```sh
projects view my-project --json | jq -r '.dir'
```

### Set up and use multi-account folders
```sh
projects folder add work --account work-org --json
projects create --title "Work API" --folder work --json
projects push work-api -m "Initial scaffold" --json
# → Automatically pushes under the work-org GitHub account
```

### Move a project to a folder
```sh
projects move my-project --folder work --json
```

For the full command specification with all flags, arguments, validation rules, and JSON schemas, see [references/REFERENCE.md](references/REFERENCE.md).
