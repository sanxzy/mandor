package issue

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/service"
	"mandor/internal/util"
)

var (
	createProjectID     string
	createName          string
	createGoal          string
	createType          string
	createPriority      string
	createDependsOn     string
	createAffectedFiles string
	createAffectedTests string
	createImplSteps     string
	createLibraries     string
	createYes           bool
)

func NewCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create <name> --project <id> --type <type> --goal <text> --affected-files <files> --affected-tests <tests> --implementation-steps <steps> [--priority <P0-P5>] [--depends-on <ids>] [--library-needs <libs>] [-y]",
		Short: "Create a new issue",
		Long:  "Create a new issue in the specified project with the given details.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := service.NewIssueService()
			if err != nil {
				return err
			}

			if !svc.WorkspaceInitialized() {
				return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
			}

			if createProjectID == "" {
				return domain.NewValidationError("Project ID is required (--project).")
			}

			if createType == "" {
				return domain.NewValidationError("Issue type is required (--type).")
			}

			if createGoal == "" {
				return domain.NewValidationError("Issue goal is required (--goal).")
			}

			affectedFiles := splitByPipe(createAffectedFiles)
			if len(affectedFiles) == 0 || (len(affectedFiles) == 1 && affectedFiles[0] == "") {
				return domain.NewValidationError("Affected files are required (--affected-files).")
			}

			affectedTests := splitByPipe(createAffectedTests)
			if len(affectedTests) == 0 || (len(affectedTests) == 1 && affectedTests[0] == "") {
				return domain.NewValidationError("Affected tests are required (--affected-tests).")
			}

			implSteps := splitByPipe(createImplSteps)
			if len(implSteps) == 0 || (len(implSteps) == 1 && implSteps[0] == "") {
				return domain.NewValidationError("Implementation steps are required (--implementation-steps).")
			}

			libraries := splitByPipe(createLibraries)
			if len(libraries) == 0 || (len(libraries) == 1 && libraries[0] == "") {
				libraries = []string{}
			}

			var dependsOnList []string
			if createDependsOn != "" {
				dependsOnList = splitByPipe(createDependsOn)
			}

			input := &domain.IssueCreateInput{
				ProjectID:           createProjectID,
				Name:                args[0],
				Goal:                createGoal,
				IssueType:           createType,
				Priority:            createPriority,
				DependsOn:           dependsOnList,
				AffectedFiles:       affectedFiles,
				AffectedTests:       affectedTests,
				ImplementationSteps: implSteps,
				LibraryNeeds:        libraries,
			}

			if err := svc.ValidateCreateInput(input); err != nil {
				return err
			}

			issue, err := svc.CreateIssue(input)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "Issue created: %s\n", issue.ID)
			fmt.Fprintf(out, "  Name:               %s\n", issue.Name)
			fmt.Fprintf(out, "  Type:               %s\n", issue.IssueType)
			fmt.Fprintf(out, "  Priority:           %s\n", issue.Priority)
			fmt.Fprintf(out, "  Status:             %s\n", issue.Status)
			fmt.Fprintf(out, "  Goal:               %s\n", truncate(issue.Goal, 60))
			fmt.Fprintf(out, "  Affected Files:     %d\n", len(issue.AffectedFiles))
			fmt.Fprintf(out, "  Affected Tests:     %d\n", len(issue.AffectedTests))
			fmt.Fprintf(out, "  Implementation Steps: %d\n", len(issue.ImplementationSteps))
			fmt.Fprintf(out, "  Library Needs:      %d\n", len(issue.LibraryNeeds))
			if len(issue.DependsOn) > 0 {
				fmt.Fprintf(out, "  Depends on:         %d issue(s)\n", len(issue.DependsOn))
			}

			_, warning := util.GetGitUsernameWithWarning()
			if warning != "" {
				fmt.Fprintln(out)
				fmt.Fprintln(out, warning)
				fmt.Fprintln(out, "  Run: git config user.name \"Your Name\"")
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&createProjectID, "project", "p", "", "Project ID (required, use -p or --project)")
	cmd.Flags().StringVarP(&createType, "type", "t", "", "Issue type: bug, improvement, debt, security, performance (required, use -t or --type)")
	cmd.Flags().StringVar(&createName, "name", "", "Issue name (required for CLI, or use positional argument)")
	cmd.Flags().StringVarP(&createGoal, "goal", "g", "", "Issue goal (required, min 200 chars, include problem description, impact analysis, and acceptance criteria)")
	cmd.Flags().StringVar(&createPriority, "priority", "", "Priority (P0-P5, default from config)")
	cmd.Flags().StringVar(&createDependsOn, "depends-on", "", "Pipe-separated issue IDs this issue depends on")
	cmd.Flags().StringVar(&createAffectedFiles, "affected-files", "", "Pipe-separated affected files (required)")
	cmd.Flags().StringVar(&createAffectedTests, "affected-tests", "", "Pipe-separated affected tests (required)")
	cmd.Flags().StringVar(&createImplSteps, "implementation-steps", "", "Pipe-separated implementation steps (required)")
	cmd.Flags().StringVar(&createLibraries, "library-needs", "", "Pipe-separated required libraries (optional)")
	cmd.Flags().BoolVarP(&createYes, "yes", "y", false, "Skip confirmation prompts")

	return cmd
}

func splitByPipe(s string) []string {
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, "|")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
