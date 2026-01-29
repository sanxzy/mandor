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
)

func TestReadyCmd_NoWorkspace(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mandor-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	cmd := task.NewReadyCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err = cmd.RunE(cmd, []string{})
	if err == nil {
		t.Error("Expected error for uninitialized workspace, got nil")
	}

	if !strings.Contains(err.Error(), "Workspace not initialized") {
		t.Errorf("Expected 'Workspace not initialized' error, got: %v", err)
	}
}

func TestReadyCmd_NoReadyTasks(t *testing.T) {
	tmpDir := setupTestTaskWorkspace(t)
	defer os.RemoveAll(tmpDir)

	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	projectID := "test-proj"
	featureID := "test-proj-feature-abc123def456"
	writeTestProjectForTaskCmd(t, tmpDir, projectID)
	writeTestFeatureForTaskCmd(t, tmpDir, projectID, featureID, "active")

	cmd := task.NewReadyCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err := cmd.RunE(cmd, []string{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "No ready tasks found") {
		t.Errorf("Expected 'No ready tasks found' message, got: %s", output)
	}
}

func TestReadyCmd_WithReadyTasks(t *testing.T) {
	tmpDir := setupTestTaskWorkspace(t)
	defer os.RemoveAll(tmpDir)

	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	projectID := "test-proj"
	featureID := "test-proj-feature-abc123def456"
	writeTestProjectForTaskCmd(t, tmpDir, projectID)
	writeTestFeatureForTaskCmd(t, tmpDir, projectID, featureID, "active")

	// Create a ready task
	taskID := featureID + "-task-xyz789uvw123"
	taskData := map[string]interface{}{
		"id":         taskID,
		"feature_id": featureID,
		"project_id": projectID,
		"name":       "Ready Task",
		"status":     "ready",
		"priority":   "P2",
		"created_at": time.Now().UTC().Format(time.RFC3339),
		"updated_at": time.Now().UTC().Format(time.RFC3339),
		"created_by": "testuser",
		"updated_by": "testuser",
	}

	tasksPath := filepath.Join(tmpDir, ".mandor", "projects", projectID, "tasks.jsonl")
	tasksFile, err := os.Create(tasksPath)
	if err != nil {
		t.Fatalf("Failed to create tasks.jsonl: %v", err)
	}
	defer tasksFile.Close()

	jsonBytes, _ := json.Marshal(taskData)
	tasksFile.WriteString(string(jsonBytes) + "\n")

	cmd := task.NewReadyCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err = cmd.RunE(cmd, []string{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Ready tasks") {
		t.Errorf("Expected 'Ready tasks' in output, got: %s", output)
	}
	if !strings.Contains(output, taskID) {
		t.Errorf("Expected task ID in output, got: %s", output)
	}
	if !strings.Contains(output, "Ready Task") {
		t.Errorf("Expected task name in output, got: %s", output)
	}
}

func TestReadyCmd_WithProjectFilter(t *testing.T) {
	tmpDir := setupTestTaskWorkspace(t)
	defer os.RemoveAll(tmpDir)

	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	projectID := "test-proj"
	featureID := "test-proj-feature-abc123def456"
	writeTestProjectForTaskCmd(t, tmpDir, projectID)
	writeTestFeatureForTaskCmd(t, tmpDir, projectID, featureID, "active")

	cmd := task.NewReadyCmd()
	cmd.Flags().Set("project", projectID)
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err := cmd.RunE(cmd, []string{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "No ready tasks found") {
		t.Errorf("Expected 'No ready tasks found' message with project filter, got: %s", output)
	}
}

func TestReadyCmd_WithFeatureFilter(t *testing.T) {
	tmpDir := setupTestTaskWorkspace(t)
	defer os.RemoveAll(tmpDir)

	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	projectID := "test-proj"
	featureID := "test-proj-feature-abc123def456"
	writeTestProjectForTaskCmd(t, tmpDir, projectID)
	writeTestFeatureForTaskCmd(t, tmpDir, projectID, featureID, "active")

	cmd := task.NewReadyCmd()
	cmd.Flags().Set("feature", featureID)
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err := cmd.RunE(cmd, []string{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "No ready tasks found") {
		t.Errorf("Expected 'No ready tasks found' message with feature filter, got: %s", output)
	}
}

func TestReadyCmd_JSONOutput(t *testing.T) {
	tmpDir := setupTestTaskWorkspace(t)
	defer os.RemoveAll(tmpDir)

	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	projectID := "test-proj"
	featureID := "test-proj-feature-abc123def456"
	writeTestProjectForTaskCmd(t, tmpDir, projectID)
	writeTestFeatureForTaskCmd(t, tmpDir, projectID, featureID, "active")

	taskID := featureID + "-task-xyz789uvw123"
	taskData := map[string]interface{}{
		"id":         taskID,
		"feature_id": featureID,
		"project_id": projectID,
		"name":       "Ready Task",
		"status":     "ready",
		"priority":   "P1",
		"created_at": time.Now().UTC().Format(time.RFC3339),
		"updated_at": time.Now().UTC().Format(time.RFC3339),
		"created_by": "testuser",
		"updated_by": "testuser",
	}

	tasksPath := filepath.Join(tmpDir, ".mandor", "projects", projectID, "tasks.jsonl")
	tasksFile, err := os.Create(tasksPath)
	if err != nil {
		t.Fatalf("Failed to create tasks.jsonl: %v", err)
	}
	defer tasksFile.Close()

	jsonBytes, _ := json.Marshal(taskData)
	tasksFile.WriteString(string(jsonBytes) + "\n")

	cmd := task.NewReadyCmd()
	cmd.Flags().Set("json", "true")
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err = cmd.RunE(cmd, []string{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	output := buf.String()
	var result map[string]interface{}
	err = json.Unmarshal([]byte(output), &result)
	if err != nil {
		t.Errorf("Failed to parse JSON output: %v", err)
	}

	if _, ok := result["tasks"]; !ok {
		t.Error("Expected 'tasks' key in JSON output")
	}
	if _, ok := result["total"]; !ok {
		t.Error("Expected 'total' key in JSON output")
	}
}

func TestReadyCmd_InvalidPriority(t *testing.T) {
	tmpDir := setupTestTaskWorkspace(t)
	defer os.RemoveAll(tmpDir)

	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	projectID := "test-proj"
	writeTestProjectForTaskCmd(t, tmpDir, projectID)

	cmd := task.NewReadyCmd()
	cmd.Flags().Set("priority", "INVALID")
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err := cmd.RunE(cmd, []string{})
	if err == nil {
		t.Error("Expected error for invalid priority, got nil")
	}

	if !strings.Contains(err.Error(), "Invalid priority") {
		t.Errorf("Expected 'Invalid priority' error, got: %v", err)
	}
}
