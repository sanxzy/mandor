package brief

import (
	"fmt"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate <brief-id>",
	Short: "Validate a Brief",
	Long:  "Check that a Brief has valid structure and all required sections",
	Args:  cobra.ExactArgs(1),
	RunE:  runValidate,
}

func runValidate(cmd *cobra.Command, args []string) error {
	briefID := args[0]
	
	// TODO: Load Brief
	// TODO: Validate structure (Why, WhatChanges, Impact, Capabilities)
	// TODO: Validate Why length (100-5000 chars)
	// TODO: Validate capability ID format
	// TODO: Report validation results
	
	fmt.Printf("Brief validation: %s\n", briefID)
	fmt.Println("(Implementation pending)")
	
	return nil
}

func GetValidateCmd() *cobra.Command {
	return validateCmd
}
