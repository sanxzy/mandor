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
	blockedProjectID string
	blockedType      string
	blockedPriority  string
	blockedJSON      bool
)

func NewBlockedCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "blocked [--project <id>] [--type <type>] [--priority <priority>] [--json]",
		Short: "List blocked issues",
		Long:  "List all issues with status='blocked' that are waiting for dependencies to complete.",
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := service.NewIssueService()
			if err != nil {
				return err
			}

			if !svc.WorkspaceInitialized() {
				return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
			}

			projectID := blockedProjectID
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

			if blockedType != "" && !domain.ValidateIssueType(blockedType) {
				return domain.NewValidationError("Invalid issue type. Valid types: bug, improvement, debt, security, performance")
			}

			if blockedPriority != "" && !domain.ValidatePriority(blockedPriority) {
				return domain.NewValidationError("Invalid priority. Valid options: P0, P1, P2, P3, P4, P5")
			}

			input := &domain.IssueListInput{
				ProjectID:      projectID,
				IssueType:      blockedType,
				Status:         domain.IssueStatusBlocked,
				Priority:       blockedPriority,
				IncludeDeleted: false,
				JSON:           blockedJSON,
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

			if blockedJSON {
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
				fmt.Fprintln(out, "No blocked issues found.")
				fmt.Fprintf(out, "\nUnblock issues: mandor issue update <id> --depends <id,...>\n")
				return nil
			}

			if blockedType != "" {
				fmt.Fprintf(out, "Blocked issues of type '%s' in project %s:\n", blockedType, projectID)
			} else {
				fmt.Fprintf(out, "Blocked issues in project %s:\n", projectID)
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

	cmd.Flags().StringVarP(&blockedProjectID, "project", "p", "", "Project ID filter")
	cmd.Flags().StringVar(&blockedType, "type", "", "Filter by issue type (bug, improvement, debt, security, performance)")
	cmd.Flags().StringVar(&blockedPriority, "priority", "", "Filter by priority (P0-P5)")
	cmd.Flags().BoolVar(&blockedJSON, "json", false, "Output as JSON")

	return cmd
}
