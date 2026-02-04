package session

import (
	"github.com/spf13/cobra"
)

// NewSessionCmd creates the session command
func NewSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "session",
		Short: "Session management for AI progress tracking",
		Long: `Session management commands for AI agents to track progress across sessions.

Use session notes to record what was done and what's next, helping maintain
context between AI work sessions.`,
	}

	cmd.AddCommand(NewNoteCmd())

	return cmd
}
