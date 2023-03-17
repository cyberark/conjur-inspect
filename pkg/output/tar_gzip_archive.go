package output

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
)

// TarGzipArchive archives an output store as a Gzipped Tar archive.
type TarGzipArchive struct {
	OutputDir string
}

// Archive writes the gzipped tar archive for the given Store using the
// provided name.
func (archive *TarGzipArchive) Archive(
	name string,
	store Store,
) error {

	archiveName := fmt.Sprintf("%s.tar.gz", name)
	archivePath := path.Join(archive.OutputDir, archiveName)

	// Create output file
	out, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Create new Writers for gzip and tar
	// These writers are chained. Writing to the tar writer will
	// write to the gzip writer which in turn will write to
	// the "buf" writer
	gzipWriter := gzip.NewWriter(out)
	defer gzipWriter.Close()
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	items, err := store.Items()
	if err != nil {
		return err
	}

	for _, item := range items {
		err = archiveItem(tarWriter, name, item)
		if err != nil {
			return err
		}
	}

	return nil
}

func archiveItem(tarWriter *tar.Writer, prefix string, item StoreItem) error {
	fileInfo, err := item.Info()
	if err != nil {
		return err
	}

	// Get the local file information
	header, err := tar.FileInfoHeader(fileInfo, "")
	if err != nil {
		return err
	}

	// Add the archive prefix to the item name
	header.Name = path.Join(prefix, fileInfo.Name())

	// Write the header for this tar entry
	err = tarWriter.WriteHeader(header)
	if err != nil {
		return err
	}

	itemReader, cleanup, err := item.Open()
	if err != nil {
		return err
	}
	defer cleanup()

	_, err = io.Copy(tarWriter, itemReader)
	if err != nil {
		return err
	}

	return nil
}
