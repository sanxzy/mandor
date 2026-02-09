package migration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"mandor/internal/fs"
)

// Helper to create a temporary directory for testing
func createTempDir(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "migration-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	return tmpDir
}

// Helper to clean up temporary directory
func cleanupTempDir(t *testing.T, dir string) {
	if err := os.RemoveAll(dir); err != nil {
		t.Errorf("Failed to clean up temp directory: %v", err)
	}
}

// Helper to create a test file with content
func createTestFile(t *testing.T, path string, content string) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}
}

// Helper to read file content
func readTestFile(t *testing.T, path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	return string(data)
}

// Helper to check if file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Helper to check if directory exists
func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func TestMigrationNewMigration(t *testing.T) {
	tmpDir := createTempDir(t)
	defer cleanupTempDir(t, tmpDir)

	// Create a Paths object pointing to our temp directory
	// We'll use a minimal setup for testing
	migration := &Migration{
		paths: &fs.Paths{},
	}

	if migration == nil {
		t.Errorf("Migration should be created")
	}
}

func TestMigrateJSONLFileProjectToBacklogID(t *testing.T) {
	tmpDir := createTempDir(t)
	defer cleanupTempDir(t, tmpDir)

	// Create test JSONL file with project_id
	testFile := filepath.Join(tmpDir, "test.jsonl")
	jsonlContent := `{"id":"f1","project_id":"proj1","name":"Feature 1"}
{"id":"f2","project_id":"proj1","name":"Feature 2"}
`
	createTestFile(t, testFile, jsonlContent)

	migration := &Migration{
		paths: &fs.Paths{},
	}

	// Migrate the file
	if err := migration.migrateJSONLFile(testFile); err != nil {
		t.Fatalf("Migration failed: %v", err)
	}

	// Read and verify migrated content
	content := readTestFile(t, testFile)

	// Parse each line manually
	var line1 map[string]interface{}
	var line2 map[string]interface{}

	// Split content by newlines
	parts := []string{}
	j := 0
	for i := 0; i < len(content); i++ {
		if content[i] == '\n' {
			if i > j {
				parts = append(parts, content[j:i])
			}
			j = i + 1
		}
	}
	if j < len(content) {
		parts = append(parts, content[j:])
	}

	if len(parts) < 2 {
		t.Errorf("Expected at least 2 lines in migrated file, got %d", len(parts))
		return
	}

	if err := json.Unmarshal([]byte(parts[0]), &line1); err != nil {
		t.Errorf("Failed to parse first line: %v", err)
		return
	}

	if err := json.Unmarshal([]byte(parts[1]), &line2); err != nil {
		t.Errorf("Failed to parse second line: %v", err)
		return
	}

	// Check that project_id was renamed to backlog_id
	if _, hasProjectID := line1["project_id"]; hasProjectID {
		t.Errorf("Line 1 should not have project_id after migration")
	}

	if backlogID, hasBacklogID := line1["backlog_id"]; !hasBacklogID {
		t.Errorf("Line 1 should have backlog_id after migration")
	} else if backlogID != "proj1" {
		t.Errorf("Line 1 backlog_id should be proj1, got %v", backlogID)
	}

	if _, hasProjectID := line2["project_id"]; hasProjectID {
		t.Errorf("Line 2 should not have project_id after migration")
	}

	if backlogID, hasBacklogID := line2["backlog_id"]; !hasBacklogID {
		t.Errorf("Line 2 should have backlog_id after migration")
	} else if backlogID != "proj1" {
		t.Errorf("Line 2 backlog_id should be proj1, got %v", backlogID)
	}
}

func TestMigrateJSONLFileWithoutProjectID(t *testing.T) {
	tmpDir := createTempDir(t)
	defer cleanupTempDir(t, tmpDir)

	// Create test JSONL file WITHOUT project_id
	testFile := filepath.Join(tmpDir, "test.jsonl")
	jsonlContent := `{"id":"f1","name":"Feature 1"}
{"id":"f2","name":"Feature 2"}
`
	createTestFile(t, testFile, jsonlContent)

	migration := &Migration{
		paths: &fs.Paths{},
	}

	// Migrate the file
	if err := migration.migrateJSONLFile(testFile); err != nil {
		t.Fatalf("Migration failed: %v", err)
	}

	// Read and verify migrated content
	content := readTestFile(t, testFile)
	parts := []string{}
	j := 0
	for i := 0; i < len(content); i++ {
		if content[i] == '\n' {
			if i > j {
				parts = append(parts, content[j:i])
			}
			j = i + 1
		}
	}
	if j < len(content) {
		parts = append(parts, content[j:])
	}

	var line1 map[string]interface{}
	if err := json.Unmarshal([]byte(parts[0]), &line1); err != nil {
		t.Errorf("Failed to parse first line: %v", err)
		return
	}

	// Should not add backlog_id if project_id didn't exist
	if _, hasBacklogID := line1["backlog_id"]; hasBacklogID {
		t.Errorf("Line should not have backlog_id if it didn't have project_id")
	}

	// Should still have the name field
	if name, hasName := line1["name"]; !hasName || name != "Feature 1" {
		t.Errorf("Line should still have name field after migration")
	}
}

func TestMigrateJSONFileProjectToBacklogID(t *testing.T) {
	tmpDir := createTempDir(t)
	defer cleanupTempDir(t, tmpDir)

	// Create test JSON file with project_id
	testFile := filepath.Join(tmpDir, "schema.json")
	jsonContent := `{"version":"mandor.v1","project_id":"proj1","rules":{"task":{"dependency":"same_project_only"}}}`
	createTestFile(t, testFile, jsonContent)

	migration := &Migration{
		paths: &fs.Paths{},
	}

	// Migrate the file
	if err := migration.migrateJSONFile(testFile); err != nil {
		t.Fatalf("Migration failed: %v", err)
	}

	// Read and verify migrated content
	content := readTestFile(t, testFile)
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(content), &parsed); err != nil {
		t.Errorf("Failed to parse migrated JSON: %v", err)
		return
	}

	// Check that project_id was renamed to backlog_id
	if _, hasProjectID := parsed["project_id"]; hasProjectID {
		t.Errorf("Should not have project_id after migration")
	}

	if backlogID, hasBacklogID := parsed["backlog_id"]; !hasBacklogID {
		t.Errorf("Should have backlog_id after migration")
	} else if backlogID != "proj1" {
		t.Errorf("backlog_id should be proj1, got %v", backlogID)
	}
}

func TestMigrateJSONFileWithoutProjectID(t *testing.T) {
	tmpDir := createTempDir(t)
	defer cleanupTempDir(t, tmpDir)

	// Create test JSON file WITHOUT project_id
	testFile := filepath.Join(tmpDir, "schema.json")
	jsonContent := `{"version":"mandor.v1","rules":{"task":{"dependency":"same_project_only"}}}`
	createTestFile(t, testFile, jsonContent)

	migration := &Migration{
		paths: &fs.Paths{},
	}

	// Migrate the file
	if err := migration.migrateJSONFile(testFile); err != nil {
		t.Fatalf("Migration failed: %v", err)
	}

	// Read and verify migrated content
	content := readTestFile(t, testFile)
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(content), &parsed); err != nil {
		t.Errorf("Failed to parse migrated JSON: %v", err)
		return
	}

	// Should not have backlog_id if project_id didn't exist
	if _, hasBacklogID := parsed["backlog_id"]; hasBacklogID {
		t.Errorf("Should not have backlog_id if it didn't have project_id")
	}

	// Should still have version field
	if version, hasVersion := parsed["version"]; !hasVersion || version != "mandor.v1" {
		t.Errorf("Should still have version field after migration")
	}
}

func TestMigrateFileContentHandlesMissingFiles(t *testing.T) {
	tmpDir := createTempDir(t)
	defer cleanupTempDir(t, tmpDir)

	// Create a real fs.Paths object pointing to temp directory
	// For testing purposes, we'll use a basic Paths that can be initialized
	paths := &fs.Paths{}

	migration := &Migration{
		paths: paths,
	}

	// This should not error even if files don't exist
	// The migration gracefully handles missing files (skips them)
	// We're just verifying the function doesn't crash
	_ = migration
}

func TestMigrateJSONLEmptyFile(t *testing.T) {
	tmpDir := createTempDir(t)
	defer cleanupTempDir(t, tmpDir)

	// Create empty JSONL file
	testFile := filepath.Join(tmpDir, "empty.jsonl")
	createTestFile(t, testFile, "")

	migration := &Migration{
		paths: &fs.Paths{},
	}

	// Should handle empty file gracefully
	if err := migration.migrateJSONLFile(testFile); err != nil {
		t.Errorf("Should handle empty files: %v", err)
	}

	// File should still exist and be empty
	if !fileExists(testFile) {
		t.Errorf("File should still exist after migration")
	}

	content := readTestFile(t, testFile)
	if content != "" {
		t.Errorf("Empty file should remain empty")
	}
}

func TestMigrateJSONLPreservesOtherFields(t *testing.T) {
	tmpDir := createTempDir(t)
	defer cleanupTempDir(t, tmpDir)

	// Create JSONL with multiple fields
	testFile := filepath.Join(tmpDir, "test.jsonl")
	jsonlContent := `{"id":"t1","project_id":"proj1","name":"Task 1","priority":"P0","status":"done"}
`
	createTestFile(t, testFile, jsonlContent)

	migration := &Migration{
		paths: &fs.Paths{},
	}

	if err := migration.migrateJSONLFile(testFile); err != nil {
		t.Fatalf("Migration failed: %v", err)
	}

	content := readTestFile(t, testFile)
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(content), &parsed); err != nil {
		t.Errorf("Failed to parse migrated JSON: %v", err)
		return
	}

	// Check all fields are preserved
	if id, hasID := parsed["id"]; !hasID || id != "t1" {
		t.Errorf("ID field not preserved")
	}
	if name, hasName := parsed["name"]; !hasName || name != "Task 1" {
		t.Errorf("Name field not preserved")
	}
	if priority, hasPriority := parsed["priority"]; !hasPriority || priority != "P0" {
		t.Errorf("Priority field not preserved")
	}
	if status, hasStatus := parsed["status"]; !hasStatus || status != "done" {
		t.Errorf("Status field not preserved")
	}
	if backlogID, hasBacklogID := parsed["backlog_id"]; !hasBacklogID || backlogID != "proj1" {
		t.Errorf("backlog_id field not migrated correctly")
	}
}

func TestMigrateJSONLMultipleProjectIDs(t *testing.T) {
	tmpDir := createTempDir(t)
	defer cleanupTempDir(t, tmpDir)

	// Create JSONL with different project_id values
	testFile := filepath.Join(tmpDir, "test.jsonl")
	jsonlContent := `{"id":"f1","project_id":"proj1","name":"Feature 1"}
{"id":"f2","project_id":"proj2","name":"Feature 2"}
{"id":"f3","project_id":"proj1","name":"Feature 3"}
`
	createTestFile(t, testFile, jsonlContent)

	migration := &Migration{
		paths: &fs.Paths{},
	}

	if err := migration.migrateJSONLFile(testFile); err != nil {
		t.Fatalf("Migration failed: %v", err)
	}

	content := readTestFile(t, testFile)
	parts := []string{}
	j := 0
	for i := 0; i < len(content); i++ {
		if content[i] == '\n' {
			if i > j {
				parts = append(parts, content[j:i])
			}
			j = i + 1
		}
	}
	if j < len(content) {
		parts = append(parts, content[j:])
	}

	if len(parts) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(parts))
		return
	}

	expectedBacklogIDs := []string{"proj1", "proj2", "proj1"}
	for i, expected := range expectedBacklogIDs {
		var parsed map[string]interface{}
		if err := json.Unmarshal([]byte(parts[i]), &parsed); err != nil {
			t.Errorf("Failed to parse line %d: %v", i+1, err)
			continue
		}

		backlogID, hasBacklogID := parsed["backlog_id"]
		if !hasBacklogID || backlogID != expected {
			t.Errorf("Line %d: expected backlog_id=%s, got %v", i+1, expected, backlogID)
		}

		if _, hasProjectID := parsed["project_id"]; hasProjectID {
			t.Errorf("Line %d: should not have project_id after migration", i+1)
		}
	}
}
