package domain

import (
	"io"
)

type (
	// Fsys represents operation to filesystem.
	Fsys interface {
		Exist(string) bool
		Getwd() (string, error)
		Open(string) (io.ReadCloser, error)
		Write(string, []byte) error
	}
)
