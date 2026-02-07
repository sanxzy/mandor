package feature

import (
	"fmt"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
	"mandor/internal/util"
)

func NewCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create <name> [--backlog <id>] --capability <cap-id> --spec-id <spec-id> --goal <text>",
		Short: "Create a new feature",
		Long:  "Create a new feature in the specified backlog with ONE-TO-ONE Spec mapping (immutable)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := service.NewFeatureService()
			if err != nil {
				return err
			}

			if !svc.WorkspaceInitialized() {
				return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
			}

			// Get backlog ID from flags or parent commands
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

			// Get capability ID from flags
			capID, _ := cmd.Flags().GetString("capability")
			if capID == "" {
				return domain.NewValidationError("Capability ID is required (--capability).")
			}

			// Get spec ID from flags
			sID, _ := cmd.Flags().GetString("spec-id")
			if sID == "" {
				return domain.NewValidationError("Spec ID is required (--spec-id).")
			}

			// Get goal from flags
			goalText, _ := cmd.Flags().GetString("goal")
			if goalText == "" {
				return domain.NewValidationError("Feature goal is required (--goal).")
			}

			// Get optional fields from flags
			featureName, _ := cmd.Flags().GetString("name")
			if featureName == "" {
				featureName = args[0]
			}

			featureScope, _ := cmd.Flags().GetString("scope")
			featurePriority, _ := cmd.Flags().GetString("priority")
			dependsOnStr, _ := cmd.Flags().GetString("depends")

			var dependsOnList []string
			if dependsOnStr != "" {
				dependsOnList = splitDependsOn(dependsOnStr)
			}

			input := &domain.FeatureCreateInput{
				BacklogID:    backlogID,
				CapabilityID: capID,
				SpecID:       sID,
				Name:         featureName,
				Goal:         goalText,
				Scope:        featureScope,
				Priority:     featurePriority,
				DependsOn:    dependsOnList,
			}

			if err := svc.ValidateCreateInput(input); err != nil {
				return err
			}

			feature, err := svc.CreateFeature(input)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "✓ Feature created: %s\n", feature.ID)
			fmt.Fprintf(out, "  Name:      %s\n", feature.Name)
			fmt.Fprintf(out, "  Backlog:  %s\n", feature.BacklogID)
			fmt.Fprintf(out, "  Capability: %s\n", feature.CapabilityID)
			fmt.Fprintf(out, "  Spec ID:   %s\n", feature.SpecID)
			fmt.Fprintf(out, "  Goal:      %s\n", feature.Goal)
			fmt.Fprintf(out, "  Scope:     %s\n", feature.Scope)
			fmt.Fprintf(out, "  Priority:  %s\n", feature.Priority)
			fmt.Fprintf(out, "  Status:    %s\n", feature.Status)

			// AI-friendly next steps
			fmt.Fprintln(out)
			fmt.Fprintln(out, "  ─────────────────────────────────────────────────────────")
			fmt.Fprintln(out, "  NEXT STEPS:")
			fmt.Fprintln(out, "  ─────────────────────────────────────────────────────────")
			fmt.Fprintf(out, "  • Create tasks: mandor task create %s <task-name> --spec-id %s --iae-scenarios <req:scenario> ...\n", feature.ID, feature.SpecID)
			fmt.Fprintf(out, "  • View feature details: mandor feature detail %s --backlog %s\n", feature.ID, backlogID)
			fmt.Fprintf(out, "  • List all features: mandor feature list --backlog %s\n", backlogID)

			_, warning := util.GetGitUsernameWithWarning()
			if warning != "" {
				fmt.Fprintln(out)
				fmt.Fprintln(out, warning)
				fmt.Fprintln(out, "  Run: git config user.name \"Your Name\"")
			}

			return nil
		},
	}

	cmd.Flags().StringP("backlog", "b", "", "Backlog ID (required, use -b or --backlog)")
	cmd.Flags().String("capability", "", "Capability ID from Brief (required)")
	cmd.Flags().String("spec-id", "", "Spec ID (required, ONE-TO-ONE immutable mapping)")
	cmd.Flags().StringP("goal", "g", "", "Feature goal (required, min 300 chars, include technical user flow and complete requirements)")
	cmd.Flags().StringP("name", "n", "", "Feature name (alternative to positional)")
	cmd.Flags().String("scope", "", "Feature scope (frontend, backend, fullstack, cli, desktop, android, flutter, react-native, ios, swift)")
	cmd.Flags().String("priority", "", "Priority (P0-P5, default from config)")
	cmd.Flags().String("depends", "", "Pipe-separated feature IDs this feature depends on")

	cmd.MarkFlagRequired("capability")
	cmd.MarkFlagRequired("spec-id")
	cmd.MarkFlagRequired("goal")

	return cmd
}

func splitDependsOn(s string) []string {
	if s == "" {
		return nil
	}
	var result []string
	for _, part := range splitByPipe(s) {
		trimmed := trimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func splitByPipe(s string) []string {
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '|' {
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
