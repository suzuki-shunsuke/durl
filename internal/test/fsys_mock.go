package test

// Don't edit this file.
// This file is generated by gomic 0.5.2.
// https://github.com/suzuki-shunsuke/gomic

import (
	"io"
	testing "testing"

	gomic "github.com/suzuki-shunsuke/gomic/gomic"
)

type (
	// Fsys is a mock.
	Fsys struct {
		t                      *testing.T
		name                   string
		callbackNotImplemented gomic.CallbackNotImplemented
		impl                   struct {
			Exist func(p0 string) bool
			Getwd func() (string, error)
			Open  func(p0 string) (io.ReadCloser, error)
			Write func(p0 string, p1 []byte) error
		}
	}
)

// NewFsys returns Fsys .
func NewFsys(t *testing.T, cb gomic.CallbackNotImplemented) *Fsys {
	return &Fsys{
		t: t, name: "Fsys", callbackNotImplemented: cb}
}

// Exist is a mock method.
func (mock Fsys) Exist(p0 string) bool {
	methodName := "Exist" // nolint: goconst
	if mock.impl.Exist != nil {
		return mock.impl.Exist(p0)
	}
	if mock.callbackNotImplemented != nil {
		mock.callbackNotImplemented(mock.t, mock.name, methodName)
	} else {
		gomic.DefaultCallbackNotImplemented(mock.t, mock.name, methodName)
	}
	return mock.fakeZeroExist(p0)
}

// SetFuncExist sets a method and returns the mock.
func (mock *Fsys) SetFuncExist(impl func(p0 string) bool) *Fsys {
	mock.impl.Exist = impl
	return mock
}

// SetReturnExist sets a fake method.
func (mock *Fsys) SetReturnExist(r0 bool) *Fsys {
	mock.impl.Exist = func(string) bool {
		return r0
	}
	return mock
}

// fakeZeroExist is a fake method which returns zero values.
func (mock Fsys) fakeZeroExist(p0 string) bool {
	var (
		r0 bool
	)
	return r0
}

// Getwd is a mock method.
func (mock Fsys) Getwd() (string, error) {
	methodName := "Getwd" // nolint: goconst
	if mock.impl.Getwd != nil {
		return mock.impl.Getwd()
	}
	if mock.callbackNotImplemented != nil {
		mock.callbackNotImplemented(mock.t, mock.name, methodName)
	} else {
		gomic.DefaultCallbackNotImplemented(mock.t, mock.name, methodName)
	}
	return mock.fakeZeroGetwd()
}

// SetFuncGetwd sets a method and returns the mock.
func (mock *Fsys) SetFuncGetwd(impl func() (string, error)) *Fsys {
	mock.impl.Getwd = impl
	return mock
}

// SetReturnGetwd sets a fake method.
func (mock *Fsys) SetReturnGetwd(r0 string, r1 error) *Fsys {
	mock.impl.Getwd = func() (string, error) {
		return r0, r1
	}
	return mock
}

// fakeZeroGetwd is a fake method which returns zero values.
func (mock Fsys) fakeZeroGetwd() (string, error) {
	var (
		r0 string
		r1 error
	)
	return r0, r1
}

// Open is a mock method.
func (mock Fsys) Open(p0 string) (io.ReadCloser, error) {
	methodName := "Open" // nolint: goconst
	if mock.impl.Open != nil {
		return mock.impl.Open(p0)
	}
	if mock.callbackNotImplemented != nil {
		mock.callbackNotImplemented(mock.t, mock.name, methodName)
	} else {
		gomic.DefaultCallbackNotImplemented(mock.t, mock.name, methodName)
	}
	return mock.fakeZeroOpen(p0)
}

// SetFuncOpen sets a method and returns the mock.
func (mock *Fsys) SetFuncOpen(impl func(p0 string) (io.ReadCloser, error)) *Fsys {
	mock.impl.Open = impl
	return mock
}

// SetReturnOpen sets a fake method.
func (mock *Fsys) SetReturnOpen(r0 io.ReadCloser, r1 error) *Fsys {
	mock.impl.Open = func(string) (io.ReadCloser, error) {
		return r0, r1
	}
	return mock
}

// fakeZeroOpen is a fake method which returns zero values.
func (mock Fsys) fakeZeroOpen(p0 string) (io.ReadCloser, error) {
	var (
		r0 io.ReadCloser
		r1 error
	)
	return r0, r1
}

// Write is a mock method.
func (mock Fsys) Write(p0 string, p1 []byte) error {
	methodName := "Write" // nolint: goconst
	if mock.impl.Write != nil {
		return mock.impl.Write(p0, p1)
	}
	if mock.callbackNotImplemented != nil {
		mock.callbackNotImplemented(mock.t, mock.name, methodName)
	} else {
		gomic.DefaultCallbackNotImplemented(mock.t, mock.name, methodName)
	}
	return mock.fakeZeroWrite(p0, p1)
}

// SetFuncWrite sets a method and returns the mock.
func (mock *Fsys) SetFuncWrite(impl func(p0 string, p1 []byte) error) *Fsys {
	mock.impl.Write = impl
	return mock
}

// SetReturnWrite sets a fake method.
func (mock *Fsys) SetReturnWrite(r0 error) *Fsys {
	mock.impl.Write = func(string, []byte) error {
		return r0
	}
	return mock
}

// fakeZeroWrite is a fake method which returns zero values.
func (mock Fsys) fakeZeroWrite(p0 string, p1 []byte) error {
	var (
		r0 error
	)
	return r0
}
