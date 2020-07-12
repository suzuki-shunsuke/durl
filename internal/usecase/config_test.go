package usecase_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suzuki-shunsuke/durl/internal/usecase"
)

func TestNewCfgReader(t *testing.T) {
	require.NotNil(t, usecase.NewCfgReader(nil))
}
