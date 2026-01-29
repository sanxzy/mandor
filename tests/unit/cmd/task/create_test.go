package task_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"mandor/internal/cmd/task"
	"mandor/internal/domain"
)

func setupTestTaskWorkspace(t *testing.T) string {
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

	return tmpDir
}

func writeTestProjectForTaskCmd(t *testing.T, tmpDir, projectID string) {
	t.Helper()

	projectDir := filepath.Join(tmpDir, ".mandor", "projects", projectID)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatalf("Failed to create project dir: %v", err)
	}

	projectPath := filepath.Join(projectDir, "project.jsonl")
	projectData := []byte("{\"id\":\"" + projectID + "\",\"name\":\"Test Project\",\"goal\":\"test goal\",\"status\":\"active\",\"strict\":false}")
	if err := os.WriteFile(projectPath, projectData, 0644); err != nil {
		t.Fatalf("Failed to write project.jsonl: %v", err)
	}

	schemaPath := filepath.Join(projectDir, "schema.json")
	schemaData := []byte("{\"version\":\"mandor.v1\",\"schema\":\"https://json-schema.org/draft/2020-12/schema\",\"rules\":{\"task\":{\"dependency\":\"same_project_only\",\"cycle\":\"disallowed\"},\"feature\":{\"dependency\":\"cross_project_allowed\",\"cycle\":\"disallowed\"},\"issue\":{\"dependency\":\"same_project_only\",\"cycle\":\"disallowed\"},\"priority\":{\"levels\":[\"P0\",\"P1\",\"P2\",\"P3\",\"P4\",\"P5\"],\"default\":\"P3\"}}}")
	if err := os.WriteFile(schemaPath, schemaData, 0644); err != nil {
		t.Fatalf("Failed to write schema.json: %v", err)
	}

	featuresPath := filepath.Join(projectDir, "features.jsonl")
	if err := os.WriteFile(featuresPath, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create features.jsonl: %v", err)
	}

	tasksPath := filepath.Join(projectDir, "tasks.jsonl")
	if err := os.WriteFile(tasksPath, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create tasks.jsonl: %v", err)
	}

	eventsPath := filepath.Join(projectDir, "events.jsonl")
	if err := os.WriteFile(eventsPath, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create events.jsonl: %v", err)
	}
}

func writeTestFeatureForTaskCmd(t *testing.T, tmpDir, projectID, featureID string, status string) {
	t.Helper()

	projectDir := filepath.Join(tmpDir, ".mandor", "projects", projectID)
	featuresPath := filepath.Join(projectDir, "features.jsonl")

	feature := &domain.Feature{
		ID:        featureID,
		ProjectID: projectID,
		Name:      "Test Feature",
		Goal:      "Test goal",
		Scope:     "fullstack",
		Priority:  "P3",
		Status:    status,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		CreatedBy: "testuser",
		UpdatedBy: "testuser",
	}

	file, err := os.OpenFile(featuresPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("Failed to open features.jsonl: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(feature)
	if err != nil {
		t.Fatalf("Failed to write feature: %v", err)
	}
}

func TestTaskCreateCmd_NotInitialized(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mandor-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	cmd := task.NewCreateCmd()
	cmd.SetArgs([]string{"Test Task", "--feature", "testproject-feature-abc", "--goal", "Test goal", "--implementation-steps", "step1,step2", "--test-cases", "test1,test2", "--derivable-files", "file1,file2", "--library-needs", "lib1,lib2"})

	err = cmd.Execute()
	if err == nil {
		t.Error("Expected error for non-initialized workspace")
	}

	me, ok := err.(*domain.MandorError)
	if !ok {
		t.Fatal("Expected MandorError")
	}

	if me.Code != domain.ExitValidationError {
		t.Errorf("Expected exit code %d, got %d", domain.ExitValidationError, me.Code)
	}
}

func TestTaskCreateCmd_MissingFeature(t *testing.T) {
	tmpDir := setupTestTaskWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	cmd := task.NewCreateCmd()
	cmd.SetArgs([]string{"Test Task", "--goal", "Test goal", "--implementation-steps", "step1,step2", "--test-cases", "test1,test2", "--derivable-files", "file1,file2", "--library-needs", "lib1,lib2"})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for missing feature")
	}
}

func TestTaskCreateCmd_MissingGoal(t *testing.T) {
	tmpDir := setupTestTaskWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	cmd := task.NewCreateCmd()
	cmd.SetArgs([]string{"Test Task", "--feature", "testproject-feature-abc", "--implementation-steps", "step1,step2", "--test-cases", "test1,test2", "--derivable-files", "file1,file2", "--library-needs", "lib1,lib2"})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for missing goal")
	}
}

func TestTaskCreateCmd_MissingImplSteps(t *testing.T) {
	tmpDir := setupTestTaskWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	cmd := task.NewCreateCmd()
	cmd.SetArgs([]string{"Test Task", "--feature", "testproject-feature-abc", "--goal", "Test goal", "--test-cases", "test1,test2", "--derivable-files", "file1,file2", "--library-needs", "lib1,lib2"})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for missing implementation steps")
	}
}

func TestTaskCreateCmd_InvalidPriority(t *testing.T) {
	tmpDir := setupTestTaskWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)
	writeTestProjectForTaskCmd(t, tmpDir, "testproject")
	writeTestFeatureForTaskCmd(t, tmpDir, "testproject", "testproject-feature-abc", domain.FeatureStatusActive)

	cmd := task.NewCreateCmd()
	cmd.SetArgs([]string{"Test Task", "--feature", "testproject-feature-abc", "--goal", "Test goal", "--implementation-steps", "step1", "--test-cases", "test1", "--derivable-files", "file1", "--library-needs", "lib1", "--priority", "P6"})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for invalid priority")
	}
}

func TestTaskCreateCmd_Success(t *testing.T) {
	tmpDir := setupTestTaskWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)
	writeTestProjectForTaskCmd(t, tmpDir, "testproject")
	writeTestFeatureForTaskCmd(t, tmpDir, "testproject", "testproject-feature-abc", domain.FeatureStatusActive)

	cmd := task.NewCreateCmd()
	cmd.SetArgs([]string{
		"New Task",
		"--feature", "testproject-feature-abc",
		"--goal", "This is a test task goal",
		"--implementation-steps", "step1,step2",
		"--test-cases", "test1,test2",
		"--derivable-files", "file1,file2",
		"--library-needs", "lib1,lib2",
	})

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	output := buf.String()
	if len(output) == 0 {
		t.Error("Expected output, got empty")
	}

	if !strings.Contains(output, "Task created:") {
		t.Error("Expected 'Task created:' in output")
	}
}

func TestNewTaskCreateCmd(t *testing.T) {
	cmd := task.NewCreateCmd()
	if cmd == nil {
		t.Fatal("Expected command, got nil")
	}
	if !strings.Contains(cmd.Use, "create") {
		t.Errorf("Expected 'create' in use, got %q", cmd.Use)
	}
}
