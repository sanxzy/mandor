package domain

import (
	"time"

	"mandor/internal/util"
)

const (
	FeatureStatusDraft     = "draft"
	FeatureStatusActive    = "active"
	FeatureStatusDone      = "done"
	FeatureStatusBlocked   = "blocked"
	FeatureStatusCancelled = "cancelled"
)

const (
	FeatureGoalMinLength            = 300
	FeatureGoalMinLengthDevelopment = 2
)

type Feature struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	Name      string    `json:"name"`
	Goal      string    `json:"goal"`
	Scope     string    `json:"scope,omitempty"`
	Priority  string    `json:"priority"`
	Status    string    `json:"status"`
	DependsOn []string  `json:"depends_on,omitempty"`
	Reason    string    `json:"reason,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedBy string    `json:"updated_by"`
}

type FeatureEvent struct {
	Layer   string    `json:"layer"`
	Type    string    `json:"type"`
	ID      string    `json:"id"`
	By      string    `json:"by"`
	Ts      time.Time `json:"ts"`
	Changes []string  `json:"changes,omitempty"`
}

type FeatureCreateInput struct {
	ProjectID string
	Name      string
	Goal      string
	Scope     string
	Priority  string
	DependsOn []string
}

type FeatureListInput struct {
	ProjectID      string
	IncludeDeleted bool
	JSON           bool
}

type FeatureDetailInput struct {
	ProjectID      string
	FeatureID      string
	JSON           bool
	IncludeDeleted bool
}

type FeatureUpdateInput struct {
	ProjectID string
	FeatureID string
	Name      *string
	Goal      *string
	Scope     *string
	Priority  *string
	Status    *string
	Reason    *string
	DependsOn *[]string
	Reopen    bool
	Cancel    bool
	Force     bool
	DryRun    bool
}

type FeatureListItem struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Goal      string `json:"goal,omitempty"`
	Scope     string `json:"scope,omitempty"`
	Priority  string `json:"priority"`
	Status    string `json:"status"`
	DependsOn int    `json:"depends_on_count"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type FeatureListOutput struct {
	Features []FeatureListItem `json:"features"`
	Total    int               `json:"total"`
	Deleted  int               `json:"deleted,omitempty"`
}

type FeatureDetailOutput struct {
	ID        string   `json:"id"`
	ProjectID string   `json:"project_id"`
	Name      string   `json:"name"`
	Goal      string   `json:"goal"`
	Scope     string   `json:"scope,omitempty"`
	Priority  string   `json:"priority"`
	Status    string   `json:"status"`
	DependsOn []string `json:"depends_on"`
	Reason    string   `json:"reason,omitempty"`
	Events    int      `json:"events"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	CreatedBy string   `json:"created_by"`
	UpdatedBy string   `json:"updated_by"`
}

func ValidateFeatureID(id string) bool {
	if len(id) == 0 {
		return false
	}
	if len(id) < 12 {
		return false
	}
	return true
}

func ValidateScope(scope string) bool {
	validScopes := []string{"frontend", "backend", "fullstack", "cli", "desktop", "android", "flutter", "react-native", "ios", "swift", ""}
	for _, s := range validScopes {
		if scope == s {
			return true
		}
	}
	return false
}

func ValidateFeatureStatus(status string) bool {
	validStatuses := []string{FeatureStatusDraft, FeatureStatusActive, FeatureStatusDone, FeatureStatusBlocked, FeatureStatusCancelled}
	for _, s := range validStatuses {
		if status == s {
			return true
		}
	}
	return false
}

func ValidateFeatureGoalLength(goal string) bool {
	minLength := FeatureGoalMinLength
	if util.IsDevelopment() {
		minLength = FeatureGoalMinLengthDevelopment
	}
	return len(goal) >= minLength
}
