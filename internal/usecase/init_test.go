package usecase_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suzuki-shunsuke/gomic/gomic"

	"github.com/suzuki-shunsuke/durl/internal/domain"
	"github.com/suzuki-shunsuke/durl/internal/test"
	"github.com/suzuki-shunsuke/durl/internal/usecase"
)

func TestInit(t *testing.T) {
	data := []struct {
		title string
		isErr bool
		fsys  domain.Fsys
	}{{
		"file exist", false, test.NewFsys(t, gomic.DoNothing).
			SetReturnExist(true),
	}, {
		"succeed to write a file", false, test.NewFsys(t, gomic.DoNothing),
	}}
	for _, tt := range data {
		tt := tt
		t.Run(tt.title, func(t *testing.T) {
			err := usecase.Init(tt.fsys, ".durl.yml")
			if tt.isErr {
				require.NotNil(t, err)
				return
			}
			require.Nil(t, err)
		})
	}
}
