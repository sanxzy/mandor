package issue

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
)

var (
	readyProjectID string
	readyType      string
	readyPriority  string
	readyJSON      bool
)

func NewReadyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ready [--project <id>] [--type <type>] [--priority <priority>] [--json]",
		Short: "List ready issues",
		Long:  "List all issues with status='ready' that are available to work on (no blocking dependencies).",
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := service.NewIssueService()
			if err != nil {
				return err
			}

			if !svc.WorkspaceInitialized() {
				return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
			}

			projectID := readyProjectID
			if projectID == "" {
				ws, err := svc.GetWorkspace()
				if err != nil {
					return domain.NewValidationError("No project specified and no default project set.")
				}
				projectID = ws.Config.DefaultProject
				if projectID == "" {
					return domain.NewValidationError("No project specified and no default project set.")
				}
			}

			if !svc.ProjectExists(projectID) {
				return domain.NewValidationError("Project not found: " + projectID)
			}

			if readyType != "" && !domain.ValidateIssueType(readyType) {
				return domain.NewValidationError("Invalid issue type. Valid types: bug, improvement, debt, security, performance")
			}

			if readyPriority != "" && !domain.ValidatePriority(readyPriority) {
				return domain.NewValidationError("Invalid priority. Valid options: P0, P1, P2, P3, P4, P5")
			}

			input := &domain.IssueListInput{
				ProjectID:      projectID,
				IssueType:      readyType,
				Status:         domain.IssueStatusReady,
				Priority:       readyPriority,
				IncludeDeleted: false,
				JSON:           readyJSON,
				Sort:           "priority",
				Order:          "asc",
			}

			output, err := svc.ListIssues(input)
			if err != nil {
				return err
			}

			issues := output.Issues

			// Sort by priority (P0 first)
			sort.Slice(issues, func(i, j int) bool {
				return comparePriority(issues[i].Priority, issues[j].Priority) < 0
			})

			out := cmd.OutOrStdout()

			if readyJSON {
				result := map[string]interface{}{
					"issues": issues,
					"total":  output.Total,
				}
				jsonBytes, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal JSON: %w", err)
				}
				fmt.Fprintln(out, string(jsonBytes))
				return nil
			}

			if len(issues) == 0 {
				fmt.Fprintln(out, "No ready issues found.")
				fmt.Fprintf(out, "\nCreate an issue: mandor issue create <name> --project %s --type <type>\n", projectID)
				return nil
			}

			if readyType != "" {
				fmt.Fprintf(out, "Ready issues of type '%s' in project %s:\n", readyType, projectID)
			} else {
				fmt.Fprintf(out, "Ready issues in project %s:\n", projectID)
			}

			fmt.Fprintf(out, "%-24s %-14s %-8s %s\n", "ID", "TYPE", "PRIORITY", "NAME")
			fmt.Fprintln(out, strings.Repeat("-", 80))

			for _, i := range issues {
				name := i.Name
				if len(name) > 30 {
					name = name[:27] + "..."
				}
				fmt.Fprintf(out, "%-24s %-14s %-8s %s\n", i.ID, i.IssueType, i.Priority, name)
			}

			fmt.Fprintf(out, "\nTotal: %d\n", len(issues))

			return nil
		},
	}

	cmd.Flags().StringVarP(&readyProjectID, "project", "p", "", "Project ID filter")
	cmd.Flags().StringVar(&readyType, "type", "", "Filter by issue type (bug, improvement, debt, security, performance)")
	cmd.Flags().StringVar(&readyPriority, "priority", "", "Filter by priority (P0-P5)")
	cmd.Flags().BoolVar(&readyJSON, "json", false, "Output as JSON")

	return cmd
}
