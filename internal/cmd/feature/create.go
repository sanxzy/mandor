package feature

import (
	"fmt"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
	"mandor/internal/util"
)

var (
	projectID string
	name      string
	goal      string
	scope     string
	priority  string
	dependsOn string
)

func NewCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create <name> --project <id> --goal <text>",
		Short: "Create a new feature",
		Long:  "Create a new feature in the specified project with the given name and goal.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := service.NewFeatureService()
			if err != nil {
				return err
			}

			if !svc.WorkspaceInitialized() {
				return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
			}

			projectID, _ := cmd.Flags().GetString("project")
			if projectID == "" {
				return domain.NewValidationError("Project ID is required (--project).")
			}

			if goal == "" {
				return domain.NewValidationError("Feature goal is required (--goal).")
			}

			var dependsOnList []string
			if dependsOn != "" {
				dependsOnList = splitDependsOn(dependsOn)
			}

			input := &domain.FeatureCreateInput{
				ProjectID: projectID,
				Name:      args[0],
				Goal:      goal,
				Scope:     scope,
				Priority:  priority,
				DependsOn: dependsOnList,
			}

			if err := svc.ValidateCreateInput(input); err != nil {
				return err
			}

			feature, err := svc.CreateFeature(input)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "Feature created: %s\n", feature.ID)
			fmt.Fprintf(out, "  Name:     %s\n", feature.Name)
			fmt.Fprintf(out, "  Project:  %s\n", feature.ProjectID)
			fmt.Fprintf(out, "  Goal:     %s\n", feature.Goal)
			fmt.Fprintf(out, "  Scope:    %s\n", feature.Scope)
			fmt.Fprintf(out, "  Priority: %s\n", feature.Priority)
			fmt.Fprintf(out, "  Status:   %s\n", feature.Status)

			_, warning := util.GetGitUsernameWithWarning()
			if warning != "" {
				fmt.Fprintln(out)
				fmt.Fprintln(out, warning)
				fmt.Fprintln(out, "  Run: git config user.name \"Your Name\"")
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectID, "project", "p", "", "Project ID (required)")
	cmd.Flags().StringVarP(&goal, "goal", "g", "", "Feature goal (required)")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Feature name (alternative to positional)")
	cmd.Flags().StringVar(&scope, "scope", "", "Feature scope (frontend, backend, fullstack, cli, desktop, android, flutter, react-native, ios, swift)")
	cmd.Flags().StringVar(&priority, "priority", "P3", "Priority (P0-P5)")
	cmd.Flags().StringVar(&dependsOn, "depends", "", "Comma-separated feature IDs this feature depends on")

	return cmd
}

func splitDependsOn(s string) []string {
	if s == "" {
		return nil
	}
	var result []string
	for _, part := range splitByComma(s) {
		trimmed := trimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func splitByComma(s string) []string {
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == ',' {
			result = append(result, s[start:i])
			start = i + 1
		}
	}
	result = append(result, s[start:])
	return result
}

func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}
