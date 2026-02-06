package spec

import (
	"fmt"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate <spec-id>",
	Short: "Validate a Spec",
	Long:  "Check that a Spec has valid structure (min 1 requirement with min 1 IAE scenario)",
	Args:  cobra.ExactArgs(1),
	RunE:  runValidate,
}

func runValidate(cmd *cobra.Command, args []string) error {
	specID := args[0]
	
	// TODO: Load Spec
	// TODO: Validate structure:
	//   - Min 1 requirement
	//   - Each requirement has min 1 IAE scenario
	//   - Each scenario has Intent, Action, Expect (all non-empty)
	//   - Requirement ID format: req-XXXX
	//   - Scenario ID format: scenario-XXXX
	// TODO: Report validation results
	
	fmt.Printf("Spec validation: %s\n", specID)
	fmt.Println("(Implementation pending)")
	
	return nil
}

func GetValidateCmd() *cobra.Command {
	return validateCmd
}
