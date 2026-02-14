# 005: Hugo docs site + GitHub Pages deploy for tsk

## Objective
Create a Hugo docs site and GitHub Actions deployment workflow for the `tsk` CLI tool.

## Context
`tsk` is a Go CLI task tracker at `zarldev/tsk` (spec 004). It needs a public-facing landing page and documentation deployed to GitHub Pages.

## Requirements

### Hugo Site in `docs/`
- Initialize Hugo site in `docs/` directory
- Use a minimal theme (e.g. Hugo Book, PaperMod, or Ananke — pick one that supports landing + docs)
- `docs/` must be self-contained — `hugo` builds from within that directory

### Landing Page
- Hero: what `tsk` is (simple CLI task tracker, zero dependencies)
- Quick install: `go install github.com/zarldev/tsk/cmd/tsk@latest`
- Key features (3-4 bullets)
- Link to docs

### Documentation Pages
- **Getting Started**: install, first task, basic workflow
- **Commands Reference**: `add`, `list` (with `--done`/`--pending`), `done`, `rm` — with examples and sample output
- **Configuration**: `~/.tasks.json` location and file format

### Sample Output
Include realistic examples, e.g.:
```
$ tsk add "buy milk"
added task 1: buy milk

$ tsk list
  1 [ ] buy milk  (just now)

$ tsk done 1
task 1 marked done

$ tsk list
  1 [x] buy milk  (2m ago)
```

### Hugo Config
- `baseURL`: `https://zarldev.github.io/tsk/`
- Title: `tsk`

### GitHub Actions Deployment
- `.github/workflows/deploy-docs.yml`
- Triggers on push to `main` with path filter: `docs/**`
- Also supports `workflow_dispatch`
- Steps: checkout → setup Hugo → build from `docs/` → deploy to Pages
- Use `actions/deploy-pages` + `actions/upload-pages-artifact`
- Permissions: `pages: write`, `id-token: write`
- Do NOT modify the existing `.github/workflows/ci.yml`

### Theme Installation
- Add the theme as a git submodule under `docs/themes/`
- Ensure the theme is referenced in `docs/hugo.toml`

## Target Repo
zarldev/tsk

## Agent Role
frontend

## Files to Create
- `docs/hugo.toml`
- `docs/content/_index.md` (landing page)
- `docs/content/docs/getting-started.md`
- `docs/content/docs/commands.md`
- `docs/content/docs/configuration.md`
- `docs/themes/<theme>` (git submodule)
- `.github/workflows/deploy-docs.yml`

## Dependencies
- 004 (tsk CLI) must be merged

## Notes
- Content quality matters more than theme perfection.
- Keep it simple and professional.
- Write actual content based on the tsk implementation — read the code if needed.
- Hugo must be available on the agent's system (`brew install hugo` or `go install`).
