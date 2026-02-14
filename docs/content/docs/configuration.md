---
title: "Configuration"
weight: 3
---

# Configuration

`tsk` requires no configuration. It works out of the box with sensible defaults.

## Data File

Tasks are stored in a single JSON file:

```
~/.tasks.json
```

This file is created automatically the first time you add a task. If the file doesn't exist or is empty, `tsk` starts fresh with no tasks.

## File Format

The JSON file has this structure:

```json
{
  "next_id": 4,
  "tasks": [
    {
      "id": 1,
      "title": "Buy groceries",
      "done": false,
      "created_at": "2026-02-14T10:30:00Z"
    },
    {
      "id": 2,
      "title": "Fix login timeout bug",
      "done": true,
      "created_at": "2026-02-13T09:15:00Z"
    },
    {
      "id": 3,
      "title": "Update deployment docs",
      "done": false,
      "created_at": "2026-02-14T14:00:00Z"
    }
  ]
}
```

### Fields

| Field | Type | Description |
|-------|------|-------------|
| `next_id` | integer | The ID that will be assigned to the next task created |
| `tasks` | array | List of all tasks |
| `tasks[].id` | integer | Unique task identifier |
| `tasks[].title` | string | Task description |
| `tasks[].done` | boolean | Whether the task is completed |
| `tasks[].created_at` | string | RFC 3339 timestamp of when the task was created |

## Syncing

Because the data is a single JSON file, you can sync it across machines using any file synchronization tool:

- Symlink to a Dropbox/iCloud/Google Drive folder
- Include in a dotfiles repo
- Use `rsync` or similar

### Example: Symlink to iCloud

```bash
mv ~/.tasks.json ~/Library/Mobile\ Documents/com~apple~CloudDocs/tasks.json
ln -s ~/Library/Mobile\ Documents/com~apple~CloudDocs/tasks.json ~/.tasks.json
```

## Backup

The file is human-readable JSON with pretty printing (2-space indent), so you can easily inspect or edit it manually if needed.

```bash
# quick backup
cp ~/.tasks.json ~/.tasks.json.bak
```

## Limitations

- No custom file path via CLI flag or environment variable (always `~/.tasks.json`)
- No multi-user or locking support — designed for single-user use
- Task IDs are never reused within a session (the `next_id` counter always increments)
