package fs

import (
	"os"
	"path/filepath"
)

const (
	MandorDir     = ".mandor"
	BacklogsDir   = "backlogs"
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

// SessionNotesPath returns the path to session-notes.jsonl
func (p *Paths) SessionNotesPath() string {
	return filepath.Join(p.MandorDirPath(), "session-notes.jsonl")
}

// BacklogDirExists checks if a backlog directory exists
func (p *Paths) BacklogDirExists(backlogID string) bool {
	_, err := os.Stat(p.BacklogDirPath(backlogID))
	return err == nil
}

// BriefPath returns the path to brief.md
func (p *Paths) BriefPath(backlogID string) string {
	return filepath.Join(p.BacklogDirPath(backlogID), "brief.md")
}

// SpecsDirPath returns the path to specs directory
func (p *Paths) SpecsDirPath(backlogID string) string {
	return filepath.Join(p.BacklogDirPath(backlogID), "specs")
}

// SpecPath returns the path to a specific spec file
func (p *Paths) SpecPath(backlogID, specID string) string {
	return filepath.Join(p.SpecsDirPath(backlogID), specID+".md")
}

// BlueprintPath returns the path to blueprint.md
func (p *Paths) BlueprintPath(backlogID string) string {
	return filepath.Join(p.BacklogDirPath(backlogID), "blueprint.md")
}
