package spec

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
)

var detailCmd = &cobra.Command{
	Use:   "detail <spec-id>",
	Short: "Display Spec details",
	Long:  "Show complete Spec with all requirements and scenarios",
	Args:  cobra.ExactArgs(1),
	RunE:  runDetail,
}

func runDetail(cmd *cobra.Command, args []string) error {
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

	spec, err := svc.ReadSpec(backlogID, specID)
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()

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
			for _, ac := range req.AcceptanceCriteria {
				fmt.Fprintf(out, "         - %s\n", ac)
			}
		}
		fmt.Fprintf(out, "       IAE Scenarios: %d\n", len(req.IAEScenarios))
		for j, iae := range req.IAEScenarios {
			fmt.Fprintf(out, "         %d. Intent: %s\n", j+1, iae.Intent)
			fmt.Fprintf(out, "            Action: %s\n", iae.Action)
			fmt.Fprintf(out, "            Expect: %s\n", iae.Expect)
		}
	}

	fmt.Fprintf(out, "\n  Created:      %s by %s\n", spec.CreatedAt.Format("2006-01-02 15:04"), spec.CreatedBy)
	fmt.Fprintf(out, "  Updated:      %s by %s\n", spec.UpdatedAt.Format("2006-01-02 15:04"), spec.UpdatedBy)

	return nil
}

func GetDetailCmd() *cobra.Command {
	return detailCmd
}
