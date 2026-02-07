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
	rejectChangeID  string
	rejectReason    string
	rejectBacklogID string
	rejectJSON      bool
)

// NewRejectCmd creates the reject subcommand
func NewRejectCmd(paths *fs.Paths) *cobra.Command {
	rejectCmd := &cobra.Command{
		Use:   "reject",
		Short: "Reject a pending change",
		Long:  "Reject a change if impact is too high or remediation is not feasible",
		RunE: func(c *cobra.Command, args []string) error {
			return runReject(paths)
		},
	}

	rejectCmd.Flags().StringVar(&rejectChangeID, "change-id", "", "ID of the change to reject")
	rejectCmd.Flags().StringVar(&rejectReason, "reason", "", "Reason for rejection (minimum 10 characters)")
	rejectCmd.Flags().StringVar(&rejectBacklogID, "backlog", "", "Backlog ID (optional)")
	rejectCmd.Flags().BoolVar(&rejectJSON, "json", false, "Output as JSON")

	rejectCmd.MarkFlagRequired("change-id")
	rejectCmd.MarkFlagRequired("reason")

	return rejectCmd
}

func runReject(paths *fs.Paths) error {
	// Validate input
	if strings.TrimSpace(rejectChangeID) == "" {
		return fmt.Errorf("change ID is required")
	}

	if strings.TrimSpace(rejectReason) == "" {
		return fmt.Errorf("rejection reason is required (minimum 10 characters)")
	}

	if len(rejectReason) < 10 {
		return fmt.Errorf("rejection reason too short (minimum 10 characters, got %d)", len(rejectReason))
	}

	if rejectJSON {
		response := map[string]interface{}{
			"status":    "rejected",
			"change_id": rejectChangeID,
			"message":   fmt.Sprintf("Change %s rejected", rejectChangeID),
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		}
		data, _ := json.MarshalIndent(response, "", "  ")
		fmt.Println(string(data))
	} else {
		fmt.Printf("✗ Change %s rejected\n", rejectChangeID)
		fmt.Printf("  Status: Rejected\n")
		fmt.Printf("  Reason: %s\n", rejectReason)
	}

	return nil
}
