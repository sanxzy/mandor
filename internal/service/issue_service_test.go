package service

import (
	"os"
	"testing"
	"time"

	"mandor/internal/domain"
	"mandor/internal/fs"
)

func setupTestIssueService(t *testing.T) (*IssueService, string) {
	tmpDir := t.TempDir()

	paths := &fs.Paths{WorkspaceRoot: tmpDir}
	writer := fs.NewWriter(paths)
	reader := fs.NewReader(paths)

	os.MkdirAll(tmpDir+"/.mandor/backlogs", 0755)

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

	service := &IssueService{
		reader: reader,
		writer: writer,
		paths:  paths,
	}

	return service, tmpDir
}

func setupTestIssueServiceWithBacklog(t *testing.T, backlogID, issueDepRule string) (*IssueService, *BacklogService) {
	tmpDir := t.TempDir()

	paths := &fs.Paths{WorkspaceRoot: tmpDir}
	writer := fs.NewWriter(paths)
	reader := fs.NewReader(paths)

	os.MkdirAll(tmpDir+"/.mandor/backlogs", 0755)

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
		ID:       backlogID,
		Name:     "Test Backlog",
		Goal:     "This is a test goal for the backlog that has enough characters to pass validation requirements for the backlog entity",
		IssueDep: issueDepRule,
	}
	if err := backlogService.CreateBacklog(backlogInput); err != nil {
		t.Fatalf("Failed to create backlog: %v", err)
	}

	issueService := &IssueService{
		reader: reader,
		writer: writer,
		paths:  paths,
	}

	return issueService, backlogService
}

func TestNewIssueService(t *testing.T) {
	service, err := NewIssueService()
	if err != nil {
		t.Fatalf("NewIssueService failed: %v", err)
	}

	if service == nil {
		t.Error("NewIssueService returned nil")
	}
	if service.paths == nil {
		t.Error("Service paths should not be nil")
	}
	if service.reader == nil {
		t.Error("Service reader is nil")
	}
	if service.writer == nil {
		t.Error("Service writer is nil")
	}
}

func TestNewIssueServiceWithPaths(t *testing.T) {
	tmpDir := t.TempDir()
	paths := &fs.Paths{WorkspaceRoot: tmpDir}

	service := NewIssueServiceWithPaths(paths)

	if service == nil {
		t.Error("NewIssueServiceWithPaths returned nil")
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

func TestIssueWorkspaceInitialized(t *testing.T) {
	service, _ := setupTestIssueService(t)

	if !service.WorkspaceInitialized() {
		t.Error("Workspace should be initialized")
	}
}

func TestIssueValidateCreateInput_BacklogNotFound(t *testing.T) {
	service, _ := setupTestIssueService(t)

	input := &domain.IssueCreateInput{
		BacklogID:           "nonexistent",
		Name:                "Test Issue",
		Goal:                "This is a test goal for the issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "bug",
		Priority:            "P2",
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected validation error for nonexistent backlog, got nil")
	}
}

func TestValidateCreateInput_NameRequired(t *testing.T) {
	service, _ := setupTestIssueServiceWithBacklog(t, "test-backlog", "same_backlog_only")

	input := &domain.IssueCreateInput{
		BacklogID:           "test-backlog",
		Name:                "",
		Goal:                "This is a test goal for the issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "bug",
		Priority:            "P2",
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected validation error for empty name, got nil")
	}
}

func TestValidateCreateInput_GoalRequired(t *testing.T) {
	service, _ := setupTestIssueServiceWithBacklog(t, "test-backlog", "same_backlog_only")

	input := &domain.IssueCreateInput{
		BacklogID:           "test-backlog",
		Name:                "Test Issue",
		Goal:                "",
		IssueType:           "bug",
		Priority:            "P2",
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected validation error for empty goal, got nil")
	}
}

func TestValidateCreateInput_GoalTooShort(t *testing.T) {
	service, _ := setupTestIssueServiceWithBacklog(t, "test-backlog", "same_backlog_only")

	input := &domain.IssueCreateInput{
		BacklogID:           "test-backlog",
		Name:                "Test Issue",
		Goal:                "Short goal",
		IssueType:           "bug",
		Priority:            "P2",
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected validation error for short goal, got nil")
	}
}

func TestValidateCreateInput_IssueTypeRequired(t *testing.T) {
	service, _ := setupTestIssueServiceWithBacklog(t, "test-backlog", "same_backlog_only")

	input := &domain.IssueCreateInput{
		BacklogID:           "test-backlog",
		Name:                "Test Issue",
		Goal:                "This is a test goal for the issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "",
		Priority:            "P2",
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected validation error for empty issue type, got nil")
	}
}

func TestValidateCreateInput_IssueTypeInvalid(t *testing.T) {
	service, _ := setupTestIssueServiceWithBacklog(t, "test-backlog", "same_backlog_only")

	input := &domain.IssueCreateInput{
		BacklogID:           "test-backlog",
		Name:                "Test Issue",
		Goal:                "This is a test goal for the issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "invalid_type",
		Priority:            "P2",
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected validation error for invalid issue type, got nil")
	}
}

func TestValidateCreateInput_AffectedFilesRequired(t *testing.T) {
	service, _ := setupTestIssueServiceWithBacklog(t, "test-backlog", "same_backlog_only")

	input := &domain.IssueCreateInput{
		BacklogID:           "test-backlog",
		Name:                "Test Issue",
		Goal:                "This is a test goal for the issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "bug",
		Priority:            "P2",
		AffectedFiles:       []string{},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected validation error for empty affected files, got nil")
	}
}

func TestValidateCreateInput_AffectedTestsRequired(t *testing.T) {
	service, _ := setupTestIssueServiceWithBacklog(t, "test-backlog", "same_backlog_only")

	input := &domain.IssueCreateInput{
		BacklogID:           "test-backlog",
		Name:                "Test Issue",
		Goal:                "This is a test goal for the issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "bug",
		Priority:            "P2",
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{},
		ImplementationSteps: []string{"Step 1"},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected validation error for empty affected tests, got nil")
	}
}

func TestValidateCreateInput_ImplementationStepsRequired(t *testing.T) {
	service, _ := setupTestIssueServiceWithBacklog(t, "test-backlog", "same_backlog_only")

	input := &domain.IssueCreateInput{
		BacklogID:           "test-backlog",
		Name:                "Test Issue",
		Goal:                "This is a test goal for the issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "bug",
		Priority:            "P2",
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected validation error for empty implementation steps, got nil")
	}
}

func TestValidateCreateInput_PriorityValid(t *testing.T) {
	service, _ := setupTestIssueServiceWithBacklog(t, "test-backlog", "same_backlog_only")

	input := &domain.IssueCreateInput{
		BacklogID:           "test-backlog",
		Name:                "Test Issue",
		Goal:                "This is a test goal for the issue that has enough characters to pass validation requirements and exceed minimum length requirement of 200 characters for issue goals in the Mandor system. This goal provides clear context and intent for the issue resolution process.",
		IssueType:           "bug",
		Priority:            "P0",
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
	}

	err := service.ValidateCreateInput(input)
	if err != nil {
		t.Errorf("Unexpected error for valid priority: %v", err)
	}
}

func TestValidateCreateInput_PriorityInvalid(t *testing.T) {
	service, _ := setupTestIssueServiceWithBacklog(t, "test-backlog", "same_backlog_only")

	input := &domain.IssueCreateInput{
		BacklogID:           "test-backlog",
		Name:                "Test Issue",
		Goal:                "This is a test goal for the issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "bug",
		Priority:            "P6",
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected validation error for invalid priority, got nil")
	}
}

func TestValidateCreateInput_DependencyNotFound(t *testing.T) {
	service, _ := setupTestIssueServiceWithBacklog(t, "test-backlog", "same_backlog_only")

	input := &domain.IssueCreateInput{
		BacklogID:           "test-backlog",
		Name:                "Test Issue",
		Goal:                "This is a test goal for the issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "bug",
		Priority:            "P2",
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
		DependsOn:           []string{"test-backlog-issue-nonexistent"},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected validation error for nonexistent dependency, got nil")
	}
}

func TestValidateCreateInput_DependencyCancelled(t *testing.T) {
	service, _ := setupTestIssueServiceWithBacklog(t, "test-backlog", "same_backlog_only")

	existingIssue := &domain.Issue{
		ID:                  "test-backlog-issue-existing",
		BacklogID:           "test-backlog",
		Name:                "Existing Issue",
		Goal:                "This is a test goal for the existing issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "bug",
		Priority:            "P2",
		Status:              domain.IssueStatusCancelled,
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
		CreatedAt:           time.Now(),
		LastUpdatedAt:       time.Now(),
		CreatedBy:           "test-user",
		LastUpdatedBy:       "test-user",
	}
	if err := service.writer.WriteIssue("test-backlog", existingIssue); err != nil {
		t.Fatalf("Failed to write existing issue: %v", err)
	}

	input := &domain.IssueCreateInput{
		BacklogID:           "test-backlog",
		Name:                "Test Issue",
		Goal:                "This is a test goal for the issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "bug",
		Priority:            "P2",
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
		DependsOn:           []string{"test-backlog-issue-existing"},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected validation error for cancelled dependency, got nil")
	}
}

func TestValidateCreateInput_CrossBacklogDependencyAllowed(t *testing.T) {
	service, backlogSvc := setupTestIssueServiceWithBacklog(t, "backlog1", "cross_backlog_allowed")

	otherBacklogInput := &domain.BacklogCreateInput{
		ID:       "backlog2",
		Name:     "Other Backlog",
		Goal:     "This is a test goal for the other backlog that has enough characters to pass validation requirements",
		IssueDep: "cross_backlog_allowed",
	}
	if err := backlogSvc.CreateBacklog(otherBacklogInput); err != nil {
		t.Fatalf("Failed to create other backlog: %v", err)
	}

	existingIssue := &domain.Issue{
		ID:                  "backlog2-issue-existing",
		BacklogID:           "backlog2",
		Name:                "Existing Issue",
		Goal:                "This is a test goal for the existing issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "bug",
		Priority:            "P2",
		Status:              domain.IssueStatusResolved,
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
		CreatedAt:           time.Now(),
		LastUpdatedAt:       time.Now(),
		CreatedBy:           "test-user",
		LastUpdatedBy:       "test-user",
	}
	if err := service.writer.WriteIssue("backlog2", existingIssue); err != nil {
		t.Fatalf("Failed to write existing issue: %v", err)
	}

	input := &domain.IssueCreateInput{
		BacklogID:           "backlog1",
		Name:                "Test Issue",
		Goal:                "This is a test goal for the issue that has enough characters to pass validation requirements and exceed minimum length requirement of 200 characters for issue goals in the Mandor system. This goal provides clear context and intent for the issue resolution process.",
		IssueType:           "bug",
		Priority:            "P2",
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
		DependsOn:           []string{"backlog2-issue-existing"},
	}

	err := service.ValidateCreateInput(input)
	if err != nil {
		t.Errorf("Unexpected error for allowed cross-backlog dependency: %v", err)
	}
}

func TestValidateCreateInput_CrossBacklogDependencyNotAllowed(t *testing.T) {
	service, backlogSvc := setupTestIssueServiceWithBacklog(t, "backlog1", "same_backlog_only")

	otherBacklogInput := &domain.BacklogCreateInput{
		ID:       "backlog2",
		Name:     "Other Backlog",
		Goal:     "This is a test goal for the other backlog that has enough characters to pass validation requirements and exceed minimum length",
		IssueDep: "same_backlog_only",
	}
	if err := backlogSvc.CreateBacklog(otherBacklogInput); err != nil {
		t.Fatalf("Failed to create other backlog: %v", err)
	}

	existingIssue := &domain.Issue{
		ID:                  "backlog2-issue-existing",
		BacklogID:           "backlog2",
		Name:                "Existing Issue",
		Goal:                "This is a test goal for the existing issue that has enough characters to pass validation requirements and exceed minimum length requirement",
		IssueType:           "bug",
		Priority:            "P2",
		Status:              domain.IssueStatusResolved,
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
		CreatedAt:           time.Now(),
		LastUpdatedAt:       time.Now(),
		CreatedBy:           "test-user",
		LastUpdatedBy:       "test-user",
	}
	if err := service.writer.WriteIssue("backlog2", existingIssue); err != nil {
		t.Fatalf("Failed to write existing issue: %v", err)
	}

	input := &domain.IssueCreateInput{
		BacklogID:           "backlog1",
		Name:                "Test Issue",
		Goal:                "This is a test goal for the issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "bug",
		Priority:            "P2",
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
		DependsOn:           []string{"backlog2-issue-existing"},
	}

	err := service.ValidateCreateInput(input)
	if err == nil {
		t.Error("Expected validation error for disallowed cross-backlog dependency, got nil")
	}
}

func TestCreateIssue_Success(t *testing.T) {
	service, _ := setupTestIssueServiceWithBacklog(t, "test-backlog", "same_backlog_only")

	input := &domain.IssueCreateInput{
		BacklogID:           "test-backlog",
		Name:                "Test Issue",
		Goal:                "This is a test goal for the issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "bug",
		Priority:            "P2",
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
	}

	issue, err := service.CreateIssue(input)
	if err != nil {
		t.Fatalf("CreateIssue failed: %v", err)
	}

	if issue == nil {
		t.Fatal("CreateIssue returned nil")
	}

	if issue.ID == "" {
		t.Error("Issue ID should not be empty")
	}
	if issue.BacklogID != "test-backlog" {
		t.Errorf("Issue BacklogID = %q, want %q", issue.BacklogID, "test-backlog")
	}
	if issue.Name != "Test Issue" {
		t.Errorf("Issue Name = %q, want %q", issue.Name, "Test Issue")
	}
	if issue.Status != domain.IssueStatusReady {
		t.Errorf("Issue Status = %q, want %q", issue.Status, domain.IssueStatusReady)
	}
}

func TestCreateIssue_BlockedDependencies(t *testing.T) {
	service, _ := setupTestIssueServiceWithBacklog(t, "test-backlog", "same_backlog_only")

	depIssue := &domain.Issue{
		ID:                  "test-backlog-issue-dep",
		BacklogID:           "test-backlog",
		Name:                "Dependency Issue",
		Goal:                "This is a test goal for the dependency issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "bug",
		Priority:            "P2",
		Status:              domain.IssueStatusOpen,
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
		CreatedAt:           time.Now(),
		LastUpdatedAt:       time.Now(),
		CreatedBy:           "test-user",
		LastUpdatedBy:       "test-user",
	}
	if err := service.writer.WriteIssue("test-backlog", depIssue); err != nil {
		t.Fatalf("Failed to write dependency issue: %v", err)
	}

	input := &domain.IssueCreateInput{
		BacklogID:           "test-backlog",
		Name:                "Test Issue",
		Goal:                "This is a test goal for the issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "bug",
		Priority:            "P2",
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
		DependsOn:           []string{"test-backlog-issue-dep"},
	}

	issue, err := service.CreateIssue(input)
	if err != nil {
		t.Fatalf("CreateIssue failed: %v", err)
	}

	if issue.Status != domain.IssueStatusBlocked {
		t.Errorf("Issue Status = %q, want %q", issue.Status, domain.IssueStatusBlocked)
	}
}

func TestCreateIssue_ResolvedDependency(t *testing.T) {
	service, _ := setupTestIssueServiceWithBacklog(t, "test-backlog", "same_backlog_only")

	depIssue := &domain.Issue{
		ID:                  "test-backlog-issue-dep",
		BacklogID:           "test-backlog",
		Name:                "Dependency Issue",
		Goal:                "This is a test goal for the dependency issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "bug",
		Priority:            "P2",
		Status:              domain.IssueStatusResolved,
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
		CreatedAt:           time.Now(),
		LastUpdatedAt:       time.Now(),
		CreatedBy:           "test-user",
		LastUpdatedBy:       "test-user",
	}
	if err := service.writer.WriteIssue("test-backlog", depIssue); err != nil {
		t.Fatalf("Failed to write dependency issue: %v", err)
	}

	input := &domain.IssueCreateInput{
		BacklogID:           "test-backlog",
		Name:                "Test Issue",
		Goal:                "This is a test goal for the issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "bug",
		Priority:            "P2",
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
		DependsOn:           []string{"test-backlog-issue-dep"},
	}

	issue, err := service.CreateIssue(input)
	if err != nil {
		t.Fatalf("CreateIssue failed: %v", err)
	}

	if issue.Status != domain.IssueStatusReady {
		t.Errorf("Issue Status = %q, want %q", issue.Status, domain.IssueStatusReady)
	}
}

func TestListIssues_Empty(t *testing.T) {
	service, _ := setupTestIssueServiceWithBacklog(t, "test-backlog", "same_backlog_only")

	input := &domain.IssueListInput{
		BacklogID: "test-backlog",
	}

	result, err := service.ListIssues(input)
	if err != nil {
		t.Fatalf("ListIssues failed: %v", err)
	}

	if len(result.Issues) != 0 {
		t.Errorf("Issues count = %d, want 0", len(result.Issues))
	}
}

func TestListIssues_WithIssues(t *testing.T) {
	service, _ := setupTestIssueServiceWithBacklog(t, "test-backlog", "same_backlog_only")

	issue := &domain.Issue{
		ID:                  "test-backlog-issue-1",
		BacklogID:           "test-backlog",
		Name:                "Test Issue 1",
		Goal:                "This is a test goal for the issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "bug",
		Priority:            "P2",
		Status:              domain.IssueStatusOpen,
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
		CreatedAt:           time.Now(),
		LastUpdatedAt:       time.Now(),
		CreatedBy:           "test-user",
		LastUpdatedBy:       "test-user",
	}
	if err := service.writer.WriteIssue("test-backlog", issue); err != nil {
		t.Fatalf("Failed to write issue: %v", err)
	}

	input := &domain.IssueListInput{
		BacklogID: "test-backlog",
	}

	result, err := service.ListIssues(input)
	if err != nil {
		t.Fatalf("ListIssues failed: %v", err)
	}

	if len(result.Issues) != 1 {
		t.Errorf("Issues count = %d, want 1", len(result.Issues))
	}
}

func TestListIssues_FilterByType(t *testing.T) {
	service, _ := setupTestIssueServiceWithBacklog(t, "test-backlog", "same_backlog_only")

	issue1 := &domain.Issue{
		ID:                  "test-backlog-issue-1",
		BacklogID:           "test-backlog",
		Name:                "Bug Issue",
		Goal:                "This is a test goal for the issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "bug",
		Priority:            "P2",
		Status:              domain.IssueStatusOpen,
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
		CreatedAt:           time.Now(),
		LastUpdatedAt:       time.Now(),
		CreatedBy:           "test-user",
		LastUpdatedBy:       "test-user",
	}
	issue2 := &domain.Issue{
		ID:                  "test-backlog-issue-2",
		BacklogID:           "test-backlog",
		Name:                "Security Issue",
		Goal:                "This is a test goal for the issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "security",
		Priority:            "P1",
		Status:              domain.IssueStatusOpen,
		AffectedFiles:       []string{"file2.go"},
		AffectedTests:       []string{"TestFile2"},
		ImplementationSteps: []string{"Step 1"},
		CreatedAt:           time.Now(),
		LastUpdatedAt:       time.Now(),
		CreatedBy:           "test-user",
		LastUpdatedBy:       "test-user",
	}
	if err := service.writer.WriteIssue("test-backlog", issue1); err != nil {
		t.Fatalf("Failed to write issue 1: %v", err)
	}
	if err := service.writer.WriteIssue("test-backlog", issue2); err != nil {
		t.Fatalf("Failed to write issue 2: %v", err)
	}

	input := &domain.IssueListInput{
		BacklogID: "test-backlog",
		IssueType: "security",
	}

	result, err := service.ListIssues(input)
	if err != nil {
		t.Fatalf("ListIssues failed: %v", err)
	}

	if len(result.Issues) != 1 {
		t.Errorf("Issues count = %d, want 1", len(result.Issues))
	}
	if result.Issues[0].IssueType != "security" {
		t.Errorf("Issue Type = %q, want %q", result.Issues[0].IssueType, "security")
	}
}

func TestGetIssueDetail_Success(t *testing.T) {
	service, _ := setupTestIssueServiceWithBacklog(t, "test-backlog", "same_backlog_only")

	issue := &domain.Issue{
		ID:                  "test-backlog-issue-1",
		BacklogID:           "test-backlog",
		Name:                "Test Issue",
		Goal:                "This is a test goal for the issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "bug",
		Priority:            "P2",
		Status:              domain.IssueStatusOpen,
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
		CreatedAt:           time.Now(),
		LastUpdatedAt:       time.Now(),
		CreatedBy:           "test-user",
		LastUpdatedBy:       "test-user",
	}
	if err := service.writer.WriteIssue("test-backlog", issue); err != nil {
		t.Fatalf("Failed to write issue: %v", err)
	}

	input := &domain.IssueDetailInput{
		BacklogID: "test-backlog",
		IssueID:   "test-backlog-issue-1",
	}

	result, err := service.GetIssueDetail(input)
	if err != nil {
		t.Fatalf("GetIssueDetail failed: %v", err)
	}

	if result.ID != "test-backlog-issue-1" {
		t.Errorf("Issue ID = %q, want %q", result.ID, "test-backlog-issue-1")
	}
	if result.Name != "Test Issue" {
		t.Errorf("Issue Name = %q, want %q", result.Name, "Test Issue")
	}
}

func TestGetIssueDetail_NotFound(t *testing.T) {
	service, _ := setupTestIssueServiceWithBacklog(t, "test-backlog", "same_backlog_only")

	input := &domain.IssueDetailInput{
		BacklogID: "test-backlog",
		IssueID:   "test-backlog-issue-nonexistent",
	}

	_, err := service.GetIssueDetail(input)
	if err == nil {
		t.Error("Expected error for nonexistent issue, got nil")
	}
}

func TestGetIssueDetail_CancelledNotIncluded(t *testing.T) {
	service, _ := setupTestIssueServiceWithBacklog(t, "test-backlog", "same_backlog_only")

	issue := &domain.Issue{
		ID:                  "test-backlog-issue-1",
		BacklogID:           "test-backlog",
		Name:                "Cancelled Issue",
		Goal:                "This is a test goal for the issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "bug",
		Priority:            "P2",
		Status:              domain.IssueStatusCancelled,
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
		CreatedAt:           time.Now(),
		LastUpdatedAt:       time.Now(),
		CreatedBy:           "test-user",
		LastUpdatedBy:       "test-user",
	}
	if err := service.writer.WriteIssue("test-backlog", issue); err != nil {
		t.Fatalf("Failed to write issue: %v", err)
	}

	input := &domain.IssueDetailInput{
		BacklogID:      "test-backlog",
		IssueID:        "test-backlog-issue-1",
		IncludeDeleted: false,
	}

	_, err := service.GetIssueDetail(input)
	if err == nil {
		t.Error("Expected error for cancelled issue without IncludeDeleted, got nil")
	}
}

func TestValidateUpdateInput_IssueNotFound(t *testing.T) {
	service, _ := setupTestIssueServiceWithBacklog(t, "test-backlog", "same_backlog_only")

	input := &domain.IssueUpdateInput{
		BacklogID: "test-backlog",
		IssueID:   "test-backlog-issue-nonexistent",
		Name:      strPtr("Updated Name"),
	}

	err := service.ValidateUpdateInput(input)
	if err == nil {
		t.Error("Expected validation error for nonexistent issue, got nil")
	}
}

func TestValidateUpdateInput_CancelledIssue(t *testing.T) {
	service, _ := setupTestIssueServiceWithBacklog(t, "test-backlog", "same_backlog_only")

	issue := &domain.Issue{
		ID:                  "test-backlog-issue-1",
		BacklogID:           "test-backlog",
		Name:                "Cancelled Issue",
		Goal:                "This is a test goal for the issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "bug",
		Priority:            "P2",
		Status:              domain.IssueStatusCancelled,
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
		CreatedAt:           time.Now(),
		LastUpdatedAt:       time.Now(),
		CreatedBy:           "test-user",
		LastUpdatedBy:       "test-user",
	}
	if err := service.writer.WriteIssue("test-backlog", issue); err != nil {
		t.Fatalf("Failed to write issue: %v", err)
	}

	input := &domain.IssueUpdateInput{
		BacklogID: "test-backlog",
		IssueID:   "test-backlog-issue-1",
		Name:      strPtr("Updated Name"),
	}

	err := service.ValidateUpdateInput(input)
	if err == nil {
		t.Error("Expected validation error for cancelled issue, got nil")
	}
}

func TestValidateUpdateInput_NameEmpty(t *testing.T) {
	service, _ := setupTestIssueServiceWithBacklog(t, "test-backlog", "same_backlog_only")

	issue := &domain.Issue{
		ID:                  "test-backlog-issue-1",
		BacklogID:           "test-backlog",
		Name:                "Original Name",
		Goal:                "This is a test goal for the issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "bug",
		Priority:            "P2",
		Status:              domain.IssueStatusOpen,
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
		CreatedAt:           time.Now(),
		LastUpdatedAt:       time.Now(),
		CreatedBy:           "test-user",
		LastUpdatedBy:       "test-user",
	}
	if err := service.writer.WriteIssue("test-backlog", issue); err != nil {
		t.Fatalf("Failed to write issue: %v", err)
	}

	input := &domain.IssueUpdateInput{
		BacklogID: "test-backlog",
		IssueID:   "test-backlog-issue-1",
		Name:      strPtr(""),
	}

	err := service.ValidateUpdateInput(input)
	if err == nil {
		t.Error("Expected validation error for empty name, got nil")
	}
}

func TestValidateUpdateInput_IssueTypeInvalid(t *testing.T) {
	service, _ := setupTestIssueServiceWithBacklog(t, "test-backlog", "same_backlog_only")

	issue := &domain.Issue{
		ID:                  "test-backlog-issue-1",
		BacklogID:           "test-backlog",
		Name:                "Test Issue",
		Goal:                "This is a test goal for the issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "bug",
		Priority:            "P2",
		Status:              domain.IssueStatusOpen,
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
		CreatedAt:           time.Now(),
		LastUpdatedAt:       time.Now(),
		CreatedBy:           "test-user",
		LastUpdatedBy:       "test-user",
	}
	if err := service.writer.WriteIssue("test-backlog", issue); err != nil {
		t.Fatalf("Failed to write issue: %v", err)
	}

	input := &domain.IssueUpdateInput{
		BacklogID: "test-backlog",
		IssueID:   "test-backlog-issue-1",
		IssueType: strPtr("invalid_type"),
	}

	err := service.ValidateUpdateInput(input)
	if err == nil {
		t.Error("Expected validation error for invalid issue type, got nil")
	}
}

func TestValidateUpdateInput_PriorityInvalid(t *testing.T) {
	service, _ := setupTestIssueServiceWithBacklog(t, "test-backlog", "same_backlog_only")

	issue := &domain.Issue{
		ID:                  "test-backlog-issue-1",
		BacklogID:           "test-backlog",
		Name:                "Test Issue",
		Goal:                "This is a test goal for the issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "bug",
		Priority:            "P2",
		Status:              domain.IssueStatusOpen,
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
		CreatedAt:           time.Now(),
		LastUpdatedAt:       time.Now(),
		CreatedBy:           "test-user",
		LastUpdatedBy:       "test-user",
	}
	if err := service.writer.WriteIssue("test-backlog", issue); err != nil {
		t.Fatalf("Failed to write issue: %v", err)
	}

	input := &domain.IssueUpdateInput{
		BacklogID: "test-backlog",
		IssueID:   "test-backlog-issue-1",
		Priority:  strPtr("P99"),
	}

	err := service.ValidateUpdateInput(input)
	if err == nil {
		t.Error("Expected validation error for invalid priority, got nil")
	}
}

func TestValidateUpdateInput_StatusInvalid(t *testing.T) {
	service, _ := setupTestIssueServiceWithBacklog(t, "test-backlog", "same_backlog_only")

	issue := &domain.Issue{
		ID:                  "test-backlog-issue-1",
		BacklogID:           "test-backlog",
		Name:                "Test Issue",
		Goal:                "This is a test goal for the issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "bug",
		Priority:            "P2",
		Status:              domain.IssueStatusOpen,
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
		CreatedAt:           time.Now(),
		LastUpdatedAt:       time.Now(),
		CreatedBy:           "test-user",
		LastUpdatedBy:       "test-user",
	}
	if err := service.writer.WriteIssue("test-backlog", issue); err != nil {
		t.Fatalf("Failed to write issue: %v", err)
	}

	input := &domain.IssueUpdateInput{
		BacklogID: "test-backlog",
		IssueID:   "test-backlog-issue-1",
		Status:    strPtr("invalid_status"),
	}

	err := service.ValidateUpdateInput(input)
	if err == nil {
		t.Error("Expected validation error for invalid status, got nil")
	}
}

func TestUpdateIssue_DryRun(t *testing.T) {
	service, _ := setupTestIssueServiceWithBacklog(t, "test-backlog", "same_backlog_only")

	issue := &domain.Issue{
		ID:                  "test-backlog-issue-1",
		BacklogID:           "test-backlog",
		Name:                "Original Name",
		Goal:                "This is a test goal for the issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "bug",
		Priority:            "P2",
		Status:              domain.IssueStatusOpen,
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
		CreatedAt:           time.Now(),
		LastUpdatedAt:       time.Now(),
		CreatedBy:           "test-user",
		LastUpdatedBy:       "test-user",
	}
	if err := service.writer.WriteIssue("test-backlog", issue); err != nil {
		t.Fatalf("Failed to write issue: %v", err)
	}

	input := &domain.IssueUpdateInput{
		BacklogID: "test-backlog",
		IssueID:   "test-backlog-issue-1",
		Name:      strPtr("Updated Name"),
		DryRun:    true,
	}

	changes, err := service.UpdateIssue(input)
	if err != nil {
		t.Fatalf("UpdateIssue failed: %v", err)
	}

	if len(changes) != 1 || changes[0] != "[DRY RUN] Would update issue: test-backlog-issue-1" {
		t.Errorf("Unexpected dry run output: %v", changes)
	}
}

func TestUpdateIssue_Name(t *testing.T) {
	service, _ := setupTestIssueServiceWithBacklog(t, "test-backlog", "same_backlog_only")

	issue := &domain.Issue{
		ID:                  "test-backlog-issue-1",
		BacklogID:           "test-backlog",
		Name:                "Original Name",
		Goal:                "This is a test goal for the issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "bug",
		Priority:            "P2",
		Status:              domain.IssueStatusOpen,
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
		CreatedAt:           time.Now(),
		LastUpdatedAt:       time.Now(),
		CreatedBy:           "test-user",
		LastUpdatedBy:       "test-user",
	}
	if err := service.writer.WriteIssue("test-backlog", issue); err != nil {
		t.Fatalf("Failed to write issue: %v", err)
	}

	input := &domain.IssueUpdateInput{
		BacklogID: "test-backlog",
		IssueID:   "test-backlog-issue-1",
		Name:      strPtr("Updated Name"),
	}

	changes, err := service.UpdateIssue(input)
	if err != nil {
		t.Fatalf("UpdateIssue failed: %v", err)
	}

	found := false
	for _, c := range changes {
		if c == "name" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected 'name' in changes, got %v", changes)
	}

	updated, _ := service.GetIssueDetail(&domain.IssueDetailInput{
		BacklogID:      "test-backlog",
		IssueID:        "test-backlog-issue-1",
		IncludeDeleted: true,
	})
	if updated.Name != "Updated Name" {
		t.Errorf("Issue Name = %q, want %q", updated.Name, "Updated Name")
	}
}

func TestUpdateIssue_Cancel(t *testing.T) {
	service, _ := setupTestIssueServiceWithBacklog(t, "test-backlog", "same_backlog_only")

	issue := &domain.Issue{
		ID:                  "test-backlog-issue-1",
		BacklogID:           "test-backlog",
		Name:                "Test Issue",
		Goal:                "This is a test goal for the issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "bug",
		Priority:            "P2",
		Status:              domain.IssueStatusOpen,
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
		CreatedAt:           time.Now(),
		LastUpdatedAt:       time.Now(),
		CreatedBy:           "test-user",
		LastUpdatedBy:       "test-user",
	}
	if err := service.writer.WriteIssue("test-backlog", issue); err != nil {
		t.Fatalf("Failed to write issue: %v", err)
	}

	input := &domain.IssueUpdateInput{
		BacklogID: "test-backlog",
		IssueID:   "test-backlog-issue-1",
		Cancel:    true,
		Reason:    strPtr("No longer needed"),
	}

	changes, err := service.UpdateIssue(input)
	if err != nil {
		t.Fatalf("UpdateIssue failed: %v", err)
	}

	found := false
	for _, c := range changes {
		if c == "status" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected 'status' in changes, got %v", changes)
	}
}

func TestUpdateIssue_Start(t *testing.T) {
	service, _ := setupTestIssueServiceWithBacklog(t, "test-backlog", "same_backlog_only")

	issue := &domain.Issue{
		ID:                  "test-backlog-issue-1",
		BacklogID:           "test-backlog",
		Name:                "Test Issue",
		Goal:                "This is a test goal for the issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "bug",
		Priority:            "P2",
		Status:              domain.IssueStatusReady,
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
		CreatedAt:           time.Now(),
		LastUpdatedAt:       time.Now(),
		CreatedBy:           "test-user",
		LastUpdatedBy:       "test-user",
	}
	if err := service.writer.WriteIssue("test-backlog", issue); err != nil {
		t.Fatalf("Failed to write issue: %v", err)
	}

	input := &domain.IssueUpdateInput{
		BacklogID: "test-backlog",
		IssueID:   "test-backlog-issue-1",
		Start:     true,
	}

	changes, err := service.UpdateIssue(input)
	if err != nil {
		t.Fatalf("UpdateIssue failed: %v", err)
	}

	found := false
	for _, c := range changes {
		if c == "status" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected 'status' in changes, got %v", changes)
	}
}

func TestUpdateIssue_Resolve(t *testing.T) {
	service, _ := setupTestIssueServiceWithBacklog(t, "test-backlog", "same_backlog_only")

	issue := &domain.Issue{
		ID:                  "test-backlog-issue-1",
		BacklogID:           "test-backlog",
		Name:                "Test Issue",
		Goal:                "This is a test goal for the issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "bug",
		Priority:            "P2",
		Status:              domain.IssueStatusInProgress,
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
		CreatedAt:           time.Now(),
		LastUpdatedAt:       time.Now(),
		CreatedBy:           "test-user",
		LastUpdatedBy:       "test-user",
	}
	if err := service.writer.WriteIssue("test-backlog", issue); err != nil {
		t.Fatalf("Failed to write issue: %v", err)
	}

	input := &domain.IssueUpdateInput{
		BacklogID: "test-backlog",
		IssueID:   "test-backlog-issue-1",
		Resolve:   true,
	}

	changes, err := service.UpdateIssue(input)
	if err != nil {
		t.Fatalf("UpdateIssue failed: %v", err)
	}

	found := false
	for _, c := range changes {
		if c == "status" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected 'status' in changes, got %v", changes)
	}
}

func TestUpdateIssue_WontFix(t *testing.T) {
	service, _ := setupTestIssueServiceWithBacklog(t, "test-backlog", "same_backlog_only")

	issue := &domain.Issue{
		ID:                  "test-backlog-issue-1",
		BacklogID:           "test-backlog",
		Name:                "Test Issue",
		Goal:                "This is a test goal for the issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "bug",
		Priority:            "P2",
		Status:              domain.IssueStatusInProgress,
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
		CreatedAt:           time.Now(),
		LastUpdatedAt:       time.Now(),
		CreatedBy:           "test-user",
		LastUpdatedBy:       "test-user",
	}
	if err := service.writer.WriteIssue("test-backlog", issue); err != nil {
		t.Fatalf("Failed to write issue: %v", err)
	}

	input := &domain.IssueUpdateInput{
		BacklogID: "test-backlog",
		IssueID:   "test-backlog-issue-1",
		WontFix:   true,
		Reason:    strPtr("Not a bug"),
	}

	changes, err := service.UpdateIssue(input)
	if err != nil {
		t.Fatalf("UpdateIssue failed: %v", err)
	}

	found := false
	for _, c := range changes {
		if c == "status" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected 'status' in changes, got %v", changes)
	}
}

func TestUpdateIssue_Reopen(t *testing.T) {
	service, _ := setupTestIssueServiceWithBacklog(t, "test-backlog", "same_backlog_only")

	issue := &domain.Issue{
		ID:                  "test-backlog-issue-1",
		BacklogID:           "test-backlog",
		Name:                "Test Issue",
		Goal:                "This is a test goal for the issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "bug",
		Priority:            "P2",
		Status:              domain.IssueStatusResolved,
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
		CreatedAt:           time.Now(),
		LastUpdatedAt:       time.Now(),
		CreatedBy:           "test-user",
		LastUpdatedBy:       "test-user",
	}
	if err := service.writer.WriteIssue("test-backlog", issue); err != nil {
		t.Fatalf("Failed to write issue: %v", err)
	}

	input := &domain.IssueUpdateInput{
		BacklogID: "test-backlog",
		IssueID:   "test-backlog-issue-1",
		Reopen:    true,
	}

	changes, err := service.UpdateIssue(input)
	if err != nil {
		t.Fatalf("UpdateIssue failed: %v", err)
	}

	found := false
	for _, c := range changes {
		if c == "status" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected 'status' in changes, got %v", changes)
	}
}

func TestFindDependents(t *testing.T) {
	service, _ := setupTestIssueServiceWithBacklog(t, "test-backlog", "same_backlog_only")

	depIssue := &domain.Issue{
		ID:                  "test-backlog-issue-dep",
		BacklogID:           "test-backlog",
		Name:                "Dependency Issue",
		Goal:                "This is a test goal for the dependency issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "bug",
		Priority:            "P2",
		Status:              domain.IssueStatusOpen,
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
		CreatedAt:           time.Now(),
		LastUpdatedAt:       time.Now(),
		CreatedBy:           "test-user",
		LastUpdatedBy:       "test-user",
	}
	if err := service.writer.WriteIssue("test-backlog", depIssue); err != nil {
		t.Fatalf("Failed to write dependency issue: %v", err)
	}

	dependentIssue := &domain.Issue{
		ID:                  "test-backlog-issue-dependent",
		BacklogID:           "test-backlog",
		Name:                "Dependent Issue",
		Goal:                "This is a test goal for the dependent issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "bug",
		Priority:            "P3",
		Status:              domain.IssueStatusBlocked,
		DependsOn:           []string{"test-backlog-issue-dep"},
		AffectedFiles:       []string{"file2.go"},
		AffectedTests:       []string{"TestFile2"},
		ImplementationSteps: []string{"Step 1"},
		CreatedAt:           time.Now(),
		LastUpdatedAt:       time.Now(),
		CreatedBy:           "test-user",
		LastUpdatedBy:       "test-user",
	}
	if err := service.writer.WriteIssue("test-backlog", dependentIssue); err != nil {
		t.Fatalf("Failed to write dependent issue: %v", err)
	}

	dependents, err := service.FindDependents("test-backlog", "test-backlog-issue-dep")
	if err != nil {
		t.Fatalf("FindDependents failed: %v", err)
	}

	if len(dependents) != 1 {
		t.Errorf("Dependents count = %d, want 1", len(dependents))
	}
	if len(dependents) > 0 && dependents[0] != "test-backlog-issue-dependent" {
		t.Errorf("Dependent ID = %q, want %q", dependents[0], "test-backlog-issue-dependent")
	}
}

func TestUnblockDependents(t *testing.T) {
	service, _ := setupTestIssueServiceWithBacklog(t, "test-backlog", "same_backlog_only")

	resolvedIssue := &domain.Issue{
		ID:                  "test-backlog-issue-resolved",
		BacklogID:           "test-backlog",
		Name:                "Resolved Issue",
		Goal:                "This is a test goal for the resolved issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "bug",
		Priority:            "P2",
		Status:              domain.IssueStatusResolved,
		AffectedFiles:       []string{"file1.go"},
		AffectedTests:       []string{"TestFile1"},
		ImplementationSteps: []string{"Step 1"},
		CreatedAt:           time.Now(),
		LastUpdatedAt:       time.Now(),
		CreatedBy:           "test-user",
		LastUpdatedBy:       "test-user",
	}
	if err := service.writer.WriteIssue("test-backlog", resolvedIssue); err != nil {
		t.Fatalf("Failed to write resolved issue: %v", err)
	}

	blockedIssue := &domain.Issue{
		ID:                  "test-backlog-issue-blocked",
		BacklogID:           "test-backlog",
		Name:                "Blocked Issue",
		Goal:                "This is a test goal for the blocked issue that has enough characters to pass validation requirements and exceed minimum length",
		IssueType:           "bug",
		Priority:            "P3",
		Status:              domain.IssueStatusBlocked,
		DependsOn:           []string{"test-backlog-issue-resolved"},
		AffectedFiles:       []string{"file2.go"},
		AffectedTests:       []string{"TestFile2"},
		ImplementationSteps: []string{"Step 1"},
		CreatedAt:           time.Now(),
		LastUpdatedAt:       time.Now(),
		CreatedBy:           "test-user",
		LastUpdatedBy:       "test-user",
	}
	if err := service.writer.WriteIssue("test-backlog", blockedIssue); err != nil {
		t.Fatalf("Failed to write blocked issue: %v", err)
	}

	unblocked, err := service.unblockDependents("test-backlog", "test-backlog-issue-resolved")
	if err != nil {
		t.Fatalf("unblockDependents failed: %v", err)
	}

	if !unblocked {
		t.Error("Expected unblocked = true")
	}

	updated, _ := service.GetIssueDetail(&domain.IssueDetailInput{
		BacklogID:      "test-backlog",
		IssueID:        "test-backlog-issue-blocked",
		IncludeDeleted: true,
	})
	if updated.Status != domain.IssueStatusReady {
		t.Errorf("Issue Status = %q, want %q", updated.Status, domain.IssueStatusReady)
	}
}

func TestBacklogExists(t *testing.T) {
	service, _ := setupTestIssueServiceWithBacklog(t, "test-backlog", "same_backlog_only")

	if !service.BacklogExists("test-backlog") {
		t.Error("Backlog should exist")
	}

	if service.BacklogExists("nonexistent") {
		t.Error("Nonexistent backlog should not exist")
	}
}
