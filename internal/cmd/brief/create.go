package brief

import (
	"fmt"
	"mandor/internal/domain"
	"mandor/internal/service"
	"mandor/internal/util"
	"strings"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create --name <name> --why <text> --capabilities <cap1:desc1>|<cap2:desc2>",
	Short: "Create a new Brief",
	Long:  "Create a Brief that defines the overall intent and scope for a project",
	RunE:  runCreate,
}

var (
	createName          string
	createWhy           string
	createCapabilities  string
	createImpactStack   string
	createImpactSystems string
	createDependencies  string
)

func init() {
	createCmd.Flags().StringVar(&createName, "name", "", "Brief name (required)")
	createCmd.Flags().StringVar(&createWhy, "why", "", "Problem statement and motivation (required, 100-5000 chars)")
	createCmd.Flags().StringVar(&createCapabilities, "capabilities", "", "New capabilities (format: 'name1:desc1|name2:desc2', required)")
	createCmd.Flags().StringVar(&createImpactStack, "tech-stack", "", "Technical stack choices (comma-separated)")
	createCmd.Flags().StringVar(&createImpactSystems, "affected-systems", "", "Affected systems (comma-separated)")
	createCmd.Flags().StringVar(&createDependencies, "dependencies", "", "External dependencies (comma-separated)")

	createCmd.MarkFlagRequired("name")
	createCmd.MarkFlagRequired("why")
	createCmd.MarkFlagRequired("capabilities")
}

func runCreate(cmd *cobra.Command, args []string) error {
	svc, err := service.NewBriefService()
	if err != nil {
		return err
	}

	if !svc.WorkspaceInitialized() {
		return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
	}

	var backlogID string
	for p := cmd; p != nil; p = p.Parent() {
		if val, err := p.Flags().GetString("backlog"); err == nil && val != "" {
			backlogID = val
			break
		}
	}
	if backlogID == "" {
		return domain.NewValidationError("Backlog ID is required (--backlog).")
	}

	// Parse capabilities
	if createCapabilities == "" {
		return domain.NewValidationError("Error: at least one capability is required\nSolution: use --capabilities 'name1:desc1' or 'name1:desc1|name2:desc2'")
	}

	capabilities, err := parseCapabilities(createCapabilities)
	if err != nil {
		return err
	}

	// Parse impact if provided
	impact := &domain.BriefImpact{
		TechnicalStack:  parseCSV(createImpactStack),
		AffectedSystems: parseCSV(createImpactSystems),
		Dependencies:    parseCSV(createDependencies),
	}

	input := &domain.BriefCreateInput{
		BacklogID:    backlogID,
		Name:         createName,
		Why:          createWhy,
		Capabilities: capabilities,
		Impact:       impact,
	}

	if err := svc.ValidateCreateInput(input); err != nil {
		return err
	}

	brief, err := svc.CreateBrief(input)
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "✓ Brief created: %s\n", brief.ID)
	fmt.Fprintf(out, "  ID: %s\n", brief.ID)
	fmt.Fprintf(out, "  Status: %s\n", brief.Status)
	fmt.Fprintf(out, "  Capabilities: %d\n", len(brief.NewCapabilities))

	// AI-friendly next steps
	fmt.Fprintln(out)
	fmt.Fprintln(out, "  ─────────────────────────────────────────────────────────")
	fmt.Fprintln(out, "  NEXT STEPS:")
	fmt.Fprintln(out, "  ─────────────────────────────────────────────────────────")
	fmt.Fprintf(out, "  • Create specs for capabilities: mandor spec create --backlog %s --capability <cap-id> ...\n", backlogID)
	fmt.Fprintf(out, "  • View brief details: mandor brief read %s --backlog %s\n", brief.ID, backlogID)
	fmt.Fprintf(out, "  • List all briefs: mandor brief list --backlog %s\n", backlogID)

	_, warning := util.GetGitUsernameWithWarning()
	if warning != "" {
		fmt.Fprintln(out)
		fmt.Fprintln(out, warning)
		fmt.Fprintln(out, "  Run: git config user.name \"Your Name\"")
	}

	return nil
}

// parseCapabilities parses "name:desc|name:desc" format into CapabilityInput
func parseCapabilities(input string) ([]domain.CapabilityInput, error) {
	var capabilities []domain.CapabilityInput

	caps := strings.Split(input, "|")
	for _, cap := range caps {
		parts := strings.SplitN(cap, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("Error: Invalid capability format '%s'\nSolution: use 'name:description' or 'name1:desc1|name2:desc2'", cap)
		}

		name := strings.TrimSpace(parts[0])
		desc := strings.TrimSpace(parts[1])

		if name == "" || desc == "" {
			return nil, fmt.Errorf("Error: capability name and description cannot be empty\nSolution: ensure format is 'name:description'")
		}

		capID := util.ToSlug(name)
		if !domain.ValidateCapabilityID(capID) {
			return nil, fmt.Errorf("Error: Invalid capability name '%s'\nSolution: use alphanumeric characters and hyphens only", name)
		}

		capabilities = append(capabilities, domain.CapabilityInput{
			Name:        name,
			Description: desc,
			Modified:    false,
		})
	}

	return capabilities, nil
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
