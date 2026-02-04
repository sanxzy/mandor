package workspace

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
)

// NewConfigCmd creates the config command
func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "View and modify workspace configuration",
		Long: `View and modify workspace configuration settings.

Available keys:
  - default_priority: Default priority for new entities (P0-P5, default: P3)
  - strict_mode: Enforce strict validation rules (true/false, default: false)
  - goal.lengths.project: Min chars for project goal (default: 500)
  - goal.lengths.feature: Min chars for feature goal (default: 300)
  - goal.lengths.task: Min chars for task goal (default: 500)
  - goal.lengths.issue: Min chars for issue goal (default: 200)`,
	}

	cmd.AddCommand(newConfigGetCmd())
	cmd.AddCommand(newConfigSetCmd())
	cmd.AddCommand(newConfigListCmd())
	cmd.AddCommand(newConfigResetCmd())

	return cmd
}

func newConfigGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get [key]",
		Short: "Display configuration value(s)",
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := service.NewWorkspaceService()
			if err != nil {
				return err
			}

			ws, err := svc.GetWorkspace()
			if err != nil {
				return err
			}

			if len(args) == 0 {
				// Show all config with workspace metadata
				fmt.Println("Workspace Information")
				fmt.Println("═════════════════════")
				fmt.Printf("Name:             %s\n", ws.Name)
				fmt.Printf("ID:               %s\n", ws.ID)
				fmt.Printf("Version:          %s\n", ws.Version)
				fmt.Printf("Created:          %s\n", ws.CreatedAt.Format("2006-01-02T15:04:05Z"))
				fmt.Printf("Created By:       %s\n", ws.CreatedBy)
				fmt.Printf("Last Updated:     %s\n", ws.LastUpdatedAt.Format("2006-01-02T15:04:05Z"))
				fmt.Println()
				fmt.Println("Configuration Settings")
				fmt.Println("══════════════════════")
				fmt.Println()
				fmt.Printf("default_priority  %s\n", ws.Config.DefaultPriority)
				fmt.Printf("strict_mode       %v\n", ws.Config.StrictMode)
				fmt.Println()
				fmt.Println("Goal Lengths (min chars)")
				fmt.Println("─────────────────────────")
				fmt.Printf("goal.lengths.project  %d\n", ws.Config.GoalLengths.Project)
				fmt.Printf("goal.lengths.feature  %d\n", ws.Config.GoalLengths.Feature)
				fmt.Printf("goal.lengths.task     %d\n", ws.Config.GoalLengths.Task)
				fmt.Printf("goal.lengths.issue    %d\n", ws.Config.GoalLengths.Issue)
				fmt.Println()
				fmt.Println("Project Dependency Rules")
				fmt.Println("════════════════════════")
				fmt.Printf("Task:              (configured per-project)\n")
				fmt.Printf("Feature:           (configured per-project)\n")
				fmt.Printf("Issue:             (configured per-project)\n")
				fmt.Println()
				fmt.Println("Use 'mandor config list' for detailed configuration information.")
				return nil
			}

			// Show specific key
			key := args[0]
			value, err := svc.GetConfigValue(key)
			if err != nil {
				return err
			}

			fmt.Printf("%s = %v\n", key, value)
			return nil
		},
	}
}

func newConfigSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set configuration value",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			valueStr := args[1]

			svc, err := service.NewWorkspaceService()
			if err != nil {
				return err
			}

			// Parse value based on key
			var value interface{}
			switch key {
			case "default_priority":
				value = strings.ToUpper(valueStr)
			case "strict_mode":
				boolValue, err := parseBool(valueStr)
				if err != nil {
					return domain.NewValidationError(
						"Invalid value for strict_mode.\nUse: true, false, yes, no, 1, or 0",
					)
				}
				value = boolValue
			case "goal.lengths.project", "goal.lengths.feature", "goal.lengths.task", "goal.lengths.issue":
				var length int
				_, err := fmt.Sscanf(valueStr, "%d", &length)
				if err != nil || length < 0 {
					return domain.NewValidationError(
						"Invalid value for goal.lengths.\nUse a non-negative integer (0-9999)",
					)
				}
				value = length
			default:
				return domain.NewValidationError(
					fmt.Sprintf("Unknown configuration key: %s\n\nAvailable keys:\n  - default_priority\n  - strict_mode\n  - goal.lengths.project\n  - goal.lengths.feature\n  - goal.lengths.task\n  - goal.lengths.issue", key),
				)
			}

			// Update config for non-goal.lengths keys
			if err := svc.UpdateWorkspaceConfig(key, value); err != nil {
				return err
			}

			ws, _ := svc.GetWorkspace()
			fmt.Printf("✓ Updated: %s = %v\n", key, value)
			fmt.Printf("  (workspace.json updated %s)\n", ws.LastUpdatedAt.Format("2006-01-02T15:04:05Z"))

			return nil
		},
	}
}

func newConfigListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all configuration keys with descriptions",
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := service.NewWorkspaceService()
			if err != nil {
				return err
			}

			ws, err := svc.GetWorkspace()
			if err != nil {
				return err
			}

			fmt.Println("Configuration Keys")
			fmt.Println("══════════════════")
			fmt.Println()

			// default_priority
			fmt.Println("default_priority")
			fmt.Println("  Type:     string")
			fmt.Printf("  Current:  %s\n", ws.Config.DefaultPriority)
			fmt.Println("  Default:  P3")
			fmt.Println("  Options:  P0, P1, P2, P3, P4, P5")
			fmt.Println("  Desc:     Default priority level for new entities")
			fmt.Println()

			// strict_mode
			fmt.Println("strict_mode")
			fmt.Println("  Type:     boolean")
			fmt.Printf("  Current:  %v\n", ws.Config.StrictMode)
			fmt.Println("  Default:  false")
			fmt.Println("  Options:  true, false")
			fmt.Println("  Desc:     Enforce strict validation rules")
			fmt.Println()

			// goal.lengths.project
			fmt.Println("goal.lengths.project")
			fmt.Println("  Type:     integer")
			fmt.Printf("  Current:  %d\n", ws.Config.GoalLengths.Project)
			fmt.Println("  Default:  500")
			fmt.Println("  Options:  0-9999")
			fmt.Println("  Desc:     Minimum characters for project goal")
			fmt.Println()

			// goal.lengths.feature
			fmt.Println("goal.lengths.feature")
			fmt.Println("  Type:     integer")
			fmt.Printf("  Current:  %d\n", ws.Config.GoalLengths.Feature)
			fmt.Println("  Default:  300")
			fmt.Println("  Options:  0-9999")
			fmt.Println("  Desc:     Minimum characters for feature goal")
			fmt.Println()

			// goal.lengths.task
			fmt.Println("goal.lengths.task")
			fmt.Println("  Type:     integer")
			fmt.Printf("  Current:  %d\n", ws.Config.GoalLengths.Task)
			fmt.Println("  Default:  500")
			fmt.Println("  Options:  0-9999")
			fmt.Println("  Desc:     Minimum characters for task goal")
			fmt.Println()

			// goal.lengths.issue
			fmt.Println("goal.lengths.issue")
			fmt.Println("  Type:     integer")
			fmt.Printf("  Current:  %d\n", ws.Config.GoalLengths.Issue)
			fmt.Println("  Default:  200")
			fmt.Println("  Options:  0-9999")
			fmt.Println("  Desc:     Minimum characters for issue goal")
			fmt.Println()

			fmt.Println("Use 'mandor config get <key>' for value.")
			fmt.Println("Use 'mandor config set <key> <value>' to update.")

			return nil
		},
	}
}

func newConfigResetCmd() *cobra.Command {
	var skipConfirm bool

	cmd := &cobra.Command{
		Use:   "reset [key]",
		Short: "Reset configuration to defaults",
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := service.NewWorkspaceService()
			if err != nil {
				return err
			}

			if len(args) == 0 {
				// Reset all
				if !skipConfirm {
					fmt.Print("Reset all configuration to defaults? [y/N] ")
					var response string
					fmt.Scanln(&response)
					if strings.ToLower(response) != "y" {
						return nil
					}
				}

				if err := svc.UpdateWorkspaceConfig("default_priority", "P3"); err != nil {
					return err
				}
				if err := svc.UpdateWorkspaceConfig("strict_mode", false); err != nil {
					return err
				}
				if err := svc.ResetAllGoalLengths(); err != nil {
					return err
				}

				fmt.Println("✓ Reset all configuration to defaults")
				fmt.Println("  - default_priority = P3")
				fmt.Println("  - strict_mode = false")
				fmt.Println("  - goal.lengths.project = 500")
				fmt.Println("  - goal.lengths.feature = 300")
				fmt.Println("  - goal.lengths.task = 500")
				fmt.Println("  - goal.lengths.issue = 200")
				return nil
			}

			// Reset specific key
			key := args[0]

			switch key {
			case "default_priority":
				if !skipConfirm {
					fmt.Print("Reset default_priority to default (P3)? [y/N] ")
					var response string
					fmt.Scanln(&response)
					if strings.ToLower(response) != "y" {
						return nil
					}
				}
				if err := svc.UpdateWorkspaceConfig("default_priority", "P3"); err != nil {
					return err
				}
				fmt.Printf("✓ Reset: default_priority = P3 (default)\n")
				return nil

			case "strict_mode":
				if !skipConfirm {
					fmt.Print("Reset strict_mode to default (false)? [y/N] ")
					var response string
					fmt.Scanln(&response)
					if strings.ToLower(response) != "y" {
						return nil
					}
				}
				if err := svc.UpdateWorkspaceConfig("strict_mode", false); err != nil {
					return err
				}
				fmt.Printf("✓ Reset: strict_mode = false (default)\n")
				return nil

			case "goal.lengths.project", "goal.lengths.feature", "goal.lengths.task", "goal.lengths.issue":
				entity := strings.Split(key, ".")[2]
				if !skipConfirm {
					fmt.Printf("Reset %s to default? [y/N] ", key)
					var response string
					fmt.Scanln(&response)
					if strings.ToLower(response) != "y" {
						return nil
					}
				}
				if err := svc.ResetGoalLength(entity); err != nil {
					return err
				}
				fmt.Printf("✓ Reset: %s (default)\n", key)
				return nil

			default:
				return domain.NewValidationError(
					fmt.Sprintf("Unknown configuration key: %s", key),
				)
			}
		},
	}

	cmd.Flags().BoolVarP(&skipConfirm, "yes", "y", false, "Skip confirmation prompt")
	return cmd
}

// parseBool parses boolean values with multiple formats
func parseBool(value string) (bool, error) {
	lowerValue := strings.ToLower(value)
	switch lowerValue {
	case "true", "yes", "1", "on":
		return true, nil
	case "false", "no", "0", "off":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean value: %s", value)
	}
}
