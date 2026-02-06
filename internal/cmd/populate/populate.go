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
║        Event-Based Task Manager for AI Agent Workflows             ║
║     Brief → Spec → Blueprint → Feature → Task Workflow             ║
╚════════════════════════════════════════════════════════════════════╝

WHY MANDOR
═══════════════════════════════════════════════════════════════════════

Stop writing markdown plans that go stale. Mandor provides:
  ✓ Gate-enforced task progression (brief, spec, session notes read gates)
  ✓ Automatic dependency resolution with auto-blocking & auto-unblocking
  ✓ Cross-feature task dependencies with cascading transitions
  ✓ Structured Brief → Spec → Blueprint → Feature → Task workflow
  ✓ Real-time status visibility with track commands
  ✓ Full audit trail of all changes
  ✓ Schema-driven definitions with goal length validation
  ✓ CLI-native, works in scripts and CI/CD

WORKFLOW OVERVIEW
═══════════════════════════════════════════════════════════════════════

1. Create Project
   mandor project create api --name "API Service" --goal "..."

2. Create Brief (project intent & capabilities)
   mandor brief create --project api --name "Auth System" --why "..." \
     --capabilities "jwt-auth:JWT authentication|refresh:Token refresh"

3. Create Specs (for each capability, with requirements & IAE scenarios)
   mandor spec create --project api --capability jwt-auth \
     --summary "JWT authentication spec" \
     --requirements "req1:intent1:action1:expect1|req2:..."

4. Create Blueprint (technical architecture linking brief+specs)
   mandor blueprint create --project api --brief brief-id \
     --problem "Secure authentication" \
     --decisions "Use JWT|Store in HTTP-only|Refresh on expiry"

5. Create Feature (from spec with one-to-one mapping)
   mandor feature create "JWT Authentication" --project api \
     --capability jwt-auth --spec-id spec-id --goal "..."

6. Create Tasks (with IAE scenarios from spec requirements)
   mandor task create feature-id "Setup JWT" --spec-id spec-id \
     --iae-scenarios "req-0001:scenario-0001|req-0002:scenario-0001" \
     --goal "..." --implementation-steps "s1|s2" --test-cases "t1|t2"

7. Set Gates & Start Work
   mandor task set-gate task-id --is-read-brief
   mandor task set-gate task-id --is-read-spec
   mandor task set-gate task-id --is-read-session-notes
   mandor task update task-id --status in_progress

8. Track Progress
   mandor track feature feature-id

════════════════════════════════════════════════════════════════════════
 TABLE OF CONTENTS
════════════════════════════════════════════════════════════════════════

   1. Workspace Setup        7. Track Commands
   2. Project Management     8. Task Management
   3. Brief Management       9. Issue Management
   4. Spec Management       10. Gate Enforcement
   5. Blueprint Management  11. Status Transitions
   6. Feature Management    12. Best Practices
                           13. Quick Workflows

════════════════════════════════════════════════════════════════════════
 1. WORKSPACE SETUP
════════════════════════════════════════════════════════════════════════

▶ mandor init [--workspace-name <name>] [-y]
  Initialize workspace in current directory
  
  Flags:
    --workspace-name <text>   Custom name (default: current directory)
    -y, --yes                 Skip confirmation
  
  Example:
    mandor init -y
    mandor init --workspace-name "My Project" -y

▶ mandor config [get|set|list|reset] [key] [value]
  Manage workspace configuration
  
  Available keys:
    default_priority       Default priority (P0-P5, default: P3)
    strict_mode           Strict validation (true/false, default: false)
    goal.lengths.project  Min chars for project goal (default: 500)
    goal.lengths.feature  Min chars for feature goal (default: 300)
    goal.lengths.task    Min chars for task goal (default: 500)
    goal.lengths.issue   Min chars for issue goal (default: 200)
  
  Examples:
    mandor config set default_priority P2
    mandor config get default_priority
    mandor config list

════════════════════════════════════════════════════════════════════════
 2. PROJECT MANAGEMENT
════════════════════════════════════════════════════════════════════════

▶ mandor project create <id> --name <text> --goal <text>
  Create a new project
  
  Required:
    <id>           Project identifier (lowercase, hyphens)
    --name <text>  Display name
    --goal <text>  Goal description (min 500 chars)
  
  Example:
    mandor project create api --name "API Service" \
      --goal "REST API with JWT auth, user management, and data endpoints"

▶ mandor project detail <id>
  Show project details
  
  Flags:
    --json    Machine-readable output
  
  Example:
    mandor project detail api --json

▶ mandor project update <id> [FLAGS]
  Update project metadata
  
  Flags:
    --name <text>           New name
    --goal <text>           New goal
  
  Example:
    mandor project update api --goal "Enhanced API with new features"

▶ mandor project delete <id> [--hard]
  Delete project (soft delete by default)
  
  Flags:
    --hard    Permanently delete
    -y, --yes Skip confirmation
  
  Example:
    mandor project delete legacy --hard -y

▶ mandor project reopen <id>
  Restore soft-deleted project
  
  Example:
    mandor project reopen legacy

════════════════════════════════════════════════════════════════════════
 3. BRIEF MANAGEMENT
════════════════════════════════════════════════════════════════════════

▶ mandor brief create -p <project> --name <text> --why <text> \
                      --capabilities <cap:desc>|<cap:desc>
  Create a Brief (project intent & capabilities)
  
  Required flags:
    -p, --project <id>    Project ID
    --name <text>         Brief name
    --why <text>          Problem statement & motivation (100-5000 chars)
    --capabilities <text> Format: name:description|name:description
  
  Optional flags:
    --tech-stack <text>      Technical stack (comma-separated)
    --affected-systems <text> Affected systems (comma-separated)
    --dependencies <text>     Dependencies (comma-separated)
  
  Example:
    mandor brief create -p api --name "JWT Authentication" \
      --why "Need secure authentication for API with token refresh" \
      --capabilities "jwt-auth:JWT-based login flow|refresh:Token refresh endpoint" \
      --tech-stack "golang,jwt,redis"

▶ mandor brief read -p <project> <brief-id>
  Read a Brief (display its contents)
  
  Example:
    mandor brief read -p api brief-jwt-auth

▶ mandor brief update -p <project> <brief-id> [FLAGS]
  Update a Brief
  
  Example:
    mandor brief update -p api brief-jwt-auth --name "Enhanced Auth"

▶ mandor brief delete -p <project> <brief-id> [-y]
  Delete a Brief
  
  Example:
    mandor brief delete -p api brief-jwt-auth -y

▶ mandor brief validate -p <project> <brief-id>
  Validate Brief structure
  
  Example:
    mandor brief validate -p api brief-jwt-auth

════════════════════════════════════════════════════════════════════════
 4. SPEC MANAGEMENT (Specifications)
════════════════════════════════════════════════════════════════════════

▶ mandor spec create -p <project> --capability <cap-id> \
                     --summary <text> --requirements <req:intent:action:expect>|...
  Create a Spec for a Brief capability
  
  Required flags:
    -p, --project <id>        Project ID
    --capability <id>         Capability ID from Brief
    --summary <text>          Brief spec description
    --requirements <text>     Format: summary:intent:action:expect|summary:...
  
  Note: Each requirement automatically generates req-0001, req-0002, etc.
        Scenarios can be added with composite IDs: req-0001:scenario-0001
  
  Example:
    mandor spec create -p api --capability jwt-auth \
      --summary "JWT authentication specification" \
      --requirements "Setup JWT:Enable auth:Import library|Validate:Verify token|Refresh:Update token"
    # Creates spec-id (from summary) with auto-generated requirement IDs

▶ mandor spec detail -p <project> <spec-id>
  Display Spec details
  
  Flags:
    --json    Machine-readable output
  
  Example:
    mandor spec detail -p api jwt-auth-spec --json

▶ mandor spec update -p <project> <spec-id> [FLAGS]
  Update a Spec
  
  Example:
    mandor spec update -p api jwt-auth-spec --summary "Enhanced spec"

▶ mandor spec delete -p <project> <spec-id> [-y]
  Delete a Spec
  
  Example:
    mandor spec delete -p api jwt-auth-spec -y

▶ mandor spec validate -p <project> <spec-id>
  Validate Spec structure and requirements
  
  Example:
    mandor spec validate -p api jwt-auth-spec

════════════════════════════════════════════════════════════════════════
 5. BLUEPRINT MANAGEMENT (Technical Architecture)
════════════════════════════════════════════════════════════════════════

▶ mandor blueprint create -p <project> --brief <brief-id> \
                          --problem <text> --decisions <decision>|...
  Create Blueprint linking Brief and its Specs
  
  Required flags:
    -p, --project <id>     Project ID
    --brief <id>           Brief ID (all capabilities must have valid Specs)
    --problem <text>       Problem statement
    --decisions <text>     Format: title:rationale|title:rationale
  
  Optional flags:
    --user-types <text>       User types (comma-separated)
    --goals-in-scope <text>   In-scope goals (comma-separated)
    --goals-out-scope <text>  Out-of-scope goals (comma-separated)
    --constraints <text>      Constraints (comma-separated)
    --implementation <text>   Implementation strategy
    --risks <text>            Format: description:mitigation|...
  
  Example:
    mandor blueprint create -p api --brief brief-jwt-auth \
      --problem "Secure API authentication and authorization" \
      --decisions "Use JWT for stateless auth|HTTP-only cookies for tokens|Refresh on expiry" \
      --user-types "Admin,User,Guest" \
      --goals-in-scope "Authentication,Authorization,Token refresh"

▶ mandor blueprint detail -p <project> <blueprint-id>
  Display Blueprint details
  
  Flags:
    --json    Machine-readable output
  
  Example:
    mandor blueprint detail -p api api-blueprint --json

▶ mandor blueprint update -p <project> <blueprint-id> [FLAGS]
  Update Blueprint
  
  Example:
    mandor blueprint update -p api api-blueprint --problem "Updated problem"

▶ mandor blueprint delete -p <project> <blueprint-id> [-y]
  Delete Blueprint
  
  Example:
    mandor blueprint delete -p api api-blueprint -y

▶ mandor blueprint validate -p <project> <blueprint-id>
  Validate Blueprint structure
  
  Example:
    mandor blueprint validate -p api api-blueprint

════════════════════════════════════════════════════════════════════════
 6. FEATURE MANAGEMENT
════════════════════════════════════════════════════════════════════════

▶ mandor feature create <name> -p <project> --capability <cap-id> \
                        --spec-id <spec-id> -g <goal> [FLAGS]
  Create feature from Spec (ONE-TO-ONE immutable mapping)
  
  Required:
    <name>              Feature name (positional)
    -p, --project <id>  Project ID
    --capability <id>   Capability ID from Brief
    --spec-id <id>      Spec ID (ONE-TO-ONE mapping to feature)
    -g, --goal <text>   Goal (min 300 chars, include user flow & requirements)
  
  Optional:
    --scope <scope>     frontend|backend|fullstack|cli|desktop|android|ios|etc
    --priority <P0-P5>  Priority (default: from config)
    --depends <ids>     Pipe-separated feature IDs this depends on
  
  Example:
    mandor feature create "JWT Authentication" -p api \
      --capability jwt-auth --spec-id jwt-auth-spec \
      --scope backend --priority P0 \
      -g "Implement JWT-based authentication with login and token refresh flows"

▶ mandor feature detail <feature-id>
  Show feature details
  
  Flags:
    --json    Machine-readable output
  
  Example:
    mandor feature detail jwt-auth-feature --json

▶ mandor feature update <feature-id> [FLAGS]
  Update feature or change status
  
  Flags:
    --name <text>      New name
    --goal <text>      New goal
    --scope <scope>    New scope
    --priority <P0-P5> New priority
    --status <status>  Status: draft|active|done|cancelled
  
  Example:
    mandor feature update jwt-auth-feature --status done

▶ mandor feature delete <feature-id> [-y]
  Delete feature
  
  Example:
    mandor feature delete jwt-auth-feature -y

════════════════════════════════════════════════════════════════════════
 7. TRACK COMMANDS (Real-Time Visibility)
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
  
  Examples:
    mandor track
    mandor track --json

▶ mandor track project <project-id>
  Show project features and issues
  
  Flags:
    --json, --csv, --tree, --graph, --verbose, --group-by
  
  Example:
    mandor track project api

▶ mandor track feature <feature-id>
  Show feature with all tasks
  
  Example:
    mandor track feature jwt-auth-feature --verbose

▶ mandor track task <task-id>
  Show single task details (including gate status)
  
  Flags:
    --json, --verbose
  
  Example:
    mandor track task jwt-auth-feature-task-a7K2

▶ mandor track issue <issue-id>
  Show single issue details
  
  Example:
    mandor track issue api-issue-0001

════════════════════════════════════════════════════════════════════════
 8. TASK MANAGEMENT
════════════════════════════════════════════════════════════════════════

▶ mandor task create <feature-id> <name> --spec-id <spec-id> \
                     --iae-scenarios <req-X:scenario-Y>|... \
                     -g <goal> --implementation-steps <steps> \
                     --test-cases <cases> [FLAGS]
  Create task with IAE scenarios from Spec requirements
  
  Required:
    <feature-id>                Feature ID (positional)
    <name>                      Task name (positional)
    --spec-id <id>              Spec ID (must match Feature's spec_id)
    --iae-scenarios <text>      IAE scenarios: req-0001:scenario-0001|req-0002:scenario-0001
    -g, --goal <text>           Goal (min 500 chars)
    --implementation-steps <s>  Pipe-separated steps
    --test-cases <cases>        Pipe-separated test cases
  
  Optional:
    --library-needs <libs>      Pipe-separated libraries (or "none")
    --priority <P0-P5>          Priority (default: from config)
    --depends-on <ids>          Pipe-separated task IDs
    -y, --yes                   Skip confirmation
  
  Task Status Rules:
    • New task without depends-on: status = ready, all gates = false
    • New task with depends-on: status = blocked (auto-assigned)
  
  Example:
    mandor task create jwt-auth-feature "Setup JWT Library" \
      --spec-id jwt-auth-spec \
      --iae-scenarios "req-0001:scenario-0001|req-0001:scenario-0002" \
      -g "Setup golang-jwt library with basic configuration and validation" \
      --implementation-steps "Import library|Configure settings|Add validation" \
      --test-cases "JWT creates|JWT validates|JWT expiry" \
      --library-needs "golang-jwt"

▶ mandor task detail <task-id>
  Show task details
  
  Flags:
    --json        Machine-readable output
    --events      Show event history
    --verbose     Show all fields
  
  Example:
    mandor task detail jwt-auth-feature-task-a7K2 --json

▶ mandor task set-gate <task-id> --is-read-brief|--is-read-spec|--is-read-session-notes
  Set read gates to allow in_progress transition
  
  Gates (all three required for in_progress):
    --is-read-brief         Mark brief as read
    --is-read-spec          Mark spec as read
    --is-read-session-notes Mark session notes as read
  
  Usage:
    Gates are ONLY required for ready → in_progress transition.
    Other transitions (ready → cancelled, in_progress → done) work without gates.
    Each flag sets one gate to true.
    Error messages show unmet gates with remediation steps.
  
  Examples:
    mandor task set-gate jwt-auth-feature-task-a7K2 --is-read-brief
    mandor task set-gate jwt-auth-feature-task-a7K2 --is-read-spec
    mandor task set-gate jwt-auth-feature-task-a7K2 --is-read-session-notes

▶ mandor task read-gates <task-id>
  Display task gate status
  
  Example:
    mandor track task jwt-auth-feature-task-a7K2  # Shows gate status

▶ mandor task update <task-id> [FLAGS]
  Update task or change status
  
  Flags:
    --name <text>                   New name
    --goal <text>                   New goal
    --priority <P0-P5>              New priority
    --status <status>               Status: ready|in_progress|done|blocked|cancelled
    --implementation-steps <steps>  Update steps (pipe-separated)
    --test-cases <cases>            Update test cases (pipe-separated)
    --library-needs <libs>          Update libraries (pipe-separated)
    --depends-add <ids>             Add dependencies (pipe-separated)
    --depends-remove <ids>          Remove dependencies (pipe-separated)
    --force                         Force operation (skip checks)
    --dry-run                       Preview changes
    -y, --yes                       Skip confirmation
  
  Gate Enforcement:
    ready → in_progress: REQUIRES all three gates set via set-gate
    ready → cancelled: allowed without gates
    in_progress → done: allowed without gates
    blocked → ready: auto-transitions when dependencies complete
    done: immutable (cannot be changed)
  
  Examples:
    mandor task update jwt-auth-feature-task-a7K2 --status in_progress
    mandor task update jwt-auth-feature-task-a7K2 --status done
    mandor task update jwt-auth-feature-task-a7K2 --status cancelled --reason "Superseded"

════════════════════════════════════════════════════════════════════════
 9. ISSUE MANAGEMENT
════════════════════════════════════════════════════════════════════════

▶ mandor issue create <name> -p <project> -t <type> -g <goal> \
                      --affected-files <files> --affected-tests <tests> \
                      --implementation-steps <steps> [FLAGS]
  Create issue in project
  
  Required:
    <name>                          Issue name (positional)
    -p, --project <id>              Project ID
    -t, --type <type>               Type: bug|improvement|debt|security|performance
    -g, --goal <text>               Goal (min 200 chars)
    --affected-files <files>        Pipe-separated file paths
    --affected-tests <tests>        Pipe-separated test files
    --implementation-steps <steps>  Pipe-separated steps
  
  Optional:
    --priority <P0-P5>              Priority (default: from config)
    --depends-on <ids>              Pipe-separated issue IDs
    --library-needs <libs>          Pipe-separated libraries
    -y, --yes                       Skip confirmation
  
  Example:
    mandor issue create "Memory leak in auth handler" -p api -t bug \
      --priority P0 \
      -g "Goroutine not cleaned up in token refresh handler" \
      --affected-files "src/handlers/auth.go|src/middleware/auth.go" \
      --affected-tests "src/handlers/auth_test.go" \
      --implementation-steps "Identify|Fix|Add tests|Verify" \
      --library-needs "none"

▶ mandor issue detail <issue-id>
  Show issue details
  
  Flags:
    --json        Machine-readable output
    --events      Show event history
  
  Example:
    mandor issue detail api-issue-0001 --json

▶ mandor issue update <issue-id> [FLAGS]
  Update issue or change status
  
  Flags:
    --name <text>                   New name
    --goal <text>                   New goal
    --type <type>                   New type
    --priority <P0-P5>              New priority
    --status <status>               Status: open|ready|in_progress|resolved|wontfix|blocked|cancelled
    --reason <text>                 Reason for status change
    --affected-files <files>        Pipe-separated files
    --affected-tests <tests>        Pipe-separated tests
    --implementation-steps <steps>  Pipe-separated steps
    --library-needs <libs>          Pipe-separated libraries
    --depends-add <ids>             Add dependencies (pipe-separated)
    --depends-remove <ids>          Remove dependencies (pipe-separated)
    --force                         Force operation
    --dry-run                       Preview changes
    -y, --yes                       Skip confirmation
  
  Examples:
    mandor issue update api-issue-0001 --status resolved
    mandor issue update api-issue-0001 --status wontfix --reason "Working as intended"

════════════════════════════════════════════════════════════════════════
 10. GATE ENFORCEMENT SYSTEM
════════════════════════════════════════════════════════════════════════

Gate System Overview:
  • Gates enforce that you've read Brief, Spec, and Session Notes before starting work
  • All three gates MUST be true before ready → in_progress transition
  • Gates are ONLY required for in_progress transition
  • New tasks start with all gates = false
  • Error messages show which gates are unmet with solution steps

Gate Types:
  IsReadBrief         Have you read the Brief document?
  IsReadSpec          Have you read the Spec with requirements?
  IsReadSessionNotes  Have you read session notes from previous work?

Setting Gates:
  1. mandor task set-gate <task-id> --is-read-brief
  2. mandor task set-gate <task-id> --is-read-spec
  3. mandor task set-gate <task-id> --is-read-session-notes
  4. mandor task update <task-id> --status in_progress  # Now allowed

Gate Verification:
  mandor track task <task-id>  # Shows gate status in task details

════════════════════════════════════════════════════════════════════════
 11. STATUS TRANSITIONS & DEPENDENCY AUTO-RESOLUTION
════════════════════════════════════════════════════════════════════════

FEATURE STATUS FLOW:
  draft ──→ active ──→ done
    ↓
    blocked (on dependencies)
    cancelled (with reason)

TASK STATUS FLOW WITH GATE ENFORCEMENT:

  Task Creation:
    • New task: status = ready (if no dependencies)
    • New task with depends-on: status = blocked (auto-assigned)

  Allowed Transitions:
    ready ──→ in_progress ──→ done
      ↑           ↓
      └───────────┴─ blocked (when dependencies exist)
    ready ──→ cancelled (no gates required)
    blocked ──→ ready (auto-transition when dependencies done)
    done: immutable (cannot transition from done)

  Gate Enforcement:
    • ready → in_progress: REQUIRES all three gates set
    • ready → cancelled: allowed without gates
    • in_progress → done: allowed without gates
    • Unmet gates: error shows specific remediation steps
    • Gates NOT required for: cancelled, done, or blocked transitions

ISSUE STATUS FLOW:
  open ──→ ready ──→ in_progress ──→ resolved
   ↓        ↓           ↓
   blocked  blocked  (returns to ready when dependency resolves)
   wontfix (with reason)
   cancelled (with reason)

DEPENDENCY AUTO-RESOLUTION:
  • When task marked done: all blocked dependents auto-transition to ready
  • When issue marked resolved: all blocked dependents auto-transition to ready
  • Works across features: dependencies can span different features
  • Cascading unblocking: Task A done → unblocks Task B → Task B done → unblocks Task C

════════════════════════════════════════════════════════════════════════
 12. BEST PRACTICES
════════════════════════════════════════════════════════════════════════

1. FOLLOW THE WORKFLOW SEQUENCE
   Brief → Spec → Blueprint → Feature → Task
   Complete each phase before moving to the next.

2. WRITE COMPREHENSIVE BRIEFS
   Include problem statement, capabilities, technical stack, affected systems.
   This provides context for all downstream specifications.

3. DETAILED SPECS WITH IAE SCENARIOS
   • Intent (what should happen)
   • Action (what the user does)
   • Expectation (what the system does)
   Specs become references for task gates.

4. MEANINGFUL ARCHITECTURE DECISIONS
   Blueprint documents technical decisions with rationale.
   Helps teams understand design constraints.

5. ONE-TO-ONE SPEC MAPPING FOR FEATURES
   Each feature maps to exactly one Spec (immutable).
   This ensures features have clear requirements.

6. LINK TASKS TO SPEC IAE SCENARIOS
   Tasks reference specific requirement-scenario pairs: req-0001:scenario-0001
   This traces implementation back to requirements.

7. GATE ENFORCEMENT DISCIPLINE
   Always read Brief, Spec, and Session Notes before starting work.
   Gates prevent starting work on incomplete context.
   
   Workflow:
     1. Read Brief, Spec, Session Notes
     2. Set three gates
     3. Transition to in_progress
     4. Implement work
     5. Mark done (auto-unblocks dependents)

8. DEPENDENCY MANAGEMENT
   • Create dependent tasks with --depends-on <task-id>
   • Dependent task auto-assigned status=blocked
   • When dependency done: dependent auto-transitions blocked → ready
   • Cross-feature dependencies: Task A (Feature 1) depends on Task B (Feature 2)

9. CONFIGURATION FOR YOUR TEAM
   Set defaults early:
     mandor config set default_priority P2
     mandor config set goal.lengths.task 500
   Rarely change these.

10. DOCUMENT REASONS FOR STATUS CHANGES
    Always explain why you're cancelling or changing status:
    mandor task update task-id --status cancelled --reason "Superseded by task X"

11. USE TRACK FOR VISIBILITY
    Always check status before starting:
    mandor track feature <feature-id>   # See ready, blocked, in_progress tasks
    mandor track task <task-id>         # See gate status and dependencies

12. PIPE SEPARATORS FOR LISTS
    Use | for multiple values:
    --implementation-steps "Step 1|Step 2|Step 3"
    --depends-on "task-1|task-2|task-3"
    --iae-scenarios "req-0001:scenario-0001|req-0002:scenario-0001"

════════════════════════════════════════════════════════════════════════
 13. QUICK WORKFLOWS
════════════════════════════════════════════════════════════════════════

SETUP NEW PROJECT:
   1. mandor init -y
   2. mandor project create api --name "API Service" --goal "..."
   3. mandor brief create -p api --name "Auth" --why "..." \
        --capabilities "jwt:JWT auth|refresh:Token refresh"
   4. mandor spec create -p api --capability jwt --summary "JWT spec" \
        --requirements "Setup:Enable auth:Import lib|Validate:Verify token"
   5. mandor blueprint create -p api --brief brief-id --problem "..." \
        --decisions "Use JWT|HTTP-only cookies"
   6. mandor feature create "JWT Auth" -p api --capability jwt \
        --spec-id spec-id -g "..."
   7. mandor track project api                          # Verify setup

TASK WITH DEPENDENCIES:
   1. Create Task A:
        mandor task create jwt-feature "Setup JWT" --spec-id spec-id \
          --iae-scenarios "req-0001:scenario-0001" \
          -g "..." --implementation-steps "s1|s2" --test-cases "t1|t2"
   
   2. Create Task B (depends on Task A):
        mandor task create jwt-feature "Validate Tokens" --spec-id spec-id \
          --iae-scenarios "req-0002:scenario-0001" \
          --depends-on jwt-feature-task-a7K2 \
          -g "..." --implementation-steps "s1|s2" --test-cases "t1|t2"
   
   3. Task B status = blocked (waiting for Task A)
   4. mandor track feature jwt-feature  # See status changes
   5. When Task A done: Task B auto-transitions blocked → ready

GATE ENFORCEMENT WORKFLOW:
   1. mandor track task jwt-feature-task-a7K2          # Verify status=ready
   2. mandor task set-gate jwt-feature-task-a7K2 --is-read-brief
   3. mandor task set-gate jwt-feature-task-a7K2 --is-read-spec
   4. mandor task set-gate jwt-feature-task-a7K2 --is-read-session-notes
   5. mandor task update jwt-feature-task-a7K2 --status in_progress  # Now allowed
   6. Work on implementation...
   7. mandor task update jwt-feature-task-a7K2 --status done  # Auto-unblocks Task B
   8. Task B now: status = ready
   9. Repeat gates for Task B, set it to in_progress

TRACK PROGRESS ACROSS WORKFLOW:
   1. mandor track                                       # Workspace overview
   2. mandor track project api                          # Project features
   3. mandor track feature jwt-feature                  # Feature tasks
   4. mandor track task jwt-feature-task-a7K2           # Task with gate status

════════════════════════════════════════════════════════════════════════
 SESSION MANAGEMENT (AI Agent Progress Tracking)
════════════════════════════════════════════════════════════════════════

▶ mandor session note [text]
  Add a timestamped note about work completed or in progress
  
  Flags:
    -r, --read              Read recent notes instead of adding
    -o, --offset <count>    Number of notes to show (default: 50)
  
  Examples:
    mandor session note "Completed JWT setup and testing"
    mandor session note "Started token validation - blocked on spec review"
    mandor session note --read                 # Show last 50 notes
    mandor session note --read --offset 100    # Show last 100 notes

════════════════════════════════════════════════════════════════════════
 AI COMMANDS
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

For more information:
  • GitHub: https://github.com/budisantoso/mandor
  • Run: mandor [command] --help for detailed flag information
  • Configuration: mandor config list to see all available settings

Built for AI Agent Workflows
`)
	return nil
}
