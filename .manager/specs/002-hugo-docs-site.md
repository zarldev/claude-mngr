# 002: Hugo docs site for tsk

## Objective
Create a Hugo static site in `docs/` with a landing page and documentation for the `tsk` CLI tool.

## Context
The `tsk` CLI tool (001) is complete. It needs a public-facing site deployed via GitHub Pages. Hugo was chosen as the static site generator.

## Requirements

### Hugo Setup
- Initialize Hugo site in `docs/` directory
- Use a minimal, landing-page-style theme (e.g. Hugo Book, PaperMod, or similar — pick one that supports both a landing page and docs pages)
- Keep it clean and simple

### Landing Page
- Hero section: what `tsk` is (simple CLI task tracker)
- Quick install instructions (`go install github.com/zarldev/claude-mngr/cmd/tsk@latest`)
- Key features (3-4 bullet points)
- Link to full documentation

### Documentation Pages
- **Getting Started**: install, first task, basic workflow
- **Commands Reference**: `add`, `list`, `done`, `rm` with examples and flags
- **Configuration**: where the JSON file lives (`~/.tasks.json`), file format

### Content
- Write actual content based on the `tsk` implementation (see spec 001 or the code in `cmd/tsk/main.go` and `internal/task/`)
- Include realistic examples with sample output

## Agent Role
frontend

## Files to Create
- `docs/` — Hugo site root
- `docs/hugo.toml` — Hugo config
- `docs/content/` — Markdown content files
- `docs/layouts/` — Any custom layout overrides (if needed)

## Dependencies
- 001 (tsk CLI) should be merged or at least complete so docs reference accurate commands

## Notes
- No over-engineering. Content is more important than theme perfection.
- The theme should look professional but doesn't need to be custom-designed.
- Make sure `docs/` is self-contained — `hugo` should build from within that directory.
