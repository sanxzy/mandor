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
║        Deterministic Task Manager for AI Agent Workflows           ║
║              Structured Storage | Dependency-Aware                 ║
╚════════════════════════════════════════════════════════════════════╝

WHY MANDOR
═══════════════════════════════════════════════════════════════════════

Stop writing markdown plans that go stale. Mandor provides:
  ✓ Single source of truth in structured JSONL format
  ✓ Automatic dependency resolution
  ✓ Real-time status visibility with track commands
  ✓ Full audit trail of all changes
  ✓ Schema-driven task definition
  ✓ CLI-native, works in scripts and CI/CD

OLD APPROACH:
  # PLAN.md (becomes outdated)
  - [ ] Task A
  - [ ] Task B (depends on A, but nobody tracks this)
  Status: Last updated 3 days ago

NEW APPROACH:
   mandor feature create "Authentication" --project api --goal "..."
   mandor task create api-feature-xxxx "Task A" --goal "..." \
     --implementation-steps "Step 1|Step 2" \
     --test-cases "Test A|Test B" \
     --library-needs "golang-jwt"
   mandor task create api-feature-xxxx "Task B" --goal "..." \
     --implementation-steps "Step 1|Step 2" \
     --test-cases "Test A|Test B" \
     --depends-on api-feature-xxxx-task-xxxx
  
  mandor track feature api-feature-xxxx    # Always current

════════════════════════════════════════════════════════════════════════
 TABLE OF CONTENTS
════════════════════════════════════════════════════════════════════════

  1. Workspace Setup        6. Task Management
  2. Project Management     7. Issue Management
  3. Feature Management     8. Track Commands
  4. Configuration          9. Session Management
  5. Status Transitions    10. AI Commands
                          11. Best Practices

════════════════════════════════════════════════════════════════════════
 1. WORKSPACE SETUP
════════════════════════════════════════════════════════════════════════

▶ mandor init [--workspace-name <name>] [-y]
  Initialize workspace in current directory
  
  Flags:
    --workspace-name <text>   Custom name (default: current directory)
    -y, --yes                 Skip confirmation
  
  Example:
    mandor init "My Project" -y

▶ mandor track [scope] [id] [FLAGS]
   View workspace, project, feature, task, or issue status (see section 8)
   
   Scopes:
     (none or workspace)    Workspace overview
     project <id>           Project issues
     feature <id>           Feature with tasks
     task <id>              Single task details
     issue <id>             Single issue details
   
   Example:
     mandor track --json
     mandor track project api

▶ mandor config [get|set|list|reset] [key] [value]
  Manage workspace configuration
  
  Available keys:
    default_priority      Default priority (P0-P5, default: P3)
    strict_mode          Strict validation (true/false, default: false)
    goal.lengths.project  Min chars for project goal (default: 500)
    goal.lengths.feature Min chars for feature goal (default: 300)
    goal.lengths.task    Min chars for task goal (default: 500)
    goal.lengths.issue   Min chars for issue goal (default: 200)
  
  Example:
    mandor config set default_priority P2
    mandor config get default_priority
    mandor config list

════════════════════════════════════════════════════════════════════════
 2. PROJECT MANAGEMENT
════════════════════════════════════════════════════════════════════════

▶ mandor project create <id> --name <text> --goal <text> [FLAGS]
  Create a new project
  
  Required:
    <id>               Project identifier (lowercase, hyphens)
    --name <text>      Display name
    --goal <text>      Goal description (min 500 chars)
  
  Optional:
    --task-dep <rule>       Task dependency rule
                           (same_project_only | cross_project_allowed | disabled)
                           Default: same_project_only
    --feature-dep <rule>    Feature dependency rule
                           (same_project_only | cross_project_allowed | disabled)
                           Default: cross_project_allowed
    --issue-dep <rule>      Issue dependency rule
                           (same_project_only | cross_project_allowed | disabled)
                           Default: same_project_only
    -y, --yes          Non-interactive mode
  
  Example:
    mandor project create api --name "API Service" \
      --goal "REST API with auth, user, and data endpoints" \
      --task-dep cross_project_allowed

▶ mandor project detail <id>
  Show project details (metadata, counts, dependencies)
  
  Flags:
    --json    Machine-readable output
  
  Example:
    mandor project detail api --json

▶ mandor track [workspace]
   List all projects (use track for visibility, see section 8)
   
   Example:
     mandor track
     mandor track --json

▶ mandor project update <id> [FLAGS]
  Update project metadata
  
  Flags:
    --name <text>           New name
    --goal <text>           New goal (min 500 chars)
    --task-dep <rule>       Update task dependency rule
    --feature-dep <rule>    Update feature dependency rule
    --issue-dep <rule>      Update issue dependency rule
    --strict <bool>         Toggle strict mode
  
  Example:
    mandor project update api --goal "Enhanced API with new features"

▶ mandor project delete <id> [--hard]
  Delete project (soft delete by default)
  
  Flags:
    --hard    Permanently delete (cannot be restored)
    -y, --yes Skip confirmation
  
  Example:
    mandor project delete legacy --hard -y

▶ mandor project reopen <id>
  Restore soft-deleted project
  
  Example:
    mandor project reopen legacy

════════════════════════════════════════════════════════════════════════
 3. FEATURE MANAGEMENT
════════════════════════════════════════════════════════════════════════

▶ mandor feature create <name> --project <id> --goal <text> [FLAGS]
  Create feature in project
  
  Required:
    <name>            Feature name (positional)
    --project, -p <id> Project ID
    --goal <text>      Goal description (min 300 chars)
  
  Optional:
    --scope <scope>     Feature scope
                       (frontend | backend | fullstack | cli | desktop | 
                        android | flutter | react-native | ios | swift)
    --priority <P0-P5>  Priority (default: from config)
    --depends <ids>     Pipe-separated feature IDs this depends on
  
  Example:
    mandor feature create "Authentication" --project api \
      --goal "JWT-based authentication with login and refresh flows" \
      --scope backend --priority P1

▶ mandor feature detail <id> --project <id> [FLAGS]
  Show feature details
  
  Flags:
    --project, -p <id>  Project ID
    --json              Machine-readable output
    --include-deleted   Include cancelled features
  
  Example:
    mandor feature detail api-feature-xxxx --project api --json

▶ mandor track project <id>
   List features in project (use track for visibility, see section 8)
   
   Example:
     mandor track project api --json

▶ mandor feature update <id> --project <id> [FLAGS]
  Update feature or change status
  
  Flags:
    --project, -p <id>  Project ID (required)
    --name <text>       New name
    --goal <text>       New goal
    --scope <scope>     New scope
    --priority <P0-P5>  New priority
    --status <status>   New status (draft | active | done | blocked | cancelled)
    --depends <ids>     Set dependencies (pipe-separated)
    --cancel            Cancel feature
    --reason <text>     Cancellation reason (required with --cancel)
    --reopen            Reopen cancelled feature
    --force             Force operation (skip validation)
    --dry-run           Preview changes
    -y, --yes           Skip confirmation
  
  Example:
    mandor feature update api-feature-xxxx --project api --status active
    mandor feature update api-feature-xxxx --project api --cancel --reason "Out of scope"

════════════════════════════════════════════════════════════════════════
 4. CONFIGURATION
════════════════════════════════════════════════════════════════════════

▶ mandor config list
  Show all configuration keys with descriptions
  
  Example:
    mandor config list

▶ mandor config get [key]
  Get configuration value
  
  Example:
    mandor config get default_priority

▶ mandor config set <key> <value>
  Set configuration value
  
  Example:
    mandor config set default_priority P2
    mandor config set strict_mode true

▶ mandor config reset [key]
  Reset configuration to default
  
  Example:
    mandor config reset default_priority

════════════════════════════════════════════════════════════════════════
 5. STATUS TRANSITIONS
════════════════════════════════════════════════════════════════════════

FEATURE STATUS FLOW:
  draft ──→ active ──→ done
    │         │
    └─→ blocked (on dependencies)
    └─→ cancelled (with reason)

TASK STATUS FLOW:
  (auto-pending if depends on incomplete)
     ↓
  ready ──→ in_progress ──→ done
    ↑           │
    └──────────╯ blocked (on dependencies)
    └─→ cancelled (with reason)

ISSUE STATUS FLOW:
  open ──→ ready ──→ in_progress ──→ resolved
   │        │           │
   └─→ blocked ──→ ready (when dependency resolves)
   └─→ wontfix (with reason)
   └─→ cancelled (with reason)

════════════════════════════════════════════════════════════════════════
 6. TASK MANAGEMENT
════════════════════════════════════════════════════════════════════════

▶ mandor task create <feature_id> <name> --goal <text> \
                     --implementation-steps <steps> \
                     --test-cases <cases> [FLAGS]
  Create task in feature
  
  Required:
    <feature_id>                Feature ID (positional)
    <name>                      Task name (positional)
    --goal <text>               Goal (min 500 chars)
    --implementation-steps      Pipe-separated steps
    --test-cases                Pipe-separated test cases
  
  Optional:
    --library-needs <libs>      Pipe-separated libraries (optional, use "none" if not needed)
    --priority <P0-P5>          Priority (default: from config)
    --depends-on <ids>          Pipe-separated task IDs
    -y, --yes                   Skip confirmation
  
  Example:
    mandor task create api-feature-xxxx "JWT Parser" \
      --goal "Parse and validate JWT tokens in requests" \
      --implementation-steps "Setup crypto|Add validation|Handle expiry" \
      --test-cases "Valid accepted|Expired rejected|Invalid rejected" \
      --library-needs "golang-jwt"

▶ mandor task detail <id> [FLAGS]
  Show task details
  
  Flags:
    --json            Machine-readable output
    --events          Show event history
    --dependencies    Show dependency info
    --timestamps      Show formatted timestamps
    --include-deleted Include cancelled tasks
  
  Example:
    mandor task detail api-feature-xxxx-task-xxxx --json

▶ mandor task update <id> [FLAGS]
  Update task or change status
  
  Flags:
    --name <text>                      New name
    --goal <text>                      New goal
    --priority <P0-P5>                 New priority
    --status <status>                  New status (ready | in_progress | done)
    --implementation-steps <steps>     Update steps (pipe-separated)
    --test-cases <cases>               Update test cases (pipe-separated)
    --derivable-files <files>          Update files (pipe-separated)
    --library-needs <libs>             Update libraries (pipe-separated)
    --depends <ids>                    Replace dependencies (pipe-separated)
    --depends-add <ids>                Add dependencies (pipe-separated)
    --depends-remove <ids>             Remove dependencies (pipe-separated)
    --cancel                           Cancel task
    --reason <text>                    Cancellation reason (required with --cancel)
    --reopen                           Reopen cancelled task
    --force                            Force operation (skip checks)
    --dry-run                          Preview changes
    -y, --yes                          Skip confirmation
  
  Examples:
    mandor task update api-feature-xxxx-task-xxxx --status in_progress
    mandor task update api-feature-xxxx-task-xxxx --status done
    mandor task update api-feature-xxxx-task-xxxx --cancel --reason "Superseded by task Y"

════════════════════════════════════════════════════════════════════════
 7. ISSUE MANAGEMENT
════════════════════════════════════════════════════════════════════════

▶ mandor issue create <name> --project <id> --type <type> \
                      --goal <text> --affected-files <files> \
                      --affected-tests <tests> \
                      --implementation-steps <steps> [FLAGS]
  Create issue in project
  
  Required:
    <name>                    Issue name (positional)
    --project, -p <id>        Project ID
    --type, -t <type>         Issue type (bug | improvement | debt | security | performance)
    --goal <text>             Goal (min 200 chars)
    --affected-files <files>  Pipe-separated affected file paths
    --affected-tests <tests>  Pipe-separated affected test files
    --implementation-steps    Pipe-separated steps
  
  Optional:
    --priority <P0-P5>        Priority (default: from config)
    --depends-on <ids>        Pipe-separated issue dependencies
    --library-needs <libs>    Pipe-separated required libraries
    -y, --yes                 Skip confirmation
  
  Example:
    mandor issue create "Memory leak in auth handler" \
      --project api --type bug --priority P0 \
      --goal "Goroutine not cleaned up in token refresh handler" \
      --affected-files "src/handlers/auth.go|src/middleware/auth.go" \
      --affected-tests "src/handlers/auth_test.go" \
      --implementation-steps "Identify|Fix|Add tests|Verify" \
      --library-needs "none"

▶ mandor issue detail <id> [--project <id>] [FLAGS]
  Show issue details
  
  Flags:
    --project, -p <id>  Project ID (optional, auto-extracted)
    --json              Machine-readable output
    --events            Show event history
    --include-deleted   Include cancelled issues
  
  Example:
    mandor issue detail api-issue-abc123 --json

▶ mandor issue update <id> [FLAGS]
  Update issue or change status
  
  Flags:
    --project, -p <id>              Project ID (optional)
    --name <text>                   New name
    --goal <text>                   New goal
    --type <type>                   New type (bug|improvement|debt|security|performance)
    --priority <P0-P5>              New priority
    --status <status>               New status (open|ready|in_progress|resolved|wontfix|blocked|cancelled)
    --reason <text>                 Reason for status change
    --affected-files <files>        Replace affected files (pipe-separated)
    --affected-tests <tests>        Replace affected tests (pipe-separated)
    --implementation-steps <steps>  Replace implementation steps (pipe-separated)
    --library-needs <libs>          Replace library needs (pipe-separated)
    --depends <ids>                  Replace dependencies (pipe-separated)
    --depends-add <ids>             Add dependencies (pipe-separated)
    --depends-remove <ids>          Remove dependencies (pipe-separated)
    --start                         Transition to in_progress
    --resolve                       Mark as resolved
    --wontfix                       Mark as wontfix (requires --reason)
    --reopen                        Reopen resolved/wontfix issue
    --cancel                        Cancel issue
    --force                         Force operation (skip checks)
    --dry-run                       Preview changes
    -y, --yes                       Skip confirmation
  
  Examples:
    mandor issue update api-issue-abc123 --resolve
    mandor issue update api-issue-abc123 --start
    mandor issue update api-issue-abc123 --wontfix --reason "Working as intended"

════════════════════════════════════════════════════════════════════════
 8. TRACK COMMANDS (Real-Time Visibility)
════════════════════════════════════════════════════════════════════════

▶ mandor track
  Show workspace overview (all projects and summaries)
  
  Flags:
    --json          Machine-readable output
    --csv           CSV export format
    --tree          Tree visualization
    --graph         ASCII graph
    --verbose       Show all fields with details
    --group-by <f>  Group by status or priority
  
  Example:
    mandor track
    mandor track --json

▶ mandor track workspace
  Show workspace overview (explicit)
  
  Example:
    mandor track workspace

▶ mandor track project <id>
  Show project issues and status
  
  Flags:
    --json          Machine-readable output
    --csv           CSV export
    --tree          Tree visualization
    --graph         ASCII graph
    --verbose       Show all fields
    --group-by <f>  Group by status or priority
  
  Example:
    mandor track project api --json

▶ mandor track feature <id>
  Show feature with all tasks
  
  Flags:
    --json          Machine-readable output
    --csv           CSV export
    --tree          Tree visualization
    --graph         ASCII graph
    --verbose       Show all fields
    --group-by <f>  Group by status or priority
  
  Example:
    mandor track feature api-feature-xxxx --verbose

▶ mandor track task <id>
  Show single task details
  
  Flags:
    --json    Machine-readable output
    --verbose Show all fields
  
  Example:
    mandor track task api-feature-xxxx-task-xxxx

▶ mandor track issue <id>
  Show single issue details
  
  Flags:
    --json    Machine-readable output
    --verbose Show all fields
  
  Example:
    mandor track issue api-issue-abc123

════════════════════════════════════════════════════════════════════════
 9. SESSION MANAGEMENT (AI Agent Progress Tracking)
════════════════════════════════════════════════════════════════════════

▶ mandor session note [text]
  Add a timestamped note about work completed or in progress
  
  Notes are stored in .mandor/session-notes.jsonl as NDJSON format.
  This provides a lightweight way for AI agents to track progress across sessions.
  
  Flags:
    -r, --read                  Read recent notes instead of adding
    -o, --offset <count>        Number of notes to show (default: 50)
  
  Examples:
    mandor session note "Completed v0.4.4 release and testing"
    mandor session note "Started performance optimization - blocked on benchmarks"
    mandor session note --read                 # Show last 50 notes
    mandor session note --read --offset 100    # Show last 100 notes

════════════════════════════════════════════════════════════════════════
 10. AI COMMANDS
════════════════════════════════════════════════════════════════════════

▶ mandor ai agents
  Generate AGENTS.md for multi-agent coordination
  
  Example:
    mandor ai agents > AGENTS.md

▶ mandor ai claude
  Generate CLAUDE.md for the project
  
  Example:
    mandor ai claude > CLAUDE.md

════════════════════════════════════════════════════════════════════════
 11. BEST PRACTICES
════════════════════════════════════════════════════════════════════════

1. USE MEANINGFUL IDS
   ✓ Good:  user-auth, auth-api, login-flow
   ✗ Avoid: p1, f123, feature-a

2. WRITE CLEAR GOALS (min char requirement enforced)
   ✓ Include what, why, technical requirements
   ✗ Avoid vague descriptions like "Add auth"

3. USE SCOPES FOR FEATURES
   frontend, backend, fullstack, cli, desktop, android, ios, etc.

4. KEEP DEPENDENCIES SHALLOW
   • Avoid chains deeper than 5 levels
   • Consider splitting into smaller features

5. USE ISSUES FOR BUGS, TASKS FOR FEATURES
   • Tasks: Features, implementation, refactoring
   • Issues: Bugs, improvements, tech debt, security, performance

6. DOCUMENT CANCELLATION REASONS
   Always provide reasons when cancelling:
   mandor task update id --cancel --reason "Superseded by feature X"

7. USE PIPE SEPARATORS FOR LISTS
   Pipe (|) separates multiple values:
   --implementation-steps "Step 1|Step 2|Step 3"
   --depends-on task-1|task-2|task-3

8. USE --dry-run BEFORE MAJOR CHANGES
   Preview changes without applying:
   mandor task update task-id --status done --dry-run

9. DEPENDENCY AUTO-RESOLUTION
   • Mark task done → dependents auto-transition to ready
   • Mark issue resolved → dependents auto-transition to ready
   • Manual block → must manually unblock

10. CONFIGURE EARLY, RARELY CHANGE
    Set defaults at workspace setup:
    mandor config set default_priority P2
    mandor config set strict_mode true

════════════════════════════════════════════════════════════════════════
 QUICK WORKFLOWS
════════════════════════════════════════════════════════════════════════

SETUP NEW PROJECT:
   1. mandor init "Project" -y
   2. mandor project create api --name "API Service" --goal "..."
   3. mandor feature create "Authentication" --project api --goal "..."
   4. mandor task create api-feature-xxxx "Task" --goal "..." \
        --implementation-steps "Step 1|Step 2" \
        --test-cases "Test A|Test B"

TRACK PROGRESS:
   1. mandor track                              # Workspace overview
   2. mandor track project api                  # Project issues
   3. mandor track feature api-feature-xxxx     # Feature tasks
   4. mandor task update task-id --status in_progress
   5. mandor task update task-id --status done  # Auto-unblocks dependents

MANAGE BLOCKERS:
   1. mandor track project api                  # See all issues
   2. Add dependency: mandor task update id --depends-add other-id
   3. Remove when ready: mandor task update id --depends-remove other-id
   4. View feature: mandor track feature api-feature-xxxx

CREATE DEPENDENT TASK:
   1. mandor task create api-feature-xxxx "Task A" --goal "..." \
        --implementation-steps "Step 1|Step 2" \
        --test-cases "Test A|Test B"
   2. mandor task create api-feature-xxxx "Task B" --goal "..." \
        --implementation-steps "Step 1|Step 2" \
        --test-cases "Test A|Test B" \
        --depends-on api-feature-xxxx-task-xxxx
   3. Task B auto-transitions to ready when Task A is done

════════════════════════════════════════════════════════════════════════

For more information:
  • GitHub: https://github.com/sanxzy/mandor
  • Use: mandor [command] --help for detailed flag information

Built for AI Agent Workflows
`)
	return nil
}
