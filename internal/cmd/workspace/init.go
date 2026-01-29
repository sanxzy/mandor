package workspace

import (
	"fmt"

	"github.com/spf13/cobra"
	"mandor/internal/service"
	"mandor/internal/util"
)

// NewInitCmd creates the init command
func NewInitCmd() *cobra.Command {
	var (
		workspaceName string
		skipConfirm   bool
		strict        bool
	)

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new Mandor workspace",
		Long: `Initialize a new Mandor workspace in the current directory.

Creates a .mandor/ directory with workspace metadata and project storage.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := service.NewWorkspaceService()
			if err != nil {
				return err
			}

			ws, err := svc.InitWorkspace(workspaceName)
			if err != nil {
				return err
			}

			// Output success message
			fmt.Printf("✓ Workspace initialized: %s\n", ws.Name)
			fmt.Printf("  Location: .mandor/\n")
			fmt.Printf("  ID: %s\n", ws.ID)
			fmt.Printf("  Creator: %s\n", ws.CreatedBy)
			fmt.Printf("  Created: %s\n", ws.CreatedAt.Format("2006-01-02T15:04:05Z"))

			username, warning := util.GetGitUsernameWithWarning()
			if username == "unknown" {
				fmt.Printf("\n")
				fmt.Printf("Warning: Git user not configured. Events will show 'unknown' as creator.\n")
				fmt.Printf("  Run: git config user.name \"Your Name\"\n")
			}

			fmt.Printf("\nNext steps:\n")
			fmt.Printf("  1. Create a project: mandor project create <project_id> --name \"Project Name\"\n")
			fmt.Printf("  2. View status: mandor status\n")
			fmt.Printf("  3. Check config: mandor config get\n")

			_ = warning
			return nil
		},
	}

	cmd.Flags().StringVarP(&workspaceName, "workspace-name", "", "", "Custom workspace name (default: current directory)")
	cmd.Flags().BoolVarP(&skipConfirm, "yes", "y", false, "Skip confirmation prompts")
	cmd.Flags().BoolVarP(&strict, "strict", "", false, "Enforce strict dependency rules (deprecated)")

	return cmd
}
