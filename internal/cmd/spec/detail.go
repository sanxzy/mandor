package spec

import (
	"fmt"

	"github.com/spf13/cobra"
)

var detailCmd = &cobra.Command{
	Use:   "detail <spec-id>",
	Short: "Display Spec details",
	Long:  "Show complete Spec with all requirements and scenarios",
	Args:  cobra.ExactArgs(1),
	RunE:  runDetail,
}

func runDetail(cmd *cobra.Command, args []string) error {
	specID := args[0]
	
	// TODO: Load Spec from filesystem
	// TODO: Display with full context (requirements, scenarios)
	
	fmt.Printf("Spec detail: %s\n", specID)
	fmt.Println("(Implementation pending)")
	
	return nil
}

func GetDetailCmd() *cobra.Command {
	return detailCmd
}
