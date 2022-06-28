package file

import "os"

type Reader struct {
	file *os.File
}

func (r *Reader) Read(p []byte) (n int, err error) {
	return r.file.Read(p)
}

func (r *Reader) Close() error {
	return r.file.Close()
}

func NewReader(filePath string) (*Reader, error) {
	file, err := os.Open(filePath)
	return &Reader{file: file}, err
}
