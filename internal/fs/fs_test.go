package fs

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"mandor/internal/domain"
)

// ==========================
// Paths Tests
// ==========================

func TestPathsConstants(t *testing.T) {
	if MandorDir != ".mandor" {
		t.Errorf("MandorDir = %q, want %q", MandorDir, ".mandor")
	}
	if BacklogsDir != "backlogs" {
		t.Errorf("BacklogsDir = %q, want %q", BacklogsDir, "backlogs")
	}
	if WorkspaceFile != "workspace.json" {
		t.Errorf("WorkspaceFile = %q, want %q", WorkspaceFile, "workspace.json")
	}
}

func TestNewPathsFromRoot(t *testing.T) {
	tests := []struct {
		name     string
		root     string
		expected string
	}{
		{"current dir", "/test/workspace", "/test/workspace"},
		{"home dir", "/home/user", "/home/user"},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paths, err := NewPathsFromRoot(tt.root)
			if err != nil {
				t.Errorf("NewPathsFromRoot() error = %v", err)
			}
			if paths.WorkspaceRoot != tt.expected {
				t.Errorf("WorkspaceRoot = %q, want %q", paths.WorkspaceRoot, tt.expected)
			}
		})
	}
}

func TestPathsMandorDirPath(t *testing.T) {
	paths, _ := NewPathsFromRoot("/workspace")
	expected := "/workspace/.mandor"
	if got := paths.MandorDirPath(); got != expected {
		t.Errorf("MandorDirPath() = %q, want %q", got, expected)
	}
}

func TestPathsWorkspacePath(t *testing.T) {
	paths, _ := NewPathsFromRoot("/workspace")
	expected := "/workspace/.mandor/workspace.json"
	if got := paths.WorkspacePath(); got != expected {
		t.Errorf("WorkspacePath() = %q, want %q", got, expected)
	}
}

func TestPathsBacklogsDirPath(t *testing.T) {
	paths, _ := NewPathsFromRoot("/workspace")
	expected := "/workspace/.mandor/backlogs"
	if got := paths.BacklogsDirPath(); got != expected {
		t.Errorf("BacklogsDirPath() = %q, want %q", got, expected)
	}
}

func TestPathsBacklogDirPath(t *testing.T) {
	paths, _ := NewPathsFromRoot("/workspace")
	expected := "/workspace/.mandor/backlogs/auth"
	if got := paths.BacklogDirPath("auth"); got != expected {
		t.Errorf("BacklogDirPath() = %q, want %q", got, expected)
	}
}

func TestPathsBacklogMetadataPath(t *testing.T) {
	paths, _ := NewPathsFromRoot("/workspace")
	expected := "/workspace/.mandor/backlogs/auth/backlog.jsonl"
	if got := paths.BacklogMetadataPath("auth"); got != expected {
		t.Errorf("BacklogMetadataPath() = %q, want %q", got, expected)
	}
}

func TestPathsBacklogSchemaPath(t *testing.T) {
	paths, _ := NewPathsFromRoot("/workspace")
	expected := "/workspace/.mandor/backlogs/auth/schema.json"
	if got := paths.BacklogSchemaPath("auth"); got != expected {
		t.Errorf("BacklogSchemaPath() = %q, want %q", got, expected)
	}
}

func TestPathsBacklogFeaturesPath(t *testing.T) {
	paths, _ := NewPathsFromRoot("/workspace")
	expected := "/workspace/.mandor/backlogs/auth/features.jsonl"
	if got := paths.BacklogFeaturesPath("auth"); got != expected {
		t.Errorf("BacklogFeaturesPath() = %q, want %q", got, expected)
	}
}

func TestPathsBacklogTasksPath(t *testing.T) {
	paths, _ := NewPathsFromRoot("/workspace")
	expected := "/workspace/.mandor/backlogs/auth/tasks.jsonl"
	if got := paths.BacklogTasksPath("auth"); got != expected {
		t.Errorf("BacklogTasksPath() = %q, want %q", got, expected)
	}
}

func TestPathsBacklogIssuesPath(t *testing.T) {
	paths, _ := NewPathsFromRoot("/workspace")
	expected := "/workspace/.mandor/backlogs/auth/issues.jsonl"
	if got := paths.BacklogIssuesPath("auth"); got != expected {
		t.Errorf("BacklogIssuesPath() = %q, want %q", got, expected)
	}
}

func TestPathsSessionNotesPath(t *testing.T) {
	paths, _ := NewPathsFromRoot("/workspace")
	expected := "/workspace/.mandor/session-notes.jsonl"
	if got := paths.SessionNotesPath(); got != expected {
		t.Errorf("SessionNotesPath() = %q, want %q", got, expected)
	}
}

func TestPathsBriefPath(t *testing.T) {
	paths, _ := NewPathsFromRoot("/workspace")
	expected := "/workspace/.mandor/backlogs/auth/brief.md"
	if got := paths.BriefPath("auth"); got != expected {
		t.Errorf("BriefPath() = %q, want %q", got, expected)
	}
}

func TestPathsSpecsDirPath(t *testing.T) {
	paths, _ := NewPathsFromRoot("/workspace")
	expected := "/workspace/.mandor/backlogs/auth/specs"
	if got := paths.SpecsDirPath("auth"); got != expected {
		t.Errorf("SpecsDirPath() = %q, want %q", got, expected)
	}
}

func TestPathsSpecPath(t *testing.T) {
	paths, _ := NewPathsFromRoot("/workspace")
	expected := "/workspace/.mandor/backlogs/auth/specs/test-cap.md"
	if got := paths.SpecPath("auth", "test-cap"); got != expected {
		t.Errorf("SpecPath() = %q, want %q", got, expected)
	}
}

func TestPathsBlueprintPath(t *testing.T) {
	paths, _ := NewPathsFromRoot("/workspace")
	expected := "/workspace/.mandor/backlogs/auth/blueprint.md"
	if got := paths.BlueprintPath("auth"); got != expected {
		t.Errorf("BlueprintPath() = %q, want %q", got, expected)
	}
}

// ==========================
// Paths.BacklogDirExists Tests
// ==========================

func TestBacklogDirExists(t *testing.T) {
	// Create a temporary directory structure
	tmpDir, err := os.MkdirTemp("", "mandor-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	paths, _ := NewPathsFromRoot(tmpDir)

	// Test non-existent backlog
	if paths.BacklogDirExists("nonexistent") {
		t.Error("BacklogDirExists() for nonexistent should return false")
	}

	// Create backlog directory
	backlogDir := filepath.Join(paths.BacklogsDirPath(), "test-backlog")
	if err := os.MkdirAll(backlogDir, 0755); err != nil {
		t.Fatalf("Failed to create backlog dir: %v", err)
	}

	// Test existing backlog
	if !paths.BacklogDirExists("test-backlog") {
		t.Error("BacklogDirExists() for existing should return true")
	}
}

// ==========================
// Reader Tests (with temp directory)
// ==========================

func setupTestWorkspace(t *testing.T) (string, *Paths) {
	tmpDir, err := os.MkdirTemp("", "mandor-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	paths, _ := NewPathsFromRoot(tmpDir)

	// Create mandor directory structure
	mandorDir := paths.MandorDirPath()
	backlogsDir := paths.BacklogsDirPath()
	if err := os.MkdirAll(mandorDir, 0755); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to create mandor dir: %v", err)
	}
	if err := os.MkdirAll(backlogsDir, 0755); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to create backlogs dir: %v", err)
	}

	return tmpDir, paths
}

func TestReaderWorkspaceExists(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	reader := NewReader(paths)

	// Workspace file doesn't exist yet
	if reader.WorkspaceExists() {
		t.Error("WorkspaceExists() should return false when workspace.json doesn't exist")
	}

	// Create workspace file
	workspacePath := paths.WorkspacePath()
	if err := os.WriteFile(workspacePath, []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to write workspace file: %v", err)
	}

	// Now it should exist
	if !reader.WorkspaceExists() {
		t.Error("WorkspaceExists() should return true when workspace.json exists")
	}
}

func TestReaderBacklogExists(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	reader := NewReader(paths)

	// Backlog doesn't exist yet
	if reader.BacklogExists("auth") {
		t.Error("BacklogExists() should return false for non-existent backlog")
	}

	// Create backlog directory
	backlogDir := paths.BacklogDirPath("auth")
	if err := os.MkdirAll(backlogDir, 0755); err != nil {
		t.Fatalf("Failed to create backlog dir: %v", err)
	}

	// Now it should exist
	if !reader.BacklogExists("auth") {
		t.Error("BacklogExists() should return true for existing backlog")
	}
}

func TestReaderListBacklogs(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	reader := NewReader(paths)

	// No backlogs yet
	backlogs, err := reader.ListBacklogs(false)
	if err != nil {
		t.Errorf("ListBacklogs() error = %v", err)
	}
	if len(backlogs) != 0 {
		t.Errorf("ListBacklogs() count = %d, want 0", len(backlogs))
	}

	// Create some backlogs
	os.MkdirAll(filepath.Join(paths.BacklogsDirPath(), "auth"), 0755)
	os.MkdirAll(filepath.Join(paths.BacklogsDirPath(), "billing"), 0755)
	os.MkdirAll(filepath.Join(paths.BacklogsDirPath(), "api"), 0755)

	backlogs, err = reader.ListBacklogs(false)
	if err != nil {
		t.Errorf("ListBacklogs() error = %v", err)
	}
	if len(backlogs) != 3 {
		t.Errorf("ListBacklogs() count = %d, want 3", len(backlogs))
	}
}

func TestReaderReadFile(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	reader := NewReader(paths)

	// Non-existent file
	content, err := reader.ReadFile("/nonexistent/file.txt")
	if err != nil {
		t.Errorf("ReadFile() error = %v", err)
	}
	if content != "" {
		t.Errorf("ReadFile() for non-existent should return empty string")
	}

	// Existing file
	testFile := filepath.Join(tmpDir, "test.txt")
	testContent := "Hello, World!"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	content, err = reader.ReadFile(testFile)
	if err != nil {
		t.Errorf("ReadFile() error = %v", err)
	}
	if content != testContent {
		t.Errorf("ReadFile() = %q, want %q", content, testContent)
	}
}

// ==========================
// Writer Tests (with temp directory)
// ==========================

func TestWriterMandorDirExists(t *testing.T) {
	// Create temp dir without mandor structure
	tmpDir, err := os.MkdirTemp("", "mandor-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	paths, _ := NewPathsFromRoot(tmpDir)
	writer := NewWriter(paths)

	// Mandor dir doesn't exist yet (setup doesn't create it in this test)
	if writer.MandorDirExists() {
		t.Error("MandorDirExists() should return false when .mandor doesn't exist")
	}

	// Create mandor dir
	if err := os.MkdirAll(paths.MandorDirPath(), 0755); err != nil {
		t.Fatalf("Failed to create mandor dir: %v", err)
	}

	// Now it should exist
	if !writer.MandorDirExists() {
		t.Error("MandorDirExists() should return true when .mandor exists")
	}
}

func TestWriterIsDirWritable(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	writer := NewWriter(paths)

	// Writable directory
	if !writer.IsDirWritable(tmpDir) {
		t.Error("IsDirWritable() should return true for writable directory")
	}

	// Non-writable directory (by removing write permission)
	nonWritableDir := filepath.Join(tmpDir, "readonly")
	os.MkdirAll(nonWritableDir, 0555)
	if writer.IsDirWritable(nonWritableDir) {
		t.Error("IsDirWritable() should return false for non-writable directory")
	}
}

func TestWriterCreateDirectory(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	writer := NewWriter(paths)

	newDir := filepath.Join(tmpDir, "new", "nested", "dir")
	err := writer.CreateDirectory(newDir)
	if err != nil {
		t.Errorf("CreateDirectory() error = %v", err)
	}

	// Verify directory was created
	info, err := os.Stat(newDir)
	if err != nil {
		t.Errorf("Directory not created: %v", err)
	}
	if !info.IsDir() {
		t.Error("Created path should be a directory")
	}
}

func TestWriterDeleteFile(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	writer := NewWriter(paths)

	// Non-existent file
	err := writer.DeleteFile("/nonexistent/file.txt")
	if err == nil {
		t.Error("DeleteFile() should return error for non-existent file")
	}

	// Existing file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	err = writer.DeleteFile(testFile)
	if err != nil {
		t.Errorf("DeleteFile() error = %v", err)
	}

	// Verify file was deleted
	_, err = os.Stat(testFile)
	if err == nil {
		t.Error("File should have been deleted")
	}
}

func TestWriterCheckBacklogWritable(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	writer := NewWriter(paths)

	// Non-existent backlog
	if writer.CheckBacklogWritable("nonexistent") {
		t.Error("CheckBacklogWritable() should return false for non-existent backlog")
	}

	// Create backlog directory
	backlogDir := paths.BacklogDirPath("auth")
	if err := os.MkdirAll(backlogDir, 0755); err != nil {
		t.Fatalf("Failed to create backlog dir: %v", err)
	}

	// Should be writable
	if !writer.CheckBacklogWritable("auth") {
		t.Error("CheckBacklogWritable() should return true for writable backlog")
	}
}

func TestWriterBacklogsDirWritable(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	writer := NewWriter(paths)

	// Backlogs dir should be writable
	if !writer.BacklogsDirWritable() {
		t.Error("BacklogsDirWritable() should return true for writable directory")
	}
}

// ==========================
// Reader.ReadWorkspace Tests
// ==========================

func TestReaderReadWorkspace_Success(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	reader := NewReader(paths)

	// Create workspace file
	ws := &domain.Workspace{
		ID:            "test-workspace",
		Name:          "Test Workspace",
		Version:       "1.0.0",
		SchemaVersion: "1.0",
		CreatedBy:     "test-user",
		Config:        domain.DefaultWorkspaceConfig(),
	}
	data, _ := json.MarshalIndent(ws, "", "  ")
	if err := os.WriteFile(paths.WorkspacePath(), data, 0644); err != nil {
		t.Fatalf("Failed to write workspace file: %v", err)
	}

	result, err := reader.ReadWorkspace()
	if err != nil {
		t.Errorf("ReadWorkspace() error = %v", err)
	}
	if result == nil {
		t.Fatal("ReadWorkspace() returned nil")
	}
	if result.ID != ws.ID {
		t.Errorf("ReadWorkspace() ID = %q, want %q", result.ID, ws.ID)
	}
	if result.Name != ws.Name {
		t.Errorf("ReadWorkspace() Name = %q, want %q", result.Name, ws.Name)
	}
}

func TestReaderReadWorkspace_NotFound(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	reader := NewReader(paths)

	_, err := reader.ReadWorkspace()
	if err == nil {
		t.Error("ReadWorkspace() should return error when workspace not found")
	}
	if !strings.Contains(err.Error(), "not initialized") {
		t.Errorf("ReadWorkspace() error should mention 'not initialized', got: %v", err)
	}
}

func TestReaderReadWorkspace_CorruptedJSON(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	reader := NewReader(paths)

	// Create corrupted workspace file
	if err := os.WriteFile(paths.WorkspacePath(), []byte("{ invalid json }"), 0644); err != nil {
		t.Fatalf("Failed to write corrupted workspace file: %v", err)
	}

	_, err := reader.ReadWorkspace()
	if err == nil {
		t.Error("ReadWorkspace() should return error for corrupted JSON")
	}
	if !strings.Contains(err.Error(), "corrupted") {
		t.Errorf("ReadWorkspace() error should mention 'corrupted', got: %v", err)
	}
}

// ==========================
// Reader.ReadBacklogMetadata Tests
// ==========================

func TestReaderReadBacklogMetadata_Success(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	reader := NewReader(paths)

	// Create backlog directory and metadata file
	backlog := &domain.Backlog{
		ID:        "test-backlog",
		Name:      "Test Backlog",
		Goal:      strings.Repeat("A", 500),
		Status:    domain.BacklogStatusActive,
		Strict:    false,
		CreatedBy: "test-user",
	}
	if err := os.MkdirAll(paths.BacklogDirPath("test-backlog"), 0755); err != nil {
		t.Fatalf("Failed to create backlog dir: %v", err)
	}
	data, _ := json.MarshalIndent(backlog, "", "  ")
	if err := os.WriteFile(paths.BacklogMetadataPath("test-backlog"), data, 0644); err != nil {
		t.Fatalf("Failed to write backlog metadata: %v", err)
	}

	result, err := reader.ReadBacklogMetadata("test-backlog")
	if err != nil {
		t.Errorf("ReadBacklogMetadata() error = %v", err)
	}
	if result == nil {
		t.Fatal("ReadBacklogMetadata() returned nil")
	}
	if result.ID != backlog.ID {
		t.Errorf("ReadBacklogMetadata() ID = %q, want %q", result.ID, backlog.ID)
	}
	if result.Name != backlog.Name {
		t.Errorf("ReadBacklogMetadata() Name = %q, want %q", result.Name, backlog.Name)
	}
}

func TestReaderReadBacklogMetadata_NotFound(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	reader := NewReader(paths)

	_, err := reader.ReadBacklogMetadata("nonexistent")
	if err == nil {
		t.Error("ReadBacklogMetadata() should return error for non-existent backlog")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("ReadBacklogMetadata() error should mention 'not found', got: %v", err)
	}
}

func TestReaderReadBacklogMetadata_CorruptedJSON(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	reader := NewReader(paths)

	// Create backlog dir and corrupted metadata file
	if err := os.MkdirAll(paths.BacklogDirPath("test-backlog"), 0755); err != nil {
		t.Fatalf("Failed to create backlog dir: %v", err)
	}
	if err := os.WriteFile(paths.BacklogMetadataPath("test-backlog"), []byte("{ invalid }"), 0644); err != nil {
		t.Fatalf("Failed to write corrupted metadata: %v", err)
	}

	_, err := reader.ReadBacklogMetadata("test-backlog")
	if err == nil {
		t.Error("ReadBacklogMetadata() should return error for corrupted JSON")
	}
}

// ==========================
// Reader.ReadBacklogSchema Tests
// ==========================

func TestReaderReadBacklogSchema_Success(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	reader := NewReader(paths)

	// Create backlog directory and schema file
	schema := domain.DefaultBacklogSchema(domain.DependencySameBacklogOnly, domain.DependencyCrossBacklogAllowed, domain.DependencySameBacklogOnly)
	if err := os.MkdirAll(paths.BacklogDirPath("test-backlog"), 0755); err != nil {
		t.Fatalf("Failed to create backlog dir: %v", err)
	}
	data, _ := json.MarshalIndent(schema, "", "  ")
	if err := os.WriteFile(paths.BacklogSchemaPath("test-backlog"), data, 0644); err != nil {
		t.Fatalf("Failed to write backlog schema: %v", err)
	}

	result, err := reader.ReadBacklogSchema("test-backlog")
	if err != nil {
		t.Errorf("ReadBacklogSchema() error = %v", err)
	}
	if result == nil {
		t.Fatal("ReadBacklogSchema() returned nil")
	}
	if result.Version != schema.Version {
		t.Errorf("ReadBacklogSchema() Version = %q, want %q", result.Version, schema.Version)
	}
	if result.Rules.Task.Dependency != schema.Rules.Task.Dependency {
		t.Errorf("ReadBacklogSchema() Task.Dependency = %q, want %q", result.Rules.Task.Dependency, schema.Rules.Task.Dependency)
	}
}

func TestReaderReadBacklogSchema_NotFound(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	reader := NewReader(paths)

	// Create backlog dir without schema
	if err := os.MkdirAll(paths.BacklogDirPath("test-backlog"), 0755); err != nil {
		t.Fatalf("Failed to create backlog dir: %v", err)
	}

	_, err := reader.ReadBacklogSchema("test-backlog")
	if err == nil {
		t.Error("ReadBacklogSchema() should return error when schema not found")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("ReadBacklogSchema() error should mention 'not found', got: %v", err)
	}
}

// ==========================
// Reader.CountLines and CountEntityLines Tests
// ==========================

func TestReaderCountLines(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	reader := NewReader(paths)

	// Non-existent file
	count, err := reader.CountLines("/nonexistent/file.txt")
	if err != nil {
		t.Errorf("CountLines() error for non-existent = %v", err)
	}
	if count != 0 {
		t.Errorf("CountLines() for non-existent = %d, want 0", count)
	}

	// File with 3 lines
	testFile := filepath.Join(tmpDir, "test.txt")
	content := "line1\nline2\nline3\n"
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	count, err = reader.CountLines(testFile)
	if err != nil {
		t.Errorf("CountLines() error = %v", err)
	}
	// Count includes newline characters, so 3 lines = 4 counts (including final newline)
	if count < 3 {
		t.Errorf("CountLines() = %d, want at least 3", count)
	}
}

func TestReaderCountEntityLines(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	reader := NewReader(paths)

	// Non-existent file
	count, err := reader.CountEntityLines("/nonexistent/file.txt")
	if err != nil {
		t.Errorf("CountEntityLines() error for non-existent = %v", err)
	}
	if count != 0 {
		t.Errorf("CountEntityLines() for non-existent = %d, want 0", count)
	}

	// NDJSON file with 3 entries
	testFile := filepath.Join(tmpDir, "test.jsonl")
	content := `{"id":"1"}
{"id":"2"}
{"id":"3"}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	count, err = reader.CountEntityLines(testFile)
	if err != nil {
		t.Errorf("CountEntityLines() error = %v", err)
	}
	if count != 3 {
		t.Errorf("CountEntityLines() = %d, want 3", count)
	}
}

// ==========================
// Reader.ReadNDJSON Tests
// ==========================

func TestReaderReadNDJSON_Success(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	reader := NewReader(paths)

	// NDJSON file with 3 entries
	testFile := filepath.Join(tmpDir, "test.jsonl")
	content := `{"id":"1","name":"first"}
{"id":"2","name":"second"}
{"id":"3","name":"third"}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	var ids []string
	err := reader.ReadNDJSON(testFile, func(raw []byte) error {
		var entry map[string]string
		if err := json.Unmarshal(raw, &entry); err != nil {
			return err
		}
		if id, ok := entry["id"]; ok {
			ids = append(ids, id)
		}
		return nil
	})
	if err != nil {
		t.Errorf("ReadNDJSON() error = %v", err)
	}
	if len(ids) != 3 {
		t.Errorf("ReadNDJSON() collected %d entries, want 3", len(ids))
	}
}

func TestReaderReadNDJSON_EmptyFile(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	reader := NewReader(paths)

	// Empty NDJSON file
	testFile := filepath.Join(tmpDir, "empty.jsonl")
	if err := os.WriteFile(testFile, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	var count int
	err := reader.ReadNDJSON(testFile, func(raw []byte) error {
		count++
		return nil
	})
	if err != nil {
		t.Errorf("ReadNDJSON() error on empty file = %v", err)
	}
	if count != 0 {
		t.Errorf("ReadNDJSON() collected %d entries, want 0", count)
	}
}

func TestReaderReadNDJSON_NotFound(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	reader := NewReader(paths)

	err := reader.ReadNDJSON("/nonexistent/file.jsonl", func(raw []byte) error {
		return nil
	})
	if err != nil {
		t.Errorf("ReadNDJSON() should not return error for non-existent file, got: %v", err)
	}
}

// ==========================
// Reader Entity Tests (ReadFeature, ReadTask, ReadIssue)
// ==========================

func TestReaderReadFeature_Success(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	reader := NewReader(paths)

	// Create backlog and feature file
	if err := os.MkdirAll(paths.BacklogDirPath("test-backlog"), 0755); err != nil {
		t.Fatalf("Failed to create backlog dir: %v", err)
	}

	feature := &domain.Feature{
		ID:           "feature-test-123",
		CapabilityID: "cap-001",
		SpecID:       "spec-001",
		BacklogID:    "test-backlog",
		Name:         "Test Feature",
		Goal:         strings.Repeat("Goal for test feature ", 20),
		Priority:     "P1",
		Status:       domain.FeatureStatusActive,
	}
	data, _ := json.MarshalIndent(feature, "", "  ")
	if err := os.WriteFile(paths.BacklogFeaturesPath("test-backlog"), data, 0644); err != nil {
		t.Fatalf("Failed to write feature file: %v", err)
	}

	result, err := reader.ReadFeature("test-backlog", "feature-test-123")
	if err != nil {
		t.Errorf("ReadFeature() error = %v", err)
	}
	if result == nil {
		t.Fatal("ReadFeature() returned nil")
	}
	if result.ID != feature.ID {
		t.Errorf("ReadFeature() ID = %q, want %q", result.ID, feature.ID)
	}
}

func TestReaderReadFeature_NotFound(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	reader := NewReader(paths)

	// Create empty feature file
	if err := os.MkdirAll(paths.BacklogDirPath("test-backlog"), 0755); err != nil {
		t.Fatalf("Failed to create backlog dir: %v", err)
	}
	if err := os.WriteFile(paths.BacklogFeaturesPath("test-backlog"), []byte(""), 0644); err != nil {
		t.Fatalf("Failed to write feature file: %v", err)
	}

	_, err := reader.ReadFeature("test-backlog", "nonexistent")
	if err == nil {
		t.Error("ReadFeature() should return error for non-existent feature")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("ReadFeature() error should mention 'not found', got: %v", err)
	}
}

func TestReaderReadTask_Success(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	reader := NewReader(paths)

	// Create backlog and task file
	if err := os.MkdirAll(paths.BacklogDirPath("test-backlog"), 0755); err != nil {
		t.Fatalf("Failed to create backlog dir: %v", err)
	}

	task := &domain.Task{
		ID:        "task-test-123456",
		FeatureID: "feature-001",
		SpecID:    "spec-001",
		BacklogID: "test-backlog",
		Name:      "Test Task",
		Goal:      strings.Repeat("Goal for test task ", 20),
		Priority:  "P2",
		Status:    domain.TaskStatusReady,
		ReadGates: domain.ReadGates{IsReadBrief: true, IsReadSpec: true},
	}
	data, _ := json.MarshalIndent(task, "", "  ")
	if err := os.WriteFile(paths.BacklogTasksPath("test-backlog"), data, 0644); err != nil {
		t.Fatalf("Failed to write task file: %v", err)
	}

	result, err := reader.ReadTask("test-backlog", "task-test-123456")
	if err != nil {
		t.Errorf("ReadTask() error = %v", err)
	}
	if result == nil {
		t.Fatal("ReadTask() returned nil")
	}
	if result.ID != task.ID {
		t.Errorf("ReadTask() ID = %q, want %q", result.ID, task.ID)
	}
}

func TestReaderReadTask_NotFound(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	reader := NewReader(paths)

	if err := os.MkdirAll(paths.BacklogDirPath("test-backlog"), 0755); err != nil {
		t.Fatalf("Failed to create backlog dir: %v", err)
	}
	if err := os.WriteFile(paths.BacklogTasksPath("test-backlog"), []byte(""), 0644); err != nil {
		t.Fatalf("Failed to write task file: %v", err)
	}

	_, err := reader.ReadTask("test-backlog", "nonexistent")
	if err == nil {
		t.Error("ReadTask() should return error for non-existent task")
	}
}

func TestReaderReadIssue_Success(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	reader := NewReader(paths)

	if err := os.MkdirAll(paths.BacklogDirPath("test-backlog"), 0755); err != nil {
		t.Fatalf("Failed to create backlog dir: %v", err)
	}

	issue := &domain.Issue{
		ID:            "issue-test-123456",
		BacklogID:     "test-backlog",
		Name:          "Test Issue",
		Goal:          strings.Repeat("Goal for test issue with enough characters to pass validation ", 5),
		IssueType:     domain.IssueTypeBug,
		Priority:      "P1",
		Status:        domain.IssueStatusOpen,
		AffectedFiles: []string{"src/main.go"},
	}
	data, _ := json.MarshalIndent(issue, "", "  ")
	if err := os.WriteFile(paths.BacklogIssuesPath("test-backlog"), data, 0644); err != nil {
		t.Fatalf("Failed to write issue file: %v", err)
	}

	result, err := reader.ReadIssue("test-backlog", "issue-test-123456")
	if err != nil {
		t.Errorf("ReadIssue() error = %v", err)
	}
	if result == nil {
		t.Fatal("ReadIssue() returned nil")
	}
	if result.ID != issue.ID {
		t.Errorf("ReadIssue() ID = %q, want %q", result.ID, issue.ID)
	}
}

func TestReaderReadIssue_NotFound(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	reader := NewReader(paths)

	if err := os.MkdirAll(paths.BacklogDirPath("test-backlog"), 0755); err != nil {
		t.Fatalf("Failed to create backlog dir: %v", err)
	}
	if err := os.WriteFile(paths.BacklogIssuesPath("test-backlog"), []byte(""), 0644); err != nil {
		t.Fatalf("Failed to write issue file: %v", err)
	}

	_, err := reader.ReadIssue("test-backlog", "nonexistent")
	if err == nil {
		t.Error("ReadIssue() should return error for non-existent issue")
	}
}

// ==========================
// Writer.CreateMandorDir Tests
// ==========================

func TestWriterCreateMandorDir_Success(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mandor-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	paths, _ := NewPathsFromRoot(tmpDir)
	writer := NewWriter(paths)

	err = writer.CreateMandorDir()
	if err != nil {
		t.Errorf("CreateMandorDir() error = %v", err)
	}

	// Verify directories were created
	if _, err := os.Stat(paths.MandorDirPath()); err != nil {
		t.Errorf("Mandor directory not created: %v", err)
	}
	if _, err := os.Stat(paths.BacklogsDirPath()); err != nil {
		t.Errorf("Backlogs directory not created: %v", err)
	}
}

// ==========================
// Writer.WriteWorkspace Tests
// ==========================

func TestWriterWriteWorkspace_Success(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	writer := NewWriter(paths)

	ws := &domain.Workspace{
		ID:            "test-ws",
		Name:          "Test Workspace",
		Version:       "1.0.0",
		SchemaVersion: "1.0",
		CreatedBy:     "test-user",
		Config:        domain.DefaultWorkspaceConfig(),
	}

	err := writer.WriteWorkspace(ws)
	if err != nil {
		t.Errorf("WriteWorkspace() error = %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(paths.WorkspacePath()); err != nil {
		t.Errorf("Workspace file not created: %v", err)
	}
}

// ==========================
// Writer.WriteJSON Tests
// ==========================

func TestWriterWriteJSON_Success(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	writer := NewWriter(paths)

	testData := map[string]string{"key": "value"}
	err := writer.WriteJSON(filepath.Join(tmpDir, "subdir", "test.json"), testData)
	if err != nil {
		t.Errorf("WriteJSON() error = %v", err)
	}

	// Verify file was created
	filePath := filepath.Join(tmpDir, "subdir", "test.json")
	if _, err := os.Stat(filePath); err != nil {
		t.Errorf("JSON file not created: %v", err)
	}
}

// ==========================
// Writer.AppendNDJSON Tests
// ==========================

func TestWriterAppendNDJSON_Success(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	writer := NewWriter(paths)

	testFile := filepath.Join(tmpDir, "test.jsonl")
	err := writer.AppendNDJSON(testFile, map[string]string{"id": "1"})
	if err != nil {
		t.Errorf("AppendNDJSON() error = %v", err)
	}

	// Append second entry
	err = writer.AppendNDJSON(testFile, map[string]string{"id": "2"})
	if err != nil {
		t.Errorf("AppendNDJSON() second entry error = %v", err)
	}

	// Verify file has both entries
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read NDJSON file: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	if len(lines) != 2 {
		t.Errorf("NDJSON file has %d lines, want 2", len(lines))
	}
}

// ==========================
// Writer.WriteBacklogMetadata Tests
// ==========================

func TestWriterWriteBacklogMetadata_Success(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	writer := NewWriter(paths)

	backlog := &domain.Backlog{
		ID:        "test-backlog",
		Name:      "Test Backlog",
		Goal:      strings.Repeat("A", 500),
		Status:    domain.BacklogStatusActive,
		Strict:    false,
		CreatedBy: "test-user",
	}

	err := writer.WriteBacklogMetadata("test-backlog", backlog)
	if err != nil {
		t.Errorf("WriteBacklogMetadata() error = %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(paths.BacklogMetadataPath("test-backlog")); err != nil {
		t.Errorf("Backlog metadata file not created: %v", err)
	}
}

// ==========================
// Writer.WriteBacklogSchema Tests
// ==========================

func TestWriterWriteBacklogSchema_Success(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	writer := NewWriter(paths)

	schema := domain.DefaultBacklogSchema(domain.DependencySameBacklogOnly, domain.DependencyCrossBacklogAllowed, domain.DependencySameBacklogOnly)

	err := writer.WriteBacklogSchema("test-backlog", &schema)
	if err != nil {
		t.Errorf("WriteBacklogSchema() error = %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(paths.BacklogSchemaPath("test-backlog")); err != nil {
		t.Errorf("Backlog schema file not created: %v", err)
	}
}

// ==========================
// Writer.CreateBacklogDir Tests
// ==========================

func TestWriterCreateBacklogDir_Success(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	writer := NewWriter(paths)

	err := writer.CreateBacklogDir("new-backlog")
	if err != nil {
		t.Errorf("CreateBacklogDir() error = %v", err)
	}

	// Verify backlog directory was created
	if _, err := os.Stat(paths.BacklogDirPath("new-backlog")); err != nil {
		t.Errorf("Backlog directory not created: %v", err)
	}

	// Verify entity files were created
	entityFiles := []string{"events.jsonl", "features.jsonl", "tasks.jsonl", "issues.jsonl"}
	for _, file := range entityFiles {
		filePath := filepath.Join(paths.BacklogDirPath("new-backlog"), file)
		if _, err := os.Stat(filePath); err != nil {
			t.Errorf("Entity file %s not created: %v", file, err)
		}
	}
}

// ==========================
// Writer.DeleteBacklogDir Tests
// ==========================

func TestWriterDeleteBacklogDir_Success(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	writer := NewWriter(paths)

	// Create a backlog first
	if err := writer.CreateBacklogDir("to-delete"); err != nil {
		t.Fatalf("Failed to create backlog dir: %v", err)
	}

	err := writer.DeleteBacklogDir("to-delete")
	if err != nil {
		t.Errorf("DeleteBacklogDir() error = %v", err)
	}

	// Verify directory was deleted
	if _, err := os.Stat(paths.BacklogDirPath("to-delete")); err == nil {
		t.Error("Backlog directory should have been deleted")
	}
}

// ==========================
// Writer Entity Write Tests (WriteFeature, WriteTask, WriteIssue)
// ==========================

func TestWriterWriteFeature_Success(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	writer := NewWriter(paths)

	if err := os.MkdirAll(paths.BacklogDirPath("test-backlog"), 0755); err != nil {
		t.Fatalf("Failed to create backlog dir: %v", err)
	}

	feature := &domain.Feature{
		ID:           "feature-new-123",
		CapabilityID: "cap-001",
		SpecID:       "spec-001",
		BacklogID:    "test-backlog",
		Name:         "New Feature",
		Goal:         strings.Repeat("Goal for new feature ", 20),
		Priority:     "P1",
		Status:       domain.FeatureStatusActive,
	}

	err := writer.WriteFeature("test-backlog", feature)
	if err != nil {
		t.Errorf("WriteFeature() error = %v", err)
	}

	// Verify file has the feature
	reader := NewReader(paths)
	result, err := reader.ReadFeature("test-backlog", "feature-new-123")
	if err != nil {
		t.Errorf("Failed to read back written feature: %v", err)
	}
	if result.ID != feature.ID {
		t.Errorf("Written feature ID = %q, want %q", result.ID, feature.ID)
	}
}

func TestWriterWriteTask_Success(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	writer := NewWriter(paths)

	if err := os.MkdirAll(paths.BacklogDirPath("test-backlog"), 0755); err != nil {
		t.Fatalf("Failed to create backlog dir: %v", err)
	}

	task := &domain.Task{
		ID:        "task-new-123456",
		FeatureID: "feature-001",
		SpecID:    "spec-001",
		BacklogID: "test-backlog",
		Name:      "New Task",
		Goal:      strings.Repeat("Goal for new task ", 20),
		Priority:  "P2",
		Status:    domain.TaskStatusReady,
	}

	err := writer.WriteTask("test-backlog", task)
	if err != nil {
		t.Errorf("WriteTask() error = %v", err)
	}

	reader := NewReader(paths)
	result, err := reader.ReadTask("test-backlog", "task-new-123456")
	if err != nil {
		t.Errorf("Failed to read back written task: %v", err)
	}
	if result.ID != task.ID {
		t.Errorf("Written task ID = %q, want %q", result.ID, task.ID)
	}
}

func TestWriterWriteIssue_Success(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	writer := NewWriter(paths)

	if err := os.MkdirAll(paths.BacklogDirPath("test-backlog"), 0755); err != nil {
		t.Fatalf("Failed to create backlog dir: %v", err)
	}

	issue := &domain.Issue{
		ID:        "issue-new-123456",
		BacklogID: "test-backlog",
		Name:      "New Issue",
		Goal:      strings.Repeat("Goal for new issue with enough characters ", 5),
		IssueType: domain.IssueTypeBug,
		Priority:  "P1",
		Status:    domain.IssueStatusOpen,
	}

	err := writer.WriteIssue("test-backlog", issue)
	if err != nil {
		t.Errorf("WriteIssue() error = %v", err)
	}

	reader := NewReader(paths)
	result, err := reader.ReadIssue("test-backlog", "issue-new-123456")
	if err != nil {
		t.Errorf("Failed to read back written issue: %v", err)
	}
	if result.ID != issue.ID {
		t.Errorf("Written issue ID = %q, want %q", result.ID, issue.ID)
	}
}

// ==========================
// Writer Entity Replace Tests (ReplaceFeature, ReplaceTask, ReplaceIssue)
// ==========================

func TestWriterReplaceFeature_Success(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	writer := NewWriter(paths)

	if err := os.MkdirAll(paths.BacklogDirPath("test-backlog"), 0755); err != nil {
		t.Fatalf("Failed to create backlog dir: %v", err)
	}

	// Write initial feature
	feature := &domain.Feature{
		ID:           "feature-replace-123",
		CapabilityID: "cap-001",
		SpecID:       "spec-001",
		BacklogID:    "test-backlog",
		Name:         "Original Feature",
		Goal:         strings.Repeat("Original goal ", 20),
		Priority:     "P1",
		Status:       domain.FeatureStatusActive,
	}
	if err := writer.WriteFeature("test-backlog", feature); err != nil {
		t.Fatalf("Failed to write initial feature: %v", err)
	}

	// Replace with updated feature
	updated := &domain.Feature{
		ID:           "feature-replace-123",
		CapabilityID: "cap-001",
		SpecID:       "spec-001",
		BacklogID:    "test-backlog",
		Name:         "Updated Feature",
		Goal:         strings.Repeat("Updated goal ", 20),
		Priority:     "P2",
		Status:       domain.FeatureStatusDone,
	}
	err := writer.ReplaceFeature("test-backlog", updated)
	if err != nil {
		t.Errorf("ReplaceFeature() error = %v", err)
	}

	// Verify updated feature
	reader := NewReader(paths)
	result, err := reader.ReadFeature("test-backlog", "feature-replace-123")
	if err != nil {
		t.Errorf("Failed to read updated feature: %v", err)
	}
	if result.Name != updated.Name {
		t.Errorf("Feature name = %q, want %q", result.Name, updated.Name)
	}
	if result.Status != updated.Status {
		t.Errorf("Feature status = %q, want %q", result.Status, updated.Status)
	}
}

func TestWriterReplaceTask_Success(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	writer := NewWriter(paths)

	if err := os.MkdirAll(paths.BacklogDirPath("test-backlog"), 0755); err != nil {
		t.Fatalf("Failed to create backlog dir: %v", err)
	}

	task := &domain.Task{
		ID:        "task-replace-123456",
		FeatureID: "feature-001",
		SpecID:    "spec-001",
		BacklogID: "test-backlog",
		Name:      "Original Task",
		Goal:      strings.Repeat("Original task goal ", 20),
		Priority:  "P2",
		Status:    domain.TaskStatusReady,
	}
	if err := writer.WriteTask("test-backlog", task); err != nil {
		t.Fatalf("Failed to write initial task: %v", err)
	}

	updated := &domain.Task{
		ID:        "task-replace-123456",
		FeatureID: "feature-001",
		SpecID:    "spec-001",
		BacklogID: "test-backlog",
		Name:      "Updated Task",
		Goal:      strings.Repeat("Updated task goal ", 20),
		Priority:  "P1",
		Status:    domain.TaskStatusInProgress,
	}
	err := writer.ReplaceTask("test-backlog", updated)
	if err != nil {
		t.Errorf("ReplaceTask() error = %v", err)
	}

	reader := NewReader(paths)
	result, err := reader.ReadTask("test-backlog", "task-replace-123456")
	if err != nil {
		t.Errorf("Failed to read updated task: %v", err)
	}
	if result.Name != updated.Name {
		t.Errorf("Task name = %q, want %q", result.Name, updated.Name)
	}
}

func TestWriterReplaceIssue_Success(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	writer := NewWriter(paths)

	if err := os.MkdirAll(paths.BacklogDirPath("test-backlog"), 0755); err != nil {
		t.Fatalf("Failed to create backlog dir: %v", err)
	}

	issue := &domain.Issue{
		ID:        "issue-replace-123456",
		BacklogID: "test-backlog",
		Name:      "Original Issue",
		Goal:      strings.Repeat("Original issue goal with enough characters ", 4),
		IssueType: domain.IssueTypeBug,
		Priority:  "P1",
		Status:    domain.IssueStatusOpen,
	}
	if err := writer.WriteIssue("test-backlog", issue); err != nil {
		t.Fatalf("Failed to write initial issue: %v", err)
	}

	updated := &domain.Issue{
		ID:        "issue-replace-123456",
		BacklogID: "test-backlog",
		Name:      "Updated Issue",
		Goal:      strings.Repeat("Updated issue goal with enough characters ", 4),
		IssueType: domain.IssueTypeBug,
		Priority:  "P2",
		Status:    domain.IssueStatusResolved,
	}
	err := writer.ReplaceIssue("test-backlog", updated)
	if err != nil {
		t.Errorf("ReplaceIssue() error = %v", err)
	}

	reader := NewReader(paths)
	result, err := reader.ReadIssue("test-backlog", "issue-replace-123456")
	if err != nil {
		t.Errorf("Failed to read updated issue: %v", err)
	}
	if result.Status != updated.Status {
		t.Errorf("Issue status = %q, want %q", result.Status, updated.Status)
	}
}

// ==========================
// Writer Batch Replace Tests (ReplaceTasks, ReplaceFeatures, ReplaceIssues)
// ==========================

func TestWriterReplaceTasks_Success(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	writer := NewWriter(paths)

	if err := os.MkdirAll(paths.BacklogDirPath("test-backlog"), 0755); err != nil {
		t.Fatalf("Failed to create backlog dir: %v", err)
	}

	// Write initial tasks
	task1 := &domain.Task{
		ID:        "batch-task-1",
		FeatureID: "feature-001",
		SpecID:    "spec-001",
		BacklogID: "test-backlog",
		Name:      "Task 1",
		Goal:      strings.Repeat("Goal 1 ", 20),
		Priority:  "P2",
		Status:    domain.TaskStatusReady,
	}
	task2 := &domain.Task{
		ID:        "batch-task-2",
		FeatureID: "feature-001",
		SpecID:    "spec-001",
		BacklogID: "test-backlog",
		Name:      "Task 2",
		Goal:      strings.Repeat("Goal 2 ", 20),
		Priority:  "P3",
		Status:    domain.TaskStatusPending,
	}
	if err := writer.WriteTask("test-backlog", task1); err != nil {
		t.Fatalf("Failed to write task1: %v", err)
	}
	if err := writer.WriteTask("test-backlog", task2); err != nil {
		t.Fatalf("Failed to write task2: %v", err)
	}

	// Read all tasks
	allTasks := []*domain.Task{task1, task2}

	// Update task1
	tasksToUpdate := map[string]*domain.Task{
		"batch-task-1": {
			ID:        "batch-task-1",
			FeatureID: "feature-001",
			SpecID:    "spec-001",
			BacklogID: "test-backlog",
			Name:      "Updated Task 1",
			Goal:      strings.Repeat("Updated goal 1 ", 20),
			Priority:  "P1",
			Status:    domain.TaskStatusDone,
		},
	}

	err := writer.ReplaceTasks("test-backlog", allTasks, tasksToUpdate)
	if err != nil {
		t.Errorf("ReplaceTasks() error = %v", err)
	}

	// Verify task1 was updated
	reader := NewReader(paths)
	result1, err := reader.ReadTask("test-backlog", "batch-task-1")
	if err != nil {
		t.Fatalf("Failed to read task1: %v", err)
	}
	if result1.Status != domain.TaskStatusDone {
		t.Errorf("Task1 status = %q, want %q", result1.Status, domain.TaskStatusDone)
	}
}

func TestWriterReplaceFeatures_Success(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	writer := NewWriter(paths)

	if err := os.MkdirAll(paths.BacklogDirPath("test-backlog"), 0755); err != nil {
		t.Fatalf("Failed to create backlog dir: %v", err)
	}

	feat1 := &domain.Feature{
		ID:           "batch-feat-1",
		CapabilityID: "cap-001",
		SpecID:       "spec-001",
		BacklogID:    "test-backlog",
		Name:         "Feature 1",
		Goal:         strings.Repeat("Goal 1 ", 20),
		Priority:     "P1",
		Status:       domain.FeatureStatusActive,
	}
	feat2 := &domain.Feature{
		ID:           "batch-feat-2",
		CapabilityID: "cap-002",
		SpecID:       "spec-002",
		BacklogID:    "test-backlog",
		Name:         "Feature 2",
		Goal:         strings.Repeat("Goal 2 ", 20),
		Priority:     "P2",
		Status:       domain.FeatureStatusDraft,
	}
	if err := writer.WriteFeature("test-backlog", feat1); err != nil {
		t.Fatalf("Failed to write feat1: %v", err)
	}
	if err := writer.WriteFeature("test-backlog", feat2); err != nil {
		t.Fatalf("Failed to write feat2: %v", err)
	}

	allFeatures := []*domain.Feature{feat1, feat2}
	featuresToUpdate := map[string]*domain.Feature{
		"batch-feat-1": {
			ID:           "batch-feat-1",
			CapabilityID: "cap-001",
			SpecID:       "spec-001",
			BacklogID:    "test-backlog",
			Name:         "Updated Feature 1",
			Goal:         strings.Repeat("Updated goal 1 ", 20),
			Priority:     "P0",
			Status:       domain.FeatureStatusDone,
		},
	}

	err := writer.ReplaceFeatures("test-backlog", allFeatures, featuresToUpdate)
	if err != nil {
		t.Errorf("ReplaceFeatures() error = %v", err)
	}

	reader := NewReader(paths)
	result1, err := reader.ReadFeature("test-backlog", "batch-feat-1")
	if err != nil {
		t.Fatalf("Failed to read feat1: %v", err)
	}
	if result1.Status != domain.FeatureStatusDone {
		t.Errorf("Feature1 status = %q, want %q", result1.Status, domain.FeatureStatusDone)
	}
}

func TestWriterReplaceIssues_Success(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	writer := NewWriter(paths)

	if err := os.MkdirAll(paths.BacklogDirPath("test-backlog"), 0755); err != nil {
		t.Fatalf("Failed to create backlog dir: %v", err)
	}

	issue1 := &domain.Issue{
		ID:        "batch-issue-1",
		BacklogID: "test-backlog",
		Name:      "Issue 1",
		Goal:      strings.Repeat("Goal 1 with enough characters ", 4),
		IssueType: domain.IssueTypeBug,
		Priority:  "P1",
		Status:    domain.IssueStatusOpen,
	}
	issue2 := &domain.Issue{
		ID:        "batch-issue-2",
		BacklogID: "test-backlog",
		Name:      "Issue 2",
		Goal:      strings.Repeat("Goal 2 with enough characters ", 4),
		IssueType: domain.IssueTypeImprovement,
		Priority:  "P3",
		Status:    domain.IssueStatusReady,
	}
	if err := writer.WriteIssue("test-backlog", issue1); err != nil {
		t.Fatalf("Failed to write issue1: %v", err)
	}
	if err := writer.WriteIssue("test-backlog", issue2); err != nil {
		t.Fatalf("Failed to write issue2: %v", err)
	}

	allIssues := []*domain.Issue{issue1, issue2}
	issuesToUpdate := map[string]*domain.Issue{
		"batch-issue-1": {
			ID:        "batch-issue-1",
			BacklogID: "test-backlog",
			Name:      "Updated Issue 1",
			Goal:      strings.Repeat("Updated goal 1 with enough characters ", 4),
			IssueType: domain.IssueTypeBug,
			Priority:  "P0",
			Status:    domain.IssueStatusResolved,
		},
	}

	err := writer.ReplaceIssues("test-backlog", allIssues, issuesToUpdate)
	if err != nil {
		t.Errorf("ReplaceIssues() error = %v", err)
	}

	reader := NewReader(paths)
	result1, err := reader.ReadIssue("test-backlog", "batch-issue-1")
	if err != nil {
		t.Fatalf("Failed to read issue1: %v", err)
	}
	if result1.Status != domain.IssueStatusResolved {
		t.Errorf("Issue1 status = %q, want %q", result1.Status, domain.IssueStatusResolved)
	}
}

// ==========================
// Writer.WriteFile Tests
// ==========================

func TestWriterWriteFile_Success(t *testing.T) {
	tmpDir, paths := setupTestWorkspace(t)
	defer os.RemoveAll(tmpDir)

	writer := NewWriter(paths)

	content := "Hello, World!"
	err := writer.WriteFile(filepath.Join(tmpDir, "subdir", "test.txt"), content)
	if err != nil {
		t.Errorf("WriteFile() error = %v", err)
	}

	// Verify file was created
	filePath := filepath.Join(tmpDir, "subdir", "test.txt")
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Errorf("Failed to read written file: %v", err)
	}
	if string(data) != content {
		t.Errorf("File content = %q, want %q", string(data), content)
	}
}
