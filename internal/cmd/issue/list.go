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
	listProjectID string
	listType      string
	listStatus    string
	listPriority  string
	listJSON      bool
	listSort      string
	listOrder     string
	listVerbose   bool
)

func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [--project <id>] [--type <type>] [--status <status>] [--priority <priority>] [--json] [--sort <field>] [--order <asc|desc>]",
		Short: "List issues",
		Long:  "List issues in the specified project with optional filters.",
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := service.NewIssueService()
			if err != nil {
				return err
			}

			if !svc.WorkspaceInitialized() {
				return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
			}

			projectID := listProjectID
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

			input := &domain.IssueListInput{
				ProjectID:      projectID,
				IssueType:      listType,
				Status:         listStatus,
				Priority:       listPriority,
				IncludeDeleted: false,
				JSON:           listJSON,
				Sort:           listSort,
				Order:          listOrder,
			}

			if listType != "" && !domain.ValidateIssueType(listType) {
				return domain.NewValidationError("Invalid issue type. Valid types: bug, improvement, debt, security, performance")
			}

			if listStatus != "" && !domain.ValidateIssueStatus(listStatus) {
				return domain.NewValidationError("Invalid status. Valid options: open, ready, in_progress, blocked, resolved, wontfix, cancelled")
			}

			if listPriority != "" && !domain.ValidatePriority(listPriority) {
				return domain.NewValidationError("Invalid priority. Valid options: P0, P1, P2, P3, P4, P5")
			}

			output, err := svc.ListIssues(input)
			if err != nil {
				return err
			}

			issues := output.Issues

			if listSort != "" {
				sortField := listSort
				if listSort == "updated_at" {
					sortField = "last_updated_at"
				}
				sort.Slice(issues, func(i, j int) bool {
					switch sortField {
					case "priority":
						return comparePriority(issues[i].Priority, issues[j].Priority) < 0
					case "name":
						return issues[i].Name < issues[j].Name
					case "created_at", "last_updated_at":
						return issues[i].CreatedAt < issues[j].CreatedAt
					default:
						return issues[i].LastUpdatedAt < issues[j].LastUpdatedAt
					}
				})
				if listOrder == "asc" {
					if listSort == "priority" {
						sort.Slice(issues, func(i, j int) bool {
							return comparePriority(issues[i].Priority, issues[j].Priority) > 0
						})
					} else {
						for i, j := 0, len(issues)-1; i < j; i, j = i+1, j-1 {
							issues[i], issues[j] = issues[j], issues[i]
						}
					}
				}
			}

			out := cmd.OutOrStdout()

			if listJSON {
				result := map[string]interface{}{
					"issues":  issues,
					"total":   output.Total,
					"deleted": output.Deleted,
				}
				jsonBytes, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal JSON: %w", err)
				}
				fmt.Fprintln(out, string(jsonBytes))
				return nil
			}

			if len(issues) == 0 {
				fmt.Fprintf(out, "No issues found in project '%s'.\n", projectID)
				fmt.Fprintf(out, "\nRun: mandor issue create <name> --project %s --type <type> --goal <text>\n", projectID)
				return nil
			}

			if listVerbose {
				fmt.Fprintf(out, "%-24s %-14s %-8s %-12s %-10s %-5s %-5s %-5s %s\n",
					"ISSUES", "TYPE", "PRIORITY", "STATUS", "UPDATED", "FILES", "TESTS", "STEPS", "NAME")
				fmt.Fprintf(out, "%s\n", strings.Repeat("-", 100))
				for _, i := range issues {
					updated := i.LastUpdatedAt
					if len(updated) >= 10 {
						updated = updated[:10]
					}
					name := i.Name
					if len(name) > 30 {
						name = name[:27] + "..."
					}
					fmt.Fprintf(out, "%-24s %-14s %-8s %-12s %-10s %-5d %-5d %-5d %s\n",
						i.ID, i.IssueType, i.Priority, i.Status, updated,
						i.AffectedFilesCount, i.AffectedTestsCount, i.ImplementationStepsCount, name)
				}
			} else {
				fmt.Fprintf(out, "%-24s %-14s %-8s %-12s %-10s %-5s %-5s %-5s\n",
					"ISSUES", "TYPE", "PRIORITY", "STATUS", "UPDATED", "FILES", "TESTS", "STEPS")
				fmt.Fprintf(out, "%s\n", strings.Repeat("-", 90))
				for _, i := range issues {
					updated := i.LastUpdatedAt
					if len(updated) >= 10 {
						updated = updated[:10]
					}
					fmt.Fprintf(out, "%-24s %-14s %-8s %-12s %-10s %-5d %-5d %-5d\n",
						i.ID, i.IssueType, i.Priority, i.Status, updated,
						i.AffectedFilesCount, i.AffectedTestsCount, i.ImplementationStepsCount)
				}
			}

			fmt.Fprintf(out, "\nTotal: %d", output.Total)
			statusCounts := make(map[string]int)
			for _, i := range issues {
				statusCounts[i.Status]++
			}
			for _, status := range []string{"open", "ready", "in_progress", "blocked", "resolved", "wontfix", "cancelled"} {
				if count, ok := statusCounts[status]; ok {
					fmt.Fprintf(out, " | %s: %d", status, count)
				}
			}
			fmt.Fprintln(out)

			return nil
		},
	}

	cmd.Flags().StringVarP(&listProjectID, "project", "p", "", "Project ID filter")
	cmd.Flags().StringVar(&listType, "type", "", "Filter by issue type")
	cmd.Flags().StringVar(&listStatus, "status", "", "Filter by status")
	cmd.Flags().StringVar(&listPriority, "priority", "", "Filter by priority (P0-P5)")
	cmd.Flags().BoolVar(&listJSON, "json", false, "Output as JSON")
	cmd.Flags().StringVar(&listSort, "sort", "last_updated_at", "Sort field (created_at, last_updated_at, priority, name)")
	cmd.Flags().StringVar(&listOrder, "order", "desc", "Sort order (asc, desc)")
	cmd.Flags().BoolVar(&listVerbose, "verbose", false, "Show issue names in table output")

	return cmd
}

func comparePriority(a, b string) int {
	order := []string{"P0", "P1", "P2", "P3", "P4", "P5"}
	ai, bj := -1, -1
	for i, p := range order {
		if a == p {
			ai = i
		}
		if b == p {
			bj = i
		}
	}
	if ai < bj {
		return -1
	}
	if ai > bj {
		return 1
	}
	return 0
}
