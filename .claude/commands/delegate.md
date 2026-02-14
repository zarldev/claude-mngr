You are acting as a PM — launching sub-agents to execute work items.

The user wants to delegate: $ARGUMENTS

Rules:
- NEVER write application code yourself
- You handle ALL git operations (commit, push, PR) — sub-agents do NOT
- Use the zarlbot PAT for all GitHub operations: `GH_TOKEN=$(cat /Users/bruno/.zarlbot/.ghpat)`
- Use the Task tool to launch sub-agents (NOT `claude -p` which can't nest)

## Process

### If `$ARGUMENTS` is "all"
Find all specs in `.manager/specs/` and launch agents for all items that have no unmet dependencies.

### If `$ARGUMENTS` is a specific ID (e.g., "001")
Launch the agent for that single item.

### For each item to delegate:

#### Step 1: Read the spec
Read `.manager/specs/<id>-<name>.md` to get:
- Agent role
- Target repo (e.g. `zarlbot/tsk`)
- Requirements
- GitHub issue number (check via `gh issue list`)

#### Step 2: Setup the target repo

**If the repo doesn't exist on GitHub:**
```bash
GH_TOKEN=$(cat /Users/bruno/.zarlbot/.ghpat) gh repo create zarlbot/<repo-name> --public --description "<description>"
```

Then scaffold it with a CI workflow:
```bash
mkdir -p ~/src/zarlbot/<repo-name>/.github/workflows
```

Create `~/src/zarlbot/<repo-name>/.github/workflows/ci.yml` with a basic Go CI pipeline:
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
cd ~/src/zarlbot/<repo-name>
git init
git -c user.name=zarlbot -c user.email=zarlbot@users.noreply.github.com commit --allow-empty -m "initial commit"
git remote add origin https://github.com/zarlbot/<repo-name>.git
git push -u origin main
```

Then commit and push the CI scaffold:
```bash
git add .github/
git -c user.name=zarlbot -c user.email=zarlbot@users.noreply.github.com commit -m "add CI workflow"
git push
```

**If the repo exists but isn't cloned locally:**
```bash
cd ~/src/zarlbot
git clone https://github.com/zarlbot/<repo-name>.git
```

**If already cloned:**
```bash
cd ~/src/zarlbot/<repo-name>
git checkout main
git pull
```

#### Step 3: Create working branch
```bash
cd ~/src/zarlbot/<repo-name>
git checkout -b work/<id>-<name>
```

If there are parallel agents on the same repo, use a worktree instead:
```bash
git branch work/<id>-<name>
git worktree add .worktrees/<id>-<name> work/<id>-<name>
```

The working directory for the sub-agent is:
- Default: `~/src/zarlbot/<repo-name>/`
- Parallel: `~/src/zarlbot/<repo-name>/.worktrees/<id>-<name>/`

#### Step 4: Comment on GitHub issue
```bash
GH_TOKEN=$(cat /Users/bruno/.zarlbot/.ghpat) gh issue comment <issue-number> --repo <target-repo> --body "Sub-agent launched. Role: <role>. Branch: work/<id>-<name>"
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
# Set up auth for push
cd <working-directory>
export GH_TOKEN=$(cat /Users/bruno/.zarlbot/.ghpat)

# Stage and commit all changes
git add -A
git -c user.name=zarlbot -c user.email=zarlbot@users.noreply.github.com commit -m "<commit message based on what was built>"

# Push
git push -u origin work/<id>-<name>

# Create PR
GH_TOKEN=$(cat /Users/bruno/.zarlbot/.ghpat) gh pr create \
  --repo zarlbot/<repo-name> \
  --title "<id>: <title>" \
  --body "Closes #<issue-number>

Spec: .manager/specs/<id>-<name>.md" \
  --base main

# Comment on issue
GH_TOKEN=$(cat /Users/bruno/.zarlbot/.ghpat) gh issue comment <issue-number> \
  --repo zarlbot/<repo-name> \
  --body "PR created: <pr-url>"
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
