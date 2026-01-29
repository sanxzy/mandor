package service_test

import (
	"os"
	"path/filepath"
	"testing"

	"mandor/internal/domain"
	"mandor/internal/fs"
	"mandor/internal/service"
)

func TestIssueService_ValidateCreateInput(t *testing.T) {
	paths, err := fs.NewPathsFromRoot("/tmp/mandor-test-" + randomString(8))
	if err != nil {
		t.Fatalf("Failed to create paths: %v", err)
	}

	ws := &domain.Workspace{
		ID:            "test-workspace",
		Name:          "Test Workspace",
		Version:       "1.0.0",
		SchemaVersion: "1.0.0",
		Config: domain.WorkspaceConfig{
			DefaultPriority: "P3",
			StrictMode:      false,
			DefaultProject:  "auth",
		},
	}

	writer := fs.NewWriter(paths)
	if err := writer.CreateMandorDir(); err != nil {
		t.Fatalf("Failed to create mandor dir: %v", err)
	}
	if err := writer.WriteWorkspace(ws); err != nil {
		t.Fatalf("Failed to write workspace: %v", err)
	}

	if err := writer.CreateProjectDir("auth"); err != nil {
		t.Fatalf("Failed to create project dir: %v", err)
	}

	svc := service.NewIssueServiceWithPaths(paths)

	tests := []struct {
		name        string
		input       *domain.IssueCreateInput
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid input",
			input: &domain.IssueCreateInput{
				ProjectID:           "auth",
				Name:                "Test Issue",
				Goal:                "Test goal description",
				IssueType:           "bug",
				Priority:            "P2",
				AffectedFiles:       []string{"src/file1.ts", "src/file2.ts"},
				AffectedTests:       []string{"tests/file1.test.ts"},
				ImplementationSteps: []string{"Step 1", "Step 2"},
				LibraryNeeds:        []string{},
			},
			expectError: false,
		},
		{
			name: "project not found",
			input: &domain.IssueCreateInput{
				ProjectID:           "nonexistent",
				Name:                "Test Issue",
				Goal:                "Test goal",
				IssueType:           "bug",
				AffectedFiles:       []string{"src/file1.ts"},
				AffectedTests:       []string{"tests/file1.test.ts"},
				ImplementationSteps: []string{"Step 1"},
			},
			expectError: true,
			errorMsg:    "Project not found: nonexistent",
		},
		{
			name: "missing name",
			input: &domain.IssueCreateInput{
				ProjectID:           "auth",
				Name:                "",
				Goal:                "Test goal",
				IssueType:           "bug",
				AffectedFiles:       []string{"src/file1.ts"},
				AffectedTests:       []string{"tests/file1.test.ts"},
				ImplementationSteps: []string{"Step 1"},
			},
			expectError: true,
			errorMsg:    "Issue name is required",
		},
		{
			name: "missing goal",
			input: &domain.IssueCreateInput{
				ProjectID:           "auth",
				Name:                "Test Issue",
				Goal:                "",
				IssueType:           "bug",
				AffectedFiles:       []string{"src/file1.ts"},
				AffectedTests:       []string{"tests/file1.test.ts"},
				ImplementationSteps: []string{"Step 1"},
			},
			expectError: true,
			errorMsg:    "Issue goal is required",
		},
		{
			name: "missing type",
			input: &domain.IssueCreateInput{
				ProjectID:           "auth",
				Name:                "Test Issue",
				Goal:                "Test goal",
				IssueType:           "",
				AffectedFiles:       []string{"src/file1.ts"},
				AffectedTests:       []string{"tests/file1.test.ts"},
				ImplementationSteps: []string{"Step 1"},
			},
			expectError: true,
			errorMsg:    "Issue type is required",
		},
		{
			name: "invalid type",
			input: &domain.IssueCreateInput{
				ProjectID:           "auth",
				Name:                "Test Issue",
				Goal:                "Test goal",
				IssueType:           "invalid",
				AffectedFiles:       []string{"src/file1.ts"},
				AffectedTests:       []string{"tests/file1.test.ts"},
				ImplementationSteps: []string{"Step 1"},
			},
			expectError: true,
			errorMsg:    "Invalid issue type",
		},
		{
			name: "missing affected files",
			input: &domain.IssueCreateInput{
				ProjectID:           "auth",
				Name:                "Test Issue",
				Goal:                "Test goal",
				IssueType:           "bug",
				AffectedFiles:       []string{},
				AffectedTests:       []string{"tests/file1.test.ts"},
				ImplementationSteps: []string{"Step 1"},
			},
			expectError: true,
			errorMsg:    "Affected files are required",
		},
		{
			name: "missing affected tests",
			input: &domain.IssueCreateInput{
				ProjectID:           "auth",
				Name:                "Test Issue",
				Goal:                "Test goal",
				IssueType:           "bug",
				AffectedFiles:       []string{"src/file1.ts"},
				AffectedTests:       []string{},
				ImplementationSteps: []string{"Step 1"},
			},
			expectError: true,
			errorMsg:    "Affected tests are required",
		},
		{
			name: "missing implementation steps",
			input: &domain.IssueCreateInput{
				ProjectID:           "auth",
				Name:                "Test Issue",
				Goal:                "Test goal",
				IssueType:           "bug",
				AffectedFiles:       []string{"src/file1.ts"},
				AffectedTests:       []string{"tests/file1.test.ts"},
				ImplementationSteps: []string{},
			},
			expectError: true,
			errorMsg:    "Implementation steps are required",
		},
		{
			name: "invalid priority",
			input: &domain.IssueCreateInput{
				ProjectID:           "auth",
				Name:                "Test Issue",
				Goal:                "Test goal",
				IssueType:           "bug",
				Priority:            "P10",
				AffectedFiles:       []string{"src/file1.ts"},
				AffectedTests:       []string{"tests/file1.test.ts"},
				ImplementationSteps: []string{"Step 1"},
			},
			expectError: true,
			errorMsg:    "Invalid priority",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.ValidateCreateInput(tt.input)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error '%s' but got nil", tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}

	os.RemoveAll(filepath.Dir(paths.MandorDirPath()))
}

func TestIssueService_CreateIssue(t *testing.T) {
	paths, err := fs.NewPathsFromRoot("/tmp/mandor-test-" + randomString(8))
	if err != nil {
		t.Fatalf("Failed to create paths: %v", err)
	}

	ws := &domain.Workspace{
		ID:            "test-workspace",
		Name:          "Test Workspace",
		Version:       "1.0.0",
		SchemaVersion: "1.0.0",
		Config: domain.WorkspaceConfig{
			DefaultPriority: "P3",
			StrictMode:      false,
		},
	}

	writer := fs.NewWriter(paths)
	if err := writer.CreateMandorDir(); err != nil {
		t.Fatalf("Failed to create mandor dir: %v", err)
	}
	if err := writer.WriteWorkspace(ws); err != nil {
		t.Fatalf("Failed to write workspace: %v", err)
	}

	if err := writer.CreateProjectDir("auth"); err != nil {
		t.Fatalf("Failed to create project dir: %v", err)
	}

	svc := service.NewIssueServiceWithPaths(paths)

	input := &domain.IssueCreateInput{
		ProjectID:           "auth",
		Name:                "Test Issue",
		Goal:                "Test goal description",
		IssueType:           "bug",
		Priority:            "P2",
		AffectedFiles:       []string{"src/file1.ts", "src/file2.ts"},
		AffectedTests:       []string{"tests/file1.test.ts"},
		ImplementationSteps: []string{"Step 1", "Step 2"},
		LibraryNeeds:        []string{"lib1", "lib2"},
	}

	issue, err := svc.CreateIssue(input)
	if err != nil {
		t.Fatalf("Failed to create issue: %v", err)
	}

	if issue.Status != domain.IssueStatusReady {
		t.Errorf("Expected status 'ready' but got '%s'", issue.Status)
	}

	if issue.ProjectID != "auth" {
		t.Errorf("Expected project 'auth' but got '%s'", issue.ProjectID)
	}

	if len(issue.DependsOn) != 0 {
		t.Errorf("Expected no dependencies but got %d", len(issue.DependsOn))
	}

	os.RemoveAll(filepath.Dir(paths.MandorDirPath()))
}

func TestIssueService_CreateIssueWithDeps(t *testing.T) {
	paths, err := fs.NewPathsFromRoot("/tmp/mandor-test-" + randomString(8))
	if err != nil {
		t.Fatalf("Failed to create paths: %v", err)
	}

	ws := &domain.Workspace{
		ID:            "test-workspace",
		Name:          "Test Workspace",
		Version:       "1.0.0",
		SchemaVersion: "1.0.0",
		Config: domain.WorkspaceConfig{
			DefaultPriority: "P3",
			StrictMode:      false,
		},
	}

	writer := fs.NewWriter(paths)
	if err := writer.CreateMandorDir(); err != nil {
		t.Fatalf("Failed to create mandor dir: %v", err)
	}
	if err := writer.WriteWorkspace(ws); err != nil {
		t.Fatalf("Failed to write workspace: %v", err)
	}

	if err := writer.CreateProjectDir("auth"); err != nil {
		t.Fatalf("Failed to create project dir: %v", err)
	}

	svc := service.NewIssueServiceWithPaths(paths)

	depIssue, err := svc.CreateIssue(&domain.IssueCreateInput{
		ProjectID:           "auth",
		Name:                "Dependency Issue",
		Goal:                "Test goal",
		IssueType:           "bug",
		AffectedFiles:       []string{"src/file1.ts"},
		AffectedTests:       []string{"tests/file1.test.ts"},
		ImplementationSteps: []string{"Step 1"},
	})
	if err != nil {
		t.Fatalf("Failed to create dependency issue: %v", err)
	}

	issue, err := svc.CreateIssue(&domain.IssueCreateInput{
		ProjectID:           "auth",
		Name:                "Test Issue",
		Goal:                "Test goal",
		IssueType:           "bug",
		AffectedFiles:       []string{"src/file2.ts"},
		AffectedTests:       []string{"tests/file2.test.ts"},
		ImplementationSteps: []string{"Step 1"},
		DependsOn:           []string{depIssue.ID},
	})
	if err != nil {
		t.Fatalf("Failed to create issue with deps: %v", err)
	}

	if issue.Status != domain.IssueStatusBlocked {
		t.Errorf("Expected status 'blocked' (has unresolved dependency) but got '%s'", issue.Status)
	}

	os.RemoveAll(filepath.Dir(paths.MandorDirPath()))
}

func TestIssueService_ValidateDependency(t *testing.T) {
	paths, err := fs.NewPathsFromRoot("/tmp/mandor-test-" + randomString(8))
	if err != nil {
		t.Fatalf("Failed to create paths: %v", err)
	}

	ws := &domain.Workspace{
		ID:            "test-workspace",
		Name:          "Test Workspace",
		Version:       "1.0.0",
		SchemaVersion: "1.0.0",
		Config: domain.WorkspaceConfig{
			DefaultPriority: "P3",
			StrictMode:      false,
		},
	}

	writer := fs.NewWriter(paths)
	if err := writer.CreateMandorDir(); err != nil {
		t.Fatalf("Failed to create mandor dir: %v", err)
	}
	if err := writer.WriteWorkspace(ws); err != nil {
		t.Fatalf("Failed to write workspace: %v", err)
	}

	if err := writer.CreateProjectDir("auth"); err != nil {
		t.Fatalf("Failed to create project dir: %v", err)
	}

	svc := service.NewIssueServiceWithPaths(paths)

	issue1, err := svc.CreateIssue(&domain.IssueCreateInput{
		ProjectID:           "auth",
		Name:                "Issue 1",
		Goal:                "Test goal",
		IssueType:           "bug",
		AffectedFiles:       []string{"src/file1.ts"},
		AffectedTests:       []string{"tests/file1.test.ts"},
		ImplementationSteps: []string{"Step 1"},
	})
	if err != nil {
		t.Fatalf("Failed to create issue 1: %v", err)
	}

	tests := []struct {
		name        string
		depID       string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "self dependency",
			depID:       issue1.ID,
			expectError: true,
			errorMsg:    "Self-dependency detected",
		},
		{
			name:        "nonexistent dependency",
			depID:       "auth-issue-nonexistent",
			expectError: true,
			errorMsg:    "Dependency not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &domain.IssueCreateInput{
				ProjectID:           "auth",
				Name:                "Test Issue",
				Goal:                "Test goal",
				IssueType:           "bug",
				AffectedFiles:       []string{"src/file.ts"},
				AffectedTests:       []string{"tests/file.test.ts"},
				ImplementationSteps: []string{"Step 1"},
				DependsOn:           []string{tt.depID},
			}
			err := svc.ValidateCreateInput(input)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error '%s' but got nil", tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}

	os.RemoveAll(filepath.Dir(paths.MandorDirPath()))
}

func TestIssueService_ListIssues(t *testing.T) {
	paths, err := fs.NewPathsFromRoot("/tmp/mandor-test-" + randomString(8))
	if err != nil {
		t.Fatalf("Failed to create paths: %v", err)
	}

	ws := &domain.Workspace{
		ID:            "test-workspace",
		Name:          "Test Workspace",
		Version:       "1.0.0",
		SchemaVersion: "1.0.0",
		Config: domain.WorkspaceConfig{
			DefaultPriority: "P3",
			StrictMode:      false,
		},
	}

	writer := fs.NewWriter(paths)
	if err := writer.CreateMandorDir(); err != nil {
		t.Fatalf("Failed to create mandor dir: %v", err)
	}
	if err := writer.WriteWorkspace(ws); err != nil {
		t.Fatalf("Failed to write workspace: %v", err)
	}

	if err := writer.CreateProjectDir("auth"); err != nil {
		t.Fatalf("Failed to create project dir: %v", err)
	}

	svc := service.NewIssueServiceWithPaths(paths)

	_, err = svc.CreateIssue(&domain.IssueCreateInput{
		ProjectID:           "auth",
		Name:                "Bug Issue",
		Goal:                "Test goal",
		IssueType:           "bug",
		AffectedFiles:       []string{"src/file1.ts"},
		AffectedTests:       []string{"tests/file1.test.ts"},
		ImplementationSteps: []string{"Step 1"},
		Priority:            "P1",
	})
	if err != nil {
		t.Fatalf("Failed to create bug issue: %v", err)
	}

	_, err = svc.CreateIssue(&domain.IssueCreateInput{
		ProjectID:           "auth",
		Name:                "Security Issue",
		Goal:                "Test goal",
		IssueType:           "security",
		AffectedFiles:       []string{"src/file2.ts"},
		AffectedTests:       []string{"tests/file2.test.ts"},
		ImplementationSteps: []string{"Step 1"},
		Priority:            "P0",
	})
	if err != nil {
		t.Fatalf("Failed to create security issue: %v", err)
	}

	output, err := svc.ListIssues(&domain.IssueListInput{
		ProjectID: "auth",
	})
	if err != nil {
		t.Fatalf("Failed to list issues: %v", err)
	}

	if output.Total != 2 {
		t.Errorf("Expected 2 issues but got %d", output.Total)
	}

	output, err = svc.ListIssues(&domain.IssueListInput{
		ProjectID: "auth",
		IssueType: "bug",
	})
	if err != nil {
		t.Fatalf("Failed to list issues by type: %v", err)
	}

	if output.Total != 1 {
		t.Errorf("Expected 1 bug issue but got %d", output.Total)
	}

	output, err = svc.ListIssues(&domain.IssueListInput{
		ProjectID: "auth",
		Priority:  "P0",
	})
	if err != nil {
		t.Fatalf("Failed to list issues by priority: %v", err)
	}

	if output.Total != 1 {
		t.Errorf("Expected 1 P0 issue but got %d", output.Total)
	}

	os.RemoveAll(filepath.Dir(paths.MandorDirPath()))
}

func TestIssueService_GetIssueDetail(t *testing.T) {
	paths, err := fs.NewPathsFromRoot("/tmp/mandor-test-" + randomString(8))
	if err != nil {
		t.Fatalf("Failed to create paths: %v", err)
	}

	ws := &domain.Workspace{
		ID:            "test-workspace",
		Name:          "Test Workspace",
		Version:       "1.0.0",
		SchemaVersion: "1.0.0",
		Config: domain.WorkspaceConfig{
			DefaultPriority: "P3",
			StrictMode:      false,
		},
	}

	writer := fs.NewWriter(paths)
	if err := writer.CreateMandorDir(); err != nil {
		t.Fatalf("Failed to create mandor dir: %v", err)
	}
	if err := writer.WriteWorkspace(ws); err != nil {
		t.Fatalf("Failed to write workspace: %v", err)
	}

	if err := writer.CreateProjectDir("auth"); err != nil {
		t.Fatalf("Failed to create project dir: %v", err)
	}

	svc := service.NewIssueServiceWithPaths(paths)

	issue, err := svc.CreateIssue(&domain.IssueCreateInput{
		ProjectID:           "auth",
		Name:                "Test Issue",
		Goal:                "Test goal description",
		IssueType:           "bug",
		Priority:            "P2",
		AffectedFiles:       []string{"src/file1.ts", "src/file2.ts"},
		AffectedTests:       []string{"tests/file1.test.ts"},
		ImplementationSteps: []string{"Step 1", "Step 2"},
		LibraryNeeds:        []string{"lib1"},
	})
	if err != nil {
		t.Fatalf("Failed to create issue: %v", err)
	}

	detail, err := svc.GetIssueDetail(&domain.IssueDetailInput{
		ProjectID: "auth",
		IssueID:   issue.ID,
	})
	if err != nil {
		t.Fatalf("Failed to get issue detail: %v", err)
	}

	if detail.Name != "Test Issue" {
		t.Errorf("Expected name 'Test Issue' but got '%s'", detail.Name)
	}

	if detail.IssueType != "bug" {
		t.Errorf("Expected type 'bug' but got '%s'", detail.IssueType)
	}

	if len(detail.AffectedFiles) != 2 {
		t.Errorf("Expected 2 affected files but got %d", len(detail.AffectedFiles))
	}

	os.RemoveAll(filepath.Dir(paths.MandorDirPath()))
}

func TestIssueService_UpdateIssue(t *testing.T) {
	paths, err := fs.NewPathsFromRoot("/tmp/mandor-test-" + randomString(8))
	if err != nil {
		t.Fatalf("Failed to create paths: %v", err)
	}

	ws := &domain.Workspace{
		ID:            "test-workspace",
		Name:          "Test Workspace",
		Version:       "1.0.0",
		SchemaVersion: "1.0.0",
		Config: domain.WorkspaceConfig{
			DefaultPriority: "P3",
			StrictMode:      false,
		},
	}

	writer := fs.NewWriter(paths)
	if err := writer.CreateMandorDir(); err != nil {
		t.Fatalf("Failed to create mandor dir: %v", err)
	}
	if err := writer.WriteWorkspace(ws); err != nil {
		t.Fatalf("Failed to write workspace: %v", err)
	}

	if err := writer.CreateProjectDir("auth"); err != nil {
		t.Fatalf("Failed to create project dir: %v", err)
	}

	svc := service.NewIssueServiceWithPaths(paths)

	issue, err := svc.CreateIssue(&domain.IssueCreateInput{
		ProjectID:           "auth",
		Name:                "Original Name",
		Goal:                "Original goal",
		IssueType:           "bug",
		AffectedFiles:       []string{"src/file1.ts"},
		AffectedTests:       []string{"tests/file1.test.ts"},
		ImplementationSteps: []string{"Step 1"},
	})
	if err != nil {
		t.Fatalf("Failed to create issue: %v", err)
	}

	newName := "Updated Name"
	changes, err := svc.UpdateIssue(&domain.IssueUpdateInput{
		ProjectID: "auth",
		IssueID:   issue.ID,
		Name:      &newName,
	})
	if err != nil {
		t.Fatalf("Failed to update issue: %v", err)
	}

	if len(changes) != 1 || changes[0] != "name" {
		t.Errorf("Expected name change but got: %v", changes)
	}

	detail, err := svc.GetIssueDetail(&domain.IssueDetailInput{
		ProjectID: "auth",
		IssueID:   issue.ID,
	})
	if err != nil {
		t.Fatalf("Failed to get issue detail: %v", err)
	}

	if detail.Name != "Updated Name" {
		t.Errorf("Expected name 'Updated Name' but got '%s'", detail.Name)
	}

	os.RemoveAll(filepath.Dir(paths.MandorDirPath()))
}

func TestIssueService_ReopenIssue(t *testing.T) {
	paths, err := fs.NewPathsFromRoot("/tmp/mandor-test-" + randomString(8))
	if err != nil {
		t.Fatalf("Failed to create paths: %v", err)
	}

	ws := &domain.Workspace{
		ID:            "test-workspace",
		Name:          "Test Workspace",
		Version:       "1.0.0",
		SchemaVersion: "1.0.0",
		Config: domain.WorkspaceConfig{
			DefaultPriority: "P3",
			StrictMode:      false,
		},
	}

	writer := fs.NewWriter(paths)
	if err := writer.CreateMandorDir(); err != nil {
		t.Fatalf("Failed to create mandor dir: %v", err)
	}
	if err := writer.WriteWorkspace(ws); err != nil {
		t.Fatalf("Failed to write workspace: %v", err)
	}

	if err := writer.CreateProjectDir("auth"); err != nil {
		t.Fatalf("Failed to create project dir: %v", err)
	}

	svc := service.NewIssueServiceWithPaths(paths)

	issue, err := svc.CreateIssue(&domain.IssueCreateInput{
		ProjectID:           "auth",
		Name:                "Test Issue",
		Goal:                "Test goal",
		IssueType:           "bug",
		AffectedFiles:       []string{"src/file1.ts"},
		AffectedTests:       []string{"tests/file1.test.ts"},
		ImplementationSteps: []string{"Step 1"},
	})
	if err != nil {
		t.Fatalf("Failed to create issue: %v", err)
	}

	reason := "No longer needed"
	_, err = svc.UpdateIssue(&domain.IssueUpdateInput{
		ProjectID: "auth",
		IssueID:   issue.ID,
		Cancel:    true,
		Reason:    &reason,
	})
	if err != nil {
		t.Fatalf("Failed to cancel issue: %v", err)
	}

	detail, _ := svc.GetIssueDetail(&domain.IssueDetailInput{
		ProjectID:      "auth",
		IssueID:        issue.ID,
		IncludeDeleted: true,
	})

	if detail.Status != domain.IssueStatusCancelled {
		t.Errorf("Expected status 'cancelled' but got '%s'", detail.Status)
	}

	_, err = svc.UpdateIssue(&domain.IssueUpdateInput{
		ProjectID: "auth",
		IssueID:   issue.ID,
		Reopen:    true,
	})
	if err != nil {
		t.Fatalf("Failed to reopen issue: %v", err)
	}

	detail, _ = svc.GetIssueDetail(&domain.IssueDetailInput{
		ProjectID:      "auth",
		IssueID:        issue.ID,
		IncludeDeleted: true,
	})

	if detail.Status != domain.IssueStatusOpen {
		t.Errorf("Expected status 'open' after reopen but got '%s'", detail.Status)
	}

	os.RemoveAll(filepath.Dir(paths.MandorDirPath()))
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[i%len(letters)]
	}
	return string(b)
}
