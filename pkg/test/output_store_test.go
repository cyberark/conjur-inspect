package test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOutputStore(t *testing.T) {
	store := NewOutputStore()
	assert.NotNil(t, store)

	t.Run("Save and retrieve", func(t *testing.T) {
		err := store.Save("test.txt", bytes.NewBufferString("Hello World!"))
		assert.NoError(t, err)

		items, err := store.Items()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(items))

		item := items[0]
		fi, err := item.Info()
		assert.NoError(t, err)
		assert.NotNil(t, fi)
		assert.Equal(t, "test.txt", fi.Name())

		reader, cleanup, err := item.Open()
		assert.NoError(t, err)
		assert.NotNil(t, reader)

		buf := make([]byte, 12)
		n, err := reader.Read(buf)
		assert.NoError(t, err)
		assert.Equal(t, 12, n)
		assert.Equal(t, "Hello World!", string(buf))

		err = cleanup()
		assert.NoError(t, err)

		err = store.Cleanup()
		assert.NoError(t, err)
	})

	t.Run("Empty store", func(t *testing.T) {
		items, err := store.Items()
		assert.NoError(t, err)
		assert.Equal(t, 0, len(items))
	})

	t.Run("Cleanup", func(t *testing.T) {
		err := store.Save("test.txt", bytes.NewBufferString("Hello World!"))
		assert.NoError(t, err)

		err = store.Cleanup()
		assert.NoError(t, err)

		items, err := store.Items()
		assert.NoError(t, err)
		assert.Equal(t, 0, len(items))
	})
}
