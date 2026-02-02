package issue

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
)

func NewSummaryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "summary <project_id>",
		Short: "Display issue summary for a project",
		Long:  "Display a summary of issues in a project grouped by status with priority and type information.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectID := args[0]

			svc, err := service.NewIssueService()
			if err != nil {
				return err
			}

			if !svc.WorkspaceInitialized() {
				return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
			}

			// Group issues by status
			issuesByStatus := make(map[string][]domain.IssueListItem)
			statusOrder := []string{
				domain.IssueStatusOpen,
				domain.IssueStatusReady,
				domain.IssueStatusInProgress,
				domain.IssueStatusBlocked,
				domain.IssueStatusResolved,
				domain.IssueStatusWontFix,
				domain.IssueStatusCancelled,
			}

			input := &domain.IssueListInput{
				ProjectID: projectID,
			}

			output, err := svc.ListIssues(input)
			if err != nil {
				return err
			}

			if output == nil || output.Total == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No issues in project %s.\n", projectID)
				return nil
			}

			// Group by status
			for _, issue := range output.Issues {
				issuesByStatus[issue.Status] = append(issuesByStatus[issue.Status], issue)
			}

			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "Issues in %s:\n\n", projectID)

			// Print by status
			for _, status := range statusOrder {
				issues := issuesByStatus[status]
				if len(issues) == 0 {
					continue
				}

				statusLabel := getIssueStatusLabel(status)
				fmt.Fprintf(out, "%s (%d issues)\n", statusLabel, len(issues))
				fmt.Fprintf(out, "| # | Issue ID | Name | Type | Priority | Status |\n")
				fmt.Fprintf(out, "|---|----------|------|------|----------|--------|\n")

				sort.Slice(issues, func(i, j int) bool {
					return issues[i].CreatedAt < issues[j].CreatedAt
				})

				for idx, issue := range issues {
					name := issue.Name
					if len(name) > 35 {
						name = name[:32] + "..."
					}
					fmt.Fprintf(out, "| %d | %s | %s | %s | %s | %s |\n",
						idx+1,
						issue.ID,
						name,
						issue.IssueType,
						issue.Priority,
						issue.Status)
				}

				fmt.Fprintln(out)
			}

			return nil
		},
	}

	return cmd
}

func getIssueStatusLabel(status string) string {
	labels := map[string]string{
		domain.IssueStatusOpen:       "Open",
		domain.IssueStatusReady:      "Ready",
		domain.IssueStatusInProgress: "In Progress",
		domain.IssueStatusBlocked:    "Blocked",
		domain.IssueStatusResolved:   "Resolved",
		domain.IssueStatusWontFix:    "Won't Fix",
		domain.IssueStatusCancelled:  "Cancelled",
	}
	if label, ok := labels[status]; ok {
		return label
	}
	return "Unknown"
}
