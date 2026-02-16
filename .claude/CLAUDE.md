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

CI is the merge gate, not GitHub approvals. PR review is automated — the `/delegate` pipeline launches a review agent (`.manager/agents/reviewer.md`) after creating the PR. The reviewer reads the diff, checks spec compliance and coding standards, comments findings, and returns a verdict. If approved and CI is green, merge. If changes needed, report back without merging.

`/review <id>` is available as a manual override for re-reviewing after changes.

## Directory Structure

```
.claude/commands/     # slash commands (manager operations)
.manager/
  agents/             # sub-agent persona files (backend, frontend, proto, testing, reviewer)
  specs/              # work item specifications (created by /plan)
  blockers/           # blocker reports (written by stuck sub-agents)
  staging/            # temp: .claude/ files staged by sub-agents (copied during git-ops)
```

## Conventions

### Specs
Work items live in `.manager/specs/<id>-<name>.md`. The spec is the contract — sub-agents must deliver exactly what the spec describes, nothing more. Each spec includes a `Target Repo` field indicating which zarldev repo the work targets.

### Sub-Agent Launching
Sub-agents are launched via the Task tool (`subagent_type: "general-purpose"`, `run_in_background: true`). They cannot nest `claude` CLI sessions. Sub-agents only write code and run tests — they do NOT run git commands, create PRs, or comment on issues.

### Git Operations
The manager handles all git operations after a sub-agent completes:
1. Copy staging files (`.manager/staging/` → `.claude/`) if present
2. `git add` + `git commit`
3. `git push`
4. `gh pr create`
5. Review agent evaluates diff against spec and standards
6. `gh pr merge --squash --delete-branch` (only if review approved)

### Blockers
If a sub-agent gets stuck, it writes a blocker file to `.manager/blockers/<id>-<name>.md` using the Write tool. The manager picks these up via `/status`.

### Branches and Worktrees
Sub-agents work on branches named `work/<id>-<name>`. **ALWAYS use git worktrees** — never switch branches on the main checkout. This keeps the main checkout on `main` and avoids dirty-state issues.

```bash
git branch work/<id>-<name>
git worktree add .worktrees/<id>-<name> work/<id>-<name>
```

Copy `settings.local.json` to the worktree (it's gitignored):
```bash
cp <repo>/.claude/settings.local.json <worktree>/.claude/settings.local.json
```

Never `git checkout` to a work branch on the main clone. The main checkout must always stay on `main`.

### .claude/ Write Protection
Sub-agents cannot write to `.claude/` directories (built-in Claude Code security boundary). If a work item requires `.claude/` changes, the sub-agent writes to `.manager/staging/` instead. The manager copies these during git-ops.

### No Co-Authored-By
Commits never include co-authored-by lines.

## Slash Commands

| Command | Purpose |
|---------|---------|
| `/discuss <topic>` | Gather requirements, explore edge cases |
| `/plan` | Decompose work into specs and GH issues |
| `/delegate <id>` | Launch sub-agent, auto-review, git ops |
| `/status` | Check progress, blockers, active agents |
| `/review <id>` | Manual re-review of sub-agent output |

### Mandatory Slash Command Usage

**ALWAYS use slash commands for orchestration operations.** Never perform these actions manually:

- **Creating specs**: ALWAYS use `/plan`. Never write spec files directly.
- **Launching sub-agents**: ALWAYS use `/delegate <id>`. Never call the Task tool directly to launch sub-agents.
- **Reviewing output**: ALWAYS use `/review <id>`. Never manually read diffs and comment.
- **Checking progress**: ALWAYS use `/status`. Never manually check agent output files.

The slash commands encode the full workflow — branch creation, persona loading, git ops, issue management. Bypassing them skips critical steps and breaks the orchestration contract.
