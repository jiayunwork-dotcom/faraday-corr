package metal

import (
	"bytes"
	"io"
)

// specSink buffers a Faraday spec encode and flushes on the first
// Close. A second Close is a no-op so a caller that defers Close
// after an explicit Close does not rewind the destination.
type specSink struct {
	dst    io.Writer
	buf    bytes.Buffer
	closed bool
}

func (s *specSink) Write(p []byte) (int, error) {
	return s.buf.Write(p)
}

func (s *specSink) Close() error {
	if s.closed {
		return nil
	}
	s.closed = true
	_, err := s.dst.Write(s.buf.Bytes())
	return err
}
