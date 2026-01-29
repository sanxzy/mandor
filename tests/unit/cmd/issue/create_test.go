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
	"mandor/internal/domain"
)

func setupTestIssueWorkspace(t *testing.T) string {
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

func writeTestProjectForIssueCmd(t *testing.T, tmpDir, projectID string) {
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

	issuesPath := filepath.Join(projectDir, "issues.jsonl")
	if err := os.WriteFile(issuesPath, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create issues.jsonl: %v", err)
	}

	eventsPath := filepath.Join(projectDir, "events.jsonl")
	if err := os.WriteFile(eventsPath, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create events.jsonl: %v", err)
	}
}

func TestCreateCmd_Validation(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		flags       map[string]string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "no project flag",
			args:        []string{"Test Issue"},
			flags:       map[string]string{"type": "bug", "goal": "Test goal", "affected-files": "src/file.ts", "affected-tests": "tests/test.ts", "implementation-steps": "Step 1"},
			expectError: true,
			errorMsg:    "Project ID is required",
		},
		{
			name:        "no type flag",
			args:        []string{"Test Issue"},
			flags:       map[string]string{"project": "auth", "goal": "Test goal", "affected-files": "src/file.ts", "affected-tests": "tests/test.ts", "implementation-steps": "Step 1"},
			expectError: true,
			errorMsg:    "Issue type is required",
		},
		{
			name:        "no goal flag",
			args:        []string{"Test Issue"},
			flags:       map[string]string{"project": "auth", "type": "bug", "affected-files": "src/file.ts", "affected-tests": "tests/test.ts", "implementation-steps": "Step 1"},
			expectError: true,
			errorMsg:    "Issue goal is required",
		},
		{
			name:        "no affected files",
			args:        []string{"Test Issue"},
			flags:       map[string]string{"project": "auth", "type": "bug", "goal": "Test goal", "affected-tests": "tests/test.ts", "implementation-steps": "Step 1"},
			expectError: true,
			errorMsg:    "Affected files are required",
		},
		{
			name:        "no affected tests",
			args:        []string{"Test Issue"},
			flags:       map[string]string{"project": "auth", "type": "bug", "goal": "Test goal", "affected-files": "src/file.ts", "implementation-steps": "Step 1"},
			expectError: true,
			errorMsg:    "Affected tests are required",
		},
		{
			name:        "no implementation steps",
			args:        []string{"Test Issue"},
			flags:       map[string]string{"project": "auth", "type": "bug", "goal": "Test goal", "affected-files": "src/file.ts", "affected-tests": "tests/test.ts"},
			expectError: true,
			errorMsg:    "Implementation steps are required",
		},
		{
			name:        "invalid type",
			args:        []string{"Test Issue"},
			flags:       map[string]string{"project": "auth", "type": "invalid", "goal": "Test goal", "affected-files": "src/file.ts", "affected-tests": "tests/test.ts", "implementation-steps": "Step 1"},
			expectError: true,
			errorMsg:    "Invalid issue type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := setupTestIssueWorkspace(t)
			defer os.RemoveAll(tmpDir)

			os.Chdir(tmpDir)

			cmd := issue.NewCreateCmd()
			cmd.SetArgs(append(tt.args, flattenFlags(tt.flags)...))

			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)

			err := cmd.Execute()
			if tt.expectError {
				if err == nil {
					output := buf.String()
					if !strings.Contains(output, tt.errorMsg) {
						t.Errorf("Expected error containing '%s' but got: %s", tt.errorMsg, output)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestCreateCmd_Success(t *testing.T) {
	tmpDir := setupTestIssueWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)
	writeTestProjectForIssueCmd(t, tmpDir, "auth")

	cmd := issue.NewCreateCmd()
	cmd.SetArgs([]string{
		"Test Issue",
		"--project", "auth",
		"--type", "bug",
		"--goal", "Test goal description",
		"--affected-files", "src/file1.ts,src/file2.ts",
		"--affected-tests", "tests/file1.test.ts",
		"--implementation-steps", "Step 1,Step 2",
		"--priority", "P1",
	})

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := cmd.Execute()
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Issue created:") {
		t.Errorf("Expected 'Issue created:' in output but got: %s", output)
	}

	if !strings.Contains(output, "auth-issue-") {
		t.Errorf("Expected issue ID with 'auth-issue-' prefix in output but got: %s", output)
	}

	if !strings.Contains(output, "P1") {
		t.Errorf("Expected priority 'P1' in output but got: %s", output)
	}
}

func TestCreateCmd_WithDependency(t *testing.T) {
	tmpDir := setupTestIssueWorkspace(t)
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)
	writeTestProjectForIssueCmd(t, tmpDir, "auth")

	depIssue := &domain.Issue{
		ID:                  "auth-issue-dep123",
		ProjectID:           "auth",
		Name:                "Dependency Issue",
		Goal:                "Test goal",
		IssueType:           "bug",
		Priority:            "P2",
		Status:              domain.IssueStatusReady,
		AffectedFiles:       []string{"src/dep.ts"},
		AffectedTests:       []string{"tests/dep.test.ts"},
		ImplementationSteps: []string{"Step 1"},
		CreatedAt:           time.Now().UTC(),
		LastUpdatedAt:       time.Now().UTC(),
		CreatedBy:           "testuser",
		LastUpdatedBy:       "testuser",
	}

	issuesPath := filepath.Join(tmpDir, ".mandor", "projects", "auth", "issues.jsonl")
	file, err := os.OpenFile(issuesPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("Failed to open issues.jsonl: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(depIssue); err != nil {
		t.Fatalf("Failed to write dependency issue: %v", err)
	}

	cmd := issue.NewCreateCmd()
	cmd.SetArgs([]string{
		"Test Issue",
		"--project", "auth",
		"--type", "security",
		"--goal", "Test goal description",
		"--affected-files", "src/file1.ts",
		"--affected-tests", "tests/file1.test.ts",
		"--implementation-steps", "Step 1",
		"--depends-on", "auth-issue-dep123",
	})

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err = cmd.Execute()
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Issue created:") {
		t.Errorf("Expected 'Issue created:' in output but got: %s", output)
	}

	if !strings.Contains(output, "Depends on") {
		t.Errorf("Expected 'Depends on' in output but got: %s", output)
	}
}

func TestNewIssueCreateCmd(t *testing.T) {
	cmd := issue.NewCreateCmd()
	if cmd == nil {
		t.Fatal("Expected command, got nil")
	}
	if !strings.Contains(cmd.Use, "create") {
		t.Errorf("Expected 'create' in use, got %q", cmd.Use)
	}
}

func flattenFlags(flags map[string]string) []string {
	var result []string
	for k, v := range flags {
		result = append(result, "--"+k, v)
	}
	return result
}
