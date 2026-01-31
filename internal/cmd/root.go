package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"mandor/internal/cmd/ai"
	"mandor/internal/cmd/feature"
	"mandor/internal/cmd/issue"
	"mandor/internal/cmd/populate"
	"mandor/internal/cmd/project"
	"mandor/internal/cmd/task"
	"mandor/internal/cmd/workspace"
	"mandor/internal/domain"
)

// NewRootCmd creates the root command
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "mandor",
		Short: "Mandor - Event-Based Task Manager for AI Agent Workflows",
		Long: `Mandor is a CLI tool for deterministic, streaming-native task management.
It provides schema-driven, event-based task management with dependency tracking.

For more information, visit: https://github.com/budisantoso/mandor`,
	}

	// Add workspace commands
	rootCmd.AddCommand(workspace.NewInitCmd())
	rootCmd.AddCommand(workspace.NewStatusCmd())
	rootCmd.AddCommand(workspace.NewConfigCmd())

	// Add project commands
	rootCmd.AddCommand(project.NewProjectCmd())

	// Add feature commands
	rootCmd.AddCommand(feature.NewFeatureCmd())

	// Add task commands
	rootCmd.AddCommand(task.NewTaskCmd())

	// Add issue commands
	rootCmd.AddCommand(issue.NewIssueCmd())

	// Add completion command
	rootCmd.AddCommand(NewCompletionCmd(rootCmd))

	// Add populate command
	rootCmd.AddCommand(populate.NewPopulateCmd())

	// Add AI commands
	rootCmd.AddCommand(ai.NewAICmd())

	// Add version command
	rootCmd.AddCommand(NewVersionCmd())

	return rootCmd
}

// ExecuteWithCode executes the command and returns the appropriate exit code
func ExecuteWithCode() int {
	rootCmd := NewRootCmd()
	// Disable automatic error printing from cobra
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true

	err := rootCmd.Execute()

	if err == nil {
		return 0
	}

	// Check if it's a MandorError
	if me, ok := err.(*domain.MandorError); ok {
		fmt.Fprintf(os.Stderr, "%v\n", me.Error())
		return int(me.Code)
	}

	// Other errors
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	return 1
}
