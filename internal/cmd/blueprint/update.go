package blueprint

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
)

var updateCmd = &cobra.Command{
	Use:   "update <blueprint-id>",
	Short: "Update a Blueprint",
	Long:  "Update Blueprint fields (decisions, data models, risks, etc.)",
	Args:  cobra.ExactArgs(1),
	RunE:  runUpdate,
}

var (
	updateProblem        string
	updateDecisions      string
	updateImplementation string
	updateStatus         string
	updateDryRun         bool
)

func init() {
	updateCmd.Flags().StringVar(&updateProblem, "problem", "", "Updated problem statement")
	updateCmd.Flags().StringVar(&updateDecisions, "decisions", "", "Updated decisions")
	updateCmd.Flags().StringVar(&updateImplementation, "implementation", "", "Updated implementation strategy")
	updateCmd.Flags().StringVar(&updateStatus, "status", "", "draft | active | archived")
	updateCmd.Flags().BoolVar(&updateDryRun, "dry-run", false, "Show what would be changed without making changes")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	blueprintID := args[0]

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

	// Load current blueprint
	blueprint, err := svc.ReadBlueprint(backlogID)
	if err != nil {
		return err
	}

	// Track changes
	var changes []string

	// Apply updates
	if updateProblem != "" {
		blueprint.ProblemStatement = updateProblem
		changes = append(changes, "problem")
	}

	if updateStatus != "" {
		blueprint.Status = updateStatus
		changes = append(changes, "status")
	}

	if updateDecisions != "" {
		changes = append(changes, "decisions")
	}

	if updateImplementation != "" {
		blueprint.ImplementationStrategy = updateImplementation
		changes = append(changes, "implementation")
	}

	out := cmd.OutOrStdout()

	if updateDryRun {
		fmt.Fprintf(out, "[DRY RUN] Would update blueprint: %s\n", blueprintID)
		if len(changes) > 0 {
			fmt.Fprintf(out, "  Changes: %s\n", strings.Join(changes, ", "))
		}
		return nil
	}

	// Save the updated blueprint
	if err := svc.UpdateBlueprint(backlogID, blueprint); err != nil {
		return err
	}

	fmt.Fprintf(out, "Blueprint updated: %s\n", blueprintID)
	if len(changes) > 0 {
		fmt.Fprintf(out, "  Changes: %s\n", strings.Join(changes, ", "))
	}

	// Show full blueprint state (following issue update pattern)
	fmt.Fprintln(out)
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
	}

	if len(blueprint.UserTypes) > 0 {
		fmt.Fprintf(out, "\n  User Types: %d\n", len(blueprint.UserTypes))
	}

	fmt.Fprintf(out, "\n  Architecture Decisions: %d\n", len(blueprint.ArchitectureDecisions))
	for i, dec := range blueprint.ArchitectureDecisions {
		fmt.Fprintf(out, "    %d. %s\n", i+1, dec.Title)
	}

	if len(blueprint.DataModels) > 0 {
		fmt.Fprintf(out, "\n  Data Models: %d\n", len(blueprint.DataModels))
	}

	if len(blueprint.Risks) > 0 {
		fmt.Fprintf(out, "\n  Risks: %d\n", len(blueprint.Risks))
	}

	fmt.Fprintf(out, "\n  Created:      %s by %s\n", blueprint.CreatedAt.Format("2006-01-02 15:04"), blueprint.CreatedBy)
	fmt.Fprintf(out, "  Updated:      %s by %s\n", blueprint.UpdatedAt.Format("2006-01-02 15:04"), blueprint.UpdatedBy)

	return nil
}

func GetUpdateCmd() *cobra.Command {
	return updateCmd
}
