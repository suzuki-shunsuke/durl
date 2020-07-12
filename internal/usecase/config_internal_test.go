package usecase

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suzuki-shunsuke/gomic/gomic"

	"github.com/suzuki-shunsuke/durl/internal/domain"
	"github.com/suzuki-shunsuke/durl/internal/test"
)

func TestNewCfgReader(t *testing.T) {
	require.NotNil(t, NewCfgReader(nil))
}

func Test_cfgReaderFindCfg(t *testing.T) {
	data := []struct {
		title string
		fsys  domain.Fsys
		isErr bool
		exp   string
	}{{
		"success", test.NewFsys(t, gomic.DoNothing).
			SetReturnExist(true),
		false, ".durl.yml",
	}, {
		"failed to get a current directory", test.NewFsys(t, gomic.DoNothing).
			SetReturnGetwd("", fmt.Errorf("failed to get a current directory")),
		true, "",
	}, {
		"failed to find a configuration file", test.NewFsys(t, gomic.DoNothing).
			SetReturnGetwd("/", nil),
		true, "",
	}}
	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			reader := &cfgReader{
				fsys: d.fsys,
			}
			p, err := reader.FindCfg()
			if d.isErr {
				require.NotNil(t, err)
				return
			}
			require.Nil(t, err)
			require.Equal(t, d.exp, p)
		})
	}
}

func Test_cfgReaderInitCfg(t *testing.T) {
	reader := &cfgReader{}
	cfg, err := reader.InitCfg(domain.Cfg{})
	require.Nil(t, err)
	require.Equal(t, "head,get", cfg.HTTPMethod)
	require.Equal(t, domain.DefaultTimeout, cfg.HTTPRequestTimeout)
	require.Equal(t, domain.DefaultMaxRequestCount, cfg.MaxRequestCount)
}

func Test_cfgReaderReadCfg(t *testing.T) {
	data := []struct {
		title   string
		isErr   bool
		cfgPath string
		mock    domain.CfgReader
		fsys    domain.Fsys
	}{{
		"failed to find a configuration file", true,
		"", test.NewCfgReader(t, nil).
			SetReturnFindCfg("", fmt.Errorf("failed to find a configuration file")), nil,
	}, {
		"failed to open a configuration file", true,
		"", test.NewCfgReader(t, nil).
			SetReturnFindCfg("/.durl.yml", nil),
		test.NewFsys(t, nil).SetReturnOpen(nil, fmt.Errorf("failed to open a configuration file")),
	}}
	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			reader := &cfgReader{
				reader: d.mock,
				fsys:   d.fsys,
			}
			_, err := reader.ReadCfg(d.cfgPath)
			if d.isErr {
				require.NotNil(t, err)
				return
			}
			require.Nil(t, err)
		})
	}
}
