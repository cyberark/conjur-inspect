// Package checks defines all of the possible Conjur Inspect checks that can
// be run.
package checks

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/checks/sanitize"
	"github.com/cyberark/conjur-inspect/pkg/log"
)

// userHomeDirFunc is a mockable version of os.UserHomeDir for testing
var userHomeDirFunc = os.UserHomeDir

// CommandHistory collects recent command history from the host machine
type CommandHistory struct{}

// Describe provides a textual description of what this check gathers info on
func (*CommandHistory) Describe() string {
	return "Command History"
}

// Run performs the command history collection
func (ch *CommandHistory) Run(runContext *check.RunContext) []check.Result {
	homeDir, err := userHomeDirFunc()
	if err != nil {
		log.Warn("failed to determine home directory: %s", err)
		return []check.Result{{
			Title:   "Command History",
			Status:  check.StatusInfo,
			Message: "Unable to determine home directory",
		}}
	}

	// Try to read history files in order of preference
	historyPaths := []string{
		filepath.Join(homeDir, ".zsh_history"),
		filepath.Join(homeDir, ".bash_history"),
	}

	var historyContent strings.Builder
	historyFound := false

	for _, historyPath := range historyPaths {
		content, err := readAndLimitHistoryFile(historyPath, 100)
		if err == nil && len(content) > 0 {
			historyContent.WriteString(content)
			historyFound = true
			break
		}
		if err != nil && !os.IsNotExist(err) {
			log.Warn("error reading history file %s: %s", historyPath, err)
		}
	}

	// If no history files were found or readable, return INFO result
	if !historyFound {
		return []check.Result{{
			Title:   "Command History",
			Status:  check.StatusInfo,
			Message: "No command history files found or accessible",
		}}
	}

	// Sanitize the history to redact sensitive values
	redactor := sanitize.NewRedactor()
	sanitizedContent := redactor.RedactLines(historyContent.String())

	// Save history to output store
	_, err = runContext.OutputStore.Save(
		"host-command-history.txt",
		strings.NewReader(sanitizedContent),
	)
	if err != nil {
		log.Warn("failed to save command history output: %s", err)
	}

	// Return empty results on success (output is saved)
	return []check.Result{}
}

// readAndLimitHistoryFile reads a history file and returns the last N lines
func readAndLimitHistoryFile(filePath string, maxLines int) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error scanning file: %w", err)
	}

	// Return the last maxLines lines
	startIndex := 0
	if len(lines) > maxLines {
		startIndex = len(lines) - maxLines
	}

	return strings.Join(lines[startIndex:], "\n"), nil
}
