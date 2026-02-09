package track

import (
	"encoding/json"
	"time"
)

// Global flags for track command
var globalFlags = struct {
	JSON    bool
	CSV     bool
	Tree    bool
	Graph   bool
	Verbose bool
	GroupBy string
}{}

// TrackResponse is the unified response structure for track command
type TrackResponse struct {
	Scope                   string             `json:"scope"`
	ID                      string             `json:"id,omitempty"`
	Name                    string             `json:"name,omitempty"`
	Projects                []ProjectTrackItem `json:"projects,omitempty"`
	Features                []FeatureTrackItem `json:"features,omitempty"`
	Tasks                   []TaskTrackItem    `json:"tasks,omitempty"`
	Issues                  []IssueTrackItem   `json:"issues,omitempty"`
	Summary                 SummaryStats       `json:"summary"`
	RecommendedNextCommands []string           `json:"recommendedNextCommands"`
}

// ProjectTrackItem represents a project in track output
type ProjectTrackItem struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Status       string         `json:"status"`
	CreatedAt    string         `json:"created_at"`
	UpdatedAt    string         `json:"updated_at"`
	Features     int            `json:"features,omitempty"`
	Tasks        int            `json:"tasks,omitempty"`
	Issues       int            `json:"issues,omitempty"`
	ByStatus     map[string]int `json:"by_status,omitempty"`
	BlockedCount int            `json:"blocked_count,omitempty"`
	Description  string         `json:"description,omitempty"`
	BlockedBy    []string       `json:"blocked_by,omitempty"`
	Blocks       []string       `json:"blocks,omitempty"`
	RelatedTo    []string       `json:"related_to,omitempty"`
}

// FeatureTrackItem represents a feature in track output
type FeatureTrackItem struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	Status        string         `json:"status"`
	Priority      string         `json:"priority"`
	Scope         string         `json:"scope,omitempty"`
	CreatedAt     string         `json:"created_at"`
	UpdatedAt     string         `json:"updated_at"`
	Tasks         int            `json:"tasks,omitempty"`
	ByStatus      map[string]int `json:"by_status,omitempty"`
	BlockedCount  int            `json:"blocked_count,omitempty"`
	CompletionPct int            `json:"completion_percent,omitempty"`
	Goal          string         `json:"goal,omitempty"`
	BlockedBy     []string       `json:"blocked_by,omitempty"`
	Blocks        []string       `json:"blocks,omitempty"`
	RelatedTo     []string       `json:"related_to,omitempty"`
}

// TaskTrackItem represents a task in track output
type TaskTrackItem struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	Status              string   `json:"status"`
	Priority            string   `json:"priority"`
	FeatureID           string   `json:"feature_id"`
	CreatedAt           string   `json:"created_at"`
	UpdatedAt           string   `json:"updated_at"`
	Goal                string   `json:"goal,omitempty"`
	ImplementationSteps []string `json:"implementation_steps,omitempty"`
	TestCases           []string `json:"test_cases,omitempty"`
	Derivables          []string `json:"derivables,omitempty"`
	BlockedBy           []string `json:"blocked_by,omitempty"`
	Blocks              []string `json:"blocks,omitempty"`
	RelatedTo           []string `json:"related_to,omitempty"`
}

// IssueTrackItem represents an issue in track output
type IssueTrackItem struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Type        string   `json:"type"`
	Status      string   `json:"status"`
	Priority    string   `json:"priority"`
	BacklogID   string   `json:"project_id"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
	Description string   `json:"description,omitempty"`
	BlockedBy   []string `json:"blocked_by,omitempty"`
	Blocks      []string `json:"blocks,omitempty"`
	RelatedTo   []string `json:"related_to,omitempty"`
}

// SummaryStats provides aggregated statistics
type SummaryStats struct {
	Total             int            `json:"total"`
	ByStatus          map[string]int `json:"by_status"`
	CompletionPercent int            `json:"completion_percent"`
}

// TreeNode represents a node in tree visualization
type TreeNode struct {
	ID        string
	Name      string
	Status    string
	Priority  string
	Type      string // "task", "issue", "feature", "project"
	Children  []*TreeNode
	BlockedBy []string
	Blocks    []string
	Level     int
}

// ParseFlags validates output format flags
func (tr *TrackResponse) ParseFlags() error {
	// Count mutually exclusive output format flags
	formatCount := 0
	if globalFlags.JSON {
		formatCount++
	}
	if globalFlags.CSV {
		formatCount++
	}
	if globalFlags.Tree {
		formatCount++
	}
	if globalFlags.Graph {
		formatCount++
	}

	// Multiple output formats not allowed
	if formatCount > 1 {
		return NewFlagError("Only one output format can be specified (--json, --csv, --tree, --graph)")
	}

	// Validate group-by values
	if globalFlags.GroupBy != "" && globalFlags.GroupBy != "status" && globalFlags.GroupBy != "priority" {
		return NewFlagError("--group-by must be 'status' or 'priority'")
	}

	return nil
}

// FlagError represents a flag validation error
type FlagError struct {
	Message string
}

func (e *FlagError) Error() string {
	return e.Message
}

func NewFlagError(msg string) *FlagError {
	return &FlagError{Message: msg}
}

// Timestamps for JSON marshaling
func timeToString(t time.Time) string {
	return t.Format(time.RFC3339)
}

// UnmarshalJSON parses JSON with relaxed timestamp parsing
func (tr *TrackResponse) UnmarshalJSON(data []byte) error {
	type Alias TrackResponse
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(tr),
	}
	return json.Unmarshal(data, &aux)
}
