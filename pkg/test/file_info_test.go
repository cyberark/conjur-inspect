package test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFileInfo(t *testing.T) {
	m := &FileInfo{
		name:    "test.txt",
		size:    1024,
		mode:    os.ModePerm,
		modTime: time.Now(),
		isDir:   false,
		sys:     nil,
	}

	assert.Equal(t, "test.txt", m.Name())
	assert.Equal(t, int64(1024), m.Size())
	assert.Equal(t, os.ModePerm, m.Mode())
	assert.Equal(t, m.modTime, m.ModTime())
	assert.Equal(t, false, m.IsDir())
	assert.Nil(t, m.Sys())
}
