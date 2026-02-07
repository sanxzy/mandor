package change

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"mandor/internal/fs"
)

var (
	listBacklogID  string
	listStatus     string
	listEntityType string
	listJSON       bool
)

// NewListCmd creates the list subcommand
func NewListCmd(paths *fs.Paths) *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all changes for a backlog",
		Long:  "List all changes (pending, approved, rejected) for a backlog with optional filtering",
		RunE: func(c *cobra.Command, args []string) error {
			return runList(paths)
		},
	}

	listCmd.Flags().StringVar(&listBacklogID, "backlog", "", "Backlog ID (optional)")
	listCmd.Flags().StringVar(&listStatus, "status", "all", "Filter by status: pending_validation, approved, rejected, or all")
	listCmd.Flags().StringVar(&listEntityType, "entity-type", "all", "Filter by entity type: brief, spec, blueprint, or all")
	listCmd.Flags().BoolVar(&listJSON, "json", false, "Output as JSON")

	return listCmd
}

func runList(paths *fs.Paths) error {
	// Validate input
	if strings.TrimSpace(listBacklogID) == "" {
		return fmt.Errorf("backlog ID is required")
	}

	// Validate status filter
	if listStatus != "" && listStatus != "all" {
		validStatuses := []string{"pending_validation", "approved", "rejected"}
		found := false
		for _, s := range validStatuses {
			if listStatus == s {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("invalid status: %s (must be pending_validation, approved, rejected, or all)", listStatus)
		}
	}

	// Validate entity type filter
	if listEntityType != "" && listEntityType != "all" {
		validTypes := []string{"brief", "spec", "blueprint"}
		found := false
		for _, t := range validTypes {
			if listEntityType == t {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("invalid entity type: %s (must be brief, spec, blueprint, or all)", listEntityType)
		}
	}

	// For now, return empty list (audit persistence not yet implemented)
	response := map[string]interface{}{
		"changes":   []map[string]interface{}{},
		"total":     0,
		"filtered":  0,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	if listJSON {
		data, _ := json.MarshalIndent(response, "", "  ")
		fmt.Println(string(data))
	} else {
		fmt.Printf("Total Changes: %d\n", response["total"])
		fmt.Printf("Filtered Results: %d\n", response["filtered"])
		fmt.Println("\nNo changes found.")
	}

	return nil
}
