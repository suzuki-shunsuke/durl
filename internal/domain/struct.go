package domain

type (
	// Cfg represents configuration.
	Cfg struct {
		IgnoreURLs         []string `yaml:"ignore_urls"`
		IgnoreHosts        []string `yaml:"ignore_hosts"`
		HTTPMethod         string   `yaml:"http_method"`
		MaxRequestCount    int      `yaml:"max_request_count"`
		HTTPRequestTimeout int      `yaml:"http_request_timeout"`
	}
)
