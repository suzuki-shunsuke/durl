package handler

import (
	"os"

	"github.com/urfave/cli"

	"github.com/suzuki-shunsuke/durl/internal/domain"
)

// Main calls a command.
func Main() {
	app := cli.NewApp()
	app.Name = "durl"
	app.Version = domain.Version
	app.Author = "suzuki-shunsuke https://github.com/suzuki-shunsuke"
	app.Usage = "check whether dead urls are included in files"

	app.Commands = []cli.Command{
		InitCommand,
		CheckCommand,
	}
	app.Run(os.Args)
}
