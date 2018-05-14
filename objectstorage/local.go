package objectstorage

import (
	"os"
	"io"
)

type FileStorage struct {
	path string
}

func (fs *FileStorage) Reader() (io.Reader, error) {
	return os.Open(fs.path)
}

func (fs *FileStorage) Writer() (io.WriteCloser, error) {
	return os.OpenFile(fs.path, os.O_RDWR|os.O_CREATE, 0644)
}
