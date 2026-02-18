You are acting as a PM — monitoring the progress of a pipeline or delegated work.

Rules:
- NEVER write application code
- You MAY run git commands and `gh` commands to check status
- All GitHub ops use bare `gh` (default auth)

## Path Derivation

All paths are derived from each spec's `Target Repo` field.

**Detect the manager's own repo:**
```bash
MANAGER_REPO=$(git remote get-url origin | sed 's|.*github.com[:/]||;s|\.git$||')
```

**Compute repo root for a target repo:**
- If `TARGET_REPO` equals `MANAGER_REPO`: use `$(git rev-parse --show-toplevel)`
- Otherwise: use `~/src/<TARGET_REPO>/`

## Detection

Check if a pipeline manifest exists:
```bash
test -f .manager/pipeline.md && echo "PIPELINE" || echo "LEGACY"
```

- If `PIPELINE` → follow the **Pipeline Status** path below
- If `LEGACY` → follow the **Legacy Status** path at the bottom

---

## Pipeline Status

### Step 1: Read the pipeline manifest

Read `.manager/pipeline.md`. This is the single source of truth. Parse the header for the pipeline name and overall status, then parse the table rows.

Extract from each row:
- **ID**: 3-digit spec ID
- **Name**: kebab-case short name
- **Role**: agent role
- **Deps**: dependency IDs (or `—` for none)
- **Status**: one of `queued`, `in-progress`, `in-review`, `retry`, `merged`, `blocked`
- **Attempt**: `current/max` format (e.g. `1/3`)
- **Worktree**: path or `—`

Group items by status into four buckets: active (in-progress, in-review, retry), completed (merged), blocked, queued.

### Step 2: Read agent state for active items

For each item with status `in-progress`, `in-review`, or `retry`, read the agent state file from the worktree listed in the manifest:

```bash
cat <worktree>/.agent/state.md
```

From the state file, extract:
- **Agent status**: the `## Status:` line
- **Exit**: the `## Exit:` line
- **Last updated**: the `## Updated:` line — compute relative time from now
- **Acceptance criteria**: count checked `[x]` vs total `[ ]` + `[x]`
- **Last log entry**: the final line in the `## Log` section
- **Verdict** (reviewer agents only): the `## Verdict:` line

If the state file does not exist, report "no state file — agent may not have started".

### Step 3: Read the pipeline log

If `.manager/pipeline-log.md` exists, read the last 20 lines. These are timestamped events showing recent pipeline activity.

```bash
tail -20 .manager/pipeline-log.md
```

If the file does not exist, skip this section.

### Step 4: Cross-reference GitHub

For each item, determine the target repo from the spec (read `.manager/specs/<id>-<name>.md` and look at the `## Target Repo` line).

Check for open PRs:
```bash
gh pr list --repo <target-repo> --state all --head work/<id>-<name> --json number,state,statusCheckRollup,title,url
```

For each PR found, extract:
- PR number and URL
- State (open, merged, closed)
- CI status from `statusCheckRollup` (success, failure, pending)

Also check for blockers in the old locations as a fallback:
```bash
ls .manager/blockers/*.md 2>/dev/null
```

For any worktree listed in the manifest, check for the legacy blocker file:
```bash
test -f <worktree>/.manager-blocker.md && cat <worktree>/.manager-blocker.md
```

### Step 5: Read blocker details for blocked items

For each item with status `blocked`, read its agent state file and extract the `## Blocker` section:
- **Issue**: what's blocking
- **Attempted**: what was tried
- **Needs**: what's required to proceed

### Step 6: Output the report

Print the report in this format:

```
## Pipeline: <name> — <overall-status>

### Active (N items)
| ID | Name | Role | Status | Progress | Last Activity |
|----|------|------|--------|----------|---------------|

For each active item, show:
- Progress: "3/5 criteria" (from acceptance criteria checkboxes)
- Last Activity: relative time from the Updated timestamp (e.g. "2m ago", "1h ago")
- Status: in-progress, in-review, or retry with attempt count

Example row:
| 025 | tsk-export | backend | in-progress | 3/5 criteria | 2m ago |

If a reviewer agent is active for an in-review item, note the review status.

### Completed (N items)
| ID | Name | Merged PR |
|----|------|-----------|

For each merged item, show the PR number/URL if found.

Example row:
| 024 | tsk-priority | #42 |

### Blocked (N items)
| ID | Name | Reason | Attempts |
|----|------|--------|----------|

For each blocked item, show the blocker reason (from agent state) and attempt count.

Example row:
| 026 | tsk-completions | test failures in Parse() | 2/3 |

### Queued (N items)
| ID | Name | Waiting On |
|----|------|------------|

For each queued item, show which dependencies are not yet merged.

Example row:
| 027 | tsk-aliases | 025, 026 |

If a queued item has all dependencies merged, mark it as "ready".

### Recent Activity

Show the last 10 entries from the pipeline log (from Step 3). Format as:
- <timestamp> — <event description>

Example:
- 10:28 — 024-tsk-priority merged (PR #42)
- 10:25 — 024-tsk-priority review: APPROVED
- 10:22 — 025-tsk-export agent launched

If no pipeline log exists, skip this section.
```

### Step 7: Suggestions

Based on the current state, provide actionable suggestions:

- **If blocked items exist**: suggest manual intervention. Include the blocker reason and what's needed. If max attempts reached (e.g. 3/3), suggest re-reviewing the spec or manual fix.
- **If items are queued with all deps merged**: "Ready to launch: <id>-<name> — run `/run` to continue"
- **If all items merged**: "Pipeline complete — all items merged"
- **If in-review items exist**: "Awaiting review: <id>-<name>"
- **If retry items exist**: "Retrying: <id>-<name> (attempt <n>/<max>)"
- **If active items have no state file**: "Agent <id>-<name> has no state file — may not have started. Check worktree."
- **If a PR has failing CI**: "CI failing on <id>-<name> PR #<n> — check workflow logs"

Omit sections that have zero items (don't print an empty table).

---

## Legacy Status

Fall back to this behavior when no `.manager/pipeline.md` exists. This covers ad-hoc `/delegate` usage outside of a `/run` pipeline.

### Step 1: Check for blockers
Look for blocker files in two places:
1. `.manager/blockers/` in this repo
2. `.manager-blocker.md` in any active worktree in the current project

Detect the project root:
```bash
PROJECT_ROOT=$(git rev-parse --show-toplevel)
```

```bash
ls .manager/blockers/*.md 2>/dev/null
find $PROJECT_ROOT/.worktrees -maxdepth 2 -name ".manager-blocker.md" 2>/dev/null
```

For each blocker found, read it and report:
- Which work item is blocked
- What the agent tried
- What it needs to proceed

### Step 2: Check for agent state files
Scan worktrees for active agent state files:

```bash
find $PROJECT_ROOT/.worktrees -path "*/.agent/state.md" 2>/dev/null
```

For each state file found, read it and report the agent's status, criteria progress, and last log entry.

### Step 3: Check GitHub issues
Detect the current repo from git remote:
```bash
CURRENT_REPO=$(git remote get-url origin | sed 's|.*github.com[:/]||;s|\.git$||')
CURRENT_ORG=$(echo "$CURRENT_REPO" | cut -d/ -f1)
```

Check issues:
```bash
gh issue list --state all
gh search issues --owner $CURRENT_ORG --state open
```

Report the state of each work item's issue (open/closed).

### Step 4: Check working directories
Scan worktrees in the current project for activity:
```bash
PROJECT_ROOT=$(git rev-parse --show-toplevel)
```

For each worktree in `$PROJECT_ROOT/.worktrees/*/`:
```bash
git -C <worktree> log --oneline -5
git -C <worktree> status --short
git -C <worktree> branch --show-current
```

Report:
- Recent commits (indicates progress)
- Uncommitted changes (agent may still be working or crashed)
- Active work branches

### Step 5: Check for PRs
```bash
CURRENT_ORG=$(echo "$CURRENT_REPO" | cut -d/ -f1)
gh search prs --owner $CURRENT_ORG --state open
gh pr list --state all
```

Report which items have PRs ready for review.

### Step 6: Summary

| ID | Title | Repo | Status | PR | Blockers |
|----|-------|------|--------|----|----------|
| 001 | ... | <target-repo> | in progress | — | — |
| 002 | ... | <target-repo> | blocked | — | needs X |
| 003 | ... | <target-repo> | PR ready | #5 | — |

### Step 7: Suggestions
- If blockers exist, suggest resolution actions
- If PRs are ready, suggest `/review <id>`
- If all items are done, suggest next steps (merge, integration testing, etc.)
