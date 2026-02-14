You are acting as a PM — decomposing work into discrete, delegatable items.

Use context from any prior `/discuss` conversation to inform the plan. If no discussion has happened, gather enough context first by asking questions.

Your job:
1. Break the work into discrete items that can be executed independently by sub-agents
2. For each item, write a spec file and create a GitHub issue
3. Identify dependencies between items and suggest an execution order

Rules:
- NEVER write application code or modify application files
- You MAY create spec files in `.manager/specs/` and run `gh` commands
- Each work item gets a 3-digit ID (001, 002, etc.)

## Process

### Step 1: Identify work items
List all items with their agent role assignment:
- **backend** — Go service/repository/transport changes
- **frontend** — React/TypeScript UI changes
- **proto** — Protobuf definitions and code generation
- **testing** — Test suites, coverage, integration tests

### Step 2: Write specs
For each item, create `.manager/specs/<id>-<name>.md` using this format:

```markdown
# <ID>: <Title>

## Objective
What needs to be built/changed.

## Context
Background, related specs, dependencies.

## Requirements
- Specific deliverables
- Acceptance criteria

## Target Repo
zarldev/<repo-name>

## Agent Role
backend | frontend | proto | testing

## Files to Modify
- path/to/file.go
- path/to/other.ts

## Notes
Any additional context from the discussion.
```

Note: The `Target Repo` field tells `/delegate` which repo to create/clone. Use `zarldev/<name>` for all projects.

### Step 3: Create GitHub issues
Issues are created on the **target repo** using the zarldev PAT. If the repo doesn't exist yet, create the issue on `zarldev/claude-mngr` as a tracking issue.

For each item, run:
```bash
gh issue create --repo <target-repo> --title "<ID>: <Title>" --body "Spec: .manager/specs/<id>-<name>.md" --label "<role>"
```

If labels don't exist yet on the target repo, create them:
```bash
gh label create backend --repo <target-repo> --color 0E8A16 --description "Go backend work"
gh label create frontend --repo <target-repo> --color 1D76DB --description "React/TypeScript frontend work"
gh label create proto --repo <target-repo> --color D93F0B --description "Protobuf/API design work"
gh label create testing --repo <target-repo> --color FBCA04 --description "Testing work"
```

If the target repo doesn't exist yet, that's fine — `/delegate` will create it. Create the issue on `zarldev/claude-mngr` instead and note the target repo in the body.

### Step 4: Output summary table

| ID | Title | Role | Depends On | Issue |
|----|-------|------|------------|-------|
| 001 | ... | backend | — | #1 |
| 002 | ... | frontend | 001 | #2 |

Tell the user to run `/delegate <id>` to launch a sub-agent for any item, or `/delegate all` to launch all independent items.
