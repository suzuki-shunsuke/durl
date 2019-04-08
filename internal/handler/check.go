package handler

import (
	"os"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/suzuki-shunsuke/go-cliutil"
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
	logic := usecase.NewLogic(infra.Fsys{})
	if terminal.IsTerminal(0) {
		return cliutil.ConvErrToExitError(logic.Check(nil, cfgPath))
	}
	return cliutil.ConvErrToExitError(logic.Check(os.Stdin, cfgPath))
}
