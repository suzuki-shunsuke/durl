package domain

import (
	"io"
)

type (
	Fsys interface {
		Open(string) (io.ReadCloser, error)
	}
)
