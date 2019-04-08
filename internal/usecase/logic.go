package usecase

import (
	"github.com/suzuki-shunsuke/durl/internal/domain"
)

type (
	logic struct {
		logic domain.Logic
		cfg   domain.Cfg
		fsys  domain.Fsys
	}
)

// NewLogic returns a domain.Logic .
func NewLogic(cfg domain.Cfg, fsys domain.Fsys) domain.Logic {
	lgc := &logic{
		cfg:  cfg,
		fsys: fsys,
	}
	lgc.logic = lgc
	return lgc
}
