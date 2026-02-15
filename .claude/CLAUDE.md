# claude-mngr — Manager/Sub-Agent Orchestration Hub

This repository is the orchestration layer for a two-tier agent system. It contains no application code — only manager configuration, slash commands, agent personas, and work item specs.

## Architecture

- **Manager** (interactive) decomposes work via slash commands, delegates to sub-agents, reviews output, and handles all git operations.
- **Sub-agents** (headless) write code and run tests. They do NOT handle git ops — the manager does that.

### Repo Separation

| Repo | Purpose |
|------|---------|
| `zarldev/claude-mngr` | Orchestration hub (this repo) |
| `zarldev/<project>` | Application code (e.g. `zarldev/tsk`) |

### Workspaces

| Path | Purpose |
|------|---------|
| `~/src/claude-mngr/` | Manager repo (orchestration) |
| `~/src/zarldev/<repo>/` | Sub-agent working directories |

### Identity

Single identity — `zarldev` (Bruno's GitHub account). All operations use the default `gh` auth.

- **Git author**: default git config (Bruno's identity)
- **All GitHub ops**: bare `gh` (no token overrides needed)
- No bot accounts, no PATs, no token env vars

### Review Strategy

CI is the merge gate, not GitHub approvals. The `/review` command does the actual code review (reads diff, checks spec compliance) and comments findings on the PR. If review passes and CI is green, merge directly. No formal GitHub approval required.

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
Work items live in `.manager/specs/<id>-<name>.md`. The spec is the contract — sub-agents must deliver exactly what the spec describes, nothing more. Each spec includes a `Target Repo` field indicating which zarldev repo the work targets.

### Sub-Agent Launching
Sub-agents are launched via the Task tool (`subagent_type: "general-purpose"`, `run_in_background: true`). They cannot nest `claude` CLI sessions. Sub-agents only write code and run tests — they do NOT run git commands, create PRs, or comment on issues.

### Git Operations
The manager handles all git operations after a sub-agent completes:
1. `git add` + `git commit`
2. `git push`
3. `gh pr create`
4. `gh issue comment` (review findings)
5. `gh pr merge --squash --delete-branch`

### Blockers
If a sub-agent gets stuck, it writes a blocker file to `.manager/blockers/<id>-<name>.md` using the Write tool. The manager picks these up via `/status`.

### Branches
Sub-agents work on branches named `work/<id>-<name>`. For parallel agents on the same repo, worktrees are used within the clone at `~/src/zarldev/<repo>/`.

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

### Mandatory Slash Command Usage

**ALWAYS use slash commands for orchestration operations.** Never perform these actions manually:

- **Creating specs**: ALWAYS use `/plan`. Never write spec files directly.
- **Launching sub-agents**: ALWAYS use `/delegate <id>`. Never call the Task tool directly to launch sub-agents.
- **Reviewing output**: ALWAYS use `/review <id>`. Never manually read diffs and comment.
- **Checking progress**: ALWAYS use `/status`. Never manually check agent output files.

The slash commands encode the full workflow — branch creation, persona loading, git ops, issue management. Bypassing them skips critical steps and breaks the orchestration contract.
