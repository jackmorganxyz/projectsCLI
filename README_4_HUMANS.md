# projects — README for Humans 🧑‍💻

> Less chaos, more shipping.

A terminal-native project manager with a nice TUI, git superpowers, and just the right amount of personality. Built with [Charmbracelet](https://charm.sh) libraries because your terminal deserves nice things.

---

## Table of Contents

- [The Pitch](#-the-pitch)
- [Installation](#-installation)
- [Getting Started](#-getting-started)
- [What You Get](#-what-you-get)
- [Command Reference](#-command-reference)
- [Multi-Account Folders](#-multi-account-folders)
- [Configuration](#-configuration)
- [The Personality](#-the-personality)
- [Environment Variables](#-environment-variables)
- [Tips & Tricks](#-tips--tricks)

---

## 🎯 The Pitch

You have projects everywhere. Some in `~/code`, some in `~/Desktop/random-idea`, some you forgot about entirely. Half of them don't have a README. None of them have a consistent structure. It's fine. We've all been there. (It's not fine.)

**projects gives every project a home.**

One command to scaffold. One command to push. A dashboard to see everything at a glance. And when you pipe it to another tool, it automatically switches from pretty TUI output to clean JSON — because machines have feelings too.

---

## 📦 Installation

### Homebrew (recommended)

```sh
brew install jackmorganxyz/tap/projects
```

### Quick install script (for the impatient)

```sh
curl -sSL https://raw.githubusercontent.com/jackmorganxyz/projectsCLI/main/install.sh | sh
```

Auto-detects your OS and architecture. Installs to `/usr/local/bin`. It's smarter than it looks.

### Build from source (respect)

```sh
git clone https://github.com/jackmorganxyz/projectsCLI.git
cd projectsCLI
make build
make install  # copies to /usr/local/bin
```

---

## 🏁 Getting Started

Let's walk through the full lifecycle. From zero to pushed-to-GitHub in about 60 seconds. Grab a coffee — actually, you won't need it. This is fast.

### 1. Create a project

```sh
projects create --title "My API" --tags "go,api" --description "The API that does the thing"
```

The slug is auto-generated from the title (`my-api`). Or provide one explicitly:

```sh
projects create my-api --title "My API" --tags "go,api" --description "The API that does the thing"
```

You'll see something like:

```
 Created project "my-api" — Fresh project, who dis?
  Directory  ~/.projects/projects/my-api
  Created    2025-02-25

Tip: 'projects push <slug>' handles git init, commit, and GitHub in one step.
```

Behind the scenes, this created a whole scaffold for you (you're welcome):

```
~/.projects/projects/my-api/
  PROJECT.md            # Your project's identity card (YAML frontmatter + markdown)
  USAGE.md              # How this workspace works — read-me-first for agents and humans
  docs/README.md        # Documentation starts here
  memory/MEMORY.md      # Persistent notes (great for AI agent context)
  context/CONTEXT.md    # Architecture decisions and context
  tasks/TODO.md         # Task tracking
  code/                 # Your actual code
  private/              # Gitignored — secrets, drafts, local stuff
  .gitignore
```

Git was automatically initialized and the scaffold was committed.

**Got Claude Code or Codex CLI installed?** After scaffolding, you'll be asked if you want to spawn an AI agent to fill out the template files. Give it a prompt like *"Build a REST API service in Go with user authentication"* and it'll get to work in your new project directory. Or just skip it — your call.

### 2. See all your projects

```sh
projects ls
```

This launches an interactive dashboard in your terminal. Navigate with `j`/`k` (or arrow keys), press `Enter` to select, `q` to quit. You'll see a table with your project slugs, titles, statuses, and creation dates — all styled in a violet-and-emerald color palette that frankly goes unreasonably hard for a CLI tool.

**Select a project and do stuff.** When you press Enter on a project, you'll get a dropdown of actions — view, edit, open, status, push, update, move, or delete. Pick one and it runs. No more memorizing slugs and typing separate commands.

No projects yet? You'll get a gentle nudge:

```
It's quiet in here... too quiet. Run 'projects create <slug>' to fix that.
```

### 3. View project details

```sh
projects view my-api
```

In your terminal, this opens a scrollable view with all your project metadata and the markdown body from `PROJECT.md`. Page up/down to scroll, `q` to exit.

### 4. Edit your project

```sh
projects edit my-api
```

This launches an interactive file browser — navigate folders, pick any text file in your project. Then choose how you want to edit it:

```
  How would you like to edit this file?
  > Manual edit (open in an editor)
    Agent edit (AI-assisted editing)
```

**Manual edit** opens the file in your preferred editor. On first run you'll pick from installed editors (Cursor, VS Code, TextEdit, Vim, etc.) — now clearly labeled as `(terminal)` or `(GUI)` so you know what you're getting. The choice is saved so you won't be asked again. Want to re-pick? Use `--editor-picker`.

**Agent edit** spawns Claude Code or Codex CLI to edit the file for you. Describe what you want changed, and the AI gets to work. This option only appears if you have `claude` or `codex` installed.

```
  Describe the changes you want
  The agent will edit: TODO.md
  > Add tasks for building the authentication module
```

Want to skip all the pickers? Use `--editor`:

```sh
projects edit my-api --editor vim
```

In non-interactive mode (piped stdin/stdout), `edit` defaults to opening `PROJECT.md` with your saved editor.

### 5. Open the project folder

```sh
projects open my-api
```

Opens the project directory in Finder (macOS), Explorer (Windows), or your default file manager (Linux). Browse your files, drag and drop, do your thing.

### 6. Push to GitHub

```sh
projects push my-api -m "Initial commit"
```

This is the magic one. ✨ It handles the entire git workflow:

1. Ensures git is initialized
2. Stages all changes
3. Commits with your message
4. If there's no remote: creates a GitHub repo using `gh` CLI
5. Pushes to the remote

You'll see:

```
 Changes committed. Shipped it!
 Repository created: https://github.com/you/my-api
```

### 7. Check the health of everything

```sh
projects status
```

Your morning standup, minus the standing. Get a health check table across all your projects:

```
 Project Health

Slug       Status   Git   Remote   Clean
my-api     active   yes   yes      clean
side-proj  paused   yes   no       dirty
old-thing  archived no    -        -
```

Green means good, amber means needs attention. At a glance, you know what needs a push. No more "wait, did I commit that?" anxiety.

---

## 🗂️ What You Get

Every project created by projects follows the same structure. Consistency is a love language.

| File/Directory | Purpose |
|---|---|
| `PROJECT.md` | Source of truth. YAML frontmatter (metadata) + markdown body (notes, plans, anything) |
| `USAGE.md` | The "how this workspace works" guide — conventions, directory roles, agent instructions |
| `docs/` | Project documentation |
| `memory/MEMORY.md` | Persistent notes — great for AI agents that need project context |
| `context/CONTEXT.md` | Architecture decisions, design rationale |
| `tasks/TODO.md` | Task tracking with markdown checkboxes |
| `code/` | Your actual code goes here |
| `private/` | Gitignored by default. Local secrets, drafts, scratch files |

The `PROJECT.md` frontmatter looks like this:

```yaml
---
title: "My API"
slug: "my-api"
status: "active"
tags:
  - go
  - api
description: "The API that does the thing"
created_at: "2025-02-25T00:00:00Z"
updated_at: "2025-02-25T00:00:00Z"
git_remote: "https://github.com/you/my-api"
---

# My API

Your markdown content here. Plans, notes, links, whatever you want.
```

---

## 🔧 Command Reference

### `create [slug]`

Scaffold a new project. Another masterpiece begins.

```sh
projects create --title "My Project"                                          # slug auto-generated: my-project
projects create my-project --title "My Project" --tags "web,react"            # explicit slug
projects create my-project --status paused                                     # start it paused if you're not ready yet
```

**Flags:** `--title`, `--description`, `--tags` (comma-separated), `--status` (active/paused/archived)

After scaffolding, if Claude Code or Codex CLI is installed, you'll be asked if you want to spawn an AI agent to fill out the template files. Give it a prompt and let it do the heavy lifting.

### `list` (alias: `ls`)

See all your projects. Feel the satisfaction.

```sh
projects ls           # interactive TUI dashboard — select a project to act on it
projects ls --json    # JSON array of all projects
projects ls | jq '.'  # auto-switches to JSON when piped
```

In the dashboard, press Enter on a project to get a dropdown of actions (view, edit, open, etc.).

### `view <slug>`

See project details.

```sh
projects view my-project          # scrollable TUI
projects view my-project --json   # full project JSON
```

### `edit <slug>`

Browse and edit any file in a project — manually or with an AI agent.

```sh
projects edit my-project                    # file browser + edit mode picker
projects edit my-project --editor vim       # skip all pickers, use vim
projects edit my-project --editor-picker    # force re-show the editor picker
```

Interactively browse files, then choose "Manual edit" (opens in your editor) or "Agent edit" (spawns Claude Code or Codex CLI with a prompt). The editor choice is auto-detected from installed apps — now labeled as `(terminal)` or `(GUI)` — and saved to config on first pick.

**Flags:** `--editor` (editor command, bypasses all pickers), `--editor-picker` (force re-pick editor)

### `open <slug>`

Open the project folder in your file manager.

```sh
projects open my-project
```

Opens in Finder (macOS), Explorer (Windows), or your default file manager (Linux).

### `load <slug>`

Export project data for scripts and agents.

```sh
projects load my-project --json      # structured JSON (default)
projects load my-project --export    # shell export statements
projects load my-project --bash      # eval-able bash variables
```

Use it in scripts:

```sh
eval $(projects load my-project --export)
echo "Working on $PROJECT_TITLE in $PROJECT_DIR"
```

### `delete <slug>` (alias: `rm`)

Delete a project and its entire directory. We'll make sure you really mean it.

```sh
projects delete my-project          # asks for confirmation with a fun prompt
projects rm my-project --force      # skip confirmation (brave)
```

The confirmation prompt is delightfully dramatic:

```
Nuke "my-project" from orbit? This cannot be undone.
```

Cancel and you'll get reassurance: *"Crisis averted."*

### `status`

Health check across every project.

```sh
projects status          # colored table
projects status --json   # structured health data
```

Shows git init status, remote configuration, and whether there are uncommitted changes.

### `push <slug>`

The full git workflow in one command. Chef's kiss. 🤌

```sh
projects push my-project -m "Add user auth"
projects push my-project --private=false    # create a public repo
projects push my-project --no-github        # skip GitHub repo creation
```

**Flags:** `-m` (commit message), `--private` (default `true`), `--no-github`

Requires `gh` CLI for GitHub repo creation. If you already have a remote, it just pushes. If the project is in a folder with a GitHub account, `gh auth` is switched automatically before pushing.

### `update <slug>`

Update project metadata without opening an editor.

```sh
projects update my-project --title "New Name"
projects update my-project --status paused
projects update my-project --tags "go,api,v2" --description "Updated description"
```

**Flags:** `--title`, `--description`, `--status` (active/paused/archived), `--tags` (comma-separated)

At least one flag is required. The `updated_at` timestamp is set automatically.

### `folder add <name>`

Create a named folder tied to a GitHub account.

```sh
projects folder add work --account work-username       # explicit account
projects folder add personal                            # interactive: picks from gh auth accounts
```

If you omit `--account` and you have multiple accounts in `gh auth`, you'll get an interactive picker. If you only have one account, it's auto-selected.

### `folder list` (alias: `ls`)

See all configured folders.

```sh
projects folder list
```

### `folder remove <name>` (alias: `rm`)

Remove a folder from config. The directory and its projects are **not** deleted.

```sh
projects folder remove work
```

### `move <slug>`

Move a project between folders, or back to the top level.

```sh
projects move my-project --folder work       # move into a folder
projects move my-project --folder ""         # move to top level
```

---

## 📂 Multi-Account Folders

Got a work GitHub account and a personal one? Folders let you organize projects by account. When you push, the CLI switches `gh auth` to the right account automatically.

### Setup

```sh
# Add folders for each account
projects folder add work --account work-org
projects folder add personal --account my-personal-gh

# Or let the CLI pick from your gh keychain
projects folder add work
# → Interactive picker shows your authenticated accounts
```

### Creating projects in folders

```sh
projects create --title "Company API" --folder work
projects create --title "Side Project" --folder personal
```

### Pushing — auth switches automatically

```sh
projects push company-api -m "Add endpoints"
# → Switches to "work-org" account, creates repo under that account, pushes
```

### Moving existing projects

```sh
projects move old-project --folder personal
```

### How it works

- Each folder is a subdirectory under your projects directory (e.g. `~/.projects/projects/work/`)
- The `--folder` flag works on any command to scope it to that folder
- On `push`, the CLI runs `gh auth switch --user <account>` before creating repos or pushing
- `list` and `status` show a Folder column when you have folders configured
- Projects without a folder live at the top level — everything is backward compatible

---

## ⚙️ Configuration

Config lives at `~/.projects/config.toml`:

```toml
# Where your projects live
projects_dir = "~/.projects/projects"

# Editor for `edit` command (auto-saved on first pick)
editor = "cursor"

# GitHub username for repo creation (prompted during first-run setup)
github_username = "my-username"

# Automatically git init new projects (prompted during first-run setup)
auto_git_init = true

# Multi-account folders (added via `projects folder add`)
[[folders]]
name = "work"
github_account = "work-org"

[[folders]]
name = "personal"
github_account = "my-personal-gh"
```

The config file is created automatically the first time you run any command. You'll be prompted for your GitHub username and git preference during setup. All fields are optional — sensible defaults are built in. Folders are managed via the `folder` command — you don't need to edit TOML by hand.

---

## 🎭 The Personality

projects believes developer tools should spark joy, not existential dread. You'll encounter random messages throughout:

**When you create a project:**
- *"Fresh project, who dis?"*
- *"Another masterpiece begins."*
- *"Built different."*
- *"The world needs this."*

**When you push:**
- *"Shipped it!"*
- *"To the cloud and beyond!"*
- *"Chef's kiss."*

**When you delete (and cancel):**
- *"Crisis averted."*
- *"Good call, that one's a keeper."*

**When nothing needs committing:**
- *"Squeaky clean."*
- *"All caught up. Maybe go touch grass?"*

**While loading:**
- *"Summoning the code elves..."*
- *"Reticulating splines..."*
- *"Warming up the flux capacitor..."*

Random tips appear ~30% of the time after commands. They're genuinely useful, not just noise. We respect your terminal real estate.

Don't want emoji? Set `NO_EMOJI=1` or `TERM=dumb`. We won't judge. (We'll judge a little.)

---

## 🌍 Environment Variables

| Variable | What it does |
|----------|-------------|
| `NO_EMOJI` | Set to any value to disable emoji output |
| `TERM=dumb` | Also disables emoji |

---

## 💡 Tips & Tricks

**Pipe anything to jq.** When projects detects it's not writing to a terminal, it auto-switches to JSON. It's psychic like that:

```sh
projects ls | jq '.[].meta.slug'
```

**Use `load --export` in shell scripts** to get project metadata as environment variables:

```sh
eval $(projects load my-project --export)
cd "$PROJECT_DIR/code"
```

**The `--json` flag works on every command.** Even `create` and `delete` return structured output you can parse. Everything is an API if you believe hard enough.

**Status is your morning check-in.** Run `projects status` to see which projects need attention — dirty repos, missing remotes, all at a glance. Better than standup.

**The `private/` directory is your safe space.** It's gitignored by default. Throw API keys, draft notes, scratch files in there. They'll never accidentally get committed. Your secrets are safe with us.

---

Built with [Cobra](https://github.com/spf13/cobra), [Bubble Tea](https://github.com/charmbracelet/bubbletea), [Lip Gloss](https://github.com/charmbracelet/lipgloss), and mass quantities of caffeine. ☕
