package issue_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"mandor/internal/cmd/issue"
)

func TestBlockedCmd_NoWorkspace(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mandor-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	cmd := issue.NewBlockedCmd()
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

func TestBlockedCmd_NoBlockedIssues(t *testing.T) {
	tmpDir := setupTestIssueWorkspace(t)
	defer os.RemoveAll(tmpDir)

	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	projectID := "test-proj"
	writeTestProjectForIssueCmd(t, tmpDir, projectID)

	cmd := issue.NewBlockedCmd()
	cmd.Flags().Set("project", projectID)
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err := cmd.RunE(cmd, []string{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "No blocked issues found") {
		t.Errorf("Expected 'No blocked issues found' message, got: %s", output)
	}
}

func TestBlockedCmd_WithBlockedIssues(t *testing.T) {
	tmpDir := setupTestIssueWorkspace(t)
	defer os.RemoveAll(tmpDir)

	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	projectID := "test-proj"
	writeTestProjectForIssueCmd(t, tmpDir, projectID)

	// Create a blocked issue
	issueID := projectID + "-issue-abc123def456"
	issueData := map[string]interface{}{
		"id":              issueID,
		"project_id":      projectID,
		"name":            "Blocked Issue",
		"issue_type":      "bug",
		"status":          "blocked",
		"priority":        "P2",
		"created_at":      time.Now().UTC().Format(time.RFC3339),
		"last_updated_at": time.Now().UTC().Format(time.RFC3339),
		"created_by":      "testuser",
		"last_updated_by": "testuser",
	}

	issuesPath := filepath.Join(tmpDir, ".mandor", "projects", projectID, "issues.jsonl")
	issuesFile, err := os.Create(issuesPath)
	if err != nil {
		t.Fatalf("Failed to create issues.jsonl: %v", err)
	}
	defer issuesFile.Close()

	jsonBytes, _ := json.Marshal(issueData)
	issuesFile.WriteString(string(jsonBytes) + "\n")

	cmd := issue.NewBlockedCmd()
	cmd.Flags().Set("project", projectID)
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err = cmd.RunE(cmd, []string{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Blocked issues") {
		t.Errorf("Expected 'Blocked issues' in output, got: %s", output)
	}
	if !strings.Contains(output, issueID) {
		t.Errorf("Expected issue ID in output, got: %s", output)
	}
	if !strings.Contains(output, "Blocked Issue") {
		t.Errorf("Expected issue name in output, got: %s", output)
	}
}

func TestBlockedCmd_WithTypeFilter(t *testing.T) {
	tmpDir := setupTestIssueWorkspace(t)
	defer os.RemoveAll(tmpDir)

	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	projectID := "test-proj"
	writeTestProjectForIssueCmd(t, tmpDir, projectID)

	cmd := issue.NewBlockedCmd()
	cmd.Flags().Set("project", projectID)
	cmd.Flags().Set("type", "bug")
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err := cmd.RunE(cmd, []string{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "No blocked issues found") {
		t.Errorf("Expected 'No blocked issues found' message with type filter, got: %s", output)
	}
}

func TestBlockedCmd_WithPriorityFilter(t *testing.T) {
	tmpDir := setupTestIssueWorkspace(t)
	defer os.RemoveAll(tmpDir)

	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	projectID := "test-proj"
	writeTestProjectForIssueCmd(t, tmpDir, projectID)

	cmd := issue.NewBlockedCmd()
	cmd.Flags().Set("project", projectID)
	cmd.Flags().Set("priority", "P0")
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err := cmd.RunE(cmd, []string{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "No blocked issues found") {
		t.Errorf("Expected 'No blocked issues found' message with priority filter, got: %s", output)
	}
}

func TestBlockedCmd_JSONOutput(t *testing.T) {
	tmpDir := setupTestIssueWorkspace(t)
	defer os.RemoveAll(tmpDir)

	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	projectID := "test-proj"
	writeTestProjectForIssueCmd(t, tmpDir, projectID)

	issueID := projectID + "-issue-xyz789uvw123"
	issueData := map[string]interface{}{
		"id":              issueID,
		"project_id":      projectID,
		"name":            "Blocked Issue",
		"issue_type":      "improvement",
		"status":          "blocked",
		"priority":        "P1",
		"created_at":      time.Now().UTC().Format(time.RFC3339),
		"last_updated_at": time.Now().UTC().Format(time.RFC3339),
		"created_by":      "testuser",
		"last_updated_by": "testuser",
	}

	issuesPath := filepath.Join(tmpDir, ".mandor", "projects", projectID, "issues.jsonl")
	issuesFile, err := os.Create(issuesPath)
	if err != nil {
		t.Fatalf("Failed to create issues.jsonl: %v", err)
	}
	defer issuesFile.Close()

	jsonBytes, _ := json.Marshal(issueData)
	issuesFile.WriteString(string(jsonBytes) + "\n")

	cmd := issue.NewBlockedCmd()
	cmd.Flags().Set("project", projectID)
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

	if _, ok := result["issues"]; !ok {
		t.Error("Expected 'issues' key in JSON output")
	}
	if _, ok := result["total"]; !ok {
		t.Error("Expected 'total' key in JSON output")
	}
}

func TestBlockedCmd_InvalidType(t *testing.T) {
	tmpDir := setupTestIssueWorkspace(t)
	defer os.RemoveAll(tmpDir)

	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	projectID := "test-proj"
	writeTestProjectForIssueCmd(t, tmpDir, projectID)

	cmd := issue.NewBlockedCmd()
	cmd.Flags().Set("project", projectID)
	cmd.Flags().Set("type", "invalid_type")
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err := cmd.RunE(cmd, []string{})
	if err == nil {
		t.Error("Expected error for invalid issue type, got nil")
	}

	if !strings.Contains(err.Error(), "Invalid issue type") {
		t.Errorf("Expected 'Invalid issue type' error, got: %v", err)
	}
}

func TestBlockedCmd_InvalidPriority(t *testing.T) {
	tmpDir := setupTestIssueWorkspace(t)
	defer os.RemoveAll(tmpDir)

	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	projectID := "test-proj"
	writeTestProjectForIssueCmd(t, tmpDir, projectID)

	cmd := issue.NewBlockedCmd()
	cmd.Flags().Set("project", projectID)
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

func TestBlockedCmd_NoDefaultProjectAndNoFilter(t *testing.T) {
	tmpDir := setupTestIssueWorkspace(t)
	defer os.RemoveAll(tmpDir)

	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	cmd := issue.NewBlockedCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err := cmd.RunE(cmd, []string{})
	if err == nil {
		t.Error("Expected error when no project specified and no default set, got nil")
	}

	if !strings.Contains(err.Error(), "No project specified") {
		t.Errorf("Expected 'No project specified' error, got: %v", err)
	}
}
