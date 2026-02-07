package brief

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
)

var readCmd = &cobra.Command{
	Use:   "read <brief-id>",
	Short: "Read a Brief",
	Long:  "Display the complete Brief document",
	Args:  cobra.ExactArgs(1),
	RunE:  runRead,
}

func runRead(cmd *cobra.Command, args []string) error {
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

	svc, err := service.NewBriefService()
	if err != nil {
		return err
	}

	if !svc.WorkspaceInitialized() {
		return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
	}

	brief, err := svc.ReadBrief(backlogID)
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()

	fmt.Fprintf(out, "BRIEF: %s\n", brief.ID)
	fmt.Fprintln(out, strings.Repeat("-", 60))

	fmt.Fprintf(out, "  Status:       %s\n", brief.Status)
	fmt.Fprintf(out, "  Backlog:     %s\n", brief.BacklogID)

	fmt.Fprintln(out, "\n  Why:")
	for _, line := range strings.Split(brief.Why, "\n") {
		fmt.Fprintf(out, "    %s\n", line)
	}

	fmt.Fprintf(out, "\n  What Changes: %d items\n", len(brief.WhatChanges))
	for _, change := range brief.WhatChanges {
		fmt.Fprintf(out, "    - %s\n", change)
	}

	if len(brief.NewCapabilities) > 0 {
		fmt.Fprintf(out, "\n  New Capabilities: %d\n", len(brief.NewCapabilities))
		for _, cap := range brief.NewCapabilities {
			fmt.Fprintf(out, "    - %s (%s)\n", cap.Name, cap.ID)
		}
	}

	if len(brief.ModifiedCapabilities) > 0 {
		fmt.Fprintf(out, "\n  Modified Capabilities: %d\n", len(brief.ModifiedCapabilities))
		for _, cap := range brief.ModifiedCapabilities {
			fmt.Fprintf(out, "    - %s (%s)\n", cap.Name, cap.ID)
		}
	}

	if len(brief.Impact.TechnicalStack) > 0 {
		fmt.Fprintf(out, "\n  Technical Stack: %d\n", len(brief.Impact.TechnicalStack))
		for _, tech := range brief.Impact.TechnicalStack {
			fmt.Fprintf(out, "    - %s\n", tech)
		}
	}

	if len(brief.Impact.AffectedSystems) > 0 {
		fmt.Fprintf(out, "\n  Affected Systems: %d\n", len(brief.Impact.AffectedSystems))
		for _, sys := range brief.Impact.AffectedSystems {
			fmt.Fprintf(out, "    - %s\n", sys)
		}
	}

	if len(brief.Impact.Dependencies) > 0 {
		fmt.Fprintf(out, "\n  Dependencies: %d\n", len(brief.Impact.Dependencies))
		for _, dep := range brief.Impact.Dependencies {
			fmt.Fprintf(out, "    - %s\n", dep)
		}
	}

	fmt.Fprintf(out, "\n  Created:      %s by %s\n", brief.CreatedAt.Format("2006-01-02 15:04"), brief.CreatedBy)
	fmt.Fprintf(out, "  Updated:      %s by %s\n", brief.UpdatedAt.Format("2006-01-02 15:04"), brief.UpdatedBy)

	return nil
}

func GetReadCmd() *cobra.Command {
	return readCmd
}
