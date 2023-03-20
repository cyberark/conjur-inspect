package test

import (
	"os"
	"time"
)

// FileInfo is a mock implementation of the os.FileInfo interface for unit testing purposes.
type FileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
	sys     interface{}
}

// Name returns the mock name
func (info *FileInfo) Name() string {
	return info.name
}

// Size returns the mock size
func (info *FileInfo) Size() int64 {
	return info.size
}

// Mode returns the mock node
func (info *FileInfo) Mode() os.FileMode {
	return info.mode
}

// ModTime returns the mock modification time
func (info *FileInfo) ModTime() time.Time {
	return info.modTime
}

// IsDir reports the mock directory value
func (info *FileInfo) IsDir() bool {
	return info.isDir
}

// Sys returns the mock sys
func (info *FileInfo) Sys() interface{} {
	return info.sys
}
