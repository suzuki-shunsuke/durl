package handler

import (
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/suzuki-shunsuke/go-cliutil"
	"github.com/urfave/cli/v2"

	"github.com/suzuki-shunsuke/durl/internal/infra"
	"github.com/suzuki-shunsuke/durl/internal/usecase"
)

// checkCommand is the sub command "check".
var checkCommand = cli.Command{ //nolint:gochecknoglobals
	Name:   "check",
	Usage:  "check files",
	Action: check,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "config, c",
			Usage: "configuration file path",
			Value: "",
		},
	},
}

func check(c *cli.Context) error {
	cfgPath := c.String("config")
	fsys := infra.Fsys{}
	cfgReader := usecase.NewCfgReader(fsys)
	cfg, err := cfgReader.ReadCfg(cfgPath)
	if err != nil {
		return cliutil.ConvErrToExitError(err)
	}
	logic := usecase.NewLogic(
		cfg, fsys, &http.Client{
			Timeout: time.Duration(cfg.HTTPRequestTimeout) * time.Second,
		})
	if terminal.IsTerminal(0) {
		return cliutil.ConvErrToExitError(logic.Check(nil, cfgPath))
	}
	return cliutil.ConvErrToExitError(logic.Check(os.Stdin, cfgPath))
}
