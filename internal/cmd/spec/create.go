package spec

import (
	"fmt"
	"mandor/internal/domain"
	"mandor/internal/service"
	"mandor/internal/util"
	"strings"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create --capability <capability-id> --summary <text> --requirements <req1:intent1:action1:expect1>|<req2:...>",
	Short: "Create a new Spec",
	Long:  "Create a Specification for a Brief capability with minimum 1 requirement containing minimum 1 IAE scenario",
	RunE:  runCreate,
}

var (
	createCapability   string
	createSummary      string
	createRequirements string
)

func init() {
	createCmd.Flags().StringVar(&createCapability, "capability", "", "Capability ID from Brief (required)")
	createCmd.Flags().StringVar(&createSummary, "summary", "", "Brief spec description (required)")
	createCmd.Flags().StringVar(&createRequirements, "requirements", "", "Requirements in format 'summary:intent:action:expect|...' (required, min 1)")
	
	createCmd.MarkFlagRequired("capability")
	createCmd.MarkFlagRequired("summary")
	createCmd.MarkFlagRequired("requirements")
}

func runCreate(cmd *cobra.Command, args []string) error {
	svc, err := service.NewSpecService()
	if err != nil {
		return err
	}

	if !svc.WorkspaceInitialized() {
		return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
	}

	var projID string
	for p := cmd; p != nil; p = p.Parent() {
		if val, err := p.Flags().GetString("project"); err == nil && val != "" {
			projID = val
			break
		}
	}
	if projID == "" {
		return domain.NewValidationError("Project ID is required (--project).")
	}

	// Parse requirements
	if createRequirements == "" {
		return domain.NewValidationError("Error: at least one requirement is required\nSolution: use --requirements 'summary:intent:action:expect'")
	}

	requirements, err := parseRequirements(createRequirements)
	if err != nil {
		return err
	}

	input := &domain.SpecCreateInput{
		ProjectID:    projID,
		CapabilityID: createCapability,
		Summary:      createSummary,
		Requirements: requirements,
	}

	if err := svc.ValidateCreateInput(input); err != nil {
		return err
	}

	spec, err := svc.CreateSpec(input)
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "✓ Spec created: %s\n", spec.ID)
	fmt.Fprintf(out, "  Capability: %s\n", createCapability)
	fmt.Fprintf(out, "  Requirements: %d\n", len(spec.Requirements))
	for i, req := range spec.Requirements {
		fmt.Fprintf(out, "    Req %d: %s (scenarios: %d)\n", i+1, req.Summary, len(req.IAEScenarios))
	}

	_, warning := util.GetGitUsernameWithWarning()
	if warning != "" {
		fmt.Fprintln(out)
		fmt.Fprintln(out, warning)
		fmt.Fprintln(out, "  Run: git config user.name \"Your Name\"")
	}

	return nil
}

// parseRequirements parses "summary:intent:action:expect|..." format
func parseRequirements(input string) ([]domain.RequirementInput, error) {
	var requirements []domain.RequirementInput
	
	reqStrs := strings.Split(input, "|")
	for _, reqStr := range reqStrs {
		parts := strings.SplitN(reqStr, ":", 4)
		if len(parts) < 4 {
			return nil, fmt.Errorf("Error: Invalid requirement format '%s'\nSolution: use 'summary:intent:action:expect' or multiple separated by |", reqStr)
		}
		
		summary := strings.TrimSpace(parts[0])
		intent := strings.TrimSpace(parts[1])
		action := strings.TrimSpace(parts[2])
		expect := strings.TrimSpace(parts[3])
		
		if summary == "" || intent == "" || action == "" || expect == "" {
			return nil, fmt.Errorf("Error: Requirement fields cannot be empty\nSolution: ensure all fields (summary, intent, action, expect) have values")
		}
		
		scenario := domain.IAEScenarioInput{
			Intent: intent,
			Action: action,
			Expect: expect,
		}
		
		requirement := domain.RequirementInput{
			Summary:      summary,
			IAEScenarios: []domain.IAEScenarioInput{scenario},
		}
		
		requirements = append(requirements, requirement)
	}
	
	return requirements, nil
}

func GetCreateCmd() *cobra.Command {
	return createCmd
}
