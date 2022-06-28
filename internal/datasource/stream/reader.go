package stream

import "io"

type Reader struct {
	source io.Reader
}

func (r *Reader) Read(p []byte) (n int, err error) {
	return r.source.Read(p)
}

func (r *Reader) Close() error {
	return nil
}

func NewReader(source io.Reader) *Reader {
	return &Reader{source: source}
}
