package brief

import (
	"fmt"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <brief-id>",
	Short: "Delete a Brief",
	Long:  "Archive or permanently delete a Brief",
	Args:  cobra.ExactArgs(1),
	RunE:  runDelete,
}

var (
	deleteHard bool
	deleteYes  bool
)

func init() {
	deleteCmd.Flags().BoolVar(&deleteHard, "hard", false, "Permanently delete instead of archiving")
	deleteCmd.Flags().BoolVarP(&deleteYes, "yes", "y", false, "Skip confirmation")
}

func runDelete(cmd *cobra.Command, args []string) error {
	briefID := args[0]
	
	// TODO: Load Brief
	// TODO: Confirm deletion if not --yes
	// TODO: Archive or hard delete
	
	fmt.Printf("Brief deleted: %s\n", briefID)
	fmt.Println("(Implementation pending)")
	
	return nil
}

func GetDeleteCmd() *cobra.Command {
	return deleteCmd
}
