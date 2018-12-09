package usecase

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"github.com/scylladb/go-set/strset"
	"github.com/stretchr/testify/assert"

	"github.com/suzuki-shunsuke/durl/internal/test"
)

func Test_extractURLsFromFiles(t *testing.T) {
	type (
		File struct {
			buf []byte
			err error
		}
	)
	data := []struct {
		title    string
		files    map[string]File
		checkErr func(assert.TestingT, interface{}, ...interface{}) bool
		set      *strset.Set
	}{{
		"no url", map[string]File{
			"foo.txt": File{[]byte(`foo`), nil},
		}, assert.Nil, strset.New(),
	}, {
		"normal", map[string]File{
			"foo.txt": File{[]byte(`foo`), nil},
			"bar.txt": File{[]byte(`http://example.com`), nil},
		}, assert.Nil, strset.New("http://example.com"),
	}, {
		"error", map[string]File{
			"bar.txt": File{[]byte(`http://example.com`), nil},
			"foo.txt": File{nil, fmt.Errorf("failed to read a file")},
		}, assert.NotNil, nil,
	}}
	for _, tt := range data {
		t.Run(tt.title, func(t *testing.T) {
			fsys := test.NewFsys(t, nil).
				SetFuncOpen(func(p string) (io.ReadCloser, error) {
					if f, ok := tt.files[p]; ok {
						if f.buf != nil {
							return ioutil.NopCloser(bytes.NewBuffer(f.buf)), f.err
						}
						return nil, f.err
					}
					return nil, fmt.Errorf("file is not found: %s", p)
				})
			files := strset.New()
			for k := range tt.files {
				files.Add(k)
			}
			set, err := extractURLsFromFiles(fsys, files)
			tt.checkErr(t, err)
			if err == nil {
				if !set.IsEqual(tt.set) {
					t.Fatalf("set = %v, wanted %v", set, tt.set)
				}
			}
		})
	}
}

func Test_extractURLsFromFile(t *testing.T) {
	data := []struct {
		title    string
		buf      []byte
		err      error
		checkErr func(assert.TestingT, interface{}, ...interface{}) bool
		set      *strset.Set
		p        string
	}{{
		"no url", []byte(`foo
bar`), nil, assert.Nil, strset.New(), "foo.txt",
	}, {
		"normal", []byte(`http://example.com`), nil, assert.Nil, strset.New("http://example.com"), "foo.txt",
	}, {
		"error", nil, fmt.Errorf("failed to read a file"), assert.NotNil, nil, "foo.txt",
	}}
	for _, tt := range data {
		t.Run(tt.title, func(t *testing.T) {
			var rc io.ReadCloser
			if tt.buf != nil {
				rc = ioutil.NopCloser(bytes.NewBuffer(tt.buf))
			}
			fsys := test.NewFsys(t, nil).
				SetReturnOpen(rc, tt.err)
			set, err := extractURLsFromFile(context.Background(), fsys, tt.p)
			tt.checkErr(t, err)
			if err == nil {
				if !set.IsEqual(tt.set) {
					t.Fatalf("set = %v, wanted %v", set, tt.set)
				}
			}
		})
	}
}

func Test_getFiles(t *testing.T) {
	data := []struct {
		title    string
		in       string
		checkErr func(assert.TestingT, interface{}, ...interface{}) bool
		arr      *strset.Set
	}{{
		"normal", `foo
bar`, assert.Nil, strset.New("foo", "bar"),
	}, {
		"spaces", `
  foo
bar
`, assert.Nil, strset.New("foo", "bar"),
	}}
	for _, tt := range data {
		t.Run(tt.title, func(t *testing.T) {
			arr, err := getFiles(bytes.NewBufferString(tt.in))
			tt.checkErr(t, err)
			if err != nil {
				if !arr.IsEqual(tt.arr) {
					t.Fatalf("arr = %v, wanted %v", arr, tt.arr)
				}
			}
		})
	}
}
