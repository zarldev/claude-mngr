# claude-mngr — Manager/Sub-Agent Orchestration Hub

This repository is the orchestration layer for a two-tier agent system. It contains no application code — only manager configuration, slash commands, agent personas, and work item specs.

## Architecture

- **Manager** (interactive) decomposes work via slash commands, delegates to sub-agents, reviews output, and handles all git operations.
- **Sub-agents** (headless) write code and run tests. They do NOT handle git ops — the manager does that.
- **zarlbot** is the bot GitHub account. Sub-agents commit, push, and create PRs as zarlbot.

### Repo Separation

| Repo | Purpose | Owner |
|------|---------|-------|
| `zarldev/claude-mngr` | Orchestration hub (this repo) | Bruno |
| `zarlbot/<project>` | Application code | zarlbot |

Application code lives in dedicated repos under the `zarlbot/` GitHub account. Each project gets its own repo (e.g. `zarlbot/tsk`, `zarlbot/next-tool`).

### Workspaces

| Path | Purpose |
|------|---------|
| `~/src/claude-mngr/` | Manager repo (orchestration) |
| `~/src/zarlbot/<repo>/` | Sub-agent working directories (clones of zarlbot repos) |

### Identity

- **zarlbot token**: set `ZARLBOT_TOKEN` environment variable (e.g. in shell profile)
- **Git author**: `zarlbot <zarlbot@users.noreply.github.com>`
- All git ops (commit, push, PR, issue comments) use zarlbot identity via `GH_TOKEN=$ZARLBOT_TOKEN`

## Directory Structure

```
.claude/commands/     # slash commands (manager operations)
.manager/
  agents/             # sub-agent persona files (backend, frontend, proto, testing)
  specs/              # work item specifications (created by /plan)
  blockers/           # blocker reports (written by stuck sub-agents)
```

## Conventions

### Specs
Work items live in `.manager/specs/<id>-<name>.md`. The spec is the contract — sub-agents must deliver exactly what the spec describes, nothing more. Each spec includes a `Target Repo` field indicating which zarlbot repo the work targets.

### Sub-Agent Launching
Sub-agents are launched via the Task tool (`subagent_type: "general-purpose"`, `run_in_background: true`). They cannot nest `claude` CLI sessions. Sub-agents only write code and run tests — they do NOT run git commands, create PRs, or comment on issues.

### Git Operations
The manager handles all git operations after a sub-agent completes:
1. `git add` + `git commit` (as zarlbot author)
2. `git push` (using zarlbot PAT)
3. `gh pr create` (using zarlbot PAT)
4. `gh issue comment` (using zarlbot PAT)

### Blockers
If a sub-agent gets stuck, it writes a blocker file to `.manager/blockers/<id>-<name>.md` using the Write tool. The manager picks these up via `/status`.

### Branches
Sub-agents work on branches named `work/<id>-<name>`. For parallel agents on the same repo, worktrees are used within the clone at `~/src/zarlbot/<repo>/`.

### No Co-Authored-By
Commits never include co-authored-by lines.

## Slash Commands

| Command | Purpose |
|---------|---------|
| `/discuss <topic>` | Gather requirements, explore edge cases |
| `/plan` | Decompose work into specs and GH issues |
| `/delegate <id>` | Launch sub-agent, auto-finish git ops |
| `/status` | Check progress, blockers, active agents |
| `/review <id>` | Review sub-agent output against spec |
