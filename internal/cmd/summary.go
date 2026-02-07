package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
)

func NewSummaryCmd() *cobra.Command {
	var backlogID string

	cmd := &cobra.Command{
		Use:        "summary [--backlog <id>]",
		Short:      "Display workspace summary",
		Deprecated: "Use 'mandor track' instead for unified workspace/backlog/feature/task/issue tracking. See 'mandor track --help' for options.",
		Long:       "Display a summary of all features grouped by priority with task counts and status overview.\n\nDEPRECATED: This command is superseded by 'mandor track'. For equivalent functionality use 'mandor track workspace' or 'mandor track backlog <id>'.",
		RunE: func(cmd *cobra.Command, args []string) error {
			fsvc, err := service.NewFeatureService()
			if err != nil {
				return err
			}

			if !fsvc.WorkspaceInitialized() {
				return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
			}

			// Find all backlogs in workspace
			backlogsDir := ".mandor/backlogs"
			entries, err := os.ReadDir(backlogsDir)
			if err != nil {
				fmt.Fprintf(cmd.OutOrStdout(), "No backlogs in workspace.\n")
				return nil
			}

			var backlogs []string
			for _, entry := range entries {
				if entry.IsDir() {
					backlogs = append(backlogs, entry.Name())
				}
			}

			if len(backlogs) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No backlogs in workspace.\n")
				return nil
			}

			// Group features by priority across all backlogs
			featuresByPriority := make(map[string][]domain.FeatureListItem)
			priorityOrder := []string{"P0", "P1", "P2", "P3", "P4", "P5"}

			// Iterate through each backlog and gather features
			for _, bID := range backlogs {
				// Skip if filtering by backlog and this isn't it
				if backlogID != "" && bID != backlogID {
					continue
				}

				input := &domain.FeatureListInput{
					BacklogID: bID,
				}

				output, err := fsvc.ListFeatures(input)
				if err != nil {
					continue
				}

				if output != nil {
					// Group by priority
					for _, feature := range output.Features {
						featuresByPriority[feature.Priority] = append(featuresByPriority[feature.Priority], feature)
					}
				}
			}

			totalFeatures := 0
			for _, features := range featuresByPriority {
				totalFeatures += len(features)
			}

			if totalFeatures == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No features in workspace.\n")
				return nil
			}

			out := cmd.OutOrStdout()

			// Print by priority
			for _, priority := range priorityOrder {
				features := featuresByPriority[priority]
				if len(features) == 0 {
					continue
				}

				priorityLabel := getPriorityLabel(priority)
				fmt.Fprintf(out, "%s - %s (%d features)\n", priority, priorityLabel, len(features))
				fmt.Fprintf(out, "| # | Feature ID | Name | Goal Summary |\n")
				fmt.Fprintf(out, "|---|------------|------|---------------|\n")

				sort.Slice(features, func(i, j int) bool {
					return features[i].CreatedAt < features[j].CreatedAt
				})

				for idx, feature := range features {
					goalSummary := truncateGoal(feature.Goal, 60)
					fmt.Fprintf(out, "| %d | %s | %s | %s |\n",
						idx+1,
						feature.ID,
						feature.Name,
						goalSummary)
				}

				fmt.Fprintln(out)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&backlogID, "backlog", "b", "", "Filter by backlog ID")

	return cmd
}

func getPriorityLabel(priority string) string {
	labels := map[string]string{
		"P0": "Critical",
		"P1": "High",
		"P2": "Medium",
		"P3": "Normal",
		"P4": "Low",
		"P5": "Minimal",
	}
	if label, ok := labels[priority]; ok {
		return label
	}
	return "Unknown"
}

func truncateGoal(goal string, maxLen int) string {
	if len(goal) <= maxLen {
		return goal
	}
	return goal[:maxLen-3] + "..."
}
