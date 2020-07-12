package infra

import (
	"io"
	"io/ioutil"
	"os"
)

type (
	// Fsys represents operation to filesystem.
	Fsys struct{}
)

// Open opens a file.
func (fsys Fsys) Open(p string) (io.ReadCloser, error) {
	return os.Open(p)
}

// Exist returns whether file exists.
func (fsys Fsys) Exist(dst string) bool {
	_, err := os.Stat(dst)
	return err == nil
}

// Write writes data to a file.
func (fsys Fsys) Write(dst string, data []byte) error {
	return ioutil.WriteFile(dst, data, 0o644) //nolint:gosec
}

// Getwd returns a current directory path.
func (fsys Fsys) Getwd() (string, error) {
	return os.Getwd()
}
