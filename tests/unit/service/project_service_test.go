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

func setupTestProjectService(t *testing.T) (*service.ProjectService, string) {
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

	svc := service.NewProjectServiceWithPaths(paths)

	return svc, tmpDir
}

func writeTestProject(t *testing.T, tmpDir, projectID string, status string) {
	t.Helper()

	paths, err := fs.NewPathsFromRoot(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create paths: %v", err)
	}
	writer := fs.NewWriter(paths)

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
}

func TestWorkspaceInitialized(t *testing.T) {
	svc, tmpDir := setupTestProjectService(t)
	defer os.RemoveAll(tmpDir)

	if !svc.WorkspaceInitialized() {
		t.Error("Expected workspace to be initialized")
	}
}

func TestValidateCreateInput_Valid(t *testing.T) {
	svc, tmpDir := setupTestProjectService(t)
	defer os.RemoveAll(tmpDir)

	input := &domain.ProjectCreateInput{
		ID:         "testproject",
		Name:       "Test Project",
		Goal:       string(make([]byte, 500)),
		TaskDep:    "same_project_only",
		FeatureDep: "cross_project_allowed",
		IssueDep:   "same_project_only",
		Strict:     false,
	}

	err := svc.ValidateCreateInput(input)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestValidateCreateInput_InvalidID(t *testing.T) {
	svc, tmpDir := setupTestProjectService(t)
	defer os.RemoveAll(tmpDir)

	input := &domain.ProjectCreateInput{
		ID:   "123invalid",
		Name: "Test Project",
		Goal: string(make([]byte, 500)),
	}

	err := svc.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected validation error for invalid ID")
	}

	me, ok := err.(*domain.MandorError)
	if !ok {
		t.Fatal("Expected MandorError")
	}

	if me.Code != domain.ExitValidationError {
		t.Errorf("Expected exit code %d, got %d", domain.ExitValidationError, me.Code)
	}
}

func TestValidateCreateInput_ProjectExists(t *testing.T) {
	svc, tmpDir := setupTestProjectService(t)
	defer os.RemoveAll(tmpDir)

	writeTestProject(t, tmpDir, "existing", domain.ProjectStatusInitial)

	input := &domain.ProjectCreateInput{
		ID:   "existing",
		Name: "Existing Project",
		Goal: string(make([]byte, 500)),
	}

	err := svc.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected validation error for existing project")
	}

	me, ok := err.(*domain.MandorError)
	if !ok {
		t.Fatal("Expected MandorError")
	}

	if me.Code != domain.ExitValidationError {
		t.Errorf("Expected exit code %d, got %d", domain.ExitValidationError, me.Code)
	}
}

func TestValidateUpdateInput_Valid(t *testing.T) {
	svc, tmpDir := setupTestProjectService(t)
	defer os.RemoveAll(tmpDir)

	writeTestProject(t, tmpDir, "testproj", domain.ProjectStatusInitial)

	name := "Updated Name"
	input := &domain.ProjectUpdateInput{
		ID:   "testproj",
		Name: &name,
	}

	err := svc.ValidateUpdateInput(input)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestValidateUpdateInput_DeletedProject(t *testing.T) {
	svc, tmpDir := setupTestProjectService(t)
	defer os.RemoveAll(tmpDir)

	writeTestProject(t, tmpDir, "deletedproj", domain.ProjectStatusDeleted)

	name := "Updated Name"
	input := &domain.ProjectUpdateInput{
		ID:   "deletedproj",
		Name: &name,
	}

	err := svc.ValidateUpdateInput(input)
	if err == nil {
		t.Error("Expected validation error for deleted project")
	}
}

func TestValidateDeleteInput_SoftDeleteOnDeleted(t *testing.T) {
	svc, tmpDir := setupTestProjectService(t)
	defer os.RemoveAll(tmpDir)

	writeTestProject(t, tmpDir, "already-deleted", domain.ProjectStatusDeleted)

	input := &domain.ProjectDeleteInput{
		ID:   "already-deleted",
		Hard: false,
	}

	err := svc.ValidateDeleteInput(input)
	if err == nil {
		t.Error("Expected validation error for soft delete on already deleted project")
	}
}

func TestValidateDeleteInput_HardDeleteOnDeleted(t *testing.T) {
	svc, tmpDir := setupTestProjectService(t)
	defer os.RemoveAll(tmpDir)

	writeTestProject(t, tmpDir, "already-deleted", domain.ProjectStatusDeleted)

	input := &domain.ProjectDeleteInput{
		ID:   "already-deleted",
		Hard: true,
	}

	err := svc.ValidateDeleteInput(input)
	if err != nil {
		t.Errorf("Expected no error for hard delete on already deleted project, got: %v", err)
	}
}

func TestValidateReopenInput_Valid(t *testing.T) {
	svc, tmpDir := setupTestProjectService(t)
	defer os.RemoveAll(tmpDir)

	writeTestProject(t, tmpDir, "toreopen", domain.ProjectStatusDeleted)

	input := &domain.ProjectReopenInput{
		ID:  "toreopen",
		Yes: true,
	}

	err := svc.ValidateReopenInput(input)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestValidateReopenInput_NotDeleted(t *testing.T) {
	svc, tmpDir := setupTestProjectService(t)
	defer os.RemoveAll(tmpDir)

	writeTestProject(t, tmpDir, "notdeleted", domain.ProjectStatusInitial)

	input := &domain.ProjectReopenInput{
		ID:  "notdeleted",
		Yes: true,
	}

	err := svc.ValidateReopenInput(input)
	if err == nil {
		t.Error("Expected validation error for non-deleted project")
	}
}
