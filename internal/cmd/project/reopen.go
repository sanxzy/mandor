package project

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
)

var yesReopen bool

func NewReopenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reopen <id>",
		Short: "Reopen a soft-deleted project",
		Long:  "Reopen a soft-deleted project, restoring it to active state.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := service.NewProjectService()
			if err != nil {
				return err
			}

			if !svc.WorkspaceInitialized() {
				return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
			}

			input := &domain.ProjectReopenInput{
				ID:  args[0],
				Yes: yesReopen,
			}

			if err := svc.ValidateReopenInput(input); err != nil {
				return err
			}

			project, err := svc.GetProject(input.ID)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			if !yesReopen {
				fmt.Fprintf(out, "Reopen project: %s\n", input.ID)
				fmt.Fprintf(out, "  Name: %s\n", project.Name)
				fmt.Fprintf(out, "  Status: deleted\n\n")
				fmt.Fprint(out, "Confirm reopen? [y/N]: ")
				scanner := bufio.NewScanner(os.Stdin)
				if scanner.Scan() {
					confirm := scanner.Text()
					if confirm != "y" && confirm != "Y" {
						fmt.Fprintln(out, "Cancelled. Project not reopened.")
						return nil
					}
				}
			}

			result, err := svc.ReopenProject(input)
			if err != nil {
				return err
			}

			fmt.Fprintln(out, result)
			fmt.Fprintf(out, "  Name:   %s\n", project.Name)
			fmt.Fprintf(out, "  Status: initial\n")
			fmt.Fprintf(out, "  By:     %s\n", project.UpdatedBy)

			return nil
		},
	}

	cmd.Flags().BoolVarP(&yesReopen, "yes", "y", false, "Skip confirmation prompt")

	return cmd
}
