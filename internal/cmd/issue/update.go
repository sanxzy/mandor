package issue

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
)

var (
	updateProjectID  string
	updateName       string
	updateGoal       string
	updateType       string
	updatePriority   string
	updateStatus     string
	updateReason     string
	updateDependsOn  string
	updateDependsAdd string
	updateDependsRemove string
	updateAffectedFiles string
	updateAffectedTests string
	updateImplSteps  string
	updateLibraries  string
	updateStart      bool
	updateResolve    bool
	updateWontFix    bool
	updateReopen     bool
	updateCancel     bool
	updateForce      bool
	updateDryRun     bool
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
				deps := splitByComma(updateDependsOn)
				input.DependsOn = &deps
			}

			if updateDependsAdd != "" {
				deps := splitByComma(updateDependsAdd)
				input.DependsAdd = &deps
			}

			if updateDependsRemove != "" {
				deps := splitByComma(updateDependsRemove)
				input.DependsRemove = &deps
			}

			if updateAffectedFiles != "" {
				files := splitByComma(updateAffectedFiles)
				input.AffectedFiles = &files
			}

			if updateAffectedTests != "" {
				tests := splitByComma(updateAffectedTests)
				input.AffectedTests = &tests
			}

			if updateImplSteps != "" {
				steps := splitByComma(updateImplSteps)
				input.ImplementationSteps = &steps
			}

			if updateLibraries != "" {
				libs := splitByComma(updateLibraries)
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
