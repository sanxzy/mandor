package feature

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
)

var (
	listProjectID string
	listJSON      bool
)

func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [--project <id>]",
		Short: "List features",
		Long:  "List all features in the specified project.",
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := service.NewFeatureService()
			if err != nil {
				return err
			}

			if !svc.WorkspaceInitialized() {
				return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
			}

			projectID := listProjectID
			if projectID == "" {
				return domain.NewValidationError("Project ID is required (--project).")
			}

			input := &domain.FeatureListInput{
				ProjectID:      projectID,
				IncludeDeleted: false,
				JSON:           listJSON,
			}

			output, err := svc.ListFeatures(input)
			if err != nil {
				return err
			}

			if listJSON {
				out := cmd.OutOrStdout()
				encoder := json.NewEncoder(out)
				encoder.SetIndent("", "  ")
				return encoder.Encode(output)
			}

			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "Features in %s:\n", projectID)
			fmt.Fprintf(out, "%-30s %-6s %-10s %s\n", "ID", "Priority", "Status", "Name")
			fmt.Fprintln(out, strings.Repeat("-", 80))

			for _, f := range output.Features {
				name := f.Name
				if len(name) > 40 {
					name = name[:37] + "..."
				}
				fmt.Fprintf(out, "%-30s %-6s %-10s %s\n", f.ID, f.Priority, f.Status, name)
			}

			fmt.Fprintf(out, "\nTotal: %d\n", output.Total)

			return nil
		},
	}

	cmd.Flags().StringVarP(&listProjectID, "project", "p", "", "Project ID (required)")
	cmd.Flags().BoolVar(&listJSON, "json", false, "Output as JSON")

	return cmd
}
