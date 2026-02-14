# 003: GitHub Actions workflow for docs deployment

## Objective
Create a GitHub Actions workflow that auto-builds and deploys the Hugo site to GitHub Pages when docs change.

## Context
The Hugo site lives in `docs/` (spec 002). We need automated deployment so pushing docs changes to main triggers a build and publish.

## Requirements

### Workflow File
- `.github/workflows/deploy-docs.yml`
- Triggers on push to `main` with path filter: `docs/**`
- Also supports manual trigger (`workflow_dispatch`)

### Build Steps
1. Checkout repo
2. Setup Hugo (use `peaceiris/actions-hugo` or equivalent)
3. Build site from `docs/` directory
4. Deploy to GitHub Pages

### Deployment Strategy
- Use GitHub's built-in Pages deployment (`actions/deploy-pages`) with `actions/upload-pages-artifact`
- Configure Pages source as GitHub Actions (not `gh-pages` branch) — this is the modern approach
- Set appropriate permissions in the workflow (`pages: write`, `id-token: write`)

### GitHub Pages Setup
- Enable GitHub Pages on the repo via settings or `gh` CLI if possible
- Base URL should work with `https://zarldev.github.io/claude-mngr/`

## Agent Role
backend

## Files to Create
- `.github/workflows/deploy-docs.yml`

## Dependencies
- 002 (Hugo site) must exist for the workflow to have something to build

## Notes
- Keep the workflow simple and standard.
- Use pinned action versions for security.
- The Hugo config in `docs/hugo.toml` should have `baseURL` set to `https://zarldev.github.io/claude-mngr/`.
