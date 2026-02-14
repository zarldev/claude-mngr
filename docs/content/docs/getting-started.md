---
title: "Getting Started"
weight: 1
---

# Getting Started

## Installation

### Using Go Install

If you have Go installed (1.23 or later):

```bash
go install github.com/zarldev/claude-mngr/cmd/tsk@latest
```

Make sure `$GOPATH/bin` (usually `~/go/bin`) is in your `PATH`.

### From Source

```bash
git clone https://github.com/zarldev/claude-mngr.git
cd claude-mngr
go build -o tsk ./cmd/tsk
mv tsk /usr/local/bin/  # or anywhere in your PATH
```

## Your First Task

Add a task:

```bash
$ tsk add Buy groceries
added task 1: Buy groceries
```

That's it. The task is stored in `~/.tasks.json` and will persist across sessions.

## Basic Workflow

A typical workflow looks like this:

```bash
# start the day — add what you need to do
$ tsk add Review PR #42
added task 2: Review PR #42

$ tsk add Fix login timeout bug
added task 3: Fix login timeout bug

$ tsk add Update deployment docs
added task 4: Update deployment docs

# check your list
$ tsk list
  2  [ ]  Review PR #42                             5m ago
  3  [ ]  Fix login timeout bug                     3m ago
  4  [ ]  Update deployment docs                    1m ago

# finish a task
$ tsk done 2
task 2 marked done

# see what's left
$ tsk list --pending
  3  [ ]  Fix login timeout bug                     10m ago
  4  [ ]  Update deployment docs                    8m ago

# check what you've completed
$ tsk list --done
  2  [x]  Review PR #42                             12m ago

# remove a task that's no longer relevant
$ tsk rm 4
task 4 removed
```

## Next Steps

- See the [Commands Reference]({{< relref "/docs/commands" >}}) for all available commands and flags
- Read about [Configuration]({{< relref "/docs/configuration" >}}) to understand where your data lives
