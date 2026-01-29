package domain

import (
	"fmt"
)

// ExitCode represents the exit code for the CLI
type ExitCode int

const (
	ExitSuccess         ExitCode = 0 // Success
	ExitSystemError     ExitCode = 1 // System error (permission denied, disk full, etc.)
	ExitValidationError ExitCode = 2 // Validation error (invalid input, already exists, etc.)
	ExitPermissionError ExitCode = 3 // Permission error (cannot write to directory)
)

// MandorError represents an error in the Mandor system
type MandorError struct {
	Code    ExitCode
	Message string
	Cause   error
}

// Error implements the error interface
func (e *MandorError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("Error: %s\n%v", e.Message, e.Cause)
	}
	return fmt.Sprintf("Error: %s", e.Message)
}

// NewSystemError creates a system error
func NewSystemError(message string, cause error) *MandorError {
	return &MandorError{
		Code:    ExitSystemError,
		Message: message,
		Cause:   cause,
	}
}

// NewValidationError creates a validation error
func NewValidationError(message string) *MandorError {
	return &MandorError{
		Code:    ExitValidationError,
		Message: message,
	}
}

// NewPermissionError creates a permission error
func NewPermissionError(message string) *MandorError {
	return &MandorError{
		Code:    ExitPermissionError,
		Message: message,
	}
}
