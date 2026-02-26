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
| `--folder` | string | `""` | Target a specific folder (for multi-account setups) |
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
  "created_at": "2025-02-25T00:00:00Z",
  "folder": "work"
}
```

`folder` is included when `--folder` is used.

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

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--field` | string | `""` | Extract a specific field (e.g. `--field dir`, `--field meta.title`) |

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
    "dir": "/Users/you/.projects/projects/my-project",
    "folder": "work"
  }
]
```

`folder` is included when the project lives in a configured folder.

**Behavior:**
- Interactive TTY: launches TUI dashboard
- Non-TTY or `--json`: outputs JSON array
- When folders are configured, `list` and `status` include a Folder column in table output

---

### `view <slug>`

Display project details.

**Arguments:**

| Arg | Required | Type |
|-----|----------|------|
| `slug` | yes | string |

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--field` | string | `""` | Extract a specific field (e.g. `--field dir`, `--field meta.title`) |

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

Interactively browse project files and open the selected file in a chosen editor.

**Arguments:**

| Arg | Required | Type |
|-----|----------|------|
| `slug` | yes | string |

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--editor` | string | `""` | Editor command to use (bypasses interactive picker) |

**Behavior:**
- **Interactive**: Shows a file browser to navigate the project directory, then opens the selected file in the user's preferred editor. On first run, prompts the user to pick from detected installed editors (Cursor, VS Code, Vim, etc.) and saves the choice to `config.editor`.
- **Non-interactive**: Opens `PROJECT.md` with the saved `config.editor` (defaults to `$EDITOR` or `vim`).
- **`--editor` flag**: Overrides the saved editor for a single invocation without re-saving.

**Editor detection:** Auto-detects installed GUI editors (macOS via Spotlight, Linux/Windows via PATH) and terminal editors (nvim, vim, nano, emacs, micro, hx).

**Side effects:** Saves the editor choice to `config.editor` on first interactive pick. Terminal editors run in the foreground (blocking); GUI editors launch in the background.

**Note:** For non-interactive agent use, prefer direct file reads/writes instead.

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

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--field` | string | `""` | Extract a specific field (e.g. `--field slug`, `--field status`) |

**JSON output:**

```json
[
  {
    "slug": "my-project",
    "folder": "work",
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
- `folder`: which folder the project belongs to (omitted when empty)
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
4. If project is in a folder with a GitHub account: runs `gh auth switch --user <account>`
5. If no remote and `--no-github` is false: creates GitHub repo via `gh` CLI under the folder's account
6. If remote exists: pushes to remote with `--set-upstream`

**Requirements:** `gh` CLI must be installed and authenticated for GitHub repo creation.

---

### `update <slug>`

Update project metadata.

**Arguments:**

| Arg | Required | Type |
|-----|----------|------|
| `slug` | yes | string |

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--title` | string | `""` | New title |
| `--description` | string | `""` | New description |
| `--status` | string | `""` | New status (`active`, `paused`, `archived`) |
| `--tags` | string | `""` | New tags (comma-separated, replaces existing) |

At least one flag is required.

**JSON output:**

```json
{
  "status": "updated",
  "slug": "my-project",
  "updated_at": "2025-02-25T00:00:00Z"
}
```

**Side effects:**
- Updates `PROJECT.md` frontmatter with new values
- Sets `updated_at` to current time
- Regenerates `PROJECTS.md` registry

---

### `folder`

Manage named folders for multi-account GitHub setups. Subcommands: `add`, `list` (`ls`), `remove` (`rm`).

#### `folder add <name>`

**Arguments:**

| Arg | Required | Type | Validation |
|-----|----------|------|------------|
| `name` | yes | string | Same slug rules: `^[a-z0-9]+(?:-[a-z0-9]+)*$` |

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account` | string | `""` | GitHub username for this folder (interactive picker if omitted) |

**JSON output:**

```json
{
  "status": "created",
  "folder": "work",
  "github_account": "work-org",
  "path": "/Users/you/.projects/projects/work"
}
```

**Behavior:**
- If `--account` is omitted and `gh` is authenticated: shows interactive account picker (auto-selects if only one account)
- If `--account` is omitted and `gh` is not available: returns error
- Warns (does not block) if the account isn't found in `gh auth` — the folder is still created
- Creates the folder subdirectory under the projects directory

#### `folder list` (alias: `ls`)

**JSON output:**

```json
[
  {"name": "work", "github_account": "work-org"},
  {"name": "personal", "github_account": "my-gh-user"}
]
```

#### `folder remove <name>` (alias: `rm`)

**JSON output:**

```json
{
  "status": "removed",
  "folder": "work"
}
```

**Note:** Only removes the folder from config. Does not delete the directory or its projects.

---

### `move <slug>`

Move a project between folders or to the top level.

**Arguments:**

| Arg | Required | Type |
|-----|----------|------|
| `slug` | yes | string |

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--folder` | string | (required) | Target folder name, or `""` for top level |

**JSON output:**

```json
{
  "status": "moved",
  "slug": "my-project",
  "from": "/Users/you/.projects/projects/my-project",
  "to": "/Users/you/.projects/projects/work/my-project",
  "to_folder": "work"
}
```

`from_folder` and `to_folder` are included when applicable.

**Side effects:**
- Moves the project directory via `os.Rename`
- Regenerates `PROJECTS.md` registry

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
projects_dir = "~/.projects/projects"
editor = "vim"
github_username = "my-username"
auto_git_init = true

[[folders]]
name = "work"
github_account = "work-org"

[[folders]]
name = "personal"
github_account = "my-gh-user"
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `projects_dir` | string | `~/.projects/projects` | Where projects are stored |
| `editor` | string | `$EDITOR` or `"vim"` | Editor command for `edit` (auto-saved on first interactive pick) |
| `github_username` | string | `""` | Default GitHub username for `push` repo creation |
| `auto_git_init` | bool | `true` | Auto git init on project create |
| `folders` | array | `[]` | Named folders with associated GitHub accounts |
| `folders[].name` | string | — | Folder name (used as subdirectory name) |
| `folders[].github_account` | string | — | GitHub account for this folder (used by `push`) |

`github_username` and `auto_git_init` are prompted interactively during first-run setup. Folders are managed via `projects folder add/remove`.

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

### Update project metadata
```sh
projects update my-project --status archived --json
projects update my-project --title "New Name" --tags "go,v2" --json
```

### Extract specific fields without jq
```sh
projects ls --field dir                  # one directory path per line
projects view my-project --field meta.title   # just the title
projects status --field slug             # one slug per line
```

### Read project files directly
```sh
# The PROJECT.md is the source of truth
cat ~/.projects/projects/my-project/PROJECT.md

# Memory file for persistent context
cat ~/.projects/projects/my-project/memory/MEMORY.md
```

### Create a project in a folder
```sh
projects create --title "Work API" --folder work --json
```

### List projects in a specific folder
```sh
projects ls --folder work --json
```

### Set up multi-account folders
```sh
projects folder add work --account work-org --json
projects folder add personal --account my-gh-user --json
```

### Move a project to a folder
```sh
projects move my-project --folder work --json
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
