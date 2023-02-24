package maybe

import (
	"io"
)

// Writer wraps an io.Writer with the Maybe
// pattern to allow for fewer error checks when
// writing many strings in a row.
type Writer struct {
	maybe Maybe[io.Writer]
}

// NewWriter creates a new Writer wrapping the given io.Writer
func NewWriter(writer io.Writer) *Writer {
	return &Writer{
		maybe: NewSuccess(writer),
	}
}

func (maybeWriter *Writer) WriteString(str string) {
	maybeWriter.maybe = Bind(
		maybeWriter.maybe,
		func(writer io.Writer) (io.Writer, error) {
			_, err := io.WriteString(writer, str)
			return writer, err
		},
	)
}

func (maybeWriter *Writer) Error() error {
	return maybeWriter.maybe.Error()
}
