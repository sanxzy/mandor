package blueprint

import (
	"fmt"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <blueprint-id>",
	Short: "Delete a Blueprint",
	Long:  "Archive or permanently delete a Blueprint",
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
	blueprintID := args[0]
	
	// TODO: Load Blueprint
	// TODO: Confirm deletion if not --yes
	// TODO: Archive or hard delete
	
	fmt.Printf("Blueprint deleted: %s\n", blueprintID)
	fmt.Println("(Implementation pending)")
	
	return nil
}

func GetDeleteCmd() *cobra.Command {
	return deleteCmd
}
