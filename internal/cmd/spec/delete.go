package spec

import (
	"fmt"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <spec-id>",
	Short: "Delete a Spec",
	Long:  "Archive or permanently delete a Spec",
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
	specID := args[0]
	
	// TODO: Load Spec
	// TODO: Confirm deletion if not --yes
	// TODO: Archive or hard delete
	
	fmt.Printf("Spec deleted: %s\n", specID)
	fmt.Println("(Implementation pending)")
	
	return nil
}

func GetDeleteCmd() *cobra.Command {
	return deleteCmd
}
