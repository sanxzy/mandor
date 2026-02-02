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
║                   Stop Writing Markdown Plans                      ║
╚════════════════════════════════════════════════════════════════════╝

THE PROBLEM WITH MARKDOWN PLANS
════════════════════════════════════════════════════════════════════

Traditional workflows scatter task state across files:
  • Markdown plan files (PLAN.md, TASKS.md) go stale
  • Dependencies are manual and error-prone
  • Status is fiction until code review
  • No audit trail of changes
  • Blocking tasks are invisible
  • Parsing and automation is regex hell

MANDOR'S SOLUTION
════════════════════════════════════════════════════════════════════

Single source of truth with deterministic state:
  ✓ Event-sourced: All changes in immutable events.jsonl
  ✓ Automatic dependencies: Mark task done → dependents auto-unblock
  ✓ Real-time queries: mandor task ready, blocked
  ✓ Audit trail: Full history of every status change
  ✓ Schema-driven: Enforce implementation steps and tests upfront
  ✓ CLI-native: Works in terminals, scripts, and CI/CD

EXAMPLE: Replace This...
════════════════════════════════════════════════════════════════════

  # PLAN.md
  ## Phase 1: Auth System
  - [ ] JWT Parser (depends on crypto)
  - [ ] Login Endpoint (depends on JWT Parser)
  - [ ] Refresh Token (depends on JWT Parser)
  
  Status: Last updated 3 days ago (probably stale!)

...WITH THIS
════════════════════════════════════════════════════════════════════

  mandor feature create "Auth System" --project api \
    --goal "Implement JWT and login endpoints"
  
  mandor task create "JWT Parser" --feature api-feature-xxx \
    --goal "..." --implementation-steps "..." --test-cases "..." \
    --derivable-files "..." --library-needs "..." --priority P1
  
  mandor task create "Login" --feature api-feature-xxx \
    --goal "..." --depends-on jwt-task-id --priority P1
  
  # Real-time progress
  mandor task ready --feature api-feature-xxx          # See what's available NOW
  mandor task blocked --feature api-feature-xxx        # See what's waiting

═════════════════════════════════════════════════════════════════════════
 TABLE OF CONTENTS
═════════════════════════════════════════════════════════════════════════

  1. Workspace Management
  2. Project Management
  3. Feature Management
  4. Task Management
  5. Issue Management
  6. AI Commands
  7. Utility Commands
  8. Configuration Reference
  9. Status Transitions
  10. Dependency Rules
  11. Best Practices
  12. Common Workflows

═════════════════════════════════════════════════════════════════════════
 1. WORKSPACE COMMANDS
═════════════════════════════════════════════════════════════════════════

▶ mandor init [--workspace-name <name>] [-y]
  Initialize a new Mandor workspace in the current directory
  
  Creates .mandor/ directory with workspace metadata and project storage.
  
  Flags:
    --workspace-name <text>   Custom workspace name (default: current directory)
    -y, --yes                 Skip confirmation prompts
    --strict                  Enforce strict dependency rules (deprecated)
  
  Example:
    mandor init "My Project" -y

───────────────────────────────────────────────────────────────────────

▶ mandor status [--project <id>] [--summary] [--json]
  Display the current status of the workspace and all projects
  
  Shows comprehensive statistics including entity counts, status breakdown,
  priority distribution, dependency information, and timeline metrics.
  
  Flags:
    --project, -p <id>    Show status for specific project only
    --summary, -s         Show summary view (one line per entity type)
    --json, -j            Output in JSON format (machine-readable)
    --include-deleted     Include deleted projects in output
  
  Examples:
    mandor status                    # Full workspace overview
    mandor status --project api      # Single project details
    mandor status --summary          # Count summary
    mandor status --json             # JSON output

───────────────────────────────────────────────────────────────────────

▶ mandor summary [--project <id>]
  Display a summary of all features grouped by priority with task counts
  
  Shows features with task counts and status overview.
  
  Flags:
    --project, -p <id>   Filter by project ID
  
  Examples:
    mandor summary                  # All features
    mandor summary --project api    # Project features only

───────────────────────────────────────────────────────────────────────

▶ mandor config [get|set|list|reset] [<key>] [<value>]
  View and modify workspace configuration settings
  
  Available keys:
    - default_priority: Default priority for new entities (P0-P5, default: P3)
    - strict_mode: Enforce strict validation rules (true/false, default: false)
  
  Subcommands:
    get <key>              Display configuration value(s)
    set <key> <value>      Set configuration value
    list                   List all configuration keys with descriptions
    reset <key>            Reset configuration to defaults
  
  Examples:
    mandor config list                      # Show all config
    mandor config set default_priority P2   # Set default priority
    mandor config get default_priority      # Get priority default
    mandor config reset default_priority    # Reset to default

═════════════════════════════════════════════════════════════════════════
 2. PROJECT COMMANDS
═════════════════════════════════════════════════════════════════════════

▶ mandor project create <id> --name <name> --goal <goal> [OPTIONS]
  Create a new project in the workspace
  
  Required Arguments:
    <id>                         Unique project identifier (alphanumeric, hyphens)
  
  Required Flags:
    --name, -n <text>            Project display name
    --goal, -g <text>            Project goal/objectives (min 500 characters)
  
  Optional Flags:
    --task-dep <rule>            Task dependency rule
                                 Values: same_project_only (default) | cross_project_allowed | disabled
    --feature-dep <rule>         Feature dependency rule
                                 Values: cross_project_allowed (default) | same_project_only | disabled
    --issue-dep <rule>           Issue dependency rule
                                 Values: same_project_only (default) | cross_project_allowed | disabled
    --strict                     Enforce strict dependency rules
    -y, --yes                    Non-interactive mode
  
  Example:
    mandor project create api \
      --name "REST API Service" \
      --goal "Build production REST API with authentication, authorization..." \
      --task-dep cross_project_allowed

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

▶ mandor project update <project_id> --name <name> [--goal <goal>]
  Update project properties
  
  Flags:
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
  Create a new feature in the specified project
  
  Required Arguments:
    <name>                   Feature name (positional or --name flag)
  
  Required Flags:
    --project, -p <id>       Project ID
    --goal, -g <text>        Feature goal (min 300 chars)
  
  Optional Flags:
    --scope <scope>          Feature scope
                             Values: frontend, backend, fullstack, cli, desktop, android, flutter, react-native, ios, swift
    --priority <P0-P5>       Priority level (default: from config)
    --depends <ids>          Pipe-separated feature IDs this depends on
  
  Example:
    mandor feature create "Authentication" --project api \
      --goal "Implement JWT-based authentication with login and refresh flows for secure API access" \
      --scope backend \
      --priority P1

───────────────────────────────────────────────────────────────────────

▶ mandor feature list --project <id> [--json] [--include-deleted]
  List features in a project
  
  Flags:
    --project, -p <id>    Project ID (required)
    --json, -j            Machine-readable JSON output
    --include-deleted     Include deleted features
  
  Example:
    mandor feature list --project api
    mandor feature list --project api --json

───────────────────────────────────────────────────────────────────────

▶ mandor feature detail <feature_id> --project <id>
  Show feature details
  
  Flags:
    --project, -p <id>       Project ID (optional, extracted from feature ID)
    --include-deleted        Include deleted features
  
  Example:
    mandor feature detail api-feature-xxx --project api

───────────────────────────────────────────────────────────────────────

▶ mandor feature update <feature_id> --project <id> [OPTIONS]
  Update feature properties, change status, cancel, or reopen
  
  Flags:
    --project, -p <id>    Project ID (required)
    --name <text>         New feature name
    --goal <text>         New feature goal
    --scope <scope>       New scope (frontend, backend, fullstack, etc.)
    --priority <P0-P5>    New priority
    --status <status>     New status (draft, active, done, blocked, cancelled)
    --depends <ids>       Set all dependencies (pipe-separated)
    --cancel              Cancel the feature
    --reason <text>       Cancellation reason (required with --cancel)
    --reopen              Reopen a cancelled feature
    --force               Force operation (skip validation)
    --dry-run             Show what would be changed
    -y, --yes             Skip confirmation
  
  Examples:
    mandor feature update api-feature-xxx --project api --name "New Name"
    mandor feature update api-feature-xxx --project api --status active
    mandor feature update api-feature-xxx --project api --cancel --reason "Not needed"
    mandor feature update api-feature-xxx --project api --reopen

═════════════════════════════════════════════════════════════════════════
 4. TASK COMMANDS
═════════════════════════════════════════════════════════════════════════

▶ mandor task create <name> --feature <id> --goal <goal> --implementation-steps <steps> --test-cases <cases> --derivable-files <files> --library-needs <libs> [OPTIONS]
  Create a new task in the specified feature
  
  Required Arguments:
    <name>                                 Task name (positional)
  
  Required Flags:
    --feature, -f <id>                     Feature ID
    --goal, -g <text>                      Task goal (min 500 chars)
    --implementation-steps <steps>         Pipe-separated implementation steps
    --test-cases <cases>                   Pipe-separated test cases
    --derivable-files <files>              Pipe-separated files to create
    --library-needs <libs>                 Pipe-separated libraries (use "none" if not needed)
  
  Optional Flags:
    --priority <P0-P5>                     Priority (default: from config)
    --depends-on <ids>                     Pipe-separated task IDs this depends on
    -y, --yes                              Skip confirmation prompts
  
  Example:
    mandor task create "JWT Parser" --feature api-feature-xxx \
      --goal "Parse and validate JWT tokens in incoming requests..." \
      --implementation-steps "Setup crypto|Add validation|Handle expiry" \
      --test-cases "Valid token accepted|Expired rejected|Invalid rejected" \
      --derivable-files "jwt_validator.go|jwt_test.go" \
      --library-needs "golang-jwt" \
      --priority P1

───────────────────────────────────────────────────────────────────────

▶ mandor task list --feature <id> [--status <status>] [--priority <priority>] [--json] [--sort <field>] [--order <asc|desc>]
  List tasks in a feature
  
  Required Flags:
    --feature, -f <id>    Feature ID
  
  Optional Flags:
    --status <status>     Filter by status (pending, ready, in_progress, blocked, done, cancelled)
    --priority <P0-P5>    Filter by priority
    --json, -j            Machine-readable JSON output
    --include-deleted     Include deleted tasks
    --sort <field>        Sort field (priority, created_at, name) (default: priority)
    --order <asc|desc>    Sort order (default: desc)
    --project, -p <id>    Filter by project ID
  
  Examples:
    mandor task list --feature api-feature-xxx
    mandor task list --feature api-feature-xxx --status ready
    mandor task list --feature api-feature-xxx --priority P1 --json

───────────────────────────────────────────────────────────────────────

▶ mandor task ready --feature <id> [--priority <priority>] [--json]
  List tasks with status='ready' that are available to work on
  
  Required Flags:
    --feature, -f <id>    Feature ID
  
  Optional Flags:
    --priority <P0-P5>    Filter by priority
    --json, -j            Machine-readable JSON output
    --project, -p <id>    Filter by project ID
  
  Examples:
    mandor task ready --feature api-feature-xxx
    mandor task ready --feature api-feature-xxx --priority P0

───────────────────────────────────────────────────────────────────────

▶ mandor task blocked --feature <id> [--priority <priority>] [--json]
  List tasks with status='blocked' that are waiting for dependencies
  
  Required Flags:
    --feature, -f <id>    Feature ID
  
  Optional Flags:
    --priority <P0-P5>    Filter by priority
    --json, -j            Machine-readable JSON output
    --project, -p <id>    Filter by project ID
  
  Example:
    mandor task blocked --feature api-feature-xxx

───────────────────────────────────────────────────────────────────────

▶ mandor task detail <task_id>
  Show task details
  
  Example:
    mandor task detail api-task-xxx-001

───────────────────────────────────────────────────────────────────────

▶ mandor task update <task_id> [OPTIONS]
  Update task properties, change status, cancel, or reopen
  
  Flags:
    --name <text>                      New task name
    --goal <text>                      New task goal
    --priority <P0-P5>                 New priority
    --implementation-steps <steps>     Update implementation steps (pipe-separated)
    --test-cases <cases>               Update test cases (pipe-separated)
    --derivable-files <files>          Update derivable files (pipe-separated)
    --library-needs <libs>             Update library needs (pipe-separated)
    --status <status>                  New status (ready, in_progress, done, blocked, cancelled)
    --depends <ids>                    Set all dependencies (pipe-separated)
    --depends-add <ids>                Add dependencies (pipe-separated)
    --depends-remove <ids>             Remove dependencies (pipe-separated)
    --cancel                           Cancel the task
    --reason <text>                    Cancellation reason (required with --cancel)
    --reopen                           Reopen a cancelled task
    --force                            Force operation (skip validation)
    --dry-run                          Show what would be changed
    -y, --yes                          Skip confirmation
  
  Examples:
    mandor task update api-task-xxx-001 --status in_progress
    mandor task update api-task-xxx-001 --status done
    mandor task update api-task-xxx-001 --priority P0
    mandor task update api-task-xxx-001 --cancel --reason "Superseded by task X"

═════════════════════════════════════════════════════════════════════════
 5. ISSUE COMMANDS
═════════════════════════════════════════════════════════════════════════

▶ mandor issue create <name> --project <id> --type <type> --goal <goal> --affected-files <files> --affected-tests <tests> --implementation-steps <steps> [OPTIONS]
  Create an issue (bug, improvement, debt, security, performance)
  
  Required Arguments:
    <name>                             Issue name (positional)
  
  Required Flags:
    --project, -p <id>                 Project ID
    --type, -t <type>                  Issue type (bug, improvement, debt, security, performance)
    --goal, -g <text>                  Issue goal (min 200 chars)
    --affected-files <files>           Pipe-separated affected file paths
    --affected-tests <tests>           Pipe-separated affected test files
    --implementation-steps <steps>     Pipe-separated implementation steps
  
  Optional Flags:
    --priority <P0-P5>                 Priority (default: from config)
    --depends-on <ids>                 Pipe-separated issue IDs for dependencies
    --library-needs <libs>             Pipe-separated required libraries
    -y, --yes                          Skip confirmation
  
  Example:
    mandor issue create "Fix memory leak in auth handler" \
      --project api \
      --type bug \
      --priority P0 \
      --goal "Goroutine not properly cleaned up in token refresh handler causing memory accumulation..." \
      --affected-files "src/handlers/auth.go|src/middleware/auth.go" \
      --affected-tests "src/handlers/auth_test.go" \
      --implementation-steps "Identify leak|Add cleanup|Add tests|Verify"

───────────────────────────────────────────────────────────────────────

▶ mandor issue list [--project <id>] [--type <type>] [--status <status>] [--priority <priority>] [--json] [--sort <field>] [--order <asc|desc>]
  List issues in the specified project with optional filters
  
  Flags:
    --project, -p <id>    Project ID filter
    --type, -t <type>     Filter by type (bug, improvement, debt, security, performance)
    --status <status>     Filter by status (open, ready, in_progress, resolved, wontfix, blocked, cancelled)
    --priority <P0-P5>    Filter by priority
    --json, -j            Machine-readable JSON output
    --sort <field>        Sort field (created_at, last_updated_at, priority, name) (default: last_updated_at)
    --order <asc|desc>    Sort order (default: desc)
    --verbose             Show issue names in table output
  
  Examples:
    mandor issue list --project api
    mandor issue list --project api --type bug --priority P0
    mandor issue list --status open

───────────────────────────────────────────────────────────────────────

▶ mandor issue ready [--project <id>] [--type <type>] [--priority <priority>] [--json]
  List all issues with status='ready' that are available to work on
  
  Flags:
    --project, -p <id>    Project ID filter
    --type, -t <type>     Filter by type (bug, improvement, debt, security, performance)
    --priority <P0-P5>    Filter by priority
    --json, -j            Machine-readable JSON output
  
  Example:
    mandor issue ready --project api --type bug --priority P0

───────────────────────────────────────────────────────────────────────

▶ mandor issue blocked [--project <id>] [--type <type>] [--priority <priority>] [--json]
  List all issues with status='blocked' that are waiting for dependencies
  
  Flags:
    --project, -p <id>    Project ID filter
    --type, -t <type>     Filter by type
    --priority <P0-P5>    Filter by priority
    --json, -j            Machine-readable JSON output
  
  Example:
    mandor issue blocked --project api --type security

───────────────────────────────────────────────────────────────────────

▶ mandor issue detail <issue_id> [--project <id>]
  Show issue details
  
  Flags:
    --project, -p <id>    Project ID (optional, auto-extracted)
  
  Example:
    mandor issue detail api-issue-abc123 --project api

───────────────────────────────────────────────────────────────────────

▶ mandor issue update <issue_id> [--project <id>] [OPTIONS]
  Update an issue's metadata, status, or dependencies
  
  Flags:
    --project, -p <id>              Project ID (optional, extracted from issue ID)
    --name <text>                   Update issue name
    --goal <text>                   Update goal
    --type <type>                   Change issue type
    --priority <P0-P5>              Update priority
    --status <status>               Set status (open, ready, in_progress, resolved, wontfix, blocked, cancelled)
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
    --wontfix                       Mark as wontfix (requires --reason)
    --reopen                        Reopen resolved/wontfix issue
    --cancel                        Cancel issue
    --force                         Force operation (skip validation)
    --dry-run                       Show what would change
  
  Examples:
    mandor issue update api-issue-abc123 --resolve
    mandor issue update api-issue-abc123 --start
    mandor issue update api-issue-abc123 --wontfix --reason "Working as intended"
    mandor issue update api-issue-abc123 --cancel --reason "Duplicate of #123"

═════════════════════════════════════════════════════════════════════════
 6. AI COMMANDS
═════════════════════════════════════════════════════════════════════════

▶ mandor ai agents
  Generate AGENTS.md for multi-agent coordination
  
  Creates documentation for coordinating multiple AI agents on project tasks.
  
  Example:
    mandor ai agents > AGENTS.md

───────────────────────────────────────────────────────────────────────

▶ mandor ai claude
  Generate CLAUDE.md for the project
  
  Creates Claude-specific project documentation.
  
  Example:
    mandor ai claude > CLAUDE.md

═════════════════════════════════════════════════════════════════════════
 7. UTILITY COMMANDS
═════════════════════════════════════════════════════════════════════════

▶ mandor completion [bash|zsh|fish]
  Generate shell completion scripts
  
  Supported Shells: bash, zsh, fish
  
  Setup:
    # Bash
    mandor completion bash > /usr/local/etc/bash_completion.d/mandor
    
    # Zsh
    mandor completion zsh > "${fpath[1]}/_mandor"
    
    # Fish
    mandor completion fish > ~/.config/fish/completions/mandor.fish

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
 8. CONFIGURATION REFERENCE
═════════════════════════════════════════════════════════════════════════

Configuration Directory: .mandor/config.json

Available Keys:
  default_priority      Default priority for new entities
                        Valid values: P0, P1, P2, P3, P4, P5
                        Default: P3

  strict_mode           Enable strict dependency validation
                        Valid values: true, false
                        Default: false

Set Configuration:
  $ mandor config set default_priority P2
  $ mandor config set strict_mode true

Get Configuration:
  $ mandor config get default_priority

List All:
  $ mandor config list

Reset to Default:
  $ mandor config reset default_priority

═════════════════════════════════════════════════════════════════════════
 9. STATUS TRANSITIONS
═════════════════════════════════════════════════════════════════════════

FEATURE STATUS FLOW:
  draft ─→ active ─→ done
    │
    └─→ blocked (waiting on dependencies)
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

  - pending:      Waiting on dependencies (auto-transitioned)
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

  - open:         Newly created, waiting on dependencies
  - ready:        Available to work on (no blocking dependencies)
  - in_progress:  Being worked on
  - resolved:     Problem fixed
  - wontfix:      Decided not to fix (permanent, can reopen)
  - blocked:      Waiting on dependencies
  - cancelled:    Abandoned or duplicate

═════════════════════════════════════════════════════════════════════════
 10. DEPENDENCY RULES
═════════════════════════════════════════════════════════════════════════

TASK DEPENDENCIES:
  • Tasks in same feature can depend on other tasks
  • Dependency rule controlled by project config (--task-dep flag)
  • When dependency is marked done, dependent tasks auto-transition to ready
  • Circular dependencies are detected and prevented
  • Failed dependency creates blocked status

FEATURE DEPENDENCIES:
  • Features can depend on other features
  • Dependency rule controlled by project config (--feature-dep flag)
  • When dependency feature is done, dependent auto-transitions to active (if draft)
  • Can be cross-project (if configured)

ISSUE DEPENDENCIES:
  • Issues can depend on other issues
  • Dependency rule controlled by project config (--issue-dep flag)
  • When dependency issue is resolved, dependent auto-transitions to ready
  • Blocking detection is automatic

BLOCKING SCENARIOS:
  • Task depends on incomplete task → Task becomes blocked
  • Feature depends on incomplete feature → Feature becomes blocked
  • Issue depends on unresolved issue → Issue becomes blocked
  • Mark dependency complete → Dependent auto-transitions to ready

═════════════════════════════════════════════════════════════════════════
 11. BEST PRACTICES
═════════════════════════════════════════════════════════════════════════

1. USE MEANINGFUL IDs
   
   Project and feature IDs should be:
   • Short but descriptive
   • Lowercase with hyphens
   • Consistent naming convention
   
   ✓ Good:   mandor project create user-auth
   ✗ Avoid:  mandor project create p1

2. WRITE CLEAR GOALS
   
   Goals should include:
   • What is being built/fixed
   • Why it matters
   • Technical requirements
   • Acceptance criteria
   
   ✓ Good:   "Implement JWT-based authentication with login and refresh flows for secure API access"
   ✗ Avoid:  "Add authentication"

3. USE SCOPES FOR FEATURES
   
   Organize by scope:
   • frontend, backend, fullstack
   • cli, desktop, android, flutter, react-native, ios, swift
   
   Example:
     mandor feature create "Login UI" --project api --scope frontend
     mandor feature create "Login API" --project api --scope backend

4. KEEP DEPENDENCIES SHALLOW
   
   • Deep chains (>5 levels) are hard to manage
   • Consider breaking into smaller features
   • Use --depends-on sparingly

5. USE ISSUES FOR BUGS, TASKS FOR FEATURES
   
   • Tasks: Feature work, implementation, refactoring
   • Issues: Bugs, improvements, technical debt, security, performance
   
   Example:
     mandor task create "Add OAuth2" --feature api-auth
     mandor issue create "Fix auth timeout" --project api --type bug

6. DOCUMENT CANCELLATION REASONS
   
   Always provide clear reasons when cancelling:
   
     mandor task update <id> --cancel --reason "Superseded by feature X"
     mandor feature update <id> --project api --cancel --reason "Sticking with JWT"

7. USE PIPE-SEPARATED LISTS
   
   For flags accepting multiple values, use pipes:
   
     --implementation-steps "Step 1|Step 2|Step 3"
     --test-cases "Case 1|Case 2|Case 3"
     --depends-on task-1|task-2|task-3

8. USE --dry-run FOR PREVIEW
   
   Before making significant updates, preview with --dry-run:
   
     mandor task update task-id --status done --dry-run
     mandor feature update feature-id --project api --cancel --reason "..." --dry-run

9. SET CONFIGURATION EARLY
   
   Configure workspace defaults at the start:
   
     mandor init "Project Name"
     mandor config set default_priority P2
     mandor config set strict_mode true

10. REVIEW STATUS REGULARLY
    
    Keep team synchronized:
    
      mandor status                                # Workspace overview
      mandor status --project api                  # Project summary
      mandor summary --project api                 # Feature priorities
      mandor task ready --feature feature-id       # See available work
      mandor issue ready --project api             # See ready issues
      mandor task blocked --feature feature-id     # See blockers
      mandor issue blocked --project api           # See blocked issues

═════════════════════════════════════════════════════════════════════════
 12. COMMON WORKFLOWS
═════════════════════════════════════════════════════════════════════════

SETUP NEW PROJECT
═══════════════════════════════════════════════════════════════════════

1. Initialize workspace:
   mandor init "Project Name" -y

2. Configure defaults:
   mandor config set default_priority P2
   mandor config set strict_mode true

3. Create project:
   mandor project create api \
     --name "API Service" \
     --goal "Build REST API with authentication and data endpoints"

4. Create feature:
   mandor feature create "Authentication" --project api \
     --goal "Implement JWT with login and refresh flows"

TRACK FEATURE WORK
═══════════════════════════════════════════════════════════════════════

1. Create tasks:
   mandor task create "JWT Parser" --feature api-feature-xxx \
     --goal "..." --implementation-steps "..." --test-cases "..." \
     --derivable-files "..." --library-needs "..." --priority P1

2. Add dependent task:
   mandor task create "Login Endpoint" --feature api-feature-xxx \
     --goal "..." --depends-on jwt-task-id --priority P1

3. Check what's available:
   mandor task ready --feature api-feature-xxx

4. Start working:
   mandor task update task-id --status in_progress

5. Complete task:
   mandor task update task-id --status done
   # Dependents auto-transition to ready

6. Check progress:
   mandor task ready --feature api-feature-xxx

ISSUE TRACKING
═══════════════════════════════════════════════════════════════════════

1. Create bug issue:
   mandor issue create "Fix timeout in auth" \
     --project api --type bug --priority P1 \
     --goal "..." --affected-files "..." --affected-tests "..." \
     --implementation-steps "..."

2. List open issues:
   mandor issue list --project api --status open

3. See what's ready:
   mandor issue ready --project api --type bug

4. Start work:
   mandor issue update issue-id --start

5. Mark resolved:
   mandor issue update issue-id --resolve

DEPENDENCY MANAGEMENT
═══════════════════════════════════════════════════════════════════════

1. View all projects:
   mandor status

2. Check blockers:
   mandor task blocked --feature feature-id
   mandor issue blocked --project api

3. Add dependency:
   mandor task update task-id --depends-add dependent-task-id

4. Remove dependency:
   mandor task update task-id --depends-remove dependent-task-id

5. View feature dependencies:
   mandor feature list --project api

CANCEL AND REOPEN
═══════════════════════════════════════════════════════════════════════

1. Cancel with reason:
   mandor task update task-id --cancel --reason "Out of scope"
   mandor feature update feature-id --project api --cancel --reason "Not needed"

2. Reopen:
   mandor task update task-id --reopen
   mandor feature update feature-id --project api --reopen

CROSS-PROJECT DEPENDENCIES
═══════════════════════════════════════════════════════════════════════

1. Create projects:
   mandor project create frontend --name "Frontend" --goal "..." \
     --task-dep cross_project_allowed
   mandor project create backend --name "Backend" --goal "..." \
     --task-dep cross_project_allowed

2. Create task in backend:
   mandor feature create "API" --project backend --goal "..."
   mandor task create "Auth Endpoint" --feature backend-feature-xxx \
     --goal "..." --priority P1

3. Create dependent task in frontend:
   mandor feature create "Login" --project frontend --goal "..."
   mandor task create "Login Form" --feature frontend-feature-yyy \
     --goal "..." --depends-on backend-task-xxx --priority P1

4. Backend task completion auto-unblocks frontend:
   mandor task update backend-task-xxx --status done
   mandor task list --feature frontend-feature-yyy
   # "Login Form" now shows as ready

═════════════════════════════════════════════════════════════════════════

For more information, visit: https://github.com/sanxzy/mandor

Built for AI Agent Workflows
`)

	return nil
}
