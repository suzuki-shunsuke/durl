package handler

import (
	"github.com/urfave/cli"

	"github.com/suzuki-shunsuke/go-cliutil"

	"github.com/suzuki-shunsuke/durl/internal/infra"
	"github.com/suzuki-shunsuke/durl/internal/usecase"
)

// InitCommand is the sub command "init".
var InitCommand = cli.Command{
	Name:   "init",
	Usage:  "create a configuration file if it doesn't exist",
	Action: initCmd,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "dest, d",
			Usage: "created configuration file path",
			Value: ".durl.yml",
		},
	},
}

// initCmd is the sub command "init".
func initCmd(c *cli.Context) error {
	return cliutil.ConvErrToExitError(
		usecase.Init(infra.Fsys{}, c.String("dest")))
}
