---
title: "Commands Reference"
weight: 2
---

# Commands Reference

`tsk` has four commands: `add`, `list`, `done`, and `rm`.

## Usage

```
tsk <command> [args]
```

Running `tsk` without arguments prints the usage summary.

---

## add

Create a new task.

```
tsk add <title>
```

The title is everything after `add`, joined with spaces. No quotes needed.

### Examples

```bash
$ tsk add Buy groceries
added task 1: Buy groceries

$ tsk add Fix the login timeout bug in auth service
added task 2: Fix the login timeout bug in auth service
```

### Notes

- Each task gets an auto-incrementing integer ID
- New tasks start with `Done = false`
- The creation timestamp is recorded automatically

---

## list

Show tasks.

```
tsk list [--done|--pending]
```

### Flags

| Flag | Description |
|------|-------------|
| (none) | Show all tasks |
| `--done` | Show only completed tasks |
| `--pending` | Show only incomplete tasks |

### Output Format

Each line shows:
- Task ID (right-aligned, 3 characters)
- Status checkbox (`[ ]` for pending, `[x]` for done)
- Title (left-aligned, 40 characters)
- Relative age (e.g., `just now`, `5m ago`, `2h ago`, `3d ago`)

### Examples

```bash
# show all tasks
$ tsk list
  1  [ ]  Buy groceries                            2h ago
  2  [x]  Fix login timeout bug                    1d ago
  3  [ ]  Update deployment docs                   5m ago

# show only pending tasks
$ tsk list --pending
  1  [ ]  Buy groceries                            2h ago
  3  [ ]  Update deployment docs                   5m ago

# show only completed tasks
$ tsk list --done
  2  [x]  Fix login timeout bug                    1d ago

# when there are no tasks
$ tsk list
no tasks
```

---

## done

Mark a task as completed.

```
tsk done <id>
```

### Examples

```bash
$ tsk done 1
task 1 marked done
```

### Errors

```bash
# task doesn't exist
$ tsk done 99
tsk: task 99 not found

# invalid id
$ tsk done abc
invalid id: abc
```

---

## rm

Remove a task permanently.

```
tsk rm <id>
```

This deletes the task from the JSON file. It cannot be undone.

### Examples

```bash
$ tsk rm 3
task 3 removed
```

### Errors

```bash
# task doesn't exist
$ tsk rm 99
tsk: task 99 not found

# invalid id
$ tsk rm abc
invalid id: abc
```

---

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Error (invalid args, task not found, file I/O failure) |
