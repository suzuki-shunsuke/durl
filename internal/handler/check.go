package handler

import (
	"os"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/urfave/cli"

	"github.com/suzuki-shunsuke/durl/internal/infra"
	"github.com/suzuki-shunsuke/durl/internal/usecase"
)

// CheckCommand is the sub command "check".
var CheckCommand = cli.Command{
	Name:   "check",
	Usage:  "check files",
	Action: check,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "configuration file path",
			Value: "",
		},
	},
}

func check(c *cli.Context) error {
	cfgPath := c.String("config")
	if terminal.IsTerminal(0) {
		return wrapUsecase(
			usecase.Check(infra.Fsys{}, nil, cfgPath))
	}
	return wrapUsecase(
		usecase.Check(infra.Fsys{}, os.Stdin, cfgPath))
}
