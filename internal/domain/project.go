package domain

import (
	"strings"
	"time"

	"mandor/internal/util"
)

const (
	ProjectStatusInitial = "initial"
	ProjectStatusActive  = "active"
	ProjectStatusDone    = "done"
	ProjectStatusDeleted = "deleted"
)

const (
	DependencySameProjectOnly     = "same_project_only"
	DependencyCrossProjectAllowed = "cross_project_allowed"
	DependencyDisabled            = "disabled"
)

const (
	CycleDisallowed = "disallowed"
	CycleAllowed    = "allowed"
)

const (
	GoalMinLength            = 500
	GoalMinLengthDevelopment = 2
)

type Project struct {
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

func ValidateGoalLength(goal string) bool {
	minLength := GoalMinLength
	if util.IsDevelopment() {
		minLength = GoalMinLengthDevelopment
	}
	return len(goal) >= minLength
}

type ProjectEvent struct {
	Layer   string    `json:"layer"`
	Type    string    `json:"type"`
	ID      string    `json:"id"`
	By      string    `json:"by"`
	Ts      time.Time `json:"ts"`
	Changes []string  `json:"changes,omitempty"`
}

type ProjectSchema struct {
	Version string       `json:"version"`
	Schema  string       `json:"$schema"`
	Rules   ProjectRules `json:"rules"`
}

type ProjectRules struct {
	Task     DependencyRule `json:"task"`
	Feature  DependencyRule `json:"feature"`
	Issue    DependencyRule `json:"issue"`
	Priority PriorityConfig `json:"priority"`
}

type DependencyRule struct {
	Dependency string `json:"dependency"`
	Cycle      string `json:"cycle"`
}

type PriorityConfig struct {
	Levels  []string `json:"levels"`
	Default string   `json:"default"`
}

func ValidateProjectID(id string) bool {
	if len(id) == 0 {
		return false
	}
	firstChar := id[0]
	if !((firstChar >= 'a' && firstChar <= 'z') || (firstChar >= 'A' && firstChar <= 'Z')) {
		return false
	}
	for _, c := range id {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_') {
			return false
		}
	}
	return true
}

func ValidateDependencyRule(rule string) bool {
	return rule == DependencySameProjectOnly || rule == DependencyCrossProjectAllowed || rule == DependencyDisabled
}

func ValidateBooleanValue(value string) bool {
	lower := strings.ToLower(value)
	return lower == "true" || lower == "false" || lower == "yes" || lower == "no" || lower == "1" || lower == "0"
}

func ParseBooleanValue(value string) bool {
	lower := strings.ToLower(value)
	return lower == "true" || lower == "yes" || lower == "1"
}

func DefaultProjectSchema(taskDep, featureDep, issueDep string) ProjectSchema {
	if taskDep == "" {
		taskDep = DependencySameProjectOnly
	}
	if featureDep == "" {
		featureDep = DependencyCrossProjectAllowed
	}
	if issueDep == "" {
		issueDep = DependencySameProjectOnly
	}

	return ProjectSchema{
		Version: "mandor.v1",
		Schema:  "https://json-schema.org/draft/2020-12/schema",
		Rules: ProjectRules{
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

type ProjectStats struct {
	Features   EntityStats   `json:"features"`
	Tasks      EntityStats   `json:"tasks"`
	Issues     EntityStats   `json:"issues"`
	Timestamps TimelineStats `json:"timeline"`
}

type EntityStats struct {
	Total        int            `json:"total"`
	ByStatus     map[string]int `json:"by_status"`
	ByType       map[string]int `json:"by_type,omitempty"`
	AvgPriority  string         `json:"avg_priority"`
	BlockedCount int            `json:"blocked_count,omitempty"`
}

type TimelineStats struct {
	OldestCreated time.Time `json:"oldest_created"`
	NewestCreated time.Time `json:"newest_created"`
	DaysActive    int       `json:"days_active"`
}

type ProjectListItem struct {
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

type ProjectListOutput struct {
	Projects []ProjectListItem `json:"projects"`
	Total    int               `json:"total"`
	Deleted  int               `json:"deleted"`
}

type ProjectDetailOutput struct {
	ID        string        `json:"id"`
	Name      string        `json:"name"`
	Goal      string        `json:"goal"`
	Status    string        `json:"status"`
	Strict    bool          `json:"strict"`
	Schema    ProjectSchema `json:"schema"`
	Stats     ProjectStats  `json:"stats"`
	Activity  ActivityInfo  `json:"activity"`
	CreatedAt string        `json:"created_at"`
	UpdatedAt string        `json:"updated_at"`
	CreatedBy string        `json:"created_by"`
	UpdatedBy string        `json:"updated_by"`
}

type ActivityInfo struct {
	TotalEvents  int    `json:"total_events"`
	LastActivity string `json:"last_event_at"`
}

type ProjectCreateInput struct {
	ID         string
	Name       string
	Goal       string
	TaskDep    string
	FeatureDep string
	IssueDep   string
	Strict     bool
}

type ProjectUpdateInput struct {
	ID         string
	Name       *string
	Goal       *string
	TaskDep    *string
	FeatureDep *string
	IssueDep   *string
	Strict     *bool
}

type ProjectDeleteInput struct {
	ID     string
	Hard   bool
	DryRun bool
	Yes    bool
}

type ProjectReopenInput struct {
	ID  string
	Yes bool
}
