# projectsCLI — README for Humans 🧑‍💻

> Less chaos, more shipping.

A terminal-native project manager with a nice TUI, git superpowers, and just the right amount of personality. Built with [Charmbracelet](https://charm.sh) libraries because your terminal deserves nice things.

---

## Table of Contents

- [The Pitch](#-the-pitch)
- [Installation](#-installation)
- [Getting Started](#-getting-started)
- [What You Get](#-what-you-get)
- [Command Reference](#-command-reference)
- [Configuration](#-configuration)
- [The Personality](#-the-personality)
- [Environment Variables](#-environment-variables)
- [Tips & Tricks](#-tips--tricks)

---

## 🎯 The Pitch

You have projects everywhere. Some in `~/code`, some in `~/Desktop/random-idea`, some you forgot about entirely. Half of them don't have a README. None of them have a consistent structure. It's fine. We've all been there. (It's not fine.)

**projectsCLI gives every project a home.**

One command to scaffold. One command to push. A dashboard to see everything at a glance. And when you pipe it to another tool, it automatically switches from pretty TUI output to clean JSON — because machines have feelings too.

---

## 📦 Installation

### Homebrew (recommended)

```sh
brew install jackmorganxyz/tap/projectsCLI
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
projectsCLI create my-api --title "My API" --tags "go,api" --description "The API that does the thing"
```

You'll see something like:

```
 Created project "my-api" — Fresh project, who dis?
  Directory  ~/.projects/projects/my-api
  Created    2025-02-25

Tip: 'projectsCLI push <slug>' handles git init, commit, and GitHub in one step.
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

### 2. See all your projects

```sh
projectsCLI ls
```

This launches an interactive dashboard in your terminal. Navigate with `j`/`k` (or arrow keys), press `Enter` to select, `q` to quit. You'll see a table with your project slugs, titles, statuses, and creation dates — all styled in a violet-and-emerald color palette that frankly goes unreasonably hard for a CLI tool.

No projects yet? You'll get a gentle nudge:

```
It's quiet in here... too quiet. Run 'projectsCLI create <slug>' to fix that.
```

### 3. View project details

```sh
projectsCLI view my-api
```

In your terminal, this opens a scrollable view with all your project metadata and the markdown body from `PROJECT.md`. Page up/down to scroll, `q` to exit.

### 4. Edit your project

```sh
projectsCLI edit my-api
```

Opens `PROJECT.md` in your configured editor (defaults to `$EDITOR`, falls back to `vim`). The frontmatter is YAML — edit the title, status, tags, description. The body is freeform markdown.

### 5. Push to GitHub

```sh
projectsCLI push my-api -m "Initial commit"
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

### 6. Check the health of everything

```sh
projectsCLI status
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

Every project created by projectsCLI follows the same structure. Consistency is a love language.

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

### `create <slug>`

Scaffold a new project. Another masterpiece begins.

```sh
projectsCLI create my-project
projectsCLI create my-project --title "My Project" --tags "web,react" --description "A cool thing"
projectsCLI create my-project --status paused  # start it paused if you're not ready yet
```

**Flags:** `--title`, `--description`, `--tags` (comma-separated), `--status` (active/paused/archived)

### `list` (alias: `ls`)

See all your projects. Feel the satisfaction.

```sh
projectsCLI ls           # interactive TUI dashboard
projectsCLI ls --json    # JSON array of all projects
projectsCLI ls | jq '.'  # auto-switches to JSON when piped
```

### `view <slug>`

See project details.

```sh
projectsCLI view my-project          # scrollable TUI
projectsCLI view my-project --json   # full project JSON
```

### `edit <slug>`

Open `PROJECT.md` in your editor.

```sh
projectsCLI edit my-project
```

Uses the `editor` from your config, then `$EDITOR`, then `vim`.

### `load <slug>`

Export project data for scripts and agents.

```sh
projectsCLI load my-project --json      # structured JSON (default)
projectsCLI load my-project --export    # shell export statements
projectsCLI load my-project --bash      # eval-able bash variables
```

Use it in scripts:

```sh
eval $(projectsCLI load my-project --export)
echo "Working on $PROJECT_TITLE in $PROJECT_DIR"
```

### `delete <slug>` (alias: `rm`)

Delete a project and its entire directory. We'll make sure you really mean it.

```sh
projectsCLI delete my-project          # asks for confirmation with a fun prompt
projectsCLI rm my-project --force      # skip confirmation (brave)
```

The confirmation prompt is delightfully dramatic:

```
Nuke "my-project" from orbit? This cannot be undone.
```

Cancel and you'll get reassurance: *"Crisis averted."*

### `status`

Health check across every project.

```sh
projectsCLI status          # colored table
projectsCLI status --json   # structured health data
```

Shows git init status, remote configuration, and whether there are uncommitted changes.

### `push <slug>`

The full git workflow in one command. Chef's kiss. 🤌

```sh
projectsCLI push my-project -m "Add user auth"
projectsCLI push my-project --private=false    # create a public repo
projectsCLI push my-project --no-github        # skip GitHub repo creation
```

**Flags:** `-m` (commit message), `--private` (default `true`), `--no-github`

Requires `gh` CLI for GitHub repo creation. If you already have a remote, it just pushes.

---

## ⚙️ Configuration

Config lives at `~/.projects/config.toml`:

```toml
# Where your projects live
projects_dir = "~/.projects/projects"

# Which editor 'edit' opens (defaults to $EDITOR, then vim)
editor = "nvim"

# GitHub org for repo creation (optional — omit for personal repos)
github_org = "my-org"

# Automatically git init new projects (default: true)
auto_git_init = true
```

The config file is created automatically the first time you run any command. All fields are optional — sensible defaults are built in. We're not here to make you configure things.

---

## 🎭 The Personality

projectsCLI believes developer tools should spark joy, not existential dread. You'll encounter random messages throughout:

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
| `EDITOR` | Default editor for the `edit` command |
| `NO_EMOJI` | Set to any value to disable emoji output |
| `TERM=dumb` | Also disables emoji |

---

## 💡 Tips & Tricks

**Pipe anything to jq.** When projectsCLI detects it's not writing to a terminal, it auto-switches to JSON. It's psychic like that:

```sh
projectsCLI ls | jq '.[].meta.slug'
```

**Use `load --export` in shell scripts** to get project metadata as environment variables:

```sh
eval $(projectsCLI load my-project --export)
cd "$PROJECT_DIR/code"
```

**The `--json` flag works on every command.** Even `create` and `delete` return structured output you can parse. Everything is an API if you believe hard enough.

**Status is your morning check-in.** Run `projectsCLI status` to see which projects need attention — dirty repos, missing remotes, all at a glance. Better than standup.

**The `private/` directory is your safe space.** It's gitignored by default. Throw API keys, draft notes, scratch files in there. They'll never accidentally get committed. Your secrets are safe with us.

---

Built with [Cobra](https://github.com/spf13/cobra), [Bubble Tea](https://github.com/charmbracelet/bubbletea), [Lip Gloss](https://github.com/charmbracelet/lipgloss), and mass quantities of caffeine. ☕
