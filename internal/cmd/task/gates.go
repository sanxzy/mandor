package task

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
)

// NewReadGatesCmd displays the current gate status for a task
func NewReadGatesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read-gates <task-id>",
		Short: "Display task read-gates status",
		Long:  "Show which gates (Brief, Spec, SessionNotes) have been read for this task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := service.NewTaskService()
			if err != nil {
				return err
			}

			if !svc.WorkspaceInitialized() {
				return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
			}

			taskID := args[0]
			
			// Parse task ID to get project and feature
			projectID, featureID, err := svc.ParseTaskID(taskID)
			if err != nil {
				return err
			}

			// Load task from storage
			task, err := svc.ReadTask(projectID, featureID, taskID)
			if err != nil {
				return err
			}

			// Display gate status
			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "✓ Task read-gates: %s\n", task.ID)
			fmt.Fprintf(out, "  Name:                  %s\n", task.Name)
			fmt.Fprintf(out, "  Feature:               %s\n", task.FeatureID)
			fmt.Fprintln(out)
			fmt.Fprintln(out, "Gates Status:")
			fmt.Fprintf(out, "  ✓ is-read-brief:       %v\n", task.ReadGates.IsReadBrief)
			fmt.Fprintf(out, "  ✓ is-read-spec:        %v\n", task.ReadGates.IsReadSpec)
			fmt.Fprintf(out, "  ✓ is-read-session-notes: %v\n", task.ReadGates.IsReadSessionNotes)
			
			// Show remediation steps if any gates are false
			var unmetGates []string
			if !task.ReadGates.IsReadBrief {
				unmetGates = append(unmetGates, "--is-read-brief")
			}
			if !task.ReadGates.IsReadSpec {
				unmetGates = append(unmetGates, "--is-read-spec")
			}
			if !task.ReadGates.IsReadSessionNotes {
				unmetGates = append(unmetGates, "--is-read-session-notes")
			}
			
			if len(unmetGates) > 0 {
				fmt.Fprintln(out)
				fmt.Fprintln(out, "Unmet Gates:")
				fmt.Fprintf(out, "  Set gates with: mandor task set-gate %s %s true\n", taskID, strings.Join(unmetGates, " true "))
			} else {
				fmt.Fprintln(out)
				fmt.Fprintln(out, "All gates satisfied! Ready for in_progress transition.")
			}
			
			return nil
		},
	}

	return cmd
}

// NewSetGateCmd sets a specific gate to true
func NewSetGateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-gate <task-id> --is-read-brief|--is-read-spec|--is-read-session-notes true",
		Short: "Set a task execution gate",
		Long:  "Mark that you have read Brief, Spec, or SessionNotes before starting task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := service.NewTaskService()
			if err != nil {
				return err
			}

			if !svc.WorkspaceInitialized() {
				return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
			}

			taskID := args[0]
			
			// Get which gate to set
			isReadBrief, _ := cmd.Flags().GetBool("is-read-brief")
			isReadSpec, _ := cmd.Flags().GetBool("is-read-spec")
			isReadSessionNotes, _ := cmd.Flags().GetBool("is-read-session-notes")
			
			if !isReadBrief && !isReadSpec && !isReadSessionNotes {
				return domain.NewValidationError("At least one gate flag must be set (--is-read-brief, --is-read-spec, or --is-read-session-notes)")
			}

			// Parse task ID to get project and feature
			projectID, featureID, err := svc.ParseTaskID(taskID)
			if err != nil {
				return err
			}

			// Load task
			task, err := svc.ReadTask(projectID, featureID, taskID)
			if err != nil {
				return err
			}
			
			// Update specified gates
			if isReadBrief {
				task.ReadGates.IsReadBrief = true
			}
			if isReadSpec {
				task.ReadGates.IsReadSpec = true
			}
			if isReadSessionNotes {
				task.ReadGates.IsReadSessionNotes = true
			}
			
			// Save task
			if err := svc.SaveTask(projectID, featureID, task); err != nil {
				return err
			}
			
			var gatesSet []string
			if isReadBrief {
				gatesSet = append(gatesSet, "is_read_brief")
			}
			if isReadSpec {
				gatesSet = append(gatesSet, "is_read_spec")
			}
			if isReadSessionNotes {
				gatesSet = append(gatesSet, "is_read_session_notes")
			}
			
			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "✓ Task gates updated: %s\n", taskID)
			fmt.Fprintf(out, "  Gates set: %s\n", strings.Join(gatesSet, ", "))
			
			// Check if all gates are now satisfied
			if task.ReadGates.IsReadBrief && task.ReadGates.IsReadSpec && task.ReadGates.IsReadSessionNotes {
				fmt.Fprintln(out)
				fmt.Fprintf(out, "  All gates satisfied! Ready to transition to in_progress with: mandor task update %s --status in_progress\n", taskID)
			}

			return nil
		},
	}

	cmd.Flags().Bool("is-read-brief", false, "Mark Brief as read")
	cmd.Flags().Bool("is-read-spec", false, "Mark Spec as read")
	cmd.Flags().Bool("is-read-session-notes", false, "Mark SessionNotes as read")

	return cmd
}

// CheckGatesBeforeInProgress validates that all gates are true before allowing in_progress transition
func CheckGatesBeforeInProgress(task *domain.Task) error {
	var unmetGates []string
	
	if !task.ReadGates.IsReadBrief {
		unmetGates = append(unmetGates, "is_read_brief")
	}
	if !task.ReadGates.IsReadSpec {
		unmetGates = append(unmetGates, "is_read_spec")
	}
	if !task.ReadGates.IsReadSessionNotes {
		unmetGates = append(unmetGates, "is_read_session_notes")
	}
	
	if len(unmetGates) > 0 {
		msg := fmt.Sprintf("Error: Cannot transition to in_progress: %s not satisfied\n", strings.Join(unmetGates, ", "))
		msg += fmt.Sprintf("Solution: Set each gate using 'mandor task set-gate %s --<gate-name> true'\n", task.ID)
		msg += "  Available gates: --is-read-brief, --is-read-spec, --is-read-session-notes"
		return domain.NewValidationError(msg)
	}
	
	return nil
}
