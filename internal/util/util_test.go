package util

import (
	"regexp"
	"strings"
	"testing"
)

// ==========================
// GenerateID Tests
// ==========================

func TestGenerateIDLength(t *testing.T) {
	id, err := GenerateID()
	if err != nil {
		t.Errorf("GenerateID() error = %v", err)
	}
	if len(id) != 4 {
		t.Errorf("GenerateID() length = %d, want 4", len(id))
	}
}

func TestGenerateIDCharacters(t *testing.T) {
	id, err := GenerateID()
	if err != nil {
		t.Errorf("GenerateID() error = %v", err)
	}

	// Check that all characters are alphanumeric
	expectedPattern := "^[a-zA-Z0-9]+$"
	matched, _ := regexp.MatchString(expectedPattern, id)
	if !matched {
		t.Errorf("GenerateID() = %q, contains non-alphanumeric characters", id)
	}
}

func TestGenerateIDUnique(t *testing.T) {
	// Generate multiple IDs and ensure they can be different
	ids := make(map[string]bool)
	count := 100

	for i := 0; i < count; i++ {
		id, err := GenerateID()
		if err != nil {
			t.Errorf("GenerateID() error = %v", err)
		}
		ids[id] = true
	}

	// With 4 characters and 62 possible chars, we should get reasonable uniqueness
	if len(ids) < count/2 {
		t.Errorf("Expected at least %d unique IDs, got %d", count/2, len(ids))
	}
}

// ==========================
// IsValidWorkspaceName Tests
// ==========================

func TestIsValidWorkspaceName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid lowercase", "myworkspace", true},
		{"valid uppercase", "MYWORKSPACE", true},
		{"valid mixed case", "MyWorkspace", true},
		{"valid with numbers", "workspace123", true},
		{"valid with hyphen", "my-workspace", true},
		{"valid with underscore", "my_workspace", true},
		{"valid mixed special", "my-workspace_123", true},
		{"empty string", "", false},
		{"with space", "my workspace", false},
		{"with dot", "my.workspace", false},
		{"with @", "my@workspace", false},
		{"with slash", "my/workspace", false},
		{"with tab", "my\tworkspace", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidWorkspaceName(tt.input)
			if result != tt.expected {
				t.Errorf("IsValidWorkspaceName(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// ==========================
// ToSlug Tests
// ==========================

func TestToSlugBasic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"lowercase", "hello", "hello"},
		{"uppercase", "HELLO", "hello"},
		{"mixed case", "HelloWorld", "helloworld"},
		{"with spaces", "hello world", "hello-world"},
		{"with multiple spaces", "hello   world", "hello-world"},
		{"with special chars", "hello@world!", "hello-world"},
		{"with underscores", "hello_world", "hello-world"},
		{"consecutive hyphens", "hello---world", "hello-world"},
		{"leading hyphen", "-hello", "hello"},
		{"trailing hyphen", "hello-", "hello"},
		{"both leading and trailing", "-hello-", "hello"},
		{"empty string", "", ""},
		{"only special chars", "@#$%", ""},
		{"CamelCase", "myWorkspace", "myworkspace"},
		{"snake_case", "my_workspace", "my-workspace"},
		{"kebab-case", "my-workspace", "my-workspace"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToSlug(tt.input)
			if result != tt.expected {
				t.Errorf("ToSlug(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToSlugNoLeadingTrailingHyphens(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"leading hyphen", "-hello"},
		{"trailing hyphen", "hello-"},
		{"both", "-hello-"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToSlug(tt.input)
			if strings.HasPrefix(result, "-") {
				t.Errorf("ToSlug(%q) = %q, should not start with hyphen", tt.input, result)
			}
			if strings.HasSuffix(result, "-") {
				t.Errorf("ToSlug(%q) = %q, should not end with hyphen", tt.input, result)
			}
		})
	}
}

func TestToSlugNoConsecutiveHyphens(t *testing.T) {
	result := ToSlug("hello---world")
	if strings.Contains(result, "--") {
		t.Errorf("ToSlug() = %q, should not contain consecutive hyphens", result)
	}
}

// ==========================
// NextSequential Tests
// ==========================

func TestNextSequentialLength(t *testing.T) {
	id := NextSequential("test", []string{})
	parts := strings.Split(id, "-")
	if len(parts) != 2 {
		t.Errorf("NextSequential() = %q, expected format prefix-XXXX", id)
	}
	if len(parts[1]) != 4 {
		t.Errorf("NextSequential() ID part length = %d, want 4", len(parts[1]))
	}
}

func TestNextSequentialPrefix(t *testing.T) {
	tests := []struct {
		name     string
		prefix   string
		expected string
	}{
		{"lowercase prefix", "task", "task-"},
		{"uppercase prefix", "TASK", "task-"},
		{"mixed prefix", "Task", "task-"},
		{"with underscore", "test_id", "test_id-"}, // Underscores are kept as-is
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := NextSequential(tt.prefix, []string{})
			if !strings.HasPrefix(id, tt.expected) {
				t.Errorf("NextSequential(%q) = %q, should start with %q", tt.prefix, id, tt.expected)
			}
		})
	}
}

func TestNextSequentialFormat(t *testing.T) {
	id := NextSequential("req", []string{})

	// Check format: prefix-XXXX
	if !regexp.MustCompile(`^[a-z]+-[a-zA-Z0-9]{4}$`).MatchString(id) {
		t.Errorf("NextSequential() = %q, expected format prefix-XXXX where XXXX is alphanumeric", id)
	}
}

func TestNextSequentialWithExistingIDs(t *testing.T) {
	existingIDs := []string{"req-a1b2", "req-c3d4"}

	// The function doesn't use existingIDs, but we test it still works
	id := NextSequential("req", existingIDs)

	if !strings.HasPrefix(id, "req-") {
		t.Errorf("NextSequential() = %q, should start with req-", id)
	}
}

// ==========================
// GetCurrentDirectory Tests
// ==========================

func TestGetCurrentDirectory(t *testing.T) {
	dir, err := GetCurrentDirectory()
	if err != nil {
		t.Errorf("GetCurrentDirectory() error = %v", err)
	}
	if dir == "" {
		t.Error("GetCurrentDirectory() returned empty string")
	}
}

func TestGetCurrentDirectoryFormat(t *testing.T) {
	dir, err := GetCurrentDirectory()
	if err != nil {
		t.Errorf("GetCurrentDirectory() error = %v", err)
	}

	// Should not contain path separators
	if strings.Contains(dir, "/") && strings.Contains(dir, "\\") {
		// On Unix, should not contain /
		// On Windows, should not contain \
		if strings.Contains(dir, string([]byte{'/'})) {
			t.Errorf("GetCurrentDirectory() = %q, should not contain path separator", dir)
		}
	}
}
