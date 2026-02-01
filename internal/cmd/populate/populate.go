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
║                   MANDOR CLI COMMAND REFERENCE                     ║
║              Event-Based Task Manager for AI Workflows             ║
╚════════════════════════════════════════════════════════════════════╝

═════════════════════════════════════════════════════════════════════════
 TABLE OF CONTENTS
═════════════════════════════════════════════════════════════════════════

  1. Workspace Management
  2. Project Management
  3. Feature Management
  4. Task Management
  5. Issue Management
  6. Utility Commands
  7. Configuration Reference
  8. Status Transitions
  9. Dependency Rules
  10. Input Formats
  11. Best Practices
  12. Common Workflows

═════════════════════════════════════════════════════════════════════════
 1. WORKSPACE COMMANDS
═════════════════════════════════════════════════════════════════════════

▶ mandor init <workspace_name> [--yes]
  Initialize a new workspace
  
  Flags:
    --yes, -y         Skip confirmation prompts
  
  Example:
    mandor init "AI Agent Project" -y

  Best Practice:
    - Run once per project to set up .mandor/ directory
    - Use descriptive workspace names
    - Set up config defaults after initialization

───────────────────────────────────────────────────────────────────────

▶ mandor status [--project <id>] [--summary] [--json]
  Display workspace and project status overview
  
  Flags:
    --project, -p     Filter by specific project ID
    --summary, -s     Show compact summary view (counts only)
    --json, -j        Machine-readable JSON output
  
  Examples:
    mandor status                    # Full workspace overview
    mandor status --project api      # Single project details
    mandor status --summary          # Count summary
    mandor status --json             # JSON output

───────────────────────────────────────────────────────────────────────

▶ mandor config [get|set|list|reset] [<key>] [<value>]
  Manage workspace configuration
  
  Subcommands:
    get <key>              Get configuration value
    set <key> <value>      Set configuration value
    list                   List all configuration
    reset <key>            Reset to default value
  
  Configuration Keys:
    priority.default       Default priority for entities (P0-P5)
    strictMode             Enable strict dependency checking (true|false)
  
  Examples:
    mandor config list                      # Show all config
    mandor config set priority.default P2   # Set default priority
    mandor config get priority.default      # Get priority default
    mandor config reset priority.default    # Reset to default

  Best Practice:
    - Set priority.default early for consistency
    - Enable strictMode for production workflows

═════════════════════════════════════════════════════════════════════════
 2. PROJECT COMMANDS
═════════════════════════════════════════════════════════════════════════

▶ mandor project create <project_id> --name <name> --goal <goal> [OPTIONS]
  Create a new project
  
  Required Arguments:
    <project_id>                   Unique project identifier (alphanumeric, hyphens)
  
  Required Flags:
    --name, -n <text>              Project display name
    --goal, -g <text>              Project objective/goal (500+ chars, or 2+ in dev mode)
  
  Optional Flags:
    --task-dep <rule>              Task dependency rule (default: same_project_only)
                                   Values: same_project_only | cross_project_allowed | disabled
    --feature-dep <rule>           Feature dependency rule (default: cross_project_allowed)
                                   Values: same_project_only | cross_project_allowed | disabled
    --issue-dep <rule>             Issue dependency rule (default: same_project_only)
                                   Values: same_project_only | cross_project_allowed | disabled
    --strict                       Enable strict dependency enforcement
    --yes, -y                      Skip confirmation prompts
  
  Example:
    mandor project create api \
      --name "REST API Service" \
      --goal "Build production REST API with authentication, authorization..." \
      --task-dep cross_project_allowed \
      --feature-dep cross_project_allowed

───────────────────────────────────────────────────────────────────────

▶ mandor project list [--json]
  List all projects in workspace
  
  Flags:
    --json, -j         Machine-readable JSON output
  
  Example:
    mandor project list
    mandor project list --json

───────────────────────────────────────────────────────────────────────

▶ mandor project detail <project_id>
  Show detailed project information
  
  Displays:
    - Project metadata (name, goal, priority rules)
    - Feature/task/issue counts
    - Dependency configuration
    - Creation/update timestamps
    - Creator information
  
  Example:
    mandor project detail api

───────────────────────────────────────────────────────────────────────

▶ mandor project update <project_id> [--name <text>] [--goal <text>]
  Update project properties
  
  Optional Flags:
    --name, -n <text>     Update project name
    --goal, -g <text>     Update project goal
  
  Example:
    mandor project update api --goal "Enhanced API with new features..."

───────────────────────────────────────────────────────────────────────

▶ mandor project delete <project_id> [--hard]
  Delete a project (soft delete by default)
  
  Flags:
    --hard             Permanently delete (cannot be restored)
  
  Default: Soft delete (can be reopened with 'reopen')
  
  Example:
    mandor project delete legacy          # Soft delete
    mandor project delete legacy --hard   # Permanent delete

───────────────────────────────────────────────────────────────────────

▶ mandor project reopen <project_id>
  Restore a soft-deleted project
  
  Example:
    mandor project reopen legacy

═════════════════════════════════════════════════════════════════════════
 3. FEATURE COMMANDS
═════════════════════════════════════════════════════════════════════════

▶ mandor feature create <name> --project <id> --goal <goal> [OPTIONS]
  Create a feature (high-level functionality epic)
  
  Required Arguments:
    <name>                         Feature name (positional)
  
  Required Flags:
    --project, -p <id>             Project ID (required)
    --goal, -g <text>              Feature goal/description (300+ chars, or 2+ in dev)
  
  Optional Flags:
    --scope <scope>                Feature scope (frontend|backend|fullstack|cli|desktop|
                                                  android|flutter|react-native|ios|swift)
    --priority <P0-P5>             Priority level (default from config)
    --depends <ids>                Pipe-separated feature IDs for dependencies
    --yes, -y                      Skip confirmation
  
  Example:
    mandor feature create "User Authentication" \
      --project api \
      --goal "Implement login, logout, registration, password reset flows..." \
      --scope backend \
      --priority P1 \
      --depends api-feature-security

  Initial Status: draft (becomes active when marked complete)

───────────────────────────────────────────────────────────────────────

▶ mandor feature list [--project <id>] [--json]
  List features
  
  Flags:
    --project, -p <id>    Filter by project
    --json, -j            JSON output
  
  Example:
    mandor feature list --project api

───────────────────────────────────────────────────────────────────────

▶ mandor feature detail <feature_id> --project <id>
  Show complete feature details
  
  Flags:
    --project, -p <id>    Project ID (required)
  
  Example:
    mandor feature detail api-feature-abc123 --project api

───────────────────────────────────────────────────────────────────────

▶ mandor feature update <feature_id> --project <id> [OPTIONS]
  Update feature properties or status
  
  Required Flags:
    --project, -p <id>    Project ID
  
  Optional Flags:
    --name <text>               Update feature name
    --goal <text>               Update goal description
    --scope <scope>             Update scope
    --priority <P0-P5>          Update priority
    --status <status>           Set status (draft|active|done|blocked|cancelled)
    --depends <ids>             Update dependencies (pipe-separated IDs)
    --cancel --reason <text>    Cancel with reason (audit trail)
    --reopen                    Reopen cancelled feature
    --force                     Force cancel even with dependents
    --dry-run                   Preview changes without saving
  
  Example:
    mandor feature update api-feature-abc123 \
      --project api \
      --priority P0 \
      --status active

───────────────────────────────────────────────────────────────────────

═════════════════════════════════════════════════════════════════════════
 4. TASK COMMANDS
═════════════════════════════════════════════════════════════════════════

▶ mandor task create <name> --feature <id> --goal <goal> [REQUIRED] [OPTIONS]
  Create a task (individual work item)
  
  Required Arguments:
    <name>                         Task name (positional)
  
  Required Flags:
    --feature, -f <id>             Feature ID (format: project-feature-xxx)
    --goal, -g <text>              Task goal (max 500 chars in prod, 2+ in dev)
    --implementation-steps <steps> Pipe-separated implementation steps
    --test-cases <cases>           Pipe-separated test cases to validate
    --derivable-files <files>      Pipe-separated output files to be created
    --library-needs <libs>         Pipe-separated required libraries (use "none" if N/A)
  
  Optional Flags:
    --priority <P0-P5>             Priority level (default from config)
    --depends-on <ids>             Pipe-separated task IDs for dependencies
    --yes, -y                      Skip confirmation
  
  Example:
    mandor task create "Setup Password Hashing" \
      --feature api-feature-auth \
      --goal "Implement bcrypt password hashing utility" \
      --implementation-steps "Install bcrypt|Create password utility|Add salt management|Write security tests" \
      --test-cases "Hash validation|Verify comparison|Salt uniqueness|Performance check" \
      --derivable-files "src/utils/password.go|src/utils/password_test.go" \
      --library-needs "bcrypt|golang-jwt"
  
  Initial Status: pending (transitions to ready when dependencies complete)

───────────────────────────────────────────────────────────────────────

▶ mandor task list [--feature <id>] [--project <id>] [--status <status>] [OPTIONS]
  List tasks with filtering
  
  Flags:
    --feature, -f <id>    Filter by feature
    --project, -p <id>    Filter by project
    --status <status>     Filter by status (pending|ready|in_progress|done|blocked|cancelled)
    --priority <P0-P5>    Filter by priority
    --json, -j            JSON output
  
  Examples:
    mandor task list --project api
    mandor task list --feature api-feature-auth --status ready
    mandor task list --project api --priority P0

───────────────────────────────────────────────────────────────────────

▶ mandor task detail <task_id>
  Show task details with implementation plan
  
  Displays:
    - Task metadata (name, goal, priority)
    - Implementation steps and test cases
    - Derivable files and library needs
    - Status and dependencies
    - Creation/update information
  
  Example:
    mandor task detail api-feature-auth-task-abc123

───────────────────────────────────────────────────────────────────────

▶ mandor task update <task_id> [OPTIONS]
  Update task properties or status
  
  Optional Flags:
    --status <status>               Change status (pending|ready|in_progress|done|blocked|cancelled)
    --priority <P0-P5>              Update priority
    --name <text>                   Update task name
    --goal <text>                   Update task goal
    --implementation-steps <steps>  Update implementation steps (pipe-separated)
    --test-cases <cases>            Update test cases (pipe-separated)
    --derivable-files <files>       Update output files (pipe-separated)
    --library-needs <libs>          Update library requirements
    --depends-on <ids>              Set dependencies (replace all)
    --depends-add <ids>             Add dependencies (additive)
    --depends-remove <ids>          Remove dependencies
    --cancel --reason <text>        Cancel with reason
    --reopen                        Reopen cancelled task
    --dry-run                       Preview without saving
  
  Example:
    mandor task update api-feature-auth-task-abc123 \
      --status in_progress \
      --priority P0

───────────────────────────────────────────────────────────────────────

▶ mandor task ready [--project <id>] [--feature <id>] [--priority <P0-P5>] [OPTIONS]
  List tasks with status='ready' (available to work on)
  
  Flags:
    --project, -p <id>    Filter by project
    --feature, -f <id>    Filter by feature
    --priority <P0-P5>    Filter by priority
    --json, -j            JSON output
  
  Example:
    mandor task ready --project api --priority P0

───────────────────────────────────────────────────────────────────────

▶ mandor task blocked [--project <id>] [--feature <id>] [--priority <P0-P5>] [OPTIONS]
  List tasks with status='blocked' (waiting on dependencies)
  
  Flags:
    --project, -p <id>    Filter by project
    --feature, -f <id>    Filter by feature
    --priority <P0-P5>    Filter by priority
    --json, -j            JSON output
  
  Example:
    mandor task blocked --project api

═════════════════════════════════════════════════════════════════════════
 5. ISSUE COMMANDS
═════════════════════════════════════════════════════════════════════════

▶ mandor issue create <name> --project <id> --type <type> --goal <goal> [REQUIRED] [OPTIONS]
  Create an issue (bug, improvement, debt, security, performance)
  
  Required Arguments:
    <name>                         Issue name (positional)
  
  Required Flags:
    --project, -p <id>             Project ID
    --type, -t <type>              Issue type (bug|improvement|debt|security|performance)
    --goal, -g <text>              Issue goal (200+ chars, or 2+ in dev)
    --affected-files <files>       Pipe-separated affected file paths
    --affected-tests <tests>       Pipe-separated affected test files
    --implementation-steps <steps> Pipe-separated implementation steps
  
  Optional Flags:
    --priority <P0-P5>             Priority level (default: P2)
    --depends-on <ids>             Pipe-separated issue IDs for dependencies
    --library-needs <libs>         Pipe-separated required libraries
    --yes, -y                      Skip confirmation
  
  Example:
    mandor issue create "Fix memory leak in auth handler" \
      --project api \
      --type bug \
      --priority P0 \
      --goal "Goroutine not properly cleaned up in token refresh handler causing memory accumulation" \
      --affected-files "src/handlers/auth.go|src/middleware/auth.go" \
      --affected-tests "src/handlers/auth_test.go" \
      --implementation-steps "Identify leak source|Add defer cleanup|Add tests|Verify with pprof"
  
  Initial Status: open (transitions based on dependencies)

───────────────────────────────────────────────────────────────────────

▶ mandor issue list [--project <id>] [--type <type>] [--status <status>] [OPTIONS]
  List issues with filtering
  
  Flags:
    --project, -p <id>    Filter by project
    --type, -t <type>     Filter by type (bug|improvement|debt|security|performance)
    --status <status>     Filter by status (open|ready|in_progress|resolved|wontfix|cancelled)
    --priority <P0-P5>    Filter by priority
    --json, -j            JSON output
  
  Examples:
    mandor issue list --project api
    mandor issue list --project api --type bug --priority P0
    mandor issue list --status open

───────────────────────────────────────────────────────────────────────

▶ mandor issue detail <issue_id> [--project <id>]
  Show issue details
  
  Optional Flags:
    --project, -p <id>    Project ID (auto-extracted if omitted)
  
  Example:
    mandor issue detail api-issue-abc123 --project api

───────────────────────────────────────────────────────────────────────

▶ mandor issue update <issue_id> [--project <id>] [OPTIONS]
  Update issue properties or status
  
  Optional Flags:
    --project, -p <id>              Project ID (auto-extracted if omitted)
    --name <text>                   Update issue name
    --goal <text>                   Update goal
    --type <type>                   Change issue type
    --priority <P0-P5>              Update priority
    --status <status>               Set status (open|ready|in_progress|resolved|wontfix|cancelled)
    --reason <text>                 Reason for status change
    --depends-on <ids>              Set dependencies (replace)
    --depends-add <ids>             Add dependencies
    --depends-remove <ids>          Remove dependencies
    --affected-files <files>        Update affected files
    --affected-tests <tests>        Update affected tests
    --implementation-steps <steps>  Update implementation steps
    --library-needs <libs>          Update library needs
    --start                         Transition to in_progress
    --resolve                       Mark as resolved
    --wontfix --reason <text>       Mark as wontfix with reason
    --reopen                        Reopen issue
    --cancel --reason <text>        Cancel issue
    --dry-run                       Preview without saving
  
  Examples:
    mandor issue update api-issue-abc123 --resolve
    mandor issue update api-issue-abc123 --start
    mandor issue update api-issue-abc123 --wontfix --reason "Working as intended"

───────────────────────────────────────────────────────────────────────

▶ mandor issue ready [--project <id>] [--type <type>] [--priority <P0-P5>] [OPTIONS]
  List issues with status='ready' (available to fix)
  
  Flags:
    --project, -p <id>    Filter by project
    --type, -t <type>     Filter by type
    --priority <P0-P5>    Filter by priority
    --json, -j            JSON output
  
  Example:
    mandor issue ready --project api --type bug --priority P0

───────────────────────────────────────────────────────────────────────

▶ mandor issue blocked [--project <id>] [--type <type>] [--priority <P0-P5>] [OPTIONS]
  List issues with status='blocked' (waiting on dependencies)
  
  Flags:
    --project, -p <id>    Filter by project
    --type, -t <type>     Filter by type
    --priority <P0-P5>    Filter by priority
    --json, -j            JSON output
  
  Example:
    mandor issue blocked --project api --type security

═════════════════════════════════════════════════════════════════════════
 6. UTILITY COMMANDS
═════════════════════════════════════════════════════════════════════════

▶ mandor completion [bash|zsh|fish]
  Generate shell completion scripts
  
  Supported Shells: bash, zsh, fish
  
  Setup:
    eval "$(mandor completion bash)"           # Bash
    eval "$(mandor completion zsh)"            # Zsh
    mandor completion fish | source            # Fish
  
  Permanent Setup:
    mandor completion bash > ~/.bash_completion.d/mandor
    mandor completion zsh > ~/.zsh/completions/_mandor

───────────────────────────────────────────────────────────────────────

▶ mandor populate [--markdown|--json]
  Display this command reference
  
  Flags:
    --markdown, -m    Output in Markdown format
    --json, -j        Output in JSON format

───────────────────────────────────────────────────────────────────────

▶ mandor version
  Display version information

═════════════════════════════════════════════════════════════════════════
 7. CONFIGURATION REFERENCE
═════════════════════════════════════════════════════════════════════════

Configuration Directory: .mandor/config.json

Available Keys:
  priority.default      Default priority for new entities
                        Valid values: P0, P1, P2, P3, P4, P5
                        Default: P3

  strictMode            Enable strict dependency validation
                        Valid values: true, false
                        Default: false

Set Configuration:
  $ mandor config set priority.default P2
  $ mandor config set strictMode true

Get Configuration:
  $ mandor config get priority.default

List All:
  $ mandor config list

Reset to Default:
  $ mandor config reset priority.default

═════════════════════════════════════════════════════════════════════════
 8. STATUS TRANSITIONS
═════════════════════════════════════════════════════════════════════════

FEATURE STATUS FLOW:
  draft ─→ active ─→ done
    │
    └─→ cancelled (with reason)

  - draft:     Initial state, can add tasks
  - active:    Work in progress
  - done:      Completed
  - blocked:   Waiting on dependencies
  - cancelled: Abandoned (can reopen)

TASK STATUS FLOW:
  pending ─→ ready ─→ in_progress ─→ done
    │         │         │
    │         └─→ blocked ─→ ready
    │
    └─→ cancelled (with reason)

  - pending:      Waiting on dependencies
  - ready:        Available to start (auto-transitioned when deps complete)
  - in_progress:  Currently being worked on
  - done:         Completed (auto-unblocks dependents)
  - blocked:      Manually blocked on external factors
  - cancelled:    Abandoned (can reopen)

ISSUE STATUS FLOW:
  open ─→ ready ─→ in_progress ─→ resolved
   │       │         │
   │       └─→ blocked ─→ ready
   │
   └─→ wontfix (with reason)
   └─→ cancelled (with reason)

  - open:         Newly reported
  - ready:        Available to work on (deps satisfied)
  - in_progress:  Currently being fixed
  - resolved:     Fixed and verified
  - blocked:      Waiting on dependencies
  - wontfix:      Intentionally not fixing (reason recorded)
  - cancelled:    Duplicate or no longer relevant

═════════════════════════════════════════════════════════════════════════
 9. DEPENDENCY RULES
═════════════════════════════════════════════════════════════════════════

Each project defines rules for each entity type:

SAME_PROJECT_ONLY:
  ✓ Feature A → Feature B (same project)
  ✗ Feature A → Feature B (different projects)
  Auto-transitioned when all dependencies complete or are cancelled.

CROSS_PROJECT_ALLOWED:
  ✓ Feature A → Feature B (same project)
  ✓ Feature A → Feature B (different projects)
  Allows cross-project dependency chains.

DISABLED:
  ✗ No dependencies allowed
  Cannot set depends-on flag.

Default Rules:
  - Features:  cross_project_allowed
  - Tasks:     same_project_only
  - Issues:    same_project_only

═════════════════════════════════════════════════════════════════════════
 10. INPUT FORMATS
═════════════════════════════════════════════════════════════════════════

PIPE-SEPARATED LISTS:
  Used for: implementation-steps, test-cases, derivable-files, library-needs,
            affected-files, affected-tests, depends-on
  
  Format: "item1|item2|item3|..."
  
  Example:
    --implementation-steps "Step 1|Step 2|Step 3"
    --test-cases "Test A|Test B|Test C"
    --library-needs "lodash|axios|express"
  
  Spaces: Trimmed automatically (leading/trailing)
  Empty strings: Filtered out

MULTIPLE IDS (for dependencies):
  Format: "id1|id2|id3|..."
  
  Example:
    --depends "feature-1|feature-2"
    --depends-on "task-1|task-2|task-3"

GOAL TEXT:
  Character limits (in production):
    - Feature: 300+ characters
    - Task: 500 characters max
    - Issue: 200+ characters
    - Project: 500+ characters
  
  Character limits (in development mode - MANDOR_ENV=development):
    - Feature: 2+ characters
    - Task: 500 characters max
    - Issue: 2+ characters
    - Project: 2+ characters
  
  Development mode enabled via: export MANDOR_ENV=development

PRIORITY LEVELS:
  Values: P0, P1, P2, P3, P4, P5
  
  P0 = Critical   (must do immediately)
  P1 = High       (important, soon)
  P2 = Medium     (should do)
  P3 = Normal     (default priority)
  P4 = Low        (nice to have)
  P5 = Minimal    (can defer)

═════════════════════════════════════════════════════════════════════════
 11. BEST PRACTICES
═════════════════════════════════════════════════════════════════════════

WORKFLOW DESIGN:
  1. Create workspace with mandor init
  2. Configure defaults (priority, strict mode)
  3. Define projects with clear goals
  4. Break projects into features
  5. Decompose features into tasks
  6. Report and track issues
  7. Use dependencies to enforce order
  8. Monitor ready/blocked queues

FEATURE CREATION:
  ✓ Use clear, descriptive names
  ✓ Define scope (frontend/backend/etc)
  ✓ Set appropriate priority (P0 for critical)
  ✓ Add dependencies for ordering
  ✓ Write detailed goals (300+ chars)
  ✗ Create features for everything (use for epics only)
  ✗ Set priority P5 for critical work

TASK CREATION:
  ✓ Define implementation steps FIRST
  ✓ Write test cases BEFORE implementation (TDD)
  ✓ List all derivable files
  ✓ Keep completable in one session
  ✓ Add dependencies for blocking
  ✗ Skip test cases
  ✗ Create tasks larger than 1 session
  ✗ Leave library-needs empty (use "none" if N/A)

ISSUE TRACKING:
  ✓ Categorize by type (bug/improvement/debt/etc)
  ✓ List affected files for debugging
  ✓ Include affected tests
  ✓ Set P0 for critical bugs
  ✓ Use --wontfix with reason for rejections
  ✗ Create issues without type
  ✗ Forget affected-files
  ✗ Cancel without reason

DEPENDENCY MANAGEMENT:
  ✓ Use dependencies to enforce order
  ✓ Create linear chains (A → B → C)
  ✓ Use fan-in for convergence (M depends on A, B, C)
  ✓ Use fan-out for one unblocks many
  ✓ Check for cycles (system prevents them)
  ✗ Create circular dependencies (prevented by system)
  ✗ Set self-dependencies (prevented)

STATUS MANAGEMENT:
  ✓ Use auto-transitions (ready when deps complete)
  ✓ Mark done/resolved to unblock dependents
  ✓ Use --cancel --reason for audit trail
  ✓ Use --reopen for accidental cancellations
  ✓ Review blocked queue regularly
  ✗ Skip status transitions
  ✗ Cancel without reason
  ✗ Leave entities stuck in blocked

NAMING CONVENTIONS:
  Features:   "User Authentication", "Admin Dashboard", "Payment Processing"
  Tasks:      "Setup bcrypt hashing", "Write auth tests", "Create JWT middleware"
  Issues:     "Fix memory leak in auth", "Improve error messages", "Refactor database layer"
  Projects:   "api", "frontend", "infrastructure", "mobile-app"

═════════════════════════════════════════════════════════════════════════
 12. COMMON WORKFLOWS
═════════════════════════════════════════════════════════════════════════

WORKFLOW 1: Start a New Project
  1. mandor init "My Project" -y
  2. mandor config set priority.default P3
  3. mandor project create api --name "REST API" --goal "..."
  4. mandor feature create "Auth" --project api --goal "..."
  5. mandor task create "Setup" --feature api-feature-xxx --goal "..."

WORKFLOW 2: Find and Start Work
  1. mandor task ready --project api                # Find ready tasks
  2. mandor task detail <task-id>                   # Review details
  3. mandor task update <task-id> --status in_progress   # Mark started
  4. (do work)
  5. mandor task update <task-id> --status done     # Mark complete

WORKFLOW 3: Fix a Bug
  1. mandor issue create "Bug name" --project api --type bug --goal "..."
  2. mandor issue detail <issue-id>                 # Review
  3. mandor issue update <issue-id> --start         # Start fixing
  4. (fix code, run tests)
  5. mandor issue update <issue-id> --resolve       # Mark resolved

WORKFLOW 4: Handle Blocked Task
  1. mandor task blocked --project api              # Find blocked
  2. mandor task detail <blocked-task>              # See what's blocking
  3. (complete blocking task first)
  4. System auto-transitions blocked task to ready
  5. mandor task ready --project api                # Verify transition

WORKFLOW 5: Cancel and Reopen
  1. mandor feature update <id> --project api --cancel --reason "Postponed for now"
  2. (later...)
  3. mandor feature update <id> --project api --reopen
  4. Feature returns to draft status

WORKFLOW 6: Track Dependencies
  1. Create independent tasks (no dependencies) - start as ready
  2. Create tasks with dependencies - start as blocked
  3. Complete independent task: dependent tasks auto-transition to ready
  4. Use mandor task blocked to monitor unblock progress

═════════════════════════════════════════════════════════════════════════
 QUICK REFERENCE
═════════════════════════════════════════════════════════════════════════

Start Here:
  mandor init "My Project"
  mandor config set priority.default P3

Basic Work:
  mandor project create <id> --name X --goal "..."
  mandor feature create <name> --project <id> --goal "..."
  mandor task create <name> --feature <id> --goal "..." --implementation-steps "..." --test-cases "..." --derivable-files "..." --library-needs "..."
  mandor task update <id> --status in_progress
  mandor task update <id> --status done

Monitor Progress:
  mandor status
  mandor task ready --project <id>
  mandor task blocked --project <id>
  mandor issue ready --project <id> --type bug

Exit Codes:
  0 = Success
  1 = System error
  2 = Validation error
  3 = Permission error

═════════════════════════════════════════════════════════════════════════

For detailed help on any command, use:
  mandor <command> --help
  mandor <command> <subcommand> --help

Example:
  mandor task --help
  mandor task create --help

═════════════════════════════════════════════════════════════════════════
`)
	return nil
}
