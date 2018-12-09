package domain

type (
	// Cfg represents configuration.
	Cfg struct {
		IgnoreURLs []string `yaml:"ignore_urls"`
	}
)
