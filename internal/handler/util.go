package handler

import (
	"github.com/urfave/cli"
)

func wrapUsecase(err error) error {
	if err == nil {
		return nil
	}
	return cli.NewExitError(err.Error(), 1)
}
