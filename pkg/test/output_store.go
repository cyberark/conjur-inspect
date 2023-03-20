package test

import (
	"bytes"
	"io"
	"os"
	"time"

	"github.com/cyberark/conjur-inspect/pkg/output"
)

// OutputStore is a mock implementation of the output.Store interface for
// unit testing purposes.
type OutputStore struct {
	items map[string]OutputStoreItem
}

// OutputStoreItem is a mock implementation of the output.StoreItem interface
// for unit testing purposes.
type OutputStoreItem struct {
	name string
	data []byte
}

// NewOutputStore returns a new mock (in-memory) output store. The intended use
// for this is in unit testing.
func NewOutputStore() *OutputStore {
	return &OutputStore{
		items: make(map[string]OutputStoreItem),
	}
}

// Save stores a given output to the directory as a file
func (store *OutputStore) Save(name string, reader io.Reader) error {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(reader)
	if err != nil {
		return err
	}

	store.items[name] = OutputStoreItem{
		name: name,
		data: buf.Bytes(),
	}

	return nil
}

// Items returns the collection of outputs store in this directory
func (store *OutputStore) Items() ([]output.StoreItem, error) {
	items := make([]output.StoreItem, 0, len(store.items))

	for _, value := range store.items {
		items = append(items, &value)
	}

	return items, nil
}

// Cleanup removes the directory and files used for this output store
func (store *OutputStore) Cleanup() error {
	store.items = make(map[string]OutputStoreItem)
	return nil
}

// Info returns the stat results for a given item (file) in a DirectoryStore
func (item *OutputStoreItem) Info() (os.FileInfo, error) {
	info := &FileInfo{
		name:    item.name,
		size:    int64(len(item.data)),
		isDir:   false,
		modTime: time.Now(),
		mode:    0644,
	}

	return info, nil
}

// Open returns the io.Reader to the file for the DirectoryStoreItem
func (item *OutputStoreItem) Open() (io.Reader, func() error, error) {
	reader := bytes.NewReader(item.data)
	cleanupFunc := func() error { return nil }

	return reader, cleanupFunc, nil
}
