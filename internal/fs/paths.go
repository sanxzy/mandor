package fs

import (
	"os"
	"path/filepath"
)

const (
	MandorDir     = ".mandor"
	ProjectsDir   = "projects"
	WorkspaceFile = "workspace.json"
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

// ProjectsDirPath returns the path to projects directory
func (p *Paths) ProjectsDirPath() string {
	return filepath.Join(p.MandorDirPath(), ProjectsDir)
}

// ProjectDirPath returns the path to a specific project directory
func (p *Paths) ProjectDirPath(projectID string) string {
	return filepath.Join(p.ProjectsDirPath(), projectID)
}

// ProjectMetadataPath returns the path to project.jsonl
func (p *Paths) ProjectMetadataPath(projectID string) string {
	return filepath.Join(p.ProjectDirPath(projectID), "project.jsonl")
}

// ProjectSchemaPath returns the path to schema.json
func (p *Paths) ProjectSchemaPath(projectID string) string {
	return filepath.Join(p.ProjectDirPath(projectID), "schema.json")
}

// ProjectEventsPath returns the path to events.jsonl (append-only audit trail)
func (p *Paths) ProjectEventsPath(projectID string) string {
	return filepath.Join(p.ProjectDirPath(projectID), "events.jsonl")
}

// ProjectFeaturesPath returns the path to features.jsonl
func (p *Paths) ProjectFeaturesPath(projectID string) string {
	return filepath.Join(p.ProjectDirPath(projectID), "features.jsonl")
}

// ProjectTasksPath returns the path to tasks.jsonl
func (p *Paths) ProjectTasksPath(projectID string) string {
	return filepath.Join(p.ProjectDirPath(projectID), "tasks.jsonl")
}

// ProjectIssuesPath returns the path to issues.jsonl
func (p *Paths) ProjectIssuesPath(projectID string) string {
	return filepath.Join(p.ProjectDirPath(projectID), "issues.jsonl")
}

// ProjectDirExists checks if a project directory exists
func (p *Paths) ProjectDirExists(projectID string) bool {
	_, err := os.Stat(p.ProjectDirPath(projectID))
	return err == nil
}
