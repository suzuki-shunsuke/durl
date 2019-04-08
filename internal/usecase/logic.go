package usecase

import (
	"github.com/suzuki-shunsuke/durl/internal/domain"
)

type (
	logic struct {
		logic  domain.Logic
		cfg    domain.Cfg
		fsys   domain.Fsys
		client domain.HTTPClient
	}
)

// NewLogic returns a domain.Logic .
func NewLogic(cfg domain.Cfg, fsys domain.Fsys, client domain.HTTPClient) domain.Logic {
	lgc := &logic{
		cfg:    cfg,
		fsys:   fsys,
		client: client,
	}
	lgc.logic = lgc
	return lgc
}
