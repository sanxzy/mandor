package backlog

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/fs"
	"mandor/internal/service"
)

// TestIntegration sets up a temporary workspace for integration testing
type TestIntegration struct {
	tmpDir string
	paths  *fs.Paths
	svc    *service.BacklogService
	t      *testing.T
}

// NewTestIntegration creates a temporary workspace for integration tests
func NewTestIntegration(t *testing.T) *TestIntegration {
	tmpDir, err := os.MkdirTemp("", "mandor-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Change to temp directory for test
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)

	paths, err := fs.NewPaths()
	if err != nil {
		t.Fatalf("Failed to create paths: %v", err)
	}

	// Initialize workspace using workspace service
	wsSvc, err := service.NewWorkspaceService()
	if err != nil {
		t.Fatalf("Failed to create workspace service: %v", err)
	}

	_, err = wsSvc.InitWorkspace("test-workspace")
	if err != nil {
		t.Fatalf("Failed to initialize workspace: %v", err)
	}

	// Recreate paths after workspace init
	paths, err = fs.NewPaths()
	if err != nil {
		t.Fatalf("Failed to create paths after init: %v", err)
	}

	svc := service.NewBacklogServiceWithPaths(paths)

	t.Cleanup(func() {
		os.Chdir(originalDir)
		os.RemoveAll(tmpDir)
	})

	return &TestIntegration{
		tmpDir: tmpDir,
		paths:  paths,
		svc:    svc,
		t:      t,
	}
}

// ExecuteCommand runs a cobra command and captures output
func ExecuteCommand(cmd *cobra.Command, args ...string) (output string, err error) {
	cmd.SetArgs(args)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	if err := cmd.Execute(); err != nil {
		return buf.String(), err
	}

	return buf.String(), nil
}

// TestCreateBacklogCommand tests the backlog create command
func TestCreateBacklogCommand(t *testing.T) {
	ti := NewTestIntegration(t)

	cmd := NewCreateCmd()

	// Test successful creation
	args := []string{
		"test-backlog",
		"--name", "Test Backlog",
		"--goal", "This backlog represents our effort to deliver a comprehensive and robust system that meets all stakeholder requirements and business objectives. We are committed to ensuring high quality, maintainability, and scalability while continuously improving our processes and practices throughout the project lifecycle. Our focus is on delivering value incrementally and gathering feedback from users to inform our direction. We aim to exceed expectations and create something truly exceptional and meaningful for everyone in the organization.",
		"--yes",
	}

	output, err := ExecuteCommand(cmd, args...)
	if err != nil {
		t.Errorf("CreateCmd failed: %v\nOutput: %s", err, output)
	}

	if !bytes.Contains([]byte(output), []byte("Backlog created: test-backlog")) {
		t.Errorf("Expected success message in output, got: %s", output)
	}

	// Verify backlog was created
	backlog, err := ti.svc.GetBacklog("test-backlog")
	if err != nil {
		t.Errorf("Failed to retrieve created backlog: %v", err)
	}

	if backlog.ID != "test-backlog" {
		t.Errorf("Expected backlog ID 'test-backlog', got '%s'", backlog.ID)
	}

	if backlog.Name != "Test Backlog" {
		t.Errorf("Expected backlog name 'Test Backlog', got '%s'", backlog.Name)
	}
}

// TestCreateBacklogCommand_InvalidID tests creation with invalid backlog ID
func TestCreateBacklogCommand_InvalidID(t *testing.T) {
	ti := NewTestIntegration(t)

	cmd := NewCreateCmd()

	// Test with invalid ID (starts with number)
	args := []string{
		"123-backlog",
		"--name", "Invalid Backlog",
		"--goal", "This backlog represents our effort to deliver a comprehensive and robust system that meets all stakeholder requirements and business objectives. We are committed to ensuring high quality, maintainability, and scalability while continuously improving our processes and practices throughout the project lifecycle. Our focus is on delivering value incrementally and gathering feedback from users to inform our direction. We aim to exceed expectations and create something truly exceptional and meaningful for everyone in the organization.",
		"--yes",
	}

	output, err := ExecuteCommand(cmd, args...)
	if err == nil {
		t.Errorf("Expected error for invalid ID, got: %s", output)
	}

	// Verify backlog was NOT created
	_, err = ti.svc.GetBacklog("123-backlog")
	if err == nil {
		t.Error("Expected error when retrieving invalid backlog, but got none")
	}
}

// TestCreateBacklogCommand_MissingGoal tests creation without goal
func TestCreateBacklogCommand_MissingGoal(t *testing.T) {
	_ = NewTestIntegration(t)

	cmd := NewCreateCmd()

	args := []string{
		"no-goal-backlog",
		"--name", "No Goal Backlog",
		"--yes",
	}

	output, err := ExecuteCommand(cmd, args...)
	if err == nil {
		t.Errorf("Expected error for missing goal, got: %s", output)
	}

	if !bytes.Contains([]byte(output), []byte("goal is required")) {
		t.Errorf("Expected 'goal is required' in error message, got: %s", output)
	}
}

// TestCreateBacklogCommand_GoalTooShort tests creation with goal below minimum length
func TestCreateBacklogCommand_GoalTooShort(t *testing.T) {
	_ = NewTestIntegration(t)

	cmd := NewCreateCmd()

	args := []string{
		"short-goal-backlog",
		"--name", "Short Goal",
		"--goal", "Too short",
		"--yes",
	}

	output, err := ExecuteCommand(cmd, args...)
	if err == nil {
		t.Errorf("Expected error for short goal, got: %s", output)
	}

	if !bytes.Contains([]byte(output), []byte("at least")) {
		t.Errorf("Expected 'at least' in error message, got: %s", output)
	}
}

// TestDetailBacklogCommand tests the backlog detail command
func TestDetailBacklogCommand(t *testing.T) {
	ti := NewTestIntegration(t)

	// First create a backlog
	input := &domain.BacklogCreateInput{
		ID:         "detail-test",
		Name:       "Detail Test Backlog",
		Goal:       "This backlog represents our effort to deliver a comprehensive and robust system that meets all stakeholder requirements and business objectives. We are committed to ensuring high quality, maintainability, and scalability while continuously improving our processes and practices throughout the project lifecycle. Our focus is on delivering value incrementally and gathering feedback from users to inform our direction. We aim to exceed expectations and create something truly exceptional and meaningful for everyone in the organization.",
		TaskDep:    "same_backlog_only",
		FeatureDep: "cross_backlog_allowed",
		IssueDep:   "disabled",
		Strict:     false,
	}
	if err := ti.svc.ValidateCreateInput(input); err != nil {
		t.Fatalf("Failed to validate input: %v", err)
	}
	if err := ti.svc.CreateBacklog(input); err != nil {
		t.Fatalf("Failed to create backlog: %v", err)
	}

	cmd := NewDetailCmd()
	args := []string{"detail-test"}

	output, err := ExecuteCommand(cmd, args...)
	if err != nil {
		t.Errorf("DetailCmd failed: %v\nOutput: %s", err, output)
	}

	if !bytes.Contains([]byte(output), []byte("Detail Test Backlog")) {
		t.Errorf("Expected backlog name in output, got: %s", output)
	}

	if !bytes.Contains([]byte(output), []byte("detail-test")) {
		t.Errorf("Expected backlog ID in output, got: %s", output)
	}
}

// TestDetailBacklogCommand_NotFound tests detail command with non-existent backlog
func TestDetailBacklogCommand_NotFound(t *testing.T) {
	_ = NewTestIntegration(t)

	cmd := NewDetailCmd()
	args := []string{"non-existent"}

	output, err := ExecuteCommand(cmd, args...)
	if err == nil {
		t.Errorf("Expected error for non-existent backlog, got: %s", output)
	}
}

// TestUpdateBacklogCommand tests the backlog update command
func TestUpdateBacklogCommand(t *testing.T) {
	ti := NewTestIntegration(t)

	// First create a backlog
	input := &domain.BacklogCreateInput{
		ID:         "update-test",
		Name:       "Original Name",
		Goal:       "This backlog represents our effort to deliver a comprehensive and robust system that meets all stakeholder requirements and business objectives. We are committed to ensuring high quality, maintainability, and scalability while continuously improving our processes and practices throughout the project lifecycle. Our focus is on delivering value incrementally and gathering feedback from users to inform our direction. We aim to exceed expectations and create something truly exceptional and meaningful for everyone in the organization.",
		TaskDep:    "same_backlog_only",
		FeatureDep: "cross_backlog_allowed",
		IssueDep:   "disabled",
		Strict:     false,
	}
	if err := ti.svc.ValidateCreateInput(input); err != nil {
		t.Fatalf("Failed to validate input: %v", err)
	}
	if err := ti.svc.CreateBacklog(input); err != nil {
		t.Fatalf("Failed to create backlog: %v", err)
	}

	cmd := NewUpdateCmd()
	args := []string{
		"update-test",
		"--name", "Updated Name",
	}

	output, err := ExecuteCommand(cmd, args...)
	if err != nil {
		t.Errorf("UpdateCmd failed: %v\nOutput: %s", err, output)
	}

	if !bytes.Contains([]byte(output), []byte("Updated Name")) {
		t.Errorf("Expected updated name in output, got: %s", output)
	}

	// Verify update was applied
	detail, err := ti.svc.GetBacklogDetail("update-test")
	if err != nil {
		t.Errorf("Failed to retrieve updated backlog: %v", err)
	}

	if detail.Name != "Updated Name" {
		t.Errorf("Expected name 'Updated Name', got '%s'", detail.Name)
	}
}

// TestDeleteBacklogCommand tests the backlog delete command
func TestDeleteBacklogCommand(t *testing.T) {
	ti := NewTestIntegration(t)

	// First create a backlog
	input := &domain.BacklogCreateInput{
		ID:         "delete-test",
		Name:       "Delete Test",
		Goal:       "This backlog represents our effort to deliver a comprehensive and robust system that meets all stakeholder requirements and business objectives. We are committed to ensuring high quality, maintainability, and scalability while continuously improving our processes and practices throughout the project lifecycle. Our focus is on delivering value incrementally and gathering feedback from users to inform our direction. We aim to exceed expectations and create something truly exceptional and meaningful for everyone in the organization.",
		TaskDep:    "same_backlog_only",
		FeatureDep: "cross_backlog_allowed",
		IssueDep:   "disabled",
		Strict:     false,
	}
	if err := ti.svc.ValidateCreateInput(input); err != nil {
		t.Fatalf("Failed to validate input: %v", err)
	}
	if err := ti.svc.CreateBacklog(input); err != nil {
		t.Fatalf("Failed to create backlog: %v", err)
	}

	cmd := NewDeleteCmd()
	args := []string{"delete-test", "--yes"}

	output, err := ExecuteCommand(cmd, args...)
	if err != nil {
		t.Errorf("DeleteCmd failed: %v\nOutput: %s", err, output)
	}

	if !bytes.Contains([]byte(output), []byte("deleted")) {
		t.Errorf("Expected 'deleted' in output, got: %s", output)
	}

	// Verify backlog was soft-deleted (exists but marked as deleted)
	backlog, err := ti.svc.GetBacklog("delete-test")
	if err != nil {
		t.Errorf("Failed to retrieve deleted backlog: %v", err)
	}

	if backlog.Status != domain.BacklogStatusDeleted {
		t.Errorf("Expected status 'deleted', got '%s'", backlog.Status)
	}
}

// TestReopenBacklogCommand tests the backlog reopen command
func TestReopenBacklogCommand(t *testing.T) {
	ti := NewTestIntegration(t)

	// First create and delete a backlog
	input := &domain.BacklogCreateInput{
		ID:         "reopen-test",
		Name:       "Reopen Test",
		Goal:       "This backlog represents our effort to deliver a comprehensive and robust system that meets all stakeholder requirements and business objectives. We are committed to ensuring high quality, maintainability, and scalability while continuously improving our processes and practices throughout the project lifecycle. Our focus is on delivering value incrementally and gathering feedback from users to inform our direction. We aim to exceed expectations and create something truly exceptional and meaningful for everyone in the organization.",
		TaskDep:    "same_backlog_only",
		FeatureDep: "cross_backlog_allowed",
		IssueDep:   "disabled",
		Strict:     false,
	}
	if err := ti.svc.ValidateCreateInput(input); err != nil {
		t.Fatalf("Failed to validate input: %v", err)
	}
	if err := ti.svc.CreateBacklog(input); err != nil {
		t.Fatalf("Failed to create backlog: %v", err)
	}

	deleteInput := &domain.BacklogDeleteInput{
		ID: "reopen-test",
	}
	if _, err := ti.svc.DeleteBacklog(deleteInput); err != nil {
		t.Fatalf("Failed to delete backlog: %v", err)
	}

	cmd := NewReopenCmd()
	args := []string{"reopen-test", "--yes"}

	output, err := ExecuteCommand(cmd, args...)
	if err != nil {
		t.Errorf("ReopenCmd failed: %v\nOutput: %s", err, output)
	}

	if !bytes.Contains([]byte(output), []byte("reopened")) {
		t.Errorf("Expected 'reopened' in output, got: %s", output)
	}

	// Verify backlog was reopened
	backlog, err := ti.svc.GetBacklog("reopen-test")
	if err != nil {
		t.Errorf("Failed to retrieve reopened backlog: %v", err)
	}

	if backlog.Status == domain.BacklogStatusDeleted {
		t.Errorf("Expected status not 'deleted', got '%s'", backlog.Status)
	}
}

// TestMultipleBacklogCreation tests creating multiple backlogs
func TestMultipleBacklogCreation(t *testing.T) {
	ti := NewTestIntegration(t)

	// Create multiple backlogs
	for i := 1; i <= 3; i++ {
		name := "Backlog" + string(rune(48+i))
		input := &domain.BacklogCreateInput{
			ID:         "backlog-" + string(rune(48+i)),
			Name:       name,
			Goal:       "This backlog represents our effort to deliver a comprehensive and robust system that meets all stakeholder requirements and business objectives. We are committed to ensuring high quality, maintainability, and scalability while continuously improving our processes and practices throughout the project lifecycle. Our focus is on delivering value incrementally and gathering feedback from users to inform our direction. We aim to exceed expectations and create something truly exceptional and meaningful for everyone in the organization.",
			TaskDep:    "same_backlog_only",
			FeatureDep: "cross_backlog_allowed",
			IssueDep:   "disabled",
			Strict:     false,
		}
		if err := ti.svc.ValidateCreateInput(input); err != nil {
			t.Fatalf("Failed to validate input: %v", err)
		}
		if err := ti.svc.CreateBacklog(input); err != nil {
			t.Fatalf("Failed to create backlog: %v", err)
		}
	}

	// Verify all backlogs were created by checking each one
	output, err := ti.svc.ListBacklogs(false, false)
	if err != nil {
		t.Fatalf("Failed to list backlogs: %v", err)
	}

	if output.Total != 3 {
		t.Errorf("Expected 3 backlogs, got %d", output.Total)
	}

	backlogIDs := make(map[string]bool)
	for _, item := range output.Backlogs {
		backlogIDs[item.ID] = true
	}

	for i := 1; i <= 3; i++ {
		expected := "backlog-" + string(rune(48+i))
		if !backlogIDs[expected] {
			t.Errorf("Expected backlog '%s' not found in list", expected)
		}
	}
}

// TestBacklogCommandIntegration tests the full workflow
func TestBacklogCommandIntegration(t *testing.T) {
	ti := NewTestIntegration(t)

	// Create a backlog
	createCmd := NewCreateCmd()
	_, err := ExecuteCommand(createCmd,
		"workflow-test",
		"--name", "Workflow Test",
		"--goal", "This backlog represents our effort to deliver a comprehensive and robust system that meets all stakeholder requirements and business objectives. We are committed to ensuring high quality, maintainability, and scalability while continuously improving our processes and practices throughout the project lifecycle. Our focus is on delivering value incrementally and gathering feedback from users to inform our direction. We aim to exceed expectations and create something truly exceptional and meaningful for everyone in the organization.",
		"--yes",
	)
	if err != nil {
		t.Fatalf("Failed to create backlog: %v", err)
	}

	// Get backlog details
	detailCmd := NewDetailCmd()
	output, err := ExecuteCommand(detailCmd, "workflow-test")
	if err != nil {
		t.Fatalf("Failed to get backlog details: %v", err)
	}

	if !bytes.Contains([]byte(output), []byte("workflow-test")) {
		t.Fatalf("Backlog not found in details output")
	}

	// Update backlog
	updateCmd := NewUpdateCmd()
	_, err = ExecuteCommand(updateCmd,
		"workflow-test",
		"--name", "Updated Workflow Test",
	)
	if err != nil {
		t.Fatalf("Failed to update backlog: %v", err)
	}

	// Verify update
	backlog, err := ti.svc.GetBacklog("workflow-test")
	if err != nil {
		t.Fatalf("Failed to retrieve backlog: %v", err)
	}

	if backlog.Name != "Updated Workflow Test" {
		t.Fatalf("Name not updated correctly: %s", backlog.Name)
	}

	// Verify we can list the backlog via service
	listOutput, err := ti.svc.ListBacklogs(false, false)
	if err != nil {
		t.Fatalf("Failed to list backlogs: %v", err)
	}

	found := false
	for _, item := range listOutput.Backlogs {
		if item.ID == "workflow-test" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("Backlog 'workflow-test' not found in list")
	}
}
