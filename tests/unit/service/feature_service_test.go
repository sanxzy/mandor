package service_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"mandor/internal/domain"
	"mandor/internal/fs"
	"mandor/internal/service"
)

func setupTestFeatureService(t *testing.T) (*service.FeatureService, string) {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "mandor-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	mandorDir := filepath.Join(tmpDir, ".mandor")
	projectsDir := filepath.Join(mandorDir, "projects")

	if err := os.MkdirAll(projectsDir, 0755); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to create .mandor directory: %v", err)
	}

	workspacePath := filepath.Join(mandorDir, "workspace.json")
	wsData := []byte("{\"version\":\"v1.0.0\"}")
	if err := os.WriteFile(workspacePath, wsData, 0644); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to write workspace.json: %v", err)
	}

	paths, err := fs.NewPathsFromRoot(tmpDir)
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to create paths: %v", err)
	}

	svc := service.NewFeatureServiceWithPaths(paths)

	return svc, tmpDir
}

func writeTestProjectForFeature(t *testing.T, tmpDir, projectID string, status string) {
	t.Helper()

	paths, err := fs.NewPathsFromRoot(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create paths: %v", err)
	}
	writer := fs.NewWriter(paths)

	if err := writer.CreateProjectDir(projectID); err != nil {
		t.Fatalf("Failed to create project directory: %v", err)
	}

	project := &domain.Project{
		ID:        projectID,
		Name:      "Test Project",
		Goal:      string(make([]byte, 500)),
		Status:    status,
		Strict:    false,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		CreatedBy: "testuser",
		UpdatedBy: "testuser",
	}

	if err := writer.WriteProjectMetadata(projectID, project); err != nil {
		t.Fatalf("Failed to write project metadata: %v", err)
	}

	schema := domain.DefaultProjectSchema("same_project_only", "cross_project_allowed", "same_project_only")
	if err := writer.WriteProjectSchema(projectID, &schema); err != nil {
		t.Fatalf("Failed to write project schema: %v", err)
	}
}

func writeTestFeature(t *testing.T, tmpDir, projectID, featureID string, status string, dependsOn []string) {
	t.Helper()

	paths, err := fs.NewPathsFromRoot(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create paths: %v", err)
	}
	writer := fs.NewWriter(paths)

	feature := &domain.Feature{
		ID:        featureID,
		ProjectID: projectID,
		Name:      "Test Feature",
		Goal:      "Test goal for feature",
		Scope:     "fullstack",
		Priority:  "P3",
		Status:    status,
		DependsOn: dependsOn,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		CreatedBy: "testuser",
		UpdatedBy: "testuser",
	}

	if err := writer.WriteFeature(projectID, feature); err != nil {
		t.Fatalf("Failed to write feature: %v", err)
	}
}

func TestFeatureWorkspaceInitialized(t *testing.T) {
	svc, tmpDir := setupTestFeatureService(t)
	defer os.RemoveAll(tmpDir)

	if !svc.WorkspaceInitialized() {
		t.Error("Expected workspace to be initialized")
	}
}

func TestFeatureValidateCreateInput_Valid(t *testing.T) {
	svc, tmpDir := setupTestFeatureService(t)
	defer os.RemoveAll(tmpDir)

	writeTestProjectForFeature(t, tmpDir, "testproject", domain.ProjectStatusInitial)

	input := &domain.FeatureCreateInput{
		ProjectID: "testproject",
		Name:      "Test Feature",
		Goal:      "Test goal for feature",
		Scope:     "fullstack",
		Priority:  "P3",
		DependsOn: nil,
	}

	err := svc.ValidateCreateInput(input)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestFeatureValidateCreateInput_MissingName(t *testing.T) {
	svc, tmpDir := setupTestFeatureService(t)
	defer os.RemoveAll(tmpDir)

	writeTestProjectForFeature(t, tmpDir, "testproject", domain.ProjectStatusInitial)

	input := &domain.FeatureCreateInput{
		ProjectID: "testproject",
		Name:      "",
		Goal:      "Test goal",
	}

	err := svc.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected validation error for missing name")
	}

	me, ok := err.(*domain.MandorError)
	if !ok {
		t.Fatal("Expected MandorError")
	}

	if me.Code != domain.ExitValidationError {
		t.Errorf("Expected exit code %d, got %d", domain.ExitValidationError, me.Code)
	}
}

func TestFeatureValidateCreateInput_MissingGoal(t *testing.T) {
	svc, tmpDir := setupTestFeatureService(t)
	defer os.RemoveAll(tmpDir)

	writeTestProjectForFeature(t, tmpDir, "testproject", domain.ProjectStatusInitial)

	input := &domain.FeatureCreateInput{
		ProjectID: "testproject",
		Name:      "Test Feature",
		Goal:      "",
	}

	err := svc.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected validation error for missing goal")
	}
}

func TestFeatureValidateCreateInput_InvalidScope(t *testing.T) {
	svc, tmpDir := setupTestFeatureService(t)
	defer os.RemoveAll(tmpDir)

	writeTestProjectForFeature(t, tmpDir, "testproject", domain.ProjectStatusInitial)

	input := &domain.FeatureCreateInput{
		ProjectID: "testproject",
		Name:      "Test Feature",
		Goal:      "Test goal",
		Scope:     "invalid-scope",
	}

	err := svc.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected validation error for invalid scope")
	}
}

func TestFeatureValidateCreateInput_InvalidPriority(t *testing.T) {
	svc, tmpDir := setupTestFeatureService(t)
	defer os.RemoveAll(tmpDir)

	writeTestProjectForFeature(t, tmpDir, "testproject", domain.ProjectStatusInitial)

	input := &domain.FeatureCreateInput{
		ProjectID: "testproject",
		Name:      "Test Feature",
		Goal:      "Test goal",
		Priority:  "P6",
	}

	err := svc.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected validation error for invalid priority")
	}
}

func TestFeatureValidateCreateInput_ProjectNotFound(t *testing.T) {
	svc, tmpDir := setupTestFeatureService(t)
	defer os.RemoveAll(tmpDir)

	input := &domain.FeatureCreateInput{
		ProjectID: "nonexistent",
		Name:      "Test Feature",
		Goal:      "Test goal",
	}

	err := svc.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected validation error for non-existent project")
	}
}

func TestFeatureCreate(t *testing.T) {
	svc, tmpDir := setupTestFeatureService(t)
	defer os.RemoveAll(tmpDir)

	writeTestProjectForFeature(t, tmpDir, "testproject", domain.ProjectStatusInitial)

	input := &domain.FeatureCreateInput{
		ProjectID: "testproject",
		Name:      "New Feature",
		Goal:      "This is a test feature goal",
		Scope:     "frontend",
		Priority:  "P2",
	}

	feature, err := svc.CreateFeature(input)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if feature.ID == "" {
		t.Error("Expected feature ID to be set")
	}

	if feature.Status != domain.FeatureStatusDraft {
		t.Errorf("Expected status draft, got: %s", feature.Status)
	}

	if feature.ProjectID != "testproject" {
		t.Errorf("Expected project ID testproject, got: %s", feature.ProjectID)
	}
}

func TestFeatureCreate_BlockedStatus(t *testing.T) {
	svc, tmpDir := setupTestFeatureService(t)
	defer os.RemoveAll(tmpDir)

	writeTestProjectForFeature(t, tmpDir, "testproject", domain.ProjectStatusInitial)

	writeTestFeature(t, tmpDir, "testproject", "testproject-feature-dep123", domain.FeatureStatusDraft, nil)

	input := &domain.FeatureCreateInput{
		ProjectID: "testproject",
		Name:      "Blocked Feature",
		Goal:      "This feature depends on another",
		DependsOn: []string{"testproject-feature-dep123"},
	}

	feature, err := svc.CreateFeature(input)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if feature.Status != domain.FeatureStatusBlocked {
		t.Errorf("Expected status blocked, got: %s", feature.Status)
	}
}

func TestFeatureList(t *testing.T) {
	svc, tmpDir := setupTestFeatureService(t)
	defer os.RemoveAll(tmpDir)

	writeTestProjectForFeature(t, tmpDir, "testproject", domain.ProjectStatusInitial)

	writeTestFeature(t, tmpDir, "testproject", "testproject-feature-abc123", domain.FeatureStatusDraft, nil)
	writeTestFeature(t, tmpDir, "testproject", "testproject-feature-def456", domain.FeatureStatusActive, nil)

	input := &domain.FeatureListInput{
		ProjectID:      "testproject",
		IncludeDeleted: false,
	}

	output, err := svc.ListFeatures(input)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if output.Total != 2 {
		t.Errorf("Expected 2 features, got: %d", output.Total)
	}
}

func TestFeatureDetail(t *testing.T) {
	svc, tmpDir := setupTestFeatureService(t)
	defer os.RemoveAll(tmpDir)

	writeTestProjectForFeature(t, tmpDir, "testproject", domain.ProjectStatusInitial)
	writeTestFeature(t, tmpDir, "testproject", "testproject-feature-abc123", domain.FeatureStatusActive, nil)

	input := &domain.FeatureDetailInput{
		ProjectID:      "testproject",
		FeatureID:      "testproject-feature-abc123",
		IncludeDeleted: false,
	}

	output, err := svc.GetFeatureDetail(input)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if output.ID != "testproject-feature-abc123" {
		t.Errorf("Expected feature ID testproject-feature-abc123, got: %s", output.ID)
	}

	if output.Status != domain.FeatureStatusActive {
		t.Errorf("Expected status active, got: %s", output.Status)
	}
}

func TestFeatureDetail_NotFound(t *testing.T) {
	svc, tmpDir := setupTestFeatureService(t)
	defer os.RemoveAll(tmpDir)

	writeTestProjectForFeature(t, tmpDir, "testproject", domain.ProjectStatusInitial)

	input := &domain.FeatureDetailInput{
		ProjectID:      "testproject",
		FeatureID:      "nonexistent",
		IncludeDeleted: false,
	}

	_, err := svc.GetFeatureDetail(input)
	if err == nil {
		t.Error("Expected error for non-existent feature")
	}
}

func TestFeatureValidateUpdateInput_Valid(t *testing.T) {
	svc, tmpDir := setupTestFeatureService(t)
	defer os.RemoveAll(tmpDir)

	writeTestProjectForFeature(t, tmpDir, "testproject", domain.ProjectStatusInitial)
	writeTestFeature(t, tmpDir, "testproject", "testproject-feature-abc123", domain.FeatureStatusDraft, nil)

	name := "Updated Name"
	input := &domain.FeatureUpdateInput{
		ProjectID: "testproject",
		FeatureID: "testproject-feature-abc123",
		Name:      &name,
	}

	err := svc.ValidateUpdateInput(input)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestFeatureValidateUpdateInput_InvalidPriority(t *testing.T) {
	svc, tmpDir := setupTestFeatureService(t)
	defer os.RemoveAll(tmpDir)

	writeTestProjectForFeature(t, tmpDir, "testproject", domain.ProjectStatusInitial)
	writeTestFeature(t, tmpDir, "testproject", "testproject-feature-abc123", domain.FeatureStatusDraft, nil)

	priority := "P10"
	input := &domain.FeatureUpdateInput{
		ProjectID: "testproject",
		FeatureID: "testproject-feature-abc123",
		Priority:  &priority,
	}

	err := svc.ValidateUpdateInput(input)
	if err == nil {
		t.Error("Expected validation error for invalid priority")
	}
}

func TestFeatureUpdate_Cancel(t *testing.T) {
	svc, tmpDir := setupTestFeatureService(t)
	defer os.RemoveAll(tmpDir)

	writeTestProjectForFeature(t, tmpDir, "testproject", domain.ProjectStatusInitial)
	writeTestFeature(t, tmpDir, "testproject", "testproject-feature-abc123", domain.FeatureStatusDraft, nil)

	reason := "No longer needed"
	input := &domain.FeatureUpdateInput{
		ProjectID: "testproject",
		FeatureID: "testproject-feature-abc123",
		Cancel:    true,
		Reason:    &reason,
	}

	changes, err := svc.UpdateFeature(input)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(changes) == 0 {
		t.Error("Expected changes to be recorded")
	}
}

func TestFeatureUpdate_Reopen(t *testing.T) {
	svc, tmpDir := setupTestFeatureService(t)
	defer os.RemoveAll(tmpDir)

	writeTestProjectForFeature(t, tmpDir, "testproject", domain.ProjectStatusInitial)
	writeTestFeature(t, tmpDir, "testproject", "testproject-feature-abc123", domain.FeatureStatusCancelled, nil)

	input := &domain.FeatureUpdateInput{
		ProjectID: "testproject",
		FeatureID: "testproject-feature-abc123",
		Reopen:    true,
	}

	changes, err := svc.UpdateFeature(input)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(changes) == 0 {
		t.Error("Expected changes to be recorded")
	}
}

func TestFeatureValidateUpdateInput_CancelledFeature(t *testing.T) {
	svc, tmpDir := setupTestFeatureService(t)
	defer os.RemoveAll(tmpDir)

	writeTestProjectForFeature(t, tmpDir, "testproject", domain.ProjectStatusInitial)
	writeTestFeature(t, tmpDir, "testproject", "testproject-feature-abc123", domain.FeatureStatusCancelled, nil)

	name := "Updated Name"
	input := &domain.FeatureUpdateInput{
		ProjectID: "testproject",
		FeatureID: "testproject-feature-abc123",
		Name:      &name,
	}

	err := svc.ValidateUpdateInput(input)
	if err == nil {
		t.Error("Expected validation error for cancelled feature")
	}
}
