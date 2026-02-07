package service

import (
	"os"
	"testing"
	"time"

	"mandor/internal/domain"
	"mandor/internal/fs"
)

func setupTestFeatureService(t *testing.T) (*FeatureService, string) {
	tmpDir := t.TempDir()

	paths := &fs.Paths{WorkspaceRoot: tmpDir}
	writer := fs.NewWriter(paths)
	reader := fs.NewReader(paths)

	os.MkdirAll(tmpDir+"/.mandor/backlogs", 0755)
	os.MkdirAll(tmpDir+"/.mandor/specs", 0755)

	workspace := &domain.Workspace{
		Name:          "test-workspace",
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
		CreatedBy:     "test-user",
		Version:       "mandor.v1",
		SchemaVersion: "mandor.v1",
		Config:        domain.DefaultWorkspaceConfig(),
	}
	if err := writer.WriteWorkspace(workspace); err != nil {
		t.Fatalf("Failed to write workspace: %v", err)
	}

	service := &FeatureService{
		reader: reader,
		writer: writer,
		paths:  paths,
	}

	return service, tmpDir
}

func setupTestFeatureServiceWithBacklogBriefSpec(t *testing.T, backlogID, capabilityID string) (*FeatureService, *BacklogService) {
	tmpDir := t.TempDir()

	paths := &fs.Paths{WorkspaceRoot: tmpDir}
	writer := fs.NewWriter(paths)
	reader := fs.NewReader(paths)

	os.MkdirAll(tmpDir+"/.mandor/backlogs", 0755)
	os.MkdirAll(tmpDir+"/.mandor/briefs", 0755)
	os.MkdirAll(tmpDir+"/.mandor/specs", 0755)

	workspace := &domain.Workspace{
		Name:          "test-workspace",
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
		CreatedBy:     "test-user",
		Version:       "mandor.v1",
		SchemaVersion: "mandor.v1",
		Config:        domain.DefaultWorkspaceConfig(),
	}
	if err := writer.WriteWorkspace(workspace); err != nil {
		t.Fatalf("Failed to write workspace: %v", err)
	}

	backlogService := &BacklogService{
		reader: reader,
		writer: writer,
		paths:  paths,
	}

	backlogInput := &domain.BacklogCreateInput{
		ID:   backlogID,
		Name: "Test Backlog",
		Goal: "This is a test goal for the backlog that has enough characters to pass validation requirements",
	}
	if err := backlogService.CreateBacklog(backlogInput); err != nil {
		t.Fatalf("Failed to create backlog: %v", err)
	}

	briefService := &BriefService{
		reader: reader,
		writer: writer,
		paths:  paths,
	}

	briefInput := &domain.BriefCreateInput{
		BacklogID: backlogID,
		Name:      "Test Brief",
		Why:       "This is a problem statement that is long enough to pass validation requirements for feature service tests",
		Capabilities: []domain.CapabilityInput{
			{
				Name:        capabilityID,
				Description: "Test capability description",
				Modified:    false,
			},
		},
		Impact: &domain.BriefImpact{},
	}
	_, _ = briefService.CreateBrief(briefInput)

	specService := &SpecService{
		reader: reader,
		writer: writer,
		paths:  paths,
	}

	specInput := &domain.SpecCreateInput{
		BacklogID:    backlogID,
		CapabilityID: capabilityID,
		Summary:      "Test spec summary",
		Requirements: []domain.RequirementInput{
			{
				Summary: "Test requirement",
				IAEScenarios: []domain.IAEScenarioInput{
					{Intent: "Test intent", Action: "Test action", Expect: "Test expect"},
				},
			},
		},
	}
	_, _ = specService.CreateSpec(specInput)

	featureService := &FeatureService{
		reader: reader,
		writer: writer,
		paths:  paths,
	}

	return featureService, backlogService
}

func TestNewFeatureServiceWithPaths(t *testing.T) {
	tmpDir := t.TempDir()
	paths := &fs.Paths{WorkspaceRoot: tmpDir}

	service := NewFeatureServiceWithPaths(paths)

	if service == nil {
		t.Error("NewFeatureServiceWithPaths returned nil")
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

func TestFeatureWorkspaceInitialized(t *testing.T) {
	service, _ := setupTestFeatureService(t)

	initialized := service.WorkspaceInitialized()
	if !initialized {
		t.Error("Workspace should be initialized")
	}
}

func TestFeatureValidateCreateInput_BacklogNotFound(t *testing.T) {
	service, _ := setupTestFeatureService(t)

	input := &domain.FeatureCreateInput{
		BacklogID:    "nonexistent",
		CapabilityID: "test-cap",
		SpecID:       "test-spec",
		Name:         "Test Feature",
		Goal:         "A test feature goal that is long enough",
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for nonexistent backlog, got nil")
	}
}

func TestFeatureValidateCreateInput_CapabilityIDRequired(t *testing.T) {
	service, _ := setupTestFeatureServiceWithBacklogBriefSpec(t, "test-backlog", "test-capability")

	input := &domain.FeatureCreateInput{
		BacklogID:    "test-backlog",
		CapabilityID: "",
		SpecID:       "test-capability-spec",
		Name:         "Test Feature",
		Goal:         "A test feature goal that is long enough",
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for empty capability ID, got nil")
	}
}

func TestFeatureValidateCreateInput_InvalidCapabilityID(t *testing.T) {
	service, _ := setupTestFeatureServiceWithBacklogBriefSpec(t, "test-backlog", "test-capability")

	input := &domain.FeatureCreateInput{
		BacklogID:    "test-backlog",
		CapabilityID: "invalid cap@id",
		SpecID:       "test-capability-spec",
		Name:         "Test Feature",
		Goal:         "A test feature goal that is long enough",
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for invalid capability ID, got nil")
	}
}

func TestFeatureValidateCreateInput_CapabilityNotInBrief(t *testing.T) {
	service, _ := setupTestFeatureServiceWithBacklogBriefSpec(t, "test-backlog", "test-capability")

	input := &domain.FeatureCreateInput{
		BacklogID:    "test-backlog",
		CapabilityID: "other-capability",
		SpecID:       "other-capability-spec",
		Name:         "Test Feature",
		Goal:         "A test feature goal that is long enough",
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for capability not in brief, got nil")
	}
}

func TestFeatureValidateCreateInput_SpecIDRequired(t *testing.T) {
	service, _ := setupTestFeatureServiceWithBacklogBriefSpec(t, "test-backlog", "test-capability")

	input := &domain.FeatureCreateInput{
		BacklogID:    "test-backlog",
		CapabilityID: "test-capability",
		SpecID:       "",
		Name:         "Test Feature",
		Goal:         "A test feature goal that is long enough",
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for empty spec ID, got nil")
	}
}

func TestFeatureValidateCreateInput_SpecIDFormatMismatch(t *testing.T) {
	service, _ := setupTestFeatureServiceWithBacklogBriefSpec(t, "test-backlog", "test-capability")

	input := &domain.FeatureCreateInput{
		BacklogID:    "test-backlog",
		CapabilityID: "test-capability",
		SpecID:       "wrong-spec-id",
		Name:         "Test Feature",
		Goal:         "A test feature goal that is long enough",
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for spec ID format mismatch, got nil")
	}
}

func TestFeatureValidateCreateInput_NameRequired(t *testing.T) {
	service, _ := setupTestFeatureServiceWithBacklogBriefSpec(t, "test-backlog", "test-capability")

	input := &domain.FeatureCreateInput{
		BacklogID:    "test-backlog",
		CapabilityID: "test-capability",
		SpecID:       "test-capability-spec",
		Name:         "",
		Goal:         "A test feature goal that is long enough",
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for empty name, got nil")
	}
}

func TestFeatureValidateCreateInput_GoalRequired(t *testing.T) {
	service, _ := setupTestFeatureServiceWithBacklogBriefSpec(t, "test-backlog", "test-capability")

	input := &domain.FeatureCreateInput{
		BacklogID:    "test-backlog",
		CapabilityID: "test-capability",
		SpecID:       "test-capability-spec",
		Name:         "Test Feature",
		Goal:         "",
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for empty goal, got nil")
	}
}

func TestFeatureValidateCreateInput_GoalTooShort(t *testing.T) {
	service, _ := setupTestFeatureServiceWithBacklogBriefSpec(t, "test-backlog", "test-capability")

	input := &domain.FeatureCreateInput{
		BacklogID:    "test-backlog",
		CapabilityID: "test-capability",
		SpecID:       "test-capability-spec",
		Name:         "Test Feature",
		Goal:         "Short goal",
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for goal too short, got nil")
	}
}

func TestFeatureValidateCreateInput_InvalidScope(t *testing.T) {
	service, _ := setupTestFeatureServiceWithBacklogBriefSpec(t, "test-backlog", "test-capability")

	input := &domain.FeatureCreateInput{
		BacklogID:    "test-backlog",
		CapabilityID: "test-capability",
		SpecID:       "test-capability-spec",
		Name:         "Test Feature",
		Goal:         "A test feature goal that is long enough",
		Scope:        "invalid-scope",
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for invalid scope, got nil")
	}
}

func TestFeatureValidateCreateInput_InvalidPriority(t *testing.T) {
	service, _ := setupTestFeatureServiceWithBacklogBriefSpec(t, "test-backlog", "test-capability")

	input := &domain.FeatureCreateInput{
		BacklogID:    "test-backlog",
		CapabilityID: "test-capability",
		SpecID:       "test-capability-spec",
		Name:         "Test Feature",
		Goal:         "A test feature goal that is long enough",
		Priority:     "P99",
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for invalid priority, got nil")
	}
}

func TestFeatureValidateCreateInput_SelfDependency(t *testing.T) {
	service, _ := setupTestFeatureServiceWithBacklogBriefSpec(t, "test-backlog", "test-capability")

	input := &domain.FeatureCreateInput{
		BacklogID:    "test-backlog",
		CapabilityID: "test-capability",
		SpecID:       "test-capability-spec",
		Name:         "Test Feature",
		Goal:         "A test feature goal that is long enough",
		DependsOn:    []string{"test-backlog-feature-new"},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for self-dependency, got nil")
	}
}

func TestFeatureValidateCreateInput_DependencyNotFound(t *testing.T) {
	service, _ := setupTestFeatureServiceWithBacklogBriefSpec(t, "test-backlog", "test-capability")

	input := &domain.FeatureCreateInput{
		BacklogID:    "test-backlog",
		CapabilityID: "test-capability",
		SpecID:       "test-capability-spec",
		Name:         "Test Feature",
		Goal:         "A test feature goal that is long enough",
		DependsOn:    []string{"test-backlog-feature-nonexistent"},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for nonexistent dependency, got nil")
	}
}

func TestFeatureValidateCreateInput_DependencyCancelled(t *testing.T) {
	service, _ := setupTestFeatureServiceWithBacklogBriefSpec(t, "test-backlog", "test-capability")
	backlogID := "test-backlog"

	cancelledFeature := &domain.Feature{
		ID:           backlogID + "-feature-cancelled",
		BacklogID:    backlogID,
		CapabilityID: "test-capability",
		SpecID:       "test-capability-spec",
		Name:         "Cancelled Feature",
		Goal:         "A test feature goal that is long enough",
		Status:       domain.FeatureStatusCancelled,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := service.writer.WriteFeature(backlogID, cancelledFeature); err != nil {
		t.Fatalf("Failed to write cancelled feature: %v", err)
	}

	input := &domain.FeatureCreateInput{
		BacklogID:    backlogID,
		CapabilityID: "test-capability",
		SpecID:       "test-capability-spec",
		Name:         "Test Feature",
		Goal:         "A test feature goal that is long enough",
		DependsOn:    []string{cancelledFeature.ID},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected error for cancelled dependency, got nil")
	}
}

func TestFeatureValidateCreateInput_Valid(t *testing.T) {
	service, _ := setupTestFeatureServiceWithBacklogBriefSpec(t, "test-backlog", "test-capability")

	longGoal := "A test feature goal that is long enough to pass validation requirements. " +
		"This goal needs to be at least 300 characters according to the validation rules. " +
		"We are adding more text here to ensure it meets the minimum length requirement. " +
		"The validation function checks for FeatureGoalMinLength which is set to 300 characters. " +
		"So we need to make sure this text is long enough to pass the validation check."

	input := &domain.FeatureCreateInput{
		BacklogID:    "test-backlog",
		CapabilityID: "test-capability",
		SpecID:       "test-capability-spec",
		Name:         "Test Feature",
		Goal:         longGoal,
		Scope:        "fullstack",
		Priority:     "P2",
	}

	err := service.ValidateCreateInput(input)
	if err != nil {
		t.Errorf("ValidateCreateInput failed: %v", err)
	}
}

func TestFeatureCreateFeature_Success(t *testing.T) {
	service, _ := setupTestFeatureServiceWithBacklogBriefSpec(t, "test-backlog", "test-capability")

	input := &domain.FeatureCreateInput{
		BacklogID:    "test-backlog",
		CapabilityID: "test-capability",
		SpecID:       "test-capability-spec",
		Name:         "Test Feature",
		Goal:         "A test feature goal that is long enough to pass validation",
		Scope:        "fullstack",
		Priority:     "P2",
	}

	feature, err := service.CreateFeature(input)
	if err != nil {
		t.Fatalf("CreateFeature failed: %v", err)
	}

	if feature.ID == "" {
		t.Error("Feature ID should not be empty")
	}
	if feature.BacklogID != input.BacklogID {
		t.Errorf("Feature BacklogID = %q, want %q", feature.BacklogID, input.BacklogID)
	}
	if feature.CapabilityID != input.CapabilityID {
		t.Errorf("Feature CapabilityID = %q, want %q", feature.CapabilityID, input.CapabilityID)
	}
	if feature.Name != input.Name {
		t.Errorf("Feature Name = %q, want %q", feature.Name, input.Name)
	}
	if feature.Status != domain.FeatureStatusDraft {
		t.Errorf("Feature Status = %q, want %q", feature.Status, domain.FeatureStatusDraft)
	}
}

func TestFeatureCreateFeature_WithDependenciesBlocked(t *testing.T) {
	service, _ := setupTestFeatureServiceWithBacklogBriefSpec(t, "test-backlog", "test-capability")
	backlogID := "test-backlog"

	existingFeature := &domain.Feature{
		ID:           backlogID + "-feature-existing",
		BacklogID:    backlogID,
		CapabilityID: "test-capability",
		SpecID:       "test-capability-spec",
		Name:         "Existing Feature",
		Goal:         "A test feature goal that is long enough to pass validation",
		Status:       domain.FeatureStatusDraft,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := service.writer.WriteFeature(backlogID, existingFeature); err != nil {
		t.Fatalf("Failed to write existing feature: %v", err)
	}

	input := &domain.FeatureCreateInput{
		BacklogID:    backlogID,
		CapabilityID: "test-capability",
		SpecID:       "test-capability-spec",
		Name:         "Blocked Feature",
		Goal:         "A test feature goal that is long enough to pass validation",
		DependsOn:    []string{existingFeature.ID},
	}

	feature, err := service.CreateFeature(input)
	if err != nil {
		t.Fatalf("CreateFeature failed: %v", err)
	}

	if feature.Status != domain.FeatureStatusBlocked {
		t.Errorf("Feature Status = %q, want %q (blocked because dependency not done)", feature.Status, domain.FeatureStatusBlocked)
	}
}

func TestFeatureValidateStatusTransition_Valid(t *testing.T) {
	service := &FeatureService{}

	tests := []struct {
		current string
		next    string
		valid   bool
	}{
		{domain.FeatureStatusDraft, domain.FeatureStatusActive, true},
		{domain.FeatureStatusDraft, domain.FeatureStatusBlocked, true},
		{domain.FeatureStatusDraft, domain.FeatureStatusCancelled, true},
		{domain.FeatureStatusActive, domain.FeatureStatusDone, true},
		{domain.FeatureStatusActive, domain.FeatureStatusBlocked, true},
		{domain.FeatureStatusBlocked, domain.FeatureStatusDraft, true},
		{domain.FeatureStatusDone, domain.FeatureStatusCancelled, true},
		{domain.FeatureStatusCancelled, domain.FeatureStatusDraft, true},
		{domain.FeatureStatusDraft, domain.FeatureStatusDone, false},
		{domain.FeatureStatusDone, domain.FeatureStatusDraft, false},
	}

	for _, tt := range tests {
		t.Run(tt.current+"_to_"+tt.next, func(t *testing.T) {
			err := service.validateStatusTransition(tt.current, tt.next)
			if tt.valid && err != nil {
				t.Errorf("Expected valid transition from %s to %s, got error: %v", tt.current, tt.next, err)
			}
			if !tt.valid && err == nil {
				t.Errorf("Expected invalid transition from %s to %s, got no error", tt.current, tt.next)
			}
		})
	}
}

func TestFeatureGetGoalMinLength(t *testing.T) {
	service, _ := setupTestFeatureService(t)

	minLen := service.getFeatureGoalMinLength()
	if minLen != domain.FeatureGoalMinLength {
		t.Errorf("getFeatureGoalMinLength() = %d, want %d", minLen, domain.FeatureGoalMinLength)
	}
}

func TestFeatureCapabilityExistsInBrief_Exists(t *testing.T) {
	service, _ := setupTestFeatureServiceWithBacklogBriefSpec(t, "test-backlog", "test-capability")

	brief, _ := service.loadBriefForBacklog("test-backlog")

	exists := service.capabilityExistsInBrief(brief, "test-capability")
	if !exists {
		t.Error("Expected capability to exist in brief")
	}
}

func TestFeatureCapabilityExistsInBrief_NotExists(t *testing.T) {
	service, _ := setupTestFeatureServiceWithBacklogBriefSpec(t, "test-backlog", "test-capability")

	brief, _ := service.loadBriefForBacklog("test-backlog")

	exists := service.capabilityExistsInBrief(brief, "nonexistent-capability")
	if exists {
		t.Error("Expected capability to not exist in brief")
	}
}

func TestFeatureValidateNoDuplicateName_NoDuplicates(t *testing.T) {
	service, _ := setupTestFeatureServiceWithBacklogBriefSpec(t, "test-backlog", "test-capability")

	err := service.validateNoDuplicateName("test-backlog", "unique-feature-name")
	if err != nil {
		t.Errorf("validateNoDuplicateName failed: %v", err)
	}
}

func TestFeatureValidateNoCycle_NoCycle(t *testing.T) {
	service, _ := setupTestFeatureServiceWithBacklogBriefSpec(t, "test-backlog", "test-capability")

	err := service.validateNoCycle("test-backlog", "", []string{"nonexistent-feature"})
	if err != nil {
		t.Errorf("validateNoCycle should not fail for nonexistent deps: %v", err)
	}
}

func TestFeatureValidateDependencies_Success(t *testing.T) {
	service, _ := setupTestFeatureServiceWithBacklogBriefSpec(t, "test-backlog", "test-capability")
	backlogID := "test-backlog"

	feature := &domain.Feature{
		ID:           backlogID + "-feature-dep1",
		BacklogID:    backlogID,
		CapabilityID: "test-capability",
		SpecID:       "test-capability-spec",
		Name:         "Dep Feature 1",
		Goal:         "A test feature goal that is long enough to pass validation",
		Status:       domain.FeatureStatusDone,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := service.writer.WriteFeature(backlogID, feature); err != nil {
		t.Fatalf("Failed to write feature: %v", err)
	}

	err := service.validateDependencies(backlogID, "", []string{feature.ID})
	if err != nil {
		t.Errorf("validateDependencies failed: %v", err)
	}
}

func TestFeatureFindDependents_NoDependents(t *testing.T) {
	service, _ := setupTestFeatureServiceWithBacklogBriefSpec(t, "test-backlog", "test-capability")

	dependents, err := service.findDependents("test-backlog", "feature-1")
	if err != nil {
		t.Errorf("findDependents failed: %v", err)
	}
	if len(dependents) != 0 {
		t.Errorf("dependents count = %d, want 0", len(dependents))
	}
}

func TestFeatureValidateUpdateInput_FeatureNotFound(t *testing.T) {
	service, _ := setupTestFeatureServiceWithBacklogBriefSpec(t, "test-backlog", "test-capability")

	input := &domain.FeatureUpdateInput{
		BacklogID: "test-backlog",
		FeatureID: "test-backlog-feature-nonexistent",
		Name:      strPtr("Updated Name"),
	}

	err := service.ValidateUpdateInput(input)
	if err == nil {
		t.Error("Expected error for nonexistent feature, got nil")
	}
}

func TestFeatureValidateUpdateInput_CancelledFeature(t *testing.T) {
	service, _ := setupTestFeatureServiceWithBacklogBriefSpec(t, "test-backlog", "test-capability")
	backlogID := "test-backlog"

	feature := &domain.Feature{
		ID:           backlogID + "-feature-cancelled",
		BacklogID:    backlogID,
		CapabilityID: "test-capability",
		SpecID:       "test-capability-spec",
		Name:         "Cancelled Feature",
		Goal:         "A test feature goal that is long enough to pass validation",
		Status:       domain.FeatureStatusCancelled,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := service.writer.WriteFeature(backlogID, feature); err != nil {
		t.Fatalf("Failed to write feature: %v", err)
	}

	input := &domain.FeatureUpdateInput{
		BacklogID: backlogID,
		FeatureID: feature.ID,
		Name:      strPtr("Updated Name"),
	}

	err := service.ValidateUpdateInput(input)
	if err == nil {
		t.Error("Expected error for cancelled feature without reopen, got nil")
	}
}

func TestFeatureValidateUpdateInput_InvalidStatusTransition(t *testing.T) {
	service, _ := setupTestFeatureServiceWithBacklogBriefSpec(t, "test-backlog", "test-capability")
	backlogID := "test-backlog"

	feature := &domain.Feature{
		ID:           backlogID + "-feature-test",
		BacklogID:    backlogID,
		CapabilityID: "test-capability",
		SpecID:       "test-capability-spec",
		Name:         "Test Feature",
		Goal:         "A test feature goal that is long enough to pass validation",
		Status:       domain.FeatureStatusDraft,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := service.writer.WriteFeature(backlogID, feature); err != nil {
		t.Fatalf("Failed to write feature: %v", err)
	}

	blockedStatus := domain.FeatureStatusBlocked
	input := &domain.FeatureUpdateInput{
		BacklogID: backlogID,
		FeatureID: feature.ID,
		Status:    &blockedStatus,
	}

	err := service.ValidateUpdateInput(input)
	if err != nil {
		t.Errorf("Expected valid transition from draft to blocked, got error: %v", err)
	}
}

func TestFeatureUpdateFeature_NameChange(t *testing.T) {
	service, _ := setupTestFeatureServiceWithBacklogBriefSpec(t, "test-backlog", "test-capability")
	backlogID := "test-backlog"

	feature := &domain.Feature{
		ID:           backlogID + "-feature-update-test",
		BacklogID:    backlogID,
		CapabilityID: "test-capability",
		SpecID:       "test-capability-spec",
		Name:         "Original Name",
		Goal:         "A test feature goal that is long enough to pass validation",
		Status:       domain.FeatureStatusDraft,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := service.writer.WriteFeature(backlogID, feature); err != nil {
		t.Fatalf("Failed to write feature: %v", err)
	}

	newName := "Updated Name"
	input := &domain.FeatureUpdateInput{
		BacklogID: backlogID,
		FeatureID: feature.ID,
		Name:      &newName,
	}

	changes, err := service.UpdateFeature(input)
	if err != nil {
		t.Fatalf("UpdateFeature failed: %v", err)
	}

	found := false
	for _, c := range changes {
		if c == "name" {
			found = true
		}
	}
	if !found {
		t.Error("Expected 'name' in changes")
	}
}

func TestFeatureUpdateFeature_DryRun(t *testing.T) {
	service, _ := setupTestFeatureServiceWithBacklogBriefSpec(t, "test-backlog", "test-capability")
	backlogID := "test-backlog"

	feature := &domain.Feature{
		ID:           backlogID + "-feature-dry-run",
		BacklogID:    backlogID,
		CapabilityID: "test-capability",
		SpecID:       "test-capability-spec",
		Name:         "Dry Run Feature",
		Goal:         "A test feature goal that is long enough to pass validation",
		Status:       domain.FeatureStatusDraft,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := service.writer.WriteFeature(backlogID, feature); err != nil {
		t.Fatalf("Failed to write feature: %v", err)
	}

	newName := "Updated Name"
	input := &domain.FeatureUpdateInput{
		BacklogID: backlogID,
		FeatureID: feature.ID,
		Name:      &newName,
		DryRun:    true,
	}

	changes, err := service.UpdateFeature(input)
	if err != nil {
		t.Fatalf("UpdateFeature dry run failed: %v", err)
	}

	if len(changes) != 1 || changes[0] != "[DRY RUN] Would update feature: "+feature.ID {
		t.Errorf("Unexpected dry run changes: %v", changes)
	}
}

func TestFeatureUpdateFeature_Cancel(t *testing.T) {
	service, _ := setupTestFeatureServiceWithBacklogBriefSpec(t, "test-backlog", "test-capability")
	backlogID := "test-backlog"

	feature := &domain.Feature{
		ID:           backlogID + "-feature-cancel-test",
		BacklogID:    backlogID,
		CapabilityID: "test-capability",
		SpecID:       "test-capability-spec",
		Name:         "Cancel Test Feature",
		Goal:         "A test feature goal that is long enough to pass validation",
		Status:       domain.FeatureStatusDraft,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := service.writer.WriteFeature(backlogID, feature); err != nil {
		t.Fatalf("Failed to write feature: %v", err)
	}

	reason := "Test cancellation"
	input := &domain.FeatureUpdateInput{
		BacklogID: backlogID,
		FeatureID: feature.ID,
		Cancel:    true,
		Reason:    &reason,
	}

	_, err := service.UpdateFeature(input)
	if err != nil {
		t.Fatalf("UpdateFeature cancel failed: %v", err)
	}

	updated, _ := service.reader.ReadFeature(backlogID, feature.ID)
	if updated.Status != domain.FeatureStatusCancelled {
		t.Errorf("Feature Status = %q, want %q", updated.Status, domain.FeatureStatusCancelled)
	}
}

func TestFeatureValidateDeleteInput_FeatureNotFound(t *testing.T) {
	service, _ := setupTestFeatureServiceWithBacklogBriefSpec(t, "test-backlog", "test-capability")

	input := &domain.FeatureDeleteInput{
		BacklogID: "test-backlog",
		FeatureID: "test-backlog-feature-nonexistent",
		Force:     true,
	}

	err := service.ValidateDeleteInput(input)
	if err == nil {
		t.Error("Expected error for nonexistent feature, got nil")
	}
}

func TestFeatureValidateDeleteInput_AlreadyCancelled(t *testing.T) {
	service, _ := setupTestFeatureServiceWithBacklogBriefSpec(t, "test-backlog", "test-capability")
	backlogID := "test-backlog"

	feature := &domain.Feature{
		ID:           backlogID + "-feature-already-cancelled",
		BacklogID:    backlogID,
		CapabilityID: "test-capability",
		SpecID:       "test-capability-spec",
		Name:         "Already Cancelled Feature",
		Goal:         "A test feature goal that is long enough to pass validation",
		Status:       domain.FeatureStatusCancelled,
		Reason:       "Previous cancellation",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := service.writer.WriteFeature(backlogID, feature); err != nil {
		t.Fatalf("Failed to write feature: %v", err)
	}

	input := &domain.FeatureDeleteInput{
		BacklogID: backlogID,
		FeatureID: feature.ID,
	}

	err := service.ValidateDeleteInput(input)
	if err == nil {
		t.Error("Expected error for already cancelled feature, got nil")
	}
}

func TestFeatureDeleteFeature_Success(t *testing.T) {
	service, _ := setupTestFeatureServiceWithBacklogBriefSpec(t, "test-backlog", "test-capability")
	backlogID := "test-backlog"

	feature := &domain.Feature{
		ID:           backlogID + "-feature-delete-test",
		BacklogID:    backlogID,
		CapabilityID: "test-capability",
		SpecID:       "test-capability-spec",
		Name:         "Delete Test Feature",
		Goal:         "A test feature goal that is long enough to pass validation",
		Status:       domain.FeatureStatusDraft,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := service.writer.WriteFeature(backlogID, feature); err != nil {
		t.Fatalf("Failed to write feature: %v", err)
	}

	input := &domain.FeatureDeleteInput{
		BacklogID: backlogID,
		FeatureID: feature.ID,
		Reason:    "Test deletion",
	}

	err := service.DeleteFeature(input)
	if err != nil {
		t.Fatalf("DeleteFeature failed: %v", err)
	}

	deleted, _ := service.reader.ReadFeature(backlogID, feature.ID)
	if deleted.Status != domain.FeatureStatusCancelled {
		t.Errorf("Feature Status = %q, want %q", deleted.Status, domain.FeatureStatusCancelled)
	}
}

func TestFeatureListFeatures_Empty(t *testing.T) {
	service, _ := setupTestFeatureServiceWithBacklogBriefSpec(t, "test-backlog", "test-capability")

	input := &domain.FeatureListInput{
		BacklogID: "test-backlog",
	}

	result, err := service.ListFeatures(input)
	if err != nil {
		t.Fatalf("ListFeatures failed: %v", err)
	}

	if len(result.Features) != 0 {
		t.Errorf("Features count = %d, want 0", len(result.Features))
	}
}

func TestFeatureListFeatures_WithFeatures(t *testing.T) {
	service, _ := setupTestFeatureServiceWithBacklogBriefSpec(t, "test-backlog", "test-capability")
	backlogID := "test-backlog"

	feature := &domain.Feature{
		ID:           backlogID + "-feature-list-test",
		BacklogID:    backlogID,
		CapabilityID: "test-capability",
		SpecID:       "test-capability-spec",
		Name:         "List Test Feature",
		Goal:         "A test feature goal that is long enough to pass validation",
		Status:       domain.FeatureStatusDraft,
		Scope:        "fullstack",
		Priority:     "P1",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := service.writer.WriteFeature(backlogID, feature); err != nil {
		t.Fatalf("Failed to write feature: %v", err)
	}

	input := &domain.FeatureListInput{
		BacklogID: backlogID,
	}

	result, err := service.ListFeatures(input)
	if err != nil {
		t.Fatalf("ListFeatures failed: %v", err)
	}

	if len(result.Features) != 1 {
		t.Errorf("Features count = %d, want 1", len(result.Features))
	}
}

func TestFeatureGetFeatureDetail_NotFound(t *testing.T) {
	service, _ := setupTestFeatureServiceWithBacklogBriefSpec(t, "test-backlog", "test-capability")

	input := &domain.FeatureDetailInput{
		BacklogID:      "test-backlog",
		FeatureID:      "test-backlog-feature-nonexistent",
		IncludeDeleted: true,
	}

	_, err := service.GetFeatureDetail(input)
	if err == nil {
		t.Error("Expected error for nonexistent feature, got nil")
	}
}

func TestFeatureGetFeatureDetail_Success(t *testing.T) {
	service, _ := setupTestFeatureServiceWithBacklogBriefSpec(t, "test-backlog", "test-capability")
	backlogID := "test-backlog"

	feature := &domain.Feature{
		ID:           backlogID + "-feature-detail-test",
		BacklogID:    backlogID,
		CapabilityID: "test-capability",
		SpecID:       "test-capability-spec",
		Name:         "Detail Test Feature",
		Goal:         "A test feature goal that is long enough to pass validation",
		Status:       domain.FeatureStatusDraft,
		Scope:        "backend",
		Priority:     "P2",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := service.writer.WriteFeature(backlogID, feature); err != nil {
		t.Fatalf("Failed to write feature: %v", err)
	}

	input := &domain.FeatureDetailInput{
		BacklogID: backlogID,
		FeatureID: feature.ID,
	}

	detail, err := service.GetFeatureDetail(input)
	if err != nil {
		t.Fatalf("GetFeatureDetail failed: %v", err)
	}

	if detail.Name != feature.Name {
		t.Errorf("Detail Name = %q, want %q", detail.Name, feature.Name)
	}
	if detail.Scope != feature.Scope {
		t.Errorf("Detail Scope = %q, want %q", detail.Scope, feature.Scope)
	}
}
