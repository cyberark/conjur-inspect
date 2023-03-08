package output

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestDirectoryStoreItem_Info(t *testing.T) {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	expected, err := os.Stat(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	item := DirectoryStoreItem{path: file.Name()}
	info, err := item.Info()
	if err != nil {
		t.Fatal(err)
	}

	if info.Name() != expected.Name() {
		t.Errorf("Expected name %s, but got %s", expected.Name(), info.Name())
	}

	if info.Size() != expected.Size() {
		t.Errorf("Expected size %d, but got %d", expected.Size(), info.Size())
	}

	if info.Mode() != expected.Mode() {
		t.Errorf("Expected mode %v, but got %v", expected.Mode(), info.Mode())
	}
}

func TestDirectoryStoreItem_Open(t *testing.T) {
	content := "Hello, World!"
	file, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	if _, err := file.WriteString(content); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	item := DirectoryStoreItem{path: file.Name()}
	reader, cleanup, err := item.Open()
	if err != nil {
		t.Fatal(err)
	}

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != content {
		t.Errorf("Expected '%s', but got '%s'", content, string(data))
	}

	if err := cleanup(); err != nil {
		t.Fatalf("Error on cleanup: %v", err)
	}
}
