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
	analyzeEntityType string
	analyzeEntityID   string
	analyzeFields     string
	analyzeBacklogID  string
	analyzeJSON       bool
)

// NewAnalyzeCmd creates the analyze subcommand
func NewAnalyzeCmd(paths *fs.Paths) *cobra.Command {
	analyzeCmd := &cobra.Command{
		Use:   "analyze",
		Short: "Analyze impact of a proposed change",
		Long:  "Analyze the impact of changing specific fields in a Brief, Spec, or Blueprint",
		RunE: func(c *cobra.Command, args []string) error {
			return runAnalyze(paths)
		},
	}

	analyzeCmd.Flags().StringVar(&analyzeEntityType, "entity-type", "", "Type of entity: brief, spec, or blueprint")
	analyzeCmd.Flags().StringVar(&analyzeEntityID, "entity-id", "", "ID of the entity to analyze")
	analyzeCmd.Flags().StringVar(&analyzeFields, "fields", "", "Comma-separated fields that will change")
	analyzeCmd.Flags().StringVar(&analyzeBacklogID, "backlog", "", "Backlog ID (optional)")
	analyzeCmd.Flags().BoolVar(&analyzeJSON, "json", false, "Output as JSON")

	analyzeCmd.MarkFlagRequired("entity-type")
	analyzeCmd.MarkFlagRequired("entity-id")
	analyzeCmd.MarkFlagRequired("fields")

	return analyzeCmd
}

func runAnalyze(paths *fs.Paths) error {
	// Convert comma-separated fields to list
	fieldsList := []string{}
	if analyzeFields != "" {
		for _, f := range strings.Split(analyzeFields, ",") {
			f = strings.TrimSpace(f)
			if f != "" {
				fieldsList = append(fieldsList, f)
			}
		}
	}

	if len(fieldsList) == 0 {
		return fmt.Errorf("at least one field must be specified")
	}

	// Convert field list to map
	fieldsMap := make(map[string]interface{})
	for _, field := range fieldsList {
		fieldsMap[field] = true
	}

	// Perform analysis based on entity type
	governanceService := service.NewChangeGovernanceService(paths)
	var analysis *service.ChangeImpactAnalysis
	var err error

	switch strings.ToLower(analyzeEntityType) {
	case "brief":
		analysis, err = governanceService.ValidateBriefChangeBlocking(analyzeBacklogID, analyzeEntityID, fieldsMap)
	case "spec":
		analysis, err = governanceService.ValidateSpecChangeBlocking(analyzeBacklogID, analyzeEntityID, fieldsMap)
	case "blueprint":
		analysis, err = governanceService.ValidateBlueprintChangeBlocking(analyzeBacklogID, analyzeEntityID, fieldsMap)
	default:
		return fmt.Errorf("invalid entity type: %s (must be brief, spec, or blueprint)", analyzeEntityType)
	}

	if err != nil {
		return err
	}

	// Persist analysis
	if err := governanceService.PersistChangeAnalysis(analyzeBacklogID, analysis); err != nil {
		return fmt.Errorf("failed to persist analysis: %w", err)
	}

	if analyzeJSON {
		response := map[string]interface{}{
			"analysis":  analysis,
			"message":   fmt.Sprintf("Change analysis complete. Blocking status: %s", analysis.BlockingStatus),
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		}
		data, _ := json.MarshalIndent(response, "", "  ")
		fmt.Println(string(data))
	} else {
		fmt.Println(printAnalysis(analysis))
		fmt.Printf("Message: Change analysis complete. Blocking status: %s\n", analysis.BlockingStatus)
	}

	return nil
}

func printAnalysis(analysis *service.ChangeImpactAnalysis) string {
	var output strings.Builder

	output.WriteString("\n================================\n")
	output.WriteString("CHANGE IMPACT ANALYSIS\n")
	output.WriteString("================================\n\n")

	output.WriteString(fmt.Sprintf("Change ID:        %s\n", analysis.ChangeID))
	output.WriteString(fmt.Sprintf("Change Type:      %s\n", analysis.ChangeType))
	output.WriteString(fmt.Sprintf("Entity ID:        %s\n", analysis.EntityID))
	output.WriteString(fmt.Sprintf("Backlog ID:      %s\n", analysis.BacklogID))
	output.WriteString(fmt.Sprintf("Status:           %s\n", analysis.Status))
	output.WriteString(fmt.Sprintf("Blocking Status:  %s\n", analysis.BlockingStatus))
	output.WriteString(fmt.Sprintf("Timestamp:        %s\n", analysis.Timestamp))

	output.WriteString("\nFields Changed:\n")
	for _, field := range analysis.FieldsChanged {
		output.WriteString(fmt.Sprintf("  - %s\n", field))
	}

	output.WriteString(fmt.Sprintf("\nImpacted Entities:\n"))
	output.WriteString(fmt.Sprintf("  Specs:       %d\n", len(analysis.ImpactedSpecs)))
	output.WriteString(fmt.Sprintf("  Features:    %d\n", len(analysis.ImpactedFeatures)))
	output.WriteString(fmt.Sprintf("  Tasks:       %d\n", len(analysis.ImpactedTasks)))
	output.WriteString(fmt.Sprintf("  Blueprints:  %d\n", len(analysis.ImpactedBlueprints)))

	if len(analysis.RequiredActions) > 0 {
		output.WriteString("\nRequired Actions:\n")
		for i, action := range analysis.RequiredActions {
			output.WriteString(fmt.Sprintf("  %d. %s\n", i+1, action))
		}
	}

	output.WriteString(fmt.Sprintf("\nValidation Deadline: %s\n", analysis.ValidationDeadline))
	output.WriteString("================================\n\n")

	return output.String()
}
