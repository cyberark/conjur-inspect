package output

// Archive is an interface to converting an output store to a portable format,
// such as a gzipped tar file.
type Archive interface {
	Archive(name string, store Store) error
}
