# claude-mngr — Manager/Sub-Agent System

This repository uses a two-tier agent system: a **manager** (interactive) that decomposes and delegates work, and **sub-agents** (headless) that execute it.

## How It Works

- **Manager** runs interactively via slash commands (`/discuss`, `/plan`, `/delegate`, `/status`, `/review`). The manager never writes application code — it gathers requirements, creates specs, launches sub-agents, and reviews output.
- **Sub-agents** run headless in git worktrees with role-specific personas appended via `--append-system-prompt-file`. They write code, run tests, and create PRs.

Both the manager and sub-agents see this file. Manager-specific behavior is enforced by the slash commands, not here.

## Directory Structure

```
.claude/commands/     # slash commands (manager operations)
.manager/
  agents/             # sub-agent persona files (backend, frontend, proto, testing)
  specs/              # work item specifications (created by /plan)
  blockers/           # blocker reports (written by stuck sub-agents)
.worktrees/           # git worktrees for sub-agents (gitignored)
```

## Conventions (all agents)

### Specs
Work items live in `.manager/specs/<id>-<name>.md`. The spec is the contract — sub-agents must deliver exactly what the spec describes, nothing more.

### Blockers
If a sub-agent gets stuck, it writes a blocker file to `.manager/blockers/<id>-<name>.md` describing the issue, what was attempted, and what's needed. The manager picks these up via `/status`.

### GitHub Issues
Each work item has a corresponding GitHub issue. Sub-agents comment on their issue with progress updates and link their PR when done.

### Branches
Sub-agents work on branches named `work/<id>-<name>` in worktrees at `.worktrees/<id>-<name>`.

## Slash Commands

| Command | Purpose |
|---------|---------|
| `/discuss <topic>` | Gather requirements, explore edge cases |
| `/plan` | Decompose work into specs and GH issues |
| `/delegate <id>` | Launch sub-agent for a work item |
| `/status` | Check progress, blockers, active agents |
| `/review <id>` | Review sub-agent output against spec |
