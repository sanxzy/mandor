package service

import (
	"encoding/json"
	"fmt"
	"time"

	"mandor/internal/domain"
	"mandor/internal/fs"
)

// StatusService handles status and statistics operations
type StatusService struct {
	reader *fs.Reader
	paths  *fs.Paths
}

// NewStatusService creates a new status service
func NewStatusService() (*StatusService, error) {
	paths, err := fs.NewPaths()
	if err != nil {
		return nil, err
	}
	return &StatusService{
		reader: fs.NewReader(paths),
		paths:  paths,
	}, nil
}

// WorkspaceStatus represents the overall workspace status
type WorkspaceStatus struct {
	Workspace    *domain.Workspace `json:"workspace"`
	Projects     []ProjectSummary  `json:"projects"`
	Dependencies DependencySummary `json:"dependencies"`
	Totals       TotalStats        `json:"totals"`
}

// ProjectSummary represents a project in status output
type ProjectSummary struct {
	ID    string              `json:"id"`
	Name  string              `json:"name,omitempty"`
	Stats domain.ProjectStats `json:"stats"`
}

// DependencySummary represents dependency statistics
type DependencySummary struct {
	CrossProjectCount int      `json:"cross_project_count"`
	CircularDeps      int      `json:"circular_dependencies"`
	BlockingItems     []string `json:"blocking_items"`
}

// TotalStats represents workspace-wide totals
type TotalStats struct {
	Features int `json:"features"`
	Tasks    int `json:"tasks"`
	Issues   int `json:"issues"`
	Active   int `json:"active"`
	Blocked  int `json:"blocked"`
}

// GetWorkspaceStatus retrieves the complete workspace status
func (s *StatusService) GetWorkspaceStatus(projectID string) (*WorkspaceStatus, error) {
	ws, err := s.reader.ReadWorkspace()
	if err != nil {
		return nil, err
	}

	status := &WorkspaceStatus{
		Workspace:    ws,
		Projects:     []ProjectSummary{},
		Dependencies: DependencySummary{},
		Totals:       TotalStats{},
	}

	// Get projects to analyze
	var projectIDs []string
	if projectID != "" {
		// Single project
		if !s.reader.ProjectExists(projectID) {
			return nil, domain.NewValidationError(fmt.Sprintf("Project not found: %s", projectID))
		}
		projectIDs = []string{projectID}
	} else {
		// All projects
		projects, err := s.reader.ListProjects(false)
		if err != nil {
			return nil, err
		}
		projectIDs = projects
	}

	// Calculate stats for each project
	for _, pid := range projectIDs {
		metadata, err := s.reader.ReadProjectMetadata(pid)
		if err != nil {
			// Skip projects that can't be read
			continue
		}

		// For now, create empty stats (will expand with feature/task/issue counts)
		stats := domain.ProjectStats{
			Features: domain.EntityStats{
				ByStatus: make(map[string]int),
			},
			Tasks: domain.EntityStats{
				ByStatus: make(map[string]int),
			},
			Issues: domain.EntityStats{
				ByStatus: make(map[string]int),
				ByType:   make(map[string]int),
			},
		}

		summary := ProjectSummary{
			ID:    pid,
			Name:  metadata.Name,
			Stats: stats,
		}

		status.Projects = append(status.Projects, summary)
	}

	return status, nil
}

// GetProjectStatus retrieves detailed status for a single project
func (s *StatusService) GetProjectStatus(projectID string) (*ProjectSummary, error) {
	if !s.reader.ProjectExists(projectID) {
		return nil, domain.NewValidationError(fmt.Sprintf("Project not found: %s", projectID))
	}

	metadata, err := s.reader.ReadProjectMetadata(projectID)
	if err != nil {
		return nil, err
	}

	stats := domain.ProjectStats{
		Features: domain.EntityStats{
			ByStatus: make(map[string]int),
		},
		Tasks: domain.EntityStats{
			ByStatus: make(map[string]int),
		},
		Issues: domain.EntityStats{
			ByStatus: make(map[string]int),
			ByType:   make(map[string]int),
		},
		Timestamps: domain.TimelineStats{
			OldestCreated: time.Now(),
			NewestCreated: time.Time{},
		},
	}

	// Read features
	s.reader.ReadNDJSON(s.paths.ProjectFeaturesPath(projectID), func(raw []byte) error {
		var feature map[string]interface{}
		if err := json.Unmarshal(raw, &feature); err != nil {
			return nil
		}

		status, ok := feature["status"].(string)
		if ok {
			stats.Features.ByStatus[status]++
			stats.Features.Total++
		}

		return nil
	})

	// Read tasks
	s.reader.ReadNDJSON(s.paths.ProjectTasksPath(projectID), func(raw []byte) error {
		var task map[string]interface{}
		if err := json.Unmarshal(raw, &task); err != nil {
			return nil
		}

		status, ok := task["status"].(string)
		if ok {
			stats.Tasks.ByStatus[status]++
			stats.Tasks.Total++
			if status == "blocked" {
				stats.Tasks.BlockedCount++
			}
		}

		return nil
	})

	// Read issues
	s.reader.ReadNDJSON(s.paths.ProjectIssuesPath(projectID), func(raw []byte) error {
		var issue map[string]interface{}
		if err := json.Unmarshal(raw, &issue); err != nil {
			return nil
		}

		if status, ok := issue["status"].(string); ok {
			stats.Issues.ByStatus[status]++
			stats.Issues.Total++
		}

		if issueType, ok := issue["type"].(string); ok {
			stats.Issues.ByType[issueType]++
		}

		return nil
	})

	return &ProjectSummary{
		ID:    projectID,
		Name:  metadata.Name,
		Stats: stats,
	}, nil
}
