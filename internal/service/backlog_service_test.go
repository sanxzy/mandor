package service

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"mandor/internal/domain"
	"mandor/internal/fs"
)

// Helper to set up test service with temporary directory
func setupTestBacklogService(t *testing.T) (*BacklogService, string) {
	tmpDir := t.TempDir()

	// Initialize workspace first
	paths := &fs.Paths{WorkspaceRoot: tmpDir}
	writer := fs.NewWriter(paths)
	reader := fs.NewReader(paths)

	// Create necessary directories (.mandor/backlogs, .mandor/briefs, etc.)
	os.MkdirAll(filepath.Join(tmpDir, ".mandor", "backlogs"), 0755)
	os.MkdirAll(filepath.Join(tmpDir, ".mandor", "briefs"), 0755)
	os.MkdirAll(filepath.Join(tmpDir, ".mandor", "features"), 0755)

	// Write workspace file
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

	service := &BacklogService{
		reader: reader,
		writer: writer,
		paths:  paths,
	}

	return service, tmpDir
}

// ==========================
// Constructor Tests
// ==========================

func TestNewBacklogServiceWithPaths(t *testing.T) {
	tmpDir := t.TempDir()
	paths := &fs.Paths{WorkspaceRoot: tmpDir}

	service := NewBacklogServiceWithPaths(paths)

	if service == nil {
		t.Error("NewBacklogServiceWithPaths returned nil")
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

// ==========================
// Validation Tests
// ==========================

func TestValidateCreateInput_InvalidID(t *testing.T) {
	service, _ := setupTestBacklogService(t)

	tests := []struct {
		name string
		id   string
	}{
		{"ID starts with number", "123-backlog"},
		{"ID with spaces", "my backlog"},
		{"ID with special chars", "backlog!@#"},
		{"Empty ID", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &domain.BacklogCreateInput{
				ID:   tt.id,
				Name: "Test Backlog",
				Goal: "This is a test goal for backlog",
			}

			err := service.ValidateCreateInput(input)
			if err == nil {
				t.Errorf("Expected error for invalid ID %q, got nil", tt.id)
			}
		})
	}
}

func TestValidateCreateInput_ValidID(t *testing.T) {
	service, tmpDir := setupTestBacklogService(t)
	_ = tmpDir

	input := &domain.BacklogCreateInput{
		ID:   "test-backlog",
		Name: "Test Backlog",
		Goal: "This is a valid goal for testing the backlog creation system with plenty of detailed content",
	}

	err := service.ValidateCreateInput(input)
	if err != nil {
		t.Errorf("Expected no error for valid input, got: %v", err)
	}
}

func TestValidateCreateInput_BacklogExists(t *testing.T) {
	service, tmpDir := setupTestBacklogService(t)

	// Create first backlog
	createInput := &domain.BacklogCreateInput{
		ID:   "existing-backlog",
		Name: "Existing Backlog",
		Goal: "This is a test goal for testing duplicate backlog validation",
	}
	service.CreateBacklog(createInput)

	// Try to create again
	err := service.ValidateCreateInput(createInput)
	if err == nil {
		t.Error("Expected error for duplicate backlog, got nil")
	}

	_ = tmpDir
}

func TestValidateUpdateInput_BacklogNotFound(t *testing.T) {
	service, _ := setupTestBacklogService(t)

	input := &domain.BacklogUpdateInput{
		ID:   "nonexistent-backlog",
		Name: strPtr("Updated Name"),
	}

	err := service.ValidateUpdateInput(input)
	if err == nil {
		t.Error("Expected error for nonexistent backlog, got nil")
	}
}

func TestValidateUpdateInput_DeletedBacklog(t *testing.T) {
	service, _ := setupTestBacklogService(t)

	// Create and delete backlog
	createInput := &domain.BacklogCreateInput{
		ID:   "test-backlog",
		Name: "Test Backlog",
		Goal: "This is a test goal for testing deleted backlog updates",
	}
	service.CreateBacklog(createInput)

	deleteInput := &domain.BacklogDeleteInput{
		ID:     "test-backlog",
		DryRun: false,
		Hard:   false,
	}
	service.DeleteBacklog(deleteInput)

	// Try to update
	updateInput := &domain.BacklogUpdateInput{
		ID:   "test-backlog",
		Name: strPtr("Updated Name"),
	}
	err := service.ValidateUpdateInput(updateInput)
	if err == nil {
		t.Error("Expected error for deleted backlog, got nil")
	}
}

func TestValidateDeleteInput_NotDeleted(t *testing.T) {
	service, _ := setupTestBacklogService(t)

	// Create backlog
	createInput := &domain.BacklogCreateInput{
		ID:   "test-backlog",
		Name: "Test Backlog",
		Goal: "This is a test goal for testing backlog deletion validation",
	}
	service.CreateBacklog(createInput)

	// Validate deletion
	deleteInput := &domain.BacklogDeleteInput{
		ID:     "test-backlog",
		DryRun: true,
		Hard:   false,
	}
	err := service.ValidateDeleteInput(deleteInput)
	if err != nil {
		t.Errorf("Expected no error for valid deletion, got: %v", err)
	}
}

func TestValidateReopenInput_NotDeleted(t *testing.T) {
	service, _ := setupTestBacklogService(t)

	// Create backlog
	createInput := &domain.BacklogCreateInput{
		ID:   "test-backlog",
		Name: "Test Backlog",
		Goal: "This is a test goal for testing backlog reopen validation",
	}
	service.CreateBacklog(createInput)

	// Try to reopen (not deleted)
	reopenInput := &domain.BacklogReopenInput{
		ID: "test-backlog",
	}
	err := service.ValidateReopenInput(reopenInput)
	if err == nil {
		t.Error("Expected error for non-deleted backlog, got nil")
	}
}

// ==========================
// Create Tests
// ==========================

func TestCreateBacklog_Success(t *testing.T) {
	service, _ := setupTestBacklogService(t)

	input := &domain.BacklogCreateInput{
		ID:         "new-backlog",
		Name:       "New Backlog",
		Goal:       "This is a comprehensive test goal for creating a new backlog",
		TaskDep:    "",
		FeatureDep: "",
		IssueDep:   "",
		Strict:     false,
	}

	err := service.CreateBacklog(input)
	if err != nil {
		t.Fatalf("CreateBacklog failed: %v", err)
	}

	// Verify backlog was created
	exists := service.reader.BacklogExists("new-backlog")
	if !exists {
		t.Error("Backlog was not created")
	}

	// Verify metadata
	backlog, err := service.reader.ReadBacklogMetadata("new-backlog")
	if err != nil {
		t.Fatalf("Failed to read backlog metadata: %v", err)
	}

	if backlog.ID != "new-backlog" {
		t.Errorf("Backlog ID = %q, want %q", backlog.ID, "new-backlog")
	}
	if backlog.Name != "New Backlog" {
		t.Errorf("Backlog Name = %q, want %q", backlog.Name, "New Backlog")
	}
	if backlog.Status != domain.BacklogStatusInitial {
		t.Errorf("Backlog Status = %q, want %q", backlog.Status, domain.BacklogStatusInitial)
	}

	// Verify schema was created
	schema, err := service.reader.ReadBacklogSchema("new-backlog")
	if err != nil {
		t.Fatalf("Failed to read backlog schema: %v", err)
	}
	if schema == nil {
		t.Error("Backlog schema is nil")
	}
}

func TestCreateBacklog_WithStrictMode(t *testing.T) {
	service, _ := setupTestBacklogService(t)

	input := &domain.BacklogCreateInput{
		ID:     "strict-backlog",
		Name:   "Strict Backlog",
		Goal:   "This is a test goal for creating a backlog with strict mode enabled",
		Strict: true,
	}

	err := service.CreateBacklog(input)
	if err != nil {
		t.Fatalf("CreateBacklog failed: %v", err)
	}

	backlog, err := service.reader.ReadBacklogMetadata("strict-backlog")
	if err != nil {
		t.Fatalf("Failed to read backlog: %v", err)
	}

	if !backlog.Strict {
		t.Error("Backlog Strict mode not set")
	}
}

// ==========================
// Read Tests
// ==========================

func TestListBacklogs_Empty(t *testing.T) {
	service, _ := setupTestBacklogService(t)

	output, err := service.ListBacklogs(false, false)
	if err != nil {
		t.Fatalf("ListBacklogs failed: %v", err)
	}

	if output.Total != 0 {
		t.Errorf("Total = %d, want 0", output.Total)
	}
	if len(output.Backlogs) != 0 {
		t.Errorf("Backlogs count = %d, want 0", len(output.Backlogs))
	}
}

func TestListBacklogs_Multiple(t *testing.T) {
	service, _ := setupTestBacklogService(t)

	// Create multiple backlogs
	for i := 1; i <= 3; i++ {
		id := "backlog-" + string(rune(48+i))
		input := &domain.BacklogCreateInput{
			ID:   id,
			Name: "Backlog " + string(rune(48+i)),
			Goal: "This is a test goal for listing multiple backlogs successfully",
		}
		service.CreateBacklog(input)
	}

	output, err := service.ListBacklogs(false, false)
	if err != nil {
		t.Fatalf("ListBacklogs failed: %v", err)
	}

	if output.Total != 3 {
		t.Errorf("Total = %d, want 3", output.Total)
	}
	if len(output.Backlogs) != 3 {
		t.Errorf("Backlogs count = %d, want 3", len(output.Backlogs))
	}
}

func TestListBacklogs_IncludeGoal(t *testing.T) {
	service, _ := setupTestBacklogService(t)

	goal := "This is a test goal for including goals in backlog listing"
	input := &domain.BacklogCreateInput{
		ID:   "test-backlog",
		Name: "Test Backlog",
		Goal: goal,
	}
	service.CreateBacklog(input)

	output, err := service.ListBacklogs(false, true)
	if err != nil {
		t.Fatalf("ListBacklogs failed: %v", err)
	}

	if len(output.Backlogs) != 1 {
		t.Fatalf("Expected 1 backlog, got %d", len(output.Backlogs))
	}

	if output.Backlogs[0].Goal == "" {
		t.Error("Goal was not included in backlog list item")
	}
}

func TestGetBacklogDetail(t *testing.T) {
	service, _ := setupTestBacklogService(t)

	createInput := &domain.BacklogCreateInput{
		ID:   "detail-backlog",
		Name: "Detail Backlog",
		Goal: "This is a test goal for getting backlog details successfully",
	}
	service.CreateBacklog(createInput)

	detail, err := service.GetBacklogDetail("detail-backlog")
	if err != nil {
		t.Fatalf("GetBacklogDetail failed: %v", err)
	}

	if detail.ID != "detail-backlog" {
		t.Errorf("Detail.ID = %q, want %q", detail.ID, "detail-backlog")
	}
	if detail.Name != "Detail Backlog" {
		t.Errorf("Detail.Name = %q, want %q", detail.Name, "Detail Backlog")
	}
	if detail.Status != domain.BacklogStatusInitial {
		t.Errorf("Detail.Status = %q, want %q", detail.Status, domain.BacklogStatusInitial)
	}
}

func TestGetBacklog(t *testing.T) {
	service, _ := setupTestBacklogService(t)

	createInput := &domain.BacklogCreateInput{
		ID:   "get-backlog",
		Name: "Get Backlog",
		Goal: "This is a test goal for getting backlog metadata directly",
	}
	service.CreateBacklog(createInput)

	backlog, err := service.GetBacklog("get-backlog")
	if err != nil {
		t.Fatalf("GetBacklog failed: %v", err)
	}

	if backlog.ID != "get-backlog" {
		t.Errorf("Backlog.ID = %q, want %q", backlog.ID, "get-backlog")
	}
	if backlog.Name != "Get Backlog" {
		t.Errorf("Backlog.Name = %q, want %q", backlog.Name, "Get Backlog")
	}
}

// ==========================
// Update Tests
// ==========================

func TestUpdateBacklog_Name(t *testing.T) {
	service, _ := setupTestBacklogService(t)

	// Create backlog
	createInput := &domain.BacklogCreateInput{
		ID:   "update-backlog",
		Name: "Original Name",
		Goal: "This is a test goal for updating backlog name",
	}
	service.CreateBacklog(createInput)

	// Update name
	newName := "Updated Name"
	updateInput := &domain.BacklogUpdateInput{
		ID:   "update-backlog",
		Name: &newName,
	}

	changes, err := service.UpdateBacklog(updateInput)
	if err != nil {
		t.Fatalf("UpdateBacklog failed: %v", err)
	}

	// Verify changes
	if len(changes) != 1 || changes[0] != "name" {
		t.Errorf("Changes = %v, want [name]", changes)
	}

	// Verify update
	backlog, _ := service.reader.ReadBacklogMetadata("update-backlog")
	if backlog.Name != newName {
		t.Errorf("Backlog.Name = %q, want %q", backlog.Name, newName)
	}
}

func TestUpdateBacklog_Goal(t *testing.T) {
	service, _ := setupTestBacklogService(t)

	// Create backlog with long goal (minimum 500 chars)
	longGoal := "This is an original goal for testing goal update functionality. " +
		"We need a lot of detailed information about what this backlog is trying to achieve. " +
		"The goal provides essential context for all team members about the overall direction and purpose. " +
		"This backlog will enable better planning and coordination across the entire organization. " +
		"Teams will understand priorities and dependencies and be able to work together more effectively. " +
		"The comprehensive goals help align everyone with the project vision and ensure successful delivery. " +
		"Additional context and details are essential for proper planning and execution across all phases of development and beyond."

	createInput := &domain.BacklogCreateInput{
		ID:   "goal-backlog",
		Name: "Goal Backlog",
		Goal: longGoal,
	}
	service.CreateBacklog(createInput)

	// Update goal with another long goal string
	newGoal := "This is an updated goal for testing backlog goal update. " +
		"The new goal includes comprehensive details about improved functionality and enhanced capabilities. " +
		"This updated direction reflects lessons learned and best practices from the team. " +
		"It provides clear guidance for implementation and ensures all stakeholders understand the vision. " +
		"The new goal encompasses multiple dimensions and addresses various concerns and requirements. " +
		"It represents a significant improvement over previous iterations and sets the stage for success. " +
		"Additional refinements and adjustments will be made as needed based on feedback and progress."

	updateInput := &domain.BacklogUpdateInput{
		ID:   "goal-backlog",
		Goal: &newGoal,
	}

	changes, err := service.UpdateBacklog(updateInput)
	if err != nil {
		t.Fatalf("UpdateBacklog failed: %v", err)
	}

	if len(changes) != 1 || changes[0] != "goal" {
		t.Errorf("Changes = %v, want [goal]", changes)
	}

	backlog, _ := service.reader.ReadBacklogMetadata("goal-backlog")
	if backlog.Goal != newGoal {
		t.Errorf("Backlog.Goal mismatch")
	}
}

func TestUpdateBacklog_Strict(t *testing.T) {
	service, _ := setupTestBacklogService(t)

	createInput := &domain.BacklogCreateInput{
		ID:     "strict-update-backlog",
		Name:   "Strict Update Backlog",
		Goal:   "This is a test goal for updating backlog strict mode",
		Strict: false,
	}
	service.CreateBacklog(createInput)

	// Update strict
	strict := true
	updateInput := &domain.BacklogUpdateInput{
		ID:     "strict-update-backlog",
		Strict: &strict,
	}

	changes, err := service.UpdateBacklog(updateInput)
	if err != nil {
		t.Fatalf("UpdateBacklog failed: %v", err)
	}

	if len(changes) != 1 || changes[0] != "strict" {
		t.Errorf("Changes = %v, want [strict]", changes)
	}

	backlog, _ := service.reader.ReadBacklogMetadata("strict-update-backlog")
	if !backlog.Strict {
		t.Error("Backlog Strict mode not updated")
	}
}

func TestUpdateBacklog_MultipleFields(t *testing.T) {
	service, _ := setupTestBacklogService(t)

	multiGoal := "This is an original goal for multiple field testing. " +
		"We include comprehensive details about the purpose and scope. " +
		"Additional context explains what team members should know when working. " +
		"This goal provides clear direction for all stakeholders involved. " +
		"It encompasses various aspects and dimensions of the project. " +
		"Team members will have complete information to make informed decisions. " +
		"The goal ensures proper alignment and coordination throughout all phases of the project."

	createInput := &domain.BacklogCreateInput{
		ID:     "multi-update-backlog",
		Name:   "Original",
		Goal:   multiGoal,
		Strict: false,
	}
	service.CreateBacklog(createInput)

	newName := "Updated"
	newGoal := "This is an updated goal for multiple field testing with improvements and enhancements. " +
		"We include even more comprehensive details about changes and advancements to the system. " +
		"Additional context explains the new direction and updated approach for implementation. " +
		"This updated goal provides clearer guidance for all stakeholders and team members involved. " +
		"It encompasses improvements and refinements from lessons learned in previous iterations. " +
		"Team members will understand the new vision and objectives very clearly and completely. " +
		"The updated goal ensures better alignment and coordination going forward for success. " +
		"Additional details provide comprehensive information for proper execution and delivery."
	strict := true
	updateInput := &domain.BacklogUpdateInput{
		ID:     "multi-update-backlog",
		Name:   &newName,
		Goal:   &newGoal,
		Strict: &strict,
	}

	changes, err := service.UpdateBacklog(updateInput)
	if err != nil {
		t.Fatalf("UpdateBacklog failed: %v", err)
	}

	if len(changes) != 3 {
		t.Errorf("Changes count = %d, want 3", len(changes))
	}

	backlog, _ := service.reader.ReadBacklogMetadata("multi-update-backlog")
	if backlog.Name != newName {
		t.Errorf("Name not updated")
	}
	if backlog.Goal != newGoal {
		t.Errorf("Goal not updated")
	}
	if !backlog.Strict {
		t.Errorf("Strict not updated")
	}
}

func TestUpdateBacklog_InvalidGoalLength(t *testing.T) {
	service, _ := setupTestBacklogService(t)

	createInput := &domain.BacklogCreateInput{
		ID:   "invalid-goal-backlog",
		Name: "Invalid Goal Backlog",
		Goal: "This is a valid test goal for creating backlog",
	}
	service.CreateBacklog(createInput)

	// Try to update with too short goal
	shortGoal := "Too short"
	updateInput := &domain.BacklogUpdateInput{
		ID:   "invalid-goal-backlog",
		Goal: &shortGoal,
	}

	_, err := service.UpdateBacklog(updateInput)
	if err == nil {
		t.Error("Expected error for short goal, got nil")
	}
}

// ==========================
// Delete Tests
// ==========================

func TestDeleteBacklog_SoftDelete(t *testing.T) {
	service, _ := setupTestBacklogService(t)

	createInput := &domain.BacklogCreateInput{
		ID:   "soft-delete-backlog",
		Name: "Soft Delete Backlog",
		Goal: "This is a test goal for testing backlog soft delete",
	}
	service.CreateBacklog(createInput)

	deleteInput := &domain.BacklogDeleteInput{
		ID:     "soft-delete-backlog",
		DryRun: false,
		Hard:   false,
	}

	msg, err := service.DeleteBacklog(deleteInput)
	if err != nil {
		t.Fatalf("DeleteBacklog failed: %v", err)
	}

	if msg == "" {
		t.Error("Delete message is empty")
	}

	// Verify soft delete (backlog still exists but marked deleted)
	backlog, err := service.reader.ReadBacklogMetadata("soft-delete-backlog")
	if err != nil {
		t.Fatalf("Failed to read deleted backlog: %v", err)
	}

	if backlog.Status != domain.BacklogStatusDeleted {
		t.Errorf("Backlog.Status = %q, want %q", backlog.Status, domain.BacklogStatusDeleted)
	}
}

func TestDeleteBacklog_DryRun(t *testing.T) {
	service, _ := setupTestBacklogService(t)

	createInput := &domain.BacklogCreateInput{
		ID:   "dryrun-backlog",
		Name: "Dry Run Backlog",
		Goal: "This is a test goal for testing backlog dry run delete",
	}
	service.CreateBacklog(createInput)

	deleteInput := &domain.BacklogDeleteInput{
		ID:     "dryrun-backlog",
		DryRun: true,
		Hard:   false,
	}

	msg, err := service.DeleteBacklog(deleteInput)
	if err != nil {
		t.Fatalf("DeleteBacklog failed: %v", err)
	}

	if msg == "" || len(msg) < 9 || msg[0:9] != "[DRY RUN]" {
		t.Errorf("Expected dry run message, got: %s", msg)
	}

	// Verify backlog still exists and not deleted
	backlog, err := service.reader.ReadBacklogMetadata("dryrun-backlog")
	if err != nil {
		t.Fatalf("Failed to read backlog: %v", err)
	}

	if backlog.Status != domain.BacklogStatusInitial {
		t.Error("Backlog status changed during dry run")
	}
}

// ==========================
// Reopen Tests
// ==========================

func TestReopenBacklog_Success(t *testing.T) {
	service, _ := setupTestBacklogService(t)

	createInput := &domain.BacklogCreateInput{
		ID:   "reopen-backlog",
		Name: "Reopen Backlog",
		Goal: "This is a test goal for testing backlog reopen",
	}
	service.CreateBacklog(createInput)

	// Delete first
	deleteInput := &domain.BacklogDeleteInput{
		ID:     "reopen-backlog",
		DryRun: false,
		Hard:   false,
	}
	service.DeleteBacklog(deleteInput)

	// Reopen
	reopenInput := &domain.BacklogReopenInput{
		ID: "reopen-backlog",
	}

	msg, err := service.ReopenBacklog(reopenInput)
	if err != nil {
		t.Fatalf("ReopenBacklog failed: %v", err)
	}

	if msg == "" {
		t.Error("Reopen message is empty")
	}

	// Verify status
	backlog, _ := service.reader.ReadBacklogMetadata("reopen-backlog")
	if backlog.Status != domain.BacklogStatusInitial {
		t.Errorf("Backlog.Status = %q, want %q", backlog.Status, domain.BacklogStatusInitial)
	}
}

// ==========================
// Dependency Rule Tests
// ==========================

func TestUpdateBacklog_TaskDependencyRule(t *testing.T) {
	service, _ := setupTestBacklogService(t)

	createInput := &domain.BacklogCreateInput{
		ID:   "dep-backlog",
		Name: "Dependency Backlog",
		Goal: "This is a test goal for testing backlog dependency rule",
	}
	service.CreateBacklog(createInput)

	// Update with valid dependency rule
	rule := domain.DependencyCrossBacklogAllowed
	updateInput := &domain.BacklogUpdateInput{
		ID:      "dep-backlog",
		TaskDep: &rule,
	}

	changes, err := service.UpdateBacklog(updateInput)
	if err != nil {
		t.Fatalf("UpdateBacklog failed: %v", err)
	}

	if len(changes) == 0 {
		t.Error("No changes returned for dependency update")
	}

	// Verify schema was updated
	schema, _ := service.reader.ReadBacklogSchema("dep-backlog")
	if schema.Rules.Task.Dependency != rule {
		t.Errorf("Task dependency = %q, want %q", schema.Rules.Task.Dependency, rule)
	}
}

func TestUpdateBacklog_InvalidDependencyRule(t *testing.T) {
	service, _ := setupTestBacklogService(t)

	createInput := &domain.BacklogCreateInput{
		ID:   "invalid-dep-backlog",
		Name: "Invalid Dependency Backlog",
		Goal: "This is a test goal for testing invalid dependency rule",
	}
	service.CreateBacklog(createInput)

	// Update with invalid dependency rule
	invalidRule := "invalid_rule"
	updateInput := &domain.BacklogUpdateInput{
		ID:      "invalid-dep-backlog",
		TaskDep: &invalidRule,
	}

	_, err := service.UpdateBacklog(updateInput)
	if err == nil {
		t.Error("Expected error for invalid dependency rule, got nil")
	}
}

// ==========================
// Helper Functions
// ==========================

func strPtr(s string) *string {
	return &s
}

func TestWorkspaceInitialized(t *testing.T) {
	service, _ := setupTestBacklogService(t)

	initialized := service.WorkspaceInitialized()
	if !initialized {
		t.Error("Workspace should be initialized")
	}
}

func TestGetBacklogGoalMinLength(t *testing.T) {
	service, _ := setupTestBacklogService(t)

	minLen := service.GetBacklogGoalMinLength()
	if minLen <= 0 {
		t.Errorf("Goal min length = %d, want > 0", minLen)
	}
	if minLen != 500 {
		t.Errorf("Goal min length = %d, want 500", minLen)
	}
}
