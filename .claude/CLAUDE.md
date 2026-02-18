# claude-mngr — Global Orchestration System

This repository is the source of truth for a two-tier agent orchestration system. It contains slash commands, agent personas, and the install script that deploys them globally to `~/.claude/`. No application code lives here.

## Architecture

- **Manager** (interactive) decomposes work via slash commands, delegates to sub-agents, reviews output, and handles all git operations.
- **Sub-agents** (headless) write code and run tests. They do NOT handle git ops — the manager does that.

### Global Deployment

Commands and agent personas are deployed to `~/.claude/` via `install.sh`, making them available in any project. The system is org-agnostic — all paths are derived from each spec's `Target Repo` field at runtime.

### Identity

Single GitHub identity. All operations use the default `gh` auth.

- **Git author**: default git config
- **All GitHub ops**: bare `gh` (no token overrides needed)
- No bot accounts, no PATs, no token env vars

### Review Strategy

CI is the merge gate, not GitHub approvals. PR review is automated — the `/delegate` pipeline launches a review agent (`~/.claude/agents/reviewer.md`) after creating the PR. The reviewer reads the diff, checks spec compliance and coding standards, comments findings, and returns a verdict. If approved and CI is green, merge. If changes needed, report back without merging.

`/review <id>` is available as a manual override for re-reviewing after changes.

## Directory Structure

```
.claude/commands/     # slash commands (deployed to ~/.claude/commands/ by install.sh)
.manager/
  agents/             # agent personas (deployed to ~/.claude/agents/ by install.sh)
  specs/              # reserved for temp specs during pipeline runs
  blockers/           # blocker reports (written by stuck sub-agents)
  staging/            # temp: .claude/ files staged by sub-agents (copied during git-ops)
install.sh            # deploys commands and agents to ~/.claude/
```

## Conventions

### Specs
Work item specs live in individual project repos, not here. Each spec includes a `Target Repo` field — the orchestration system derives all paths from this field, making commands org-agnostic.

### Sub-Agent Launching
Sub-agents are launched via the Task tool (`subagent_type: "general-purpose"`, `run_in_background: true`). They cannot nest `claude` CLI sessions. Sub-agents only write code and run tests — they do NOT run git commands, create PRs, or comment on issues.

### Git Operations
The manager handles all git operations after a sub-agent completes:
1. Copy staging files (`.manager/staging/` -> `.claude/`) if present
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
| `/run` | Autonomous pipeline — launch waves, review, merge |
| `/status` | Check progress, blockers, active agents |
| `/review <id>` | Manual re-review of sub-agent output |

### Mandatory Slash Command Usage

**ALWAYS use slash commands for orchestration operations.** Never perform these actions manually:

- **Creating specs**: ALWAYS use `/plan`. Never write spec files directly.
- **Launching sub-agents**: ALWAYS use `/delegate <id>` or `/run`. Never call the Task tool directly to launch sub-agents.
- **Reviewing output**: ALWAYS use `/review <id>`. Never manually read diffs and comment.
- **Checking progress**: ALWAYS use `/status`. Never manually check agent output files.

The slash commands encode the full workflow — branch creation, persona loading, git ops, issue management. Bypassing them skips critical steps and breaks the orchestration contract.
