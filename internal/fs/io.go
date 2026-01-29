package fs

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"mandor/internal/domain"
)

// Reader reads filesystem resources
type Reader struct {
	paths *Paths
}

// NewReader creates a new filesystem reader
func NewReader(paths *Paths) *Reader {
	return &Reader{paths: paths}
}

// WorkspaceExists checks if workspace is initialized
func (r *Reader) WorkspaceExists() bool {
	_, err := os.Stat(r.paths.WorkspacePath())
	return err == nil
}

// ReadWorkspace reads the workspace.json file
func (r *Reader) ReadWorkspace() (*domain.Workspace, error) {
	data, err := os.ReadFile(r.paths.WorkspacePath())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
		}
		return nil, domain.NewSystemError("Cannot read workspace file", err)
	}

	var ws domain.Workspace
	if err := json.Unmarshal(data, &ws); err != nil {
		return nil, domain.NewSystemError("Cannot parse workspace config. File may be corrupted.", err)
	}
	return &ws, nil
}

// ProjectExists checks if a project directory exists
func (r *Reader) ProjectExists(projectID string) bool {
	_, err := os.Stat(r.paths.ProjectDirPath(projectID))
	return err == nil
}

// ListProjects lists all project directories
func (r *Reader) ListProjects(includeDeleted bool) ([]string, error) {
	projectsDir := r.paths.ProjectsDirPath()

	// Check if projects directory exists
	_, err := os.Stat(projectsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil // Empty workspace
		}
		return nil, domain.NewSystemError("Cannot read projects directory", err)
	}

	entries, err := os.ReadDir(projectsDir)
	if err != nil {
		return nil, domain.NewSystemError("Cannot list projects", err)
	}

	var projects []string
	for _, entry := range entries {
		if entry.IsDir() {
			projects = append(projects, entry.Name())
		}
	}
	return projects, nil
}

// ReadProjectMetadata reads project metadata from project.jsonl
func (r *Reader) ReadProjectMetadata(projectID string) (*domain.Project, error) {
	data, err := os.ReadFile(r.paths.ProjectMetadataPath(projectID))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, domain.NewValidationError(fmt.Sprintf("Project not found: %s", projectID))
		}
		return nil, domain.NewSystemError("Cannot read project metadata", err)
	}

	var project domain.Project
	if err := json.Unmarshal(data, &project); err != nil {
		return nil, domain.NewSystemError("Cannot parse project metadata", err)
	}
	return &project, nil
}

// ReadProjectSchema reads the schema.json file for a project
func (r *Reader) ReadProjectSchema(projectID string) (*domain.ProjectSchema, error) {
	data, err := os.ReadFile(r.paths.ProjectSchemaPath(projectID))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, domain.NewValidationError(fmt.Sprintf("Project schema not found: %s", projectID))
		}
		return nil, domain.NewSystemError("Cannot read project schema", err)
	}

	var schema domain.ProjectSchema
	if err := json.Unmarshal(data, &schema); err != nil {
		return nil, domain.NewSystemError("Cannot parse project schema", err)
	}
	return &schema, nil
}

// CountLines counts the number of lines in a file
func (r *Reader) CountLines(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, domain.NewSystemError("Cannot open file", err)
	}
	defer file.Close()

	count := 0
	buf := make([]byte, 32*1024)
	for {
		n, err := file.Read(buf)
		count++
		if err == io.EOF {
			break
		}
		if err != nil {
			return count, domain.NewSystemError("Cannot read file", err)
		}
		count += n
		for i := 0; i < n; i++ {
			if buf[i] == '\n' {
				count++
			}
		}
	}
	return count, nil
}

// CountEntityLines counts non-empty lines in an entity file
func (r *Reader) CountEntityLines(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, domain.NewSystemError("Cannot open entity file", err)
	}
	defer file.Close()

	count := 0
	buf := make([]byte, 32*1024)
	for {
		n, err := file.Read(buf)
		for i := 0; i < n; i++ {
			if buf[i] == '\n' {
				count++
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return count, domain.NewSystemError("Cannot read entity file", err)
		}
	}
	return count, nil
}

// CountEventLines counts event lines in events.jsonl
func (r *Reader) CountEventLines(projectID string) (int, error) {
	return r.CountEntityLines(r.paths.ProjectEventsPath(projectID))
}

// ReadNDJSON reads NDJSON file and unmarshal each line
func (r *Reader) ReadNDJSON(filepath string, processor func([]byte) error) error {
	file, err := os.Open(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // File doesn't exist yet
		}
		return domain.NewSystemError("Cannot read NDJSON file", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	for {
		var raw json.RawMessage
		if err := decoder.Decode(&raw); err != nil {
			if err == io.EOF {
				break
			}
			return domain.NewSystemError("Cannot parse NDJSON", err)
		}
		if err := processor(raw); err != nil {
			return err
		}
	}
	return nil
}

// Writer writes filesystem resources
type Writer struct {
	paths *Paths
}

// NewWriter creates a new filesystem writer
func NewWriter(paths *Paths) *Writer {
	return &Writer{paths: paths}
}

// CreateMandorDir creates the .mandor directory structure
func (w *Writer) CreateMandorDir() error {
	mandorDir := w.paths.MandorDirPath()
	if err := os.MkdirAll(mandorDir, 0755); err != nil {
		if os.IsPermission(err) {
			return domain.NewPermissionError("Permission denied. Cannot initialize workspace here.")
		}
		return domain.NewSystemError("Cannot create .mandor directory", err)
	}

	projectsDir := w.paths.ProjectsDirPath()
	if err := os.MkdirAll(projectsDir, 0755); err != nil {
		if os.IsPermission(err) {
			return domain.NewPermissionError("Permission denied. Cannot create projects directory.")
		}
		return domain.NewSystemError("Cannot create projects directory", err)
	}

	return nil
}

// WriteWorkspace writes the workspace.json file
func (w *Writer) WriteWorkspace(ws *domain.Workspace) error {
	data, err := json.MarshalIndent(ws, "", "  ")
	if err != nil {
		return domain.NewSystemError("Cannot marshal workspace", err)
	}

	path := w.paths.WorkspacePath()

	// Write to temporary file first for atomic operation
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		if os.IsPermission(err) {
			return domain.NewPermissionError("Permission denied. Cannot write to .mandor/workspace.json.")
		}
		return domain.NewSystemError("Cannot write workspace file", err)
	}

	// Atomic rename
	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return domain.NewSystemError("Cannot save workspace file", err)
	}

	return nil
}

// AppendNDJSON appends a JSON object as a new line to NDJSON file
func (w *Writer) AppendNDJSON(filepath string, obj interface{}) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return domain.NewSystemError("Cannot marshal to JSON", err)
	}

	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsPermission(err) {
			return domain.NewPermissionError("Permission denied. Cannot write to file.")
		}
		return domain.NewSystemError("Cannot open file for writing", err)
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		return domain.NewSystemError("Cannot write to file", err)
	}
	if _, err := file.WriteString("\n"); err != nil {
		return domain.NewSystemError("Cannot write newline", err)
	}

	return nil
}

// WriteJSON writes a JSON file
func (w *Writer) WriteJSON(filePath string, obj interface{}) error {
	// Create parent directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		if os.IsPermission(err) {
			return domain.NewPermissionError("Permission denied. Cannot create directory.")
		}
		return domain.NewSystemError("Cannot create directory", err)
	}

	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return domain.NewSystemError("Cannot marshal to JSON", err)
	}

	tmpPath := filePath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		if os.IsPermission(err) {
			return domain.NewPermissionError("Permission denied. Cannot write file.")
		}
		return domain.NewSystemError("Cannot write file", err)
	}

	if err := os.Rename(tmpPath, filePath); err != nil {
		os.Remove(tmpPath)
		return domain.NewSystemError("Cannot save file", err)
	}

	return nil
}

// MandorDirExists checks if .mandor directory exists
func (w *Writer) MandorDirExists() bool {
	_, err := os.Stat(w.paths.MandorDirPath())
	return err == nil
}

// CreateProjectDir creates the project directory structure
func (w *Writer) CreateProjectDir(projectID string) error {
	projectDir := w.paths.ProjectDirPath(projectID)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		if os.IsPermission(err) {
			return domain.NewPermissionError("Permission denied. Cannot create project directory.")
		}
		return domain.NewSystemError("Cannot create project directory", err)
	}

	entityFiles := []string{"events.jsonl", "features.jsonl", "tasks.jsonl", "issues.jsonl"}
	for _, file := range entityFiles {
		filePath := filepath.Join(projectDir, file)
		if err := os.WriteFile(filePath, []byte(""), 0644); err != nil {
			if os.IsPermission(err) {
				return domain.NewPermissionError(fmt.Sprintf("Permission denied. Cannot create %s.", file))
			}
			return domain.NewSystemError(fmt.Sprintf("Cannot create %s", file), err)
		}
	}

	return nil
}

// WriteProjectMetadata writes project metadata to project.jsonl (atomic)
func (w *Writer) WriteProjectMetadata(projectID string, project *domain.Project) error {
	data, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		return domain.NewSystemError("Cannot marshal project metadata", err)
	}

	path := w.paths.ProjectMetadataPath(projectID)

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		if os.IsPermission(err) {
			return domain.NewPermissionError("Permission denied. Cannot create project directory.")
		}
		return domain.NewSystemError("Cannot create project directory", err)
	}

	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		if os.IsPermission(err) {
			return domain.NewPermissionError("Permission denied. Cannot write to project.jsonl.")
		}
		return domain.NewSystemError("Cannot write project file", err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return domain.NewSystemError("Cannot save project file", err)
	}

	return nil
}

// WriteProjectSchema writes schema.json for a project
func (w *Writer) WriteProjectSchema(projectID string, schema *domain.ProjectSchema) error {
	return w.WriteJSON(w.paths.ProjectSchemaPath(projectID), schema)
}

// AppendProjectEvent appends an event to events.jsonl
func (w *Writer) AppendProjectEvent(projectID string, event *domain.ProjectEvent) error {
	return w.AppendNDJSON(w.paths.ProjectEventsPath(projectID), event)
}

// DeleteProjectDir removes the project directory and all contents
func (w *Writer) DeleteProjectDir(projectID string) error {
	projectDir := w.paths.ProjectDirPath(projectID)
	if err := os.RemoveAll(projectDir); err != nil {
		if os.IsPermission(err) {
			return domain.NewPermissionError("Permission denied. Cannot delete project directory.")
		}
		return domain.NewSystemError("Cannot delete project directory", err)
	}
	return nil
}

// IsDirWritable checks if a directory is writable
func (w *Writer) IsDirWritable(dirPath string) bool {
	testFile := filepath.Join(dirPath, ".write_test")
	defer os.Remove(testFile)

	err := os.WriteFile(testFile, []byte("test"), 0644)
	return err == nil
}

// CheckProjectWritable checks if project files are writable
func (w *Writer) CheckProjectWritable(projectID string) bool {
	testFile := filepath.Join(w.paths.ProjectDirPath(projectID), ".write_test")
	defer os.Remove(testFile)

	err := os.WriteFile(testFile, []byte("test"), 0644)
	return err == nil
}

// ProjectsDirWritable checks if projects directory is writable
func (w *Writer) ProjectsDirWritable() bool {
	return w.IsDirWritable(w.paths.ProjectsDirPath())
}

func (r *Reader) ReadFeature(projectID, featureID string) (*domain.Feature, error) {
	var feature *domain.Feature
	err := r.ReadNDJSON(r.paths.ProjectFeaturesPath(projectID), func(raw []byte) error {
		var f domain.Feature
		if err := json.Unmarshal(raw, &f); err != nil {
			return err
		}
		if f.ID == featureID {
			feature = &f
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if feature == nil {
		return nil, domain.NewValidationError("Feature not found: " + featureID)
	}
	return feature, nil
}

func (w *Writer) AppendFeatureEvent(projectID string, event *domain.FeatureEvent) error {
	return w.AppendNDJSON(w.paths.ProjectEventsPath(projectID), event)
}

func (w *Writer) WriteFeature(projectID string, feature *domain.Feature) error {
	return w.AppendNDJSON(w.paths.ProjectFeaturesPath(projectID), feature)
}

func (w *Writer) ReplaceFeature(projectID string, feature *domain.Feature) error {
	featuresPath := w.paths.ProjectFeaturesPath(projectID)

	var features []*domain.Feature
	reader := NewReader(w.paths)
	err := reader.ReadNDJSON(featuresPath, func(raw []byte) error {
		var f domain.Feature
		if err := json.Unmarshal(raw, &f); err != nil {
			return err
		}
		if f.ID != feature.ID {
			features = append(features, &f)
		}
		return nil
	})
	if err != nil {
		return err
	}

	features = append(features, feature)

	file, err := os.OpenFile(featuresPath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return domain.NewSystemError("Cannot open features file for writing", err)
	}
	defer file.Close()

	for _, f := range features {
		encoder := json.NewEncoder(file)
		if err := encoder.Encode(f); err != nil {
			return domain.NewSystemError("Cannot write feature", err)
		}
	}

	return nil
}

func (r *Reader) ReadTask(projectID, taskID string) (*domain.Task, error) {
	var task *domain.Task
	err := r.ReadNDJSON(r.paths.ProjectTasksPath(projectID), func(raw []byte) error {
		var t domain.Task
		if err := json.Unmarshal(raw, &t); err != nil {
			return err
		}
		if t.ID == taskID {
			task = &t
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, domain.NewValidationError("Task not found: " + taskID)
	}
	return task, nil
}

func (w *Writer) AppendTaskEvent(projectID string, event *domain.TaskEvent) error {
	return w.AppendNDJSON(w.paths.ProjectEventsPath(projectID), event)
}

func (w *Writer) WriteTask(projectID string, task *domain.Task) error {
	return w.AppendNDJSON(w.paths.ProjectTasksPath(projectID), task)
}

func (w *Writer) ReplaceTask(projectID string, task *domain.Task) error {
	tasksPath := w.paths.ProjectTasksPath(projectID)

	var tasks []*domain.Task
	reader := NewReader(w.paths)
	err := reader.ReadNDJSON(tasksPath, func(raw []byte) error {
		var t domain.Task
		if err := json.Unmarshal(raw, &t); err != nil {
			return err
		}
		if t.ID != task.ID {
			tasks = append(tasks, &t)
		}
		return nil
	})
	if err != nil {
		return err
	}

	tasks = append(tasks, task)

	file, err := os.OpenFile(tasksPath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return domain.NewSystemError("Cannot open tasks file for writing", err)
	}
	defer file.Close()

	for _, t := range tasks {
		encoder := json.NewEncoder(file)
		if err := encoder.Encode(t); err != nil {
			return domain.NewSystemError("Cannot write task", err)
		}
	}

	return nil
}

func (r *Reader) ReadIssue(projectID, issueID string) (*domain.Issue, error) {
	var issue *domain.Issue
	err := r.ReadNDJSON(r.paths.ProjectIssuesPath(projectID), func(raw []byte) error {
		var i domain.Issue
		if err := json.Unmarshal(raw, &i); err != nil {
			return err
		}
		if i.ID == issueID {
			issue = &i
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if issue == nil {
		return nil, domain.NewValidationError("Issue not found: " + issueID)
	}
	return issue, nil
}

func (w *Writer) AppendIssueEvent(projectID string, event *domain.IssueEvent) error {
	return w.AppendNDJSON(w.paths.ProjectEventsPath(projectID), event)
}

func (w *Writer) WriteIssue(projectID string, issue *domain.Issue) error {
	return w.AppendNDJSON(w.paths.ProjectIssuesPath(projectID), issue)
}

func (w *Writer) ReplaceIssue(projectID string, issue *domain.Issue) error {
	issuesPath := w.paths.ProjectIssuesPath(projectID)

	var issues []*domain.Issue
	reader := NewReader(w.paths)
	err := reader.ReadNDJSON(issuesPath, func(raw []byte) error {
		var i domain.Issue
		if err := json.Unmarshal(raw, &i); err != nil {
			return err
		}
		if i.ID != issue.ID {
			issues = append(issues, &i)
		}
		return nil
	})
	if err != nil {
		return err
	}

	issues = append(issues, issue)

	file, err := os.OpenFile(issuesPath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return domain.NewSystemError("Cannot open issues file for writing", err)
	}
	defer file.Close()

	for _, i := range issues {
		encoder := json.NewEncoder(file)
		if err := encoder.Encode(i); err != nil {
			return domain.NewSystemError("Cannot write issue", err)
		}
	}

	return nil
}
