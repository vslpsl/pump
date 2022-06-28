package file

import "os"

type Writer struct {
	file *os.File
}

func (w *Writer) Write(p []byte) (n int, err error) {
	return w.file.Write(p)
}

func (w *Writer) Close() error {
	return w.file.Close()
}

func NewWriter(filePath string) (*Writer, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0644)
	return &Writer{file: file}, err
}
