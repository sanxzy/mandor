package domain

import (
	"strings"
	"time"
)

const (
	BacklogStatusInitial = "initial"
	BacklogStatusActive  = "active"
	BacklogStatusDone    = "done"
	BacklogStatusDeleted = "deleted"
)

const (
	DependencySameBacklogOnly     = "same_backlog_only"
	DependencyCrossBacklogAllowed = "cross_backlog_allowed"
	DependencyDisabled            = "disabled"
)

const (
	CycleDisallowed = "disallowed"
	CycleAllowed    = "allowed"
)

const GoalMinLength = 500

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

func ValidateGoalLength(goal string, minLength int) bool {
	return len(goal) >= minLength
}

type BacklogSchema struct {
	Version string       `json:"version"`
	Schema  string       `json:"$schema"`
	Rules   BacklogRules `json:"rules"`
}

type BacklogRules struct {
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

func ValidateBacklogID(id string) bool {
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
	return rule == DependencySameBacklogOnly || rule == DependencyCrossBacklogAllowed || rule == DependencyDisabled
}

func ValidateBooleanValue(value string) bool {
	lower := strings.ToLower(value)
	return lower == "true" || lower == "false" || lower == "yes" || lower == "no" || lower == "1" || lower == "0"
}

func ParseBooleanValue(value string) bool {
	lower := strings.ToLower(value)
	return lower == "true" || lower == "yes" || lower == "1"
}

func DefaultBacklogSchema(taskDep, featureDep, issueDep string) BacklogSchema {
	if taskDep == "" {
		taskDep = DependencySameBacklogOnly
	}
	if featureDep == "" {
		featureDep = DependencyCrossBacklogAllowed
	}
	if issueDep == "" {
		issueDep = DependencySameBacklogOnly
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

type BacklogStats struct {
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
	ID        string        `json:"id"`
	Name      string        `json:"name"`
	Goal      string        `json:"goal"`
	Status    string        `json:"status"`
	Strict    bool          `json:"strict"`
	Schema    BacklogSchema `json:"schema"`
	Stats     BacklogStats  `json:"stats"`
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
