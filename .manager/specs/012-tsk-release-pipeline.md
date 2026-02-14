# 012: tsk release pipeline

## Objective
Add a `tsk version` command and a GitHub Actions workflow that builds cross-platform binaries and creates GitHub Releases on tag push.

## Context
`tsk` is a simple Go CLI with zero dependencies. No need for GoReleaser — a plain GitHub Actions workflow with `go build` and a build matrix handles everything. Version is injected at build time via ldflags.

## Requirements

### `tsk version` command
- Add `version` to the switch in `cmd/tsk/main.go`
- Print `tsk <version>` to stdout (e.g. `tsk v1.0.0`)
- Version comes from a package-level `var version = "dev"` in main.go
- At build time, inject via: `go build -ldflags "-X main.version=v1.0.0"`
- When built without ldflags (e.g. `go install`), prints `tsk dev`

### GitHub Actions release workflow
`.github/workflows/release.yml`:
- Trigger: push tag matching `v*` (e.g. `v1.0.0`)
- Build matrix:
  - `linux/amd64`, `linux/arm64`
  - `darwin/amd64`, `darwin/arm64`
  - `windows/amd64`
- Steps per matrix entry:
  1. Checkout
  2. Setup Go 1.23
  3. `go build` with ldflags injecting the tag as version
  4. Binary name: `tsk-<os>-<arch>` (`.exe` suffix for windows)
  5. Tar/gzip for linux/darwin, zip for windows
- After matrix completes:
  1. Create GitHub Release using the tag
  2. Upload all archives as release assets
  3. Auto-generate release notes from commits since last tag
- Use `gh release create` or `actions/create-release` — keep it simple
- Do NOT modify `ci.yml` or `deploy-docs.yml`

### Archive naming
```
tsk-linux-amd64.tar.gz
tsk-linux-arm64.tar.gz
tsk-darwin-amd64.tar.gz
tsk-darwin-arm64.tar.gz
tsk-windows-amd64.zip
```

Each archive contains a single binary named `tsk` (or `tsk.exe`).

### Update usage output
Update the `usage()` function to include the version command:
```
usage: tsk <command> [args]

commands:
  add <title>              add a new task
  list [--done|--pending]  list tasks
  done <id>                mark a task as done
  rm <id>                  remove a task
  version                  print version
```

### Tests
- Test that the version command doesn't panic (basic smoke test)
- Existing tests must continue to pass

## Target Repo
zarldev/tsk

## Agent Role
backend

## Files to Modify
- `cmd/tsk/main.go` (add version var, version command, update usage)

## Files to Create
- `.github/workflows/release.yml`

## Dependencies
None

## Notes
- Keep the workflow simple — no third-party actions beyond `actions/checkout` and `actions/setup-go`.
- Use `gh release create` in the workflow for release creation (gh is pre-installed on GitHub runners).
- The workflow should be readable by anyone — no clever tricks.
- First release will be `v0.1.0` (user tags manually to trigger).
