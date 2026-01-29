package workspace

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"mandor/internal/service"
)

// NewStatusCmd creates the status command
func NewStatusCmd() *cobra.Command {
	var (
		projectID      string
		summaryFormat  bool
		jsonFormat     bool
		includeDeleted bool
	)

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Display workspace and project status",
		Long: `Display the current status of the Mandor workspace and all projects.

Shows comprehensive statistics including entity counts, status breakdown,
priority distribution, dependency information, and timeline metrics.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			statusSvc, err := service.NewStatusService()
			if err != nil {
				return err
			}

			status, err := statusSvc.GetWorkspaceStatus(projectID)
			if err != nil {
				return err
			}

			// Output formatting
			if jsonFormat {
				return outputStatusJSON(status)
			}

			if summaryFormat {
				return outputStatusSummary(status)
			}

			return outputStatusDefault(status)
		},
	}

	cmd.Flags().StringVarP(&projectID, "project", "p", "", "Show status for specific project only")
	cmd.Flags().BoolVarP(&summaryFormat, "summary", "s", false, "Show summary view (one line per entity type)")
	cmd.Flags().BoolVarP(&jsonFormat, "json", "j", false, "Output in JSON format (machine-readable)")
	cmd.Flags().BoolVarP(&includeDeleted, "include-deleted", "", false, "Include deleted projects in output")

	return cmd
}

func outputStatusDefault(status *service.WorkspaceStatus) error {
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║ MANDOR WORKSPACE STATUS                                    ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	fmt.Printf("Workspace: %s\n", status.Workspace.Name)
	fmt.Printf("ID: %s\n", status.Workspace.ID)
	fmt.Printf("Created: %s by %s\n", status.Workspace.CreatedAt.Format("2006-01-02T15:04:05Z"), status.Workspace.CreatedBy)
	fmt.Printf("Updated: %s\n", status.Workspace.LastUpdatedAt.Format("2006-01-02T15:04:05Z"))
	fmt.Println()

	// Project summary
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Printf("║ PROJECT SUMMARY (%d projects)\n", len(status.Projects))
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	for i, project := range status.Projects {
		fmt.Printf("[%d] %s\n", i+1, project.ID)
		if project.Name != "" {
			fmt.Printf("    Name: %s\n", project.Name)
		}
		fmt.Printf("    Features: %d total\n", project.Stats.Features.Total)
		fmt.Printf("    Tasks: %d total\n", project.Stats.Tasks.Total)
		fmt.Printf("    Issues: %d total\n", project.Stats.Issues.Total)
		fmt.Println()
	}

	// Dependency summary
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║ DEPENDENCY SUMMARY                                         ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	fmt.Printf("Cross-project dependencies: %d\n", status.Dependencies.CrossProjectCount)
	fmt.Printf("Circular dependencies: %d ✓\n", status.Dependencies.CircularDeps)
	fmt.Println()

	// Total stats
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║ WORKSPACE STATS                                            ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	totalEntities := status.Totals.Features + status.Totals.Tasks + status.Totals.Issues
	fmt.Printf("Total Entities: %d\n", totalEntities)
	fmt.Printf("  - Features: %d\n", status.Totals.Features)
	fmt.Printf("  - Tasks: %d\n", status.Totals.Tasks)
	fmt.Printf("  - Issues: %d\n", status.Totals.Issues)
	fmt.Println()

	return nil
}

func outputStatusSummary(status *service.WorkspaceStatus) error {
	fmt.Printf("Workspace: %s\n", status.Workspace.Name)
	fmt.Printf("Projects: %d\n", len(status.Projects))

	for _, project := range status.Projects {
		fmt.Printf("  %s: %dF | %dT | %dI\n",
			project.ID,
			project.Stats.Features.Total,
			project.Stats.Tasks.Total,
			project.Stats.Issues.Total,
		)
	}

	totalFeatures := 0
	totalTasks := 0
	totalIssues := 0
	for _, project := range status.Projects {
		totalFeatures += project.Stats.Features.Total
		totalTasks += project.Stats.Tasks.Total
		totalIssues += project.Stats.Issues.Total
	}

	fmt.Printf("Total: %dF | %dT | %dI | Blocked: %d | Circular deps: %d\n",
		totalFeatures,
		totalTasks,
		totalIssues,
		status.Totals.Blocked,
		status.Dependencies.CircularDeps,
	)

	return nil
}

func outputStatusJSON(status *service.WorkspaceStatus) error {
	data, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprint(os.Stdout, string(data))
	return nil
}
