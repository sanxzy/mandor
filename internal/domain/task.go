package domain

import "time"

const (
	TaskStatusPending    = "pending"
	TaskStatusReady      = "ready"
	TaskStatusInProgress = "in_progress"
	TaskStatusBlocked    = "blocked"
	TaskStatusDone       = "done"
	TaskStatusCancelled  = "cancelled"
)

const TaskGoalMinLength = 500

// ReadGates represent execution gates that must be satisfied before task can transition to in_progress
type ReadGates struct {
	IsReadBrief        bool `json:"is_read_brief"`         // Must read Brief before starting task
	IsReadSpec         bool `json:"is_read_spec"`          // Must read Spec before starting task
	IsReadSessionNotes bool `json:"is_read_session_notes"` // Must read session notes before starting task
}

type Task struct {
	ID                  string    `json:"id"`
	FeatureID           string    `json:"feature_id"`
	SpecID              string    `json:"spec_id"` // Reference to Spec (must match Feature's spec_id)
	BacklogID           string    `json:"backlog_id"`
	Name                string    `json:"name"`
	Goal                string    `json:"goal"`
	Priority            string    `json:"priority"`
	Status              string    `json:"status"`
	DependsOn           []string  `json:"depends_on,omitempty"`
	Reason              string    `json:"reason,omitempty"`
	IAEScenarios        []string  `json:"iae_scenarios,omitempty"` // Array of "req-XXXX:scenario-YYYY" references
	ImplementationSteps []string  `json:"implementation_steps,omitempty"`
	TestCases           []string  `json:"test_cases,omitempty"`
	LibraryNeeds        []string  `json:"library_needs,omitempty"`
	ReadGates           ReadGates `json:"read_gates"` // Execution gates
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	CreatedBy           string    `json:"created_by"`
	UpdatedBy           string    `json:"updated_by"`
}

type TaskCreateInput struct {
	FeatureID           string
	SpecID              string // Must match Feature's spec_id
	Name                string
	Goal                string
	IAEScenarios        []string // Array of "req-XXXX:scenario-YYYY" references (pipe-separated in CLI)
	ImplementationSteps []string
	TestCases           []string
	LibraryNeeds        []string
	Priority            string
	DependsOn           []string
}

type TaskListInput struct {
	FeatureID      string
	BacklogID      string
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
	IAEScenarios        *[]string
	ImplementationSteps *[]string
	TestCases           *[]string
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
	BacklogID      string `json:"backlog_id"`
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
	BacklogID           string   `json:"backlog_id"`
	Name                string   `json:"name"`
	Goal                string   `json:"goal"`
	Priority            string   `json:"priority"`
	Status              string   `json:"status"`
	DependsOn           []string `json:"depends_on"`
	Reason              string   `json:"reason,omitempty"`
	ImplementationSteps []string `json:"implementation_steps"`
	TestCases           []string `json:"test_cases"`
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

func ValidateTaskGoalLength(goal string, minLength int) bool {
	return len(goal) >= minLength
}
