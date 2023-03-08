package output

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTarGzipArchive_Archive(t *testing.T) {
	testName := "test-archive"
	dir, err := ioutil.TempDir("", "test-store-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	store := NewDirectoryStore(dir)

	store.Save("test1.txt", strings.NewReader("test 1"))
	store.Save("test2.json", strings.NewReader("test 2"))

	// Run the test
	archive := &TarGzipArchive{}
	err = archive.Archive(testName, store)
	assert.Nil(t, err)

	// Verify the result
	file, err := os.Open(testName + ".tar.gz")
	assert.Nil(t, err)
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	assert.Nil(t, err)
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	// Check that both test files are part of the archive
	foundFiles := make(map[string]bool)
	for {
		hdr, err := tarReader.Next()
		if err == io.EOF {
			break // End of archive
		}
		assert.Nil(t, err)

		foundFiles[hdr.Name] = true

		body, err := ioutil.ReadAll(tarReader)
		assert.Nil(t, err)
		assert.NotEmpty(t, body)
	}

	assert.True(t, foundFiles["test-archive/test1.txt"])
	assert.True(t, foundFiles["test-archive/test2.json"])

	// Clean up
	err = os.Remove(testName + ".tar.gz")
	assert.Nil(t, err)
}
