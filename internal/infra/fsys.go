package infra

import (
	"io"
	"os"
)

type (
	Fsys struct{}
)

func (fsys Fsys) Open(p string) (io.ReadCloser, error) {
	return os.Open(p)
}
