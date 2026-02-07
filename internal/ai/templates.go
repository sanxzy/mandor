package ai

func generateClaudeTemplate(projectName, createdAt string) string {
	return "# CLAUDE.md\n\n" +
		"This file provides guidance to Claude Code when working with this project.\n\n" +
		"## Project Context\n\n" +
		"- **Project**: " + projectName + "\n" +
		"- **Created**: " + createdAt + "\n" +
		"- **Type**: Mandor workspace\n\n" +
		"## Mandor Commands\n\n" +
		"### Essential Commands\n\n" +
		"1. **mandor populate** - View all available commands\n" +
		"   ```bash\n" +
		"   mandor populate\n" +
		"   ```\n\n" +
		"2. **mandor track** - Check status of work items\n" +
		"   ```bash\n" +
		"   mandor track              # View entire workspace\n" +
		"   mandor track backlog <id>  # View specific backlog\n" +
		"   mandor track feature <id> # View feature with tasks\n" +
		"   mandor track task <id>    # View single task\n" +
		"   ```\n\n" +
		"3. **mandor session note** - Record session progress\n" +
		"   ```bash\n" +
		"   mandor session note \"Starting work on...\"\n" +
		"   mandor session note \"Completed feature X\"\n" +
		"   mandor session note --read\n" +
		"   ```\n\n" +
		"### Backlog Management\n\n" +
		"```bash\n" +
		"# Create backlog\n" +
		"mandor backlog create <backlog_id> --name \"Backlog Name\"\n\n" +
		"# View backlog details\n" +
		"mandor backlog detail <backlog_id>\n\n" +
		"# Delete backlog\n" +
		"mandor backlog delete <backlog_id>\n" +
		"```\n\n" +
		"### Feature Management\n\n" +
		"```bash\n" +
		"# Create feature\n" +
		"mandor feature create --backlog <id> --capability <id> --spec <id> --name \"Feature Name\"\n\n" +
		"# Update feature\n" +
		"mandor feature update <id> --name \"New Name\"\n\n" +
		"# Feature status: draft -> active -> done\n" +
		"mandor feature start <id>\n" +
		"mandor feature done <id>\n" +
		"```\n\n" +
		"### Task Management\n\n" +
		"```bash\n" +
		"# Create task\n" +
		"mandor task create --feature <id> --spec <id> --name \"Task Name\"\n\n" +
		"# Update task\n" +
		"mandor task update <id> --name \"New Name\"\n\n" +
		"# Task status: pending -> ready -> in_progress -> done\n" +
		"mandor task start <id>\n" +
		"mandor task done <id>\n" +
		"```\n\n" +
		"### Issue Tracking\n\n" +
		"```bash\n" +
		"# Create issue\n" +
		"mandor issue create --backlog <id> --name \"Issue Name\" --type bug|improvement|debt|security|performance\n\n" +
		"# Update issue\n" +
		"mandor issue update <id> --start|resolve|wontfix|cancel\n" +
		"```\n\n" +
		"## Workflow\n\n" +
		"1. Start session: mandor session note \"Starting work on...\"\n" +
		"2. Check what's ready: mandor track\n" +
		"3. Create work: mandor feature create or mandor task create\n" +
		"4. Track progress: mandor track feature <id>\n" +
		"5. End session: mandor session note \"Completed...\"\n"
}

func generateAgentsTemplate(projectName, createdAt string) string {
	return "# Mandor Essential Commands\n\n" +
		"Use these three commands to manage your work:\n\n" +
		"---\n\n" +
		"## 1. mandor populate\n\n" +
		"View all available commands and usage instructions.\n\n" +
		"```bash\n" +
		"mandor populate\n" +
		"```\n\n" +
		"Shows complete reference guide with examples for creating backlogs, features, tasks, issues, and managing dependencies.\n\n" +
		"---\n\n" +
		"## 2. mandor track\n\n" +
		"Check status of workspace, backlogs, features, and tasks.\n\n" +
		"```bash\n" +
		"# View entire workspace\n" +
		"mandor track\n\n" +
		"# View specific backlog\n" +
		"mandor track backlog <backlog-id>\n\n" +
		"# View feature with tasks\n" +
		"mandor track feature <feature-id>\n\n" +
		"# View single task\n" +
		"mandor track task <task-id>\n" +
		"```\n\n" +
		"Shows what's ready, blocked, in progress, or done. Use before starting work.\n\n" +
		"---\n\n" +
		"## 3. mandor session note\n\n" +
		"Record and read session progress (for AI agents).\n\n" +
		"```bash\n" +
		"# Log what you completed\n" +
		"mandor session note \"Completed v0.6.0 feature implementation\"\n\n" +
		"# Read last 50 notes\n" +
		"mandor session note --read\n\n" +
		"# Read more notes\n" +
		"mandor session note --read --offset 100\n" +
		"```\n\n" +
		"End each session with a note. Start next session by resuming from notes.\n\n" +
		"---\n\n" +
		"## Quick Workflow\n\n" +
		"1. mandor track - See what is ready to work on\n" +
		"2. mandor populate - Learn how to create/update work\n" +
		"3. mandor session note \"done\" - Log progress before ending session\n" +
		"4. mandor session note --read - Check progress when resuming\n\n" +
		"---\n\n" +
		"## Project Context\n\n" +
		"- **Project**: " + projectName + "\n" +
		"- **Created**: " + createdAt + "\n"
}
