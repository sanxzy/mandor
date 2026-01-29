package populate

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewPopulateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "populate",
		Short: "Display all commands, options, and best practices",
		Long: `Display comprehensive documentation of all available commands,
their options, flags, and best practices for effective use.

This command serves as a quick reference guide for learning Mandor
and understanding the recommended workflows.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return outputPopulate(cmd)
		},
	}

	cmd.Flags().BoolP("markdown", "m", false, "Output in Markdown format")
	cmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	return cmd
}

func outputPopulate(cmd *cobra.Command) error {
	out := cmd.OutOrStdout()

	fmt.Fprint(out, `
╔════════════════════════════════════════════════════════════════════╗
║                    MANDOR COMMAND REFERENCE                       ║
╚════════════════════════════════════════════════════════════════════╝

═════════════════════════════════════════════════════════════════════
 WORKSPACE COMMANDS
═════════════════════════════════════════════════════════════════════

▶ mandor init <name>
  Initialize a new workspace
  Best Practice: Choose a descriptive name that reflects your project scope
  Example: mandor init "AI Agent Project"

▶ mandor status [--project <id>] [--summary] [--json]
  Display workspace and project status
  Best Practice: Run regularly to track progress
  --project, -p  : Filter by specific project
  --summary, -s  : Show compact summary view
  --json, -j     : Machine-readable JSON output

▶ mandor config [get|set|list|reset] [<key>] [<value>]
  Manage workspace configuration
  Best Practice: Set defaults early for consistent behavior
  Examples:
    mandor config list
    mandor config set priority.default P2
    mandor config get priority.default

═════════════════════════════════════════════════════════════════════
 PROJECT COMMANDS
═════════════════════════════════════════════════════════════════════

▶ mandor project create <id> --name <name> --goal <goal> [--strict]
  Create a new project
  Best Practice: Define clear, measurable goals (500+ chars)
  Example:
    mandor project create api \
      --name "API Service" \
      --goal "Implement REST API with authentication..."

▶ mandor project list
  List all projects
  Best Practice: Review project structure before creating entities

▶ mandor project detail <id>
  Show detailed project information
  Best Practice: Check dependencies and entity counts

▶ mandor project update <id> [--name <name>] [--goal <goal>] [--status <status>]
  Update project metadata
  Best Practice: Update goals when project scope changes

▶ mandor project delete <id> [--hard]
  Delete a project (soft delete by default)
  Best Practice: Use soft delete first; hard delete is permanent

▶ mandor project reopen <id>
  Reopen a deleted project
  Best Practice: Restore accidentally deleted projects

═════════════════════════════════════════════════════════════════════
 FEATURE COMMANDS
═════════════════════════════════════════════════════════════════════

▶ mandor feature create <name> --project <id> --goal <goal> \
    [--scope <scope>] [--priority <P0-P5>] [--depends <ids>]
  Create a new feature (high-level functionality)
  Best Practices:
    - Keep features as epics that can be broken into tasks
    - Define clear scope (frontend/backend/fullstack/cli)
    - Set appropriate priority (P0=critical, P5=minimal)
    - Add dependencies to enforce completion order
  Example:
    mandor feature create "User Authentication" \
      --project api \
      --goal "Implement login, logout, registration..." \
      --scope backend \
      --priority P2 \
      --depends api-feature-other

▶ mandor feature list [--project <id>] [--json]
  List features with status and priority
  Best Practice: Filter by project for focused view

▶ mandor feature detail <id> --project <id>
  Show complete feature details
  Best Practice: Review before creating dependent features

▶ mandor feature update <id> --project <id> \
    [--name <name>] [--priority <P0-P5>] [--status <status>] \
    [--cancel --reason <text>] [--reopen]
  Update feature metadata or status
  Best Practices:
    - Update priority as project evolves
    - Cancel with clear reason for audit trail
    - Reopen if cancelled prematurely

═════════════════════════════════════════════════════════════════════
 TASK COMMANDS
═════════════════════════════════════════════════════════════════════

▶ mandor task create <name> --feature <id> --goal <goal> \
    --implementation-steps <steps> --test-cases <cases> \
    --derivable-files <files> --library-needs <libs> \
    [--priority <P0-P5>] [--depends-on <ids>]
  Create a new task (individual work item)
  Best Practices:
    - Each task should be completable in one session
    - Define test cases BEFORE implementation (TDD)
    - List all derivable files for clear output
    - Specify libraries needed for dependency management
  Example:
    mandor task create "Password Hashing" \
      --feature api-feature-abc \
      --goal "Implement bcrypt password hashing" \
      --implementation-steps "Install bcrypt,Create utility,Write tests" \
      --test-cases "Hash validation,Password comparison" \
      --derivable-files "src/utils/password.ts,src/types/auth.ts" \
      --library-needs "bcrypt"

▶ mandor task list [--feature <id>] [--project <id>] \
    [--status <status>] [--priority <P0-P5>] [--json]
  List tasks with filters
  Best Practice: Use combined filters for precise views

▶ mandor task detail <id>
  Show task details with implementation progress
  Best Practice: Review before starting work

▶ mandor task update <id> \
     [--status <pending|ready|in_progress|done|blocked|cancelled>] \
     [--priority <P0-P5>] [--name <name>] \
     [--cancel --reason <text>] [--reopen] \
     [--depends-add <ids>] [--depends-remove <ids>]
   Update task status and metadata
   Best Practices:
     - Set status to 'in_progress' when starting
     - Add dependencies to block until prerequisites complete
     - Cancel with reason for audit trail
     - Use 'blocked' when waiting on external factors

▶ mandor task ready [--project <id>] [--feature <id>] [--priority <P0-P5>] [--json]
   List tasks with status='ready' (available to work on)
   Best Practice: Check ready tasks to find work
   Example:
     mandor task ready --project api --priority P0

▶ mandor task blocked [--project <id>] [--feature <id>] [--priority <P0-P5>] [--json]
   List tasks with status='blocked' (waiting on dependencies)
   Best Practice: Review blocked tasks to unblock progress
   Example:
     mandor task blocked --feature api-feature-abc123

═════════════════════════════════════════════════════════════════════
 ISSUE COMMANDS
═════════════════════════════════════════════════════════════════════

▶ mandor issue create <name> --project <id> --type <type> --goal <goal> \
    --affected-files <files> --affected-tests <tests> \
    --implementation-steps <steps> [--library-needs <libs>]
  Create a new issue (bug, improvement, debt, security, performance)
  Best Practices:
    - Categorize correctly (bug/improvement/debt/security/performance)
    - List affected files for focused debugging
    - Include affected tests to verify fixes
  Example:
    mandor issue create "Fix memory leak" \
      --project api \
      --type bug \
      --goal "Fix memory leak in auth handler" \
      --affected-files "src/api/handlers/auth.go" \
      --affected-tests "src/api/handlers/auth_test.go" \
      --implementation-steps "Identify leak,Add cleanup"

▶ mandor issue list [--project <id>] [--type <type>] [--status <status>]
  List issues filtered by criteria
  Best Practice: Regular review of open issues

▶ mandor issue detail <id>
  Show issue details with affected components
  Best Practice: Understand scope before fixing

▶ mandor issue update <id> \
     [--status <open|ready|in_progress|resolved>] \
     [--type <type>] [--priority <P0-P5>] \
     [--resolve] [--wontfix --reason <text>] \
     [--cancel --reason <text>] [--reopen]
   Update issue status
   Best Practices:
     - Mark resolved when fix is complete
     - Use 'wontfix' with reason for declining
     - Cancel duplicates with clear reference

▶ mandor issue ready [--project <id>] [--type <type>] [--priority <P0-P5>] [--json]
   List issues with status='ready' (available to fix)
   Best Practice: Check ready issues to prioritize fixes
   Example:
     mandor issue ready --project api --type bug --priority P0

▶ mandor issue blocked [--project <id>] [--type <type>] [--priority <P0-P5>] [--json]
   List issues with status='blocked' (waiting on dependencies)
   Best Practice: Review blocked issues to unblock fixes
   Example:
     mandor issue blocked --project api --type security

═════════════════════════════════════════════════════════════════════
 UTILITY COMMANDS
═════════════════════════════════════════════════════════════════════

▶ mandor completion [bash|zsh|fish]
  Generate shell completion scripts
  Best Practice: Enable for faster CLI usage

═════════════════════════════════════════════════════════════════════
 BEST PRACTICES SUMMARY
═════════════════════════════════════════════════════════════════════

1. START WITH WORKSPACE
   └─ mandor init "Project Name"
   └─ Set config defaults (priority.default, strictMode)

2. DEFINE PROJECTS FIRST
   └─ Each project = distinct component/service
   └─ Set clear goals (500+ chars for clarity)

3. BREAK INTO FEATURES
   └─ Features = major functionality (epics)
   └─ Set scope and priority appropriately
   └─ Define dependencies early

4. SPLIT INTO TASKS
   └─ Tasks = individual work items
   └─ Use TDD: define test cases first
   └─ Keep tasks completable in one session

5. TRACK ISSUES
   └─ Report bugs, improvements, technical debt
   └─ Link to affected files/tests

6. USE STATUS FLOWS
   Feature: draft → active → done (or blocked/cancelled)
   Task:    pending → ready → in_progress → done (or blocked/cancelled)
   Issue:   open → ready → in_progress → resolved (or wontfix/cancelled)

7. LEVERAGE DEPENDENCIES
   └─ Block until prerequisites complete
   └─ Prevent invalid states automatically
   └─ Enforce logical completion order

8. MAINTAIN AUDIT TRAIL
   └─ Use --cancel --reason "..."
   └─ Use --wontfix --reason "..."
   └─ All changes logged in events.jsonl

9. REVIEW REGULARLY
    └─ mandor status (workspace overview)
    └─ mandor project list (project structure)
    └─ mandor task list --project <id> (task progress)
    └─ mandor task ready (find work available now)
    └─ mandor task blocked (unblock dependencies)
    └─ mandor issue ready (prioritize fixes)
    └─ mandor issue blocked (unblock issues)

═════════════════════════════════════════════════════════════════════
 EXIT CODES
═════════════════════════════════════════════════════════════════════

  0 = Success
  1 = System error (file I/O, permissions)
  2 = Validation error (not found, invalid input)
  3 = Permission error (file/directory access)

═════════════════════════════════════════════════════════════════════
 PRIORITY LEVELS
═════════════════════════════════════════════════════════════════════

  P0 = Critical - Must do immediately
  P1 = High     - Important, soon
  P2 = Medium   - Should do
  P3 = Normal   - Default priority
  P4 = Low      - Nice to have
  P5 = Minimal  - Can defer

═════════════════════════════════════════════════════════════════════
`)
	return nil
}
