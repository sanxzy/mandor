package task

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
)

var (
	detailJSON         bool
	detailEvents       bool
	detailDependencies bool
	detailTimestamps   bool
)

func NewDetailCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detail <task_id> [--json] [--events] [--dependencies] [--timestamps]",
		Short: "Show task details",
		Long:  "Show detailed information about a specific task.",
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

			input := &domain.TaskDetailInput{
				TaskID:         taskID,
				JSON:           detailJSON,
				IncludeDeleted: false,
				Events:         detailEvents,
				Dependencies:   detailDependencies,
				Timestamps:     detailTimestamps,
			}

			output, err := svc.GetTaskDetail(input)
			if err != nil {
				return err
			}

			if detailJSON {
				out := cmd.OutOrStdout()
				encoder := json.NewEncoder(out)
				encoder.SetIndent("", "  ")
				return encoder.Encode(output)
			}

			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "Task: %s\n", output.ID)
			fmt.Fprintf(out, "  Name:               %s\n", output.Name)
			fmt.Fprintf(out, "  Feature:            %s\n", output.FeatureID)
			fmt.Fprintf(out, "  Project:            %s\n", output.ProjectID)
			fmt.Fprintf(out, "  Status:             %s\n", output.Status)
			fmt.Fprintf(out, "  Priority:           %s\n", output.Priority)
			fmt.Fprintf(out, "  Goal:               %s\n", output.Goal)
			fmt.Fprintf(out, "  Implementation Steps (%d):\n", len(output.ImplementationSteps))
			for i, step := range output.ImplementationSteps {
				fmt.Fprintf(out, "    %d. %s\n", i+1, step)
			}
			fmt.Fprintf(out, "  Test Cases (%d):\n", len(output.TestCases))
			for i, tc := range output.TestCases {
				fmt.Fprintf(out, "    %d. %s\n", i+1, tc)
			}
			fmt.Fprintf(out, "  Derivable Files (%d):\n", len(output.DerivableFiles))
			for _, f := range output.DerivableFiles {
				fmt.Fprintf(out, "    - %s\n", f)
			}
			fmt.Fprintf(out, "  Library Needs (%d):\n", len(output.LibraryNeeds))
			for _, lib := range output.LibraryNeeds {
				fmt.Fprintf(out, "    - %s\n", lib)
			}
			if len(output.DependsOn) > 0 {
				fmt.Fprintf(out, "  Depends On (%d):\n", len(output.DependsOn))
				for _, dep := range output.DependsOn {
					fmt.Fprintf(out, "    - %s\n", dep)
				}
			}
			fmt.Fprintf(out, "  Created:   %s\n", output.CreatedAt)
			fmt.Fprintf(out, "  Updated:   %s\n", output.UpdatedAt)
			fmt.Fprintf(out, "  CreatedBy: %s\n", output.CreatedBy)
			fmt.Fprintf(out, "  UpdatedBy: %s\n", output.UpdatedBy)
			fmt.Fprintf(out, "  Events:    %d\n", output.Events)

			return nil
		},
	}

	cmd.Flags().BoolVar(&detailJSON, "json", false, "Output as JSON")
	cmd.Flags().BoolVar(&detailEvents, "events", false, "Include event history")
	cmd.Flags().BoolVar(&detailDependencies, "dependencies", false, "Include dependency information")
	cmd.Flags().BoolVar(&detailTimestamps, "timestamps", false, "Show formatted timestamps")

	return cmd
}
