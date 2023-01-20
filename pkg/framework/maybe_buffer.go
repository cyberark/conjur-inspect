package framework

import (
	"bytes"

	"github.com/conjurinc/conjur-preflight/pkg/maybe"
)

type MaybeBuffer struct {
	maybe maybe.Maybe[*bytes.Buffer]
}

func NewMaybeBuffer() *MaybeBuffer {
	return &MaybeBuffer{
		maybe: maybe.NewSuccess(new(bytes.Buffer)),
	}
}

func (maybeBuffer *MaybeBuffer) WriteString(str string) {
	maybeBuffer.maybe = maybe.Bind(
		maybeBuffer.maybe,
		func(buf *bytes.Buffer) (*bytes.Buffer, error) {
			_, err := buf.WriteString(str)
			return buf, err
		},
	)
}

func (maybeBuffer *MaybeBuffer) String() (string, error) {
	if maybeBuffer.maybe.Error() != nil {
		return "", maybeBuffer.maybe.Error()
	}

	return maybeBuffer.maybe.Value().String(), nil
}
