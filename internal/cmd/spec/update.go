package spec

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
)

var updateCmd = &cobra.Command{
	Use:   "update <spec-id>",
	Short: "Update a Spec",
	Long:  "Update Spec fields",
	Args:  cobra.ExactArgs(1),
	RunE:  runUpdate,
}

var (
	updateSummary      string
	updateRequirements string
	updateStatus       string
	updateDryRun       bool
)

func init() {
	updateCmd.Flags().StringVar(&updateSummary, "summary", "", "Updated description")
	updateCmd.Flags().StringVar(&updateRequirements, "requirements", "", "Updated requirements")
	updateCmd.Flags().StringVar(&updateStatus, "status", "", "draft | active | archived")
	updateCmd.Flags().BoolVar(&updateDryRun, "dry-run", false, "Show what would be changed without making changes")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	specID := args[0]

	var backlogID string
	for p := cmd; p != nil; p = p.Parent() {
		if val, err := p.Flags().GetString("backlog"); err == nil && val != "" {
			backlogID = val
			break
		}
	}
	if backlogID == "" {
		return domain.NewValidationError("Backlog ID is required (--backlog).")
	}

	svc, err := service.NewSpecService()
	if err != nil {
		return err
	}

	if !svc.WorkspaceInitialized() {
		return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
	}

	// Load current spec
	spec, err := svc.ReadSpec(backlogID, specID)
	if err != nil {
		return err
	}

	// Track changes
	var changes []string

	// Apply updates
	if updateSummary != "" {
		spec.Summary = updateSummary
		changes = append(changes, "summary")
	}

	if updateStatus != "" {
		spec.Status = updateStatus
		changes = append(changes, "status")
	}

	if updateRequirements != "" {
		changes = append(changes, "requirements")
	}

	out := cmd.OutOrStdout()

	if updateDryRun {
		fmt.Fprintf(out, "[DRY RUN] Would update spec: %s\n", specID)
		if len(changes) > 0 {
			fmt.Fprintf(out, "  Changes: %s\n", strings.Join(changes, ", "))
		}
		return nil
	}

	// Save the updated spec
	if err := svc.UpdateSpec(backlogID, spec); err != nil {
		return err
	}

	fmt.Fprintf(out, "Spec updated: %s\n", specID)
	if len(changes) > 0 {
		fmt.Fprintf(out, "  Changes: %s\n", strings.Join(changes, ", "))
	}

	// Show full spec state (following issue update pattern)
	fmt.Fprintln(out)
	fmt.Fprintf(out, "SPEC: %s\n", spec.ID)
	fmt.Fprintln(out, strings.Repeat("-", 60))

	fmt.Fprintf(out, "  Capability:   %s\n", spec.CapabilityID)
	fmt.Fprintf(out, "  Status:       %s\n", spec.Status)
	fmt.Fprintf(out, "  Backlog:     %s\n", spec.BacklogID)

	fmt.Fprintln(out, "\n  Summary:")
	for _, line := range strings.Split(spec.Summary, "\n") {
		fmt.Fprintf(out, "    %s\n", line)
	}

	fmt.Fprintf(out, "\n  Requirements: %d\n", len(spec.Requirements))
	for i, req := range spec.Requirements {
		fmt.Fprintf(out, "\n    %d. %s (%s)\n", i+1, req.Summary, req.ID)
		if req.Details != "" {
			fmt.Fprintf(out, "       Details: %s\n", req.Details)
		}
		if len(req.AcceptanceCriteria) > 0 {
			fmt.Fprintf(out, "       Acceptance Criteria: %d\n", len(req.AcceptanceCriteria))
		}
		fmt.Fprintf(out, "       IAE Scenarios: %d\n", len(req.IAEScenarios))
	}

	fmt.Fprintf(out, "\n  Created:      %s by %s\n", spec.CreatedAt.Format("2006-01-02 15:04"), spec.CreatedBy)
	fmt.Fprintf(out, "  Updated:      %s by %s\n", spec.UpdatedAt.Format("2006-01-02 15:04"), spec.UpdatedBy)

	return nil
}

func GetUpdateCmd() *cobra.Command {
	return updateCmd
}
