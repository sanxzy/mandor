package issue

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
)

var (
	updateProjectID     string
	updateName          string
	updateGoal          string
	updateType          string
	updatePriority      string
	updateStatus        string
	updateReason        string
	updateDependsOn     string
	updateDependsAdd    string
	updateDependsRemove string
	updateAffectedFiles string
	updateAffectedTests string
	updateImplSteps     string
	updateLibraries     string
	updateStart         bool
	updateResolve       bool
	updateWontFix       bool
	updateReopen        bool
	updateCancel        bool
	updateForce         bool
	updateDryRun        bool
)

func NewUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <issue_id> [--name <text>] [--goal <text>] [--type <type>] [--priority <P0-P5>] [--status <status>] [--reason <text>] [--depends-on <ids>] [--affected-files <files>] [--affected-tests <tests>] [--implementation-steps <steps>] [--library-needs <libs>] [--start] [--resolve] [--wontfix] [--reopen] [--cancel] [--force] [--dry-run]",
		Short: "Update an issue",
		Long:  "Update an issue's metadata, status, or dependencies.",
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

			projectID := updateProjectID
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

			input := &domain.IssueUpdateInput{
				ProjectID: projectID,
				IssueID:   issueID,
			}

			if updateName != "" {
				input.Name = &updateName
			}

			if updateGoal != "" {
				input.Goal = &updateGoal
			}

			if updateType != "" {
				input.IssueType = &updateType
			}

			if updatePriority != "" {
				input.Priority = &updatePriority
			}

			if updateStatus != "" {
				input.Status = &updateStatus
			}

			if updateReason != "" {
				input.Reason = &updateReason
			}

			if updateDependsOn != "" {
				deps := splitByPipe(updateDependsOn)
				input.DependsOn = &deps
			}

			if updateDependsAdd != "" {
				deps := splitByPipe(updateDependsAdd)
				input.DependsAdd = &deps
			}

			if updateDependsRemove != "" {
				deps := splitByPipe(updateDependsRemove)
				input.DependsRemove = &deps
			}

			if updateAffectedFiles != "" {
				files := splitByPipe(updateAffectedFiles)
				input.AffectedFiles = &files
			}

			if updateAffectedTests != "" {
				tests := splitByPipe(updateAffectedTests)
				input.AffectedTests = &tests
			}

			if updateImplSteps != "" {
				steps := splitByPipe(updateImplSteps)
				input.ImplementationSteps = &steps
			}

			if updateLibraries != "" {
				libs := splitByPipe(updateLibraries)
				input.LibraryNeeds = &libs
			}

			input.Start = updateStart
			input.Resolve = updateResolve
			input.WontFix = updateWontFix
			input.Reopen = updateReopen
			input.Cancel = updateCancel
			input.Force = updateForce
			input.DryRun = updateDryRun

			if err := svc.ValidateUpdateInput(input); err != nil {
				return err
			}

			changes, err := svc.UpdateIssue(input)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()

			if updateDryRun {
				if len(changes) == 0 {
					fmt.Fprintf(out, "[DRY RUN] No changes to issue: %s\n", issueID)
				} else {
					fmt.Fprintf(out, "[DRY RUN] Would update issue: %s\n", issueID)
					fmt.Fprintf(out, "  Changes: %s\n", strings.Join(changes, ", "))
				}
				return nil
			}

			fmt.Fprintf(out, "Issue updated: %s\n", issueID)
			if len(changes) > 0 {
				fmt.Fprintf(out, "  Changes: %s\n", strings.Join(changes, ", "))
			}

			// Show issue detail unless status is "resolved"
			if input.Status == nil || *input.Status != domain.IssueStatusResolved {
				fmt.Fprintln(out)
				detailInput := &domain.IssueDetailInput{
					ProjectID:      projectID,
					IssueID:        issueID,
					JSON:           false,
					IncludeDeleted: false,
					Events:         false,
					Timestamps:     false,
				}

				detailOutput, err := svc.GetIssueDetail(detailInput)
				if err == nil {
					fmt.Fprintf(out, "ISSUE: %s", detailOutput.ID)
					if detailOutput.Status == domain.IssueStatusCancelled {
						fmt.Fprint(out, " [CANCELLED]")
					}
					fmt.Fprintln(out)
					fmt.Fprintln(out, strings.Repeat("-", 60))

					fmt.Fprintf(out, "  Name:        %s\n", detailOutput.Name)
					fmt.Fprintf(out, "  Type:        %s\n", detailOutput.IssueType)
					fmt.Fprintf(out, "  Priority:    %s\n", detailOutput.Priority)
					fmt.Fprintf(out, "  Status:      %s\n", detailOutput.Status)
					fmt.Fprintf(out, "  Project:     %s\n", detailOutput.ProjectID)

					if detailOutput.Goal != "" {
						fmt.Fprintf(out, "\n  Goal:        %s\n", detailOutput.Goal)
					}

					fmt.Fprintf(out, "\n  Depends on:  %d issue(s)\n", len(detailOutput.DependsOn))
					for _, depID := range detailOutput.DependsOn {
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

					fmt.Fprintf(out, "\n  Affected Files:      %d\n", len(detailOutput.AffectedFiles))
					for _, f := range detailOutput.AffectedFiles {
						fmt.Fprintf(out, "    - %s\n", f)
					}

					fmt.Fprintf(out, "\n  Affected Tests:      %d\n", len(detailOutput.AffectedTests))
					for _, t := range detailOutput.AffectedTests {
						fmt.Fprintf(out, "    - %s\n", t)
					}

					fmt.Fprintf(out, "\n  Implementation Steps: %d\n", len(detailOutput.ImplementationSteps))
					for i, step := range detailOutput.ImplementationSteps {
						fmt.Fprintf(out, "    %d. %s\n", i+1, step)
					}

					if len(detailOutput.LibraryNeeds) > 0 {
						fmt.Fprintf(out, "\n  Library Needs:       %d\n", len(detailOutput.LibraryNeeds))
						for _, lib := range detailOutput.LibraryNeeds {
							fmt.Fprintf(out, "    - %s\n", lib)
						}
					}

					fmt.Fprintf(out, "\n  Created:     %s by %s\n", detailOutput.CreatedAt, detailOutput.CreatedBy)
					fmt.Fprintf(out, "  Updated:     %s by %s\n", detailOutput.LastUpdatedAt, detailOutput.LastUpdatedBy)

					if detailOutput.Reason != "" && detailOutput.Status == domain.IssueStatusCancelled {
						fmt.Fprintf(out, "\n  Cancellation Reason: %s\n", detailOutput.Reason)
					}

					fmt.Fprintf(out, "\n  Events:      %d\n", detailOutput.Events)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&updateProjectID, "project", "p", "", "Project ID (optional, extracted from issue ID)")
	cmd.Flags().StringVar(&updateName, "name", "", "Update issue name")
	cmd.Flags().StringVar(&updateGoal, "goal", "", "Update issue goal")
	cmd.Flags().StringVar(&updateType, "type", "", "Update issue type (bug/improvement/debt/security/performance)")
	cmd.Flags().StringVar(&updatePriority, "priority", "", "Update priority (P0-P5)")
	cmd.Flags().StringVar(&updateStatus, "status", "", "Set status directly")
	cmd.Flags().StringVar(&updateReason, "reason", "", "Reason for status change")
	cmd.Flags().StringVar(&updateDependsOn, "depends-on", "", "Replace dependencies")
	cmd.Flags().StringVar(&updateDependsAdd, "depends-add", "", "Add dependencies")
	cmd.Flags().StringVar(&updateDependsRemove, "depends-remove", "", "Remove dependencies")
	cmd.Flags().StringVar(&updateAffectedFiles, "affected-files", "", "Replace affected files")
	cmd.Flags().StringVar(&updateAffectedTests, "affected-tests", "", "Replace affected tests")
	cmd.Flags().StringVar(&updateImplSteps, "implementation-steps", "", "Replace implementation steps")
	cmd.Flags().StringVar(&updateLibraries, "library-needs", "", "Replace library needs")
	cmd.Flags().BoolVar(&updateStart, "start", false, "Start working (open/ready → in_progress)")
	cmd.Flags().BoolVar(&updateResolve, "resolve", false, "Mark as resolved")
	cmd.Flags().BoolVar(&updateWontFix, "wontfix", false, "Mark as wontfix")
	cmd.Flags().BoolVar(&updateReopen, "reopen", false, "Reopen resolved/wontfix issue")
	cmd.Flags().BoolVar(&updateCancel, "cancel", false, "Cancel issue")
	cmd.Flags().BoolVar(&updateForce, "force", false, "Force operation (skip checks)")
	cmd.Flags().BoolVar(&updateDryRun, "dry-run", false, "Show what would change")

	return cmd
}
