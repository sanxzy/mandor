package issue

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
)

var (
	detailProjectID  string
	detailJSON       bool
	detailIncludeDeleted bool
	detailEvents     bool
	detailTimestamps bool
)

func NewDetailCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detail <issue_id> [--project <id>] [--json] [--include-deleted] [--events]",
		Short: "Show issue details",
		Long:  "Show detailed information about an issue.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := service.NewIssueService()
			if err != nil {
				return err
			}

			if !svc.WorkspaceInitialized() {
				return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
			}

			issueID := args[0]

			projectID := detailProjectID
			if projectID == "" {
				parts := strings.Split(issueID, "-issue-")
				if len(parts) < 2 {
					return domain.NewValidationError("Invalid issue ID format. Expected: <project_id>-issue-<nanoid>")
				}
				projectID = parts[0]
			}

			if !svc.ProjectExists(projectID) {
				return domain.NewValidationError("Project not found: " + projectID)
			}

			input := &domain.IssueDetailInput{
				ProjectID:      projectID,
				IssueID:        issueID,
				JSON:           detailJSON,
				IncludeDeleted: detailIncludeDeleted,
				Events:         detailEvents,
				Timestamps:     detailTimestamps,
			}

			output, err := svc.GetIssueDetail(input)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()

			if detailJSON {
				jsonBytes, err := json.MarshalIndent(output, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal JSON: %w", err)
				}
				fmt.Fprintln(out, string(jsonBytes))
				return nil
			}

			fmt.Fprintf(out, "ISSUE: %s", output.ID)
			if output.Status == domain.IssueStatusCancelled {
				fmt.Fprint(out, " [CANCELLED]")
			}
			fmt.Fprintln(out)
			fmt.Fprintln(out, strings.Repeat("-", 60))

			fmt.Fprintf(out, "  Name:        %s\n", output.Name)
			fmt.Fprintf(out, "  Type:        %s\n", output.IssueType)
			fmt.Fprintf(out, "  Priority:    %s\n", output.Priority)
			fmt.Fprintf(out, "  Status:      %s\n", output.Status)
			fmt.Fprintf(out, "  Project:     %s\n", output.ProjectID)

			if output.Goal != "" {
				fmt.Fprintf(out, "\n  Goal:        %s\n", output.Goal)
			}

			fmt.Fprintf(out, "\n  Depends on:  %d issue(s)\n", len(output.DependsOn))
			for _, depID := range output.DependsOn {
				dep, err := svc.ReadDependency(projectID, depID)
				statusIcon := "○"
				if err == nil {
					switch dep.Status {
					case domain.IssueStatusResolved, domain.IssueStatusWontFix:
						statusIcon = "✓"
					case domain.IssueStatusCancelled:
						statusIcon = "✗"
					}
				}
				fmt.Fprintf(out, "    %s %s\n", statusIcon, depID)
			}

			fmt.Fprintf(out, "\n  Affected Files:      %d\n", len(output.AffectedFiles))
			for _, f := range output.AffectedFiles {
				fmt.Fprintf(out, "    - %s\n", f)
			}

			fmt.Fprintf(out, "\n  Affected Tests:      %d\n", len(output.AffectedTests))
			for _, t := range output.AffectedTests {
				fmt.Fprintf(out, "    - %s\n", t)
			}

			fmt.Fprintf(out, "\n  Implementation Steps: %d\n", len(output.ImplementationSteps))
			for i, step := range output.ImplementationSteps {
				fmt.Fprintf(out, "    %d. %s\n", i+1, step)
			}

			if len(output.LibraryNeeds) > 0 {
				fmt.Fprintf(out, "\n  Library Needs:       %d\n", len(output.LibraryNeeds))
				for _, lib := range output.LibraryNeeds {
					fmt.Fprintf(out, "    - %s\n", lib)
				}
			}

			fmt.Fprintf(out, "\n  Created:     %s by %s\n", output.CreatedAt, output.CreatedBy)
			fmt.Fprintf(out, "  Updated:     %s by %s\n", output.LastUpdatedAt, output.LastUpdatedBy)

			if output.Reason != "" && output.Status == domain.IssueStatusCancelled {
				fmt.Fprintf(out, "\n  Cancellation Reason: %s\n", output.Reason)
			}

			if detailEvents {
				fmt.Fprintf(out, "\n  Events:      %d\n", output.Events)
				events, _ := svc.GetIssueEvents(projectID, issueID)
				for _, event := range events {
					fmt.Fprintf(out, "    %s [%s] by %s\n", event.Ts.Format("2006-01-02 15:04:05"), event.Type, event.By)
				}
			} else {
				fmt.Fprintf(out, "\n  Events:      %d\n", output.Events)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&detailProjectID, "project", "p", "", "Project ID (optional, extracted from issue ID)")
	cmd.Flags().BoolVar(&detailJSON, "json", false, "Output as JSON")
	cmd.Flags().BoolVar(&detailIncludeDeleted, "include-deleted", false, "Include cancelled issues")
	cmd.Flags().BoolVar(&detailEvents, "events", false, "Show event history")
	cmd.Flags().BoolVar(&detailTimestamps, "timestamps", false, "Show all timestamps")

	return cmd
}
