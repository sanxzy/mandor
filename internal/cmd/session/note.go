package session

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"mandor/internal/domain"
	"mandor/internal/fs"
)

// SessionNote represents a single session note entry
type SessionNote struct {
	Timestamp string `json:"timestamp"`
	Note      string `json:"note"`
}

const defaultNoteLimit = 50

// NewNoteCmd creates the session note command
func NewNoteCmd() *cobra.Command {
	var readFlag bool
	var offsetFlag int

	cmd := &cobra.Command{
		Use:   "note [text]",
		Short: "Add or read session notes",
		Long: fmt.Sprintf(`Add a timestamped note about work completed, or read recent notes.

Notes are stored in .mandor/session-notes.jsonl as NDJSON (one JSON object per line).
This provides a lightweight way for AI agents to track progress across sessions.

By default, the last %d notes are shown. Use --offset to show more or fewer.`, defaultNoteLimit),
		Example: `  # Add a note
  mandor session note "Completed v0.4.4 release and testing"

  # Read last 50 notes (default)
  mandor session note --read

  # Read last 100 notes
  mandor session note --read --offset 100`,
		RunE: func(cmd *cobra.Command, args []string) error {
			paths, err := fs.NewPaths()
			if err != nil {
				return domain.NewSystemError("Cannot initialize paths", err)
			}

			// Check workspace is initialized
			reader := fs.NewReader(paths)
			if !reader.WorkspaceExists() {
				return domain.NewValidationError("Workspace not initialized. Run `mandor init` first.")
			}

			if readFlag {
				return readNotes(paths, offsetFlag)
			}

			// Add note mode - require text argument
			if len(args) == 0 {
				return domain.NewValidationError("Note text required. Usage: mandor session note \"your note here\"")
			}

			noteText := args[0]
			return addNote(paths, noteText)
		},
	}

	cmd.Flags().BoolVarP(&readFlag, "read", "r", false, "Read recent notes instead of adding")
	cmd.Flags().IntVarP(&offsetFlag, "offset", "o", defaultNoteLimit, "Number of notes to show")

	return cmd
}

// addNote appends a new note to the session notes file
func addNote(paths *fs.Paths, noteText string) error {
	note := SessionNote{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Note:      noteText,
	}

	data, err := json.Marshal(note)
	if err != nil {
		return domain.NewSystemError("Cannot marshal note", err)
	}

	notesPath := paths.SessionNotesPath()

	file, err := os.OpenFile(notesPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsPermission(err) {
			return domain.NewPermissionError("Permission denied. Cannot write to session notes.")
		}
		return domain.NewSystemError("Cannot open session notes file", err)
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		return domain.NewSystemError("Cannot write note", err)
	}
	if _, err := file.WriteString("\n"); err != nil {
		return domain.NewSystemError("Cannot write newline", err)
	}

	fmt.Printf("✓ Note added: %s\n", noteText)
	return nil
}

// readNotes reads and displays the last N notes
func readNotes(paths *fs.Paths, limit int) error {
	notesPath := paths.SessionNotesPath()

	file, err := os.Open(notesPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No session notes yet.")
			return nil
		}
		return domain.NewSystemError("Cannot read session notes", err)
	}
	defer file.Close()

	// Read all notes into a slice
	var notes []SessionNote
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var note SessionNote
		if err := json.Unmarshal([]byte(line), &note); err != nil {
			// Skip malformed lines
			continue
		}
		notes = append(notes, note)
	}

	if err := scanner.Err(); err != nil {
		return domain.NewSystemError("Cannot scan session notes", err)
	}

	if len(notes) == 0 {
		fmt.Println("No session notes yet.")
		return nil
	}

	// Show last N notes (or all if less than limit)
	startIdx := 0
	if len(notes) > limit {
		startIdx = len(notes) - limit
	}

	fmt.Println("Recent session notes:")
	fmt.Println()

	for i := startIdx; i < len(notes); i++ {
		note := notes[i]
		timestamp, err := time.Parse(time.RFC3339, note.Timestamp)
		if err != nil {
			timestamp = time.Time{}
		}
		fmt.Printf("  [%s] %s\n", timestamp.Format("2006-01-02 15:04"), note.Note)
	}

	if len(notes) > limit {
		fmt.Printf("\n  ... and %d older notes\n", len(notes)-limit)
	}

	return nil
}
