package domain

import (
	"time"

	"mandor/internal/util"
)

const (
	TaskStatusPending    = "pending"
	TaskStatusReady      = "ready"
	TaskStatusInProgress = "in_progress"
	TaskStatusBlocked    = "blocked"
	TaskStatusDone       = "done"
	TaskStatusCancelled  = "cancelled"
)

const (
	TaskGoalMinLength            = 500
	TaskGoalMinLengthDevelopment = 2
)

type Task struct {
	ID                  string    `json:"id"`
	FeatureID           string    `json:"feature_id"`
	ProjectID           string    `json:"project_id"`
	Name                string    `json:"name"`
	Goal                string    `json:"goal"`
	Priority            string    `json:"priority"`
	Status              string    `json:"status"`
	DependsOn           []string  `json:"depends_on,omitempty"`
	Reason              string    `json:"reason,omitempty"`
	ImplementationSteps []string  `json:"implementation_steps,omitempty"`
	TestCases           []string  `json:"test_cases,omitempty"`
	DerivableFiles      []string  `json:"derivable_files,omitempty"`
	LibraryNeeds        []string  `json:"library_needs,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	CreatedBy           string    `json:"created_by"`
	UpdatedBy           string    `json:"updated_by"`
}

type TaskEvent struct {
	Layer   string    `json:"layer"`
	Type    string    `json:"type"`
	ID      string    `json:"id"`
	By      string    `json:"by"`
	Ts      time.Time `json:"ts"`
	Changes []string  `json:"changes,omitempty"`
}

type TaskCreateInput struct {
	FeatureID           string
	Name                string
	Goal                string
	ImplementationSteps []string
	TestCases           []string
	DerivableFiles      []string
	LibraryNeeds        []string
	Priority            string
	DependsOn           []string
}

type TaskListInput struct {
	FeatureID      string
	ProjectID      string
	Status         string
	Priority       string
	IncludeDeleted bool
	JSON           bool
	Sort           string
	Order          string
}

type TaskDetailInput struct {
	FeatureID      string
	TaskID         string
	JSON           bool
	IncludeDeleted bool
	Events         bool
	Dependencies   bool
	Timestamps     bool
}

type TaskUpdateInput struct {
	FeatureID           string
	TaskID              string
	Name                *string
	Goal                *string
	Priority            *string
	ImplementationSteps *[]string
	TestCases           *[]string
	DerivableFiles      *[]string
	LibraryNeeds        *[]string
	Status              *string
	Reason              *string
	DependsOn           *[]string
	DependsAdd          *[]string
	DependsRemove       *[]string
	Reopen              bool
	Cancel              bool
	Force               bool
	DryRun              bool
}

type TaskListItem struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Status         string `json:"status"`
	Priority       string `json:"priority"`
	FeatureID      string `json:"feature_id"`
	ProjectID      string `json:"project_id"`
	DependsOnCount int    `json:"depends_on_count"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

type TaskListOutput struct {
	Tasks   []TaskListItem `json:"tasks"`
	Total   int            `json:"total"`
	Deleted int            `json:"deleted,omitempty"`
}

type TaskDetailOutput struct {
	ID                  string   `json:"id"`
	FeatureID           string   `json:"feature_id"`
	ProjectID           string   `json:"project_id"`
	Name                string   `json:"name"`
	Goal                string   `json:"goal"`
	Priority            string   `json:"priority"`
	Status              string   `json:"status"`
	DependsOn           []string `json:"depends_on"`
	Reason              string   `json:"reason,omitempty"`
	ImplementationSteps []string `json:"implementation_steps"`
	TestCases           []string `json:"test_cases"`
	DerivableFiles      []string `json:"derivable_files"`
	LibraryNeeds        []string `json:"library_needs"`
	Events              int      `json:"events"`
	CreatedAt           string   `json:"created_at"`
	UpdatedAt           string   `json:"updated_at"`
	CreatedBy           string   `json:"created_by"`
	UpdatedBy           string   `json:"updated_by"`
}

func ValidateTaskID(id string) bool {
	if len(id) == 0 {
		return false
	}
	if len(id) < 12 {
		return false
	}
	return true
}

func ValidateTaskStatus(status string) bool {
	validStatuses := []string{TaskStatusPending, TaskStatusReady, TaskStatusInProgress, TaskStatusBlocked, TaskStatusDone, TaskStatusCancelled}
	for _, s := range validStatuses {
		if status == s {
			return true
		}
	}
	return false
}

func ValidateTaskGoalLength(goal string) bool {
	minLength := TaskGoalMinLength
	if util.IsDevelopment() {
		minLength = TaskGoalMinLengthDevelopment
	}
	return len(goal) >= minLength
}
