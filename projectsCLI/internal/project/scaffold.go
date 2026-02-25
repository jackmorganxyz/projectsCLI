package project

import (
	"fmt"
	"os"
	"path/filepath"
)

// Scaffold creates the full directory tree and template files for a new project.
func Scaffold(projectsDir string, meta ProjectMeta) (string, error) {
	dir := filepath.Join(projectsDir, meta.Slug)

	if _, err := os.Stat(dir); err == nil {
		return "", fmt.Errorf("project directory already exists: %s", dir)
	}

	// Create directory tree.
	dirs := []string{
		dir,
		filepath.Join(dir, "docs"),
		filepath.Join(dir, "memory"),
		filepath.Join(dir, "context"),
		filepath.Join(dir, "tasks"),
		filepath.Join(dir, "code"),
		filepath.Join(dir, "private"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return "", fmt.Errorf("create directory %s: %w", d, err)
		}
	}

	// Write PROJECT.md with frontmatter.
	body := fmt.Sprintf("# %s\n\n%s\n", meta.Title, meta.Description)
	if err := WriteProjectFile(dir, meta, body); err != nil {
		return "", fmt.Errorf("write PROJECT.md: %w", err)
	}

	// Write template files.
	templates := map[string]string{
		filepath.Join(dir, "USAGE.md"):               usageTemplate(meta),
		filepath.Join(dir, "memory", "MEMORY.md"):    memoryTemplate(meta),
		filepath.Join(dir, "context", "CONTEXT.md"):  contextTemplate(meta),
		filepath.Join(dir, "tasks", "TODO.md"):       tasksTemplate(meta),
		filepath.Join(dir, "docs", "README.md"):      fmt.Sprintf("# %s\n\n%s\n", meta.Title, meta.Description),
		filepath.Join(dir, ".gitignore"):              "# Private files — never pushed to remote\nprivate/\n",
	}

	for path, content := range templates {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return "", fmt.Errorf("write %s: %w", filepath.Base(path), err)
		}
	}

	return dir, nil
}

func usageTemplate(meta ProjectMeta) string {
	return fmt.Sprintf(`# %s — Project Guide

> This file explains the structure of this project workspace.
> Read this first if you're an AI agent, a collaborator, or future-you coming back after a break.

## Directory Structure

| Directory | What goes here | Who uses it |
|-----------|---------------|-------------|
| `+"`code/`"+`| Source code, scripts, and anything that runs | You, your editor, your build tools |
| `+"`docs/`"+` | Documentation — READMEs, guides, API docs, specs | Humans and agents alike |
| `+"`memory/`"+` | Persistent notes and context that should survive between sessions | AI agents (and forgetful humans) |
| `+"`context/`"+` | Architecture decisions, design rationale, "why we did it this way" | Anyone making big decisions |
| `+"`tasks/`"+` | Task tracking, TODOs, checklists | Anyone working on the project |
| `+"`private/`"+` | Local-only files — secrets, API keys, drafts, scratch work | **Gitignored.** Never pushed. Your safe space |

## Key Files

- **PROJECT.md** — Source of truth. YAML frontmatter has all project metadata (title, status, tags). The body is freeform markdown for plans, notes, whatever.
- **USAGE.md** — You're reading it. The "how this place works" guide.
- **memory/MEMORY.md** — Append persistent notes here. Things you'd want to remember next time you (or an agent) open this project.
- **context/CONTEXT.md** — Document architectural decisions and rationale. "We chose X because Y" goes here.
- **tasks/TODO.md** — Track work with markdown checkboxes. Keep it simple.
- **docs/README.md** — The public-facing documentation for this project.

## Conventions

1. **Code goes in `+"`code/`"+`** — Not in the root, not scattered around. Keep it contained.
2. **Secrets go in `+"`private/`"+`** — API keys, .env files, credentials. It's gitignored for a reason.
3. **Notes go in `+"`memory/`"+`** — If you learn something useful about this project, write it down. Future you will thank past you.
4. **Decisions go in `+"`context/`"+`** — "Why did we use Postgres?" "Why is this a monorepo?" Document the reasoning, not just the choice.
5. **Don't fight the structure** — It's opinionated on purpose. Consistency across projects is the whole point.

## For AI Agents

If you're an AI agent working in this project:

- **Read `+"`PROJECT.md`"+` first** for project metadata and high-level context.
- **Read `+"`memory/MEMORY.md`"+`** for persistent notes from previous sessions.
- **Read `+"`context/CONTEXT.md`"+`** before making architectural decisions.
- **Check `+"`tasks/TODO.md`"+`** for current work items.
- **Write back to `+"`memory/MEMORY.md`"+`** when you learn something worth remembering.
- **Never put secrets in tracked files** — use `+"`private/`"+` for anything sensitive.
- **Put code in `+"`code/`"+`** — respect the structure.
`, meta.Title)
}

func memoryTemplate(meta ProjectMeta) string {
	return fmt.Sprintf(`# %s — Memory

Persistent notes and context for this project. Add anything here that should survive between sessions.

## Quick Facts

- **Created:** %s
- **Status:** %s

## Notes

_Nothing yet. Add notes as the project evolves._
`, meta.Title, meta.CreatedAt[:10], meta.Status)
}

func contextTemplate(meta ProjectMeta) string {
	return fmt.Sprintf(`# %s — Context & Decisions

Architecture decisions, design rationale, and the "why" behind this project.

## Overview

%s

## Decisions

_Document key decisions here as they're made. Format: what was decided, why, and what alternatives were considered._
`, meta.Title, meta.Description)
}

func tasksTemplate(meta ProjectMeta) string {
	return fmt.Sprintf(`# %s — Tasks

## Active

- [ ] Initial setup
- [ ] Define project goals in PROJECT.md

## Done

_Nothing yet. Ship something and check it off._
`, meta.Title)
}
