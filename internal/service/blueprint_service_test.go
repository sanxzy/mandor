package service

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"mandor/internal/domain"
	"mandor/internal/fs"
)

func setupTestBlueprintService(t *testing.T) (*BlueprintService, string) {
	tmpDir := t.TempDir()

	paths := &fs.Paths{WorkspaceRoot: tmpDir}
	writer := fs.NewWriter(paths)
	reader := fs.NewReader(paths)

	os.MkdirAll(filepath.Join(tmpDir, ".mandor", "backlogs"), 0755)
	os.MkdirAll(filepath.Join(tmpDir, ".mandor", "briefs"), 0755)
	os.MkdirAll(filepath.Join(tmpDir, ".mandor", "features"), 0755)
	os.MkdirAll(filepath.Join(tmpDir, ".mandor", "specs"), 0755)

	workspace := &domain.Workspace{
		Name:          "test-workspace",
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
		CreatedBy:     "test-user",
		Config: domain.WorkspaceConfig{
			DefaultPriority: "P3",
			GoalLengths: domain.GoalLengths{
				Backlog: 500,
				Feature: 300,
				Task:    500,
				Issue:   200,
			},
		},
	}
	if err := writer.WriteWorkspace(workspace); err != nil {
		t.Fatalf("Failed to write workspace: %v", err)
	}

	service := &BlueprintService{
		reader: reader,
		writer: writer,
		paths:  paths,
	}

	return service, tmpDir
}

func setupTestBlueprintServiceWithBriefAndSpecs(t *testing.T, backlogID string) *BlueprintService {
	tmpDir := t.TempDir()

	paths := &fs.Paths{WorkspaceRoot: tmpDir}
	writer := fs.NewWriter(paths)
	reader := fs.NewReader(paths)

	os.MkdirAll(filepath.Join(tmpDir, ".mandor", "backlogs"), 0755)
	os.MkdirAll(filepath.Join(tmpDir, ".mandor", "briefs"), 0755)
	os.MkdirAll(filepath.Join(tmpDir, ".mandor", "features"), 0755)
	os.MkdirAll(filepath.Join(tmpDir, ".mandor", "specs"), 0755)

	workspace := &domain.Workspace{
		Name:          "test-workspace",
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
		CreatedBy:     "test-user",
		Config: domain.WorkspaceConfig{
			DefaultPriority: "P3",
		},
	}
	writer.WriteWorkspace(workspace)

	backlogService := &BacklogService{
		reader: reader,
		writer: writer,
		paths:  paths,
	}
	backlogInput := &domain.BacklogCreateInput{
		ID:   backlogID,
		Name: "Test Backlog",
		Goal: "This is a test goal for the backlog that has enough characters to pass validation",
	}
	backlogService.CreateBacklog(backlogInput)

	briefService := &BriefService{
		reader: reader,
		writer: writer,
		paths:  paths,
	}
	briefInput := &domain.BriefCreateInput{
		BacklogID: backlogID,
		Name:      "Test Brief",
		Why:       "This is a problem statement that is long enough to pass validation requirements for blueprint tests",
		Capabilities: []domain.CapabilityInput{
			{
				Name:        "new-feature",
				Description: "This is a new feature description",
				Modified:    false,
			},
			{
				Name:        "modified-feature",
				Description: "This is a modified feature description",
				Modified:    true,
			},
		},
		Impact: &domain.BriefImpact{},
	}
	briefService.CreateBrief(briefInput)

	specService := &SpecService{
		reader: reader,
		writer: writer,
		paths:  paths,
	}
	specInput1 := &domain.SpecCreateInput{
		BacklogID:    backlogID,
		CapabilityID: "new-feature",
		Summary:      "Spec for new feature",
		Requirements: []domain.RequirementInput{
			{
				Summary: "Test requirement",
				IAEScenarios: []domain.IAEScenarioInput{
					{Intent: "Test intent", Action: "Test action", Expect: "Test expect"},
				},
			},
		},
	}
	specService.CreateSpec(specInput1)

	specInput2 := &domain.SpecCreateInput{
		BacklogID:    backlogID,
		CapabilityID: "modified-feature",
		Summary:      "Spec for modified feature",
		Requirements: []domain.RequirementInput{
			{
				Summary: "Another requirement",
				IAEScenarios: []domain.IAEScenarioInput{
					{Intent: "Another intent", Action: "Another action", Expect: "Another expect"},
				},
			},
		},
	}
	specService.CreateSpec(specInput2)

	blueprintService := &BlueprintService{
		reader: reader,
		writer: writer,
		paths:  paths,
	}

	return blueprintService
}

func TestNewBlueprintServiceWithPaths(t *testing.T) {
	tmpDir := t.TempDir()
	paths := &fs.Paths{WorkspaceRoot: tmpDir}

	service := NewBlueprintServiceWithPaths(paths)

	if service == nil {
		t.Error("NewBlueprintServiceWithPaths returned nil")
	}
	if service.paths != paths {
		t.Error("Service paths not set correctly")
	}
	if service.reader == nil {
		t.Error("Service reader is nil")
	}
	if service.writer == nil {
		t.Error("Service writer is nil")
	}
}

func TestBlueprintWorkspaceInitialized(t *testing.T) {
	service, _ := setupTestBlueprintService(t)

	initialized := service.WorkspaceInitialized()
	if !initialized {
		t.Error("Workspace should be initialized")
	}
}

func TestValidateCreateInput_NoBacklog(t *testing.T) {
	service, _ := setupTestBlueprintService(t)

	input := &domain.BlueprintCreateInput{
		BacklogID:        "nonexistent-backlog",
		BriefID:          "test-brief",
		ProblemStatement: "This is a problem statement",
		ArchitectureDecisions: []domain.ArchitectureDecisionInput{
			{
				Title:     "Decision 1",
				Decision:  "The decision made",
				Rationale: "This is a rationale that is at least fifty characters long for validation",
			},
		},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for nonexistent backlog, got nil")
	}
}

func TestValidateCreateInput_EmptyBriefID(t *testing.T) {
	service, _ := setupTestBlueprintService(t)

	input := &domain.BlueprintCreateInput{
		BacklogID:        "test-backlog",
		BriefID:          "",
		ProblemStatement: "This is a problem statement",
		ArchitectureDecisions: []domain.ArchitectureDecisionInput{
			{
				Title:     "Decision 1",
				Decision:  "The decision made",
				Rationale: "This is a rationale that is at least fifty characters long for validation",
			},
		},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for empty brief ID, got nil")
	}
}

func TestValidateCreateInput_BriefNotFound(t *testing.T) {
	service, _ := setupTestBlueprintService(t)

	input := &domain.BlueprintCreateInput{
		BacklogID:        "test-backlog",
		BriefID:          "nonexistent-brief",
		ProblemStatement: "This is a problem statement",
		ArchitectureDecisions: []domain.ArchitectureDecisionInput{
			{
				Title:     "Decision 1",
				Decision:  "The decision made",
				Rationale: "This is a rationale that is at least fifty characters long for validation",
			},
		},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for nonexistent brief, got nil")
	}
}

func TestValidateCreateInput_NoSpecForCapability(t *testing.T) {
	tmpDir := t.TempDir()

	paths := &fs.Paths{WorkspaceRoot: tmpDir}
	writer := fs.NewWriter(paths)
	reader := fs.NewReader(paths)

	os.MkdirAll(filepath.Join(tmpDir, ".mandor", "backlogs"), 0755)
	os.MkdirAll(filepath.Join(tmpDir, ".mandor", "briefs"), 0755)

	workspace := &domain.Workspace{
		Name:      "test-workspace",
		CreatedAt: time.Now(),
	}
	writer.WriteWorkspace(workspace)

	backlogService := &BacklogService{reader: reader, writer: writer, paths: paths}
	backlogInput := &domain.BacklogCreateInput{
		ID:   "test-backlog",
		Name: "Test Backlog",
		Goal: "This is a test goal for the backlog that has enough characters to pass validation",
	}
	backlogService.CreateBacklog(backlogInput)

	briefService := &BriefService{reader: reader, writer: writer, paths: paths}
	briefInput := &domain.BriefCreateInput{
		BacklogID: "test-backlog",
		Name:      "Test Brief",
		Why:       "This is a problem statement that is long enough to pass validation requirements",
		Capabilities: []domain.CapabilityInput{
			{
				Name:        "test-capability",
				Description: "Test description",
				Modified:    false,
			},
		},
		Impact: &domain.BriefImpact{},
	}
	briefService.CreateBrief(briefInput)

	service := &BlueprintService{
		reader: reader,
		writer: writer,
		paths:  paths,
	}

	input := &domain.BlueprintCreateInput{
		BacklogID:        "test-backlog",
		BriefID:          "test-brief",
		ProblemStatement: "This is a problem statement",
		ArchitectureDecisions: []domain.ArchitectureDecisionInput{
			{
				Title:     "Decision 1",
				Decision:  "The decision made",
				Rationale: "This is a rationale that is at least fifty characters long for validation",
			},
		},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error when spec is missing for capability, got nil")
	}
}

func TestValidateCreateInput_InvalidSpec(t *testing.T) {
	tmpDir := t.TempDir()

	paths := &fs.Paths{WorkspaceRoot: tmpDir}
	writer := fs.NewWriter(paths)
	reader := fs.NewReader(paths)

	workspace := &domain.Workspace{
		Name:      "test-workspace",
		CreatedAt: time.Now(),
	}
	writer.WriteWorkspace(workspace)

	backlogService := &BacklogService{reader: reader, writer: writer, paths: paths}
	backlogInput := &domain.BacklogCreateInput{
		ID:   "test-backlog",
		Name: "Test Backlog",
		Goal: "This is a test goal for the backlog that has enough characters to pass validation",
	}
	backlogService.CreateBacklog(backlogInput)

	briefService := &BriefService{reader: reader, writer: writer, paths: paths}
	briefInput := &domain.BriefCreateInput{
		BacklogID: "test-backlog",
		Name:      "Test Brief",
		Why:       "This is a problem statement that is long enough to pass validation requirements",
		Capabilities: []domain.CapabilityInput{
			{
				Name:        "test-capability",
				Description: "Test description",
				Modified:    false,
			},
		},
		Impact: &domain.BriefImpact{},
	}
	briefService.CreateBrief(briefInput)

	os.MkdirAll(filepath.Join(tmpDir, ".mandor", "specs"), 0755)
	specPath := filepath.Join(tmpDir, ".mandor", "specs", "test-capability-spec.md")
	os.WriteFile(specPath, []byte("invalid content"), 0644)

	service := &BlueprintService{
		reader: reader,
		writer: writer,
		paths:  paths,
	}

	input := &domain.BlueprintCreateInput{
		BacklogID:        "test-backlog",
		BriefID:          "test-brief",
		ProblemStatement: "This is a problem statement",
		ArchitectureDecisions: []domain.ArchitectureDecisionInput{
			{
				Title:     "Decision 1",
				Decision:  "The decision made",
				Rationale: "This is a rationale that is at least fifty characters long for validation",
			},
		},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for invalid spec, got nil")
	}
}

func TestValidateCreateInput_EmptyProblemStatement(t *testing.T) {
	service := setupTestBlueprintServiceWithBriefAndSpecs(t, "test-backlog")

	input := &domain.BlueprintCreateInput{
		BacklogID:        "test-backlog",
		BriefID:          "test-brief",
		ProblemStatement: "",
		ArchitectureDecisions: []domain.ArchitectureDecisionInput{
			{
				Title:     "Decision 1",
				Decision:  "The decision made",
				Rationale: "This is a rationale that is at least fifty characters long for validation",
			},
		},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for empty problem statement, got nil")
	}
}

func TestValidateCreateInput_NoArchitectureDecisions(t *testing.T) {
	service := setupTestBlueprintServiceWithBriefAndSpecs(t, "test-backlog")

	input := &domain.BlueprintCreateInput{
		BacklogID:             "test-backlog",
		BriefID:               "test-brief",
		ProblemStatement:      "This is a problem statement",
		ArchitectureDecisions: []domain.ArchitectureDecisionInput{},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for no architecture decisions, got nil")
	}
}

func TestValidateCreateInput_EmptyDecisionTitle(t *testing.T) {
	service := setupTestBlueprintServiceWithBriefAndSpecs(t, "test-backlog")

	input := &domain.BlueprintCreateInput{
		BacklogID:        "test-backlog",
		BriefID:          "test-brief",
		ProblemStatement: "This is a problem statement",
		ArchitectureDecisions: []domain.ArchitectureDecisionInput{
			{
				Title:     "",
				Decision:  "The decision made",
				Rationale: "This is a rationale that is at least fifty characters long for validation",
			},
		},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for empty decision title, got nil")
	}
}

func TestValidateCreateInput_EmptyDecisionStatement(t *testing.T) {
	service := setupTestBlueprintServiceWithBriefAndSpecs(t, "test-backlog")

	input := &domain.BlueprintCreateInput{
		BacklogID:        "test-backlog",
		BriefID:          "test-brief",
		ProblemStatement: "This is a problem statement",
		ArchitectureDecisions: []domain.ArchitectureDecisionInput{
			{
				Title:     "Decision 1",
				Decision:  "",
				Rationale: "This is a rationale that is at least fifty characters long for validation",
			},
		},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for empty decision statement, got nil")
	}
}

func TestValidateCreateInput_ShortRationale(t *testing.T) {
	service := setupTestBlueprintServiceWithBriefAndSpecs(t, "test-backlog")

	input := &domain.BlueprintCreateInput{
		BacklogID:        "test-backlog",
		BriefID:          "test-brief",
		ProblemStatement: "This is a problem statement",
		ArchitectureDecisions: []domain.ArchitectureDecisionInput{
			{
				Title:     "Decision 1",
				Decision:  "The decision made",
				Rationale: "Too short",
			},
		},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for short rationale, got nil")
	}
}

func TestValidateCreateInput_FullValidInput(t *testing.T) {
	service := setupTestBlueprintServiceWithBriefAndSpecs(t, "test-backlog")

	input := &domain.BlueprintCreateInput{
		BacklogID:        "test-backlog",
		BriefID:          "test-brief",
		ProblemStatement: "This is a problem statement that describes the issue clearly",
		ArchitectureDecisions: []domain.ArchitectureDecisionInput{
			{
				Title:                  "Decision 1",
				Decision:               "The decision made for this architecture",
				Rationale:              "This is a rationale that is at least fifty characters long for validation purposes",
				AlternativesConsidered: []string{"Alternative 1", "Alternative 2"},
			},
		},
	}

	err := service.ValidateCreateInput(input)
	if err != nil {
		t.Errorf("Expected no error for valid input, got: %v", err)
	}
}

func TestCreateBlueprint_Success(t *testing.T) {
	service := setupTestBlueprintServiceWithBriefAndSpecs(t, "test-backlog")

	input := &domain.BlueprintCreateInput{
		BacklogID:        "test-backlog",
		BriefID:          "test-brief",
		ProblemStatement: "This is a problem statement that describes the issue clearly",
		Goals:            &domain.BlueprintGoals{},
		ArchitectureDecisions: []domain.ArchitectureDecisionInput{
			{
				Title:     "Decision 1",
				Decision:  "The decision made for this architecture",
				Rationale: "This is a rationale that is at least fifty characters long for validation purposes",
			},
		},
	}

	blueprint, err := service.CreateBlueprint(input)
	if err != nil {
		t.Fatalf("CreateBlueprint failed: %v", err)
	}

	if blueprint == nil {
		t.Error("Blueprint returned is nil")
	}
	if blueprint.ID != "test-backlog-blueprint" {
		t.Errorf("Blueprint.ID = %q, want %q", blueprint.ID, "test-backlog-blueprint")
	}
	if blueprint.BriefID != "test-brief" {
		t.Errorf("Blueprint.BriefID = %q, want %q", blueprint.BriefID, "test-brief")
	}
	if blueprint.BacklogID != "test-backlog" {
		t.Errorf("Blueprint.BacklogID = %q, want %q", blueprint.BacklogID, "test-backlog")
	}
	if blueprint.Status != domain.BlueprintStatusDraft {
		t.Errorf("Blueprint.Status = %q, want %q", blueprint.Status, domain.BlueprintStatusDraft)
	}
	if len(blueprint.ArchitectureDecisions) != 1 {
		t.Errorf("ArchitectureDecisions count = %d, want 1", len(blueprint.ArchitectureDecisions))
	}
	if blueprint.CreatedBy == "" {
		t.Error("Blueprint.CreatedBy is empty")
	}
}

func TestCreateBlueprint_WithRisks(t *testing.T) {
	service := setupTestBlueprintServiceWithBriefAndSpecs(t, "test-backlog")

	input := &domain.BlueprintCreateInput{
		BacklogID:        "test-backlog",
		BriefID:          "test-brief",
		ProblemStatement: "This is a problem statement that describes the issue clearly",
		Goals:            &domain.BlueprintGoals{},
		ArchitectureDecisions: []domain.ArchitectureDecisionInput{
			{
				Title:     "Decision 1",
				Decision:  "The decision made for this architecture",
				Rationale: "This is a rationale that is at least fifty characters long for validation purposes",
			},
		},
		Risks: []domain.RiskInput{
			{
				Description: "Risk description",
				Mitigation:  "Risk mitigation strategy",
			},
		},
	}

	blueprint, err := service.CreateBlueprint(input)
	if err != nil {
		t.Fatalf("CreateBlueprint failed: %v", err)
	}

	if len(blueprint.Risks) != 1 {
		t.Errorf("Risks count = %d, want 1", len(blueprint.Risks))
	}
}

func TestCreateBlueprint_WithConstraintsAndUserTypes(t *testing.T) {
	service := setupTestBlueprintServiceWithBriefAndSpecs(t, "test-backlog")

	input := &domain.BlueprintCreateInput{
		BacklogID:        "test-backlog",
		BriefID:          "test-brief",
		ProblemStatement: "This is a problem statement that describes the issue clearly",
		Goals:            &domain.BlueprintGoals{},
		Constraints:      []string{"Constraint 1", "Constraint 2"},
		UserTypes:        []string{"User Type 1", "User Type 2"},
		ArchitectureDecisions: []domain.ArchitectureDecisionInput{
			{
				Title:     "Decision 1",
				Decision:  "The decision made for this architecture",
				Rationale: "This is a rationale that is at least fifty characters long for validation purposes",
			},
		},
	}

	blueprint, err := service.CreateBlueprint(input)
	if err != nil {
		t.Fatalf("CreateBlueprint failed: %v", err)
	}

	if len(blueprint.Constraints) != 2 {
		t.Errorf("Constraints count = %d, want 2", len(blueprint.Constraints))
	}
	if len(blueprint.UserTypes) != 2 {
		t.Errorf("UserTypes count = %d, want 2", len(blueprint.UserTypes))
	}
}

func TestCreateBlueprint_WithGoals(t *testing.T) {
	service := setupTestBlueprintServiceWithBriefAndSpecs(t, "test-backlog")

	goals := &domain.BlueprintGoals{
		InScope:  []string{"Goal 1", "Goal 2"},
		OutScope: []string{"Out of scope 1"},
	}

	input := &domain.BlueprintCreateInput{
		BacklogID:        "test-backlog",
		BriefID:          "test-brief",
		ProblemStatement: "This is a problem statement that describes the issue clearly",
		Goals:            goals,
		ArchitectureDecisions: []domain.ArchitectureDecisionInput{
			{
				Title:     "Decision 1",
				Decision:  "The decision made for this architecture",
				Rationale: "This is a rationale that is at least fifty characters long for validation purposes",
			},
		},
	}

	blueprint, err := service.CreateBlueprint(input)
	if err != nil {
		t.Fatalf("CreateBlueprint failed: %v", err)
	}

	if len(blueprint.Goals.InScope) != 2 {
		t.Errorf("Goals.InScope count = %d, want 2", len(blueprint.Goals.InScope))
	}
	if len(blueprint.Goals.OutScope) != 1 {
		t.Errorf("Goals.OutScope count = %d, want 1", len(blueprint.Goals.OutScope))
	}
}

func TestCreateBlueprint_WithDataModels(t *testing.T) {
	service := setupTestBlueprintServiceWithBriefAndSpecs(t, "test-backlog")

	input := &domain.BlueprintCreateInput{
		BacklogID:        "test-backlog",
		BriefID:          "test-brief",
		ProblemStatement: "This is a problem statement that describes the issue clearly",
		Goals:            &domain.BlueprintGoals{},
		DataModels: []domain.DataModel{
			{
				Name:        "User",
				Description: "User data model",
				Fields: []domain.DataModelField{
					{Name: "ID", Type: "UUID", Required: true, Description: "User ID"},
					{Name: "Email", Type: "String", Required: true, Description: "User email"},
				},
			},
		},
		ArchitectureDecisions: []domain.ArchitectureDecisionInput{
			{
				Title:     "Decision 1",
				Decision:  "The decision made for this architecture",
				Rationale: "This is a rationale that is at least fifty characters long for validation purposes",
			},
		},
	}

	blueprint, err := service.CreateBlueprint(input)
	if err != nil {
		t.Fatalf("CreateBlueprint failed: %v", err)
	}

	if len(blueprint.DataModels) != 1 {
		t.Errorf("DataModels count = %d, want 1", len(blueprint.DataModels))
	}
	if len(blueprint.DataModels[0].Fields) != 2 {
		t.Errorf("DataModel.Fields count = %d, want 2", len(blueprint.DataModels[0].Fields))
	}
}

func TestReadBlueprint_Success(t *testing.T) {
	service := setupTestBlueprintServiceWithBriefAndSpecs(t, "test-backlog")

	createInput := &domain.BlueprintCreateInput{
		BacklogID:        "test-backlog",
		BriefID:          "test-brief",
		ProblemStatement: "This is a problem statement that describes the issue clearly",
		Goals:            &domain.BlueprintGoals{},
		ArchitectureDecisions: []domain.ArchitectureDecisionInput{
			{
				Title:     "Decision 1",
				Decision:  "The decision made for this architecture",
				Rationale: "This is a rationale that is at least fifty characters long for validation purposes",
			},
		},
	}
	service.CreateBlueprint(createInput)

	blueprint, err := service.ReadBlueprint("test-backlog")
	if err != nil {
		t.Fatalf("ReadBlueprint failed: %v", err)
	}

	if blueprint == nil {
		t.Error("Blueprint returned is nil")
	}
	if blueprint.ID != "test-backlog-blueprint" {
		t.Errorf("Blueprint.ID = %q, want %q", blueprint.ID, "test-backlog-blueprint")
	}
	if blueprint.ProblemStatement != "This is a problem statement that describes the issue clearly" {
		t.Error("Blueprint.ProblemStatement not preserved")
	}
}

func TestReadBlueprint_NotFound(t *testing.T) {
	service, _ := setupTestBlueprintService(t)

	blueprint, err := service.ReadBlueprint("nonexistent-backlog")
	if err == nil {
		t.Error("Expected error for nonexistent blueprint, got nil")
	}
	if blueprint != nil {
		t.Error("Blueprint should be nil when not found")
	}
}

func TestUpdateBlueprint_Success(t *testing.T) {
	service := setupTestBlueprintServiceWithBriefAndSpecs(t, "test-backlog")

	createInput := &domain.BlueprintCreateInput{
		BacklogID:        "test-backlog",
		BriefID:          "test-brief",
		ProblemStatement: "Original problem statement",
		Goals:            &domain.BlueprintGoals{},
		ArchitectureDecisions: []domain.ArchitectureDecisionInput{
			{
				Title:     "Decision 1",
				Decision:  "The decision made",
				Rationale: "This is a rationale that is at least fifty characters long for validation purposes",
			},
		},
	}
	blueprint, _ := service.CreateBlueprint(createInput)

	newProblemStatement := "Updated problem statement"
	blueprint.ProblemStatement = newProblemStatement
	err := service.UpdateBlueprint("test-backlog", blueprint)
	if err != nil {
		t.Fatalf("UpdateBlueprint failed: %v", err)
	}

	updatedBlueprint, _ := service.ReadBlueprint("test-backlog")
	if updatedBlueprint.ProblemStatement != newProblemStatement {
		t.Errorf("Blueprint.ProblemStatement = %q, want %q", updatedBlueprint.ProblemStatement, newProblemStatement)
	}
}

func TestUpdateBlueprint_UpdatesTimestamp(t *testing.T) {
	service := setupTestBlueprintServiceWithBriefAndSpecs(t, "test-backlog")

	createInput := &domain.BlueprintCreateInput{
		BacklogID:        "test-backlog",
		BriefID:          "test-brief",
		ProblemStatement: "This is a problem statement",
		Goals:            &domain.BlueprintGoals{},
		ArchitectureDecisions: []domain.ArchitectureDecisionInput{
			{
				Title:     "Decision 1",
				Decision:  "The decision made",
				Rationale: "This is a rationale that is at least fifty characters long for validation purposes",
			},
		},
	}
	blueprint, _ := service.CreateBlueprint(createInput)
	originalUpdatedAt := blueprint.UpdatedAt

	time.Sleep(10 * time.Millisecond)
	blueprint.ProblemStatement = "Updated problem statement"
	service.UpdateBlueprint("test-backlog", blueprint)

	updatedBlueprint, _ := service.ReadBlueprint("test-backlog")
	if updatedBlueprint.UpdatedAt.Equal(originalUpdatedAt) {
		t.Error("UpdatedAt timestamp was not updated")
	}
}

func TestDeleteBlueprint_Success(t *testing.T) {
	service := setupTestBlueprintServiceWithBriefAndSpecs(t, "test-backlog")

	createInput := &domain.BlueprintCreateInput{
		BacklogID:        "test-backlog",
		BriefID:          "test-brief",
		ProblemStatement: "This is a problem statement",
		Goals:            &domain.BlueprintGoals{},
		ArchitectureDecisions: []domain.ArchitectureDecisionInput{
			{
				Title:     "Decision 1",
				Decision:  "The decision made",
				Rationale: "This is a rationale that is at least fifty characters long for validation purposes",
			},
		},
	}
	service.CreateBlueprint(createInput)

	err := service.DeleteBlueprint("test-backlog")
	if err != nil {
		t.Fatalf("DeleteBlueprint failed: %v", err)
	}

	_, err = service.ReadBlueprint("test-backlog")
	if err == nil {
		t.Error("Expected error when reading deleted blueprint")
	}
}

func TestDeleteBlueprint_NonExistent(t *testing.T) {
	service, _ := setupTestBlueprintService(t)

	err := service.DeleteBlueprint("nonexistent-backlog")
	_ = err
}

func TestBlueprintFileExists(t *testing.T) {
	service := &BlueprintService{}

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "blueprint.md")
	os.WriteFile(path, []byte("Some content"), 0644)

	if !service.blueprintFileExists(path) {
		t.Error("blueprintFileExists returned false for existing file")
	}

	emptyPath := filepath.Join(tmpDir, "empty.md")
	os.WriteFile(emptyPath, []byte(""), 0644)

	if service.blueprintFileExists(emptyPath) {
		t.Error("blueprintFileExists returned true for empty file")
	}

	if service.blueprintFileExists(filepath.Join(tmpDir, "nonexistent.md")) {
		t.Error("blueprintFileExists returned true for nonexistent file")
	}
}

func TestBlueprintToMarkdown_ContainsAllSections(t *testing.T) {
	service := &BlueprintService{}

	blueprint := &domain.Blueprint{
		ID:               "test-backlog-blueprint",
		BriefID:          "test-brief",
		BacklogID:        "test-backlog",
		Status:           domain.BlueprintStatusDraft,
		Version:          "1.0",
		ProblemStatement: "Test problem statement",
		Constraints:      []string{"Constraint 1"},
		UserTypes:        []string{"User Type 1"},
		Goals: domain.BlueprintGoals{
			InScope:  []string{"Goal 1"},
			OutScope: []string{"Out of scope 1"},
		},
		ArchitectureDecisions: []domain.ArchitectureDecision{
			{
				ID:                     "decision-001",
				Title:                  "Decision 1",
				Decision:               "The decision made",
				Rationale:              "This is a rationale that is at least fifty characters long for display",
				AlternativesConsidered: []string{"Alternative 1"},
			},
		},
		DataModels: []domain.DataModel{
			{
				Name:        "User",
				Description: "User model",
				Fields: []domain.DataModelField{
					{Name: "ID", Type: "UUID", Required: true, Description: "User ID"},
				},
			},
		},
		ImplementationStrategy: "Implementation strategy here",
		Risks: []domain.Risk{
			{
				ID:          "risk-001",
				Description: "Risk description",
				Mitigation:  "Risk mitigation",
			},
		},
		CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CreatedBy: "test-user",
		UpdatedBy: "test-user",
	}

	md := service.blueprintToMarkdown(blueprint)

	if !contains(md, "BLUEPRINT_DATA_START") {
		t.Error("Markdown missing BLUEPRINT_DATA_START marker")
	}
	if !contains(md, "BLUEPRINT_DATA_END") {
		t.Error("Markdown missing BLUEPRINT_DATA_END marker")
	}
	if !contains(md, "# Blueprint: test-backlog-blueprint") {
		t.Error("Markdown missing blueprint title")
	}
	if !contains(md, "**Status:** draft") {
		t.Error("Markdown missing status")
	}
	if !contains(md, "## Problem Statement") {
		t.Error("Markdown missing Problem Statement section")
	}
	if !contains(md, "## Architecture Decisions") {
		t.Error("Markdown missing Architecture Decisions section")
	}
	if !contains(md, "## Constraints") {
		t.Error("Markdown missing Constraints section")
	}
	if !contains(md, "## User Types") {
		t.Error("Markdown missing User Types section")
	}
	if !contains(md, "## Goals") {
		t.Error("Markdown missing Goals section")
	}
}

func TestMarkdownToBlueprint_Success(t *testing.T) {
	service := &BlueprintService{}

	md := `<!--
BLUEPRINT_DATA_START
{"id":"test-blueprint","brief_id":"test-brief","backlog_id":"test-backlog","status":"draft","version":"1.0","problem_statement":"Test problem","constraints":["Constraint 1"],"user_types":["User 1"],"goals":{"in_scope":["Goal 1"],"out_scope":[]},"architecture_decisions":[{"id":"dec-001","title":"Decision 1","decision":"The decision","rationale":"This is a rationale that is at least fifty characters long","alternatives_considered":["Alt 1"]}],"data_models":[],"implementation_strategy":"","risks":[],"created_at":"2024-01-01T00:00:00Z","updated_at":"2024-01-01T00:00:00Z","created_by":"user","updated_by":"user"}
BLUEPRINT_DATA_END
-->
# Blueprint: test-blueprint

**Version:** 1.0
**Status:** draft

## Problem Statement

Test problem

## Architecture Decisions

### Decision 1: Decision 1 (dec-001)

**Decision:** The decision

**Rationale:** This is a rationale that is at least fifty characters long

**Alternatives Considered:**
- Alt 1

**Created:** 2024-01-01T00:00:00Z by user
**Updated:** 2024-01-01T00:00:00Z by user
`

	blueprint, err := service.markdownToBlueprint(md)
	if err != nil {
		t.Fatalf("markdownToBlueprint failed: %v", err)
	}

	if blueprint.ID != "test-blueprint" {
		t.Errorf("Blueprint.ID = %q, want %q", blueprint.ID, "test-blueprint")
	}
	if blueprint.ProblemStatement != "Test problem" {
		t.Errorf("Blueprint.ProblemStatement = %q, want %q", blueprint.ProblemStatement, "Test problem")
	}
	if len(blueprint.ArchitectureDecisions) != 1 {
		t.Errorf("ArchitectureDecisions count = %d, want 1", len(blueprint.ArchitectureDecisions))
	}
}

func TestMarkdownToBlueprint_MissingData(t *testing.T) {
	service := &BlueprintService{}

	tests := []struct {
		name string
		md   string
	}{
		{"Missing start marker", "# Blueprint: test\n\nContent"},
		{"Missing end marker", "<!--\n{}nBLUEPRINT_DATA_START"},
		{"Empty content", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.markdownToBlueprint(tt.md)
			if err == nil {
				t.Error("Expected error for malformed markdown, got nil")
			}
		})
	}
}

func TestMarkdownToBlueprint_InvalidJSON(t *testing.T) {
	service := &BlueprintService{}

	md := `<!--
BLUEPRINT_DATA_START
not valid json
BLUEPRINT_DATA_END
-->
# Blueprint: test
`

	_, err := service.markdownToBlueprint(md)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestBlueprintSerialization_Roundtrip(t *testing.T) {
	service := setupTestBlueprintServiceWithBriefAndSpecs(t, "test-backlog")

	input := &domain.BlueprintCreateInput{
		BacklogID:        "test-backlog",
		BriefID:          "test-brief",
		ProblemStatement: "Complex problem statement for roundtrip testing",
		Goals:            &domain.BlueprintGoals{},
		Constraints:      []string{"Constraint 1", "Constraint 2"},
		UserTypes:        []string{"User Type 1", "User Type 2"},
		ArchitectureDecisions: []domain.ArchitectureDecisionInput{
			{
				Title:                  "Decision 1",
				Decision:               "The decision made",
				Rationale:              "This is a rationale that is at least fifty characters long for validation purposes",
				AlternativesConsidered: []string{"Alternative 1", "Alternative 2"},
			},
			{
				Title:     "Decision 2",
				Decision:  "Another decision",
				Rationale: "Another rationale that is at least fifty characters long for testing",
			},
		},
	}

	created, err := service.CreateBlueprint(input)
	if err != nil {
		t.Fatalf("CreateBlueprint failed: %v", err)
	}

	read, err := service.ReadBlueprint("test-backlog")
	if err != nil {
		t.Fatalf("ReadBlueprint failed: %v", err)
	}

	if read.ID != created.ID {
		t.Errorf("ID mismatch: %q vs %q", read.ID, created.ID)
	}
	if read.ProblemStatement != created.ProblemStatement {
		t.Error("ProblemStatement not preserved")
	}
	if len(read.ArchitectureDecisions) != 2 {
		t.Errorf("ArchitectureDecisions count = %d, want 2", len(read.ArchitectureDecisions))
	}
	if len(read.Constraints) != 2 {
		t.Errorf("Constraints count = %d, want 2", len(read.Constraints))
	}
	if len(read.UserTypes) != 2 {
		t.Errorf("UserTypes count = %d, want 2", len(read.UserTypes))
	}
}

func TestVerifyAllSpecsExist(t *testing.T) {
	service := setupTestBlueprintServiceWithBriefAndSpecs(t, "test-backlog")

	brief, _ := service.loadBrief("test-backlog")

	err := service.verifyAllSpecsExist("test-backlog", brief)
	if err != nil {
		t.Errorf("verifyAllSpecsExist failed: %v", err)
	}
}

func TestVerifyAllSpecsValid(t *testing.T) {
	service := setupTestBlueprintServiceWithBriefAndSpecs(t, "test-backlog")

	brief, _ := service.loadBrief("test-backlog")

	err := service.verifyAllSpecsValid("test-backlog", brief)
	if err != nil {
		t.Errorf("verifyAllSpecsValid failed: %v", err)
	}
}

func TestLoadBrief(t *testing.T) {
	service := setupTestBlueprintServiceWithBriefAndSpecs(t, "test-backlog")

	brief, err := service.loadBrief("test-backlog")
	if err != nil {
		t.Fatalf("loadBrief failed: %v", err)
	}

	if brief == nil {
		t.Error("Brief returned is nil")
	}
	if brief.ID != "test-brief" {
		t.Errorf("Brief.ID = %q, want %q", brief.ID, "test-brief")
	}
}

func TestBlueprintWithMultipleArchitectureDecisions(t *testing.T) {
	service := setupTestBlueprintServiceWithBriefAndSpecs(t, "test-backlog")

	input := &domain.BlueprintCreateInput{
		BacklogID:        "test-backlog",
		BriefID:          "test-brief",
		ProblemStatement: "This is a problem statement",
		Goals:            &domain.BlueprintGoals{},
		ArchitectureDecisions: []domain.ArchitectureDecisionInput{
			{
				Title:     "Decision 1",
				Decision:  "First decision",
				Rationale: "This is a rationale that is at least fifty characters long for decision one",
			},
			{
				Title:     "Decision 2",
				Decision:  "Second decision",
				Rationale: "This is a rationale that is at least fifty characters long for decision two",
			},
			{
				Title:     "Decision 3",
				Decision:  "Third decision",
				Rationale: "This is a rationale that is at least fifty characters long for decision three",
			},
		},
	}

	blueprint, err := service.CreateBlueprint(input)
	if err != nil {
		t.Fatalf("CreateBlueprint failed: %v", err)
	}

	if len(blueprint.ArchitectureDecisions) != 3 {
		t.Errorf("ArchitectureDecisions count = %d, want 3", len(blueprint.ArchitectureDecisions))
	}
}

func TestBlueprintWithImplementationStrategy(t *testing.T) {
	service := setupTestBlueprintServiceWithBriefAndSpecs(t, "test-backlog")

	input := &domain.BlueprintCreateInput{
		BacklogID:              "test-backlog",
		BriefID:                "test-brief",
		ProblemStatement:       "This is a problem statement",
		Goals:                  &domain.BlueprintGoals{},
		ImplementationStrategy: "This is a detailed implementation strategy for the project",
		ArchitectureDecisions: []domain.ArchitectureDecisionInput{
			{
				Title:     "Decision 1",
				Decision:  "The decision made",
				Rationale: "This is a rationale that is at least fifty characters long for validation purposes",
			},
		},
	}

	blueprint, err := service.CreateBlueprint(input)
	if err != nil {
		t.Fatalf("CreateBlueprint failed: %v", err)
	}

	if blueprint.ImplementationStrategy != "This is a detailed implementation strategy for the project" {
		t.Error("ImplementationStrategy not preserved")
	}
}
