package feature_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"mandor/internal/domain"
)

func setupTestFeatureWorkspace(t *testing.T) string {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "mandor-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	mandorDir := filepath.Join(tmpDir, ".mandor")
	projectsDir := filepath.Join(mandorDir, "projects")

	if err := os.MkdirAll(projectsDir, 0755); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to create .mandor directory: %v", err)
	}

	workspacePath := filepath.Join(mandorDir, "workspace.json")
	wsData := []byte("{\"version\":\"v1.0.0\"}")
	if err := os.WriteFile(workspacePath, wsData, 0644); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to write workspace.json: %v", err)
	}

	return tmpDir
}

func writeTestProjectForCmd(t *testing.T, tmpDir, projectID string) {
	t.Helper()

	projectDir := filepath.Join(tmpDir, ".mandor", "projects", projectID)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatalf("Failed to create project dir: %v", err)
	}

	projectPath := filepath.Join(projectDir, "project.jsonl")
	projectData := []byte("{\"id\":\"" + projectID + "\",\"name\":\"Test Project\",\"goal\":\"test goal\",\"status\":\"active\",\"strict\":false}")
	if err := os.WriteFile(projectPath, projectData, 0644); err != nil {
		t.Fatalf("Failed to write project.jsonl: %v", err)
	}

	schemaPath := filepath.Join(projectDir, "schema.json")
	schemaData := []byte("{\"version\":\"mandor.v1\",\"schema\":\"https://json-schema.org/draft/2020-12/schema\",\"rules\":{\"task\":{\"dependency\":\"same_project_only\",\"cycle\":\"disallowed\"},\"feature\":{\"dependency\":\"cross_project_allowed\",\"cycle\":\"disallowed\"},\"issue\":{\"dependency\":\"same_project_only\",\"cycle\":\"disallowed\"},\"priority\":{\"levels\":[\"P0\",\"P1\",\"P2\",\"P3\",\"P4\",\"P5\"],\"default\":\"P3\"}}}")
	if err := os.WriteFile(schemaPath, schemaData, 0644); err != nil {
		t.Fatalf("Failed to write schema.json: %v", err)
	}

	featuresPath := filepath.Join(projectDir, "features.jsonl")
	if err := os.WriteFile(featuresPath, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create features.jsonl: %v", err)
	}

	eventsPath := filepath.Join(projectDir, "events.jsonl")
	if err := os.WriteFile(eventsPath, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create events.jsonl: %v", err)
	}
}

func writeTestFeatureForCmd(t *testing.T, tmpDir, projectID, featureID string, status string) {
	t.Helper()

	projectDir := filepath.Join(tmpDir, ".mandor", "projects", projectID)
	featuresPath := filepath.Join(projectDir, "features.jsonl")

	feature := &domain.Feature{
		ID:        featureID,
		ProjectID: projectID,
		Name:      "Test Feature",
		Goal:      "Test goal",
		Scope:     "fullstack",
		Priority:  "P3",
		Status:    status,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		CreatedBy: "testuser",
		UpdatedBy: "testuser",
	}

	file, err := os.OpenFile(featuresPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("Failed to open features.jsonl: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(feature)
	if err != nil {
		t.Fatalf("Failed to write feature: %v", err)
	}
}
