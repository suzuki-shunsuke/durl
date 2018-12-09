package usecase

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"gopkg.in/h2non/gock.v1"

	"github.com/scylladb/go-set/strset"
	"github.com/stretchr/testify/assert"

	"github.com/suzuki-shunsuke/durl/internal/domain"
	"github.com/suzuki-shunsuke/durl/internal/test"
)

type (
	File struct {
		buf []byte
		err error
	}
)

func newFsys(t *testing.T, files map[string]File) domain.Fsys {
	return test.NewFsys(t, nil).
		SetFuncOpen(func(p string) (io.ReadCloser, error) {
			if f, ok := files[p]; ok {
				if f.buf != nil {
					return ioutil.NopCloser(bytes.NewBuffer(f.buf)), f.err
				}
				return nil, f.err
			}
			return nil, fmt.Errorf("file is not found: %s", p)
		})
}

func TestCheck(t *testing.T) {
	defer gock.Off()
	data := []struct {
		title    string
		in       string
		replies  map[string]int
		files    map[string]File
		checkErr func(assert.TestingT, interface{}, ...interface{}) bool
	}{{
		"normal", "foo.txt", map[string]int{"/foo": 200},
		map[string]File{
			"foo.txt": {[]byte("http://example.com/foo"), nil},
		}, assert.Nil,
	}, {
		"http error", "foo.txt", map[string]int{"/foo": 500},
		map[string]File{
			"foo.txt": {[]byte("http://example.com/foo"), nil},
		}, assert.NotNil,
	}, {
		"file read error", "foo.txt", map[string]int{"/foo": 200},
		map[string]File{
			"foo.txt": {nil, fmt.Errorf("failed to read a file")},
		}, assert.NotNil,
	}}
	for _, tt := range data {
		t.Run(tt.title, func(t *testing.T) {
			g := gock.New("http://example.com")
			for p, c := range tt.replies {
				g.Get(p).Reply(c)
			}
			tt.checkErr(t, Check(newFsys(t, tt.files), bytes.NewBufferString(tt.in)))
		})
	}
}

func Test_checkURLs(t *testing.T) {
	defer gock.Off()
	data := []struct {
		title    string
		replies  map[string]int
		urls     map[string]*strset.Set
		checkErr func(assert.TestingT, interface{}, ...interface{}) bool
	}{{
		"normal", map[string]int{"/foo": 200},
		map[string]*strset.Set{
			"http://example.com/foo": strset.New("foo.txt"),
		}, assert.Nil,
	}, {
		"error", map[string]int{"/foo": 200, "/bar": 500},
		map[string]*strset.Set{
			"http://example.com/foo": strset.New("foo.txt"),
			"http://example.com/bar": strset.New("bar.txt"),
		}, assert.NotNil,
	}}
	for _, tt := range data {
		t.Run(tt.title, func(t *testing.T) {
			g := gock.New("http://example.com")
			for p, c := range tt.replies {
				g.Get(p).Reply(c)
			}
			tt.checkErr(t, checkURLs(tt.urls))
		})
	}
}

func Test_checkURL(t *testing.T) {
	defer gock.Off()
	client := http.Client{
		Timeout: domain.DefaultTimeout,
	}

	data := []struct {
		title    string
		path     string
		reply    int
		checkErr func(assert.TestingT, interface{}, ...interface{}) bool
	}{{
		"normal", "/foo", 200, assert.Nil,
	}, {
		"500 error", "/bar", 500, assert.NotNil,
	}}
	for _, tt := range data {
		t.Run(tt.title, func(t *testing.T) {
			host, err := url.Parse("http://example.com")
			if err != nil {
				t.Fatal(err)
			}
			host.Path = tt.path
			gock.New("http://example.com").
				Get(tt.path).Reply(tt.reply)
			tt.checkErr(t, checkURL(context.Background(), client, host.String()))
		})
	}
}

func Test_extractURLsFromFiles(t *testing.T) {
	data := []struct {
		title    string
		files    map[string]File
		checkErr func(assert.TestingT, interface{}, ...interface{}) bool
		set      map[string]*strset.Set
	}{{
		"no url", map[string]File{
			"foo.txt": {[]byte(`foo`), nil},
		}, assert.Nil, map[string]*strset.Set{},
	}, {
		"normal", map[string]File{
			"foo.txt": {[]byte(`foo`), nil},
			"bar.txt": {[]byte(`http://example.com`), nil},
		}, assert.Nil, map[string]*strset.Set{
			"http://example.com": strset.New("bar.txt")},
	}, {
		"error", map[string]File{
			"bar.txt": {[]byte(`http://example.com`), nil},
			"foo.txt": {nil, fmt.Errorf("failed to read a file")},
		}, assert.NotNil, nil,
	}}
	for _, tt := range data {
		t.Run(tt.title, func(t *testing.T) {
			fsys := newFsys(t, tt.files)
			files := strset.New()
			for k := range tt.files {
				files.Add(k)
			}
			set, err := extractURLsFromFiles(fsys, files)
			tt.checkErr(t, err)
			if err == nil {
				assert.Equal(t, tt.set, set)
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
