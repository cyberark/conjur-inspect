package output

import (
	"io"
	"os"
)

// DirectoryStoreItem is a reference to a file stored in a DirectoryStore
type DirectoryStoreItem struct {
	path string
}

// Info returns the stat results for a given item (file) in a DirectoryStore
func (item *DirectoryStoreItem) Info() (os.FileInfo, error) {
	return os.Stat(item.path)
}

// Open returns the io.Reader to the file for the DirectoryStoreItem
func (item *DirectoryStoreItem) Open() (io.Reader, func() error, error) {
	reader, err := os.Open(item.path)
	if err != nil {
		return nil, nil, err
	}

	cleanupFunc := func() error {
		return reader.Close()
	}

	return reader, cleanupFunc, nil
}
