// Package test contains test utilities and mocks for unit testing purposes.
package test

import "github.com/cyberark/conjur-inspect/pkg/output"

// OutputArchive is a mock implementation of the output.Archive interface for
// unit testing purposes.
type OutputArchive struct {
	archiveCalled bool
}

// Archive is a no-op that records that it was called for unit test assertions.
func (oa *OutputArchive) Archive(name string, store output.Store) error {
	oa.archiveCalled = true
	return nil
}

// IsArchived returns whether the Archive method was called.
func (oa *OutputArchive) IsArchived() bool {
	return oa.archiveCalled
}
