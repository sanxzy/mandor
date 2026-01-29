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

func setupTestTaskService(t *testing.T) (*service.TaskService, string) {
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

	svc := service.NewTaskServiceWithPaths(paths)

	return svc, tmpDir
}

func writeTestProjectForTask(t *testing.T, tmpDir, projectID string, status string) {
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

func writeTestFeatureForTask(t *testing.T, tmpDir, projectID, featureID string, status string) {
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
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		CreatedBy: "testuser",
		UpdatedBy: "testuser",
	}

	if err := writer.WriteFeature(projectID, feature); err != nil {
		t.Fatalf("Failed to write feature: %v", err)
	}
}

func writeTestTask(t *testing.T, tmpDir, projectID, taskID string, status string, dependsOn []string) {
	t.Helper()

	paths, err := fs.NewPathsFromRoot(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create paths: %v", err)
	}
	writer := fs.NewWriter(paths)

	task := &domain.Task{
		ID:                  taskID,
		FeatureID:           projectID + "-feature-abc",
		ProjectID:           projectID,
		Name:                "Test Task",
		Goal:                "Test goal for task",
		Priority:            "P3",
		Status:              status,
		DependsOn:           dependsOn,
		ImplementationSteps: []string{"step1", "step2"},
		TestCases:           []string{"test1", "test2"},
		DerivableFiles:      []string{"file1"},
		LibraryNeeds:        []string{"lib1"},
		CreatedAt:           time.Now().UTC(),
		UpdatedAt:           time.Now().UTC(),
		CreatedBy:           "testuser",
		UpdatedBy:           "testuser",
	}

	if err := writer.WriteTask(projectID, task); err != nil {
		t.Fatalf("Failed to write task: %v", err)
	}
}

func TestTaskWorkspaceInitialized(t *testing.T) {
	svc, tmpDir := setupTestTaskService(t)
	defer os.RemoveAll(tmpDir)

	if !svc.WorkspaceInitialized() {
		t.Error("Expected workspace to be initialized")
	}
}

func TestTaskValidateCreateInput_Valid(t *testing.T) {
	svc, tmpDir := setupTestTaskService(t)
	defer os.RemoveAll(tmpDir)

	writeTestProjectForTask(t, tmpDir, "testproject", domain.ProjectStatusInitial)
	writeTestFeatureForTask(t, tmpDir, "testproject", "testproject-feature-abc", domain.FeatureStatusActive)

	input := &domain.TaskCreateInput{
		FeatureID:           "testproject-feature-abc",
		Name:                "Test Task",
		Goal:                "This is a test task goal for validation",
		ImplementationSteps: []string{"step1", "step2"},
		TestCases:           []string{"test1", "test2"},
		DerivableFiles:      []string{"file1", "file2"},
		LibraryNeeds:        []string{"lib1", "lib2"},
		Priority:            "P3",
	}

	err := svc.ValidateCreateInput(input)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestTaskValidateCreateInput_MissingFeature(t *testing.T) {
	svc, tmpDir := setupTestTaskService(t)
	defer os.RemoveAll(tmpDir)

	input := &domain.TaskCreateInput{
		Name:                "Test Task",
		Goal:                "Test goal",
		ImplementationSteps: []string{"step1"},
		TestCases:           []string{"test1"},
		DerivableFiles:      []string{"file1"},
		LibraryNeeds:        []string{"lib1"},
	}

	err := svc.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected validation error for missing feature")
	}
}

func TestTaskValidateCreateInput_MissingName(t *testing.T) {
	svc, tmpDir := setupTestTaskService(t)
	defer os.RemoveAll(tmpDir)

	writeTestProjectForTask(t, tmpDir, "testproject", domain.ProjectStatusInitial)
	writeTestFeatureForTask(t, tmpDir, "testproject", "testproject-feature-abc", domain.FeatureStatusActive)

	input := &domain.TaskCreateInput{
		FeatureID:           "testproject-feature-abc",
		Name:                "",
		Goal:                "Test goal",
		ImplementationSteps: []string{"step1"},
		TestCases:           []string{"test1"},
		DerivableFiles:      []string{"file1"},
		LibraryNeeds:        []string{"lib1"},
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

func TestTaskValidateCreateInput_MissingGoal(t *testing.T) {
	svc, tmpDir := setupTestTaskService(t)
	defer os.RemoveAll(tmpDir)

	writeTestProjectForTask(t, tmpDir, "testproject", domain.ProjectStatusInitial)
	writeTestFeatureForTask(t, tmpDir, "testproject", "testproject-feature-abc", domain.FeatureStatusActive)

	input := &domain.TaskCreateInput{
		FeatureID:           "testproject-feature-abc",
		Name:                "Test Task",
		Goal:                "",
		ImplementationSteps: []string{"step1"},
		TestCases:           []string{"test1"},
		DerivableFiles:      []string{"file1"},
		LibraryNeeds:        []string{"lib1"},
	}

	err := svc.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected validation error for missing goal")
	}
}

func TestTaskValidateCreateInput_MissingImplSteps(t *testing.T) {
	svc, tmpDir := setupTestTaskService(t)
	defer os.RemoveAll(tmpDir)

	writeTestProjectForTask(t, tmpDir, "testproject", domain.ProjectStatusInitial)
	writeTestFeatureForTask(t, tmpDir, "testproject", "testproject-feature-abc", domain.FeatureStatusActive)

	input := &domain.TaskCreateInput{
		FeatureID:           "testproject-feature-abc",
		Name:                "Test Task",
		Goal:                "Test goal",
		ImplementationSteps: []string{},
		TestCases:           []string{"test1"},
		DerivableFiles:      []string{"file1"},
		LibraryNeeds:        []string{"lib1"},
	}

	err := svc.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected validation error for missing implementation steps")
	}
}

func TestTaskValidateCreateInput_InvalidPriority(t *testing.T) {
	svc, tmpDir := setupTestTaskService(t)
	defer os.RemoveAll(tmpDir)

	writeTestProjectForTask(t, tmpDir, "testproject", domain.ProjectStatusInitial)
	writeTestFeatureForTask(t, tmpDir, "testproject", "testproject-feature-abc", domain.FeatureStatusActive)

	input := &domain.TaskCreateInput{
		FeatureID:           "testproject-feature-abc",
		Name:                "Test Task",
		Goal:                "Test goal",
		ImplementationSteps: []string{"step1"},
		TestCases:           []string{"test1"},
		DerivableFiles:      []string{"file1"},
		LibraryNeeds:        []string{"lib1"},
		Priority:            "P6",
	}

	err := svc.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected validation error for invalid priority")
	}
}

func TestTaskValidateCreateInput_ProjectNotFound(t *testing.T) {
	svc, tmpDir := setupTestTaskService(t)
	defer os.RemoveAll(tmpDir)

	input := &domain.TaskCreateInput{
		FeatureID:           "testproject-feature-abc",
		Name:                "Test Task",
		Goal:                "Test goal",
		ImplementationSteps: []string{"step1"},
		TestCases:           []string{"test1"},
		DerivableFiles:      []string{"file1"},
		LibraryNeeds:        []string{"lib1"},
	}

	err := svc.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected validation error for non-existent project")
	}
}

func TestTaskCreate(t *testing.T) {
	svc, tmpDir := setupTestTaskService(t)
	defer os.RemoveAll(tmpDir)

	writeTestProjectForTask(t, tmpDir, "testproject", domain.ProjectStatusInitial)
	writeTestFeatureForTask(t, tmpDir, "testproject", "testproject-feature-abc", domain.FeatureStatusActive)

	input := &domain.TaskCreateInput{
		FeatureID:           "testproject-feature-abc",
		Name:                "New Task",
		Goal:                "This is a test task goal",
		ImplementationSteps: []string{"step1", "step2"},
		TestCases:           []string{"test1", "test2"},
		DerivableFiles:      []string{"file1", "file2"},
		LibraryNeeds:        []string{"lib1", "lib2"},
		Priority:            "P2",
	}

	task, err := svc.CreateTask(input)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if task.ID == "" {
		t.Error("Expected task ID to be set")
	}

	if task.Status != domain.TaskStatusReady {
		t.Errorf("Expected status ready, got: %s", task.Status)
	}

	if task.ProjectID != "testproject" {
		t.Errorf("Expected project ID testproject, got: %s", task.ProjectID)
	}
}

func TestTaskCreate_PendingStatus(t *testing.T) {
	svc, tmpDir := setupTestTaskService(t)
	defer os.RemoveAll(tmpDir)

	writeTestProjectForTask(t, tmpDir, "testproject", domain.ProjectStatusInitial)
	writeTestFeatureForTask(t, tmpDir, "testproject", "testproject-feature-abc", domain.FeatureStatusActive)

	writeTestTask(t, tmpDir, "testproject", "testproject-feature-abc-task-dep123", domain.TaskStatusReady, nil)

	input := &domain.TaskCreateInput{
		FeatureID:           "testproject-feature-abc",
		Name:                "Pending Task",
		Goal:                "This task depends on another",
		ImplementationSteps: []string{"step1"},
		TestCases:           []string{"test1"},
		DerivableFiles:      []string{"file1"},
		LibraryNeeds:        []string{"lib1"},
		DependsOn:           []string{"testproject-feature-abc-task-dep123"},
	}

	task, err := svc.CreateTask(input)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if task.Status != domain.TaskStatusBlocked {
		t.Errorf("Expected status blocked, got: %s", task.Status)
	}
}

func TestTaskList(t *testing.T) {
	svc, tmpDir := setupTestTaskService(t)
	defer os.RemoveAll(tmpDir)

	writeTestProjectForTask(t, tmpDir, "testproject", domain.ProjectStatusInitial)
	writeTestFeatureForTask(t, tmpDir, "testproject", "testproject-feature-abc", domain.FeatureStatusActive)

	writeTestTask(t, tmpDir, "testproject", "testproject-feature-abc-task-abc123", domain.TaskStatusReady, nil)
	writeTestTask(t, tmpDir, "testproject", "testproject-feature-abc-task-def456", domain.TaskStatusPending, nil)

	input := &domain.TaskListInput{
		ProjectID:      "testproject",
		IncludeDeleted: false,
	}

	output, err := svc.ListTasks(input)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if output.Total != 2 {
		t.Errorf("Expected 2 tasks, got: %d", output.Total)
	}
}

func TestTaskDetail(t *testing.T) {
	svc, tmpDir := setupTestTaskService(t)
	defer os.RemoveAll(tmpDir)

	writeTestProjectForTask(t, tmpDir, "testproject", domain.ProjectStatusInitial)
	writeTestFeatureForTask(t, tmpDir, "testproject", "testproject-feature-abc", domain.FeatureStatusActive)
	writeTestTask(t, tmpDir, "testproject", "testproject-feature-abc-task-abc123", domain.TaskStatusReady, nil)

	input := &domain.TaskDetailInput{
		TaskID:         "testproject-feature-abc-task-abc123",
		IncludeDeleted: false,
	}

	output, err := svc.GetTaskDetail(input)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if output.ID != "testproject-feature-abc-task-abc123" {
		t.Errorf("Expected task ID testproject-feature-abc-task-abc123, got: %s", output.ID)
	}

	if output.Status != domain.TaskStatusReady {
		t.Errorf("Expected status ready, got: %s", output.Status)
	}
}

func TestTaskDetail_NotFound(t *testing.T) {
	svc, tmpDir := setupTestTaskService(t)
	defer os.RemoveAll(tmpDir)

	writeTestProjectForTask(t, tmpDir, "testproject", domain.ProjectStatusInitial)

	input := &domain.TaskDetailInput{
		TaskID:         "testproject-feature-abc-task-nonexistent",
		IncludeDeleted: false,
	}

	_, err := svc.GetTaskDetail(input)
	if err == nil {
		t.Error("Expected error for non-existent task")
	}
}

func TestTaskValidateUpdateInput_Valid(t *testing.T) {
	svc, tmpDir := setupTestTaskService(t)
	defer os.RemoveAll(tmpDir)

	writeTestProjectForTask(t, tmpDir, "testproject", domain.ProjectStatusInitial)
	writeTestFeatureForTask(t, tmpDir, "testproject", "testproject-feature-abc", domain.FeatureStatusActive)
	writeTestTask(t, tmpDir, "testproject", "testproject-feature-abc-task-abc123", domain.TaskStatusReady, nil)

	name := "Updated Name"
	input := &domain.TaskUpdateInput{
		TaskID: "testproject-feature-abc-task-abc123",
		Name:   &name,
	}

	err := svc.ValidateUpdateInput(input)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestTaskValidateUpdateInput_InvalidPriority(t *testing.T) {
	svc, tmpDir := setupTestTaskService(t)
	defer os.RemoveAll(tmpDir)

	writeTestProjectForTask(t, tmpDir, "testproject", domain.ProjectStatusInitial)
	writeTestFeatureForTask(t, tmpDir, "testproject", "testproject-feature-abc", domain.FeatureStatusActive)
	writeTestTask(t, tmpDir, "testproject", "testproject-feature-abc-task-abc123", domain.TaskStatusReady, nil)

	priority := "P10"
	input := &domain.TaskUpdateInput{
		TaskID:   "testproject-feature-abc-task-abc123",
		Priority: &priority,
	}

	err := svc.ValidateUpdateInput(input)
	if err == nil {
		t.Error("Expected validation error for invalid priority")
	}
}

func TestTaskUpdate_Cancel(t *testing.T) {
	svc, tmpDir := setupTestTaskService(t)
	defer os.RemoveAll(tmpDir)

	writeTestProjectForTask(t, tmpDir, "testproject", domain.ProjectStatusInitial)
	writeTestFeatureForTask(t, tmpDir, "testproject", "testproject-feature-abc", domain.FeatureStatusActive)
	writeTestTask(t, tmpDir, "testproject", "testproject-feature-abc-task-abc123", domain.TaskStatusReady, nil)

	reason := "No longer needed"
	input := &domain.TaskUpdateInput{
		TaskID: "testproject-feature-abc-task-abc123",
		Cancel: true,
		Reason: &reason,
	}

	changes, err := svc.UpdateTask(input)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(changes) == 0 {
		t.Error("Expected changes to be recorded")
	}
}

func TestTaskUpdate_Reopen(t *testing.T) {
	svc, tmpDir := setupTestTaskService(t)
	defer os.RemoveAll(tmpDir)

	writeTestProjectForTask(t, tmpDir, "testproject", domain.ProjectStatusInitial)
	writeTestFeatureForTask(t, tmpDir, "testproject", "testproject-feature-abc", domain.FeatureStatusActive)
	writeTestTask(t, tmpDir, "testproject", "testproject-feature-abc-task-abc123", domain.TaskStatusCancelled, nil)

	input := &domain.TaskUpdateInput{
		TaskID: "testproject-feature-abc-task-abc123",
		Reopen: true,
	}

	changes, err := svc.UpdateTask(input)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(changes) == 0 {
		t.Error("Expected changes to be recorded")
	}
}

func TestTaskValidateUpdateInput_CancelledTask(t *testing.T) {
	svc, tmpDir := setupTestTaskService(t)
	defer os.RemoveAll(tmpDir)

	writeTestProjectForTask(t, tmpDir, "testproject", domain.ProjectStatusInitial)
	writeTestFeatureForTask(t, tmpDir, "testproject", "testproject-feature-abc", domain.FeatureStatusActive)
	writeTestTask(t, tmpDir, "testproject", "testproject-feature-abc-task-abc123", domain.TaskStatusCancelled, nil)

	name := "Updated Name"
	input := &domain.TaskUpdateInput{
		TaskID: "testproject-feature-abc-task-abc123",
		Name:   &name,
	}

	err := svc.ValidateUpdateInput(input)
	if err == nil {
		t.Error("Expected validation error for cancelled task")
	}
}

func TestTaskParseTaskID(t *testing.T) {
	svc, tmpDir := setupTestTaskService(t)
	defer os.RemoveAll(tmpDir)

	projectID, featureID, err := svc.ParseTaskID("testproject-feature-abc-task-xyz123")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if projectID != "testproject" {
		t.Errorf("Expected project ID testproject, got: %s", projectID)
	}

	if featureID != "testproject-feature-abc" {
		t.Errorf("Expected feature ID testproject-feature-abc, got: %s", featureID)
	}
}

func TestTaskParseTaskID_Invalid(t *testing.T) {
	svc, tmpDir := setupTestTaskService(t)
	defer os.RemoveAll(tmpDir)

	_, _, err := svc.ParseTaskID("invalid-id")
	if err == nil {
		t.Error("Expected error for invalid task ID")
	}
}
