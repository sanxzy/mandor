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
  - strict_mode: Enforce strict validation rules (true/false, default: false)`,
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
			default:
				return domain.NewValidationError(
					fmt.Sprintf("Unknown configuration key: %s\n\nAvailable keys:\n  - default_priority\n  - strict_mode", key),
				)
			}

			// Update config
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

				fmt.Println("✓ Reset all configuration to defaults")
				fmt.Println("  - default_priority = P3")
				fmt.Println("  - strict_mode = false")
				return nil
			}

			// Reset specific key
			key := args[0]
			var defaultValue interface{}

			switch key {
			case "default_priority":
				defaultValue = "P3"
			case "strict_mode":
				defaultValue = false
			default:
				return domain.NewValidationError(
					fmt.Sprintf("Unknown configuration key: %s", key),
				)
			}

			if !skipConfirm {
				fmt.Printf("Reset %s to default (%v)? [y/N] ", key, defaultValue)
				var response string
				fmt.Scanln(&response)
				if strings.ToLower(response) != "y" {
					return nil
				}
			}

			if err := svc.UpdateWorkspaceConfig(key, defaultValue); err != nil {
				return err
			}

			fmt.Printf("✓ Reset: %s = %v (default)\n", key, defaultValue)
			return nil
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
