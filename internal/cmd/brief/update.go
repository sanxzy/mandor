package brief

import (
	"fmt"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update <brief-id>",
	Short: "Update a Brief",
	Long:  "Update Brief fields (why, capabilities, impact, etc.)",
	Args:  cobra.ExactArgs(1),
	RunE:  runUpdate,
}

var (
	updateWhy          string
	updateCapabilities string
	updateStatus       string
)

func init() {
	updateCmd.Flags().StringVar(&updateWhy, "why", "", "Updated problem statement")
	updateCmd.Flags().StringVar(&updateCapabilities, "capabilities", "", "Updated capabilities")
	updateCmd.Flags().StringVar(&updateStatus, "status", "", "draft | active | archived")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	briefID := args[0]
	
	// TODO: Load Brief from filesystem
	// TODO: Update fields
	// TODO: Validate updates
	// TODO: Store updated Brief
	
	fmt.Printf("Brief updated: %s\n", briefID)
	fmt.Println("(Implementation pending)")
	
	return nil
}

func GetUpdateCmd() *cobra.Command {
	return updateCmd
}
