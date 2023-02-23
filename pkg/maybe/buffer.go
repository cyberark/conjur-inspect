package maybe

import (
	"bytes"
)

type Buffer struct {
	maybe Maybe[*bytes.Buffer]
}

func NewBuffer() *Buffer {
	return &Buffer{
		maybe: NewSuccess(new(bytes.Buffer)),
	}
}

func (maybeBuffer *Buffer) WriteString(str string) {
	maybeBuffer.maybe = Bind(
		maybeBuffer.maybe,
		func(buf *bytes.Buffer) (*bytes.Buffer, error) {
			_, err := buf.WriteString(str)
			return buf, err
		},
	)
}

func (maybeBuffer *Buffer) String() (string, error) {
	if maybeBuffer.maybe.Error() != nil {
		return "", maybeBuffer.maybe.Error()
	}

	return maybeBuffer.maybe.Value().String(), nil
}
