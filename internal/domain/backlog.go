package domain

import (
	"time"
)

const (
	BacklogStatusInitial = "initial"
	BacklogStatusActive  = "active"
	BacklogStatusDone    = "done"
	BacklogStatusDeleted = "deleted"
)

// Note: Dependency rule constants and helper functions are defined in project.go
// They are shared between Project and Backlog for backward compatibility

// Type aliases for shared schema types
type BacklogSchema = ProjectSchema
type BacklogRules = ProjectRules

type Backlog struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Goal      string    `json:"goal"`
	Status    string    `json:"status"`
	Strict    bool      `json:"strict"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedBy string    `json:"updated_by"`
}

// Use type aliases for shared types to avoid duplication
type BacklogStats = ProjectStats

// ValidateBacklogID - same as ValidateProjectID
func ValidateBacklogID(id string) bool {
	return ValidateProjectID(id)
}

func DefaultBacklogSchema(taskDep, featureDep, issueDep string) BacklogSchema {
	if taskDep == "" {
		taskDep = DependencySameProjectOnly // uses project constants which work for both
	}
	if featureDep == "" {
		featureDep = DependencyCrossProjectAllowed
	}
	if issueDep == "" {
		issueDep = DependencySameProjectOnly
	}

	return BacklogSchema{
		Version: "mandor.v1",
		Schema:  "https://json-schema.org/draft/2020-12/schema",
		Rules: BacklogRules{
			Task: DependencyRule{
				Dependency: taskDep,
				Cycle:      CycleDisallowed,
			},
			Feature: DependencyRule{
				Dependency: featureDep,
				Cycle:      CycleDisallowed,
			},
			Issue: DependencyRule{
				Dependency: issueDep,
				Cycle:      CycleDisallowed,
			},
			Priority: PriorityConfig{
				Levels:  []string{"P0", "P1", "P2", "P3", "P4", "P5"},
				Default: "P3",
			},
		},
	}
}

type BacklogListItem struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Goal      string `json:"goal,omitempty"`
	Status    string `json:"status"`
	Features  int    `json:"features"`
	Tasks     int    `json:"tasks"`
	Issues    int    `json:"issues"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type BacklogListOutput struct {
	Backlogs []BacklogListItem `json:"backlogs"`
	Total    int               `json:"total"`
	Deleted  int               `json:"deleted"`
}

type BacklogDetailOutput struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Goal      string       `json:"goal"`
	Status    string       `json:"status"`
	Strict    bool         `json:"strict"`
	Schema    BacklogSchema `json:"schema"`
	Stats     BacklogStats `json:"stats"`
	Activity  ActivityInfo `json:"activity"`
	CreatedAt string       `json:"created_at"`
	UpdatedAt string       `json:"updated_at"`
	CreatedBy string       `json:"created_by"`
	UpdatedBy string       `json:"updated_by"`
}

type BacklogCreateInput struct {
	ID         string
	Name       string
	Goal       string
	TaskDep    string
	FeatureDep string
	IssueDep   string
	Strict     bool
}

type BacklogUpdateInput struct {
	ID         string
	Name       *string
	Goal       *string
	TaskDep    *string
	FeatureDep *string
	IssueDep   *string
	Strict     *bool
}

type BacklogDeleteInput struct {
	ID     string
	Hard   bool
	DryRun bool
	Yes    bool
}

type BacklogReopenInput struct {
	ID  string
	Yes bool
}
