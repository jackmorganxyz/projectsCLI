# projects — Agent Reference

Machine-readable reference for AI agents, scripts, and automation.

## Overview

| Property | Value |
|----------|-------|
| Binary | `projects` |
| Config file | `~/.projects/config.toml` (TOML) |
| Data directory | `~/.projects/projects/` |
| Project file | `<project-dir>/PROJECT.md` (YAML frontmatter + Markdown body) |
| Output | JSON when `--json` flag is set or stdout is not a TTY |

## Installation

```sh
# Homebrew
brew install jackmorganxyz/tap/projects

# Shell script
curl -sSL https://raw.githubusercontent.com/jackmorganxyz/projectsCLI/main/install.sh | sh
```

## Global Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--json` | bool | `false` (auto `true` when piped) | Force JSON output |
| `--config` | string | `~/.projects/config.toml` | Config file path |
| `--version` | bool | `false` | Print version and exit |

---

## Commands

### `create [slug]`

Create a new project scaffold.

**Arguments:**

| Arg | Required | Type | Validation |
|-----|----------|------|------------|
| `slug` | no | string | Regex `^[a-z0-9]+(?:-[a-z0-9]+)*$`, max 64 chars |

If slug is omitted, it is auto-generated from `--title` (e.g. `--title "My Cool Project"` produces slug `my-cool-project`). Either a slug argument or `--title` must be provided.

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--title` | string | slug value | Project title (required if slug is omitted) |
| `--description` | string | `""` | Project description |
| `--tags` | []string | `[]` | Comma-separated tags |
| `--status` | string | `"active"` | Initial status (`active`, `paused`, `archived`) |

**JSON output:**

```json
{
  "status": "created",
  "slug": "my-project",
  "dir": "/Users/you/.projects/projects/my-project",
  "created_at": "2025-02-25T00:00:00Z"
}
```

**Side effects:**
- Creates directory tree: `docs/`, `memory/`, `context/`, `tasks/`, `code/`, `private/`
- Writes `PROJECT.md` with YAML frontmatter
- Writes template files: `USAGE.md`, `memory/MEMORY.md`, `context/CONTEXT.md`, `tasks/TODO.md`, `docs/README.md`, `.gitignore`
- If `auto_git_init = true`: runs `git init`, `git add -A`, `git commit`
- Regenerates `PROJECTS.md` registry

---

### `list`

List all projects. Alias: `ls`.

**Arguments:** None.

**Flags:** None (uses global `--json`).

**JSON output:**

```json
[
  {
    "meta": {
      "title": "My Project",
      "slug": "my-project",
      "status": "active",
      "tags": ["go", "cli"],
      "description": "A brief description",
      "created_at": "2025-02-25T00:00:00Z",
      "updated_at": "2025-02-25T00:00:00Z",
      "git_remote": "https://github.com/user/my-project"
    },
    "body": "# My Project\n\nMarkdown content...\n",
    "dir": "/Users/you/.projects/projects/my-project"
  }
]
```

**Behavior:**
- Interactive TTY: launches TUI dashboard
- Non-TTY or `--json`: outputs JSON array

---

### `view <slug>`

Display project details.

**Arguments:**

| Arg | Required | Type |
|-----|----------|------|
| `slug` | yes | string |

**Flags:** None (uses global `--json`).

**JSON output:** Same as a single element from `list` output (Project object with `meta`, `body`, `dir`).

**Behavior:**
- Interactive TTY: scrollable TUI detail view
- Non-TTY: plain text fields
- `--json`: full project JSON

---

### `load <slug>`

Output project data for agent consumption.

**Arguments:**

| Arg | Required | Type |
|-----|----------|------|
| `slug` | yes | string |

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--export` | bool | `false` | Output as shell `export` statements |
| `--bash` | bool | `false` | Output as eval-able bash variables |

**Output formats:**

`--json` (default):
```json
{
  "meta": {
    "title": "My Project",
    "slug": "my-project",
    "status": "active",
    "tags": ["go"],
    "description": "Description",
    "created_at": "2025-02-25T00:00:00Z",
    "updated_at": "2025-02-25T00:00:00Z",
    "git_remote": "https://github.com/user/my-project"
  },
  "body": "# My Project\n...",
  "dir": "/Users/you/.projects/projects/my-project"
}
```

`--export`:
```sh
export PROJECT_SLUG="my-project"
export PROJECT_TITLE="My Project"
export PROJECT_STATUS="active"
export PROJECT_DIR="/Users/you/.projects/projects/my-project"
export PROJECT_DESCRIPTION="Description"
export PROJECT_TAGS="go"
export PROJECT_GIT_REMOTE="https://github.com/user/my-project"
```

`--bash`:
```sh
PROJECT_SLUG="my-project"
PROJECT_TITLE="My Project"
PROJECT_STATUS="active"
PROJECT_DIR="/Users/you/.projects/projects/my-project"
PROJECT_DESCRIPTION="Description"
PROJECT_TAGS="go"
PROJECT_GIT_REMOTE="https://github.com/user/my-project"
```

---

### `edit <slug>`

Open project's `PROJECT.md` in the OS default application.

**Arguments:**

| Arg | Required | Type |
|-----|----------|------|
| `slug` | yes | string |

**Flags:** None.

**Side effects:** Opens `PROJECT.md` with the OS default application (TextEdit on macOS, Notepad on Windows, xdg-open on Linux). The CLI returns immediately.

**Note:** This command is interactive-only. Not suitable for non-interactive agent use — use direct file reads/writes instead.

---

### `open <slug>`

Open the project directory in the OS file manager.

**Arguments:**

| Arg | Required | Type |
|-----|----------|------|
| `slug` | yes | string |

**Flags:** None.

**Side effects:** Opens the project directory in Finder (macOS), Explorer (Windows), or xdg-open (Linux). The CLI returns immediately.

**Note:** This command is interactive-only. For agents, use `projects view <slug> --json | jq -r '.dir'` to get the directory path.

---

### `delete <slug>`

Delete a project and its directory. Alias: `rm`.

**Arguments:**

| Arg | Required | Type |
|-----|----------|------|
| `slug` | yes | string |

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--force` | bool | `false` | Skip confirmation prompt |

**JSON output:**

```json
{
  "status": "deleted",
  "slug": "my-project"
}
```

**Important:** In non-interactive mode (piped/scripted), `--force` is **required**. Without it, the command returns an error.

**Side effects:**
- Removes the entire project directory (`rm -rf`)
- Regenerates `PROJECTS.md` registry

---

### `status`

Health check across all projects.

**Arguments:** None.

**Flags:** None (uses global `--json`).

**JSON output:**

```json
[
  {
    "slug": "my-project",
    "title": "My Project",
    "status": "active",
    "has_git": true,
    "has_remote": true,
    "uncommitted": false,
    "has_project_md": true
  }
]
```

**Fields:**
- `has_git`: whether the project directory is a git repository
- `has_remote`: whether a git remote is configured
- `uncommitted`: whether there are uncommitted changes
- `has_project_md`: whether `PROJECT.md` exists

---

### `push <slug>`

Full git workflow: init, stage, commit, create GitHub repo, push.

**Arguments:**

| Arg | Required | Type |
|-----|----------|------|
| `slug` | yes | string |

**Flags:**

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--message` | `-m` | string | `"Update project"` | Commit message |
| `--private` | | bool | `true` | Create private GitHub repo |
| `--no-github` | | bool | `false` | Skip GitHub repo creation |

**JSON output:**

```json
{
  "status": "pushed",
  "slug": "my-project",
  "remote": "https://github.com/user/my-project"
}
```

**Workflow:**
1. If not a git repo: runs `git init`
2. Runs `git add -A`
3. If uncommitted changes: commits with provided message
4. If no remote and `--no-github` is false: creates GitHub repo via `gh` CLI
5. If remote exists: pushes to remote with `--set-upstream`

**Requirements:** `gh` CLI must be installed and authenticated for GitHub repo creation.

---

## Data Schemas

### PROJECT.md Format

```yaml
---
title: "My Project"
slug: "my-project"
status: "active"
tags:
  - go
  - cli
description: "A brief description"
created_at: "2025-02-25T00:00:00Z"
updated_at: "2025-02-25T00:00:00Z"
git_remote: "https://github.com/user/my-project"
---

# My Project

Markdown body content here.
```

### ProjectMeta Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `title` | string | yes | Display name |
| `slug` | string | yes | Unique identifier (`^[a-z0-9]+(?:-[a-z0-9]+)*$`) |
| `status` | string | yes | One of: `active`, `paused`, `archived` |
| `tags` | []string | no | Categorization tags |
| `description` | string | no | Short description |
| `created_at` | string (RFC 3339) | yes | Creation timestamp |
| `updated_at` | string (RFC 3339) | yes | Last update timestamp |
| `git_remote` | string (URL) | no | Git remote URL, set by `push` |

### Config Schema (`~/.projects/config.toml`)

```toml
projects_dir = "~/.projects/projects"   # Path to projects directory
editor = "vim"                          # Editor binary name
github_username = "my-username"         # GitHub username for repo creation
auto_git_init = true                    # Auto-init git on `create` (default: true)
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `projects_dir` | string | `~/.projects/projects` | Where projects are stored |
| `editor` | string | `$EDITOR` or `"vim"` | Editor binary name |
| `github_username` | string | `""` | GitHub username for `push` repo creation |
| `auto_git_init` | bool | `true` | Auto git init on project create |

Both `github_username` and `auto_git_init` are prompted interactively during first-run setup.

### Directory Structure per Project

```
<slug>/
  PROJECT.md            # YAML frontmatter + markdown body (source of truth)
  USAGE.md              # Workspace guide — directory roles, conventions, agent instructions
  docs/README.md        # Project documentation
  memory/MEMORY.md      # Persistent notes and context
  context/CONTEXT.md    # Architecture decisions
  tasks/TODO.md         # Task tracking
  code/                 # Code directory
  private/              # Gitignored, never pushed
  .gitignore            # Ignores private/
```

### Registry File

`~/.projects/projects/PROJECTS.md` — Auto-generated markdown table of all projects. Regenerated on `create` and `delete`.

---

## Agent Integration Patterns

### Load project into shell environment
```sh
eval $(projects load my-project --export)
echo $PROJECT_SLUG    # my-project
echo $PROJECT_DIR     # /Users/you/.projects/projects/my-project
echo $PROJECT_STATUS  # active
```

### List all projects as JSON
```sh
projects ls --json
# Auto-detects piped output, so this also works:
projects ls | jq '.'
```

### Get a single project
```sh
projects view my-project --json | jq '.meta'
```

### Find projects with uncommitted changes
```sh
projects status --json | jq '.[] | select(.uncommitted == true) | .slug'
```

### Find projects without a remote
```sh
projects status --json | jq '.[] | select(.has_remote == false and .has_git == true) | .slug'
```

### Non-interactive delete
```sh
projects delete my-project --force --json
```

### Create a project programmatically
```sh
projects create --title "My API" --tags "go,api" --json
```

### Read project files directly
```sh
# The PROJECT.md is the source of truth
cat ~/.projects/projects/my-project/PROJECT.md

# Memory file for persistent context
cat ~/.projects/projects/my-project/memory/MEMORY.md
```

---

## Exit Codes

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | Error (message printed to stderr) |

## Environment Variables

| Variable | Effect |
|----------|--------|
| `NO_EMOJI` | Set to any value to disable emoji in TUI output |
| `TERM=dumb` | Disables emoji in TUI output |
