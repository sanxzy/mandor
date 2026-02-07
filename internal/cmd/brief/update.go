package brief

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
)

var updateCmd = &cobra.Command{
	Use:   "update <brief-id>",
	Short: "Update a Brief",
	Long:  "Update Brief fields (why, capabilities, impact, etc.)",
	Args:  cobra.ExactArgs(1),
	RunE:  runUpdate,
}

var (
	updateWhy          string
	updateCapabilities string
	updateStatus       string
	updateDryRun       bool
)

func init() {
	updateCmd.Flags().StringVar(&updateWhy, "why", "", "Updated problem statement")
	updateCmd.Flags().StringVar(&updateCapabilities, "capabilities", "", "Updated capabilities")
	updateCmd.Flags().StringVar(&updateStatus, "status", "", "draft | active | archived")
	updateCmd.Flags().BoolVar(&updateDryRun, "dry-run", false, "Show what would be changed without making changes")
}

func runUpdate(cmd *cobra.Command, args []string) error {
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

	// Load current brief
	brief, err := svc.ReadBrief(backlogID)
	if err != nil {
		return err
	}

	// Track changes
	var changes []string

	// Apply updates
	if updateWhy != "" {
		brief.Why = updateWhy
		changes = append(changes, "why")
	}

	if updateStatus != "" {
		brief.Status = updateStatus
		changes = append(changes, "status")
	}

	if updateCapabilities != "" {
		caps, err := parseCapabilities(updateCapabilities)
		if err != nil {
			return err
		}
		brief.NewCapabilities = nil
		brief.ModifiedCapabilities = nil
		for _, cap := range caps {
			capID := strings.ToLower(strings.ReplaceAll(cap.Name, " ", "-"))
			capability := domain.Capability{
				ID:          capID,
				Name:        cap.Name,
				Description: cap.Description,
			}
			if cap.Modified {
				brief.ModifiedCapabilities = append(brief.ModifiedCapabilities, capability)
			} else {
				brief.NewCapabilities = append(brief.NewCapabilities, capability)
			}
		}
		changes = append(changes, "capabilities")
	}

	out := cmd.OutOrStdout()

	if updateDryRun {
		fmt.Fprintf(out, "[DRY RUN] Would update brief: %s\n", backlogID)
		if len(changes) > 0 {
			fmt.Fprintf(out, "  Changes: %s\n", strings.Join(changes, ", "))
		}
		return nil
	}

	// Save the updated brief
	if err := svc.UpdateBrief(backlogID, brief); err != nil {
		return err
	}

	fmt.Fprintf(out, "Brief updated: %s\n", backlogID)
	if len(changes) > 0 {
		fmt.Fprintf(out, "  Changes: %s\n", strings.Join(changes, ", "))
	}

	// Show full brief state (following issue update pattern)
	fmt.Fprintln(out)
	fmt.Fprintf(out, "BRIEF: %s\n", brief.ID)
	fmt.Fprintln(out, strings.Repeat("-", 60))

	fmt.Fprintf(out, "  Status:       %s\n", brief.Status)
	fmt.Fprintf(out, "  Backlog:      %s\n", brief.BacklogID)

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

	fmt.Fprintf(out, "\n  Created:      %s by %s\n", brief.CreatedAt.Format("2006-01-02 15:04"), brief.CreatedBy)
	fmt.Fprintf(out, "  Updated:      %s by %s\n", brief.UpdatedAt.Format("2006-01-02 15:04"), brief.UpdatedBy)

	return nil
}

func GetUpdateCmd() *cobra.Command {
	return updateCmd
}
