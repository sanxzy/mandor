package brief

import (
	"fmt"

	"github.com/spf13/cobra"
)

var readCmd = &cobra.Command{
	Use:   "read <brief-id>",
	Short: "Read a Brief",
	Long:  "Display the complete Brief document",
	Args:  cobra.ExactArgs(1),
	RunE:  runRead,
}

func runRead(cmd *cobra.Command, args []string) error {
	briefID := args[0]
	
	// TODO: Load Brief from filesystem (.mandor/projects/<project-id>/brief.md)
	// TODO: Parse Markdown to Brief struct
	// TODO: Display Brief
	
	fmt.Printf("Brief: %s\n", briefID)
	fmt.Println("(Implementation pending)")
	
	return nil
}

func GetReadCmd() *cobra.Command {
	return readCmd
}
