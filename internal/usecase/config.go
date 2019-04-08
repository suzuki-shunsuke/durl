package usecase

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/suzuki-shunsuke/go-cliutil"

	"github.com/suzuki-shunsuke/durl/internal/domain"
)

type (
	cfgReader struct {
		reader domain.CfgReader
		fsys   domain.Fsys
	}
)

func NewCfgReader(fsys domain.Fsys) domain.CfgReader {
	cr := &cfgReader{
		fsys: fsys,
	}
	cr.reader = cr
	return cr
}

func (reader *cfgReader) FindCfg() (string, error) {
	wd, err := reader.fsys.Getwd()
	if err != nil {
		return "", err
	}
	return cliutil.FindFile(wd, ".durl.yml", reader.fsys.Exist)
}

func (reader *cfgReader) ReadCfg(cfgPath string) (domain.Cfg, error) {
	cfg := domain.Cfg{
		HTTPMethod: "head,get",
	}
	if cfgPath == "" {
		d, err := reader.reader.FindCfg()
		if err != nil {
			return cfg, err
		}
		cfgPath = d
	}
	rc, err := reader.fsys.Open(cfgPath)
	if err != nil {
		return cfg, err
	}
	defer rc.Close()
	if err := yaml.NewDecoder(rc).Decode(&cfg); err != nil {
		return cfg, err
	}
	return reader.reader.InitCfg(cfg)
}

func (reader *cfgReader) InitCfg(cfg domain.Cfg) (domain.Cfg, error) {
	methods := map[string]struct{}{
		"get":      {},
		"head":     {},
		"head,get": {},
	}
	if _, ok := methods[cfg.HTTPMethod]; !ok {
		return cfg, fmt.Errorf(`invalid http_method_type: %s`, cfg.HTTPMethod)
	}
	return cfg, nil
}
