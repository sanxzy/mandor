package blueprint

import (
	"fmt"
	"mandor/internal/domain"
	"mandor/internal/service"
	"mandor/internal/util"
	"strings"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create --brief <brief-id> --problem <text> --decisions <decision1>|<decision2>",
	Short: "Create a new Blueprint",
	Long:  "Create a Blueprint (technical architecture) - requires all Brief capabilities to have valid Specs",
	RunE:  runCreate,
}

var (
	createBrief           string
	createProblem         string
	createDecisions       string
	createConstraints     string
	createUserTypes       string
	createGoalsInScope    string
	createGoalsOutScope   string
	createImplementation  string
	createRisks           string
)

func init() {
	createCmd.Flags().StringVar(&createBrief, "brief", "", "Brief ID (required)")
	createCmd.Flags().StringVar(&createProblem, "problem", "", "Problem statement (required)")
	createCmd.Flags().StringVar(&createDecisions, "decisions", "", "Architecture decisions - format: 'title:rationale|...' (required, min 1)")
	createCmd.Flags().StringVar(&createConstraints, "constraints", "", "Constraints (comma-separated)")
	createCmd.Flags().StringVar(&createUserTypes, "user-types", "", "User types (comma-separated)")
	createCmd.Flags().StringVar(&createGoalsInScope, "goals-in-scope", "", "In-scope goals (comma-separated)")
	createCmd.Flags().StringVar(&createGoalsOutScope, "goals-out-scope", "", "Out-of-scope goals (comma-separated)")
	createCmd.Flags().StringVar(&createImplementation, "implementation", "", "Implementation strategy")
	createCmd.Flags().StringVar(&createRisks, "risks", "", "Risks - format: 'description:mitigation|...'")
	
	createCmd.MarkFlagRequired("brief")
	createCmd.MarkFlagRequired("problem")
	createCmd.MarkFlagRequired("decisions")
}

func runCreate(cmd *cobra.Command, args []string) error {
	svc, err := service.NewBlueprintService()
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

	// Parse architecture decisions (min 1 required)
	if createDecisions == "" {
		return domain.NewValidationError("Error: at least one architecture decision is required\nSolution: use --decisions 'title:rationale' with rationale min 50 chars")
	}

	decisions, err := parseDecisions(createDecisions)
	if err != nil {
		return err
	}

	if len(decisions) == 0 {
		return domain.NewValidationError("Error: at least one valid architecture decision is required\nSolution: each decision needs title and rationale (min 50 chars)")
	}

	// Parse risks
	risks, err := parseRisks(createRisks)
	if err != nil {
		return err
	}

	// Create input
	input := &domain.BlueprintCreateInput{
		ProjectID:             projID,
		BriefID:               createBrief,
		ProblemStatement:      createProblem,
		Constraints:           parseCSV(createConstraints),
		UserTypes:             parseCSV(createUserTypes),
		Goals: &domain.BlueprintGoals{
			InScope:  parseCSV(createGoalsInScope),
			OutScope: parseCSV(createGoalsOutScope),
		},
		ArchitectureDecisions: decisions,
		DataModels:            []domain.DataModel{},
		ImplementationStrategy: createImplementation,
		Risks:                 risks,
	}

	if err := svc.ValidateCreateInput(input); err != nil {
		return err
	}

	blueprint, err := svc.CreateBlueprint(input)
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "✓ Blueprint created: %s\n", blueprint.ID)
	fmt.Fprintf(out, "  Brief: %s\n", createBrief)
	fmt.Fprintf(out, "  Architecture Decisions: %d\n", len(blueprint.ArchitectureDecisions))

	_, warning := util.GetGitUsernameWithWarning()
	if warning != "" {
		fmt.Fprintln(out)
		fmt.Fprintln(out, warning)
		fmt.Fprintln(out, "  Run: git config user.name \"Your Name\"")
	}

	return nil
}

// parseDecisions parses "title:decision:rationale|..." format (rationale must be min 50 chars)
func parseDecisions(input string) ([]domain.ArchitectureDecisionInput, error) {
	var decisions []domain.ArchitectureDecisionInput

	decStrs := strings.Split(input, "|")
	for _, decStr := range decStrs {
		// Try to parse as title:rationale first (original format)
		parts := strings.SplitN(decStr, ":", 2)
		if len(parts) < 2 {
			return nil, fmt.Errorf("Error: Invalid decision format '%s'\nSolution: use 'title:rationale' with rationale min 50 chars", decStr)
		}

		title := strings.TrimSpace(parts[0])
		rationale := strings.TrimSpace(parts[1])

		if title == "" || rationale == "" {
			return nil, fmt.Errorf("Error: Decision title and rationale cannot be empty\nSolution: ensure format is 'title:rationale'")
		}

		if len(rationale) < 50 {
			return nil, fmt.Errorf("Error: Decision rationale must be at least 50 characters (got %d)\nSolution: expand your rationale", len(rationale))
		}

		decisions = append(decisions, domain.ArchitectureDecisionInput{
			Title:     title,
			Decision:  title, // Use title as decision for now
			Rationale: rationale,
		})
	}

	return decisions, nil
}

// parseRisks parses "description:mitigation|..." format
func parseRisks(input string) ([]domain.RiskInput, error) {
	if input == "" {
		return []domain.RiskInput{}, nil
	}

	var risks []domain.RiskInput

	riskStrs := strings.Split(input, "|")
	for _, riskStr := range riskStrs {
		parts := strings.SplitN(riskStr, ":", 2)
		if len(parts) < 2 {
			return nil, fmt.Errorf("Error: Invalid risk format '%s'\nSolution: use 'description:mitigation'", riskStr)
		}

		description := strings.TrimSpace(parts[0])
		mitigation := strings.TrimSpace(parts[1])

		if description == "" || mitigation == "" {
			return nil, fmt.Errorf("Error: Risk description and mitigation cannot be empty\nSolution: ensure format is 'description:mitigation'")
		}

		risks = append(risks, domain.RiskInput{
			Description: description,
			Mitigation:  mitigation,
		})
	}

	return risks, nil
}

// parseCSV parses comma-separated values
func parseCSV(input string) []string {
	if input == "" {
		return []string{}
	}

	var result []string
	parts := strings.Split(input, ",")
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func GetCreateCmd() *cobra.Command {
	return createCmd
}
