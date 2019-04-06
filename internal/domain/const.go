package domain

import (
	"time"
)

const (
	// DefaultTimeout is a default timeout of http request.
	DefaultTimeout = 60 * time.Second
	// CfgTpl is a template of configuration file.
	CfgTpl = `
---
# configuration file of durl, which is a CLI tool to check whether dead urls are included in files.
# https://github.com/suzuki-shunsuke/durl
ignore_urls:
- https://github.com/suzuki-shunsuke/durl
`
)

var (
	IgnoreHosts = []string{"localhost", "example.com", "example.org", "example.net"}
)
