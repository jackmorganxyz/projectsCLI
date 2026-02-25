---
name: projects-cli
description: Manage projects with projectsCLI — scaffold, list, view, edit, push, and delete projects from the terminal. Use this skill when the user wants to create a new project, organize existing projects, check project health, push to GitHub, or work with projectsCLI commands and project metadata.
license: MIT
compatibility: Requires projectsCLI binary installed. Optional gh CLI for GitHub integration.
metadata:
  author: jackmorganxyz
  version: "1.0"
  repository: "https://github.com/jackmorganxyz/projectsCLI"
---

# projectsCLI Agent Skill

You are working with **projectsCLI**, a terminal-native project manager. It scaffolds projects with a consistent directory structure, tracks metadata in YAML frontmatter, provides a TUI dashboard for humans, and outputs clean JSON for agents and scripts.

## Key Concepts

- **Binary**: `projectsCLI`
- **Config**: `~/.projects/config.toml` (TOML)
- **Projects directory**: `~/.projects/projects/` (configurable)
- **Project source of truth**: `<project-dir>/PROJECT.md` (YAML frontmatter + Markdown body)
- **JSON mode**: Enabled automatically when stdout is piped, or explicitly with `--json`

## When to Use This Skill

Use projectsCLI when the user wants to:
- Create or scaffold a new project
- List, view, or search across their projects
- Check the health/status of projects (git state, remotes, uncommitted changes)
- Push a project to GitHub (init, commit, create repo, push — all in one command)
- Load project metadata into shell variables for scripting
- Delete a project
- Edit project metadata or notes

## Installation Check

Before using projectsCLI, verify it is installed:

```sh
which projectsCLI
```

If not installed, the user can install via:

```sh
# Homebrew (recommended)
brew install jackmorganxyz/tap/projectsCLI

# Shell script
curl -sSL https://raw.githubusercontent.com/jackmorganxyz/projectsCLI/main/install.sh | sh
```

## Commands Quick Reference

| Command | Alias | Purpose |
|---------|-------|---------|
| `create <slug>` | — | Scaffold a new project |
| `list` | `ls` | List all projects (TUI or JSON) |
| `view <slug>` | — | View project details |
| `edit <slug>` | — | Open PROJECT.md in editor |
| `load <slug>` | — | Export project data (JSON, shell vars) |
| `delete <slug>` | `rm` | Delete a project |
| `status` | — | Health check across all projects |
| `push <slug>` | — | Full git workflow: init, commit, GitHub repo, push |

## Global Flags

All commands support:
- `--json` — Force JSON output (auto-enabled when piped)
- `--config <path>` — Override config file path
- `--version` — Print version

## Agent Integration Guidelines

**Always use `--json` or pipe output** to get structured data instead of TUI output:

```sh
projectsCLI ls --json
projectsCLI view my-project --json
projectsCLI status --json
```

**For non-interactive deletion**, always pass `--force`:

```sh
projectsCLI delete my-project --force --json
```

**To load project data into the environment**:

```sh
eval $(projectsCLI load my-project --export)
# Now available: $PROJECT_SLUG, $PROJECT_TITLE, $PROJECT_STATUS, $PROJECT_DIR, etc.
```

## Creating a Project

```sh
projectsCLI create <slug> [flags] --json
```

**Slug rules**: lowercase alphanumeric + hyphens only, max 64 chars, regex `^[a-z0-9]+(?:-[a-z0-9]+)*$`

**Flags**:
- `--title <string>` — Display name (defaults to slug)
- `--description <string>` — Short description
- `--tags <string>` — Comma-separated tags (e.g., `"go,api,backend"`)
- `--status <string>` — `active` (default), `paused`, or `archived`

**Example**:
```sh
projectsCLI create my-api --title "My API" --tags "go,api" --description "REST API service" --json
```

**JSON response**:
```json
{"status": "created", "slug": "my-api", "dir": "/path/to/my-api", "created_at": "2025-02-25T00:00:00Z"}
```

**Side effects**: Creates directory tree (`docs/`, `memory/`, `context/`, `tasks/`, `code/`, `private/`), writes `PROJECT.md` and template files, optionally runs `git init`.

## Listing Projects

```sh
projectsCLI ls --json
```

**JSON response**: Array of project objects, each with `meta`, `body`, and `dir` fields.

```json
[{"meta": {"title": "My API", "slug": "my-api", "status": "active", "tags": ["go"], "description": "...", "created_at": "...", "updated_at": "...", "git_remote": "..."}, "body": "# My API\n...", "dir": "/path/to/my-api"}]
```

## Viewing a Project

```sh
projectsCLI view <slug> --json
```

Returns a single project object (same shape as one element from `list`).

## Checking Project Health

```sh
projectsCLI status --json
```

**JSON response**:
```json
[{"slug": "my-api", "title": "My API", "status": "active", "has_git": true, "has_remote": true, "uncommitted": false, "has_project_md": true}]
```

**Useful queries**:
```sh
# Projects with uncommitted changes
projectsCLI status --json | jq '.[] | select(.uncommitted == true) | .slug'

# Projects without a remote
projectsCLI status --json | jq '.[] | select(.has_remote == false and .has_git == true) | .slug'
```

## Pushing to GitHub

```sh
projectsCLI push <slug> -m "commit message" --json
```

**Flags**:
- `-m`, `--message <string>` — Commit message (default: `"Update project"`)
- `--private` — Create private GitHub repo (default: `true`)
- `--no-github` — Skip GitHub repo creation

**Workflow**: git init (if needed) -> git add -A -> git commit -> gh repo create (if no remote) -> git push

**Requires**: `gh` CLI installed and authenticated for GitHub repo creation.

**JSON response**:
```json
{"status": "pushed", "slug": "my-api", "remote": "https://github.com/user/my-api"}
```

## Loading Project Data

```sh
# JSON (default)
projectsCLI load <slug> --json

# Shell export statements
projectsCLI load <slug> --export

# Eval-able bash variables
projectsCLI load <slug> --bash
```

Exported variables: `PROJECT_SLUG`, `PROJECT_TITLE`, `PROJECT_STATUS`, `PROJECT_DIR`, `PROJECT_DESCRIPTION`, `PROJECT_TAGS`, `PROJECT_GIT_REMOTE`.

## Deleting a Project

```sh
projectsCLI delete <slug> --force --json
```

**Important**: `--force` is required in non-interactive (piped/scripted) mode. Without it, the command errors.

**JSON response**:
```json
{"status": "deleted", "slug": "my-project"}
```

**This permanently removes the entire project directory.**

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

## Configuration

Config at `~/.projects/config.toml`:

```toml
projects_dir = "~/.projects/projects"
editor = "vim"
github_org = ""
auto_git_init = true
```

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
projectsCLI create my-service --title "My Service" --tags "python,api" --json
projectsCLI push my-service -m "Initial scaffold" --json
```

### Audit all projects
```sh
projectsCLI status --json | jq '.'
```

### Find and fix dirty projects
```sh
for slug in $(projectsCLI status --json | jq -r '.[] | select(.uncommitted == true) | .slug'); do
  projectsCLI push "$slug" -m "Auto-commit changes" --json
done
```

### Get a project's directory path
```sh
projectsCLI view my-project --json | jq -r '.dir'
```

For the full command specification with all flags, arguments, validation rules, and JSON schemas, see [references/REFERENCE.md](references/REFERENCE.md).
