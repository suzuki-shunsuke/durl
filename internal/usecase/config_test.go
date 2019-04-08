package usecase

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/suzuki-shunsuke/durl/internal/domain"
)

func Test_cfgReaderInitCfg(t *testing.T) {
	reader := &cfgReader{}
	cfg, err := reader.InitCfg(domain.Cfg{})
	require.Nil(t, err)
	require.Equal(t, "head,get", cfg.HTTPMethod)
	require.Equal(t, domain.DefaultTimeout, cfg.HTTPRequestTimeout)
	require.Equal(t, domain.DefaultMaxRequestCount, cfg.MaxRequestCount)
}
