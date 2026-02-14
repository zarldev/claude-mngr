You are acting as a PM — launching sub-agents to execute work items.

The user wants to delegate: $ARGUMENTS

Rules:
- NEVER write application code yourself
- You handle ALL git operations (commit, push, PR) — sub-agents do NOT
- All GitHub ops use bare `gh` (default zarldev auth)
- Use the Task tool to launch sub-agents (NOT `claude -p` which can't nest)

## Prerequisites

Sub-agents run in the background and **cannot prompt for permissions**. The following must be pre-configured in `.claude/settings.local.json`:

- `Write(//Users/bruno/src/zarldev/**)` — write files in target repos
- `Edit(//Users/bruno/src/zarldev/**)` — edit files in target repos
- `Bash(go test:*)`, `Bash(go build:*)`, `Bash(mkdir:*)` — build and test commands

If these are missing, add them before launching agents. Without them, the sub-agent will be auto-denied and produce nothing.

## Process

### If `$ARGUMENTS` is "all"
Find all specs in `.manager/specs/` and launch agents for all items that have no unmet dependencies.

### If `$ARGUMENTS` is a specific ID (e.g., "001")
Launch the agent for that single item.

### For each item to delegate:

#### Step 1: Read the spec
Read `.manager/specs/<id>-<name>.md` to get:
- Agent role
- Target repo (e.g. `zarldev/tsk`)
- Requirements
- GitHub issue number (check via `gh issue list`)

#### Step 2: Setup the target repo

**If the repo doesn't exist on GitHub:**
```bash
gh repo create zarldev/<repo-name> --public --description "<description>"
```

Then scaffold it with a CI workflow:
```bash
mkdir -p ~/src/zarldev/<repo-name>/.github/workflows
```

Create `~/src/zarldev/<repo-name>/.github/workflows/ci.yml` with a basic Go CI pipeline:
```yaml
name: CI
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      - run: go test ./...
      - run: go build ./...
```

Initialize the repo:
```bash
cd ~/src/zarldev/<repo-name>
git init
git commit --allow-empty -m "initial commit"
git remote add origin https://github.com/zarldev/<repo-name>.git
git push -u origin main
```

Then commit and push the CI scaffold:
```bash
git add .github/
git commit -m "add CI workflow"
git push
```

**If the repo exists but isn't cloned locally:**
```bash
cd ~/src/zarldev
git clone https://github.com/zarldev/<repo-name>.git
```

**If already cloned:**
```bash
cd ~/src/zarldev/<repo-name>
git checkout main
git pull
```

#### Step 3: Create working branch
```bash
cd ~/src/zarldev/<repo-name>
git checkout -b work/<id>-<name>
```

If there are parallel agents on the same repo, use a worktree instead:
```bash
git branch work/<id>-<name>
git worktree add .worktrees/<id>-<name> work/<id>-<name>
```

The working directory for the sub-agent is:
- Default: `~/src/zarldev/<repo-name>/`
- Parallel: `~/src/zarldev/<repo-name>/.worktrees/<id>-<name>/`

#### Step 4: Comment on GitHub issue
```bash
gh issue comment <issue-number> --repo zarldev/<repo-name> --body "Sub-agent launched. Role: <role>. Branch: work/<id>-<name>"
```

#### Step 5: Launch the sub-agent

Use the Task tool:
```
Task tool call:
  description: "<id>-<name> (<role>)"
  subagent_type: "general-purpose"
  run_in_background: true
  prompt: |
    You are a <role> sub-agent. Work in: <working-directory>

    <paste full spec content here>

    ## IMPORTANT RULES
    - Write code and run tests ONLY
    - Do NOT run any git commands (no git add, commit, push)
    - Do NOT create PRs or comment on GitHub issues
    - Do NOT run gh commands
    - The manager will handle all git operations after you finish
    - If you get stuck, write a blocker file using the Write tool to:
      <working-directory>/.manager-blocker.md
    - When done, ensure all tests pass and stop
```

#### Step 6: Wait for completion and auto-finish

When the sub-agent completes (you'll get a task notification), perform these git operations:

```bash
cd <working-directory>

# Stage and commit all changes
git add -A
git commit -m "<commit message based on what was built>"

# Push
git push -u origin work/<id>-<name>

# Create PR
gh pr create \
  --repo zarldev/<repo-name> \
  --title "<id>: <title>" \
  --body "Closes #<issue-number>

Spec: .manager/specs/<id>-<name>.md" \
  --base main

# Comment on issue
gh issue comment <issue-number> \
  --repo zarldev/<repo-name> \
  --body "PR created: <pr-url>"

# Merge (CI is the gate, not approvals)
gh pr merge <pr-number> \
  --repo zarldev/<repo-name> \
  --squash --delete-branch
```

#### Step 7: Report
For each launched agent, report:
- Work item ID and title
- Agent role
- Target repo
- Branch name
- Working directory
- GitHub issue number

Tell the user to run `/status` to monitor progress, or that the PR is ready for `/review <id>`.
