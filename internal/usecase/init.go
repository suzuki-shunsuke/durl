package usecase

import (
	"strings"

	"github.com/suzuki-shunsuke/durl/internal/domain"
)

// Init creates a configuration file if it doesn't exist.
func Init(fsys domain.Fsys, dst string) error {
	if fsys.Exist(dst) {
		return nil
	}
	return fsys.Write(dst, []byte(strings.Trim(domain.CfgTpl, "\n")))
}
