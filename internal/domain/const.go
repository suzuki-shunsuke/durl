package domain

const (
	// DefaultTimeout is a default timeout of http request.
	DefaultTimeout = 10
	// DefaultMaxRequestCount is a default max parallel http request count.
	DefaultMaxRequestCount = 10
	// CfgTpl is a template of configuration file.
	CfgTpl = `
---
# configuration file of durl, which is a CLI tool to check whether dead urls are included in files.
# https://github.com/suzuki-shunsuke/durl
ignore_urls:
- https://github.com/suzuki-shunsuke/durl
ignore_hosts: []
http_method: head,get
max_request_count: 10
max_failed_request_count: 5
http_request_timeout: 10
`
)

var IgnoreHosts = []string{ //nolint:gochecknoglobals
	"localhost", "example.com", "example.org", "example.net", "127.0.0.1",
}
