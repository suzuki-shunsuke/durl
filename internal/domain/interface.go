package domain

import (
	"context"
	"io"
	"net/http"

	"github.com/scylladb/go-set/strset"
)

type (
	// Fsys represents operation to filesystem.
	Fsys interface {
		Exist(string) bool
		Getwd() (string, error)
		Open(string) (io.ReadCloser, error)
		Write(string, []byte) error
	}

	// Logic represents application logic.
	Logic interface {
		Check(stdin io.Reader, cfgPath string) error
		IsIgnoredURL(uri string) bool
		CheckURLs(urls map[string]*strset.Set) error
		CheckURLWithMethod(ctx context.Context, u, method string) error
		CheckURL(ctx context.Context, u string) error
		ExtractURLsFromFiles(files *strset.Set) (map[string]*strset.Set, error)
		ExtractURLsFromFile(ctx context.Context, p string) (*strset.Set, error)
		GetFiles(stdin io.Reader) (*strset.Set, error)
	}

	// CfgReader reads and parses the configuration file.
	CfgReader interface {
		FindCfg() (string, error)
		ReadCfg(cfgPath string) (Cfg, error)
		InitCfg(cfg Cfg) (Cfg, error)
	}

	// HTTPClient abstracts *http.Client .
	HTTPClient interface {
		Do(req *http.Request) (*http.Response, error)
	}
)
