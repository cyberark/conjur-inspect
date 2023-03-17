package output

import (
	"io"
	"io/fs"
)

// Store represent an object that can save raw outputs, referencing them by name.
type Store interface {
	Save(name string, reader io.Reader) error
	Items() ([]StoreItem, error)
	Cleanup() error
}

// StoreItem represents a particular raw output that has been saved to a Store.
type StoreItem interface {
	Info() (fs.FileInfo, error)
	Open() (reader io.Reader, cleanup func() error, err error)
}
