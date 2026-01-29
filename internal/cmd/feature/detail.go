package feature

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
)

var (
	detailProjectID string
	detailJSON      bool
)

func NewDetailCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detail <feature_id> [--project <id>]",
		Short: "Show feature details",
		Long:  "Show detailed information about a specific feature.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := service.NewFeatureService()
			if err != nil {
				return err
			}

			if !svc.WorkspaceInitialized() {
				return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
			}

			projectID := detailProjectID
			if projectID == "" {
				return domain.NewValidationError("Project ID is required (--project).")
			}

			featureID := args[0]

			input := &domain.FeatureDetailInput{
				ProjectID:      projectID,
				FeatureID:      featureID,
				JSON:           detailJSON,
				IncludeDeleted: false,
			}

			output, err := svc.GetFeatureDetail(input)
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
			fmt.Fprintf(out, "Feature: %s\n", output.ID)
			fmt.Fprintf(out, "  Name:      %s\n", output.Name)
			fmt.Fprintf(out, "  Project:   %s\n", output.ProjectID)
			fmt.Fprintf(out, "  Goal:      %s\n", output.Goal)
			fmt.Fprintf(out, "  Scope:     %s\n", output.Scope)
			fmt.Fprintf(out, "  Priority:  %s\n", output.Priority)
			fmt.Fprintf(out, "  Status:    %s\n", output.Status)
			fmt.Fprintf(out, "  DependsOn: %v\n", output.DependsOn)
			fmt.Fprintf(out, "  Reason:    %s\n", output.Reason)
			fmt.Fprintf(out, "  Events:    %d\n", output.Events)
			fmt.Fprintf(out, "  Created:   %s\n", output.CreatedAt)
			fmt.Fprintf(out, "  Updated:   %s\n", output.UpdatedAt)
			fmt.Fprintf(out, "  CreatedBy: %s\n", output.CreatedBy)
			fmt.Fprintf(out, "  UpdatedBy: %s\n", output.UpdatedBy)

			return nil
		},
	}

	cmd.Flags().StringVarP(&detailProjectID, "project", "p", "", "Project ID (required)")
	cmd.Flags().BoolVar(&detailJSON, "json", false, "Output as JSON")

	return cmd
}
