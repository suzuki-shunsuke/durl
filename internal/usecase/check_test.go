package usecase

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/scylladb/go-set/strset"
	"github.com/stretchr/testify/require"
	"github.com/suzuki-shunsuke/gomic/gomic"

	"github.com/suzuki-shunsuke/durl/internal/domain"
	"github.com/suzuki-shunsuke/durl/internal/test"
)

type (
	File struct {
		buf []byte
		err error
	}
)

func newFsys(t *testing.T, files map[string]File) *test.Fsys {
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

func Test_logicCheck(t *testing.T) {
	data := []struct {
		title    string
		mock     domain.Logic
		checkErr func(require.TestingT, interface{}, ...interface{})
	}{{
		"normal", test.NewLogic(t, gomic.DoNothing), require.Nil,
	}, {
		"failed to read config", test.NewLogic(t, gomic.DoNothing).SetReturnReadCfg(domain.Cfg{}, fmt.Errorf("failed to read config")), require.NotNil,
	}, {
		"failed to get file paths", test.NewLogic(t, gomic.DoNothing).SetReturnGetFiles(nil, fmt.Errorf("failed to get file paths")), require.NotNil,
	}, {
		"failed to extract urls from files", test.NewLogic(t, gomic.DoNothing).SetReturnExtractURLsFromFiles(nil, fmt.Errorf("failed to extract urls from files")), require.NotNil,
	}}
	for _, tt := range data {
		t.Run(tt.title, func(t *testing.T) {
			lgc := &logic{
				logic: tt.mock,
			}
			tt.checkErr(t, lgc.Check(bytes.NewBufferString("stdin"), ""))
		})
	}
}

func Test_logicIsIgnoredURL(t *testing.T) {
	data := []struct {
		url string
		exp bool
		cfg domain.Cfg
	}{
		{"example.com", true, domain.Cfg{}},
		{"ldap://example.com", true, domain.Cfg{}},
		{"http://example.com", true, domain.Cfg{}},
		{"https://example.com", true, domain.Cfg{}},
		{"https://localhost.com", false, domain.Cfg{}},
		{"https://localhost.com", true, domain.Cfg{IgnoreURLs: []string{"https://localhost.com"}}},
		{"https://localhost.com", true, domain.Cfg{IgnoreHosts: []string{"localhost.com"}}},
		{"https://localhost", true, domain.Cfg{}},
		{"http://localhost", true, domain.Cfg{}},
		{"http://localhost:8000", true, domain.Cfg{}},
	}
	lgc := NewLogic(nil)
	for _, d := range data {
		if d.exp {
			require.True(t, lgc.IsIgnoredURL(d.url, d.cfg), d.url)
			continue
		}
		require.False(t, lgc.IsIgnoredURL(d.url, d.cfg), d.url)
	}
}

func Test_logicCheckURLs(t *testing.T) {
	data := []struct {
		title    string
		mock     domain.Logic
		urls     map[string]*strset.Set
		checkErr func(require.TestingT, interface{}, ...interface{})
	}{{
		"normal",
		test.NewLogic(t, gomic.DoNothing),
		map[string]*strset.Set{
			"http://example.com/foo": strset.New("foo.txt"),
		}, require.Nil,
	}, {
		"urls is empty",
		test.NewLogic(t, gomic.DoNothing), nil, require.Nil,
	}}
	cfg := domain.Cfg{HTTPMethod: "head,get"}
	for _, tt := range data {
		t.Run(tt.title, func(t *testing.T) {
			lgc := &logic{logic: tt.mock}
			tt.checkErr(t, lgc.CheckURLs(cfg, tt.urls))
		})
	}
}

func Test_logicCheckURL(t *testing.T) {
	client := http.Client{
		Timeout: domain.DefaultTimeout,
	}

	data := []struct {
		title    string
		method   string
		mock     domain.Logic
		checkErr func(require.TestingT, interface{}, ...interface{})
	}{{
		"get", "get",
		test.NewLogic(t, gomic.DoNothing),
		require.Nil,
	}, {
		"head", "head",
		test.NewLogic(t, gomic.DoNothing),
		require.Nil,
	}, {
		"head,get", "head,get",
		test.NewLogic(t, gomic.DoNothing),
		require.Nil,
	}, {
		"empty", "",
		test.NewLogic(t, gomic.DoNothing),
		require.Nil,
	}, {
		"invalid method", "invalid method",
		test.NewLogic(t, gomic.DoNothing),
		require.NotNil,
	}}
	for _, tt := range data {
		t.Run(tt.title, func(t *testing.T) {
			lgc := &logic{logic: tt.mock}
			cfg := domain.Cfg{HTTPMethod: tt.method}
			tt.checkErr(t, lgc.CheckURL(context.Background(), cfg, client, "http://example.com"))
		})
	}
}

func Test_logicExtractURLsFromFiles(t *testing.T) {
	data := []struct {
		title    string
		files    map[string]File
		checkErr func(require.TestingT, interface{}, ...interface{})
		set      map[string]*strset.Set
	}{{
		"no url", map[string]File{
			"foo.txt": {[]byte(`foo`), nil},
		}, require.Nil, map[string]*strset.Set{},
	}, {
		"normal", map[string]File{
			"foo.txt": {[]byte(`foo`), nil},
			"bar.txt": {[]byte(`http://example.com`), nil},
		}, require.Nil, map[string]*strset.Set{
			"http://example.com": strset.New("bar.txt")},
	}, {
		"error", map[string]File{
			"bar.txt": {[]byte(`http://example.com`), nil},
			"foo.txt": {nil, fmt.Errorf("failed to read a file")},
		}, require.Nil, map[string]*strset.Set{
			"http://example.com": strset.New("bar.txt")},
	}}
	for _, tt := range data {
		t.Run(tt.title, func(t *testing.T) {
			fsys := newFsys(t, tt.files)
			files := strset.New()
			for k := range tt.files {
				files.Add(k)
			}
			lgc := NewLogic(fsys)
			set, err := lgc.ExtractURLsFromFiles(files)
			tt.checkErr(t, err)
			if err == nil {
				require.Equal(t, tt.set, set)
			}
		})
	}
}

func Test_logicExtractURLsFromFile(t *testing.T) {
	data := []struct {
		title    string
		buf      []byte
		err      error
		checkErr func(require.TestingT, interface{}, ...interface{})
		set      *strset.Set
		p        string
	}{{
		"no url", []byte(`foo
bar`), nil, require.Nil, strset.New(), "foo.txt",
	}, {
		"normal", []byte(`http://example.com`), nil, require.Nil, strset.New("http://example.com"), "foo.txt",
	}, {
		"error", nil, fmt.Errorf("failed to read a file"), require.NotNil, nil, "foo.txt",
	}, {
		"too long", []byte(strings.Repeat("X", 65536)), nil, require.NotNil, nil, "foo.txt",
	}}
	for _, tt := range data {
		t.Run(tt.title, func(t *testing.T) {
			var rc io.ReadCloser
			if tt.buf != nil {
				rc = ioutil.NopCloser(bytes.NewBuffer(tt.buf))
			}
			fsys := test.NewFsys(t, nil).
				SetReturnOpen(rc, tt.err)
			lgc := NewLogic(fsys)
			set, err := lgc.ExtractURLsFromFile(context.Background(), tt.p)
			tt.checkErr(t, err)
			if err == nil {
				if !set.IsEqual(tt.set) {
					t.Fatalf("set = %v, wanted %v", set, tt.set)
				}
			}
		})
	}
}

func Test_logicGetFiles(t *testing.T) {
	data := []struct {
		title    string
		in       string
		checkErr func(require.TestingT, interface{}, ...interface{})
		arr      *strset.Set
	}{{
		"normal", `foo
bar`, require.Nil, strset.New("foo", "bar"),
	}, {
		"spaces", `
  foo
bar
`, require.Nil, strset.New("foo", "bar"),
	}}
	lgc := NewLogic(nil)
	for _, tt := range data {
		t.Run(tt.title, func(t *testing.T) {
			arr, err := lgc.GetFiles(bytes.NewBufferString(tt.in))
			tt.checkErr(t, err)
			if err != nil {
				if !arr.IsEqual(tt.arr) {
					t.Fatalf("arr = %v, wanted %v", arr, tt.arr)
				}
			}
		})
	}
}
