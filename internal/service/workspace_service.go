package service

import (
	"fmt"
	"os"
	"time"

	"mandor/internal/domain"
	"mandor/internal/fs"
	"mandor/internal/util"
)

// WorkspaceService handles workspace operations
type WorkspaceService struct {
	reader *fs.Reader
	writer *fs.Writer
	paths  *fs.Paths
}

// NewWorkspaceService creates a new workspace service
func NewWorkspaceService() (*WorkspaceService, error) {
	paths, err := fs.NewPaths()
	if err != nil {
		return nil, err
	}
	return &WorkspaceService{
		reader: fs.NewReader(paths),
		writer: fs.NewWriter(paths),
		paths:  paths,
	}, nil
}

// InitWorkspace initializes a new workspace
func (s *WorkspaceService) InitWorkspace(workspaceName string) (*domain.Workspace, error) {
	// Pre-flight checks
	if s.writer.MandorDirExists() {
		return nil, domain.NewValidationError(
			"Workspace already initialized in this directory.\nUse `mandor status` to view.",
		)
	}

	// Check write permissions
	cwd, err := os.Getwd()
	if err != nil {
		return nil, domain.NewSystemError("Cannot determine current directory", err)
	}

	testFile := fmt.Sprintf("%s/.mandor_test", cwd)
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		os.Remove(testFile)
		if os.IsPermission(err) {
			return nil, domain.NewPermissionError(
				"Permission denied. Cannot initialize workspace here.",
			)
		}
		return nil, domain.NewSystemError("Write permission check failed", err)
	}
	os.Remove(testFile)

	// Generate workspace ID
	id, err := util.GenerateNanoID()
	if err != nil {
		return nil, domain.NewSystemError("Cannot generate workspace ID", err)
	}

	// Determine workspace name
	if workspaceName == "" {
		workspaceName, err = util.GetCurrentDirectory()
		if err != nil {
			return nil, domain.NewSystemError("Cannot determine directory name", err)
		}
	}

	// Validate workspace name
	if !util.IsValidWorkspaceName(workspaceName) {
		return nil, domain.NewValidationError(
			"Invalid workspace name.\nAllowed characters: alphanumeric, hyphens (-), underscores (_)",
		)
	}

	// Get git user
	createdBy := util.GetGitUsername()

	// Create workspace structure
	now := time.Now().UTC()
	ws := &domain.Workspace{
		ID:            id,
		Name:          workspaceName,
		Version:       "mandor.v1",
		SchemaVersion: "mandor.v1",
		CreatedAt:     now,
		LastUpdatedAt: now,
		CreatedBy:     createdBy,
		Config:        domain.DefaultWorkspaceConfig(),
	}

	// Create directories
	if err := s.writer.CreateMandorDir(); err != nil {
		return nil, err
	}

	// Write workspace.json
	if err := s.writer.WriteWorkspace(ws); err != nil {
		return nil, err
	}

	return ws, nil
}

// GetWorkspace retrieves the workspace configuration
func (s *WorkspaceService) GetWorkspace() (*domain.Workspace, error) {
	return s.reader.ReadWorkspace()
}

// UpdateWorkspaceConfig updates a configuration value
func (s *WorkspaceService) UpdateWorkspaceConfig(key string, value interface{}) error {
	ws, err := s.reader.ReadWorkspace()
	if err != nil {
		return err
	}

	// Validate and set value
	switch key {
	case "default_priority":
		strValue, ok := value.(string)
		if !ok {
			return domain.NewValidationError("default_priority must be a string")
		}
		if !domain.ValidatePriority(strValue) {
			return domain.NewValidationError(
				"Invalid value for default_priority.\nUse one of: P0, P1, P2, P3, P4, P5",
			)
		}
		ws.Config.DefaultPriority = strValue

	case "strict_mode":
		boolValue, ok := value.(bool)
		if !ok {
			return domain.NewValidationError("strict_mode must be a boolean")
		}
		ws.Config.StrictMode = boolValue

	default:
		return domain.NewValidationError(
			fmt.Sprintf("Unknown configuration key: %s\n\nAvailable keys:\n  - default_priority\n  - strict_mode", key),
		)
	}

	// Update timestamp
	ws.LastUpdatedAt = time.Now().UTC()

	// Write back
	return s.writer.WriteWorkspace(ws)
}

// GetConfigValue retrieves a single configuration value
func (s *WorkspaceService) GetConfigValue(key string) (interface{}, error) {
	ws, err := s.reader.ReadWorkspace()
	if err != nil {
		return nil, err
	}

	switch key {
	case "default_priority":
		return ws.Config.DefaultPriority, nil
	case "strict_mode":
		return ws.Config.StrictMode, nil
	default:
		return nil, domain.NewValidationError(
			fmt.Sprintf("Unknown configuration key: %s", key),
		)
	}
}
