package fs

import (
	"os"
	"path/filepath"
)

const (
	MandorDir      = ".mandor"
	BacklogsDir    = "backlogs"
	WorkspaceFile  = "workspace.json"
	ProjectsDir    = "projects" // Deprecated: kept for backward compatibility
)

// Paths manages filesystem paths for the workspace
type Paths struct {
	WorkspaceRoot string
}

// NewPaths creates a new Paths instance for the current working directory
func NewPaths() (*Paths, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return &Paths{WorkspaceRoot: cwd}, nil
}

// NewPathsFromRoot creates a new Paths instance for the specified root directory
func NewPathsFromRoot(root string) (*Paths, error) {
	return &Paths{WorkspaceRoot: root}, nil
}

// MandorDirPath returns the path to .mandor directory
func (p *Paths) MandorDirPath() string {
	return filepath.Join(p.WorkspaceRoot, MandorDir)
}

// WorkspacePath returns the path to workspace.json
func (p *Paths) WorkspacePath() string {
	return filepath.Join(p.MandorDirPath(), WorkspaceFile)
}

// BacklogsDirPath returns the path to backlogs directory
func (p *Paths) BacklogsDirPath() string {
	return filepath.Join(p.MandorDirPath(), BacklogsDir)
}

// BacklogDirPath returns the path to a specific backlog directory
func (p *Paths) BacklogDirPath(backlogID string) string {
	return filepath.Join(p.BacklogsDirPath(), backlogID)
}

// BacklogMetadataPath returns the path to backlog.jsonl
func (p *Paths) BacklogMetadataPath(backlogID string) string {
	return filepath.Join(p.BacklogDirPath(backlogID), "backlog.jsonl")
}

// BacklogSchemaPath returns the path to schema.json
func (p *Paths) BacklogSchemaPath(backlogID string) string {
	return filepath.Join(p.BacklogDirPath(backlogID), "schema.json")
}

// BacklogFeaturesPath returns the path to features.jsonl
func (p *Paths) BacklogFeaturesPath(backlogID string) string {
	return filepath.Join(p.BacklogDirPath(backlogID), "features.jsonl")
}

// BacklogTasksPath returns the path to tasks.jsonl
func (p *Paths) BacklogTasksPath(backlogID string) string {
	return filepath.Join(p.BacklogDirPath(backlogID), "tasks.jsonl")
}

// BacklogIssuesPath returns the path to issues.jsonl
func (p *Paths) BacklogIssuesPath(backlogID string) string {
	return filepath.Join(p.BacklogDirPath(backlogID), "issues.jsonl")
}

// Deprecated: Use BacklogsDirPath instead
func (p *Paths) ProjectsDirPath() string {
	return filepath.Join(p.MandorDirPath(), ProjectsDir)
}

// Deprecated: Use BacklogDirPath instead
func (p *Paths) ProjectDirPath(projectID string) string {
	return filepath.Join(p.ProjectsDirPath(), projectID)
}

// Deprecated: Use BacklogMetadataPath instead
func (p *Paths) ProjectMetadataPath(projectID string) string {
	return filepath.Join(p.ProjectDirPath(projectID), "project.jsonl")
}

// Deprecated: Use BacklogSchemaPath instead
func (p *Paths) ProjectSchemaPath(projectID string) string {
	return filepath.Join(p.ProjectDirPath(projectID), "schema.json")
}

// Deprecated: Use BacklogFeaturesPath instead
func (p *Paths) ProjectFeaturesPath(projectID string) string {
	return filepath.Join(p.ProjectDirPath(projectID), "features.jsonl")
}

// Deprecated: Use BacklogTasksPath instead
func (p *Paths) ProjectTasksPath(projectID string) string {
	return filepath.Join(p.ProjectDirPath(projectID), "tasks.jsonl")
}

// Deprecated: Use BacklogIssuesPath instead
func (p *Paths) ProjectIssuesPath(projectID string) string {
	return filepath.Join(p.ProjectDirPath(projectID), "issues.jsonl")
}

// SessionNotesPath returns the path to session-notes.jsonl
func (p *Paths) SessionNotesPath() string {
	return filepath.Join(p.MandorDirPath(), "session-notes.jsonl")
}

// BacklogDirExists checks if a backlog directory exists
func (p *Paths) BacklogDirExists(backlogID string) bool {
	_, err := os.Stat(p.BacklogDirPath(backlogID))
	return err == nil
}

// Deprecated: Use BacklogDirExists instead
func (p *Paths) ProjectDirExists(projectID string) bool {
	_, err := os.Stat(p.ProjectDirPath(projectID))
	return err == nil
}
