package domain

import (
	"testing"
)

func TestValidateProjectID(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		expected bool
	}{
		{"valid simple", "auth", true},
		{"valid with numbers", "auth123", true},
		{"valid with hyphen", "my-project", true},
		{"valid with underscore", "my_project", true},
		{"valid mixed", "Auth-Service_123", true},
		{"invalid starts with number", "123auth", false},
		{"invalid starts with hyphen", "-auth", false},
		{"invalid starts with underscore", "_auth", false},
		{"empty", "", false},
		{"invalid character", "auth@123", false},
		{"invalid space", "auth 123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateProjectID(tt.id)
			if result != tt.expected {
				t.Errorf("ValidateProjectID(%q) = %v, want %v", tt.id, result, tt.expected)
			}
		})
	}
}

func TestValidateGoalLength(t *testing.T) {
	shortGoal := "This is a short goal"
	longGoal := ""
	for i := 0; i < 501; i++ {
		longGoal += "x"
	}

	tests := []struct {
		name     string
		goal     string
		expected bool
	}{
		{"exactly 500", string(make([]byte, 500)), false},
		{"501 chars", longGoal, true},
		{"empty", "", false},
		{"short goal", shortGoal, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateGoalLength(tt.goal)
			if result != tt.expected {
				t.Errorf("ValidateGoalLength(len=%d) = %v, want %v", len(tt.goal), result, tt.expected)
			}
		})
	}
}

func TestValidateDependencyRule(t *testing.T) {
	tests := []struct {
		name     string
		rule     string
		expected bool
	}{
		{"valid same_project_only", "same_project_only", true},
		{"valid cross_project_allowed", "cross_project_allowed", true},
		{"valid disabled", "disabled", true},
		{"invalid", "invalid", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateDependencyRule(tt.rule)
			if result != tt.expected {
				t.Errorf("ValidateDependencyRule(%q) = %v, want %v", tt.rule, result, tt.expected)
			}
		})
	}
}

func TestValidateBooleanValue(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"true lowercase", "true", true},
		{"true uppercase", "TRUE", true},
		{"false lowercase", "false", true},
		{"yes", "yes", true},
		{"no", "no", true},
		{"1", "1", true},
		{"0", "0", true},
		{"invalid", "maybe", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateBooleanValue(tt.value)
			if result != tt.expected {
				t.Errorf("ValidateBooleanValue(%q) = %v, want %v", tt.value, result, tt.expected)
			}
		})
	}
}

func TestParseBooleanValue(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"true", "true", true},
		{"TRUE", "TRUE", true},
		{"yes", "yes", true},
		{"1", "1", true},
		{"false", "false", false},
		{"no", "no", false},
		{"0", "0", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseBooleanValue(tt.value)
			if result != tt.expected {
				t.Errorf("ParseBooleanValue(%q) = %v, want %v", tt.value, result, tt.expected)
			}
		})
	}
}

func TestDefaultProjectSchema(t *testing.T) {
	schema := DefaultProjectSchema("", "", "")

	if schema.Version != "mandor.v1" {
		t.Errorf("Version = %q, want %q", schema.Version, "mandor.v1")
	}

	if schema.Rules.Task.Dependency != "same_project_only" {
		t.Errorf("Task.Dependency = %q, want %q", schema.Rules.Task.Dependency, "same_project_only")
	}

	if schema.Rules.Feature.Dependency != "cross_project_allowed" {
		t.Errorf("Feature.Dependency = %q, want %q", schema.Rules.Feature.Dependency, "cross_project_allowed")
	}

	if schema.Rules.Issue.Dependency != "same_project_only" {
		t.Errorf("Issue.Dependency = %q, want %q", schema.Rules.Issue.Dependency, "same_project_only")
	}

	if schema.Rules.Priority.Default != "P3" {
		t.Errorf("Priority.Default = %q, want %q", schema.Rules.Priority.Default, "P3")
	}
}

func TestDefaultProjectSchemaWithCustomDeps(t *testing.T) {
	schema := DefaultProjectSchema("cross_project_allowed", "same_project_only", "disabled")

	if schema.Rules.Task.Dependency != "cross_project_allowed" {
		t.Errorf("Task.Dependency = %q, want %q", schema.Rules.Task.Dependency, "cross_project_allowed")
	}

	if schema.Rules.Feature.Dependency != "same_project_only" {
		t.Errorf("Feature.Dependency = %q, want %q", schema.Rules.Feature.Dependency, "same_project_only")
	}

	if schema.Rules.Issue.Dependency != "disabled" {
		t.Errorf("Issue.Dependency = %q, want %q", schema.Rules.Issue.Dependency, "disabled")
	}
}
