# Mandor Essential Commands

Use these three commands to manage your work:

---

## 1. mandor populate

View all available commands and usage instructions.

```bash
mandor populate
```

Shows complete reference guide with examples for creating projects, features, tasks, issues, and managing dependencies.

---

## 2. mandor track

Check status of workspace, projects, features, and tasks.

```bash
# View entire workspace
mandor track

# View specific project
mandor track project <project-id>

# View feature with tasks
mandor track feature <feature-id>

# View single task
mandor track task <task-id>
```

Shows what's ready, blocked, in progress, or done. Use before starting work.

---

## 3. mandor session note

Record and read session progress (for AI agents).

```bash
# Log what you completed
mandor session note "Completed v0.4.4 release and testing"

# Read last 50 notes
mandor session note --read

# Read more notes
mandor session note --read --offset 100
```

End each session with a note. Start next session by reading notes to resume work.

---

## Quick Workflow

1. `mandor track` - See what's ready
2. `mandor populate` - Learn how to create/update work
3. `mandor session note "done"` - Log progress before ending session
4. `mandor session note --read` - Check progress when resuming
