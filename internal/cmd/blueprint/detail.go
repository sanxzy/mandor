package blueprint

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
)

var detailCmd = &cobra.Command{
	Use:   "detail <blueprint-id>",
	Short: "Display Blueprint details",
	Long:  "Show complete Blueprint with decisions, data models, and risks",
	Args:  cobra.ExactArgs(1),
	RunE:  runDetail,
}

func runDetail(cmd *cobra.Command, args []string) error {
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

	svc, err := service.NewBlueprintService()
	if err != nil {
		return err
	}

	if !svc.WorkspaceInitialized() {
		return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
	}

	blueprint, err := svc.ReadBlueprint(backlogID)
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()

	fmt.Fprintf(out, "BLUEPRINT: %s\n", blueprint.ID)
	fmt.Fprintln(out, strings.Repeat("-", 60))

	fmt.Fprintf(out, "  Version:      %s\n", blueprint.Version)
	fmt.Fprintf(out, "  Status:       %s\n", blueprint.Status)
	fmt.Fprintf(out, "  Backlog:      %s\n", blueprint.BacklogID)

	fmt.Fprintln(out, "\n  Problem Statement:")
	for _, line := range strings.Split(blueprint.ProblemStatement, "\n") {
		fmt.Fprintf(out, "    %s\n", line)
	}

	if len(blueprint.Constraints) > 0 {
		fmt.Fprintf(out, "\n  Constraints: %d\n", len(blueprint.Constraints))
		for _, c := range blueprint.Constraints {
			fmt.Fprintf(out, "    - %s\n", c)
		}
	}

	if len(blueprint.UserTypes) > 0 {
		fmt.Fprintf(out, "\n  User Types: %d\n", len(blueprint.UserTypes))
		for _, u := range blueprint.UserTypes {
			fmt.Fprintf(out, "    - %s\n", u)
		}
	}

	fmt.Fprintf(out, "\n  Architecture Decisions: %d\n", len(blueprint.ArchitectureDecisions))
	for i, dec := range blueprint.ArchitectureDecisions {
		fmt.Fprintf(out, "\n    %d. %s (%s)\n", i+1, dec.Title, dec.ID)
		fmt.Fprintf(out, "       Decision: %s\n", dec.Decision)
		fmt.Fprintf(out, "       Rationale: %s\n", dec.Rationale)
	}

	if len(blueprint.DataModels) > 0 {
		fmt.Fprintf(out, "\n  Data Models: %d\n", len(blueprint.DataModels))
		for _, dm := range blueprint.DataModels {
			fmt.Fprintf(out, "    - %s\n", dm.Name)
		}
	}

	if blueprint.ImplementationStrategy != "" {
		fmt.Fprintln(out, "\n  Implementation Strategy:")
		for _, line := range strings.Split(blueprint.ImplementationStrategy, "\n") {
			fmt.Fprintf(out, "    %s\n", line)
		}
	}

	if len(blueprint.Risks) > 0 {
		fmt.Fprintf(out, "\n  Risks: %d\n", len(blueprint.Risks))
		for i, risk := range blueprint.Risks {
			fmt.Fprintf(out, "    %d. %s (%s)\n", i+1, risk.Description, risk.ID)
			fmt.Fprintf(out, "       Mitigation: %s\n", risk.Mitigation)
		}
	}

	fmt.Fprintf(out, "\n  Created:      %s by %s\n", blueprint.CreatedAt.Format("2006-01-02 15:04"), blueprint.CreatedBy)
	fmt.Fprintf(out, "  Updated:      %s by %s\n", blueprint.UpdatedAt.Format("2006-01-02 15:04"), blueprint.UpdatedBy)

	return nil
}

func GetDetailCmd() *cobra.Command {
	return detailCmd
}
