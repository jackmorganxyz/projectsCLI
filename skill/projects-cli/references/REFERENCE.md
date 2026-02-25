# projects — Full Command Reference

Complete specification of every command, flag, argument, and JSON output schema.

## Global Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--json` | bool | `false` (auto `true` when piped) | Force JSON output |
| `--config` | string | `~/.projects/config.toml` | Config file path |
| `--version` | bool | `false` | Print version and exit |

---

## `create [slug]`

Create a new project scaffold.

### Arguments

| Arg | Required | Type | Validation |
|-----|----------|------|------------|
| `slug` | no | string | Regex `^[a-z0-9]+(?:-[a-z0-9]+)*$`, max 64 chars |

If slug is omitted, it is auto-generated from `--title` (e.g. `--title "My Cool Project"` produces slug `my-cool-project`). Either a slug argument or `--title` must be provided.

### Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--title` | string | slug value | Project title (required if slug is omitted) |
| `--description` | string | `""` | Project description |
| `--tags` | []string | `[]` | Comma-separated tags |
| `--status` | string | `"active"` | Initial status (`active`, `paused`, `archived`) |

### JSON Output

```json
{
  "status": "created",
  "slug": "my-project",
  "dir": "/Users/you/.projects/projects/my-project",
  "created_at": "2025-02-25T00:00:00Z"
}
```

### Side Effects

- Creates directory tree: `docs/`, `memory/`, `context/`, `tasks/`, `code/`, `private/`
- Writes `PROJECT.md` with YAML frontmatter
- Writes template files: `USAGE.md`, `memory/MEMORY.md`, `context/CONTEXT.md`, `tasks/TODO.md`, `docs/README.md`, `.gitignore`
- If `auto_git_init = true` in config: runs `git init`, `git add -A`, `git commit`
- Regenerates `PROJECTS.md` registry in the projects directory

---

## `list` (alias: `ls`)

List all projects.

### Arguments

None.

### Flags

None (uses global `--json`).

### JSON Output

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

### Behavior

- Interactive TTY: launches TUI dashboard
- Non-TTY or `--json`: outputs JSON array
- Optional fields (`tags`, `description`, `git_remote`, `body`) are omitted from JSON when empty

---

## `view <slug>`

Display project details.

### Arguments

| Arg | Required | Type |
|-----|----------|------|
| `slug` | yes | string |

### Flags

None (uses global `--json`).

### JSON Output

Same shape as a single element from `list` output — object with `meta`, `body`, `dir`.

### Behavior

- Interactive TTY: scrollable TUI detail view
- Non-TTY: plain text fields
- `--json`: full project JSON

---

## `load <slug>`

Output project data for agent/script consumption.

### Arguments

| Arg | Required | Type |
|-----|----------|------|
| `slug` | yes | string |

### Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--export` | bool | `false` | Output as shell `export` statements |
| `--bash` | bool | `false` | Output as eval-able bash variables |

### Output Formats

**JSON (default)**:
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

**`--export`**:
```sh
export PROJECT_SLUG="my-project"
export PROJECT_TITLE="My Project"
export PROJECT_STATUS="active"
export PROJECT_DIR="/Users/you/.projects/projects/my-project"
export PROJECT_DESCRIPTION="Description"
export PROJECT_TAGS="go"                                        # only if tags are non-empty
export PROJECT_GIT_REMOTE="https://github.com/user/my-project"  # only if remote is set
```

**`--bash`**:
```sh
PROJECT_SLUG="my-project"
PROJECT_TITLE="My Project"
PROJECT_STATUS="active"
PROJECT_DIR="/Users/you/.projects/projects/my-project"
PROJECT_DESCRIPTION="Description"
PROJECT_TAGS="go"                                        # only if tags are non-empty
PROJECT_GIT_REMOTE="https://github.com/user/my-project"  # only if remote is set
```

**Note:** `PROJECT_TAGS` and `PROJECT_GIT_REMOTE` are only included when the project has tags or a git remote configured, respectively.

---

## `edit <slug>`

Open project's `PROJECT.md` in the OS default application (e.g. TextEdit on macOS).

### Arguments

| Arg | Required | Type |
|-----|----------|------|
| `slug` | yes | string |

### Flags

None.

### Side Effects

Opens `PROJECT.md` with the OS default application for `.md` files:
- **macOS**: `open` (typically TextEdit)
- **Windows**: `start` (typically Notepad)
- **Linux**: `xdg-open` (user's configured default)

The CLI returns immediately after launching the application.

---

## `open <slug>`

Open a project's directory in the OS file manager.

### Arguments

| Arg | Required | Type |
|-----|----------|------|
| `slug` | yes | string |

### Flags

None.

### Side Effects

Opens the project directory with the OS file manager:
- **macOS**: Finder (via `open`)
- **Windows**: Explorer (via `start`)
- **Linux**: Default file manager (via `xdg-open`)

The CLI returns immediately after launching the file manager.

---

## `delete <slug>` (alias: `rm`)

Delete a project and its entire directory.

### Arguments

| Arg | Required | Type |
|-----|----------|------|
| `slug` | yes | string |

### Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--force` | bool | `false` | Skip confirmation prompt |

### JSON Output

```json
{
  "status": "deleted",
  "slug": "my-project"
}
```

### Important

In non-interactive mode (piped/scripted), `--force` is **required**. Without it, the command returns an error.

### Side Effects

- Removes the entire project directory (`rm -rf`)
- Regenerates `PROJECTS.md` registry

---

## `status`

Health check across all projects.

### Arguments

None.

### Flags

None (uses global `--json`).

### JSON Output

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

### Fields

| Field | Type | Description |
|-------|------|-------------|
| `slug` | string | Project identifier |
| `title` | string | Display name |
| `status` | string | `active`, `paused`, or `archived` |
| `has_git` | bool | Whether project directory is a git repository |
| `has_remote` | bool | Whether a git remote is configured |
| `uncommitted` | bool | Whether there are uncommitted changes |
| `has_project_md` | bool | Whether `PROJECT.md` exists |

---

## `push <slug>`

Full git workflow: init, stage, commit, create GitHub repo, push.

### Arguments

| Arg | Required | Type |
|-----|----------|------|
| `slug` | yes | string |

### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--message` | `-m` | string | `"Update project"` | Commit message |
| `--private` | | bool | `true` | Create private GitHub repo |
| `--no-github` | | bool | `false` | Skip GitHub repo creation |

### JSON Output

```json
{
  "status": "pushed",
  "slug": "my-project",
  "remote": "https://github.com/user/my-project"
}
```

### Workflow

1. If not a git repo: runs `git init`
2. Runs `git add -A`
3. If uncommitted changes: commits with provided message
4. If no remote and `--no-github` is false: creates GitHub repo via `gh repo create --source --push` (creates and pushes in one step)
5. Else if remote already exists: pushes to remote with `git push -u origin <branch>`

**Note:** Steps 4 and 5 are mutually exclusive — either a new repo is created (which includes the push), or an existing remote is pushed to.

### Requirements

- `gh` CLI must be installed and authenticated for GitHub repo creation
- Without `gh`, requires existing remote or `--no-github` flag

---

## Data Schemas

### ProjectMeta Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `title` | string | yes | Display name |
| `slug` | string | yes | Unique identifier (`^[a-z0-9]+(?:-[a-z0-9]+)*$`) |
| `status` | string | yes | One of: `active`, `paused`, `archived` |
| `tags` | []string | no | Categorization tags (omitted from JSON when empty) |
| `description` | string | no | Short description (omitted from JSON when empty) |
| `created_at` | string (RFC 3339) | yes | Creation timestamp |
| `updated_at` | string (RFC 3339) | yes | Last update timestamp |
| `git_remote` | string (URL) | no | Git remote URL, set by `push` (omitted from JSON when empty) |

### Config Schema (`~/.projects/config.toml`)

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `projects_dir` | string | `~/.projects/projects` | Where projects are stored |
| `editor` | string | `$EDITOR` or `"vim"` | Editor binary name |
| `github_username` | string | `""` | GitHub username for `push` repo creation |
| `auto_git_init` | bool | `true` | Auto git init on project create |

Both `github_username` and `auto_git_init` are prompted interactively during first-run setup.

### Per-Project Directory Structure

```
<slug>/
  PROJECT.md            # YAML frontmatter + markdown body (source of truth)
  USAGE.md              # Workspace guide — directory roles, conventions
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
