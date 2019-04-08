package usecase

import (
	"github.com/suzuki-shunsuke/durl/internal/domain"
)

type (
	logic struct {
		logic domain.Logic
		fsys  domain.Fsys
	}
)

func NewLogic(fsys domain.Fsys) domain.Logic {
	lgc := &logic{
		fsys: fsys,
	}
	lgc.logic = lgc
	return lgc
}
