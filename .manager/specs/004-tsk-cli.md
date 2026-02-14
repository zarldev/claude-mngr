# 004: tsk — CLI task tracker

## Objective
Build a Go CLI tool called `tsk` that tracks tasks in a local JSON file.

## Context
Standalone CLI tool. No server, no frontend, no protobufs. Simple personal task tracker. Rebuild of the original prototype, now in its own repo.

## Requirements

### Task Structure
```go
type Task struct {
    ID        int
    Title     string
    Done      bool
    CreatedAt time.Time
}
```

### Commands
- `tsk add "buy milk"` — create a task with auto-incrementing ID, Done=false, CreatedAt=now
- `tsk list` — show all tasks; support filtering (e.g. `tsk list --done`, `tsk list --pending`)
- `tsk done <id>` — mark task as done
- `tsk rm <id>` — delete a task

### Storage
- JSON file at `~/.tasks.json`
- Concrete implementation — no interface abstraction
- Read entire file, modify in memory, write back
- Create file if it doesn't exist

### CLI Framework
Use the standard library (`os.Args`) or a minimal library — keep dependencies low.

### Output Format
- `list` should show ID, done status (checkbox or similar), title, and age/date
- Keep it readable in a terminal

### Project Structure
```
cmd/tsk/main.go       # entrypoint
internal/task/         # task type, storage, operations
```

### Tests
- Test task CRUD operations (add, list, done, rm)
- Test JSON file read/write round-trip
- Use a temp file for tests, not the real ~/.tasks.json
- Table-driven tests

## Target Repo
zarldev/tsk

## Agent Role
backend

## Files to Create
- go.mod (module name: github.com/zarldev/tsk)
- cmd/tsk/main.go
- internal/task/task.go
- internal/task/store.go
- internal/task/store_test.go

## Notes
- Keep it simple. No over-engineering.
- Follow error handling conventions: direct context, no "failed to" prefixes.
- Scope-based naming: smaller scope = shorter names.
