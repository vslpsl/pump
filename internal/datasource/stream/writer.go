package stream

import "io"

type Writer struct {
	destination io.Writer
}

func (w *Writer) Write(p []byte) (n int, err error) {
	return w.destination.Write(p)
}

func (w *Writer) Close() error {
	return nil
}

func NewWriter(destination io.Writer) *Writer {
	return &Writer{destination: destination}
}
