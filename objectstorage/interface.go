package objectstorage

import (
	"io"
)

type StoredObject interface {
	Reader() (io.Reader, error)
	Writer() (io.WriteCloser, error)
}
