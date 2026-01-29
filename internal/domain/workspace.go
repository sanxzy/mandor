package domain

import "time"

// Workspace represents the root workspace configuration
type Workspace struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	Version       string          `json:"version"`
	SchemaVersion string          `json:"schema_version"`
	CreatedAt     time.Time       `json:"created_at"`
	LastUpdatedAt time.Time       `json:"last_updated_at"`
	CreatedBy     string          `json:"created_by"`
	Config        WorkspaceConfig `json:"config"`
}

// WorkspaceConfig holds workspace-level configuration
type WorkspaceConfig struct {
	DefaultPriority string `json:"default_priority"`
	StrictMode      bool   `json:"strict_mode"`
	DefaultProject  string `json:"default_project,omitempty"`
}

// DefaultWorkspaceConfig returns the default configuration
func DefaultWorkspaceConfig() WorkspaceConfig {
	return WorkspaceConfig{
		DefaultPriority: "P3",
		StrictMode:      false,
	}
}

// ValidatePriority checks if priority is valid (P0-P5)
func ValidatePriority(priority string) bool {
	validPriorities := []string{"P0", "P1", "P2", "P3", "P4", "P5"}
	for _, valid := range validPriorities {
		if priority == valid {
			return true
		}
	}
	return false
}
