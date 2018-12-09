package domain

import (
	"io"
)

type (
	Fsys interface {
		Exist(string) bool
		Getwd() (string, error)
		Open(string) (io.ReadCloser, error)
		Write(string, []byte) error
	}
)
