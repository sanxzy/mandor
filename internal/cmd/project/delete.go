package project

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
)

var (
	hardDelete   bool
	dryRunDelete bool
	yesDelete    bool
)

func NewDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a project",
		Long:  "Delete a project (soft delete by default, hard delete with --hard).",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := service.NewProjectService()
			if err != nil {
				return err
			}

			if !svc.WorkspaceInitialized() {
				return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
			}

			input := &domain.ProjectDeleteInput{
				ID:     args[0],
				Hard:   hardDelete,
				DryRun: dryRunDelete,
				Yes:    yesDelete,
			}

			if err := svc.ValidateDeleteInput(input); err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			if !input.DryRun && !input.Yes && input.Hard {
				fmt.Fprintln(out, "⚠ WARNING: Hard delete is PERMANENT and cannot be undone.")
				fmt.Fprintf(out, "All project files will be removed.\n\n")
				fmt.Fprintf(out, "Type 'HARD DELETE %s' to confirm: ", args[0])
				scanner := bufio.NewScanner(os.Stdin)
				if !scanner.Scan() {
					return domain.NewValidationError("Invalid confirmation. Delete cancelled.")
				}
				confirm := scanner.Text()
				if confirm != "HARD DELETE "+args[0] {
					return domain.NewValidationError("Invalid confirmation. Delete cancelled.")
				}
			}

			result, err := svc.DeleteProject(input)
			if err != nil {
				return err
			}

			fmt.Fprintln(out, result)

			return nil
		},
	}

	cmd.Flags().BoolVar(&hardDelete, "hard", false, "Permanently delete project and all files")
	cmd.Flags().BoolVar(&dryRunDelete, "dry-run", false, "Preview deletion without applying")
	cmd.Flags().BoolVarP(&yesDelete, "yes", "y", false, "Skip confirmation prompts")

	return cmd
}
