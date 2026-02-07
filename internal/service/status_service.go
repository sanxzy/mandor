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
	Backlogs     []BacklogSummary  `json:"backlogs"`
	Dependencies DependencySummary `json:"dependencies"`
	Totals       TotalStats        `json:"totals"`
}

// BacklogSummary represents a backlog in status output
type BacklogSummary struct {
	ID    string              `json:"id"`
	Name  string              `json:"name,omitempty"`
	Stats domain.BacklogStats `json:"stats"`
}

// DependencySummary represents dependency statistics
type DependencySummary struct {
	CrossBacklogCount int      `json:"cross_backlog_count"`
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
func (s *StatusService) GetWorkspaceStatus(backlogID string) (*WorkspaceStatus, error) {
	ws, err := s.reader.ReadWorkspace()
	if err != nil {
		return nil, err
	}

	status := &WorkspaceStatus{
		Workspace:    ws,
		Backlogs:     []BacklogSummary{},
		Dependencies: DependencySummary{},
		Totals:       TotalStats{},
	}

	// Get backlogs to analyze
	var backlogIDs []string
	if backlogID != "" {
		// Single backlog
		if !s.reader.BacklogExists(backlogID) {
			return nil, domain.NewValidationError(fmt.Sprintf("Backlog not found: %s", backlogID))
		}
		backlogIDs = []string{backlogID}
	} else {
		// All backlogs
		backlogs, err := s.reader.ListBacklogs(false)
		if err != nil {
			return nil, err
		}
		backlogIDs = backlogs
	}

	// Calculate stats for each backlog
	for _, bid := range backlogIDs {
		backlogStatus, err := s.GetBacklogStatus(bid)
		if err != nil {
			// Skip backlogs that can't be read
			continue
		}

		status.Backlogs = append(status.Backlogs, *backlogStatus)

		// Accumulate totals
		status.Totals.Features += backlogStatus.Stats.Features.Total
		status.Totals.Tasks += backlogStatus.Stats.Tasks.Total
		status.Totals.Issues += backlogStatus.Stats.Issues.Total
		status.Totals.Blocked += backlogStatus.Stats.Tasks.BlockedCount
	}

	return status, nil
}

// GetBacklogStatus retrieves detailed status for a single backlog
func (s *StatusService) GetBacklogStatus(backlogID string) (*BacklogSummary, error) {
	if !s.reader.BacklogExists(backlogID) {
		return nil, domain.NewValidationError(fmt.Sprintf("Backlog not found: %s", backlogID))
	}

	metadata, err := s.reader.ReadBacklogMetadata(backlogID)
	if err != nil {
		return nil, err
	}

	stats := domain.BacklogStats{
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
	s.reader.ReadNDJSON(s.paths.BacklogFeaturesPath(backlogID), func(raw []byte) error {
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
	s.reader.ReadNDJSON(s.paths.BacklogTasksPath(backlogID), func(raw []byte) error {
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
	s.reader.ReadNDJSON(s.paths.BacklogIssuesPath(backlogID), func(raw []byte) error {
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

	return &BacklogSummary{
		ID:    backlogID,
		Name:  metadata.Name,
		Stats: stats,
	}, nil
}
