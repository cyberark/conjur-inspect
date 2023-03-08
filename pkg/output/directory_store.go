package output

import (
	"fmt"
	"io"
	"os"
	"path"
)

// DirectoryStore is an output store implementation that stores outputs as files
// in a given directory.
type DirectoryStore struct {
	directory string
}

// NewDirectoryStore instantiates a new DirectoryStore struct
func NewDirectoryStore(directory string) *DirectoryStore {
	return &DirectoryStore{
		directory: directory,
	}
}

// Save stores a given output to the directory as a file
func (dirStore *DirectoryStore) Save(name string, reader io.Reader) error {
	path := path.Join(dirStore.directory, name)

	fmt.Println(path)

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	return err
}

// Items returns the collection of outputs store in this directory
func (dirStore *DirectoryStore) Items() ([]StoreItem, error) {
	files, err := os.ReadDir(dirStore.directory)
	if err != nil {
		return nil, err
	}

	var items []StoreItem
	for _, file := range files {
		// Don't include directories in the list, it should only be a collection
		// of flat files.
		if file.IsDir() {
			continue
		}

		// Resolve the absolute path
		filePath := path.Join(dirStore.directory, file.Name())

		// Add the FileStoreItem to the array
		items = append(items, &DirectoryStoreItem{path: filePath})
	}

	return items, nil
}

// Cleanup removes the directory and files used for this output store
func (dirStore *DirectoryStore) Cleanup() error {
	return os.RemoveAll(dirStore.directory)
}
