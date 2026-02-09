package migration

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"mandor/internal/fs"
)

// Migration handles the project-to-backlog directory and content migration
type Migration struct {
	paths *fs.Paths
}

// NewMigration creates a new migration handler
func NewMigration(paths *fs.Paths) *Migration {
	return &Migration{paths: paths}
}

// RunMigration executes the full migration: backup, directory rename, and content update
func (m *Migration) RunMigration() error {
	// Step 1: Check if migration is needed
	projectsDir := m.paths.ProjectsDirPath()
	backlogsDir := m.paths.BacklogsDirPath()

	projectsDirExists, err := m.dirExists(projectsDir)
	if err != nil {
		return fmt.Errorf("failed to check projects directory: %w", err)
	}

	backlogsDirExists, err := m.dirExists(backlogsDir)
	if err != nil {
		return fmt.Errorf("failed to check backlogs directory: %w", err)
	}

	// If backlogs already exist and projects don't, migration already done
	if !projectsDirExists && backlogsDirExists {
		return nil
	}

	// If neither exist, nothing to migrate
	if !projectsDirExists && !backlogsDirExists {
		return nil
	}

	// Step 2: Create backup
	backupDir := m.paths.MandorDirPath() + "/.projects_backup_" + time.Now().Format("20060102_150405")
	if err := m.createBackup(projectsDir, backupDir); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// Step 3: Migrate directory structure
	if err := m.MigrateDirectoryStructure(); err != nil {
		// Restore from backup on failure
		if restoreErr := m.restoreBackup(backupDir, projectsDir); restoreErr != nil {
			return fmt.Errorf("migration failed: %w (and backup restoration failed: %w)", err, restoreErr)
		}
		return fmt.Errorf("migration failed: %w", err)
	}

	// Step 4: Clean up backup (optional - keep for safety)
	// os.RemoveAll(backupDir)

	return nil
}

// MigrateDirectoryStructure renames .mandor/projects to .mandor/backlogs
func (m *Migration) MigrateDirectoryStructure() error {
	projectsDir := m.paths.ProjectsDirPath()
	backlogsDir := m.paths.BacklogsDirPath()

	// Check if projects directory exists
	exists, err := m.dirExists(projectsDir)
	if err != nil {
		return err
	}
	if !exists {
		return nil // Nothing to migrate
	}

	// Check if backlogs directory already exists
	backlogsExists, err := m.dirExists(backlogsDir)
	if err != nil {
		return err
	}
	if backlogsExists {
		return fmt.Errorf("backlogs directory already exists - cannot migrate")
	}

	// List all project directories
	entries, err := os.ReadDir(projectsDir)
	if err != nil {
		return fmt.Errorf("failed to read projects directory: %w", err)
	}

	// Ensure backlogs directory exists
	if err := os.MkdirAll(backlogsDir, 0755); err != nil {
		return fmt.Errorf("failed to create backlogs directory: %w", err)
	}

	// Move each project directory to backlogs
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		projectID := entry.Name()
		oldPath := filepath.Join(projectsDir, projectID)
		newPath := filepath.Join(backlogsDir, projectID)

		if err := os.Rename(oldPath, newPath); err != nil {
			return fmt.Errorf("failed to move %s to backlogs: %w", projectID, err)
		}

		// Migrate file content for each backlog
		if err := m.MigrateFileContent(projectID); err != nil {
			return fmt.Errorf("failed to migrate content for %s: %w", projectID, err)
		}
	}

	// Remove now-empty projects directory
	if err := os.RemoveAll(projectsDir); err != nil {
		return fmt.Errorf("failed to remove projects directory: %w", err)
	}

	return nil
}

// MigrateFileContent updates JSON field names from "project_id" to "backlog_id"
// in all JSONL files within a backlog directory
func (m *Migration) MigrateFileContent(backlogID string) error {
	// Files to migrate
	filesToMigrate := []string{
		m.paths.BacklogMetadataPath(backlogID),
		m.paths.BacklogFeaturesPath(backlogID),
		m.paths.BacklogTasksPath(backlogID),
		m.paths.BacklogIssuesPath(backlogID),
	}

	for _, filePath := range filesToMigrate {
		// Check if file exists
		if _, err := os.Stat(filePath); err != nil {
			if os.IsNotExist(err) {
				continue // File doesn't exist, skip
			}
			return err
		}

		// Migrate file content
		if err := m.migrateJSONLFile(filePath); err != nil {
			return fmt.Errorf("failed to migrate %s: %w", filePath, err)
		}
	}

	// Also migrate schema.json (single JSON file, not JSONL)
	schemaPath := m.paths.BacklogSchemaPath(backlogID)
	if _, err := os.Stat(schemaPath); err == nil {
		if err := m.migrateJSONFile(schemaPath); err != nil {
			return fmt.Errorf("failed to migrate schema %s: %w", schemaPath, err)
		}
	}

	return nil
}

// migrateJSONLFile updates project_id → backlog_id in a JSONL file
func (m *Migration) migrateJSONLFile(filePath string) error {
	// Read the JSONL file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("cannot open file: %w", err)
	}
	defer file.Close()

	// Parse and update each line
	var lines [][]byte
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		// Parse JSON
		var obj map[string]interface{}
		if err := json.Unmarshal(line, &obj); err != nil {
			return fmt.Errorf("cannot parse JSON line: %w", err)
		}

		// Rename project_id to backlog_id
		if projectID, exists := obj["project_id"]; exists {
			obj["backlog_id"] = projectID
			delete(obj, "project_id")
		}

		// Marshal back to JSON
		migratedLine, err := json.Marshal(obj)
		if err != nil {
			return fmt.Errorf("cannot marshal JSON: %w", err)
		}

		lines = append(lines, migratedLine)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// Write to temporary file
	tmpFile := filePath + ".tmp"
	tmpFileHandle, err := os.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("cannot create temp file: %w", err)
	}
	defer tmpFileHandle.Close()

	// Write migrated lines
	for i, line := range lines {
		if _, err := tmpFileHandle.Write(line); err != nil {
			os.Remove(tmpFile)
			return fmt.Errorf("cannot write to temp file: %w", err)
		}
		// Add newline after each line (except we're writing raw bytes, so add \n)
		if i < len(lines)-1 || true { // Always add newline
			if _, err := tmpFileHandle.WriteString("\n"); err != nil {
				os.Remove(tmpFile)
				return fmt.Errorf("cannot write newline: %w", err)
			}
		}
	}

	tmpFileHandle.Close()

	// Atomically replace original file
	if err := os.Rename(tmpFile, filePath); err != nil {
		os.Remove(tmpFile)
		return fmt.Errorf("cannot replace file: %w", err)
	}

	return nil
}

// migrateJSONFile updates project_id → backlog_id in a single JSON file (not JSONL)
func (m *Migration) migrateJSONFile(filePath string) error {
	// Read the JSON file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("cannot read file: %w", err)
	}

	// Parse JSON
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return fmt.Errorf("cannot parse JSON: %w", err)
	}

	// Rename project_id to backlog_id if it exists
	if projectID, exists := obj["project_id"]; exists {
		obj["backlog_id"] = projectID
		delete(obj, "project_id")
	}

	// Marshal back with indentation
	migratedData, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot marshal JSON: %w", err)
	}

	// Write to temporary file
	tmpFile := filePath + ".tmp"
	if err := os.WriteFile(tmpFile, migratedData, 0644); err != nil {
		return fmt.Errorf("cannot write temp file: %w", err)
	}

	// Atomically replace original file
	if err := os.Rename(tmpFile, filePath); err != nil {
		os.Remove(tmpFile)
		return fmt.Errorf("cannot replace file: %w", err)
	}

	return nil
}

// createBackup creates a backup of the projects directory
func (m *Migration) createBackup(source, destination string) error {
	return m.copyDir(source, destination)
}

// restoreBackup restores the projects directory from backup
func (m *Migration) restoreBackup(backupDir, destination string) error {
	// Remove the migrated backlogs directory if it was created
	backlogsDir := m.paths.BacklogsDirPath()
	if _, err := os.Stat(backlogsDir); err == nil {
		os.RemoveAll(backlogsDir)
	}

	// Restore projects from backup
	return m.copyDir(backupDir, destination)
}

// copyDir recursively copies a directory
func (m *Migration) copyDir(source, destination string) error {
	// Ensure destination parent exists
	destParent := filepath.Dir(destination)
	if err := os.MkdirAll(destParent, 0755); err != nil {
		return fmt.Errorf("failed to create destination parent: %w", err)
	}

	// Walk source directory
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate relative path
		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(destination, relPath)

		if info.IsDir() {
			// Create directory
			return os.MkdirAll(destPath, 0755)
		} else {
			// Copy file
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			return os.WriteFile(destPath, data, info.Mode())
		}
	})
}

// dirExists checks if a directory exists
func (m *Migration) dirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return info.IsDir(), nil
}

// HasProjectData checks if there are any projects that need migration
func (m *Migration) HasProjectData() (bool, error) {
	projectsDir := m.paths.ProjectsDirPath()
	exists, err := m.dirExists(projectsDir)
	if err != nil {
		return false, err
	}
	return exists, nil
}
