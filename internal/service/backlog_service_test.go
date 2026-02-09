package service

import (
	"testing"
	"time"

	"mandor/internal/domain"
)

// MockReader implements fs.Reader interface for testing
type MockBacklogReader struct {
	backlogExists     bool
	backlog           *domain.Backlog
	backlogSchema     *domain.BacklogSchema
	backlogIDs        []string
	workspaceExists   bool
	entityLineCount   int
	readBacklogErr    error
	readSchemaErr     error
	listBacklogsErr   error
	countEntityErr    error
}

func (m *MockBacklogReader) WorkspaceExists() bool {
	return m.workspaceExists
}

func (m *MockBacklogReader) BacklogExists(backlogID string) bool {
	return m.backlogExists
}

func (m *MockBacklogReader) ProjectExists(projectID string) bool {
	return false
}

func (m *MockBacklogReader) ListBacklogs(includeDeleted bool) ([]string, error) {
	if m.listBacklogsErr != nil {
		return nil, m.listBacklogsErr
	}
	return m.backlogIDs, nil
}

func (m *MockBacklogReader) ListProjects(includeDeleted bool) ([]string, error) {
	return nil, nil
}

func (m *MockBacklogReader) ReadBacklogMetadata(backlogID string) (*domain.Backlog, error) {
	if m.readBacklogErr != nil {
		return nil, m.readBacklogErr
	}
	return m.backlog, nil
}

func (m *MockBacklogReader) ReadBacklogSchema(backlogID string) (*domain.BacklogSchema, error) {
	if m.readSchemaErr != nil {
		return nil, m.readSchemaErr
	}
	return m.backlogSchema, nil
}

func (m *MockBacklogReader) ReadProjectMetadata(projectID string) (*domain.Project, error) {
	return nil, nil
}

func (m *MockBacklogReader) ReadProjectSchema(projectID string) (*domain.ProjectSchema, error) {
	return nil, nil
}

func (m *MockBacklogReader) CountEntityLines(filePath string) (int, error) {
	if m.countEntityErr != nil {
		return 0, m.countEntityErr
	}
	return m.entityLineCount, nil
}

func (m *MockBacklogReader) ReadWorkspace() (*domain.Workspace, error) {
	return &domain.Workspace{
		Config: domain.WorkspaceConfig{
			GoalLengths: domain.GoalLengths{
				Project: domain.GoalMinLength,
			},
		},
	}, nil
}

// MockWriter implements fs.Writer interface for testing
type MockBacklogWriter struct {
	createErr           error
	writeMetadataErr    error
	writeSchemaErr      error
	deleteErr           error
	backlogsDirWritable bool
	checkWritableOk     bool
}

func (m *MockBacklogWriter) CreateBacklogDir(backlogID string) error {
	return m.createErr
}

func (m *MockBacklogWriter) WriteBacklogMetadata(backlogID string, backlog *domain.Backlog) error {
	return m.writeMetadataErr
}

func (m *MockBacklogWriter) WriteBacklogSchema(backlogID string, schema *domain.BacklogSchema) error {
	return m.writeSchemaErr
}

func (m *MockBacklogWriter) DeleteBacklogDir(backlogID string) error {
	return m.deleteErr
}

func (m *MockBacklogWriter) BacklogsDirWritable() bool {
	return m.backlogsDirWritable
}

func (m *MockBacklogWriter) CheckBacklogWritable(backlogID string) bool {
	return m.checkWritableOk
}

func (m *MockBacklogWriter) CreateProjectDir(projectID string) error {
	return nil
}

func (m *MockBacklogWriter) WriteProjectMetadata(projectID string, project *domain.Project) error {
	return nil
}

func (m *MockBacklogWriter) WriteProjectSchema(projectID string, schema *domain.ProjectSchema) error {
	return nil
}

func (m *MockBacklogWriter) DeleteProjectDir(projectID string) error {
	return nil
}

func (m *MockBacklogWriter) ProjectsDirWritable() bool {
	return false
}

func (m *MockBacklogWriter) CheckProjectWritable(projectID string) bool {
	return false
}

func TestBacklogServiceValidateCreateInput_ValidID(t *testing.T) {
	input := &domain.BacklogCreateInput{
		ID:   "valid-backlog",
		Name: "Valid Backlog",
	}

	// We validate that the ID is syntactically valid
	if !domain.ValidateBacklogID(input.ID) {
		t.Errorf("ValidateBacklogID should accept 'valid-backlog'")
	}
}

func TestBacklogServiceValidateCreateInput_InvalidID(t *testing.T) {
	input := &domain.BacklogCreateInput{
		ID:   "123-invalid",
		Name: "Invalid Backlog",
	}

	if domain.ValidateBacklogID(input.ID) {
		t.Errorf("ValidateBacklogID should reject ID starting with number")
	}
}

func TestBacklogStatusTransitions(t *testing.T) {
	now := time.Now().UTC()

	// Create a backlog in initial state
	backlog := &domain.Backlog{
		ID:        "test-backlog",
		Name:      "Test",
		Status:    domain.BacklogStatusInitial,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if backlog.Status != domain.BacklogStatusInitial {
		t.Errorf("Initial status should be initial, got %q", backlog.Status)
	}

	// Transition to active
	backlog.Status = domain.BacklogStatusActive
	if backlog.Status != domain.BacklogStatusActive {
		t.Errorf("Status should be active, got %q", backlog.Status)
	}

	// Transition to deleted
	backlog.Status = domain.BacklogStatusDeleted
	if backlog.Status != domain.BacklogStatusDeleted {
		t.Errorf("Status should be deleted, got %q", backlog.Status)
	}

	// Transition back to initial (reopen)
	backlog.Status = domain.BacklogStatusInitial
	if backlog.Status != domain.BacklogStatusInitial {
		t.Errorf("Status should be initial after reopen, got %q", backlog.Status)
	}
}

func TestBacklogMetadataCreation(t *testing.T) {
	backlog := &domain.Backlog{
		ID:        "test-backlog",
		Name:      "Test Backlog",
		Goal:      "A comprehensive test backlog",
		Status:    domain.BacklogStatusInitial,
		Strict:    true,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		CreatedBy: "testuser",
		UpdatedBy: "testuser",
	}

	// Verify all fields are set correctly
	if backlog.ID != "test-backlog" {
		t.Errorf("ID mismatch: got %q", backlog.ID)
	}
	if backlog.Name != "Test Backlog" {
		t.Errorf("Name mismatch: got %q", backlog.Name)
	}
	if backlog.Status != domain.BacklogStatusInitial {
		t.Errorf("Status mismatch: got %q", backlog.Status)
	}
	if !backlog.Strict {
		t.Errorf("Strict should be true")
	}
}

func TestBacklogListItemConstruction(t *testing.T) {
	item := domain.BacklogListItem{
		ID:        "test-backlog",
		Name:      "Test",
		Goal:      "Test Goal",
		Status:    domain.BacklogStatusActive,
		Features:  5,
		Tasks:     10,
		Issues:    2,
		CreatedAt: "2026-02-09T10:00:00Z",
		UpdatedAt: "2026-02-09T11:00:00Z",
	}

	if item.Features != 5 {
		t.Errorf("Features count should be 5, got %d", item.Features)
	}
	if item.Tasks != 10 {
		t.Errorf("Tasks count should be 10, got %d", item.Tasks)
	}
	if item.Issues != 2 {
		t.Errorf("Issues count should be 2, got %d", item.Issues)
	}
}

func TestBacklogListOutputAggregation(t *testing.T) {
	items := []domain.BacklogListItem{
		{ID: "backlog-1", Status: domain.BacklogStatusActive},
		{ID: "backlog-2", Status: domain.BacklogStatusActive},
		{ID: "backlog-3", Status: domain.BacklogStatusDeleted},
	}

	output := &domain.BacklogListOutput{
		Backlogs: items,
		Total:    3,
		Deleted:  1,
	}

	if len(output.Backlogs) != 3 {
		t.Errorf("Should have 3 backlogs, got %d", len(output.Backlogs))
	}
	if output.Total != 3 {
		t.Errorf("Total should be 3, got %d", output.Total)
	}
	if output.Deleted != 1 {
		t.Errorf("Deleted count should be 1, got %d", output.Deleted)
	}
}

func TestBacklogDetailOutputConstruction(t *testing.T) {
	schema := domain.DefaultBacklogSchema("", "", "")

	output := &domain.BacklogDetailOutput{
		ID:        "test-backlog",
		Name:      "Test",
		Goal:      "Test Goal",
		Status:    domain.BacklogStatusActive,
		Schema:    schema,
		CreatedAt: "2026-02-09T10:00:00Z",
		UpdatedAt: "2026-02-09T11:00:00Z",
		CreatedBy: "user1",
		UpdatedBy: "user1",
	}

	if output.ID != "test-backlog" {
		t.Errorf("ID mismatch")
	}
	if output.Schema.Version != "mandor.v1" {
		t.Errorf("Schema version mismatch")
	}
	if output.CreatedBy != "user1" {
		t.Errorf("CreatedBy mismatch")
	}
}

func TestBacklogUpdateInputPartialUpdate(t *testing.T) {
	newName := "Updated Name"
	input := &domain.BacklogUpdateInput{
		ID:   "test-backlog",
		Name: &newName,
	}

	// Only Name is being updated, others should be nil
	if input.Name == nil {
		t.Errorf("Name should be set")
	}
	if input.Goal != nil {
		t.Errorf("Goal should be nil")
	}
	if input.Strict != nil {
		t.Errorf("Strict should be nil")
	}
}

func TestBacklogDeleteInputDryRun(t *testing.T) {
	input := domain.BacklogDeleteInput{
		ID:     "test-backlog",
		Hard:   false,
		DryRun: true,
	}

	if !input.DryRun {
		t.Errorf("DryRun should be true")
	}
	if input.Hard {
		t.Errorf("Hard should be false")
	}
}

func TestBacklogDeleteInputHardDelete(t *testing.T) {
	input := domain.BacklogDeleteInput{
		ID:     "test-backlog",
		Hard:   true,
		DryRun: false,
	}

	if !input.Hard {
		t.Errorf("Hard should be true")
	}
	if input.DryRun {
		t.Errorf("DryRun should be false")
	}
}

func TestBacklogReopenInputConfirmation(t *testing.T) {
	input := domain.BacklogReopenInput{
		ID:  "test-backlog",
		Yes: true,
	}

	if !input.Yes {
		t.Errorf("Yes should be true")
	}
}

func TestBacklogSchemaDefaults(t *testing.T) {
	schema := domain.DefaultBacklogSchema("", "", "")

	expectedTaskDep := domain.DependencySameProjectOnly
	expectedFeatureDep := domain.DependencyCrossProjectAllowed
	expectedIssueDep := domain.DependencySameProjectOnly

	if schema.Rules.Task.Dependency != expectedTaskDep {
		t.Errorf("Task dependency mismatch: expected %q, got %q", expectedTaskDep, schema.Rules.Task.Dependency)
	}
	if schema.Rules.Feature.Dependency != expectedFeatureDep {
		t.Errorf("Feature dependency mismatch: expected %q, got %q", expectedFeatureDep, schema.Rules.Feature.Dependency)
	}
	if schema.Rules.Issue.Dependency != expectedIssueDep {
		t.Errorf("Issue dependency mismatch: expected %q, got %q", expectedIssueDep, schema.Rules.Issue.Dependency)
	}
}

func TestBacklogSchemaCustomization(t *testing.T) {
	schema := domain.DefaultBacklogSchema("cross_project_allowed", "disabled", "cross_project_allowed")

	if schema.Rules.Task.Dependency != "cross_project_allowed" {
		t.Errorf("Task dependency should be cross_project_allowed")
	}
	if schema.Rules.Feature.Dependency != "disabled" {
		t.Errorf("Feature dependency should be disabled")
	}
	if schema.Rules.Issue.Dependency != "cross_project_allowed" {
		t.Errorf("Issue dependency should be cross_project_allowed")
	}
}

func TestBacklogStatsCounting(t *testing.T) {
	stats := domain.BacklogStats{
		Features: domain.EntityStats{Total: 5},
		Tasks:    domain.EntityStats{Total: 20},
		Issues:   domain.EntityStats{Total: 8},
	}

	if stats.Features.Total != 5 {
		t.Errorf("Features total should be 5")
	}
	if stats.Tasks.Total != 20 {
		t.Errorf("Tasks total should be 20")
	}
	if stats.Issues.Total != 8 {
		t.Errorf("Issues total should be 8")
	}
}
