package domain

import "time"

const (
	IssueStatusOpen       = "open"
	IssueStatusReady      = "ready"
	IssueStatusInProgress = "in_progress"
	IssueStatusBlocked    = "blocked"
	IssueStatusResolved   = "resolved"
	IssueStatusWontFix    = "wontfix"
	IssueStatusCancelled  = "cancelled"
)

const (
	IssueTypeBug         = "bug"
	IssueTypeImprovement = "improvement"
	IssueTypeDebt        = "debt"
	IssueTypeSecurity    = "security"
	IssueTypePerformance = "performance"
)

type Issue struct {
	ID                  string    `json:"id"`
	ProjectID           string    `json:"project_id"`
	Name                string    `json:"name"`
	Goal                string    `json:"goal,omitempty"`
	IssueType           string    `json:"issue_type"`
	Priority            string    `json:"priority"`
	Status              string    `json:"status"`
	DependsOn           []string  `json:"depends_on,omitempty"`
	Reason              string    `json:"reason,omitempty"`
	AffectedFiles       []string  `json:"affected_files,omitempty"`
	AffectedTests       []string  `json:"affected_tests,omitempty"`
	ImplementationSteps []string  `json:"implementation_steps,omitempty"`
	LibraryNeeds        []string  `json:"library_needs,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
	LastUpdatedAt       time.Time `json:"last_updated_at"`
	CreatedBy           string    `json:"created_by"`
	LastUpdatedBy       string    `json:"last_updated_by"`
}

type IssueEvent struct {
	Layer   string    `json:"layer"`
	Type    string    `json:"type"`
	ID      string    `json:"id"`
	By      string    `json:"by"`
	Ts      time.Time `json:"ts"`
	Changes []string  `json:"changes,omitempty"`
}

type IssueCreateInput struct {
	ProjectID           string
	Name                string
	Goal                string
	IssueType           string
	Priority            string
	DependsOn           []string
	AffectedFiles       []string
	AffectedTests       []string
	ImplementationSteps []string
	LibraryNeeds        []string
}

type IssueListInput struct {
	ProjectID      string
	IssueType      string
	Status         string
	Priority       string
	IncludeDeleted bool
	JSON           bool
	Sort           string
	Order          string
}

type IssueDetailInput struct {
	ProjectID      string
	IssueID        string
	JSON           bool
	IncludeDeleted bool
	Events         bool
	Dependencies   bool
	Timestamps     bool
}

type IssueUpdateInput struct {
	ProjectID           string
	IssueID             string
	Name                *string
	Goal                *string
	IssueType           *string
	Priority            *string
	Status              *string
	Reason              *string
	DependsOn           *[]string
	DependsAdd          *[]string
	DependsRemove       *[]string
	AffectedFiles       *[]string
	AffectedTests       *[]string
	ImplementationSteps *[]string
	LibraryNeeds        *[]string
	Start               bool
	Resolve             bool
	WontFix             bool
	Reopen              bool
	Cancel              bool
	Force               bool
	DryRun              bool
}

type IssueListItem struct {
	ID                       string `json:"id"`
	Name                     string `json:"name"`
	IssueType                string `json:"issue_type"`
	Status                   string `json:"status"`
	Priority                 string `json:"priority"`
	ProjectID                string `json:"project_id"`
	DependsOnCount           int    `json:"depends_on_count"`
	AffectedFilesCount       int    `json:"affected_files_count"`
	AffectedTestsCount       int    `json:"affected_tests_count"`
	ImplementationStepsCount int    `json:"implementation_steps_count"`
	LibraryNeedsCount        int    `json:"library_needs_count"`
	CreatedAt                string `json:"created_at"`
	LastUpdatedAt            string `json:"last_updated_at"`
}

type IssueListOutput struct {
	Issues  []IssueListItem `json:"issues"`
	Total   int             `json:"total"`
	Deleted int             `json:"deleted,omitempty"`
}

type IssueDetailOutput struct {
	ID                  string   `json:"id"`
	ProjectID           string   `json:"project_id"`
	Name                string   `json:"name"`
	Goal                string   `json:"goal,omitempty"`
	IssueType           string   `json:"issue_type"`
	Priority            string   `json:"priority"`
	Status              string   `json:"status"`
	DependsOn           []string `json:"depends_on"`
	Reason              string   `json:"reason,omitempty"`
	AffectedFiles       []string `json:"affected_files"`
	AffectedTests       []string `json:"affected_tests"`
	ImplementationSteps []string `json:"implementation_steps"`
	LibraryNeeds        []string `json:"library_needs"`
	Events              int      `json:"events"`
	CreatedAt           string   `json:"created_at"`
	LastUpdatedAt       string   `json:"last_updated_at"`
	CreatedBy           string   `json:"created_by"`
	LastUpdatedBy       string   `json:"last_updated_by"`
}

func ValidateIssueID(id string) bool {
	if len(id) == 0 {
		return false
	}
	if len(id) < 12 {
		return false
	}
	return true
}

func ValidateIssueStatus(status string) bool {
	validStatuses := []string{IssueStatusOpen, IssueStatusReady, IssueStatusInProgress, IssueStatusBlocked, IssueStatusResolved, IssueStatusWontFix, IssueStatusCancelled}
	for _, s := range validStatuses {
		if status == s {
			return true
		}
	}
	return false
}

func ValidateIssueType(issueType string) bool {
	validTypes := []string{IssueTypeBug, IssueTypeImprovement, IssueTypeDebt, IssueTypeSecurity, IssueTypePerformance}
	for _, t := range validTypes {
		if issueType == t {
			return true
		}
	}
	return false
}

func IsIssueTerminalStatus(status string) bool {
	return status == IssueStatusResolved || status == IssueStatusWontFix || status == IssueStatusCancelled
}
