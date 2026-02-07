package service

import (
	"testing"

	"mandor/internal/fs"
)

// TestScenarioA_BriefChangeBlocking tests change governance for brief modifications
func TestScenarioA_BriefChangeBlocking(t *testing.T) {
	tests := []struct {
		name            string
		changedFields   map[string]interface{}
		expectedBlock   string
		expectedActions []string
	}{
		{
			name: "Capability change is BLOCKING",
			changedFields: map[string]interface{}{
				"capabilities": true,
			},
			expectedBlock: "BLOCKING",
			expectedActions: []string{
				"Regenerate all affected Specs",
			},
		},
		{
			name: "Tech-stack change is BLOCKING",
			changedFields: map[string]interface{}{
				"tech-stack": true,
			},
			expectedBlock: "BLOCKING",
			expectedActions: []string{
				"Regenerate Blueprint",
			},
		},
		{
			name: "Timeline change is BLOCKING",
			changedFields: map[string]interface{}{
				"timeline": true,
			},
			expectedBlock: "BLOCKING",
			expectedActions: []string{
				"Revalidate all task effort estimates",
			},
		},
		{
			name: "Why change is NON_BLOCKING",
			changedFields: map[string]interface{}{
				"why": true,
			},
			expectedBlock:   "NON_BLOCKING",
			expectedActions: []string{},
		},
		{
			name: "Business-rationale change is NON_BLOCKING",
			changedFields: map[string]interface{}{
				"business-rationale": true,
			},
			expectedBlock:   "NON_BLOCKING",
			expectedActions: []string{},
		},
		{
			name: "Priority change is NON_BLOCKING",
			changedFields: map[string]interface{}{
				"priority": true,
			},
			expectedBlock:   "NON_BLOCKING",
			expectedActions: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paths, _ := fs.NewPaths()
			svc := NewChangeGovernanceService(paths)

			analysis, err := svc.ValidateBriefChangeBlocking("test-backlog", "test-brief", tt.changedFields)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if analysis.BlockingStatus != tt.expectedBlock {
				t.Errorf("expected BlockingStatus %s, got %s", tt.expectedBlock, analysis.BlockingStatus)
			}

			if tt.expectedBlock == "BLOCKING" && len(analysis.RequiredActions) == 0 {
				t.Errorf("expected required actions for BLOCKING status, got none")
			}
		})
	}
}

// TestScenarioB_SpecChangeBlocking tests change governance for spec modifications
func TestScenarioB_SpecChangeBlocking(t *testing.T) {
	tests := []struct {
		name            string
		changedFields   map[string]interface{}
		expectedBlock   string
		expectedActions []string
	}{
		{
			name: "Acceptance-criteria change is BLOCKING",
			changedFields: map[string]interface{}{
				"acceptance-criteria": true,
			},
			expectedBlock: "BLOCKING",
			expectedActions: []string{
				"Regenerate all dependent Tasks",
			},
		},
		{
			name: "Requirements change is BLOCKING",
			changedFields: map[string]interface{}{
				"requirements": true,
			},
			expectedBlock: "BLOCKING",
			expectedActions: []string{
				"Regenerate all dependent Tasks",
			},
		},
		{
			name: "IAE-scenarios change is BLOCKING",
			changedFields: map[string]interface{}{
				"iae-scenarios": true,
			},
			expectedBlock: "BLOCKING",
			expectedActions: []string{
				"Regenerate all dependent Tasks",
			},
		},
		{
			name: "Input-validation change is BLOCKING",
			changedFields: map[string]interface{}{
				"input-validation": true,
			},
			expectedBlock: "BLOCKING",
			expectedActions: []string{
				"Regenerate affected Tasks",
			},
		},
		{
			name: "Error-handling change is BLOCKING",
			changedFields: map[string]interface{}{
				"error-handling": true,
			},
			expectedBlock: "BLOCKING",
			expectedActions: []string{
				"Regenerate affected Tasks",
			},
		},
		{
			name: "Summary change is NON_BLOCKING",
			changedFields: map[string]interface{}{
				"summary": true,
			},
			expectedBlock:   "NON_BLOCKING",
			expectedActions: []string{},
		},
		{
			name: "Testing-strategy-notes change is NON_BLOCKING",
			changedFields: map[string]interface{}{
				"testing-strategy-notes": true,
			},
			expectedBlock:   "NON_BLOCKING",
			expectedActions: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paths, _ := fs.NewPaths()
			svc := NewChangeGovernanceService(paths)

			analysis, err := svc.ValidateSpecChangeBlocking("test-backlog", "test-spec", tt.changedFields)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if analysis.BlockingStatus != tt.expectedBlock {
				t.Errorf("expected BlockingStatus %s, got %s", tt.expectedBlock, analysis.BlockingStatus)
			}

			if tt.expectedBlock == "BLOCKING" && len(analysis.RequiredActions) == 0 {
				t.Errorf("expected required actions for BLOCKING status, got none")
			}
		})
	}
}

// TestScenarioC_BlueprintChangeBlocking tests change governance for blueprint modifications
func TestScenarioC_BlueprintChangeBlocking(t *testing.T) {
	tests := []struct {
		name            string
		changedFields   map[string]interface{}
		expectedBlock   string
		expectedActions []string
	}{
		{
			name: "Tech-stack change is BLOCKING",
			changedFields: map[string]interface{}{
				"tech-stack": true,
			},
			expectedBlock: "BLOCKING",
			expectedActions: []string{
				"Regenerate all dependent Tasks",
			},
		},
		{
			name: "Deployment-strategy change is BLOCKING",
			changedFields: map[string]interface{}{
				"deployment-strategy": true,
			},
			expectedBlock: "BLOCKING",
			expectedActions: []string{
				"Regenerate all dependent Tasks",
			},
		},
		{
			name: "Architecture-pattern change is BLOCKING",
			changedFields: map[string]interface{}{
				"architecture-pattern": true,
			},
			expectedBlock: "BLOCKING",
			expectedActions: []string{
				"Regenerate all dependent Tasks",
			},
		},
		{
			name: "Scalability change is BLOCKING",
			changedFields: map[string]interface{}{
				"scalability": true,
			},
			expectedBlock: "BLOCKING",
			expectedActions: []string{
				"Revalidate task acceptance criteria",
			},
		},
		{
			name: "Resilience change is BLOCKING",
			changedFields: map[string]interface{}{
				"resilience": true,
			},
			expectedBlock: "BLOCKING",
			expectedActions: []string{
				"Revalidate task acceptance criteria",
			},
		},
		{
			name: "Performance-targets change is BLOCKING",
			changedFields: map[string]interface{}{
				"performance-targets": true,
			},
			expectedBlock: "BLOCKING",
			expectedActions: []string{
				"Revalidate task acceptance criteria",
			},
		},
		{
			name: "Monitoring change is BLOCKING",
			changedFields: map[string]interface{}{
				"monitoring": true,
			},
			expectedBlock: "BLOCKING",
			expectedActions: []string{
				"Regenerate operations/monitoring tasks",
			},
		},
		{
			name: "Alerting change is BLOCKING",
			changedFields: map[string]interface{}{
				"alerting": true,
			},
			expectedBlock: "BLOCKING",
			expectedActions: []string{
				"Regenerate operations/monitoring tasks",
			},
		},
		{
			name: "Operational-requirements change is BLOCKING",
			changedFields: map[string]interface{}{
				"operational-requirements": true,
			},
			expectedBlock: "BLOCKING",
			expectedActions: []string{
				"Regenerate operations/monitoring tasks",
			},
		},
		{
			name: "Cost-analysis change is NON_BLOCKING",
			changedFields: map[string]interface{}{
				"cost-analysis": true,
			},
			expectedBlock:   "NON_BLOCKING",
			expectedActions: []string{},
		},
		{
			name: "1-year-projection change is NON_BLOCKING",
			changedFields: map[string]interface{}{
				"1-year-projection": true,
			},
			expectedBlock:   "NON_BLOCKING",
			expectedActions: []string{},
		},
		{
			name: "Rationale change is NON_BLOCKING",
			changedFields: map[string]interface{}{
				"rationale": true,
			},
			expectedBlock:   "NON_BLOCKING",
			expectedActions: []string{},
		},
		{
			name: "Decision-reversibility change is NON_BLOCKING",
			changedFields: map[string]interface{}{
				"decision-reversibility": true,
			},
			expectedBlock:   "NON_BLOCKING",
			expectedActions: []string{},
		},
		{
			name: "Assumptions change is NON_BLOCKING",
			changedFields: map[string]interface{}{
				"assumptions": true,
			},
			expectedBlock:   "NON_BLOCKING",
			expectedActions: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paths, _ := fs.NewPaths()
			svc := NewChangeGovernanceService(paths)

			analysis, err := svc.ValidateBlueprintChangeBlocking("test-backlog", "test-blueprint", tt.changedFields)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if analysis.BlockingStatus != tt.expectedBlock {
				t.Errorf("expected BlockingStatus %s, got %s", tt.expectedBlock, analysis.BlockingStatus)
			}

			if tt.expectedBlock == "BLOCKING" && len(analysis.RequiredActions) == 0 {
				t.Errorf("expected required actions for BLOCKING status, got none")
			}
		})
	}
}

// TestChangeImpactAnalysisStructure validates that change analysis has required fields
func TestChangeImpactAnalysisStructure(t *testing.T) {
	paths, _ := fs.NewPaths()
	svc := NewChangeGovernanceService(paths)

	analysis, err := svc.ValidateBriefChangeBlocking("test-backlog", "test-brief", map[string]interface{}{"why": true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Validate required fields
	if analysis.ChangeID == "" {
		t.Error("ChangeID must not be empty")
	}
	if analysis.ChangeType != "brief_modification" {
		t.Error("ChangeType must be set correctly")
	}
	if analysis.EntityID != "test-brief" {
		t.Error("EntityID must match input")
	}
	if analysis.BacklogID != "test-backlog" {
		t.Error("BacklogID must match input")
	}
	if analysis.Timestamp == "" {
		t.Error("Timestamp must be set")
	}
	if analysis.Status != "pending_validation" {
		t.Error("Status must be pending_validation on creation")
	}
	if analysis.ValidationDeadline == "" {
		t.Error("ValidationDeadline must be set")
	}
	if analysis.BlockingStatus == "" {
		t.Error("BlockingStatus must be set")
	}
}

// TestMultipleFieldChanges validates behavior when multiple fields change
func TestMultipleFieldChanges(t *testing.T) {
	tests := []struct {
		name          string
		scenario      string
		changedFields map[string]interface{}
		expectedBlock string
	}{
		{
			name:     "Brief: Multiple changes with one blocking field",
			scenario: "brief",
			changedFields: map[string]interface{}{
				"why":          true,
				"capabilities": true,
				"priority":     true,
			},
			expectedBlock: "BLOCKING",
		},
		{
			name:     "Brief: Multiple non-blocking changes",
			scenario: "brief",
			changedFields: map[string]interface{}{
				"why":                true,
				"business-rationale": true,
				"priority":           true,
			},
			expectedBlock: "NON_BLOCKING",
		},
		{
			name:     "Spec: Multiple changes with one blocking field",
			scenario: "spec",
			changedFields: map[string]interface{}{
				"summary":             true,
				"acceptance-criteria": true,
			},
			expectedBlock: "BLOCKING",
		},
		{
			name:     "Blueprint: Multiple non-blocking changes",
			scenario: "blueprint",
			changedFields: map[string]interface{}{
				"cost-analysis":          true,
				"decision-reversibility": true,
				"assumptions":            true,
			},
			expectedBlock: "NON_BLOCKING",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paths, _ := fs.NewPaths()
			svc := NewChangeGovernanceService(paths)

			var analysis *ChangeImpactAnalysis
			var err error

			switch tt.scenario {
			case "brief":
				analysis, err = svc.ValidateBriefChangeBlocking("test-backlog", "test-brief", tt.changedFields)
			case "spec":
				analysis, err = svc.ValidateSpecChangeBlocking("test-backlog", "test-spec", tt.changedFields)
			case "blueprint":
				analysis, err = svc.ValidateBlueprintChangeBlocking("test-backlog", "test-bp", tt.changedFields)
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if analysis.BlockingStatus != tt.expectedBlock {
				t.Errorf("expected %s, got %s", tt.expectedBlock, analysis.BlockingStatus)
			}
		})
	}
}
