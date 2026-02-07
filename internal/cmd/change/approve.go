package change

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"mandor/internal/fs"
	"mandor/internal/service"
)

var (
	approveChangeID  string
	approveReason    string
	approveBacklogID string
	approveJSON      bool
)

// NewApproveCmd creates the approve subcommand
func NewApproveCmd(paths *fs.Paths) *cobra.Command {
	approveCmd := &cobra.Command{
		Use:   "approve",
		Short: "Approve a pending change",
		Long:  "Approve a change after impact analysis and remediation",
		RunE: func(c *cobra.Command, args []string) error {
			return runApprove(paths)
		},
	}

	approveCmd.Flags().StringVar(&approveChangeID, "change-id", "", "ID of the change to approve")
	approveCmd.Flags().StringVar(&approveReason, "reason", "", "Reason for approval (minimum 10 characters)")
	approveCmd.Flags().StringVar(&approveBacklogID, "backlog", "", "Backlog ID (optional)")
	approveCmd.Flags().BoolVar(&approveJSON, "json", false, "Output as JSON")

	approveCmd.MarkFlagRequired("change-id")
	approveCmd.MarkFlagRequired("reason")

	return approveCmd
}

func runApprove(paths *fs.Paths) error {
	// Validate input
	if strings.TrimSpace(approveChangeID) == "" {
		return fmt.Errorf("change ID is required")
	}

	if strings.TrimSpace(approveReason) == "" {
		return fmt.Errorf("approval reason is required (minimum 10 characters)")
	}

	if len(approveReason) < 10 {
		return fmt.Errorf("approval reason too short (minimum 10 characters, got %d)", len(approveReason))
	}

	// Perform approval
	governanceService := service.NewChangeGovernanceService(paths)
	approval, err := governanceService.ApproveChange(approveBacklogID, approveChangeID, approveReason)
	if err != nil {
		return err
	}

	if approveJSON {
		response := map[string]interface{}{
			"approval":  approval,
			"message":   fmt.Sprintf("Change %s approved successfully", approveChangeID),
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		}
		data, _ := json.MarshalIndent(response, "", "  ")
		fmt.Println(string(data))
	} else {
		fmt.Printf("✓ Change %s approved successfully\n", approveChangeID)
		fmt.Printf("  Status: Approved\n")
		fmt.Printf("  Reason: %s\n", approveReason)
	}

	return nil
}
