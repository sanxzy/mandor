package cmd_test

import (
	"os"
	"testing"

	"mandor/internal/cmd"
)

func TestCompletionCmd_InvalidShell(t *testing.T) {
	rootCmd := cmd.NewRootCmd()
	completionCmd := cmd.NewCompletionCmd(rootCmd)

	completionCmd.SetArgs([]string{"invalid"})
	err := completionCmd.Execute()
	if err == nil {
		t.Error("Expected error for invalid shell")
	}
}

func TestCompletionCmd_Bash(t *testing.T) {
	rootCmd := cmd.NewRootCmd()
	completionCmd := cmd.NewCompletionCmd(rootCmd)

	completionCmd.SetArgs([]string{"bash"})
	err := completionCmd.Execute()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestCompletionCmd_Zsh(t *testing.T) {
	rootCmd := cmd.NewRootCmd()
	completionCmd := cmd.NewCompletionCmd(rootCmd)

	completionCmd.SetArgs([]string{"zsh"})
	err := completionCmd.Execute()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestCompletionCmd_Fish(t *testing.T) {
	rootCmd := cmd.NewRootCmd()
	completionCmd := cmd.NewCompletionCmd(rootCmd)

	completionCmd.SetArgs([]string{"fish"})
	err := completionCmd.Execute()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestCompletionCmd_NoArgs(t *testing.T) {
	rootCmd := cmd.NewRootCmd()
	completionCmd := cmd.NewCompletionCmd(rootCmd)

	completionCmd.SetArgs([]string{})
	err := completionCmd.Execute()
	if err == nil {
		t.Error("Expected error when no shell specified")
	}
}

func TestCompletionCmd_Help(t *testing.T) {
	rootCmd := cmd.NewRootCmd()
	completionCmd := cmd.NewCompletionCmd(rootCmd)

	completionCmd.SetArgs([]string{"--help"})
	err := completionCmd.Execute()
	if err != nil {
		t.Errorf("Expected no error for --help, got: %v", err)
	}
}

func TestNewCompletionCmd(t *testing.T) {
	rootCmd := cmd.NewRootCmd()
	cmd := cmd.NewCompletionCmd(rootCmd)
	if cmd == nil {
		t.Fatal("Expected command, got nil")
	}
	if cmd.Use != "completion [bash|zsh|fish]" {
		t.Errorf("Expected use 'completion [bash|zsh|fish]', got %q", cmd.Use)
	}
	if !cmd.DisableFlagsInUseLine {
		t.Error("Expected DisableFlagsInUseLine to be true")
	}
}

func TestCompletionCmd_Integration(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("Skipping integration test, set INTEGRATION_TEST=1 to run")
	}
}
