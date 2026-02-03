package track

import (
	"github.com/spf13/cobra"
)

// NewTrackCmd creates the track command
func NewTrackCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "track [scope] [id]",
		Short: "Track workspace, project, feature, task, or issue",
		Long: `Track provides real-time visibility into task and issue status with unified filtering and output options.

Usage:
  mandor track                                    # Workspace overview
  mandor track workspace                         # Workspace overview (explicit)
  mandor track project <project_id>              # Project issues
  mandor track feature <feature_id>              # Feature tasks
  mandor track task <task_id>                    # Single task details
  mandor track issue <issue_id>                  # Single issue details

Output formats (mutually exclusive):
  --json                                         # Machine-readable JSON
  --csv                                          # CSV export format
  --tree                                         # Tree visualization
  --graph                                        # ASCII graph visualization
  (default)                                      # Human-readable table

Additional options:
  --verbose                                      # Show all fields with details
  --summary                                      # Show aggregated counts only
  --group-by <field>                            # Group by status or priority`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Default to workspace if no scope provided
			scope := "workspace"
			id := ""

			if len(args) > 0 {
				scope = args[0]
			}
			if len(args) > 1 {
				id = args[1]
			}

			// Route to appropriate subcommand
			switch scope {
			case "workspace":
				return handleWorkspace(cmd, id)
			case "project":
				if id == "" {
					return cmd.Help()
				}
				return handleProject(cmd, id)
			case "feature":
				if id == "" {
					return cmd.Help()
				}
				return handleFeature(cmd, id)
			case "task":
				if id == "" {
					return cmd.Help()
				}
				return handleTask(cmd, id)
			case "issue":
				if id == "" {
					return cmd.Help()
				}
				return handleIssue(cmd, id)
			default:
				// Treat first arg as ID, try to resolve scope
				return handleAutoScope(cmd, scope)
			}
		},
	}

	// Add global flags
	cmd.Flags().BoolVar(&globalFlags.JSON, "json", false, "Output as JSON")
	cmd.Flags().BoolVar(&globalFlags.CSV, "csv", false, "Output as CSV")
	cmd.Flags().BoolVar(&globalFlags.Tree, "tree", false, "Output as tree")
	cmd.Flags().BoolVar(&globalFlags.Graph, "graph", false, "Output as ASCII graph")
	cmd.Flags().BoolVar(&globalFlags.Verbose, "verbose", false, "Show all fields with details")
	cmd.Flags().BoolVar(&globalFlags.Summary, "summary", false, "Show aggregated counts only")
	cmd.Flags().StringVar(&globalFlags.GroupBy, "group-by", "", "Group by field (status or priority)")

	return cmd
}
