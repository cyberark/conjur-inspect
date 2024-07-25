// Package shell defines helpers for interacting with shell commands.
package shell

import "io"

// ReadOrDefault reads the data from the provided reader and returns it as a
// string. If the reader returns an error, the default value is returned.
func ReadOrDefault(reader io.Reader, defaultValue string) string {
	if reader == nil {
		return defaultValue
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return defaultValue
	}

	return string(data)
}
