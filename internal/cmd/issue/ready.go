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
	readyBacklogID string
	readyType      string
	readyPriority  string
	readyJSON      bool
)

func NewReadyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:        "ready [--backlog <id>] [--type <type>] [--priority <priority>] [--json]",
		Short:      "List ready issues",
		Deprecated: "Use 'mandor track backlog <id>' instead. All issues are shown; use '--group-by status' to filter by status.",
		Long:       "List all issues with status='ready' that are available to work on (no blocking dependencies).\n\nDEPRECATED: Use 'mandor track backlog <id>' instead to view all issues, optionally with '--group-by status' to organize by status.",
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := service.NewIssueService()
			if err != nil {
				return err
			}

			if !svc.WorkspaceInitialized() {
				return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
			}

			backlogID := readyBacklogID
			if backlogID == "" {
				ws, err := svc.GetWorkspace()
				if err != nil {
					return domain.NewValidationError("No backlog specified and no default backlog set.")
				}
				backlogID = ws.Config.DefaultBacklog
				if backlogID == "" {
					return domain.NewValidationError("No backlog specified and no default backlog set.")
				}
			}

			if !svc.BacklogExists(backlogID) {
				return domain.NewValidationError("Backlog not found: " + backlogID)
			}

			if readyType != "" && !domain.ValidateIssueType(readyType) {
				return domain.NewValidationError("Invalid issue type. Valid types: bug, improvement, debt, security, performance")
			}

			if readyPriority != "" && !domain.ValidatePriority(readyPriority) {
				return domain.NewValidationError("Invalid priority. Valid options: P0, P1, P2, P3, P4, P5")
			}

			input := &domain.IssueListInput{
				BacklogID:      backlogID,
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
				return service.ComparePriority(issues[i].Priority, issues[j].Priority) < 0
			})

			out := cmd.OutOrStdout()

			if readyJSON {
				result := map[string]interface{}{
					" issues": issues,
					"total":   output.Total,
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
				fmt.Fprintf(out, "\nCreate an issue: mandor issue create <name> --backlog %s --type <type>\n", backlogID)
				return nil
			}

			if readyType != "" {
				fmt.Fprintf(out, "Ready issues of type '%s' in backlog %s:\n", readyType, backlogID)
			} else {
				fmt.Fprintf(out, "Ready issues in backlog %s:\n", backlogID)
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

	cmd.Flags().StringVarP(&readyBacklogID, "backlog", "b", "", "Backlog ID filter")
	cmd.Flags().StringVar(&readyType, "type", "", "Filter by issue type (bug, improvement, debt, security, performance)")
	cmd.Flags().StringVar(&readyPriority, "priority", "", "Filter by priority (P0-P5)")
	cmd.Flags().BoolVar(&readyJSON, "json", false, "Output as JSON")

	return cmd
}
