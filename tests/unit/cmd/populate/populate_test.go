package populate_test

import (
	"bytes"
	"strings"
	"testing"

	"mandor/internal/cmd/populate"
)

func TestNewPopulateCmd(t *testing.T) {
	cmd := populate.NewPopulateCmd()

	if cmd.Use != "populate" {
		t.Errorf("Expected use 'populate', got '%s'", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("Expected Short description to be set")
	}

	if cmd.Long == "" {
		t.Error("Expected Long description to be set")
	}

	if cmd.RunE == nil {
		t.Error("Expected RunE to be set")
	}

	if !cmd.Flags().HasFlags() {
		t.Error("Expected flags to be defined")
	}
}

func TestNewPopulateCmd_HasMarkdownFlag(t *testing.T) {
	cmd := populate.NewPopulateCmd()

	if cmd.Flags().Lookup("markdown") == nil {
		t.Error("Expected 'markdown' flag to exist")
	}

	markdownFlag := cmd.Flags().Lookup("markdown")
	if markdownFlag.Shorthand != "m" {
		t.Errorf("Expected 'm' shorthand flag to exist, got '%s'", markdownFlag.Shorthand)
	}
}

func TestNewPopulateCmd_HasJSONFlag(t *testing.T) {
	cmd := populate.NewPopulateCmd()

	if cmd.Flags().Lookup("json") == nil {
		t.Error("Expected 'json' flag to exist")
	}

	jsonFlag := cmd.Flags().Lookup("json")
	if jsonFlag.Shorthand != "j" {
		t.Errorf("Expected 'j' shorthand flag to exist, got '%s'", jsonFlag.Shorthand)
	}
}

func TestPopulateCmdExecution(t *testing.T) {
	cmd := populate.NewPopulateCmd()

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	if err != nil {
		t.Errorf("Execute returned error: %v", err)
	}

	output := buf.String()

	if output == "" {
		t.Error("Expected non-empty output")
	}

	if !strings.Contains(output, "MANDOR CLI COMMAND REFERENCE") {
		t.Error("Expected output to contain 'MANDOR CLI COMMAND REFERENCE'")
	}

	if !strings.Contains(output, "WORKSPACE") || !strings.Contains(output, "mandor init") {
		t.Error("Expected output to contain workspace commands")
	}

	if !strings.Contains(output, "PROJECT") || !strings.Contains(output, "mandor project") {
		t.Error("Expected output to contain project commands")
	}

	if !strings.Contains(output, "FEATURE") || !strings.Contains(output, "mandor feature") {
		t.Error("Expected output to contain feature commands")
	}

	if !strings.Contains(output, "TASK") || !strings.Contains(output, "mandor task") {
		t.Error("Expected output to contain task commands")
	}

	if !strings.Contains(output, "ISSUE") || !strings.Contains(output, "mandor issue") {
		t.Error("Expected output to contain issue commands")
	}

	if !strings.Contains(output, "BEST PRACTICES") {
		t.Error("Expected output to contain 'BEST PRACTICES'")
	}

	if !strings.Contains(output, "Exit Codes") {
		t.Error("Expected output to contain 'Exit Codes'")
	}

	if !strings.Contains(output, "Priority") {
		t.Error("Expected output to contain priority information")
	}
}

func TestPopulateCmd_ContainsCommandExamples(t *testing.T) {
	cmd := populate.NewPopulateCmd()

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{})

	cmd.Execute()
	output := buf.String()

	expectedCommands := []string{
		"mandor init",
		"mandor status",
		"mandor config",
		"mandor project create",
		"mandor feature create",
		"mandor task create",
		"mandor issue create",
		"mandor completion",
	}

	for _, expectedCmd := range expectedCommands {
		if !strings.Contains(output, expectedCmd) {
			t.Errorf("Expected output to contain '%s'", expectedCmd)
		}
	}
}

func TestPopulateCmd_ContainsBestPractices(t *testing.T) {
	cmd := populate.NewPopulateCmd()

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{})

	cmd.Execute()
	output := buf.String()

	expectedPractices := []string{
		"WORKFLOW DESIGN",
		"FEATURE CREATION",
		"TASK CREATION",
		"ISSUE TRACKING",
		"DEPENDENCY MANAGEMENT",
		"STATUS MANAGEMENT",
	}

	for _, practice := range expectedPractices {
		if !strings.Contains(output, practice) {
			t.Errorf("Expected output to contain '%s'", practice)
		}
	}
}

func TestPopulateCmd_ContainsStatusFlows(t *testing.T) {
	cmd := populate.NewPopulateCmd()

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{})

	cmd.Execute()
	output := buf.String()

	if !strings.Contains(output, "draft ─→ active ─→ done") {
		t.Error("Expected output to contain feature status flow")
	}

	if !strings.Contains(output, "pending ─→ ready ─→ in_progress ─→ done") {
		t.Error("Expected output to contain task status flow")
	}

	if !strings.Contains(output, "open ─→ ready ─→ in_progress ─→ resolved") {
		t.Error("Expected output to contain issue status flow")
	}
}

func TestPopulateCmd_PriorityLevels(t *testing.T) {
	cmd := populate.NewPopulateCmd()

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{})

	cmd.Execute()
	output := buf.String()

	priorityLevels := []string{"P0", "P1", "P2", "P3", "P4", "P5"}
	for _, p := range priorityLevels {
		if !strings.Contains(output, p) {
			t.Errorf("Expected output to contain '%s' priority level", p)
		}
	}
}

func TestPopulateCmd_ExitCodes(t *testing.T) {
	cmd := populate.NewPopulateCmd()

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{})

	cmd.Execute()
	output := buf.String()

	if !strings.Contains(output, "0 = Success") {
		t.Error("Expected output to contain exit code 0")
	}

	if !strings.Contains(output, "1 = System error") {
		t.Error("Expected output to contain exit code 1")
	}

	if !strings.Contains(output, "2 = Validation error") {
		t.Error("Expected output to contain exit code 2")
	}

	if !strings.Contains(output, "3 = Permission error") {
		t.Error("Expected output to contain exit code 3")
	}
}

func TestPopulateCmd_SubCommandStructure(t *testing.T) {
	cmd := populate.NewPopulateCmd()

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{})

	cmd.Execute()
	output := buf.String()

	subCommands := []string{
		"mandor project list",
		"mandor project detail",
		"mandor project update",
		"mandor project delete",
		"mandor project reopen",
		"mandor feature list",
		"mandor feature detail",
		"mandor feature update",
		"mandor task list",
		"mandor task detail",
		"mandor task update",
		"mandor issue list",
		"mandor issue detail",
		"mandor issue update",
	}

	for _, subCmd := range subCommands {
		if !strings.Contains(output, subCmd) {
			t.Errorf("Expected output to contain '%s'", subCmd)
		}
	}
}
