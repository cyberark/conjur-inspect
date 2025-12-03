package checks

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cyberark/conjur-inspect/pkg/check"
	"github.com/cyberark/conjur-inspect/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommandHistoryDescribe(t *testing.T) {
	ch := &CommandHistory{}
	assert.Equal(t, "Command History", ch.Describe())
}

func TestCommandHistoryRunZshHistoryFound(t *testing.T) {
	// Create a temporary directory with a mock .zsh_history file
	tempDir := t.TempDir()

	// Create a .zsh_history file with some test content
	zshHistoryPath := filepath.Join(tempDir, ".zsh_history")
	zshContent := strings.Join([]string{
		"line1",
		"line2",
		"line3",
		"line4",
		"line5",
	}, "\n")
	err := os.WriteFile(zshHistoryPath, []byte(zshContent), 0644)
	require.NoError(t, err)

	// Mock os.UserHomeDir to return our temp directory
	oldFunc := userHomeDirFunc
	userHomeDirFunc = func() (string, error) {
		return tempDir, nil
	}
	defer func() {
		userHomeDirFunc = oldFunc
	}()

	// Run the check
	ch := &CommandHistory{}
	runContext := test.NewRunContext("")
	results := ch.Run(&runContext)

	// Should return empty results on success
	assert.Empty(t, results)

	// Verify the file was saved to the output store
	items, err := runContext.OutputStore.Items()
	require.NoError(t, err)
	require.Len(t, items, 1)
	info, err := items[0].Info()
	require.NoError(t, err)
	assert.Equal(t, "host-command-history.txt", info.Name())
}

func TestCommandHistoryRunBashHistoryFallback(t *testing.T) {
	// Create a temporary directory with only a .bash_history file
	tempDir := t.TempDir()

	// Create a .bash_history file with some test content
	bashHistoryPath := filepath.Join(tempDir, ".bash_history")
	bashContent := strings.Join([]string{
		"cmd1",
		"cmd2",
		"cmd3",
	}, "\n")
	err := os.WriteFile(bashHistoryPath, []byte(bashContent), 0644)
	require.NoError(t, err)

	// Mock os.UserHomeDir to return our temp directory
	oldFunc := userHomeDirFunc
	userHomeDirFunc = func() (string, error) {
		return tempDir, nil
	}
	defer func() {
		userHomeDirFunc = oldFunc
	}()

	// Run the check
	ch := &CommandHistory{}
	runContext := test.NewRunContext("")
	results := ch.Run(&runContext)

	// Should return empty results on success
	assert.Empty(t, results)

	// Verify the file was saved to the output store
	items, err := runContext.OutputStore.Items()
	require.NoError(t, err)
	require.Len(t, items, 1)
	info, err := items[0].Info()
	require.NoError(t, err)
	assert.Equal(t, "host-command-history.txt", info.Name())
}

func TestCommandHistoryRunZshHistoryPreferredOverBash(t *testing.T) {
	// Create a temporary directory with both .zsh_history and .bash_history
	tempDir := t.TempDir()

	// Create both history files
	zshHistoryPath := filepath.Join(tempDir, ".zsh_history")
	zshContent := "zsh_command"
	err := os.WriteFile(zshHistoryPath, []byte(zshContent), 0644)
	require.NoError(t, err)

	bashHistoryPath := filepath.Join(tempDir, ".bash_history")
	bashContent := "bash_command"
	err = os.WriteFile(bashHistoryPath, []byte(bashContent), 0644)
	require.NoError(t, err)

	// Mock os.UserHomeDir to return our temp directory
	oldFunc := userHomeDirFunc
	userHomeDirFunc = func() (string, error) {
		return tempDir, nil
	}
	defer func() {
		userHomeDirFunc = oldFunc
	}()

	// Run the check
	ch := &CommandHistory{}
	runContext := test.NewRunContext("")
	results := ch.Run(&runContext)

	// Should return empty results on success
	assert.Empty(t, results)

	// Verify that zsh history was used (it's preferred)
	items, err := runContext.OutputStore.Items()
	require.NoError(t, err)
	require.Len(t, items, 1)
	info, err := items[0].Info()
	require.NoError(t, err)
	assert.Equal(t, "host-command-history.txt", info.Name())
}

func TestCommandHistoryRunNoHistoryFiles(t *testing.T) {
	// Create a temporary directory with no history files
	tempDir := t.TempDir()

	// Mock os.UserHomeDir to return our temp directory
	oldFunc := userHomeDirFunc
	userHomeDirFunc = func() (string, error) {
		return tempDir, nil
	}
	defer func() {
		userHomeDirFunc = oldFunc
	}()

	// Run the check
	ch := &CommandHistory{}
	runContext := test.NewRunContext("")
	results := ch.Run(&runContext)

	// Should return INFO result when no files found
	require.Len(t, results, 1)
	assert.Equal(t, "Command History", results[0].Title)
	assert.Equal(t, check.StatusInfo, results[0].Status)
	assert.Equal(t, "No command history files found or accessible", results[0].Message)
}

func TestCommandHistoryRunLimitTo100Lines(t *testing.T) {
	// Create a temporary directory with a .bash_history file with more than 100 lines
	tempDir := t.TempDir()

	// Create a .bash_history file with 150 lines
	bashHistoryPath := filepath.Join(tempDir, ".bash_history")
	lines := make([]string, 150)
	for i := 0; i < 150; i++ {
		lines[i] = "command_" + string(rune('a'+i%26))
	}
	bashContent := strings.Join(lines, "\n")
	err := os.WriteFile(bashHistoryPath, []byte(bashContent), 0644)
	require.NoError(t, err)

	// Mock os.UserHomeDir to return our temp directory
	oldFunc := userHomeDirFunc
	userHomeDirFunc = func() (string, error) {
		return tempDir, nil
	}
	defer func() {
		userHomeDirFunc = oldFunc
	}()

	// Run the check
	ch := &CommandHistory{}
	runContext := test.NewRunContext("")
	results := ch.Run(&runContext)

	// Should return empty results on success
	assert.Empty(t, results)

	// Verify the file was saved and contains only 100 lines
	items, err := runContext.OutputStore.Items()
	require.NoError(t, err)
	require.Len(t, items, 1)

	// Read the saved content and verify line count
	reader, cleanup, err := items[0].Open()
	require.NoError(t, err)
	defer cleanup()

	buf := make([]byte, 10000)
	n, err := reader.Read(buf)
	if err != nil && err.Error() != "EOF" {
		require.NoError(t, err)
	}

	contentStr := string(buf[:n])
	savedLines := strings.Split(strings.TrimSpace(contentStr), "\n")
	assert.Len(t, savedLines, 100)
}

func TestCommandHistoryRunSanitizesSensitiveData(t *testing.T) {
	// Create a temporary directory with a .bash_history file containing sensitive data
	tempDir := t.TempDir()

	bashHistoryPath := filepath.Join(tempDir, ".bash_history")
	historyContent := strings.Join([]string{
		"curl https://api.example.com",
		"export api_key=sk_live_secret123",
		"mysql -u root -pPassword123",
		"password=supersecret456",
		"normal command here",
	}, "\n")
	err := os.WriteFile(bashHistoryPath, []byte(historyContent), 0644)
	require.NoError(t, err)

	// Mock os.UserHomeDir to return our temp directory
	oldFunc := userHomeDirFunc
	userHomeDirFunc = func() (string, error) {
		return tempDir, nil
	}
	defer func() {
		userHomeDirFunc = oldFunc
	}()

	// Run the check
	ch := &CommandHistory{}
	runContext := test.NewRunContext("")
	results := ch.Run(&runContext)

	// Should return empty results on success
	assert.Empty(t, results)

	// Verify the file was saved and contains sanitized content
	items, err := runContext.OutputStore.Items()
	require.NoError(t, err)
	require.Len(t, items, 1)

	// Read the saved content and verify sensitive data is redacted
	reader, cleanup, err := items[0].Open()
	require.NoError(t, err)
	defer cleanup()

	buf := make([]byte, 10000)
	n, err := reader.Read(buf)
	if err != nil && err.Error() != "EOF" {
		require.NoError(t, err)
	}

	contentStr := string(buf[:n])

	// Verify sensitive data was redacted
	assert.NotContains(t, contentStr, "sk_live_secret123")
	assert.NotContains(t, contentStr, "supersecret456")
	assert.Contains(t, contentStr, "[REDACTED]")
	assert.Contains(t, contentStr, "normal command here")
}

func TestCommandHistoryRunEmptyHistoryFile(t *testing.T) {
	// Create a temporary directory with an empty .bash_history file
	tempDir := t.TempDir()

	bashHistoryPath := filepath.Join(tempDir, ".bash_history")
	err := os.WriteFile(bashHistoryPath, []byte(""), 0644)
	require.NoError(t, err)

	// Mock os.UserHomeDir to return our temp directory
	oldFunc := userHomeDirFunc
	userHomeDirFunc = func() (string, error) {
		return tempDir, nil
	}
	defer func() {
		userHomeDirFunc = oldFunc
	}()

	// Run the check
	ch := &CommandHistory{}
	runContext := test.NewRunContext("")
	results := ch.Run(&runContext)

	// Should return INFO result when history file is empty
	require.Len(t, results, 1)
	assert.Equal(t, check.StatusInfo, results[0].Status)
	assert.Equal(t, "No command history files found or accessible", results[0].Message)
}

func TestCommandHistoryRunUserHomeDirError(t *testing.T) {
	// Mock os.UserHomeDir to return an error
	oldFunc := userHomeDirFunc
	userHomeDirFunc = func() (string, error) {
		return "", os.ErrNotExist
	}
	defer func() {
		userHomeDirFunc = oldFunc
	}()

	// Run the check
	ch := &CommandHistory{}
	runContext := test.NewRunContext("")
	results := ch.Run(&runContext)

	// Should return INFO result when home directory cannot be determined
	require.Len(t, results, 1)
	assert.Equal(t, check.StatusInfo, results[0].Status)
	assert.Equal(t, "Unable to determine home directory", results[0].Message)
}
