package spec

import (
	"fmt"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update <spec-id>",
	Short: "Update a Spec",
	Long:  "Update Spec fields",
	Args:  cobra.ExactArgs(1),
	RunE:  runUpdate,
}

var (
	updateSummary      string
	updateRequirements string
	updateStatus       string
)

func init() {
	updateCmd.Flags().StringVar(&updateSummary, "summary", "", "Updated description")
	updateCmd.Flags().StringVar(&updateRequirements, "requirements", "", "Updated requirements")
	updateCmd.Flags().StringVar(&updateStatus, "status", "", "draft | active | archived")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	specID := args[0]
	
	// TODO: Load Spec
	// TODO: Update fields
	// TODO: Validate updates
	// TODO: Store updated Spec
	
	fmt.Printf("Spec updated: %s\n", specID)
	fmt.Println("(Implementation pending)")
	
	return nil
}

func GetUpdateCmd() *cobra.Command {
	return updateCmd
}
